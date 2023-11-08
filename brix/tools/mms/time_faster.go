package mms

import (
	"fmt"

	"txscheduler/brix/tools/dbg"
)

//TimeFaster :
type TimeFaster struct {
	IsStart  bool  `json:"is_start"`
	StartYMD int   `json:"start_ymd"`
	StartAt  int64 `json:"start_at"`
	SecValue int64 `json:"sec_value"`
}

//TagString :
func (TimeFaster) TagString() []string {
	return []string{
		"is_start", "시간 증가 모드 발동 여부",
		"start_ymd", "기준날짜",
		"start_at", "기준 시간",
		"sec_value", "초당 증가 폭",
	}
}

func (my TimeFaster) String() string { return dbg.ToJSONString(my) }

var (
	timeFaster = newTimeFaster(20200101, 30)
)

//FasterClone :
func FasterClone() TimeFaster {
	f := TimeFaster{
		IsStart:  timeFaster.IsStart,
		StartYMD: timeFaster.StartYMD,
		StartAt:  timeFaster.StartAt,
		SecValue: timeFaster.SecValue,
	}
	return f
}

//FasterStop :
func FasterStop() {
	timeFaster.IsStart = false
	dbg.YellowBold("#################################")
	dbg.YellowBold(" mms.TimeFaster MODE  --- END")
	dbg.YellowBold("#################################")
	fmt.Println()
}

//FasterTime :
func FasterTime(ymd int, minToDay int) {
	timeFaster = newTimeFaster(ymd, minToDay)
	timeFaster.IsStart = true
	dbg.YellowBold("#################################")
	dbg.YellowBold(" mms.TimeFaster MODE  --- START")
	dbg.YellowBold("#################################")
	dbg.YellowBold(timeFaster)
	dbg.YellowBold("#################################")
	fmt.Println()
}

//newTimeFaster :
func newTimeFaster(ymd int, minToDay int) TimeFaster {
	f := TimeFaster{
		StartYMD: ymd,
		StartAt:  int64(FromYMD(ymd)),
	}
	var day int64 = (60 * 60) * 24 * 1000
	f.SecValue = day / int64(60*1000*minToDay)
	return f
}

//Now :
func (my TimeFaster) Now() MMS {
	nt := int64(MMS(utc().UnixNano() / DivMMS))
	jt := (nt - my.StartAt) / 1000 //sec
	dt := my.StartAt + jt*my.SecValue*1000
	return MMS(dt)
}

//HelpString :
func (TimeFaster) HelpString() string {
	return ` -
		jt := (현재시간 - tm.start_at) / 1000
		mms := tm.start_at + jt * tm.sec_value * 1000
		return mms
	`
}
