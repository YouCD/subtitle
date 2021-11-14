package xunlei

import (
	"fmt"
	"testing"
)

var x = Xunlei{"http://sub.xmp.sandai.net:8000/subxl/%s.json"}

func TestXunlei_GetSubtitleInfo(t *testing.T) {
	list, _ := x.GetSubtitleInfo("/tmp/Forever.2014.S01E01.Pilot.1080p.WEB-DL.DD5.1.H.264-ECI.mkv")
	for _, v := range list {
		fmt.Println(v)
	}
}
