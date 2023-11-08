package console

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"txscheduler/brix/tools/console/lout"
	"txscheduler/brix/tools/dbg/cc"
)

const (
	LINE_LO = "--------------------------------------------------------------"
	LINE_UP = "=============================================================="
)

var (
	cmds     = make([]ICmd, 0)
	helpTaps = []string{}

	isServerStart       bool
	isRemoteClientStart bool
	remoteModeOn        bool
	remoteConsoleTag    = ""

	readyMessage   = "server_ready"
	runtimeMessage = ""

	isUserHelp bool
	userHelp   string
)

//DoneC :
func DoneC(done chan<- bool, msg ...string) {
	if len(msg) > 0 {
		Log(msg[0])
	} else {
		Log("< end >")
	}
	done <- true
}

//IsRemoteClientStart :
func IsRemoteClientStart() bool {
	return isRemoteClientStart
}

// ServerStart :
func ServerStart() {
	isServerStart = true
}

//SetReadyMessage : isServerStart == false
func SetReadyMessage(msg string) {
	readyMessage = msg
}

//SetRuntimeMessage : isServerStart == true
func SetRuntimeMessage(msg string) {
	runtimeMessage = msg
}

//SetUserHelp :
func SetUserHelp(msg string) {
	userHelp = msg
	isUserHelp = true
}

//SetCmd :
func SetCmd(cs Commands) {
	if cs == nil {
		return
	}
	for _, c := range cs {
		if checkHelp(c.ICommand()) == false {
			if c.Work != nil {
				cmds = append(cmds, c)
			}
		}
	} //for

}

//SetICmd :
func SetICmd(ics ...Ccmd) {
	for _, ic := range ics {
		if checkHelp(ic.ICommand()) == false {
			if ic.Work != nil {
				cmds = append(cmds, ic)
			}
		}
	} //for
}

//ClearCmd :
func ClearCmd() {
	cmds = make([]ICmd, 0)
}

//StartRemote : remote client func.
func StartRemote(tag ...string) {
	start(true, false)
	go remoteRun(tag...)
}

//Start :
func Start() {
	start(false, false)
}

//ReadMode :
func ReadMode() {
	start(false, true)
}

var isRun = false
var runMu sync.RWMutex

//start : remoteMode bool, readMode bool
func start(remoteMode bool, readMode bool) {
	defer runMu.Unlock()
	runMu.Lock()
	if isRun {
		return
	}
	isRun = true

	isRemoteClientStart = remoteMode
	isServerStart = false

	cmds = append(cmds, privateNanoToTime())
	cmds = append(cmds, privateMMSToTime())

	if remoteMode == false && readMode == false {
		cmds = append(cmds, privateConsoleSet())
		cmds = append(cmds, privateConsoleClear())
	}
	cmds = append(cmds, privateFileCmds()...)

	remoteSkipIndex := -1
	if isRemoteClientStart == true {
		for i := 0; i < len(cmds); i++ {
			cc := cmds[i].(Ccmd)
			if cc.Cmd == KEYWORD_REMOTE {
				remoteSkipIndex = i
				cc.Cmd = KEYWORD_REMOTE_SWITCH_CMD
				cmds[i] = cc
			}
		} //for
		cmds = append(cmds, getRemoteConsole())
	} //if
	cmds = append(cmds, getHelpCmdLO(), getHelpCmdUP())

	maxIDX := 0
	maxLen := 0
	for i, cmd := range cmds {
		if remoteSkipIndex == i {
			continue
		}
		nLen := len(cmd.ICommand())
		if nLen > maxLen {
			maxLen = nLen
			maxIDX = i
		}
	} //for
	for i, cmd := range cmds {
		if remoteSkipIndex == i {
			continue
		}
		if i == maxIDX {
			helpTaps = append(helpTaps, "")
			continue
		}
		nLen := len(cmd.ICommand())
		if nLen == maxLen {
			helpTaps = append(helpTaps, "")
		} else {
			gapLen := maxLen - nLen
			gapTap := ""
			for j := 0; j < gapLen; j++ {
				gapTap = fmt.Sprintf("%v ", gapTap)
			}
			helpTaps = append(helpTaps, gapTap)
		}
	} //for

	if isRemoteClientStart {
		// go remoteRun()
	} else {
		if readMode == false {
			go run()
		}
	}

}

//ChangeRemoteConsoleTag :
func ChangeRemoteConsoleTag(tag string) {
	remoteConsoleTag = tag
}

