package cc

import (
	"fmt"
	"strings"
)

func Println(a ...interface{}) { fmt.Println(a...) }

func Black(a ...interface{})  { printColor(D_Black, a...) }
func Red(a ...interface{})    { printColor(D_Red, a...) }
func Green(a ...interface{})  { printColor(D_Green, a...) }
func Yellow(a ...interface{}) { printColor(D_Yellow, a...) }
func Blue(a ...interface{})   { printColor(D_Blue, a...) }
func Purple(a ...interface{}) { printColor(D_Purple, a...) }
func Cyan(a ...interface{})   { printColor(D_Cyan, a...) }
func White(a ...interface{})  { printColor(D_White, a...) }

func Gray(a ...interface{}) { printColor(D_Gray, a...) }

func BlackBold(a ...interface{})  { printColor(D_BlackBold, a...) }
func RedBold(a ...interface{})    { printColor(D_RedBold, a...) }
func GreenBold(a ...interface{})  { printColor(D_GreenBold, a...) }
func YellowBold(a ...interface{}) { printColor(D_YellowBold, a...) }
func BlueBold(a ...interface{})   { printColor(D_BlueBold, a...) }
func PurpleBold(a ...interface{}) { printColor(D_PurpleBold, a...) }
func CyanBold(a ...interface{})   { printColor(D_CyanBold, a...) }
func WhiteBold(a ...interface{})  { printColor(D_WhiteBold, a...) }

func BlackBG(a ...interface{})  { printColor(D_BlackBG, a...) }
func RedBG(a ...interface{})    { printColor(D_RedBG, a...) }
func GreenBG(a ...interface{})  { printColor(D_GreenBG, a...) }
func YellowBG(a ...interface{}) { printColor(D_YellowBG, a...) }
func BlueBG(a ...interface{})   { printColor(D_BlueBG, a...) }
func PurpleBG(a ...interface{}) { printColor(D_PurpleBG, a...) }
func CyanBG(a ...interface{})   { printColor(D_CyanBG, a...) }
func WhiteBG(a ...interface{})  { printColor(D_WhiteBG, a...) }

func BlackBoldBG(a ...interface{})  { printColor(D_BlackBoldBG, a...) }
func RedBoldBG(a ...interface{})    { printColor(D_RedBoldBG, a...) }
func GreenBoldBG(a ...interface{})  { printColor(D_GreenBoldBG, a...) }
func YellowBoldBG(a ...interface{}) { printColor(D_YellowBoldBG, a...) }
func BlueBoldBG(a ...interface{})   { printColor(D_BlueBoldBG, a...) }
func PurpleBoldBG(a ...interface{}) { printColor(D_PurpleBoldBG, a...) }
func CyanBoldBG(a ...interface{})   { printColor(D_CyanBoldBG, a...) }
func WhiteBoldBG(a ...interface{})  { printColor(D_WhiteBoldBG, a...) }

func BlackItalic(a ...interface{})  { printColor(D_BlackItalic, a...) }
func RedItalic(a ...interface{})    { printColor(D_RedItalic, a...) }
func GreenItalic(a ...interface{})  { printColor(D_GreenItalic, a...) }
func YellowItalic(a ...interface{}) { printColor(D_YellowItalic, a...) }
func BlueItalic(a ...interface{})   { printColor(D_BlueItalic, a...) }
func PurpleItalic(a ...interface{}) { printColor(D_PurpleItalic, a...) }
func CyanItalic(a ...interface{})   { printColor(D_CyanItalic, a...) }
func WhiteItalic(a ...interface{})  { printColor(D_WhiteItalic, a...) }

func BlackItalicBG(a ...interface{})  { printColor(D_BlackItalicBG, a...) }
func RedItalicBG(a ...interface{})    { printColor(D_RedItalicBG, a...) }
func GreenItalicBG(a ...interface{})  { printColor(D_GreenItalicBG, a...) }
func YellowItalicBG(a ...interface{}) { printColor(D_YellowItalicBG, a...) }
func BlueItalicBG(a ...interface{})   { printColor(D_BlueItalicBG, a...) }
func PurpleItalicBG(a ...interface{}) { printColor(D_PurpleItalicBG, a...) }
func CyanItalicBG(a ...interface{})   { printColor(D_CyanItalicBG, a...) }
func WhiteItalicBG(a ...interface{})  { printColor(D_WhiteItalicBG, a...) }

func BlackUL(a ...interface{})  { printColor(D_BlackUL, a...) }
func RedUL(a ...interface{})    { printColor(D_RedUL, a...) }
func GreenUL(a ...interface{})  { printColor(D_GreenUL, a...) }
func YellowUL(a ...interface{}) { printColor(D_YellowUL, a...) }
func BlueUL(a ...interface{})   { printColor(D_BlueUL, a...) }
func PurpleUL(a ...interface{}) { printColor(D_PurpleUL, a...) }
func CyanUL(a ...interface{})   { printColor(D_CyanUL, a...) }
func WhiteUL(a ...interface{})  { printColor(D_WhiteUL, a...) }

func BlackULBG(a ...interface{})  { printColor(D_BlackULBG, a...) }
func RedULBG(a ...interface{})    { printColor(D_RedULBG, a...) }
func GreenULBG(a ...interface{})  { printColor(D_GreenULBG, a...) }
func YellowULBG(a ...interface{}) { printColor(D_YellowULBG, a...) }
func BlueULBG(a ...interface{})   { printColor(D_BlueULBG, a...) }
func PurpleULBG(a ...interface{}) { printColor(D_PurpleULBG, a...) }
func CyanULBG(a ...interface{})   { printColor(D_CyanULBG, a...) }
func WhiteULBG(a ...interface{})  { printColor(D_WhiteULBG, a...) }

//////////////////////////////////////////////

var (
	file_log_writer func(v ...interface{}) = nil
)

func SetFileLogWriter(f func(v ...interface{})) {
	file_log_writer = f
}

func printColor(color string, a ...interface{}) {

	if file_log_writer != nil {
		file_log_writer(a...)
	}

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

				aa = append(aa, color, body, D_End)
				isFormat = true
			}
		}
	}

	if !isFormat {
		aa = append(aa, color)
		aa = append(aa, a...)
		aa = append(aa, D_End)
	}
	fmt.Println(aa...)

}
