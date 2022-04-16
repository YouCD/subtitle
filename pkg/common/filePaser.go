package common

import (
	"crypto/tls"
	"fmt"
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	return !IsDir(path)
}

func FileNamePrefix(filename string) string {

	fileNameAll := path.Base(filename)
	fileSuffix := path.Ext(filename)
	filePrefix := fileNameAll[0 : len(fileNameAll)-len(fileSuffix)]

	return filePrefix
}

func FileNameSuffix(filename string) string {
	fileSuffix := path.Ext(filename)
	return fileSuffix
}

func DownLoadFile(info SubtitleInfo) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Transport: tr,
	}
	resp, err := httpClient.Get(info.Url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(info.SubtitleName, body, 0644)
	if err != nil {
		return errors.WithMessage(err, "ioutil.WriteFile")
	}

	green := color.FgGreen.Render
	fmt.Printf("The subtitle file %s is saved.\n", green(info.SubtitleName))
	return nil
}

func IsVideo(filename string) bool {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Get the content
	contentType, err := GetFileContentType(f)
	if err != nil {
		return false
	}
	if strings.Contains(contentType,"video")   {
		return true
	}
	return false
}

func GetFileContentType(out *os.File) (string, error) {

	// 只需要前 512 个字节就可以了
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func GetVideoList(files []fs.FileInfo,path string) (fileList []string) {
	for _, f1 := range files {
		if f1.IsDir() {
			continue
		}
		filePath := filepath.Join(path, f1.Name())
		var cmd string
		args := strings.Split(os.Args[0], "./")
		if len(args) > 1 {
			cmd = strings.Split(os.Args[0], "./")[1]
		} else {
			cmd = strings.Split(os.Args[0], "./")[0]
		}

		if os.Args[0] == f1.Name() || strings.Contains(filePath, cmd) {
			continue
		}
		f := filepath.Join(path, f1.Name())
		if !IsVideo(f) {
			continue
		}
		fileList = append(fileList, f1.Name())
	}


	return
}
