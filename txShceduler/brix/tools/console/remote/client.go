package remote

import (
	"encoding/json"
	"fmt"
	"strings"

	"txscheduler/brix/tools/console"
	"txscheduler/brix/tools/console/lout"
	"txscheduler/brix/tools/crypt/jaes"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/jnet/cnet"
)

var (
	serverAddress string
)

//LockCommands : prefix-lock
type LockCommands []string

//Do :
func (my LockCommands) Do(ps []string) bool {
	if len(ps) > 0 {
		text := ps[0]
		for _, v := range my {
			if v == text {
				return true
			}
		}
	}
	return false
}

//ClientRemoteConsole :
func ClientRemoteConsole(lockMsgs ...LockCommands) console.Ccmd {

	return console.Ccmd{
		Cmd:        dfREMOTE,
		SenderCmd:  true,
		HeaderFunc: func() string { return fmt.Sprintf("%v on / %v off", dfREMOTE, dfREMOTE) },
		Help:       remoteHelp(),
		Work: func(done chan<- bool, ps []string) {
			defer func() {
				console.Log("< remoteCall end >")
				done <- true
			}()

			if len(lockMsgs) > 0 {
				locks := lockMsgs[0]
				if locks.Do(ps) {
					dbg.RedBold("-------------------------------------")
					dbg.RedBold(" text contains prefix-lock Message")
					dbg.RedBold("-------------------------------------")
					return
				}
			}
			//dbg.Red("xxxxxxxxx --- ", ps)
			remoteCall(ps)
		},
	}
}

func remoteHelp() string {
	return `
{
	` + dfREMOTE + ` command....
}`
}

//TestCommand : test func
func TestCommand(msg string) {
	ss := strings.Split(msg, " ")
	remoteCall(ss)
}

func remoteCall(ps []string) {
	rQuery := ""
	for _, v := range ps {
		rQuery = fmt.Sprintf("%v%v ", rQuery, v)
	}
	queryCipher := AesKeyObject.EncryptStringString1(rQuery)
	fdata := chttp.JsonType{
		dfREMOTE: queryCipher,
	}
	if isAuth == true {
		fdata[authKey] = authValue
	}

	postSend(serverAddress, fdata, func(status int, buf []byte, err error) {
		if err != nil {
			console.Log(err)
			return
		}

		if status == chttp.StatusOK {
			result := chttp.JsonType{}
			if err := json.Unmarshal(buf, &result); err != nil {
				console.Log("Unmarshal : ", err)
				return
			}
			cipherMsg := result[dfREMOTE].(string)
			if text, err := AesKeyObject.DecryptStringString(cipherMsg); err != nil {
				console.Log("what the woo?", cipherMsg)
			} else {
				console.Log(text)

				lout.RemoteWrite(rQuery, text)
			}

		} else {
			console.Log("code:", status, ", errMessage", string(buf))
		}
	})
}

func replaceFileNameFilter(str string) string {
	flist := map[string]string{
		"\\": "",
		"/":  "",
		":":  ";",
		"*":  "",
		"?":  "",
		"\"": "",
		"<":  "",
		">":  "",
		"|":  "",
	}
	for k, v := range flist {
		str = strings.ReplaceAll(str, k, v)
	}
	str = strings.TrimSpace(str)
	return str
}

//ClientSetAddress :
func ClientSetAddress(address string) {
	serverAddress = address
}

//ClientGetAddress :
func ClientGetAddress() string {
	return serverAddress
}

//ClientInfo :
func ClientInfo() string {
	return `remote.Info:{
	address  : ` + dbg.Void(serverAddress) + `
	isAuth   : ` + dbg.Void(isAuth) + `
	auth_key : ` + dbg.Void(authKey) + `
	auth_val : ` + dbg.Void(authValue) + `
}`
}

//ClientSetAuth :
func ClientSetAuth(aeskey, key, value string) {
	isAuth = true
	AesKeyObject = jaes.New(aeskey)
	authKey = key
	authValue = value
}

//postSend :
func postSend(address string, fdata chttp.JsonType, callback func(status int, buf []byte, err error)) {

	client := cnet.New(address)
	client.SetTimeout(0)

	recvcode := 0
	recvbuf := []byte{}
	err := client.FORM(API, fdata, func(res cnet.Responser) {
		recvcode = res.StatusCode
		recvbuf = res.Bytes()
	})
	callback(recvcode, recvbuf, err)

	// fullURL := address + API
	// formData := func(d chttp.JsonType) url.Values {
	// 	u := url.Values{}
	// 	if d != nil {
	// 		for k, v := range d {
	// 			u.Set(k, fmt.Sprintf("%v", v))
	// 		}
	// 	}
	// 	return u
	// }

	// res, err := http.PostForm(fullURL, formData(fdata))
	// if err != nil {
	// 	fmt.Println("remote.post_err :", err)
	// 	code := 0
	// 	if res != nil {
	// 		code = res.StatusCode
	// 	}
	// 	callback(code, nil, err)
	// 	return
	// }
	// defer res.Body.Close()

	// buf, err := ioutil.ReadAll(res.Body)
	// callback(res.StatusCode, buf, err)
}
