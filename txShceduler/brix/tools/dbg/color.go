package dbg

import (
	"fmt"
	"strings"

	"txscheduler/brix/tools/dbg/cc"
)

var (
	isColorLogForce = false
	logWriteln      func(v ...interface{})
)

func SetColorLogForce(b bool) {
	isColorLogForce = b
}

func SetLogWriteln(f func(v ...interface{})) {
	logWriteln = f
}

func _logWriteln(v ...interface{}) {
	if logWriteln != nil {
		sl := []interface{}{}
		for _, i := range v {
			switch i.(type) {
			case cc.COLOR:
			default:
				sl = append(sl, i)
			}
		} //for
		logWriteln(sl...)
	}
}
func CallLogWriteln(v ...interface{}) {
	_logWriteln(v...)
}

func printColor(color cc.COLOR, a ...interface{}) {
	colorPrintln(
		fmt.Println,
		color,
		a...,
	)
}
func colorPrintln(
	___println func(a ...interface{}) (n int, err error),
	color cc.COLOR,
	a ...interface{},
) {
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

				aa = append(aa, color, body, cc.END)
				isFormat = true
			}
		}
	}
	if !isFormat {
		aa = append(aa, color)
		aa = append(aa, a...)
		aa = append(aa, cc.END)
	}
	___println(aa...)
	_logWriteln(a...)
}

// Color :
func Color(color cc.COLOR, a ...interface{}) { printColor(color, a...) }

func Gray(a ...interface{}) { printColor(cc.Gray, a...) }

// Black :
func Black(a ...interface{})         { printColor(cc.Black, a...) }
func BlackBold(a ...interface{})     { printColor(cc.BlackBold, a...) }
func BlackBG(a ...interface{})       { printColor(cc.BlackBG, a...) }
func BlackBoldBG(a ...interface{})   { printColor(cc.BlackBoldBG, a...) }
func BlackItalic(a ...interface{})   { printColor(cc.BlackItalic, a...) }
func BlackItalicBG(a ...interface{}) { printColor(cc.BlackItalicBG, a...) }
func BlackUL(a ...interface{})       { printColor(cc.BlackUL, a...) }
func BlackULBG(a ...interface{})     { printColor(cc.BlackULBG, a...) }

// Red :
func Red(a ...interface{})         { printColor(cc.Red, a...) }
func RedBold(a ...interface{})     { printColor(cc.RedBold, a...) }
func RedBG(a ...interface{})       { printColor(cc.RedBG, a...) }
func RedBoldBG(a ...interface{})   { printColor(cc.RedBoldBG, a...) }
func RedItalic(a ...interface{})   { printColor(cc.RedItalic, a...) }
func RedItalicBG(a ...interface{}) { printColor(cc.RedItalicBG, a...) }
func RedUL(a ...interface{})       { printColor(cc.RedUL, a...) }
func RedULBG(a ...interface{})     { printColor(cc.RedULBG, a...) }

// Green :
func Green(a ...interface{})         { printColor(cc.Green, a...) }
func GreenBold(a ...interface{})     { printColor(cc.GreenBold, a...) }
func GreenBG(a ...interface{})       { printColor(cc.GreenBG, a...) }
func GreenBoldBG(a ...interface{})   { printColor(cc.GreenBoldBG, a...) }
func GreenItalic(a ...interface{})   { printColor(cc.GreenItalic, a...) }
func GreenItalicBG(a ...interface{}) { printColor(cc.GreenItalicBG, a...) }
func GreenUL(a ...interface{})       { printColor(cc.GreenUL, a...) }
func GreenULBG(a ...interface{})     { printColor(cc.GreenULBG, a...) }

// Yellow :
func Yellow(a ...interface{})         { printColor(cc.Yellow, a...) }
func YellowBold(a ...interface{})     { printColor(cc.YellowBold, a...) }
func YellowBG(a ...interface{})       { printColor(cc.YellowBG, a...) }
func YellowBoldBG(a ...interface{})   { printColor(cc.YellowBoldBG, a...) }
func YellowItalic(a ...interface{})   { printColor(cc.YellowItalic, a...) }
func YellowItalicBG(a ...interface{}) { printColor(cc.YellowItalicBG, a...) }
func YellowUL(a ...interface{})       { printColor(cc.YellowUL, a...) }
func YellowULBG(a ...interface{})     { printColor(cc.YellowULBG, a...) }

