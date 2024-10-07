package cert

import (
	"testing"

	"txscheduler/brix/tools/jpath"
)

func TestDefaultOption(t_ *testing.T) {
	MakeOption(nil)
}

// !smm.ssl.make !ssl.make !makessl !ssl !https
// !ssl.test
func TestGMMSSLMake(t_ *testing.T) {
	opt := &Option{
		Country:            "ko",
		Organization:       "fmdot",
		OrganizationalUnit: "zenebito",
		CommonName:         "gmm_for_web_server",
		SerialNumberString: "9115883790349911111126501228181555555821869333691004588797900000455571312345678915681511554545646477",
		StartYMD:           20200101,
		EndYMD:             21200101,
		IP: []string{
			"127.0.0.1",
			"192.168.0.19",   //local-pc
			"192.168.0.52",   //local-test
			"121.140.201.65", //local-public
			"8.210.228.21",   //real-server public,
			"172.31.92.82",   //real-server private,
		},
		DNSNames: []string{
			"aga.gmm.gold",
			"aga.gmm.gold:37001",
		},
		FileName: "ssl",
		RootPath: jpath.NowPath() + "\\gmm_ssl",
	}
	//DNSNames       []string
	MakeOption(opt)

	copyTarget(opt.RootPath, "ssl", `D:\work\go\src\github.com\outsourcing\Zenebito\cmd\a_server_infra`)
}
