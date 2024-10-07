package spaceup

import "jtools/cc"

//////////////////////////////////////////////////////

type Color func(a ...interface{})

var (
	Red     = cc.Red
	RedBold = cc.RedBold
	RedBG   = cc.RedBG

	Blue     = cc.Blue
	BlueBold = cc.BlueBold
	BlueBG   = cc.BlueBG

	Green     = cc.Green
	GreenBold = cc.GreenBold
	GreenBG   = cc.GreenBG

	Yellow     = cc.Yellow
	YellowBold = cc.YellowBold
	YellowBG   = cc.YellowBG

	Purple     = cc.Purple
	PurpleBold = cc.PurpleBold
	PurpleBG   = cc.PurpleBG

	Cyan     = cc.Cyan
	CyanBold = cc.CyanBold
	CyanBG   = cc.CyanBG

	White     = cc.White
	WhiteBold = cc.WhiteBold
	WhiteBG   = cc.WhiteBG

	Gray = cc.Gray
)

func (my _table) ViewConsole(title_color Color, body_color Color, is_sort ...bool) {
	for i, v := range my.List(is_sort...) {
		if i == 0 {
			title_color(v)
		} else {
			body_color(v)
		}
	} //for
}
