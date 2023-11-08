package remote

import (
	"txscheduler/brix/tools/console"
	"txscheduler/brix/tools/crypt/jaes"
)

const (
	//DefaultKey :
	DefaultKey = "outsoaaaa123!!@"
	//DefaultValue :
	DefaultValue = "zenebito_keyxxy!@#$"
	//DefaultAESKey :
	DefaultAESKey = "AXAeS!@#$1234keyDefAULT//"
)

const (
	//API :URL :
	API = "/consoleremote"

	//dfREMOTE : "remote"
	dfREMOTE = console.KEYWORD_REMOTE
)

var (
	isAuth        bool
	authKey       string
	authValue     string
	preHandleFunc func()
)

//AesKeyObject :
var AesKeyObject jaes.Openssl
