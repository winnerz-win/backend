package dbg

import "time"

func SleepSecond() {
	time.Sleep(time.Second)
}
func Sleep500() {
	time.Sleep(time.Millisecond * 500)
}
func Sleep250() {
	time.Sleep(time.Millisecond * 250)
}
func Sleep100() {
	time.Sleep(time.Millisecond * 100)
}
func Sleep50() {
	time.Sleep(time.Millisecond * 50)
}
func Sleep10() {
	time.Sleep(time.Millisecond * 10)
}
