package common

type SubtitleInfo struct {
	Url          string
	SubtitleName string
}

type Subtitle interface {
	CalculateHash(path string) (hash string, err error)
	GetSubtitleInfo(path string) (SubtitleInfoList []SubtitleInfo, err error)
	DownloadSubtitle(info SubtitleInfo) error
}
