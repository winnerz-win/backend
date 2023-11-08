package mms

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"txscheduler/brix/tools/dbg"
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
)

func utc() time.Time {
	return time.Now().UTC()
}

//Now : mms
func Now() MMS {
	if timeFaster.IsStart {
		return timeFaster.Now()
	}
	return MMS(utc().UnixNano() / DivMMS)
}

//NowTime : time.Now().UTC()
func NowTime() time.Time {
	return utc()
}

//Zero : 2020-01-01 00:00
func Zero() MMS {
	nt := utc()
	nt = toDate(nt.Year(), nt.Month(), nt.Day(), 0, 0, 0, 0)
	return FromTime(nt)
}

//DayRemain : Seconds
func DayRemain(dt MMS) int64 {
	st := dt.ToTime()
	t := dt.Add(Day).ToTime()
	nextDay := toDate(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0)

	duration := nextDay.Sub(st)
	//fmt.Println(duration, int(duration.Seconds())+1)
	return int64(duration.Seconds()) + 1
}

//DayAfter : Seconds
func DayAfter(dt MMS) int64 {
	t := dt.ToTime()
	zeroDay := toDate(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0)
	duration := t.Sub(zeroDay)
	//fmt.Println(duration, int(duration.Seconds()))
	return int64(duration.Seconds())
}

//IsOver : Now() >= dt
func IsOver(dt MMS) bool {
	return Now() >= dt
}

//Elapsed : [ elapsed := tick - mms.Now().Sub(startAt) ]
func Elapsed(dt MMS, tick time.Duration) time.Duration {
	x := Now().Sub(dt)
	return tick - x
}

//MMS : milisecond
type MMS int64

//Value :
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

//Refact13 : 자리수를 13자리로 보정
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
		dbg.Red("Refact13 :", err)
		return my
	}
	return MMS(v)
}
func (my *MMS) Refact13Self() {
	*my = my.Refact13()
}

//Sec : mms단위를 000으로 셋팅.
func (my MMS) Sec() MMS {
	ivalue := int64(my)
	return MMS(ivalue / divMMSec * divMMSec)
}

//SecInt64 : mms를 제외한 초값을 반환
func (my MMS) SecInt64() int64 {
	ivalue := int64(my)
	return ivalue / divMMSec
}

//SecInt : mms를 제외한 초값을 반환
func (my MMS) SecInt() int { return int(my.SecInt64()) }

func (my MMS) SecIntString() string { return fmt.Sprint(my.SecInt()) }

//Int64String :
func (my MMS) Int64String() string { return fmt.Sprint(my.Int64()) }

//Int64 :
func (my MMS) Int64() int64 {
	return int64(my)
}

//String :
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

//MMSKST : UTC시간에서 9시간을 뺀 한국시로 보정한 UTC를 반환 한다.
func (my MMS) MMSKST() MMS {
	unixKST := my.Add(time.Hour * -9)
	return unixKST
}

//KST :
func (my MMS) KST() string {
	nt := my.ToTime()
	nt = nt.Add(time.Hour * 9)
	mm := fmt.Sprintf("%v", nt)
	return strings.ReplaceAll(mm, "UTC", "KST")
}

//KST :
func KST() string {
	return Now().KST()
}

//ToTime :
func (my MMS) ToTime() time.Time {
	return ToTime(my)
}

//ToZero : 00:00
func (my MMS) ToZero() MMS {
	return FromYMD(my.YMD())

	// nt := Now().ToTime()
	// nt = time.Date(nt.Year(), nt.Month(), nt.Day(), 0, 0, 0, 0, time.UTC)
	// return FromTime(nt)
}

//DayRange : 00:00 ~ 23:59
func (my MMS) DayRange() (MMS, MMS) {
	start := my.ToZero()
	end := MMS(start.Add(Day) - 1)
	return start, end
}

//YMD :
func (my MMS) YMD(hmsn ...int) int {
	return YMDN(my, hmsn...)
}

//YMDString :
func (my MMS) YMDString(hmsn ...int) string {
	return YMDString(my, hmsn...)
}

//Add :
func (my MMS) Add(d time.Duration) MMS {
	t := my.ToTime().Add(d)
	return FromTime(t)
}

//Sub : mms - mms = time.Duration
func (my MMS) Sub(m MMS) time.Duration {
	v := my - m
	return time.Duration(v.Value() * DivMMS)
}

//Duration : time.Duration
func (my MMS) Duration() time.Duration {
	return time.Duration(my.Value() * DivMMS)
}

//CompareTo :
func (my MMS) CompareTo(m MMS) int {
	if my.Value() == m.Value() {
		return 0
	} else if my.Value() > m.Value() {
		return 1
	} else {
		return -1
	}
}

//SnapShot : //지난 시간
func (my MMS) SnapShot() time.Duration {
	return Now().Sub(my)
}

///////////////////////////////////////////////////
// Public function area...
///////////////////////////////////////////////////

type Locale int

const (
	LocaleUTC = Locale(0)
	LocaleKST = Locale(1)
)

func toDate(year int, month time.Month, day, hour, min, sec, nsec int, locale ...Locale) time.Time {
	dt := time.Date(year, month, day, hour, min, sec, nsec, time.UTC)
	if len(locale) > 0 && locale[0] == LocaleKST {
		dt = dt.Add(time.Hour * -9)
	}
	return dt
}

