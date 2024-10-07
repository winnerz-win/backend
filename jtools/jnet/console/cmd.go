package console

import "strings"

/*
	CMDList : {
		{
			Name : "single-depth"
			HelpS : ""
			Action : func(ps []string){
			}
		},
		{
			Name : "multi-depth"
			HelpS : ""
			CMDS : console.CMDList{...}
		},
	}
*/
type CMDList []*CMD
type CMD struct {
	Name        string
	HelpS       string
	HelpL       string
	Pair        []string
	PairWordCut bool
	Action      func(ps []string)
	CMDS        CMDList

	////
	parent   *CMD
	helpTapS string
	nameGap  string
}

func Cmd(name string, helpS string, action func(ps []string), sub_cmds ...*CMD) *CMD {
	cmd := &CMD{
		Name:   name,
		HelpS:  helpS,
		Action: action,
	}
	if len(sub_cmds) > 0 {
		cmd.CMDS = append(cmd.CMDS, sub_cmds...)
	}
	return cmd
}

func (my *CMD) ready(p *CMD) {
	my.parent = p
	for i := range my.CMDS {
		my.CMDS[i].ready(my)
	}
	alineHelpSpace(my.CMDS, false)
}

func (my CMD) fullName() string {
	list := []string{my.Name}
	var p *CMD = my.parent
	for p != nil {
		list = append(list, p.Name)
		p = p.parent
	} //for

	rv := []string{}
	for i := len(list) - 1; i >= 0; i-- {
		rv = append(rv, list[i])
	}
	return strings.Join(rv, " ")
}
