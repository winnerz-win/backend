package crypt

import (
	"fmt"
	"sync"
)

type xCount struct {
	seq      int64
	mu       sync.RWMutex
	initFunc func(counter Counter, a ...interface{})
}

func NewCounter(f func(counter Counter, a ...interface{})) Counter {
	return &xCount{
		seq:      0,
		initFunc: f,
	}
}

func (my *xCount) InitFunc(a ...interface{}) {
	if my.initFunc == nil {
		return
	}
	my.initFunc(my, a...)
	my.initFunc = nil
}

func (my *xCount) SetCount(cnt int64) {
	my.mu.Lock()
	defer my.mu.Unlock()
	my.seq = cnt
}

func (my *xCount) GetCount() int64 {
	my.mu.RLock()
	defer my.mu.RUnlock()
	return my.seq
}

func (my *xCount) NextAction(callback func(next int64) bool) {
	if callback == nil {
		return
	}

	my.mu.Lock()
	defer my.mu.Unlock()

	defer func() {
		if e := recover(); e != nil {
			fmt.Println("crypt.xCount :", e)
		}
	}()

	ok := callback(my.seq + 1)
	if ok {
		my.seq++
	}

}

type Counter interface {
	InitFunc(a ...interface{})
	SetCount(cnt int64)
	GetCount() int64
	NextAction(callback func(next int64) bool)
}

/////////////////////////////////////////////////

type Lock interface {
	Action(f func())
	Action1(f func() interface{}) interface{}
}
type xlock struct {
	mu sync.RWMutex
}

func (my *xlock) Action(f func()) {
	if f == nil {
		return
	}
	my.mu.Lock()
	defer my.mu.Unlock()
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("crypt.xlock :", e)
		}
	}()

	f()
}

func (my *xlock) Action1(f func() interface{}) interface{} {
	if f == nil {
		return nil
	}
	my.mu.Lock()
	defer my.mu.Unlock()
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("crypt.xlock :", e)
		}
	}()

	return f()
}

func NewLock() Lock {
	return &xlock{}
}
