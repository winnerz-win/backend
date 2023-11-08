package dbg

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"txscheduler/brix/tools/dbg/cc"
)

const (
	//ENTER :
	ENTER = "\n"
)

var (
	skipLog  = false
	testMode = false
)

// SetSkipLog :
func SetSkipLog(flag bool) {
	skipLog = flag
}

// SetTestMode :
func SetTestMode(flag bool) {
	testMode = flag
}

func Exit(a ...interface{}) {
	fmt.Println(a...)
	fmt.Println(
		Stack(),
	)
	os.Exit(1)
}

func _println(a ...interface{}) (int, error) {
	_logWriteln(a...)
	return fmt.Println(a...)
}
func _printf(format string, a ...interface{}) {
	_logWriteln(fmt.Sprintf(format, a...))
	fmt.Printf(format, a...)
}

// PrintInit :
func PrintInit(a ...interface{}) {
	as := []interface{}{}
	as = append(as, cc.Yellow, "[INIT] ")
	as = append(as, a...)
	as = append(as, cc.END)

	_println(as...)
}

// Println :
func Println(a ...interface{}) (n int, err error) {
	_logWriteln(a...)

	if !isColorLogForce {
		if skipLog {
			return 0, nil
		}
	}

	return fmt.Println(a...)
}

// PrintForce :
func PrintForce(a ...interface{}) (n int, err error) {
	as := []interface{}{}
	as = append(as, cc.Green, "▷")
	as = append(as, a...)
	as = append(as, cc.END)

	return _println(as...)
}

// HandlerCall :
func HandlerCall(a ...interface{}) {
	as := []interface{}{}
	as = append(as, "[ HANDLERFUNC ] ")
	as = append(as, a...)
	fmt.Println(as...)
}

// IToString :
type IToString interface {
	ToString() string
}

// SliceToString : ToString() 이 정의된 []Struct 를 []ToString() 한다.
func SliceToString(sl []IToString) string {
	msg := ""
	for i, v := range sl {
		if i == 0 {
			msg = fmt.Sprintf("%v\n", v.ToString())
		} else {
			msg = fmt.Sprintf("%v%v\n", msg, v.ToString())
		}
	}
	return msg
}

// SlString :
func SlString(sl ...string) string {
	msg := "{\n"
	for _, v := range sl {
		msg += fmt.Sprintf("  %v\n", v)
	}
	msg += "}\n"

	return msg
}

// ToTimeString :
func ToTimeString(nano int64, mmsSkip ...bool) string {
	if nano == 0 {
		return "0"
	}
	return ToTimeStringT(time.Unix(0, nano), mmsSkip...)
}

// ToTimeStringT :
func ToTimeStringT(nt time.Time, mmsSkip ...bool) string {
	nt = nt.UTC().Add(time.Hour * 9)
	if len(mmsSkip) > 0 && mmsSkip[0] == true {
		ss := strings.Split(fmt.Sprintf("%v", nt), ".")
		return ss[0] + " KST"
	}
	msg := fmt.Sprintf("%v", nt)
	return strings.Replace(msg, "UTC", "KST", 1)
}

// Sprintln :
func Sprintln(a ...interface{}) string {
	msg := ""
	for i, c := range a {
		if i == 0 {
			msg = fmt.Sprintf("%v", c)
		} else {
			msg = fmt.Sprintf("%v %v", msg, c)
		}
	} //for
	return msg
}

// Printf :
func Printf(format string, a ...interface{}) (n int, err error) {
	if skipLog == true {
		return 0, nil
	}
	return fmt.Printf(format, a...)
}

// PrintError :
func PrintError(a ...interface{}) {
	if len(a) > 0 {
		st := fmt.Sprintf("%v", time.Now().UTC().Add(time.Hour*9))
		ss := strings.Split(st, " ") //[2019-06-28] [00:00:00.0000000] [+0000] [UTC]
		aa := []interface{}{"### ERROR ###", ss[0], ss[1]}
		for _, v := range a {
			aa = append(aa, v)
		}
		aa = append(aa)

		fmt.Println(cc.Red)
		_println(aa...)
		fmt.Println(cc.END)
	}
}

