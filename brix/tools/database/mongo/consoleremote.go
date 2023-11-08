package mongo

import "fmt"

var (
	isReceiver    = false
	remoteMessage string

	callbackFunc func(a string)         = nil
	deliverFunc  func(a ...interface{}) = nil

	isStartRemote    = false
	isViewJSONFormat = false
)

//ConsoleViewJSONFormat :
func ConsoleViewJSONFormat() {
	isViewJSONFormat = true
}

//ConsoleRemoteRegister :
func ConsoleRemoteRegister(callback func(msg string)) {
	isReceiver = true

	callbackFunc = callback
	remoteMessage = ""
}

//ConsoleRemoteDeliver :
func ConsoleRemoteDeliver(f func(a ...interface{})) {
	isReceiver = true

	deliverFunc = f
	remoteMessage = ""
}

func startRemote() {
	isStartRemote = true
}

func isRemoteMode() bool {
	return isStartRemote == true && isReceiver == true
}

func remoteMsg(a ...interface{}) {
	if isReceiver == false {
		return
	}
	for _, v := range a {
		remoteMessage = fmt.Sprintf("%v%v ", remoteMessage, v)
	} //for
	remoteMessage = fmt.Sprintf("%v\n", remoteMessage)
}

func remoteCallback() {
	isStartRemote = false

	if isReceiver == false {
		return
	}

	if callbackFunc != nil {
		callbackFunc(remoteMessage)
	}

	if deliverFunc != nil {
		deliverFunc(remoteMessage)
	}

	remoteMessage = ""
}
