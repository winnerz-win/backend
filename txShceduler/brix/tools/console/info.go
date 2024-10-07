package console

import (
	"fmt"
	"strings"

	"txscheduler/brix/tools/dbg/cc"
)

const (
	normalTag = `
┌──────────────────────────────────────────────────┐
│              CLI Console Module                  │
└──────────────────────────────────────────────────┘
 ▷ Console.Cmds Count : <cmds>
 ▷ Console.Run  Start . . .
 ──────────────────────────────────────────────────`
	remoteTag = `
 ┌──────────────────────────────────────────────────┐
 │            REMOTE-CLI Console Module             │
 └──────────────────────────────────────────────────┘
  ▷ Console.Cmds Count : <cmds>
  ▷ Console.Run  Start . . .
  ▷ remote on / remote off
  ──────────────────────────────────────────────────`
)

func getInfoTag(cmdSize int, remoteMode bool) string {
	size := fmt.Sprintf("%v", cmdSize)
	tag := ""
	if remoteMode {
		tag = strings.Replace(remoteTag, "<cmds>", size, 1)
	} else {
		tag = strings.Replace(normalTag, "<cmds>", size, 1)
	}
	return cc.Cyan.String() + tag + cc.END.String()
}

func getHelpCmdLO() Ccmd {
	cd := Ccmd{
		Cmd:      KEYWORD_HELP_LO,
		Help:     "Simple show command help!",
		NoParams: true,
		Work: func(done chan<- bool, ps []string) {
			helpView(true)
			done <- true
		},
	}
	cd.HeaderFunc = func() string { return cd.Help }
	return cd
}

func getHelpCmdUP() Ccmd {
	cd := Ccmd{
		Cmd:      KEYWORD_HELP_UP,
		Help:     "Detail show command help!",
		NoParams: true,
		Work: func(done chan<- bool, ps []string) {
			helpView(false)
			done <- true
		},
	}
	cd.HeaderFunc = func() string { return cd.Help }
	return cd
}

func getRemoteConsole() Ccmd {
	return Ccmd{
		Cmd:        KEYWORD_REMOTE,
		HeaderFunc: func() string { return fmt.Sprintf("%v on / %v off", KEYWORD_REMOTE, KEYWORD_REMOTE) },
		Help:       KEYWORD_REMOTE + " [on/off] ( RemoteMode ON/OFF. )",
		Work: func(done chan<- bool, ps []string) {
			defer DoneC(done)
			switch ps[0] {
			case "on":
				remoteModeOn = true
			case "off":
				remoteModeOn = false
			default:
				Log("RemoteModeOn :", remoteModeOn)
			} //switch
		},
	}
}