// Blue :
func Blue(a ...interface{})         { printColor(cc.Blue, a...) }
func BlueBold(a ...interface{})     { printColor(cc.BlueBold, a...) }
func BlueBG(a ...interface{})       { printColor(cc.BlueBG, a...) }
func BlueBoldBG(a ...interface{})   { printColor(cc.BlueBoldBG, a...) }
func BlueItalic(a ...interface{})   { printColor(cc.BlueItalic, a...) }
func BlueItalicBG(a ...interface{}) { printColor(cc.BlueItalicBG, a...) }
func BlueUL(a ...interface{})       { printColor(cc.BlueUL, a...) }
func BlueULBG(a ...interface{})     { printColor(cc.BlueULBG, a...) }

// Purple :
func Purple(a ...interface{})         { printColor(cc.Purple, a...) }
func PurpleBold(a ...interface{})     { printColor(cc.PurpleBold, a...) }
func PurpleBG(a ...interface{})       { printColor(cc.PurpleBG, a...) }
func PurpleBoldBG(a ...interface{})   { printColor(cc.PurpleBoldBG, a...) }
func PurpleItalic(a ...interface{})   { printColor(cc.PurpleItalic, a...) }
func PurpleItalicBG(a ...interface{}) { printColor(cc.PurpleItalicBG, a...) }
func PurpleUL(a ...interface{})       { printColor(cc.PurpleUL, a...) }
func PurpleULBG(a ...interface{})     { printColor(cc.PurpleULBG, a...) }

// Cyan :
func Cyan(a ...interface{})         { printColor(cc.Cyan, a...) }
func CyanBold(a ...interface{})     { printColor(cc.CyanBold, a...) }
func CyanBG(a ...interface{})       { printColor(cc.CyanBG, a...) }
func CyanBoldBG(a ...interface{})   { printColor(cc.CyanBoldBG, a...) }
func CyanItalic(a ...interface{})   { printColor(cc.CyanItalic, a...) }
func CyanItalicBG(a ...interface{}) { printColor(cc.CyanItalicBG, a...) }
func CyanUL(a ...interface{})       { printColor(cc.CyanUL, a...) }
func CyanULBG(a ...interface{})     { printColor(cc.CyanULBG, a...) }

// White :
func White(a ...interface{})         { printColor(cc.White, a...) }
func WhiteBold(a ...interface{})     { printColor(cc.WhiteBold, a...) }
func WhiteBG(a ...interface{})       { printColor(cc.WhiteBG, a...) }
func WhiteBoldBG(a ...interface{})   { printColor(cc.WhiteBoldBG, a...) }
func WhiteItalic(a ...interface{})   { printColor(cc.WhiteItalic, a...) }
func WhiteItalicBG(a ...interface{}) { printColor(cc.WhiteItalicBG, a...) }
func WhiteUL(a ...interface{})       { printColor(cc.WhiteUL, a...) }
func WhiteULBG(a ...interface{})     { printColor(cc.WhiteULBG, a...) }

// Log : WhiteBold
func Log(a ...interface{}) { printColor(cc.WhiteItalic, a...) }

// ///////////////////////////////////////////////////////////////////////////
type TLOG func(...interface{})

type colorLoger struct {
	isSkip bool
}

func (my colorLoger) call(color cc.COLOR, a ...interface{}) {
	if my.isSkip {
		return
	}
	colorPrintln(
		fmt.Println,
		color,
		a...,
	)
}

func (my colorLoger) Skip() bool { return my.isSkip }
func (my *colorLoger) SetSkip(f bool) {
	my.isSkip = f
}

func NewColorLoger(isSkip ...bool) IColorLoger {
	my := colorLoger{
		isSkip: IsTrue2(isSkip...),
	}
	return &my
}

// Black :
func (my colorLoger) Black(a ...interface{})         { my.call(cc.Black, a...) }
func (my colorLoger) BlackBold(a ...interface{})     { my.call(cc.BlackBold, a...) }
func (my colorLoger) BlackBG(a ...interface{})       { my.call(cc.BlackBG, a...) }
func (my colorLoger) BlackBoldBG(a ...interface{})   { my.call(cc.BlackBoldBG, a...) }
func (my colorLoger) BlackItalic(a ...interface{})   { my.call(cc.BlackItalic, a...) }
func (my colorLoger) BlackItalicBG(a ...interface{}) { my.call(cc.BlackItalicBG, a...) }
func (my colorLoger) BlackUL(a ...interface{})       { my.call(cc.BlackUL, a...) }
func (my colorLoger) BlackULBG(a ...interface{})     { my.call(cc.BlackULBG, a...) }

