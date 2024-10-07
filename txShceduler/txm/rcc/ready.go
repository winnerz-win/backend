package rcc

import (
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp"
)

func Ready(classic *chttp.Classic) {
	dbg.PrintForce("rcc.Ready ::::: Remote console command receiver")

	cmdDBClear(classic)
}
