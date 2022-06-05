package opensubtitles

import (
	"fmt"
	"testing"
)

func Test_openSubtitles_getDownloadUlr(t *testing.T) {
	subtitleInfo, err := NewOpenSubtitles().getDownloadUlr(7251340)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(subtitleInfo)
}
