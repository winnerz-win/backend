package rcc

import (
	"time"
	"txscheduler/brix/tools/console"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func cmdDBClear(classic *chttp.Classic) {
	console.AppendCmd(
		"dbclear",
		"dbclear [true] --- only testnet",
		false,
		func(ps []string) {
			if inf.Mainnet() {
				console.Log("only test-net version")
				return
			}
			if dbg.IsTrue(ps[0]) == false {
				console.Log("do must true")
				return
			}

			go func() {
				defer dbg.Red("server-exit")

				console.TestCommand("exit")
				for {
					if classic.IsExit() == false {
						time.Sleep(time.Millisecond * 10)
						continue
					}
					break
				} //for

				time.Sleep(time.Second * 2)

				model.DB(func(db mongo.DATABASE) {
					rmvfunc := func(name string) {
						i, _ := db.C(name).RemoveAll(nil)
						console.Log(name, i)
					}
					rmvfunc(inf.COLMember)
					rmvfunc(inf.COLInfoDeposit)
					rmvfunc(inf.COLInfoMaster)
					rmvfunc(inf.COLLogDeposit)
					rmvfunc(inf.COLLogWithdraw)
					rmvfunc(inf.TXETHCount)
					rmvfunc(inf.TXETHBlock)
					rmvfunc(inf.TXETHCharger)

					rmvfunc(inf.TXETHDepositLog)
					rmvfunc(inf.TXETHWithdraw)
					rmvfunc(inf.TXETHInternalCnt)
					rmvfunc(inf.TXETHTokenEx)

					msgEnd := "dbclear ---- end"
					console.Log(msgEnd)
					console.Log(msgEnd)
					console.Log(msgEnd)
					console.Log(msgEnd)
				})
			}()

			console.Log("dbclear == Call")

			/*
				txm.coin_day removeAll nil
				txm.log_deposit removeAll nil
				txm.log_tomaster removeAll nil
				txm.log_withdraw removeAll nil
				txm.member removeAll nil
				txm.tx_eth_block removeAll nil
				txm.tx_eth_deposit_log removeAll nil
				txm.tx_eth_withdraw removeAll nil
				txm.coin_sum update {"key":"info_coin_sum"} {"$set":{"coin.ETH":"0", "coin.ERCT":"0", "coin.GDG":"0", "coin.USDT":"0"}}
			*/

		},
	)
}
