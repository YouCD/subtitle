package opensubtitles

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/YouCD/subtitle/pkg/common"
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	ChunkSize = 65536
	ApiKey    = "HaSKP2QrF89J5xooZPU6HcZUgPfrDpFw"
)

var (
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
)

type attributes struct {
	//SubtitleID        string    `json:"subtitle_id"`
	//Language          string    `json:"language"`
	//DownloadCount     int       `json:"download_count"`
	//NewDownloadCount  int       `json:"new_download_count"`
	//HearingImpaired   bool      `json:"hearing_impaired"`
	//Hd                bool      `json:"hd"`
	//Fps               float64   `json:"fps"`
	//Votes             int       `json:"votes"`
	//Ratings           float64   `json:"ratings"`
	//FromTrusted       bool      `json:"from_trusted"`
	//ForeignPartsOnly  bool      `json:"foreign_parts_only"`
	//UploadDate        time.Time `json:"upload_date"`
	//AiTranslated      bool      `json:"ai_translated"`
	//MachineTranslated bool      `json:"machine_translated"`
	//Release           string    `json:"release"`
	//Comments          string    `json:"comments"`
	//LegacySubtitleID  int       `json:"legacy_subtitle_id"`
	//Uploader          struct {
	//	UploaderID int    `json:"uploader_id"`
	//	Name       string `json:"name"`
	//	Rank       string `json:"rank"`
	//} `json:"uploader"`
	//FeatureDetails struct {
	//	FeatureID       int    `json:"feature_id"`
	//	FeatureType     string `json:"feature_type"`
	//	Year            int    `json:"year"`
	//	Title           string `json:"title"`
	//	MovieName       string `json:"movie_name"`
	//	ImdbID          int    `json:"imdb_id"`
	//	TmdbID          int    `json:"tmdb_id"`
	//	SeasonNumber    int    `json:"season_number"`
	//	EpisodeNumber   int    `json:"episode_number"`
	//	ParentImdbID    int    `json:"parent_imdb_id"`
	//	ParentTitle     string `json:"parent_title"`
	//	ParentTmdbID    int    `json:"parent_tmdb_id"`
	//	ParentFeatureID int    `json:"parent_feature_id"`
	//} `json:"feature_details"`
	//URL          string `json:"url"`
	//RelatedLinks []struct {
	//	Label  string `json:"label"`
	//	URL    string `json:"url"`
	//	ImgURL string `json:"img_url,omitempty"`
	//} `json:"related_links"`
	Files []struct {
		FileID   int    `json:"file_id"`
		CdNumber int    `json:"cd_number"`
		FileName string `json:"file_name"`
	} `json:"files"`
	MoviehashMatch bool `json:"moviehash_match"`
}

type data struct {
	//ID         string     `json:"id"`
	//Type       string     `json:"type"`
	Attributes attributes `json:"attributes"`
}
type openSubtitlesData struct {
	TotalPages int `json:"total_pages"`
	TotalCount int `json:"total_count"`
	//PerPage    int    `json:"per_page"`
	Page int    `json:"page"`
	Data []data `json:"data"`
}

type openSubtitles struct {
	searchApi   string
	downloadApi string
}

func NewOpenSubtitles() *openSubtitles {
	return &openSubtitles{
		searchApi: "https://api.opensubtitles.com/api/v1/subtitles?languages=chi,zht,zhe,ace&moviehash=",
		//searchApi:   "https://api.opensubtitles.com/api/v1/subtitles?moviehash=",
		downloadApi: "https://api.opensubtitles.com/api/v1/download",
	}
}

