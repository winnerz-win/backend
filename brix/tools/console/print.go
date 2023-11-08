package console

import (
	"fmt"

	"txscheduler/brix/tools/console/lout"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/dbg/cc"
)

////////////////////////////////////////

//Println :
func Println(a ...interface{}) {
	as := []interface{}{}
	as = append(as, ":")
	as = append(as, a...)

	Log(as...)
}

//LogForce :
func LogForce(a ...interface{}) {
	as := []interface{}{}
	as = append(as, "▷")
	as = append(as, a...)
	Log(as...)
}

//Log :
func Log(a ...interface{}) {
	dbg.CallLogWriteln(a...)

	lout.Log(a...)
	if isRemoteWriteMode() {
		remoteMsg(a...)
	} else {
		fmt.Println(a...)
	}
}
func Red(a ...interface{})     { ColorLog(cc.Red, a...) }
func RedBold(a ...interface{}) { ColorLog(cc.RedBold, a...) }

func Green(a ...interface{})     { ColorLog(cc.Green, a...) }
func GreenBold(a ...interface{}) { ColorLog(cc.GreenBold, a...) }

func Yellow(a ...interface{})     { ColorLog(cc.Yellow, a...) }
func YellowBold(a ...interface{}) { ColorLog(cc.YellowBold, a...) }

func Blue(a ...interface{})     { ColorLog(cc.Blue, a...) }
func BlueBold(a ...interface{}) { ColorLog(cc.BlueBold, a...) }

func Purple(a ...interface{})     { ColorLog(cc.Purple, a...) }
func PurpleBold(a ...interface{}) { ColorLog(cc.PurpleBold, a...) }

func Cyan(a ...interface{})     { ColorLog(cc.Cyan, a...) }
func CyanBold(a ...interface{}) { ColorLog(cc.CyanBold, a...) }

func ColorLog(color cc.COLOR, a ...interface{}) {
	lout.Log(a...)
	if isRemoteWriteMode() {
		remoteMsg(a...)
		remoteMessage = color.String() + remoteMessage + cc.END.String()
	} else {
		dbg.Color(color, a...)
	}
}

const (
	ALine = "------------------------------------------------------------"
	BLine = "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"
	CLine = "============================================================"
	DLine = "............................................................"
)

//Atap : --------
func Atap() { Log(ALine) }

//Btap : ~~~~~~~~
func Btap() { Log(BLine) }

//Ctap : ========
func Ctap() { Log(CLine) }

//Dtap : ........
func Dtap() { Log(DLine) }
