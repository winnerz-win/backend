package dbg

import (
	"jtools/cc"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

// NumGoroutine :
func NumGoroutine() int {
	return runtime.NumGoroutine()
}

// ViewGo :
func ViewGo(tag ...string) {
	if len(tag) == 0 {
		cc.Yellow("go-cnt : ", runtime.NumGoroutine())
	} else {
		cc.Yellow("[", tag[0], "]go-cnt : ", runtime.NumGoroutine())
	}
}

// WaitSec : [ call() return true ] than direct return
func WaitSec(sec int, calls ...func() bool) {
	if sec <= 0 {
		sec = 1
	}
	var call func() bool = nil
	if len(calls) > 0 {
		call = calls[0]
	}

	for {
		isExit := false
		if call != nil {
			isExit = call()
		}
		time.Sleep(time.Second)
		sec--
		if sec <= 0 || isExit {
			break
		}
	} //for
}

// GoMax : runtime.GOMAXPROCS(runtime.NumCPU())
func GoMax() int {
	v := runtime.GOMAXPROCS(runtime.NumCPU())
	cc.Purple("dbg.GoMax : ", v)
	return v
}

// TestLock : channel blocking.. only use TDD
func TestLock(funcs ...func()) {
	cpuCnt := GoMax()
	cc.Purple("CPU Count :", cpuCnt)

	for _, f := range funcs {
		go f()
	} //for

	cc.Red("----- dbg.TestLock -----")
	// w := make(chan struct{}, 1)
	// <-w
	for {
		time.Sleep(time.Second)
	}
}
func TestLoop(sleepDuration time.Duration, loop_funcs ...func()) {
	cpuCnt := GoMax()
	cc.Purple("CPU Count :", cpuCnt)
	if sleepDuration <= 0 {
		sleepDuration = time.Second
	}
	for {
		for _, f := range loop_funcs {
			f()
		}
		time.Sleep(sleepDuration)
	} //for
}

// BreakLock :
func BreakLock(isFull ...bool) {
	cc.Yellow("------------ Lock Break ------------")
	dstr := string(debug.Stack())
	ss := strings.Split(dstr, "\n")
	if len(ss) > 5 {
		ss = ss[5:]
	}
	if BoolsOne(isFull...) == false && len(ss) > 2 {
		ss = ss[:2]
	}
	for _, v := range ss {
		cc.Yellow(":", v)
	}
	for {
		time.Sleep(time.Second)
	}
}

/*
flagshipdev1

flagshipcoin

flagshipdend
*/
