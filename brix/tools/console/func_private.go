package console

import "strings"

func checkHelp(cmd string) bool {
	if cmd == KEYWORD_HELP_LO || cmd == KEYWORD_HELP_UP {
		return true
	}
	return false
}

func helpView(isSimple bool) {
	line := ""
	if isSimple == true {
		line = LINE_LO
	} else {
		line = LINE_UP
	}
	Log()
	if isUserHelp == true {
		Log(line)
		Log(userHelp)
	}
	Log(line)
	Log("Command help : show cmd to help.")
	Log(line)

	bSkip := false
	for i, c := range cmds {
		if c.ICommand() == KEYWORD_REMOTE_SWITCH_CMD {
			bSkip = true
			continue
		}
		if bSkip == true {
			i--
		}

		if c.INoParams() == true {
			Log("#", c.ICommand(), helpTaps[i], c.IHeader(), " -- (no params)") //" -- (cmds widthout params.)"
		} else {
			Log("#", c.ICommand(), helpTaps[i], c.IHeader())
		}
		if isSimple == false {
			ih := c.IHelp()
			ih = strings.Trim(ih, " ")
			if ih != "" {
				Log("-> ", c.IHelp())
			}
			if i < len(cmds)-1 {
				Log(LINE_LO)
			}
		}
	} //for

	Log(line)
}

// RecoverPs :
func recoverPs(done chan<- bool) {
	if err := recover(); err != nil {
		Log("command error", err)
		done <- true
	}
}
