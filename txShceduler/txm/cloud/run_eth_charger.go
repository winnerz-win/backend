package cloud

import (
	"context"
	"jtools/cc"
	"jtools/cloud/ebcm"
	"jtools/jmath"
	"jtools/mms"
	"time"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jticker"
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

					result, _, _ := checker.TransactionByHash(item.Hash)
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

			from := inf.Charger()
			sender := get_sender_x()

			nonce, err := ebcm.MMA_GetNonce(
				sender,
				from.Address,
			)
			if err != nil {
				dbg.Red("ebcm.MMA_GetNonce ::: ", err)
				return
			}

			fromWEI := sender.Balance(from.Address)
			if jmath.CMP(fromWEI, 0) <= 0 {
				dbg.RedBG("runETHCharger::ChargerETH is ZERO")
				for _, group := range groupBy {
					group.RemoveForError(db)
				}
				return
			}

			ctx := context.Background()

			gas_price, err := sender.SuggestGasPrice(ctx)
			if err != nil {
				dbg.RedBG("charger.SuggestGasPrice ::", err)
				for _, group := range groupBy {
					group.RemoveForError(db)
				}
				return
			}

			type BOX struct {
				group model.TxChargeGroup

				TransferData
				// padbytes ebcm.PADBYTES
				// limit    uint64
				// nonce    uint64

				// wei string

				// stx ebcm.WrappedTransaction
			}

			nonceCounter := nonce
			boxlist := []BOX{}
			for _, group := range groupBy {
				if jmath.CMP(fromWEI, 0) < 0 {
					cc.RedItalic("charger.remain_coin_price is zero.")
					group.RemoveForError(db)
					continue
				}

				wei := ebcm.ETHToWei(group.Price)

				padbytes := ebcm.PadByteETH()

				limit, err := sender.EstimateGas(
					ctx,
					ebcm.MakeCallMsg(
						from.Address,
						group.Address,
						wei,
						padbytes,
					),
				)
				if err != nil {
					cc.RedItalic("charger.EstimateGas :: ", err)
					group.RemoveForError(db)
					continue
				}

				est_gas_wei := gas_price.EstimateGasFeeWEI(limit)

				fromWEI = jmath.SUB(fromWEI, est_gas_wei)
				if jmath.CMP(fromWEI, 0) < 0 {
					cc.RedItalic("charger.remain_coin_price is zero.")
					group.RemoveForError(db)
					continue
				}

				ntx := sender.NewTransaction(
					nonceCounter,
					group.Address,
					wei,
					limit,
					gas_price,
					padbytes,
				)

				stx, err := sender.SignTx(
					ntx,
					from.PrivateKey,
				)
				if err != nil {
					cc.RedItalic("charger.SignTx :::", err)
					group.RemoveForError(db)
					continue
				}

				boxlist = append(boxlist, BOX{
					group: group,
					TransferData: TransferData{
						padbytes: padbytes,
						limit:    limit,
						wei:      wei,
						nonce:    nonceCounter,
						stx:      stx,
					},
				})

				nonceCounter++
			} //for

			is_err_stop := false
			for _, box := range boxlist {

				if is_err_stop {
					box.group.RemoveForError(db)
					continue
				}

				hash, err := sender.SendTransaction(
					ctx,
					box.stx,
				)
				if err != nil {
					is_err_stop = true
					box.group.RemoveForError(db)
					continue
				}

				for _, item := range box.group.List {
					logCharger("Send [", item.UID, "]", item.Address, item.Price)
					db.C(inf.TXETHCharger).Update(
						item.Selector(),
						mongo.Bson{"$set": mongo.Bson{
							"hash":  hash,
							"state": model.TxStatePending,
						}},
					)
				} //for

			} //for

		})
	} //for
}
