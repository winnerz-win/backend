package chttp

type Starter interface {
	Start()
	Close()
	WaitStart()
	EndC() <-chan struct{}
}
type StarterList []Starter

type cStart struct {
	name   string
	seq    int8
	startC chan struct{}
	endC   chan struct{}
}

func NewStarter(name string) Starter {
	s := &cStart{
		name:   name,
		startC: make(chan struct{}),
		endC:   make(chan struct{}),
	}
	return s
}

func (my *cStart) Start() {
	if my.seq == 0 {
		close(my.startC)
		my.seq = 1
	}

}
func (my *cStart) Close() {
	if my.seq == 1 {
		close(my.endC)
		my.seq = 2
	}
}
func (my *cStart) WaitStart() {
	<-my.startC
}
func (my *cStart) EndC() <-chan struct{} {
	return my.endC
}
