package jticker

import (
	"time"
	"txscheduler/brix/tools/dbg"
)

//Ticker :
type Ticker struct {
	tick  time.Duration
	sleep time.Duration
	over  time.Time
	fSkip bool
}

//IsWait : tick 값이 안되었을경우 true , tick시간에 도달하면 false반환후 over값을 다음tick으로 갱신.
func (my *Ticker) IsWait() bool {
	if my.fSkip {
		my.fSkip = false
		return false
	}

	time.Sleep(my.sleep)
	d := time.Now().Sub(my.over)
	if d >= 0 {
		my.over = time.Now().Add(my.tick)
		return false
	}
	return true
}

//SetTick : tick-time set
func (my *Ticker) SetTick(d time.Duration) {
	my.tick = d
	my.over = time.Now().Add(my.tick)
}

//Tick : currnet-tick
func (my Ticker) Tick() time.Duration {
	return my.tick
}

//RefreshOver :
func (my *Ticker) RefreshOver() {
	my.over = time.Now().Add(my.tick)
}

//New : tick , sleep
func New(tick, sleep time.Duration, isFirstSkip ...bool) *Ticker {
	ins := &Ticker{
		tick:  tick,
		sleep: sleep,
		over:  time.Now().Add(tick),
		fSkip: dbg.BoolsOne(isFirstSkip...),
	}
	return ins
}
