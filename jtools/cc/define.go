package cc

const (
	//color = "\033[(normal/bold);30m"
	D_End = "\033[0m"

	D_Black  = "\033[0;30m"
	D_Red    = "\033[0;31m"
	D_Green  = "\033[0;32m"
	D_Yellow = "\033[0;33m"
	D_Blue   = "\033[0;34m"
	D_Purple = "\033[0;35m"
	D_Cyan   = "\033[0;36m"
	D_White  = "\033[0;37m"

	D_Gray = "\033[2;37m"

	D_BlackBold  = "\033[1;30m"
	D_RedBold    = "\033[1;31m"
	D_GreenBold  = "\033[1;32m"
	D_YellowBold = "\033[1;33m"
	D_BlueBold   = "\033[1;34m"
	D_PurpleBold = "\033[1;35m"
	D_CyanBold   = "\033[1;36m"
	D_WhiteBold  = "\033[1;37m"

	D_BlackBG  = "\033[0;40m"
	D_RedBG    = "\033[0;41m"
	D_GreenBG  = "\033[0;42m"
	D_YellowBG = "\033[0;43m"
	D_BlueBG   = "\033[0;44m"
	D_PurpleBG = "\033[0;45m"
	D_CyanBG   = "\033[0;46m"
	D_WhiteBG  = "\033[0;37m"

	D_BlackBoldBG  = "\033[1;40m"
	D_RedBoldBG    = "\033[1;41m"
	D_GreenBoldBG  = "\033[1;42m"
	D_YellowBoldBG = "\033[1;43m"
	D_BlueBoldBG   = "\033[1;44m"
	D_PurpleBoldBG = "\033[1;45m"
	D_CyanBoldBG   = "\033[1;46m"
	D_WhiteBoldBG  = "\033[1;47m"

	D_BlackItalic  = "\033[3;30m"
	D_RedItalic    = "\033[3;31m"
	D_GreenItalic  = "\033[3;32m"
	D_YellowItalic = "\033[3;33m"
	D_BlueItalic   = "\033[3;34m"
	D_PurpleItalic = "\033[3;35m"
	D_CyanItalic   = "\033[3;36m"
	D_WhiteItalic  = "\033[3;37m"

	D_BlackItalicBG  = "\033[3;30m"
	D_RedItalicBG    = "\033[3;31m"
	D_GreenItalicBG  = "\033[3;32m"
	D_YellowItalicBG = "\033[3;33m"
	D_BlueItalicBG   = "\033[3;34m"
	D_PurpleItalicBG = "\033[3;35m"
	D_CyanItalicBG   = "\033[3;36m"
	D_WhiteItalicBG  = "\033[3;37m"

	D_BlackUL  = "\033[4;30m"
	D_RedUL    = "\033[4;31m"
	D_GreenUL  = "\033[4;32m"
	D_YellowUL = "\033[4;33m"
	D_BlueUL   = "\033[4;34m"
	D_PurpleUL = "\033[4;35m"
	D_CyanUL   = "\033[4;36m"
	D_WhiteUL  = "\033[4;37m"

	D_BlackULBG  = "\033[4;40m"
	D_RedULBG    = "\033[4;41m"
	D_GreenULBG  = "\033[4;42m"
	D_YellowULBG = "\033[4;43m"
	D_BlueULBG   = "\033[4;44m"
	D_PurpleULBG = "\033[4;45m"
	D_CyanULBG   = "\033[4;46m"
	D_WhiteULBG  = "\033[4;47m"
)

///////////////////////////////////////////////////////////

type iColor struct {
	skip bool
}

func New(skips ...bool) iColor {
	skip := false
	if len(skips) > 0 {
		skip = skips[0]
	}
	return iColor{
		skip: skip,
	}
}
func (my *iColor) Skip(skips ...bool) {
	if len(skips) > 0 {
		my.skip = skips[0]
	} else {
		my.skip = true
	}
}

///////////////////////////////////////////////////////////

func (my iColor) Black(a ...interface{}) {
	if !my.skip {
		printColor(D_Black, a...)
	}
}
func (my iColor) Red(a ...interface{}) {
	if !my.skip {
		printColor(D_Red, a...)
	}
}
func (my iColor) Green(a ...interface{}) {
	if !my.skip {
		printColor(D_Green, a...)
	}
}