func remoteRun(tag ...string) {
	fmt.Println(getInfoTag(len(cmds), true))

	isTag := false
	if len(tag) > 0 {
		isTag = true
		remoteConsoleTag = tag[0]
	}

	reader := bufio.NewReader(os.Stdin)
	//OUT:
	for {
		promptMsg := ""
		if isTag == true {
			if remoteModeOn == true {
				promptMsg = remoteConsoleTag + " remote on>"
			} else {
				promptMsg = remoteConsoleTag + " remote off>"
			}
		} else {
			if remoteModeOn == true {
				promptMsg = "remote on>"
			} else {
				promptMsg = "remote off>"
			}
		}
		if lout.Do() {
			promptMsg = lout.Tag + promptMsg
		}
		if remoteModeOn {
			fmt.Print(cc.CyanBold)
			fmt.Print(promptMsg)
			fmt.Print(cc.END)
		} else {
			fmt.Print(promptMsg)
		}

		line, _ := reader.ReadString('\n')
		buf := []byte(line)
		switch runtime.GOOS {
		case "windows":
			buf = buf[:len(buf)-2]
		default:
			buf = buf[:len(buf)-1]
		} //switch

		line = string(buf)
		line = strings.Trim(line, " ")
		if line == "" {
			continue
		}
		if remoteModeOn == true {
			isSkip := false
			for _, keyword := range removePrefixKeywords {
				if strings.HasPrefix(line, keyword) {
					isSkip = true
					break
				}
			} //for

			if isSkip == false {
				if strings.HasPrefix(line, KEYWORD_REMOTE) == false {
					checkHelp := strings.ToLower(line)
					if strings.HasPrefix(checkHelp, "help") {
						line = fmt.Sprintf("%v  %v", KEYWORD_REMOTE_SWITCH_CMD, line)
					} else {
						line = fmt.Sprintf("%v %v", KEYWORD_REMOTE_SWITCH_CMD, line)
					}
				}
			}
		} //if
		//fmt.Println("->", line)
		workCmds(line)
	} //for
}

func run() {
	fmt.Println(getInfoTag(len(cmds), false))

	reader := bufio.NewReader(os.Stdin)
	//OUT:
	for {

		promptMsg := ""
		if isServerStart == false {
			promptMsg = readyMessage
		} else {
			promptMsg = runtimeMessage
		}
		if prefixCmdMessage != "" {
			promptMsg = promptMsg + "[" + prefixCmdMessage + "]"
		}
		if lout.Do() {
			promptMsg = lout.Tag + promptMsg
		}
		fmt.Print(fmt.Sprintf("%v>", promptMsg))
		//fmt.Print(promptTag)

		line, _ := reader.ReadString('\n')
		buf := []byte(line)
		switch runtime.GOOS {
		case "windows":
			buf = buf[:len(buf)-2]
		default:
			buf = buf[:len(buf)-1]
		} //switch

		line = string(buf)
		line = strings.Trim(line, " ")
		if line == "" {
			continue
		}

		upateLinePrefix(&line)
		workCmds(line)

		//break OUT
	} //for
	//dbg.PrintForce("console.run end....")
}

//TestCommand :
func TestCommand(line string, remoteFlag ...bool) {
	defer muCommand.Unlock()
	muCommand.Lock()
	//dbg.Purple("TestCommand :", line)
	//fmt.Println("TestCommand", line)
	isRemoteMode := false
	if len(remoteFlag) > 0 {
		isRemoteMode = remoteFlag[0]
		if isRemoteMode == true {
			startWriteMode()
		}
	}

	done := make(chan bool)
	ss := strings.Split(line, " ")

	for _, cmd := range cmds {
		if ss[0] == cmd.ICommand() {

			ps := ss[1:len(ss)]
			if len(ps) > 0 {
				if checkHelp(ps[0]) == true {
					Log(cc.Yellow)
					Log("---------------------------------------")
					Log(cmd.IHelp())
					Log("---------------------------------------")
					Log(":ok")
					Log(cc.END)
					break
				}
			}

			args, err := getArgs(cmd, ps)
			if err != nil {
				Log(err.Error())
				break
			}

			go cmd.IWork(done, args)

			if result := <-done; result == true {
				Log(":ok")
			} else {
				Log(":done<-false")
			}
			break
		} //if
	} //for

	if isRemoteMode == true {
		removeCallback()
	}

}

func workCmds(line string) {
	defer muCommand.Unlock()
	muCommand.Lock()
	//fmt.Println("workCmds", line)
	done := make(chan bool)
	ss := strings.Split(line, " ")

	for _, cmd := range cmds {
		if ss[0] == cmd.ICommand() {
			lout.Title(line)

			ps := ss[1:len(ss)]
			if len(ps) > 0 {
				if checkHelp(ps[0]) == true {
					Log(cc.Yellow)
					Log("---------------------------------------")
					Log(cmd.IHelp())
					Log("---------------------------------------")
					Log(":ok")
					Log(cc.END)
					break
				}
			}

			args := []string{}
			if cmd.ISenderCmd() == false {
				var err error
				args, err = getArgs(cmd, ps)
				if err != nil {
					Log(err.Error())
					break
				}
			} else {
				args = ps
			}

			go cmd.IWork(done, args)

			if isRemoteClientStart && remoteModeOn {
				<-done
			} else {
				if result := <-done; result == true {
					Log(":ok")
				} else {
					Log(":done<-false")
				}
			}
			lout.Write()
			break
		} //if
	} //for

}