//FromTime :
func FromTime(nt time.Time) MMS {
	return MMS(nt.UnixNano() / DivMMS)
}

//FromNano :
func FromNano(nano int64) MMS {
	if len(dbg.Void(nano)) != 19 {
		dbg.Red("mms.FromNano value is not 19.")
		return MMS(nano)
	}
	return MMS(nano / DivMMS)
}

//FromUnix :
func FromUnix(unix int64) MMS {
	if len(dbg.Void(unix)) != 10 {
		dbg.Red("mms.FromUnix value is not 10.")
		return MMS(unix)
	}
	return MMS(unix * divMMSec)
}

//FromString :
func FromString(str string) MMS {
	str = strings.TrimSpace(str)
	// if str == "" {
	// 	str = "0"
	// }
	v, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		dbg.Red("mms.FromString str is no number")
	}
	return MMS(v)
}

/*
10000000
20210101
*/

//FromYMD : 20210101 100000000
func FromYMD(ymd int, locale ...Locale) MMS {
	if ymd < 19700101 {
		ymd = 19700101 //19700101
	}
	str := fmt.Sprintf("%v", ymd)
	year, _ := strconv.ParseInt(str[:4], 10, 32)
	month, _ := strconv.ParseInt(str[4:6], 10, 32)
	day, _ := strconv.ParseInt(str[6:8], 10, 32)
	hour := 0
	min := 0
	sec := 0
	nsec := 0
	if len(str) >= 10 {
		v, _ := strconv.ParseInt(str[8:10], 10, 32)
		hour = int(v)
	}
	if len(str) >= 12 {
		v, _ := strconv.ParseInt(str[10:12], 10, 32)
		min = int(v)
	}
	if len(str) >= 14 {
		v, _ := strconv.ParseInt(str[12:14], 10, 32)
		sec = int(v)
	}
	if len(str) > 14 && len(str) <= 17 {
		v, _ := strconv.ParseInt(str[14:17], 10, 32)
		nsec = int(v) * DivMMS
	}
	return FromTime(toDate(int(year), time.Month(month), int(day), hour, min, sec, nsec, locale...))
}

func FromYMD3(year, month, day int, locale ...Locale) MMS {
	return FromTime(toDate(int(year), time.Month(month), int(day), 0, 0, 0, 0, locale...))
}

//
func FromYMDHM(year, month, day, hour, min int, locale ...Locale) MMS {
	return FromTime(toDate(
		int(year),
		time.Month(month),
		int(day),
		hour,
		min,
		0, 0,
		locale...,
	))
}

//ToTime :
func ToTime(mms MMS) time.Time {
	return time.Unix(0, mms.Value(true)*DivMMS).UTC()
}

//YMD :
func YMD(t time.Time, hmsn ...int) int {
	return YMDN(FromTime(t), hmsn...)
}

//YMDN : mms
func YMDN(mms MMS, hmsn ...int) int {
	ymd := YMDString(mms, hmsn...)
	v, _ := strconv.ParseInt(ymd, 10, 64)
	return int(v)
}

//YMDString : YEAR/MONTH/DAY
func YMDString(mms MMS, hmsn ...int) string {
	unix := int64(mms / divMMSec)
	t := time.Unix(unix, 0).UTC()
	//t = uts(t)

	str := t.String()
	ss := strings.Split(str, " ")
	ymd := strings.Split(ss[0], "-") //2020-01-01
	//
	// _ = hms
	// _ = ss[2] // +0000

	var result string
	result = fmt.Sprintf("%v%v%v", ymd[0], ymd[1], ymd[2])

	if len(hmsn) > 0 {
		size := hmsn[0]
		hms := strings.Split(ss[1], ":") //00:00:00
		if size >= 1 {
			result += hms[0] //hour
		}
		if size >= 2 {
			result += hms[1] //min
		}
		if size >= 3 {
			result += hms[2] //sec
		}
	}

	return result
}

//YearMonthDay : year, month , day
func (my MMS) YearMonthDay() (int, int, int) {
	return YearMonthDay(my)
}

//YearMonthDay : year, month , day
func YearMonthDay(mms MMS) (int, int, int) {
	ymdn := YMDN(mms)
	year := ymdn / 10000

	mdn := ymdn % 10000
	month := mdn / 100

	day := mdn % 100

	return year, month, day
}

//SplitYMDHM : year, month, day, hour, min
func (my MMS) SplitYMDHM() (int, int, int, int, int) {
	t := my.ToTime()
	year := t.Year()
	month := int(t.Month())
	day := t.Day()
	hour := t.Hour()
	min := t.Minute()
	return year, month, day, hour, min
}
func (my MMS) SplitYMDHMArray() []int {
	sl := []int{}
	y, mo, d, h, m := my.SplitYMDHM()
	sl = append(sl, y, mo, d, h, m)
	return sl
}
func FromYMDHMArray(sl []int, locale ...Locale) MMS {
	if len(sl) < 5 {
		return 0
	}
	return FromYMDHM(
		sl[0],
		sl[1],
		sl[2],
		sl[3],
		sl[4],
		locale...,
	)
}

func ElapsedTime(work func(), skipView ...bool) time.Duration {
	nt := Now()
	if work != nil {
		work()
	}
	du := Now().Sub(nt)

	if len(skipView) > 0 && skipView[0] {

	} else {
		dbg.YellowBoldBG("elapsedTime :", du)
	}
	return du
}
