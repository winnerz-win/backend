package remote

import (
	"fmt"
	"net/http"
	"strings"

	"txscheduler/brix/tools/console"
	"txscheduler/brix/tools/crypt/jaes"
	"txscheduler/brix/tools/jnet/chttp"
)

//serverSetAuth :
func serverSetAuth(aeskey, key, value string) {
	isAuth = true
	AesKeyObject = jaes.New(aeskey)
	authKey = key
	authValue = value
}

//ServerPreHandleFunc : Call function before remote.RouterHandle!
func ServerPreHandleFunc(f func()) {
	preHandleFunc = f
}

//ServerContext :
func ServerContext() chttp.PContext {
	return &chttp.Context{
		chttp.POST, API, postHandle(),
	}
}

//postHandle :
func postHandle() chttp.RouterHandle {

	if preHandleFunc != nil {
		preHandleFunc()
	}
	msgReceiverC := console.RemoteRegister()

	return func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
		chttp.BindingDummy(req)

		if isAuth == true {
			cVal := req.Form.Get(authKey)
			if authValue != cVal {
				chttp.ResultJSON(w, chttp.StatusOK, chttp.JsonType{dfREMOTE: "remote not allow."})
				return
			}
		}

		cipherMsg := req.Form.Get(dfREMOTE)
		remoteMsg := ""
		if text, err := AesKeyObject.DecryptStringString(cipherMsg); err != nil {
			chttp.ResultJSON(w, chttp.StatusBadRequest, "fuck the woo.")
			return
		} else {
			remoteMsg = text
		}

		remoteMsg = strings.Trim(remoteMsg, " ")
		if remoteMsg == "" {
			chttp.ResultJSON(w, chttp.StatusOK, chttp.JsonType{dfREMOTE: "empty"})
			return
		}

		console.TestCommand(remoteMsg, true)
		resultMessage := <-msgReceiverC

		cipherString := AesKeyObject.EncryptStringString1(resultMessage)

		chttp.ResultJSON(w, chttp.StatusOK, chttp.JsonType{dfREMOTE: cipherString})
	}
}

//ServerHandleContexts : ServerhandlerFunc - mongo.ConsoleRemoteDeliver(console.Log)
func ServerHandleContexts(aeskey, authKey, authValue string, serverhandlerFunc func()) []chttp.PContext {
	fmt.Println("▷ REGISTER : remote.ServerHandleContexts()")

	serverSetAuth(aeskey, authKey, authValue)

	//handlerFunc : database.ConsoleRemoteDeliver(console.Log)
	ServerPreHandleFunc(serverhandlerFunc)

	return []chttp.PContext{
		ServerContext(),
	}
}

/*
rshc := remote.ServerHandleContexts(remote.DefaultAESKey, remote.DefaultKey, remote.DefaultValue, dbconsole.RemoteDeliver)
classic.SetContextHandles(rshc)
dbconsole.SetConsole(model.DB(), "db")
*/
