package cc

import (
	"fmt"
)

const (
	//color = "\033[(normal/bold);30m"
	gray = "\033[2;37m"

	black  = "\033[0;30m"
	red    = "\033[0;31m"
	green  = "\033[0;32m"
	yellow = "\033[0;33m"
	blue   = "\033[0;34m"
	purple = "\033[0;35m"
	cyan   = "\033[0;36m"
	white  = "\033[0;37m"

	blackBold  = "\033[1;30m"
	redBold    = "\033[1;31m"
	greenBold  = "\033[1;32m"
	yellowBold = "\033[1;33m"
	blueBold   = "\033[1;34m"
	purpleBold = "\033[1;35m"
	cyanBold   = "\033[1;36m"
	whiteBold  = "\033[1;37m"

	blackBG  = "\033[0;40m"
	redBG    = "\033[0;41m"
	greenBG  = "\033[0;42m"
	yellowBG = "\033[0;43m"
	blueBG   = "\033[0;44m"
	purpleBG = "\033[0;45m"
	cyanBG   = "\033[0;46m"
	whiteBG  = "\033[0;37m"

	blackBoldBG  = "\033[1;40m"
	redBoldBG    = "\033[1;41m"
	greenBoldBG  = "\033[1;42m"
	yellowBoldBG = "\033[1;43m"
	blueBoldBG   = "\033[1;44m"
	purpleBoldBG = "\033[1;45m"
	cyanBoldBG   = "\033[1;46m"
	whiteBoldBG  = "\033[1;47m"

	blackItalic  = "\033[3;30m"
	redItalic    = "\033[3;31m"
	greenItalic  = "\033[3;32m"
	yellowItalic = "\033[3;33m"
	blueItalic   = "\033[3;34m"
	purpleItalic = "\033[3;35m"
	cyanItalic   = "\033[3;36m"
	whiteItalic  = "\033[3;37m"

	blackItalicBG  = "\033[3;40m"
	redItalicBG    = "\033[3;41m"
	greenItalicBG  = "\033[3;42m"
	yellowItalicBG = "\033[3;43m"
	blueItalicBG   = "\033[3;44m"
	purpleItalicBG = "\033[3;45m"
	cyanItalicBG   = "\033[3;46m"
	whiteItalicBG  = "\033[3;47m"

	blackUL  = "\033[4;30m"
	redUL    = "\033[4;31m"
	greenUL  = "\033[4;32m"
	yellowUL = "\033[4;33m"
	blueUL   = "\033[4;34m"
	purpleUL = "\033[4;35m"
	cyanUL   = "\033[4;36m"
	whiteUL  = "\033[4;37m"

	blackULBG  = "\033[4;40m"
	redULBG    = "\033[4;41m"
	greenULBG  = "\033[4;42m"
	yellowULBG = "\033[4;43m"
	blueULBG   = "\033[4;44m"
	purpleULBG = "\033[4;45m"
	cyanULBG   = "\033[4;46m"
	whiteULBG  = "\033[4;47m"

	end = "\033[0m"
)

// COLOR :
type COLOR string

func (my COLOR) String() string {
	if isLegacyMode {
		return ""
	}
	return string(my)
}

var (
	Gray = COLOR(gray)

	Black  = COLOR(black)
	Red    = COLOR(red)
	Green  = COLOR(green)
	Yellow = COLOR(yellow)
	Blue   = COLOR(blue)
	Purple = COLOR(purple)
	Cyan   = COLOR(cyan)
	White  = COLOR(white)

	BlackBold  = COLOR(blackBold)
	RedBold    = COLOR(redBold)
	GreenBold  = COLOR(greenBold)
	YellowBold = COLOR(yellowBold)
	BlueBold   = COLOR(blueBold)
	PurpleBold = COLOR(purpleBold)
	CyanBold   = COLOR(cyanBold)
	WhiteBold  = COLOR(whiteBold)

	BlackBG  = COLOR(blackBG)
	RedBG    = COLOR(redBG)
	GreenBG  = COLOR(greenBG)
	YellowBG = COLOR(yellowBG)
	BlueBG   = COLOR(blueBG)
	PurpleBG = COLOR(purpleBG)
	CyanBG   = COLOR(cyanBG)
	WhiteBG  = COLOR(whiteBG)

	BlackBoldBG  = COLOR(blackBoldBG)
	RedBoldBG    = COLOR(redBoldBG)
	GreenBoldBG  = COLOR(greenBoldBG)
	YellowBoldBG = COLOR(yellowBoldBG)
	BlueBoldBG   = COLOR(blueBoldBG)
	PurpleBoldBG = COLOR(purpleBoldBG)
	CyanBoldBG   = COLOR(cyanBoldBG)
	WhiteBoldBG  = COLOR(whiteBoldBG)

	BlackItalic  = COLOR(blackItalic)
	RedItalic    = COLOR(redItalic)
	GreenItalic  = COLOR(greenItalic)
	YellowItalic = COLOR(yellowItalic)
	BlueItalic   = COLOR(blueItalic)
	PurpleItalic = COLOR(purpleItalic)
	CyanItalic   = COLOR(cyanItalic)
	WhiteItalic  = COLOR(whiteItalic)

	BlackItalicBG  = COLOR(blackItalicBG)
	RedItalicBG    = COLOR(redItalicBG)
	GreenItalicBG  = COLOR(greenItalicBG)
	YellowItalicBG = COLOR(yellowItalicBG)
	BlueItalicBG   = COLOR(blueItalicBG)
	PurpleItalicBG = COLOR(purpleItalicBG)
	CyanItalicBG   = COLOR(cyanItalicBG)
	WhiteItalicBG  = COLOR(whiteItalicBG)

	BlackUL  = COLOR(blackUL)
	RedUL    = COLOR(redUL)
	GreenUL  = COLOR(greenUL)
	YellowUL = COLOR(yellowUL)
	BlueUL   = COLOR(blueUL)
	PurpleUL = COLOR(purpleUL)
	CyanUL   = COLOR(cyanUL)
	WhiteUL  = COLOR(whiteUL)

	BlackULBG  = COLOR(blackULBG)
	RedULBG    = COLOR(redULBG)
	GreenULBG  = COLOR(greenULBG)
	YellowULBG = COLOR(yellowULBG)
	BlueULBG   = COLOR(blueULBG)
	PurpleULBG = COLOR(purpleULBG)
	CyanULBG   = COLOR(cyanULBG)
	WhiteULBG  = COLOR(whiteULBG)

	END = COLOR(end)

	isLegacyMode = false
)

func initOS(osName string) {

	bar := func() {
		fmt.Println(Black.String(), "========================================", END.String())
	}
	defer bar()
	bar()

	msg := Black.String()
	msg += "..."
	msg += Red.String()
	msg += "dbg."
	msg += Green.String()
	msg += "color."
	msg += Yellow.String()
	msg += "init."
	msg += Blue.String()
	msg += "func."
	msg += Purple.String()
	msg += "call."
	msg += Cyan.String()
	msg += "for."
	msg += White.String()
	msg += osName
	msg += END.String()
	fmt.Println(msg)

}

// SetLegacyMode :
func SetLegacyMode(on bool) {
	isLegacyMode = on
}

// LegacyMode :
func LegacyMode() bool {
	return isLegacyMode
}
