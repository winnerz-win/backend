package mms

import (
	"fmt"
	"jtools/cc"
	"strconv"
	"strings"
	"time"
)

const (
	//DivUnix : unix = nano / val
	DivUnix = 1000000000 //1sec
	//DivMMS :
	DivMMS = 1000000

	divMMSec = 1000
	//1000000
	//1000000000

	//MinusOne    = time.Dration(-1)
	Day         = time.Hour * 24
	Hour        = time.Hour
	Minute      = time.Minute
	Second      = time.Second
	Millisecond = time.Millisecond
	Month       = Day * 30

	DaySec = (60 * 60) * 24 //86400

	SecMMS  = divMMSec
	MinMMS  = SecMMS * 60
	HourMMS = MinMMS * 60
	DayMMS  = HourMMS * 24

	ZERO = MMS(0)
)

// MMS : milisecond
type MMS int64

func utc() time.Time {
	return time.Now().UTC()
}

// Now : mms
func Now() MMS {
	// if timeFaster.IsStart {
	// 	return timeFaster.Now()
	// }
	return MMS(utc().UnixNano() / DivMMS)
}

// Zero : 2020-01-01 00:00
func Zero() MMS {
	nt := utc()
	nt = toDate(nt.Year(), nt.Month(), nt.Day(), 0, 0, 0, 0)
	return FromTime(nt)
}

// Value :
func (my MMS) Value(isTotime ...bool) int64 {
	ivalue := int64(my)
	if len(isTotime) > 0 && isTotime[0] {

	} else {
		str := fmt.Sprint(ivalue)
		if len(str) > 13 {
			loop := len(str) - 13
			var div int64 = 1
			for i := 0; i < loop; i++ {
				div *= 10
			}
			return ivalue / div
		} else if len(str) == 10 {
			ivalue *= 1000
		}
	}

	return ivalue
}
func (my MMS) Refact13() MMS {
	ms := my.Int64String()
	size := len(ms)
	if size < 13 {
		pushCnt := 13 - size
		for i := 0; i < pushCnt; i++ {
			ms += "0"
		}
	} else if size > 13 {
		ms = ms[:13]
	}
	v, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		cc.RedItalic("Refact13 :", err)
		return my
	}
	return MMS(v)
}
func (my *MMS) Refact13Self() {
	*my = my.Refact13()
}
func (my MMS) Sec() MMS {
	ivalue := int64(my)
	return MMS(ivalue / divMMSec * divMMSec)
}

func (my MMS) Unix() int64 { return my.SecInt64() }
func (my MMS) SecInt64() int64 {
	ivalue := int64(my)
	return ivalue / divMMSec
}
func (my MMS) SecInt() int { return int(my.SecInt64()) }

func (my MMS) SecIntString() string { return fmt.Sprint(my.SecInt()) }

func (my MMS) Int64String() string { return fmt.Sprint(my.Int64()) }

// Int64 :
func (my MMS) Int64() int64 {
	return int64(my)
}

// String :
func (my MMS) String() string {
	s := fmt.Sprintf("%v", my.ToTime())
	switch len(s) {
	case 30:
		s += "   "
	case 31:
		s += "  "
	case 32:
		s += " "
	}
	return s
}
func (my MMS) String2() string {
	s := strings.ReplaceAll(fmt.Sprintf("%v", my.ToTime()), " +0000 UTC", "")
	switch len(s) {
	case 20:
		s += "   "
	case 21:
		s += "  "
	case 22:
		s += " "
	}
	return s
}

// MMSKST : UTC시간에서 9시간을 뺀 한국시로 보정한 UTC를 반환 한다.
func (my MMS) MMSKST() MMS {
	unixKST := my.Add(time.Hour * -9)
	return unixKST
}

// KST :
func (my MMS) KST() string {
	nt := my.ToTime()
	nt = nt.Add(time.Hour * 9)
	mm := fmt.Sprintf("%v", nt)
	return strings.ReplaceAll(mm, "UTC", "KST")
}

// KST :
func KST() string {
	return Now().KST()
}

// ToTime :
func (my MMS) ToTime() time.Time {
	return ToTime(my)
}

// ToZero : 00:00
func (my MMS) ToZero() MMS {
	return FromYMD(my.YMD())

	// nt := Now().ToTime()
	// nt = time.Date(nt.Year(), nt.Month(), nt.Day(), 0, 0, 0, 0, time.UTC)
	// return FromTime(nt)
}

// DayRange : 00:00 ~ 23:59
func (my MMS) DayRange() (MMS, MMS) {
	start := my.ToZero()
	end := MMS(start.Add(Day) - 1)
	return start, end
}

// YMD :
func (my MMS) YMD(hmsn ...int) int {
	return YMDN(my, hmsn...)
}

// YMDString :
func (my MMS) YMDString(hmsn ...int) string {
	return YMDString(my, hmsn...)
}

// Add :
func (my MMS) Add(d time.Duration) MMS {
	t := my.ToTime().Add(d)
	return FromTime(t)
}

// Sub : mms - mms = time.Duration
func (my MMS) Sub(m MMS) time.Duration {
	v := my - m
	return time.Duration(v.Value() * DivMMS)
}

// Duration : time.Duration
func (my MMS) Duration() time.Duration {
	return time.Duration(my.Value() * DivMMS)
}

// CompareTo :
func (my MMS) CompareTo(m MMS) int {
	if my.Value() == m.Value() {
		return 0
	} else if my.Value() > m.Value() {
		return 1
	} else {
		return -1
	}
}

// SnapShot : //지난 시간
func (my MMS) SnapShot() time.Duration {
	return Now().Sub(my)
}