func (my iColor) Yellow(a ...interface{}) {
	if !my.skip {
		printColor(D_Yellow, a...)
	}
}
func (my iColor) Blue(a ...interface{}) {
	if !my.skip {
		printColor(D_Blue, a...)
	}
}
func (my iColor) Purple(a ...interface{}) {
	if !my.skip {
		printColor(D_Purple, a...)
	}
}
func (my iColor) Cyan(a ...interface{}) {
	if !my.skip {
		printColor(D_Cyan, a...)
	}
}
func (my iColor) White(a ...interface{}) {
	if !my.skip {
		printColor(D_White, a...)
	}
}

func (my iColor) Gray(a ...interface{}) {
	if !my.skip {
		printColor(D_Gray, a...)
	}
}

func (my iColor) BlackBold(a ...interface{}) {
	if !my.skip {
		printColor(D_BlackBold, a...)
	}
}
func (my iColor) RedBold(a ...interface{}) {
	if !my.skip {
		printColor(D_RedBold, a...)
	}
}
func (my iColor) GreenBold(a ...interface{}) {
	if !my.skip {
		printColor(D_GreenBold, a...)
	}
}
func (my iColor) YellowBold(a ...interface{}) {
	if !my.skip {
		printColor(D_YellowBold, a...)
	}
}
func (my iColor) BlueBold(a ...interface{}) {
	if !my.skip {
		printColor(D_BlueBold, a...)
	}
}
func (my iColor) PurpleBold(a ...interface{}) {
	if !my.skip {
		printColor(D_PurpleBold, a...)
	}
}
func (my iColor) CyanBold(a ...interface{}) {
	if !my.skip {
		printColor(D_CyanBold, a...)
	}
}
func (my iColor) WhiteBold(a ...interface{}) {
	if !my.skip {
		printColor(D_WhiteBold, a...)
	}
}

func (my iColor) BlackBG(a ...interface{}) {
	if !my.skip {
		printColor(D_BlackBG, a...)
	}
}
func (my iColor) RedBG(a ...interface{}) {
	if !my.skip {
		printColor(D_RedBG, a...)
	}
}
func (my iColor) GreenBG(a ...interface{}) {
	if !my.skip {
		printColor(D_GreenBG, a...)
	}
}
func (my iColor) YellowBG(a ...interface{}) {
	if !my.skip {
		printColor(D_YellowBG, a...)
	}
}
func (my iColor) BlueBG(a ...interface{}) {
	if !my.skip {
		printColor(D_BlueBG, a...)
	}
}
func (my iColor) PurpleBG(a ...interface{}) {
	if !my.skip {
		printColor(D_PurpleBG, a...)
	}
}
func (my iColor) CyanBG(a ...interface{}) {
	if !my.skip {
		printColor(D_CyanBG, a...)
	}
}
func (my iColor) WhiteBG(a ...interface{}) {
	if !my.skip {
		printColor(D_WhiteBG, a...)
	}
}

func (my iColor) BlackBoldBG(a ...interface{}) {
	if !my.skip {
		printColor(D_BlackBoldBG, a...)
	}
}
func (my iColor) RedBoldBG(a ...interface{}) {
	if !my.skip {
		printColor(D_RedBoldBG, a...)
	}
}
func (my iColor) GreenBoldBG(a ...interface{}) {
	if !my.skip {
		printColor(D_GreenBoldBG, a...)
	}
}
func (my iColor) YellowBoldBG(a ...interface{}) {
	if !my.skip {
		printColor(D_YellowBoldBG, a...)
	}
}
func (my iColor) BlueBoldBG(a ...interface{}) {
	if !my.skip {
		printColor(D_BlueBoldBG, a...)
	}
}
func (my iColor) PurpleBoldBG(a ...interface{}) {
	if !my.skip {
		printColor(D_PurpleBoldBG, a...)
	}
}
func (my iColor) CyanBoldBG(a ...interface{}) {
	if !my.skip {
		printColor(D_CyanBoldBG, a...)
	}
}
func (my iColor) WhiteBoldBG(a ...interface{}) {
	if !my.skip {
		printColor(D_WhiteBoldBG, a...)
	}
}

