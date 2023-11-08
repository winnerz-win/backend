package console

import (
	"fmt"
	"sync"
)

var (
	isReceiver    = false
	remoteMessage string
	callbackFunc  func(msg string)

	isWriteMode = false
	muCommand   sync.Mutex
)

/*
RemoteRegister : 원격 console 명령을 처리 할수 있게 한다.

RemoteRegister를 호출 한 후에
console.TestCommand(msg , true) 를 호출 하면
msgReceiverC 채널로 로그를 수신할 수 있다.
*/
func RemoteRegister() <-chan string {

	isReceiver = true

	msgReceiverC := make(chan string, 1)
	callbackFunc = func(msg string) {
		msgReceiverC <- msg
	}

	remoteMessage = ""

	return msgReceiverC
}

func startWriteMode() {
	isWriteMode = true
	// if isRemoteWriteMode() {
	// 	Log("----------- REMOTE_LOG_START -----------")
	// }
}

func isRemoteWriteMode() bool {
	return isReceiver == true && isWriteMode == true
}

func remoteMsg(a ...interface{}) {
	if isReceiver == false || isWriteMode == false {
		return
	}

	for _, v := range a {
		remoteMessage = fmt.Sprintf("%v%v ", remoteMessage, v)
	} //for
	remoteMessage = fmt.Sprintf("%v\n", remoteMessage)
}

func removeCallback() {
	if isWriteMode == false {
		return
	}

	if isReceiver == true && callbackFunc != nil {
		//Log("----------- REMOTE_LOG_END -----------")

		callbackFunc(remoteMessage)
		remoteMessage = ""
	}

	isWriteMode = false
}