// Red :
func (my colorLoger) Red(a ...interface{})         { my.call(cc.Red, a...) }
func (my colorLoger) RedBold(a ...interface{})     { my.call(cc.RedBold, a...) }
func (my colorLoger) RedBG(a ...interface{})       { my.call(cc.RedBG, a...) }
func (my colorLoger) RedBoldBG(a ...interface{})   { my.call(cc.RedBoldBG, a...) }
func (my colorLoger) RedItalic(a ...interface{})   { my.call(cc.RedItalic, a...) }
func (my colorLoger) RedItalicBG(a ...interface{}) { my.call(cc.RedItalicBG, a...) }
func (my colorLoger) RedUL(a ...interface{})       { my.call(cc.RedUL, a...) }
func (my colorLoger) RedULBG(a ...interface{})     { my.call(cc.RedULBG, a...) }

// Green :
func (my colorLoger) Green(a ...interface{})         { my.call(cc.Green, a...) }
func (my colorLoger) GreenBold(a ...interface{})     { my.call(cc.GreenBold, a...) }
func (my colorLoger) GreenBG(a ...interface{})       { my.call(cc.GreenBG, a...) }
func (my colorLoger) GreenBoldBG(a ...interface{})   { my.call(cc.GreenBoldBG, a...) }
func (my colorLoger) GreenItalic(a ...interface{})   { my.call(cc.GreenItalic, a...) }
func (my colorLoger) GreenItalicBG(a ...interface{}) { my.call(cc.GreenItalicBG, a...) }
func (my colorLoger) GreenUL(a ...interface{})       { my.call(cc.GreenUL, a...) }
func (my colorLoger) GreenULBG(a ...interface{})     { my.call(cc.GreenULBG, a...) }

// Yellow :
func (my colorLoger) Yellow(a ...interface{})         { my.call(cc.Yellow, a...) }
func (my colorLoger) YellowBold(a ...interface{})     { my.call(cc.YellowBold, a...) }
func (my colorLoger) YellowBG(a ...interface{})       { my.call(cc.YellowBG, a...) }
func (my colorLoger) YellowBoldBG(a ...interface{})   { my.call(cc.YellowBoldBG, a...) }
func (my colorLoger) YellowItalic(a ...interface{})   { my.call(cc.YellowItalic, a...) }
func (my colorLoger) YellowItalicBG(a ...interface{}) { my.call(cc.YellowItalicBG, a...) }
func (my colorLoger) YellowUL(a ...interface{})       { my.call(cc.YellowUL, a...) }
func (my colorLoger) YellowULBG(a ...interface{})     { my.call(cc.YellowULBG, a...) }

// Blue :
func (my colorLoger) Blue(a ...interface{})         { my.call(cc.Blue, a...) }
func (my colorLoger) BlueBold(a ...interface{})     { my.call(cc.BlueBold, a...) }
func (my colorLoger) BlueBG(a ...interface{})       { my.call(cc.BlueBG, a...) }
func (my colorLoger) BlueBoldBG(a ...interface{})   { my.call(cc.BlueBoldBG, a...) }
func (my colorLoger) BlueItalic(a ...interface{})   { my.call(cc.BlueItalic, a...) }
func (my colorLoger) BlueItalicBG(a ...interface{}) { my.call(cc.BlueItalicBG, a...) }
func (my colorLoger) BlueUL(a ...interface{})       { my.call(cc.BlueUL, a...) }
func (my colorLoger) BlueULBG(a ...interface{})     { my.call(cc.BlueULBG, a...) }

// Purple :
func (my colorLoger) Purple(a ...interface{})         { my.call(cc.Purple, a...) }
func (my colorLoger) PurpleBold(a ...interface{})     { my.call(cc.PurpleBold, a...) }
func (my colorLoger) PurpleBG(a ...interface{})       { my.call(cc.PurpleBG, a...) }
func (my colorLoger) PurpleBoldBG(a ...interface{})   { my.call(cc.PurpleBoldBG, a...) }
func (my colorLoger) PurpleItalic(a ...interface{})   { my.call(cc.PurpleItalic, a...) }
func (my colorLoger) PurpleItalicBG(a ...interface{}) { my.call(cc.PurpleItalicBG, a...) }
func (my colorLoger) PurpleUL(a ...interface{})       { my.call(cc.PurpleUL, a...) }
func (my colorLoger) PurpleULBG(a ...interface{})     { my.call(cc.PurpleULBG, a...) }

