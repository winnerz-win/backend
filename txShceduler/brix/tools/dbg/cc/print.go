package cc

import (
	"fmt"
	"strings"
)

//PrintRed :
func PrintRed(a ...interface{}) {
	printColor(Red, a...)
}

//Green :
func PrintGreen(a ...interface{}) {
	printColor(Green, a...)
}

//PrintYellow :
func PrintYellow(a ...interface{}) {
	printColor(Yellow, a...)
}

//PrintBlue :
func PrintBlue(a ...interface{}) {
	printColor(Blue, a...)
}

//PrintPurple :
func PrintPurple(a ...interface{}) {
	printColor(Purple, a...)
}

//PrintCyan :
func PrintCyan(a ...interface{}) {
	printColor(Cyan, a...)
}

func printColor(color COLOR, a ...interface{}) {
	// if runtime.GOOS != "windows" {
	// 	fmt.Println(a...)
	// 	return
	// }
	aa := []interface{}{}

	isFormat := false
	if len(a) > 1 {
		v, do := a[0].(string)
		if do {
			if strings.Contains(v, "%v") {
				vp := []interface{}{}
				for i, p := range a {
					if i != 0 {
						vp = append(vp, p)
					}
				} //for
				body := fmt.Sprintf(v, vp...)

				aa = append(aa, color, body, END)
				isFormat = true
			}
		}
	}

	if isFormat == false {
		aa = append(aa, color)
		aa = append(aa, a...)
		aa = append(aa, END)
	}
	fmt.Println(aa...)

}
