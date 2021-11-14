package shooter

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/YouCD/subtitle/pkg/common"
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type File struct {
	Ext  string `json:"Ext"`
	Link string `json:"Link"`
}

type ShooterSub struct {
	Desc  string `json:"Desc"`
	Delay int    `json:"Delay"`
	Files []File `json:"Files"`
}

var (
	InvalidCharacter = errors.New("invalid character 'ÿ' looking for beginning of value")
)

type Shooter struct {
	API string
}

func GetShooter() *Shooter {
	return &Shooter{API: "https://www.shooter.cn/api/subapi.php?"}
}

func (m *Shooter) GetSubtitleInfo(path string) (SubtitleInfoList []common.SubtitleInfo, err error) {

	hash, err := m.CalculateHash(path)
	if err != nil {
		return nil, err
	}

	dir, file := filepath.Split(path)
	urlTemp := fmt.Sprintf("filehash=%s&pathinfo=%s&format=%s&lang=%s", hash, file, "json", "Chn")

	a := url.PathEscape(urlTemp)

	req, err := http.NewRequest("POST", m.API+a, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "http.NewRequest")
	}
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)

	if err != nil {
		return nil, errors.WithMessage(err, "http.NewRequest")
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "ioutil.ReadAll")
	}

	var shooterSubList []ShooterSub
	err = json.Unmarshal(b, &shooterSubList)
	if err != nil {
		red := color.FgRed.Render
		fmt.Printf("Not found subtitle for video file %s may be incomplete.\n", red(file))
		return nil, InvalidCharacter
	}
	for index, v := range shooterSubList {
		for _, f := range v.Files {
			var item common.SubtitleInfo
			item.Url = f.Link

			filename := common.FileNamePrefix(path)
			SubtitlePath := filepath.Join(dir, filename)
			if index < 10 {
				if index == 0 {
					item.SubtitleName = SubtitlePath + "." + f.Ext
					SubtitleInfoList = append(SubtitleInfoList, item)
					continue
				}
				item.SubtitleName = SubtitlePath + ".0" + strconv.Itoa(index) + "." + f.Ext
			}
			SubtitleInfoList = append(SubtitleInfoList, item)
		}
	}

	return
}

func (m *Shooter) DownloadSubtitle(info common.SubtitleInfo) error {
	return common.DownLoadFile(info)
}

func (m *Shooter) CalculateHash(path string) (hash string, err error) {
	if !common.Exists(path) {
		return "", errors.New(fmt.Sprintf("%s not exists.", path))
	}
	if common.IsDir(path) {
		return "", errors.New("not video file.")
	}

	var hashlist []string

	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return "", errors.WithMessage(err, "openfile")
	}
	offset := 4096
	md5str1, err := md5Sum(offset, 0, f)
	if err != nil {
		return "", errors.WithMessage(err, "md5Sum")
	}

	hashlist = append(hashlist, md5str1)

	fi, _ := f.Stat()
	md5str2, err := md5Sum(int(fi.Size()/3*2), 0, f)
	if err != nil {
		return "", errors.WithMessage(err, "md5Sum")
	}

	hashlist = append(hashlist, md5str2)

	md5str3, err := md5Sum(int(fi.Size()/3), 0, f)
	if err != nil {
		return "", errors.WithMessage(err, "md5Sum")
	}
	hashlist = append(hashlist, md5str3)
	md5str4, err := md5Sum(int(offset*-2), 2, f)
	if err != nil {
		return "", errors.WithMessage(err, "md5Sum")
	}

	hashlist = append(hashlist, md5str4)
	hash = strings.Join(hashlist, ";")

	return
}

func md5Sum(position, whence int, file *os.File) (md5Str string, err error) {
	_, err = file.Seek(int64(position), whence)
	if err != nil {
		return "", errors.WithMessage(err, "file.Seek")
	}
	b2 := make([]byte, 4096)
	_, err = file.Read(b2)
	if err != nil {
		return "", errors.WithMessage(err, "file.Read")
	}
	has := md5.Sum(b2)
	md5str1 := fmt.Sprintf("%x", has)
	return md5str1, nil
}
