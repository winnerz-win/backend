package cc

import (
	"fmt"
	"strings"
)

func Println(a ...interface{}) { fmt.Println(a...) }

func Black(a ...interface{})  { printColor(black, a...) }
func Red(a ...interface{})    { printColor(red, a...) }
func Green(a ...interface{})  { printColor(green, a...) }
func Yellow(a ...interface{}) { printColor(yellow, a...) }
func Blue(a ...interface{})   { printColor(blue, a...) }
func Purple(a ...interface{}) { printColor(purple, a...) }
func Cyan(a ...interface{})   { printColor(cyan, a...) }
func White(a ...interface{})  { printColor(white, a...) }

func Gray(a ...interface{}) { printColor(gray, a...) }

func BlackBold(a ...interface{})  { printColor(blackBold, a...) }
func RedBold(a ...interface{})    { printColor(redBold, a...) }
func GreenBold(a ...interface{})  { printColor(greenBold, a...) }
func YellowBold(a ...interface{}) { printColor(yellowBold, a...) }
func BlueBold(a ...interface{})   { printColor(blueBold, a...) }
func PurpleBold(a ...interface{}) { printColor(purpleBold, a...) }
func CyanBold(a ...interface{})   { printColor(cyanBold, a...) }
func WhiteBold(a ...interface{})  { printColor(whiteBold, a...) }

func BlackBG(a ...interface{})  { printColor(blackBG, a...) }
func RedBG(a ...interface{})    { printColor(redBG, a...) }
func GreenBG(a ...interface{})  { printColor(greenBG, a...) }
func YellowBG(a ...interface{}) { printColor(yellowBG, a...) }
func BlueBG(a ...interface{})   { printColor(blueBG, a...) }
func PurpleBG(a ...interface{}) { printColor(purpleBG, a...) }
func CyanBG(a ...interface{})   { printColor(cyanBG, a...) }
func WhiteBG(a ...interface{})  { printColor(whiteBG, a...) }

func BlackBoldBG(a ...interface{})  { printColor(blackBoldBG, a...) }
func RedBoldBG(a ...interface{})    { printColor(redBoldBG, a...) }
func GreenBoldBG(a ...interface{})  { printColor(greenBoldBG, a...) }
func YellowBoldBG(a ...interface{}) { printColor(yellowBoldBG, a...) }
func BlueBoldBG(a ...interface{})   { printColor(blueBoldBG, a...) }
func PurpleBoldBG(a ...interface{}) { printColor(purpleBoldBG, a...) }
func CyanBoldBG(a ...interface{})   { printColor(cyanBoldBG, a...) }
func WhiteBoldBG(a ...interface{})  { printColor(whiteBoldBG, a...) }

func BlackItalic(a ...interface{})  { printColor(blackItalic, a...) }
func RedItalic(a ...interface{})    { printColor(redItalic, a...) }
func GreenItalic(a ...interface{})  { printColor(greenItalic, a...) }
func YellowItalic(a ...interface{}) { printColor(yellowItalic, a...) }
func BlueItalic(a ...interface{})   { printColor(blueItalic, a...) }
func PurpleItalic(a ...interface{}) { printColor(purpleItalic, a...) }
func CyanItalic(a ...interface{})   { printColor(cyanItalic, a...) }
func WhiteItalic(a ...interface{})  { printColor(whiteItalic, a...) }

func BlackItalicBG(a ...interface{})  { printColor(blackItalicBG, a...) }
func RedItalicBG(a ...interface{})    { printColor(redItalicBG, a...) }
func GreenItalicBG(a ...interface{})  { printColor(greenItalicBG, a...) }
func YellowItalicBG(a ...interface{}) { printColor(yellowItalicBG, a...) }
func BlueItalicBG(a ...interface{})   { printColor(blueItalicBG, a...) }
func PurpleItalicBG(a ...interface{}) { printColor(purpleItalicBG, a...) }
func CyanItalicBG(a ...interface{})   { printColor(cyanItalicBG, a...) }
func WhiteItalicBG(a ...interface{})  { printColor(whiteItalicBG, a...) }

func BlackUL(a ...interface{})  { printColor(blackUL, a...) }
func RedUL(a ...interface{})    { printColor(redUL, a...) }
func GreenUL(a ...interface{})  { printColor(greenUL, a...) }
func YellowUL(a ...interface{}) { printColor(yellowUL, a...) }
func BlueUL(a ...interface{})   { printColor(blueUL, a...) }
func PurpleUL(a ...interface{}) { printColor(purpleUL, a...) }
func CyanUL(a ...interface{})   { printColor(cyanUL, a...) }
func WhiteUL(a ...interface{})  { printColor(whiteUL, a...) }

func BlackULBG(a ...interface{})  { printColor(blackULBG, a...) }
func RedULBG(a ...interface{})    { printColor(redULBG, a...) }
func GreenULBG(a ...interface{})  { printColor(greenULBG, a...) }
func YellowULBG(a ...interface{}) { printColor(yellowULBG, a...) }
func BlueULBG(a ...interface{})   { printColor(blueULBG, a...) }
func PurpleULBG(a ...interface{}) { printColor(purpleULBG, a...) }
func CyanULBG(a ...interface{})   { printColor(cyanULBG, a...) }
func WhiteULBG(a ...interface{})  { printColor(whiteULBG, a...) }

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

				aa = append(aa, color, body, end)
				isFormat = true
			}
		}
	}

	if !isFormat {
		aa = append(aa, color)
		aa = append(aa, a...)
		aa = append(aa, end)
	}
	fmt.Println(aa...)

}
