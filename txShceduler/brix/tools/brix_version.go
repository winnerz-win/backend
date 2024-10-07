package tools

import (
	"runtime"
	"txscheduler/brix/tools/dbg/cc"
)

const (
	gversion = `
────────────────────────────────────────────────────
         ### BRIX - 2021.03.19 ###`

	//BAR : ────────────────────────────────────────────────────
	BAR = "────────────────────────────────────────────────────"
	//DOT : ....................................................
	DOT = "...................................................."
	//ENTER : \n
	ENTER = "\n"
)

var verfunc func() string

func Version() string {
	profileLog()
	if verfunc != nil {
		return cc.Yellow.String() +
			gversion + ENTER +
			"os : " + runtime.GOOS + ENTER +
			verfunc() + ENTER +
			BAR +
			cc.END.String() +
			ENTER
	}
	return cc.Yellow.String() + gversion + ENTER + BAR + cc.END.String() + ENTER
}

//SetVersion :
func SetVersion(f func() string) {
	verfunc = f
}