// Cyan :
func (my colorLoger) Cyan(a ...interface{})         { my.call(cc.Cyan, a...) }
func (my colorLoger) CyanBold(a ...interface{})     { my.call(cc.CyanBold, a...) }
func (my colorLoger) CyanBG(a ...interface{})       { my.call(cc.CyanBG, a...) }
func (my colorLoger) CyanBoldBG(a ...interface{})   { my.call(cc.CyanBoldBG, a...) }
func (my colorLoger) CyanItalic(a ...interface{})   { my.call(cc.CyanItalic, a...) }
func (my colorLoger) CyanItalicBG(a ...interface{}) { my.call(cc.CyanItalicBG, a...) }
func (my colorLoger) CyanUL(a ...interface{})       { my.call(cc.CyanUL, a...) }
func (my colorLoger) CyanULBG(a ...interface{})     { my.call(cc.CyanULBG, a...) }

// White :
func (my colorLoger) White(a ...interface{})         { my.call(cc.White, a...) }
func (my colorLoger) WhiteBold(a ...interface{})     { my.call(cc.WhiteBold, a...) }
func (my colorLoger) WhiteBG(a ...interface{})       { my.call(cc.WhiteBG, a...) }
func (my colorLoger) WhiteBoldBG(a ...interface{})   { my.call(cc.WhiteBoldBG, a...) }
func (my colorLoger) WhiteItalic(a ...interface{})   { my.call(cc.WhiteItalic, a...) }
func (my colorLoger) WhiteItalicBG(a ...interface{}) { my.call(cc.WhiteItalicBG, a...) }
func (my colorLoger) WhiteUL(a ...interface{})       { my.call(cc.WhiteUL, a...) }
func (my colorLoger) WhiteULBG(a ...interface{})     { my.call(cc.WhiteULBG, a...) }

type IColorLoger interface {
	Skip() bool
	SetSkip(f bool)

	Black(a ...interface{})
	BlackBold(a ...interface{})
	BlackBG(a ...interface{})
	BlackBoldBG(a ...interface{})
	BlackItalic(a ...interface{})
	BlackItalicBG(a ...interface{})
	BlackUL(a ...interface{})
	BlackULBG(a ...interface{})
	Red(a ...interface{})
	RedBold(a ...interface{})
	RedBG(a ...interface{})
	RedBoldBG(a ...interface{})
	RedItalic(a ...interface{})
	RedItalicBG(a ...interface{})
	RedUL(a ...interface{})
	RedULBG(a ...interface{})
	Green(a ...interface{})
	GreenBold(a ...interface{})
	GreenBG(a ...interface{})
	GreenBoldBG(a ...interface{})
	GreenItalic(a ...interface{})
	GreenItalicBG(a ...interface{})
	GreenUL(a ...interface{})
	GreenULBG(a ...interface{})
	Yellow(a ...interface{})
	YellowBold(a ...interface{})
	YellowBG(a ...interface{})
	YellowBoldBG(a ...interface{})
	YellowItalic(a ...interface{})
	YellowItalicBG(a ...interface{})
	YellowUL(a ...interface{})
	YellowULBG(a ...interface{})
	Blue(a ...interface{})
	BlueBold(a ...interface{})
	BlueBG(a ...interface{})
	BlueBoldBG(a ...interface{})
	BlueItalic(a ...interface{})
	BlueItalicBG(a ...interface{})
	BlueUL(a ...interface{})
	BlueULBG(a ...interface{})
	Purple(a ...interface{})
	PurpleBold(a ...interface{})
	PurpleBG(a ...interface{})
	PurpleBoldBG(a ...interface{})
	PurpleItalic(a ...interface{})
	PurpleItalicBG(a ...interface{})
	PurpleUL(a ...interface{})
	PurpleULBG(a ...interface{})
	Cyan(a ...interface{})
	CyanBold(a ...interface{})
	CyanBG(a ...interface{})
	CyanBoldBG(a ...interface{})
	CyanItalic(a ...interface{})
	CyanItalicBG(a ...interface{})
	CyanUL(a ...interface{})
	CyanULBG(a ...interface{})
	White(a ...interface{})
	WhiteBold(a ...interface{})
	WhiteBG(a ...interface{})
	WhiteBoldBG(a ...interface{})
	WhiteItalic(a ...interface{})
	WhiteItalicBG(a ...interface{})
	WhiteUL(a ...interface{})
	WhiteULBG(a ...interface{})
}