// CStack :
func CStack(skipDepth ...int) {
	defaultSkipLine := 5
	if len(skipDepth) > 0 {
		defaultSkipLine += skipDepth[0] * 2
	}
	_println(cc.Yellow, "--------------------------------------------------------------------")
	sl := strings.Split(string(debug.Stack()), "\n")
	for _, v := range sl {
		if defaultSkipLine > 0 {
			defaultSkipLine--
			continue
		}
		_printf("%v\n", v)
	}
	_println("--------------------------------------------------------------------", cc.END)
}

// PrintStack :
func PrintStack(a ...interface{}) {
	_println(cc.Yellow, "[STACK]")
	_println(a...)
	sl := strings.Split(string(debug.Stack()), "\n")
	for _, v := range sl {
		_printf("%v\n", v)
	}
	_println("------------------------------------", cc.END)
}

// StackError :
func StackError(a ...interface{}) string {
	subject := []interface{}{"[ StackError ]  "}
	subject = append(subject, a...)
	subject = append(subject, "          ")
	stackString := string(debug.Stack())

	RedItalicBG(subject...)
	RedItalic(stackString)

	msg := ""
	for _, s := range subject {
		msg = fmt.Sprintf("%v%v\n", msg, s)
	} //for
	msg = fmt.Sprintf("%v%v\n", msg, stackString)
	return msg

}

// PrintPanic :
func PrintPanic(panic interface{}, do func(err error)) {
	if panic != nil {
		_println(cc.Red)
		_printf("\n###### PANIC ######\n")
		_printf(": %v\n", time.Now())
		_printf(": %v\n", panic)
		_printf("-------------------\n")
		sl := strings.Split(string(debug.Stack()), "\n")
		doLog := false
		callStack := ""
		for _, v := range sl {
			if !doLog {
				if strings.Contains(v, "panic.go") {
					doLog = true
				}
			} else {
				_printf("%v\n", v)
				callStack = fmt.Sprintf("%v %v\n", callStack, v)
			}
		} //for
		_printf("###################%v\n\n", cc.END)
		//debug.PrintStack()
		if do != nil {
			do(fmt.Errorf("CALLSTACK>\n %v \n error:%v", callStack, panic))
		}
	}
}

// CRecover :
func CRecover(do func(err error)) {
	if e := recover(); e != nil {
		PrintPanic(e, do)
	}
}

// VRecover :
func VRecover(do func()) {
	viewPanic := func(panic interface{}, do func()) {
		if panic != nil {
			_println(cc.Red)
			_println("###### PANIC [VRecover] #####################################################################################")
			_printf(": %v\n", time.Now())
			_printf(": %v\n", panic)
			_println("=============================================================================================================")
			sl := strings.Split(string(debug.Stack()), "\n")
			doLog := false
			callStack := ""
			for _, v := range sl {
				if doLog == false {
					if strings.Contains(v, "panic.go") {
						doLog = true
					}
				} else {
					_printf("%v\n", v)
					callStack = fmt.Sprintf("%v %v\n", callStack, v)
				}
			} //for
			_println("#############################################################################################################")
			_println(cc.END)
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

// CPanic : Panic Test Function
func CPanic(a ...interface{}) {
	if testMode == true {
		panic(fmt.Sprint(a...))
	}
}

// CPanicHoldLoop :
func CPanicHoldLoop() {
	if e := recover(); e != nil {
		stackLog := string(debug.Stack())
		_println(cc.Red)
		_println("□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□")
		_println("□□□□□□□□□□□□□□□□□□□ P A N I C □□□□□□□□□□□□□□□□□□□□□□□□")
		_println("□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□")
		_println(stackLog)
		_println("□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□")
		_println("□ MESSAGE □", e)
		_println("□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□□")
		for {
			time.Sleep(time.Second)
			fmt.Print(".")
		} //for
	} //if
}
