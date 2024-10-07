package scv

import (
	"time"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/runtext"
)

type HKey string

const (
	KeyNone           = HKey("")
	ToMasterCallback  = HKey("ToMasterCallback")
	ExMasterCallback  = HKey("ExMasterCallback")
	MasterOutCallback = HKey("MasterOutCallback")
)

type CallbackItem struct {
	HKey       HKey
	StartMsg   string
	SleepDu    time.Duration
	LooperFunc func()
	//////////////////////////////////

	SubPackage func(classic *chttp.Classic) runtext.Starter //ready()
}

type CallbackList []CallbackItem

type START_PACKAGE_RUNNER func(classic *chttp.Classic) runtext.Starter

func (my CallbackList) StartPackageList() []START_PACKAGE_RUNNER {
	list := []START_PACKAGE_RUNNER{}
	for _, v := range my {
		if v.SubPackage != nil {
			list = append(list, v.SubPackage)
		}
	}
	return list
}

func (my CallbackList) StartLooper(rtx runtext.Runner) {
	<-rtx.WaitStart()
	for _, item := range my {
		if item.LooperFunc == nil {
			continue
		}

		go func(item CallbackItem) {
			dbg.PrintForce(item.StartMsg)
			for {
				item.LooperFunc()
				time.Sleep(item.SleepDu)
			} //for
		}(item)
	} //for

}

func (my CallbackList) HKeys() []HKey {
	list := []HKey{}
	for _, v := range my {
		if v.HKey == KeyNone || v.LooperFunc == nil {
			continue
		}
		list = append(list, v.HKey)
	} //for
	return list
}
