package cloud

import (
	"time"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
	"txscheduler/brix/tools/jticker"
	"txscheduler/brix/tools/mms"
	"txscheduler/brix/tools/runtext"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func runETHCharger(rtx runtext.Runner) {
	defer dbg.PrintForce("cloud.runETHCharger ---- END")
	<-rtx.WaitStart()
	dbg.PrintForce("cloud.runETHCharger ---- START")

	ticker := jticker.New(time.Second*5, time.Millisecond*200, true)
	for {
		if ticker.IsWait() {
			continue
		}

		model.DB(func(db mongo.DATABASE) {

			//pending-check
			list := model.TxETHChargerList{}
			db.C(inf.TXETHCharger).Find(mongo.Bson{"state": model.TxStatePending}).All(&list)
			pendingCnt := len(list)
			if pendingCnt > 0 {
				for _, item := range list {
					checker := get_sender_x()
					//result := checker.Receipt(item.Hash)

					result, _, _, _ := checker.TransactionByHash(item.Hash)
					if !result.IsReceiptedByHash {
						continue
					}
					feeGas := result.GetTransactionFee()
					model.CoinDay{}.AddChargerGas(db, feeGas, mms.Now())

					if result.IsError {
						db.C(inf.TXETHDepositLog).Update(
							mongo.Bson{"key": item.Key},
							mongo.Bson{"$set": mongo.Bson{"state": model.TxStateNone}},
						)
						db.C(inf.TXETHCharger).Remove(item.Selector())
					} else {
						db.C(inf.TXETHDepositLog).Update(
							mongo.Bson{"key": item.Key},
							mongo.Bson{"$set": mongo.Bson{"state": model.TxStateNone}},
						)
						db.C(inf.TXETHCharger).Remove(item.Selector())
					}

				} //for
				return
			}

			list = model.TxETHChargerList{}
			db.C(inf.TXETHCharger).Find(mongo.Bson{"state": model.TxStateNone}).All(&list)
			if len(list) == 0 {
				return
			}

			groupBy := list.GroupBy()
			if len(groupBy) == 0 {
				return
			}

			speed := ecsx.GasAverage
			from := inf.Charger()
			sender := get_sender_x()

			fromWEI := sender.Balance(from.Address)
			if jmath.CMP(fromWEI, 0) <= 0 {
				dbg.RedBG("runETHCharger::ChargerETH is ZERO")
				for _, group := range groupBy {
					group.RemoveForError(db)
				}
				return
			}

			type BOX struct {
				group model.TxChargeGroup
				box   ecsx.GasBoxData
			}
			boxlist := []BOX{}
			for _, group := range groupBy {
				wei := ecsx.ETHToWei(group.Price)
				box := sender.GasBox("eth", from.Address, group.Address, wei, speed)
				if box.Error != nil {
					dbg.Red("runETHCharger.box :", box.Error)
					group.RemoveForError(db)
					continue
				}
				boxlist = append(boxlist, BOX{
					group: group,
					box:   box,
				})
			} //for

			nonce, err := sender.Nonce(from.PrivateKey)
			if err != nil {
				dbg.RedItalic("charger.NONCE :", err)
				return
			}
			type NTX struct {
				group model.TxChargeGroup
				tx    *ecsx.NTX
			}
			ntxlist := []NTX{}
			var nonceCounter uint64 = 0
			for _, v := range boxlist {
				if ntx, err := nonce.BoxTx(v.box, nonceCounter); err == nil {
					ntxlist = append(ntxlist, NTX{
						group: v.group,
						tx:    ntx,
					})
					nonceCounter++
				}
			} //for

			type STX struct {
				group model.TxChargeGroup
				tx    *ecsx.STX
			}
			stxlist := []STX{}
			for _, v := range ntxlist {
				if stx, err := v.tx.Tx(); err == nil {
					stxlist = append(stxlist, STX{
						group: v.group,
						tx:    stx,
					})
				}
			} //for

			type SEND struct {
				group model.TxChargeGroup
				hash  string
			}
			sendlist := []SEND{}
			for _, v := range stxlist {
				if err := v.tx.Send(); err == nil {
					sendlist = append(sendlist, SEND{
						group: v.group,
						hash:  v.tx.Hash(),
					})
				}
			} //for

			for _, v := range sendlist {
				for _, item := range v.group.List {
					logCharger("Send [", item.UID, "]", item.Address, item.Price)
					db.C(inf.TXETHCharger).Update(
						item.Selector(),
						mongo.Bson{"$set": mongo.Bson{
							"hash":  v.hash,
							"state": model.TxStatePending,
						}},
					)
				} //for
			} //for

		})
	} //for
}
