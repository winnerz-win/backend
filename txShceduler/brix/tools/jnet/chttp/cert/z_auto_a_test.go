package cert

import (
	"testing"

	"txscheduler/brix/tools/jpath"
)

func TestMake(t_ *testing.T) {
	//Make("127.0.0.1", "ssl", 1)
	//Make("http://npt.iptime.org:37003", "ssl")
	Make("192.168.0.19", "ssl", 100)

	copyTarget(jpath.NowPath(), "ssl", `D:\work\go\src\github.com\outsourcing\Zenebito\cmd\a_server_infra`)
}

func TestXXX(t_ *testing.T) {
	Make("http://npt.iptime.org:37003", "ssl", 1)
}

func TestCopy(t_ *testing.T) {
	copyTarget(jpath.NowPath(), "ssl", `D:\work\go\src\github.com\outsourcing\Zenebito\cmd\a_server_infra`)
}
