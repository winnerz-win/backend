package runtext

import (
	"fmt"
	"sync"
)

//Runner :
type Runner interface {
	EndC() <-chan struct{}
	WaitStart() <-chan struct{}
	IsStart() bool
}

//Runtext :
type Runtext struct {
	name    string
	isStart bool

	isPause  bool
	tag      string
	tagIndex int

	firstStart bool
	startC     chan struct{}
	endC       chan struct{}
	mu         sync.RWMutex
}

//Name :
func (my Runtext) Name() string {
	return my.name
}

//Starter : start , close interface
type Starter interface {
	Start()
	Close()
}

type dummyStarter struct{}

func (*dummyStarter) Start() {}
func (*dummyStarter) Close() {}

func DummyStarter() Starter {
	return &dummyStarter{}
}

//StarterList : []Starter : start , close interface-group
type StarterList []Starter

//Append :
func (my *StarterList) Append(s ...Starter) {
	if s == nil {
		return
	}

	for _, v := range s {
		if v == nil {
			continue
		}
		*my = append(*my, v)
	}
}

//AppendList :
func (my *StarterList) AppendList(list StarterList) {
	if list == nil {
		return
	}
	for _, v := range list {
		my.Append(v)
	}
}

//New :
func New(name ...string) *Runtext {
	rName := ""
	if len(name) > 0 {
		rName = name[0]
	}
	return &Runtext{
		name:       rName,
		tag:        "",
		tagIndex:   100,
		firstStart: false,
		startC:     make(chan struct{}, 1),
		endC:       make(chan struct{}),
	}
}

//EndC :
func (my *Runtext) EndC() <-chan struct{} {
	return my.endC
}

//Start :
func (my *Runtext) Start() {
	defer my.mu.Unlock()
	my.mu.Lock()
	my.isStart = true

	if my.firstStart == false {
		if my.name != "" {
			fmt.Println(fmt.Sprintf("[Runtext.%v].Start()", my.name))
		}
		my.firstStart = true
		close(my.startC)
	}
}

//WaitStart :
func (my *Runtext) WaitStart() <-chan struct{} {
	return my.startC
}

//IsStart :
func (my *Runtext) IsStart() bool {
	defer my.mu.RUnlock()
	my.mu.RLock()
	return my.isStart
}

//Close :
func (my *Runtext) Close() {
	defer my.mu.Unlock()
	my.mu.Lock()
	if my.isStart == false {
		return
	}
	if my.name != "" {
		fmt.Println(fmt.Sprintf("[Runtext.%v].Close()", my.name))
	}
	my.isStart = false
	close(my.endC)
}
