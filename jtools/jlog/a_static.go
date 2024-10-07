package jlog

import (
	"jtools/cc"
	"jtools/dbg"
	"os"
	"strings"
)

const (
	RED    = "__31__" //31
	GREEN  = "__32__" //32
	YELLOW = "__33__" //33
	BLUE   = "__34__" //34
	PURPLE = "__35__" //35
	CYAN   = "__36__" //36
	WHITE  = "__37__" //37
)

var (
	_msg_colors = map[string]int{
		RED:    31,
		GREEN:  32,
		YELLOW: 33,
		BLUE:   34,
		PURPLE: 35,
		CYAN:   36,
		WHITE:  37,
	}
)

func prefix_message_color(message string) (string, int) {
	for tag, color := range _msg_colors {
		if strings.HasPrefix(message, tag) {
			message = message[len(tag):]
			return message, color
		}
	} //for
	return message, 0
}

var (
	logEntry *LogEntry
)

func Init(config ConfigLogYAML) {
	logEntry = new(config)
}
func GetEntry() *LogEntry {
	return logEntry
}

func Panic(args ...interface{}) {
	stack_err_msg := dbg.StackError(args...)
	if logEntry == nil {
		cc.Red(stack_err_msg)
		return
	}

	logEntry.Panic(stack_err_msg)
}

func Exit(args ...interface{}) {
	defer func() {
		recover()
		os.Exit(1)
	}()

	stack_err_msg := dbg.StackError(args...)
	if logEntry == nil {
		cc.Red(stack_err_msg)
		return
	}

	logEntry.Panic(stack_err_msg)
}

func Error(args ...interface{}) {
	if logEntry == nil {
		cc.Red(args...)
		return
	}
	logEntry.Error(args...)
}
func Warn(args ...interface{}) {
	if logEntry == nil {
		cc.Yellow(args...)
		return
	}
	logEntry.Warn(args...)
}
func Info(args ...interface{}) {
	if logEntry == nil {
		cc.Cyan(args...)
		return
	}
	logEntry.Info(args...)
}
func Debug(args ...interface{}) {
	if logEntry == nil {
		cc.Green(args...)
		return
	}
	logEntry.Debug(args...)
}
func Trace(args ...interface{}) {
	if logEntry == nil {
		cc.GreenBold(args...)
		return
	}
	logEntry.Trace(args...)
}
