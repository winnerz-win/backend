package dbg

import (
	"errors"
	"fmt"
	"jtools/cc"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

func Error(a ...interface{}) error {
	if len(a) == 0 {
		return errors.New(Stack())
	}

	sl := []string{}
	for i := 0; i < len(a); i++ {
		sl = append(sl, "%v")
	}
	sf := strings.Join(sl, " ")

	return fmt.Errorf(sf, a...)

}

func Errort(a ...interface{}) error {
	if len(a) == 0 {
		return errors.New(Stack())
	}

	text := Cat(a...)
	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ReplaceAll(text, " ", "")
	return errors.New(text)
}

func Exit(a ...interface{}) {
	Println(a...)
	Println(
		Stack(),
	)
	os.Exit(1)
}

func Stack() string {
	sl := strings.Split(string(debug.Stack()), "\n")
	if len(sl) >= 5 {
		/*
			[ 0 ] goroutine 6 [running]:
			[ 1 ] runtime/debug.Stack()
			[ 2 ] 	C:/Program Files/Go/src/runtime/debug/stack.go:24 +0x7a
			[ 3 ] jtools/dbg._print_stack()
			[ 4 ] 	d:/work/go/src/brix_pkg/jtools/dbg/error.go:32 +0x2e
		*/
		sl = sl[5:]
	}
	// for i, v := range sl {
	// 	Println("[", i, "]", v)
	// }

	return strings.Join(sl, "\n")
}

// StackError :
func StackError(a ...interface{}) string {
	subject := []interface{}{"[ StackError ]  "}
	subject = append(subject, a...)
	subject = append(subject, "-------------------------------\n")
	stackString := string(debug.Stack())

	// cc.RedItalicBG(subject...)
	// cc.RedItalic(stackString)

	msg := ""
	for _, s := range subject {
		msg = fmt.Sprintf("%v%v\n", msg, s)
	} //for
	msg = fmt.Sprintf("%v%v\n", msg, stackString)
	return msg

}

func VRecover(do func()) {
	viewPanic := func(panic interface{}, do func()) {
		if panic != nil {
			Println(cc.D_Red)
			Println("###### PANIC [VRecover] #####################################################################################")
			Printf(": %v\n", time.Now())
			Printf(": %v\n", panic)
			Println("=============================================================================================================")
			sl := strings.Split(string(debug.Stack()), "\n")
			doLog := false
			callStack := ""
			for _, v := range sl {
				if doLog == false {
					if strings.Contains(v, "panic.go") {
						doLog = true
					}
				} else {
					Printf("%v\n", v)
					callStack = fmt.Sprintf("%v %v\n", callStack, v)
				}
			} //for
			Println("#############################################################################################################")
			Println(cc.D_End)
			// fmt.Println("**** CALLSTACK ****")
			// fmt.Println(callStack)
			// fmt.Println("*******************")
			//debug.PrintStack()
			if do != nil {
				do()
			}
		}
	}
	if e := recover(); e != nil {
		viewPanic(e, do)
	}
}