func (my iColor) BlackItalic(a ...interface{}) {
	if !my.skip {
		printColor(D_BlackItalic, a...)
	}
}
func (my iColor) RedItalic(a ...interface{}) {
	if !my.skip {
		printColor(D_RedItalic, a...)
	}
}
func (my iColor) GreenItalic(a ...interface{}) {
	if !my.skip {
		printColor(D_GreenItalic, a...)
	}
}
func (my iColor) YellowItalic(a ...interface{}) {
	if !my.skip {
		printColor(D_YellowItalic, a...)
	}
}
func (my iColor) BlueItalic(a ...interface{}) {
	if !my.skip {
		printColor(D_BlueItalic, a...)
	}
}
func (my iColor) PurpleItalic(a ...interface{}) {
	if !my.skip {
		printColor(D_PurpleItalic, a...)
	}
}
func (my iColor) CyanItalic(a ...interface{}) {
	if !my.skip {
		printColor(D_CyanItalic, a...)
	}
}
func (my iColor) WhiteItalic(a ...interface{}) {
	if !my.skip {
		printColor(D_WhiteItalic, a...)
	}
}

func (my iColor) BlackItalicBG(a ...interface{}) {
	if !my.skip {
		printColor(D_BlackItalicBG, a...)
	}
}
func (my iColor) RedItalicBG(a ...interface{}) {
	if !my.skip {
		printColor(D_RedItalicBG, a...)
	}
}
func (my iColor) GreenItalicBG(a ...interface{}) {
	if !my.skip {
		printColor(D_GreenItalicBG, a...)
	}
}
func (my iColor) YellowItalicBG(a ...interface{}) {
	if !my.skip {
		printColor(D_YellowItalicBG, a...)
	}
}
func (my iColor) BlueItalicBG(a ...interface{}) {
	if !my.skip {
		printColor(D_BlueItalicBG, a...)
	}
}
func (my iColor) PurpleItalicBG(a ...interface{}) {
	if !my.skip {
		printColor(D_PurpleItalicBG, a...)
	}
}
func (my iColor) CyanItalicBG(a ...interface{}) {
	if !my.skip {
		printColor(D_CyanItalicBG, a...)
	}
}
func (my iColor) WhiteItalicBG(a ...interface{}) {
	if !my.skip {
		printColor(D_WhiteItalicBG, a...)
	}
}

func (my iColor) BlackUL(a ...interface{}) {
	if !my.skip {
		printColor(D_BlackUL, a...)
	}
}
func (my iColor) RedUL(a ...interface{}) {
	if !my.skip {
		printColor(D_RedUL, a...)
	}
}
func (my iColor) GreenUL(a ...interface{}) {
	if !my.skip {
		printColor(D_GreenUL, a...)
	}
}
func (my iColor) YellowUL(a ...interface{}) {
	if !my.skip {
		printColor(D_YellowUL, a...)
	}
}
func (my iColor) BlueUL(a ...interface{}) {
	if !my.skip {
		printColor(D_BlueUL, a...)
	}
}
func (my iColor) PurpleUL(a ...interface{}) {
	if !my.skip {
		printColor(D_PurpleUL, a...)
	}
}
func (my iColor) CyanUL(a ...interface{}) {
	if !my.skip {
		printColor(D_CyanUL, a...)
	}
}
func (my iColor) WhiteUL(a ...interface{}) {
	if !my.skip {
		printColor(D_WhiteUL, a...)
	}
}

func (my iColor) BlackULBG(a ...interface{}) {
	if !my.skip {
		printColor(D_BlackULBG, a...)
	}
}
func (my iColor) RedULBG(a ...interface{}) {
	if !my.skip {
		printColor(D_RedULBG, a...)
	}
}
func (my iColor) GreenULBG(a ...interface{}) {
	if !my.skip {
		printColor(D_GreenULBG, a...)
	}
}
func (my iColor) YellowULBG(a ...interface{}) {
	if !my.skip {
		printColor(D_YellowULBG, a...)
	}
}
func (my iColor) BlueULBG(a ...interface{}) {
	if !my.skip {
		printColor(D_BlueULBG, a...)
	}
}
func (my iColor) PurpleULBG(a ...interface{}) {
	if !my.skip {
		printColor(D_PurpleULBG, a...)
	}
}
func (my iColor) CyanULBG(a ...interface{}) {
	if !my.skip {
		printColor(D_CyanULBG, a...)
	}
}
func (my iColor) WhiteULBG(a ...interface{}) {
	if !my.skip {
		printColor(D_WhiteULBG, a...)
	}
}
