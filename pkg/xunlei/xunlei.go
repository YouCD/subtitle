package xunlei

import (
	"subtitle/pkg/common"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type XunleiSub struct {
	Scid     string `json:"scid"`
	Sname    string `json:"sname"`
	Language string `json:"language"`
	Rate     string `json:"rate"`
	Surl     string `json:"surl"`
	Svote    int    `json:"svote"`
	Roffset  int    `json:"roffset"`
}
var (
	NotFoundErr=errors.New("NotFound")
)
type XunleiSublist struct {
	Sublist []XunleiSub `json:"sublist"`
}
type Xunlei struct {
	API string
}

func GetXunlei() *Xunlei {
	return &Xunlei{API: "http://sub.xmp.sandai.net:8000/subxl/%s.json"}
}
func (m *Xunlei) GetSubtitleInfo(path string) (SubtitleInfoList []common.SubtitleInfo, err error) {

	cid, err := m.CalculateHash(path)
	if err != nil {
		return nil, errors.WithMessage(err, "Xunlei.CalculateHash")
	}

	if len(cid) == 0 {
		return nil, errors.New("cid is 0")
	}
	dir, fileName := filepath.Split(path)
	resp, err := http.Get(fmt.Sprintf(m.API, cid))
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, errors.New("resp is nil")
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var xunleiSublist XunleiSublist
	err = json.Unmarshal(b, &xunleiSublist)
	if err != nil {
		return nil, errors.WithMessage(err, "Unmarshal subList ")
	}
	if len(xunleiSublist.Sublist)==1&&xunleiSublist.Sublist[0].Sname==""{
		red := color.FgRed.Render
		fmt.Printf("Not found subtitle for video file %s may be incomplete.\n", red(fileName))
		return nil, NotFoundErr
	}

	subList := xunleiSublist.Sublist

	for index, v := range subList {
		var item common.SubtitleInfo
		filename := common.FileNamePrefix(path)
		SubtitlePath := filepath.Join(dir, filename)
		suffix := common.FileNameSuffix(v.Surl)
		if len(v.Scid) > 0 && v.Scid != "" {
			item.Url = v.Surl
			if strings.Contains(v.Language, "简体&英语") || strings.Contains(v.Language, "简体") || strings.Contains(v.Language, "未知语言") || strings.Contains(v.Language, "英语") {
				if index == 0 {
					item.SubtitleName = SubtitlePath + "." + v.Language + suffix
					SubtitleInfoList = append(SubtitleInfoList, item)
					continue
				}
				item.SubtitleName = SubtitlePath + ".0" + strconv.Itoa(index) + "." + v.Language + suffix
				SubtitleInfoList = append(SubtitleInfoList, item)
				continue
			}
			item.SubtitleName = SubtitlePath + ".0" + strconv.Itoa(index) + "." + v.Language + suffix
			SubtitleInfoList = append(SubtitleInfoList, item)
		}
	}
	return
}

func (m *Xunlei) DownloadSubtitle(info common.SubtitleInfo) error {
	return common.DownLoadFile(info)
}

func (m *Xunlei) CalculateHash(path string) (hash string, err error) {
	if !common.Exists(path) {
		return "", errors.New("not exists.")
	}
	if common.IsDir(path) {
		return "", errors.New("not video file.")
	}

	sha1Ctx := sha1.New()

	fp, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fp.Close()
	stat, err := fp.Stat()
	if err != nil {
		return "", err
	}
	fileLength := stat.Size()
	if fileLength < 0xF000 {
		return "", err
	}
	bufferSize := int64(0x5000)
	positions := []int64{0, int64(math.Floor(float64(fileLength) / 3)), fileLength - bufferSize}
	for _, position := range positions {
		var buffer = make([]byte, bufferSize)
		_, err = fp.Seek(position, 0)
		if err != nil {
			return "", err
		}
		_, err = fp.Read(buffer)
		if err != nil {
			return "", err
		}
		sha1Ctx.Write(buffer)
	}

	hash = fmt.Sprintf("%X", sha1Ctx.Sum(nil))
	return hash, nil

}
