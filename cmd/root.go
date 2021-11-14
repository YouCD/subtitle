package cmd

import (
	"fmt"
	"github.com/YouCD/subtitle/pkg/common"
	"github.com/YouCD/subtitle/pkg/shooter"
	"github.com/YouCD/subtitle/pkg/xunlei"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	Name      = "subtitle"
	VideoFile string
	VideoDir  string
	Source    string
	x         = xunlei.GetXunlei()
	s         = shooter.GetShooter()
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.Flags().StringVarP(&VideoFile, "file", "f", "", "指定电影文件,如果不指定则使用当前目录")
	rootCmd.Flags().StringVarP(&VideoDir, "dir", "d", "", "指定电影存放目录")
	rootCmd.Flags().StringVar(&Source, "source", "shooter", "指定字幕源,如果不指定则默认为 shooter")
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   Name,
	Short: fmt.Sprintf("%s 是用于获取字幕的工具", Name),
	Example: `	subtitle -f SomeVideo.mkv
	subtitle -d SomeVideDir
	subtitle --source xunlei -d SomeVideDir
    subtitle `,
	Run: func(cmd *cobra.Command, args []string) {
		if Source != "shooter" && Source != "xunlei" {
			fmt.Println("not support the source.")
			return
		}
		if Source == "shooter" {
			doExec(s)
		}
		if Source == "xunlei" {
			doExec(x)
		}

	},
}

func execDir(subtitle common.Subtitle, path string) {
	var tempPath string
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if path == "" {
		tempPath = dir
	} else {
		tempPath = path
	}

	if !common.IsDir(tempPath) {
		return
	}

	files, err := ioutil.ReadDir(tempPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	videoList := common.GetVideoList(files, tempPath)
	for _, f := range videoList {
		list, err := subtitle.GetSubtitleInfo(filepath.Join(tempPath, f))
		if err != nil {
			switch {
			case err == shooter.InvalidCharacter:
				continue
			case err == xunlei.NotFoundErr:
				continue
			default:
				fmt.Println(err)
				continue
			}
		}

		for _, item := range list {
			err = subtitle.DownloadSubtitle(item)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func execFile(subtitle common.Subtitle, filePath string) {
	list, err := subtitle.GetSubtitleInfo(filePath)
	if err != nil {
		switch {
		case err == shooter.InvalidCharacter:
			return
		case err == xunlei.NotFoundErr:
			return
		default:
			fmt.Println(err)
			return
		}
	}
	if len(list) == 0 {
		return
	}
	for _, v := range list {
		err := subtitle.DownloadSubtitle(v)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func doExec(subtitle common.Subtitle) {
	// 指定单个文件
	if VideoDir == "" && VideoFile != "" {
		execFile(subtitle, VideoFile)
		return
	}
	// 没有指定video目录时候
	if VideoDir == "" && VideoFile == "" {
		execDir(subtitle, "")
		return
	}
	// 指定文件夹
	if VideoDir != "" && VideoFile == "" {
		if !common.IsDir(VideoDir) {
			fmt.Printf("%s is  dir.\n", VideoDir)
			return
		}
		execDir(subtitle, VideoDir)

		return
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
