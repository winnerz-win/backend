package runtext

import (
	"sync"
	"time"
)

const (
	TAGPause  = "pause"
	TAGResume = "resume"
)

//GetTag :
func (my *Runtext) GetTag() string {
	defer my.mu.RUnlock()
	my.mu.RLock()
	return my.tag
}

//GetTagIndex :
func (my *Runtext) GetTagIndex() int {
	defer my.mu.RUnlock()
	my.mu.RLock()
	return my.tagIndex
}

//SetTag :
func (my *Runtext) SetTag(v string) int {
	defer my.mu.Unlock()
	my.mu.Lock()
	my.tag = v
	my.tagIndex++
	return my.tagIndex
}

//WaitTag :
func (my *Runtext) WaitTag(v string) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(tag string) {
		defer wg.Done()
		for {
			if my.GetTag() == tag {
				return
			}
			time.Sleep(time.Millisecond * 100)
		} //for
	}(v)
	wg.Wait()
}

//Pause :
func (my *Runtext) Pause() bool {
	defer my.mu.Unlock()
	my.mu.Lock()
	_isPause := !my.isPause
	my.isPause = true
	return _isPause
}

//Resume :
func (my *Runtext) Resume() bool {
	defer my.mu.Unlock()
	my.mu.Lock()
	_isResume := my.isPause
	my.isPause = false
	return _isResume
}

//IsPause :
func (my *Runtext) IsPause() bool {
	defer my.mu.RUnlock()
	my.mu.RLock()
	return my.isPause
}
