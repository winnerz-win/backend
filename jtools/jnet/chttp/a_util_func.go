package chttp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func NowTime() time.Time {
	return time.Now().UTC()
}

func YMDString(nt time.Time) string {
	_2p := func(v int) string {
		if v < 10 {
			return "0" + fmt.Sprint(v)
		}
		return fmt.Sprint(v)
	}
	y := nt.Year()
	m := int(nt.Month())
	d := nt.Day()

	return Cat(y, _2p(m), _2p(d))
}

func Cat(a ...interface{}) string {
	sl := []string{}
	for _, v := range a {
		sl = append(sl, fmt.Sprint(v))
	}
	return strings.Join(sl, "")
}

func LogColor(color string, a ...interface{}) {
	sl := []interface{}{
		color,
	}
	sl = append(sl, a...)
	sl = append(sl, "\033[0m")
	fmt.Println(sl...)
}

func LogError(a ...interface{})    { LogColor("\033[3;31m", a...) } //Red
func LogYellow(a ...interface{})   { LogColor("\033[3;33m", a...) }
func LogYellowBG(a ...interface{}) { LogColor("\033[3;43m", a...) }
func LogPurple(a ...interface{})   { LogColor("\033[3;35m", a...) }
func LogWhite(a ...interface{})    { LogColor("\033[3;37m", a...) }
func LogGray(a ...interface{})     { LogColor("\033[2;37m", a...) }

func NowPath() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

func ToJsonString(v interface{}) string {
	if v == nil {
		return ""
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil || json.Valid(b) {
		fmt.Sprintln(string(b))
	}
	return fmt.Sprint(v)
}

func IsTrue(p interface{}) bool {
	switch v := p.(type) {
	case bool:
		return v
	case []bool:
		if len(v) > 0 {
			return v[0]
		}
	case string:
		return strings.ToLower(strings.TrimSpace(v)) == "true"

	case []interface{}:
		if len(v) > 0 {
			return IsTrue(v[0])
		}
	}
	return false
}
