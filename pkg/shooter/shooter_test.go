package shooter

import "testing"

var (
	s = Shooter{"https://www.shooter.cn/api/subapi.php?"}
)

func TestShooter_GetSubtitleInfo(t *testing.T) {
	//url := "https://www.shooter.cn/api/subapi.php?lang=Chn&pathinfo=%2Fhome%2Fycd%2Faria2_download%2Faria2-downloads%2FForever.US.S01%2FForever.2014.S01E01.Pilot.1080p.WEB-DL.DD5.1.H.264-ECI.mkv&filehash=09d179e119dab0e053b0f99315e9c34b%3B0c8899a3ac10292fa1c5fcd32af1a80d%3Ba5028840052e9430e527a3998d462277%3B62424897269e3b1be315f5ff4d28f1ab&format=json"
	s.GetSubtitleInfo("/tmp/Forever.2014.S01E01.Pilot.1080p.WEB-DL.DD5.1.H.264-ECI.mkv")
}
