package dbg

import "fmt"

const (
	ENTER = "\n"
)

func Println(a ...interface{}) {
	fmt.Println(a...)
}

func Printf(foramt string, a ...interface{}) {
	fmt.Printf(foramt, a...)
}

func Print(a ...interface{}) {
	fmt.Print(a...)
}