func (o *openSubtitles) CalculateHash(path string) (hash string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hashUint64, err := hashFile(file)
	if err != nil {
		return "", err
	}
	hash = fmt.Sprintf("%x", hashUint64)
	return
}
func (o *openSubtitles) GetSubtitleInfo(path string) (SubtitleInfoList []*common.SubtitleInfo, err error) {
	hash, err := o.CalculateHash(path)
	if err != nil {
		return nil, err
	}
	dir, file := filepath.Split(path)

	req, err := http.NewRequest(http.MethodGet, o.searchApi+hash, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "http.NewRequest")
	}
	req.Header.Add("Api-Key", ApiKey)
	req.Header.Add("Content-Type", "application/json")

	httpClient := &http.Client{
		Transport: tr,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, "http.NewRequest")
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "ioutil.ReadAll")
	}

	var openSubtitlesList openSubtitlesData
	err = json.Unmarshal(b, &openSubtitlesList)
	if err != nil {
		red := color.FgRed.Render
		fmt.Printf("Not found subtitle for video file %s may be incomplete.\n", red(file))
		return nil, common.InvalidCharacter
	}
	if len(openSubtitlesList.Data) == 0 {
		red := color.FgRed.Render
		fmt.Printf("Not found subtitle for video file %s may be incomplete.\n", red(file))
		return nil, common.InvalidCharacter
	}
	for _, v := range openSubtitlesList.Data {
		for _, f := range v.Attributes.Files {
			osbDownloadInfo, err := o.getDownloadUlr(f.FileID)
			if err != nil {
				red := color.FgRed.Render
				fmt.Printf("Not found subtitle for video file %s may be incomplete.\n", red(file))
				return nil, common.InvalidCharacter
			}

			SubtitlePath := filepath.Join(dir, osbDownloadInfo.FileName)

			var sub = common.SubtitleInfo{
				Url:          osbDownloadInfo.Link,
				SubtitleName: SubtitlePath,
			}
			SubtitleInfoList = append(SubtitleInfoList, &sub)
		}
	}

	return
}
func (o *openSubtitles) DownloadSubtitle(info *common.SubtitleInfo) error {

	return common.DownLoadFile(info)

}

func hashFromBuffer(buf []byte, fileSize uint64) (hash uint64, err error) {
	// Convert to uint64, and sum.
	var nums [(ChunkSize * 2) / 8]uint64
	reader := bytes.NewReader(buf)
	err = binary.Read(reader, binary.LittleEndian, &nums)
	if err != nil {
		return 0, err
	}
	for _, num := range nums {
		hash += num
	}

	return hash + fileSize, nil
}

// Read a chunk of a file at `offset` so as to fill `buf`.
func readChunk(file *os.File, offset int64, buf []byte) (err error) {
	n, err := file.ReadAt(buf, offset)
	if err != nil {
		return
	}
	if n != ChunkSize {
		return fmt.Errorf("Invalid read %v", n)
	}
	return
}

// hashFile generates an OSDB hash for an *os.File.
func hashFile(file *os.File) (hash uint64, err error) {
	fi, err := file.Stat()
	if err != nil {
		return
	}
	if fi.Size() < ChunkSize {
		return 0, fmt.Errorf("File is too small")
	}

	// Read head and tail blocks.
	buf := make([]byte, ChunkSize*2)
	err = readChunk(file, 0, buf[:ChunkSize])
	if err != nil {
		return
	}
	err = readChunk(file, fi.Size()-ChunkSize, buf[ChunkSize:])
	if err != nil {
		return
	}

	return hashFromBuffer(buf, uint64(fi.Size()))
}

type downloadInfo struct {
	Link         string    `json:"link"`
	FileName     string    `json:"file_name"`
	Requests     int       `json:"requests"`
	Remaining    int       `json:"remaining"`
	Message      string    `json:"message"`
	ResetTime    string    `json:"reset_time"`
	ResetTimeUtc time.Time `json:"reset_time_utc"`
}

func (o *openSubtitles) getDownloadUlr(fileID int) (osbDownloadInfo *downloadInfo, err error) {
	osbDownloadInfo = new(downloadInfo)
	reader := strings.NewReader(fmt.Sprintf("{\n  \"file_id\": %d\n}", fileID))
	req, err := http.NewRequest("POST", o.downloadApi, reader)
	if err != nil {
		fmt.Println(err)
		return nil, errors.WithMessage(err, "http.NewRequest")
	}

	req.Header.Add("Api-Key", ApiKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")

	httpClient := &http.Client{
		Transport: tr,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, "http.NewRequest")
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "ioutil.ReadAll")
	}
	if err = json.Unmarshal(bytes, &osbDownloadInfo); err != nil {
		return nil, errors.WithMessage(err, "Unmarshal")
	}

	return
}
