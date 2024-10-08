package rcc

import (
	"jcloudnet/itype"
	"jtools/cloud/jeth/ecs"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp"
)

func Ready(classic *chttp.Classic) {
	dbg.PrintForce("rcc.Ready ::::: Remote console command receiver")

	cmdDBClear(classic)
}

func Finder(mainnet bool, infura_key string) *itype.IClient {
	return itype.New(ecs.RPC_URL(mainnet), false, infura_key)
}
