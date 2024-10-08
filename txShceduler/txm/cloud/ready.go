package cloud

import (
	"jcloudnet/itype"
	"jtools/cloud/ebcm"
	"time"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/runtext"
	"txscheduler/txm/inf"
)

// func get_sender_x() *ebcm.Sender {
func get_sender_x() *ebcm.Sender {
	return inf.GetSender()
}

func get_finder() *itype.IClient {
	return inf.GetFinder()
}

// Ready :
func Ready() runtext.Starter {
	rtx := runtext.New("cloud")

	if inf.LOCALMODE {
		return rtx
	}

	_cloud_master_console()

	ethDriver := func() {

		go runSyncCoin(rtx)
		go runETHCharger(rtx)
		go runETHDepositToMaster(rtx)
		go runETHDepositChn(rtx)
		go runETHCollection(rtx)
		go runSELFWithdraw(rtx)

		if inf.Master().Address != inf.Owner().Address {
			dbg.Yellow("MASTER <> OWNER  XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")

			go runETHWithdraw(rtx, false)
			if inf.IsOnwerTaskMode() {
				go runOwnerTaskLogSender()
				go runOwnerTaskAction(rtx)
			}
		} else {
			same_tag := "============= MASTER == OWNER  ==========================================="
			dbg.Yellow(same_tag)
			go func() {
				for i := 0; i < 10; i++ {
					dbg.Yellow(same_tag)
					time.Sleep(time.Second)
				}
			}()

			is_owner_master_same := true
			if !inf.IsOnwerTaskMode() {
				is_owner_master_same = false
			}
			go runETHWithdraw(rtx, is_owner_master_same)
		}

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

type TransferData struct {
	padbytes ebcm.PADBYTES
	limit    uint64
	nonce    uint64

	wei string

	stx ebcm.WrappedTransaction
}
