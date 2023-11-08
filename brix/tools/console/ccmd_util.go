package console

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"txscheduler/brix/tools/console/lout"
	"txscheduler/brix/tools/jmath"
	"txscheduler/brix/tools/mms"
)

//ClearConsole : public
func ClearConsole() Ccmd {
	return Ccmd{
		Cmd:        KEYWORD_CLS,
		NoParams:   true,
		HeaderFunc: func() string { return "clear console windows" },
		Help:       "clear console.",
		Work: func(done chan<- bool, ps []string) {
			defer DoneC(done)

			switch runtime.GOOS {
			case "linux":
				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				cmd.Run()

			case "windows":
				cmd := exec.Command("cmd", "/c", "cls")
				cmd.Stdout = os.Stdout
				cmd.Run()
			default:
				Log("Not allow commamd : ", runtime.GOOS)
			}
		},
	}
}

//privateNanoToTime :
func privateNanoToTime() Ccmd {
	return Ccmd{
		Cmd:        KEYWORD_NANOTIME,
		HeaderFunc: func() string { return "view go-time [KST]" },
		HelpFunc: func() string {
			return fmt.Sprintf("%v [%v] -- view go-time [KST]", KEYWORD_NANOTIME, time.Now().UnixNano())
		},
		Work: func(done chan<- bool, ps []string) {
			defer DoneC(done)
			var err error

			strVal := strings.ReplaceAll(ps[0], ",", "")
			val := jmath.NewBigDecimal(strVal, &err)
			if err != nil {
				Log(err)
				return
			}
			nano := val.ToBigInteger().Uint64()
			utc := time.Unix(0, int64(nano)).UTC()
			kst := utc.Add(time.Hour * 9)

			rutc := strings.ReplaceAll(fmt.Sprintf("%v", utc), "UTC", "")
			rkst := strings.ReplaceAll(fmt.Sprintf("%v", kst), "UTC", "")
			Atap()
			Log("[UTC]", rutc)
			Atap()
			Log("[KST]", rkst)
			Atap()

		},
	}
}

func privateMMSToTime() Ccmd {
	return Ccmd{
		Cmd:        KEYWORD_MMSTIME,
		HeaderFunc: func() string { return "view mms-time [UTC]" },
		HelpFunc: func() string {
			now := mms.Now()
			return fmt.Sprintf("%v [%v] -- view mms-time \n %v \n %v",
				KEYWORD_MMSTIME,
				now.Value(),
				now.String(),
				now.KST(),
			)
		},
		Work: func(done chan<- bool, ps []string) {
			defer DoneC(done)

			strVal := strings.ReplaceAll(ps[0], ",", "")
			if jmath.IsUnderZero(strVal) {
				Log(strVal)
				return
			}

			val := jmath.Int64(strVal)
			mmsLong := mms.MMS(val)

			nt := mms.Now()

			Ctap()
			Log("cur :", nt)
			Log("cur :", nt.KST())
			Atap()
			Log(":", val)
			Atap()
			Log("tar :", mmsLong)
			Log("tar :", mmsLong.KST())
			Atap()

		},
	}
}

//////////////////////////////////////////////////////////

var prefixCmdMessage = ""

func upateLinePrefix(line *string) {
	lineMessage := *line

	if prefixCmdMessage != "" {
		for _, keyword := range prefixKeywordList {
			if strings.HasPrefix(lineMessage, keyword) {
				return
			}
		} //for

		*line = prefixCmdMessage + " " + lineMessage
	}
}

func privateConsoleSet() Ccmd {
	helpMsg := KEYWORD_CONSOLE_SET + " [prefixTag]"
	return Ccmd{
		Cmd:        KEYWORD_CONSOLE_SET,
		HeaderFunc: func() string { return helpMsg },
		Help:       helpMsg + "--->  xxxx[db]> show => db show",
		Work: func(done chan<- bool, ps []string) {
			defer DoneC(done)
			prefixCmdMessage = ""
			for i := 0; i < len(ps); i++ {
				prefixCmdMessage += ps[i]
				if i+1 < len(ps) {
					prefixCmdMessage += " "
				}
			} //for
		},
	}
}
func privateConsoleClear() Ccmd {
	helpMsg := KEYWORD_CONSOLE_CLEAR + " [prefixTag] remove"
	return Ccmd{
		Cmd:        KEYWORD_CONSOLE_CLEAR,
		NoParams:   true,
		HeaderFunc: func() string { return helpMsg },
		Help:       helpMsg,
		Work: func(done chan<- bool, ps []string) {
			defer DoneC(done)
			prefixCmdMessage = ""
		},
	}
}

///////////////////////////////////////////////////////////

func privateFileCmds() []ICmd {
	startMsg := KEYWORD_FILE_START + " --- query write file start"
	stopMsg := KEYWORD_FILE_CLEAR + " --- query write file stop"
	return []ICmd{
		Ccmd{
			Cmd:        KEYWORD_FILE_START,
			NoParams:   true,
			HeaderFunc: func() string { return startMsg },
			Help:       startMsg,
			Work: func(done chan<- bool, ps []string) {
				defer DoneC(done)
				Log("> File Log Mode Start")
				lout.Set(true)
			},
		},
		Ccmd{
			Cmd:        KEYWORD_FILE_CLEAR,
			NoParams:   true,
			HeaderFunc: func() string { return stopMsg },
			Help:       stopMsg,
			Work: func(done chan<- bool, ps []string) {
				defer DoneC(done)
				lout.Set(false)
				Log("> File Log Mode Stop")
			},
		},
	}
}

func CC() {
	fmt.Println("XXXXXXXXXXXXXXXXXXXXXX")
}
