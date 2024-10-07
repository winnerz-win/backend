package unix

import (
	"fmt"
	"jtools/mms"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Time int64

const (
	ZERO    = Time(0)
	INFNITY = Time(99999999999) //5138-11-16T09:46:39Z

	Nano = 1000000000

	Second = Time(1)
	Minute = Time(Second * 60)
	Hour   = Time(Minute * 60)
	Day    = Time(Hour * 24)

	Month = Time(Day * 30)
	Year  = Time(Day * 365)
)

func Now() Time {
	return Time(time.Now().UTC().Unix())
}

func (my Time) MMS() mms.MMS {
	return mms.FromUnix(my.Int64())
}

func FromString(s string) Time {
	v, _ := strconv.ParseInt(s, 10, 64)
	return Time(v)
}

func FromMMS(v mms.MMS) Time {
	return Time(v.Unix())
}

// func (my Time) MMS() unix.Time {
// 	return unix.Time(my.Int64() * mms.SecMMS)
// }

func FromTime(v time.Time) Time {
	return Time(v.UTC().Unix())
}
func FromYMD(ymd int) Time {
	year := ymd / 10000
	val := ymd % 10000
	month := val / 100
	day := val % 100
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return FromTime(t)

}

// FromRFC3339 : 2006-01-02T15:04:05Z07:00
func FromRFC3339(text string) Time {
	v, _ := time.ParseInLocation(time.RFC3339, text, time.UTC)
	return FromTime(v)
}

func FromSplit(year, month, day, hour, min, sec int) Time {
	t := time.Date(year, time.Month(month), day, hour, min, sec, 0, time.UTC)
	return FromTime(t)
}

func (my Time) Int64() int64 { return int64(my) }

func (my Time) String() string {
	return my.Time().Format(time.RFC3339)
}

func (my Time) Int64String() string { return fmt.Sprint(my.Int64()) }

func (my Time) Time() time.Time {
	return time.Unix(my.Int64(), 0).UTC()
}

func (my Time) DateTime() primitive.DateTime {
	return primitive.NewDateTimeFromTime(my.Time())
}
func (my Time) DateTimeString() string {
	dt := primitive.NewDateTimeFromTime(my.Time())
	b, _ := dt.MarshalJSON()
	b = b[1:]
	b = b[:len(b)-1]
	return string(b)
}

func FromDateTime(dt primitive.DateTime) Time {
	return FromTime(dt.Time())
}

func (my Time) Duration() time.Duration {
	return time.Duration(my.Int64() * Nano)
}

func (my Time) ElipsedDuration() time.Duration {
	return Now().Sub(my)
}

func (my Time) YMD() int {
	year, month, day, _, _, _ := my.Split()
	return year*10000 + month*100 + day
}

func (my Time) Split() (year, month, day, hour, min, sec int) {

	t := my.Time()
	year = t.Year()
	month = int(t.Month())
	day = t.Day()
	hour = t.Hour()
	min = t.Minute()
	sec = t.Second()
	return
}

func (my Time) Sub(v Time) time.Duration {
	return (my - v).Duration()
}

func (my Time) Add(d time.Duration) Time {
	return my + Time(d/Nano)
}

// KST : 2001-01-01T00:00:00Z
func (my Time) KST() string {
	str := my.Time().UTC().Add(time.Hour * 9).String()
	spt := strings.Split(str, " ")
	sl := []string{spt[0], "T", spt[1], "KST"}
	return strings.Join(sl, "")
}

// ToZero : 00:00
func (my Time) ToZero() Time {
	return FromYMD(my.YMD())

	// nt := Now().ToTime()
	// nt = time.Date(nt.Year(), nt.Month(), nt.Day(), 0, 0, 0, 0, time.UTC)
	// return FromTime(nt)
}
