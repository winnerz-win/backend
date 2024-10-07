package console

import "fmt"

const (
	_lineA = "--------------------------------------------------------"
	_lineB = "========================================================"
	_helpS = "shot"
	_helpL = "long"
)

func (my *cConsole) Log(a ...interface{}) {
	if my.log_func != nil {
		my.log_func(a...)
	} else {
		fmt.Println(a...)
	}
}
func (my *cConsole) Error(a ...interface{}) {
	if my.err_func != nil {
		my.err_func(a...)
	} else {
		fmt.Println(a...)
	}
}

func (my *cConsole) Atap() {
	my.Log(_lineA)
}
func (my *cConsole) Btap() {
	my.Log(_lineB)
}

func (my *cConsole) initHelp() {

	alineHelpSpace(my.cmds, true)
}

func alineHelpSpace(cmds CMDList, isFirst bool) {
	max := 0
	for _, v := range cmds {
		if max < len(v.Name) {
			max = len(v.Name)
		}
	} //for
	max += 15
	getTap := func(cnt int, name string) string {
		cnt = cnt - len(name)
		return getSpace(cnt)
	}
	for i := range cmds {
		if isFirst {
			cmds[i].ready(nil)
		}
		cmds[i].nameGap = getSpace(len(cmds[i].Name))
		cmds[i].helpTapS = getTap(max, cmds[i].Name)
	} //for
}

func getSpace(cnt int) string {
	space := ""
	for i := 0; i < cnt; i++ {
		space += " "
	}
	return space
}

func checkHelp(v string) (string, bool) {
	switch v {
	case "help", "--help":
		return _helpS, true
	case "HELP", "--HELP":
		return _helpL, true

	}
	return "", false
}

func getSubAlerm(cmd *CMD) string {
	subAlerm := ""
	subCnt := len(cmd.CMDS)
	if subCnt > 0 {
		subAlerm = fmt.Sprintf("[ +%v ]", subCnt)
	}
	return subAlerm
}

func viewMainHelp(my *cConsole, v string) {

	my.Log(_lineB)
	if v == _helpS {
		for _, cmd := range my.cmds {
			my.Log(cmd.Name, cmd.helpTapS, getSubAlerm(cmd), cmd.HelpS)
		} //for
	} else {
		for _, cmd := range my.cmds {
			msg := cmd.HelpS
			if cmd.HelpL != "" {
				msg = cmd.HelpL
			}
			my.Log(cmd.Name, cmd.helpTapS, msg)
			my.Log(_lineA)
		} //for
	}
	my.Log(_lineB)
}

func viewCmdHelp(my *cConsole, cmd *CMD, isSubCheck bool) {
	hlog := my.Log
	if isSubCheck {
		hlog(_lineB)
	}

	hlog(cmd.fullName(), cmd.helpTapS, getSubAlerm(cmd), cmd.HelpS)
	hlog(_lineA)

	if cmd.HelpL != "" {
		hlog(cmd.HelpL)
		hlog(_lineA)
	}

	if isSubCheck {
		for _, sub := range cmd.CMDS {
			viewCmdHelp(my, sub, false)
		} //for
	}

}
