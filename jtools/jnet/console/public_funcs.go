package console

import "strings"

type IConsole interface {
	SetLogFunc(f func(a ...interface{})) IConsole
	SetErrFunc(f func(a ...interface{})) IConsole
	Clear()
	Start()
	Append(list ...*CMD)
	AppendList(list CMDList)
	SetTitle(v string) IConsole
	Call(line string)

	Log(a ...interface{})
	Error(a ...interface{})
	Atap() //---
	Btap() //===
}

func New() IConsole {
	return &cConsole{}
}

func (my *cConsole) SetLogFunc(f func(a ...interface{})) IConsole {
	defer my.mu.Unlock()
	my.mu.Lock()
	my.log_func = f
	return my
}
func (my *cConsole) SetErrFunc(f func(a ...interface{})) IConsole {
	defer my.mu.Unlock()
	my.mu.Lock()
	my.err_func = f
	return my
}

func (my *cConsole) Clear() {
	defer my.mu.Unlock()
	my.mu.Lock()
	if my.isStart {
		return
	}
	my.cmds = my.cmds[:0]
}

func (my *cConsole) Start() {
	defer my.mu.Unlock()
	my.mu.Lock()
	if my.isStart {
		return
	}
	go my.run()
	my.isStart = true
}

func (my *cConsole) Append(list ...*CMD) {
	defer my.mu.Unlock()
	my.mu.Lock()
	my.cmds = append(my.cmds, list...)
}
func (my *cConsole) AppendList(list CMDList) {
	defer my.mu.Unlock()
	my.mu.Lock()
	my.cmds = append(my.cmds, list...)
}

func (my *cConsole) SetTitle(v string) IConsole {
	defer my.mu.Unlock()
	my.mu.Lock()
	my.titleName = v
	return my
}

func (my *cConsole) Call(line string) {
	defer my.mu.Unlock()
	my.mu.Lock()

	line = strings.TrimSpace(line)
	name := strings.Split(line, " ")[0]
	if v, do := checkHelp(line); do {
		viewMainHelp(my, v)
		return
	}
	line = strings.TrimSpace(line[len(name):])
	ok := parseCmd(
		my,
		my.cmds,
		name,
		line,
	)
	if !ok {
		my.Error(notSupportd)
	}
}
