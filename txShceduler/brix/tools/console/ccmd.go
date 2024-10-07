package console

import (
	"fmt"
	"strings"

	"txscheduler/brix/tools/dbg/cc"
)

//ICmd :
type ICmd interface {
	ICommand() string
	IHeader() string
	IHelp() string
	IWork(done chan<- bool, ps []string)
	INoParams() bool
	ISenderCmd() bool //remote sender flag
	ArgTags() []string
}

//Ccmd : V0
type Ccmd struct {
	Cmd        string
	HeaderFunc func() string
	Help       string
	HelpFunc   func() string
	NoParams   bool
	SenderCmd  bool
	ArgTag     []string
	Work       func(done chan<- bool, ps []string)
}

//AppendCmd :
func AppendCmd(cmd string, help string, noParams bool, work func(ps []string), skipTap ...bool) {
	SetICmd(CmdFunc(cmd, help, help, noParams, work, skipTap...))
}

//CmdFunc :
func CmdFunc(cmd string, header, help string, noParams bool, work func(ps []string), skipTap ...bool) Ccmd {
	cc := Ccmd{
		Cmd:      cmd,
		NoParams: noParams,
		Work: func(done chan<- bool, ps []string) {
			defer DoneC(done)
			isTapLine := true
			if len(skipTap) > 0 && skipTap[0] {
				isTapLine = false
			}
			if isTapLine {
				defer Atap()
				Atap()
			}
			work(ps)
		},
	}
	header = strings.TrimSpace(header)
	help = strings.TrimSpace(help)
	if header != "" {
		cc.HeaderFunc = func() string {
			return header
		}
	}
	if help != "" {
		cc.Help = help
	}
	return cc
}

/*Commands : func(done chan<- bool, ps []string)
{
	Cmd, HeaderFunc, Help, HelpFunc, NoParams, Work
}
*/
type Commands []Ccmd

//ArgTags : "[" , "]"
func (c Ccmd) ArgTags() []string {
	return c.ArgTag
}

//IHeader :
func (c Ccmd) IHeader() string {
	hmsg := ""
	if c.HeaderFunc != nil {
		hmsg = c.HeaderFunc()
		if len(hmsg) > 0 {
			hmsg = fmt.Sprintf(" [ %v ]", hmsg)
		}
	} else {
		//fmt.Println(c.Cmd, "headerfunc is nil.")
	}
	return hmsg
}

//ICommand :
func (c Ccmd) ICommand() string {
	return c.Cmd
}

//IHelp :
func (c Ccmd) IHelp() string {
	if c.HelpFunc != nil {
		return c.HelpFunc()
	} else {
		return c.Help
	}

}

//INoParams :
func (c Ccmd) INoParams() bool {
	return c.NoParams
}

//ISenderCmd :
func (c Ccmd) ISenderCmd() bool {
	return c.SenderCmd
}

//IWork :
func (c Ccmd) IWork(done chan<- bool, ps []string) {
	defer recoverPs(done)

	if c.NoParams == false && checkHelp(c.Cmd) == false && len(ps) == 0 {
		Log(cc.Yellow)
		Log("")
		Log("---------------------------------------------------------")
		Log(c.IHelp())
		Log("---------------------------------------------------------")
		Log(cc.END)
		done <- false
		return
	}

	c.Work(done, ps)
}
