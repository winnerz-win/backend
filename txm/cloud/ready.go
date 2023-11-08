package cloud

import (
	"time"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/runtext"
	"txscheduler/txm/inf"
)

func get_sender_x() *ecsx.Sender {
	return ecsx.New(inf.Mainnet(), inf.InfuraKey())
}

// Ready :
func Ready() runtext.Starter {
	rtx := runtext.New("cloud")

	if inf.LOCALMODE {
		return rtx
	}

	ethDriver := func() {
		go runSyncCoin(rtx)
		go runETHWithdraw(rtx)
		go runETHCharger(rtx)
		go runETHDepositToMaster(rtx)
		go runETHDepositChn(rtx)
		go runETHCollection(rtx)
		go runSELFWithdraw(rtx)
	}
	_ = ethDriver

	if !inf.Mainnet() && inf.Args().Do("nochain") {
		skipView := func() {
			dbg.Red("################################################")
			dbg.Red("################################################")
			dbg.Red(" SKIP - CLOUD Driver")
			dbg.Red("################################################")
			dbg.Red("################################################")
		}
		go func() {
			for i := 0; i < 3; i++ {
				skipView()
				time.Sleep(time.Second)
			}
		}()
	} else {
		ethDriver()
	}

	return rtx
}

func logTag(tag string, a ...interface{}) []interface{} {
	logs := []interface{}{
		tag,
	}
	logs = append(logs, a...)
	return logs
}

func logDeposit(a ...interface{}) {
	dbg.Purple(logTag("[Deposit]", a...)...)
}

func logCollection(a ...interface{}) {
	dbg.Gray(logTag("[Collection]", a...)...)
}

func logCharger(a ...interface{}) {
	dbg.Purple(logTag("[Charger]", a...)...)
}

func logWithdraw(a ...interface{}) {
	dbg.Purple(logTag("[Withdraw]", a...)...)
}
