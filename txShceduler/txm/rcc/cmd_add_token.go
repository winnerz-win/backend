package rcc

import (
	"txscheduler/brix/tools/console"
	"txscheduler/brix/tools/dbg"
	"txscheduler/txm/inf"
)

func init() {
	cmdAddToken()
}

func cmdAddToken() {
	console.AppendCmd(
		"addtoken",
		"addtoken [contract]",
		false,
		func(ps []string) {
			contract := dbg.TrimToLower(ps[0])

			isAdd := inf.AddToken(contract)
			if isAdd == false {
				console.Log("token is not add :", contract)
				return
			}

			token := inf.TokenList().GetContract(contract)
			console.Log(token)
		},
	)
}
