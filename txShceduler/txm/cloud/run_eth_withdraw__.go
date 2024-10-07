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

func runETHWithdraw(rtx runtext.Runner) {
	defer dbg.PrintForce("cloud.runETHWithdraw ---- END")
	<-rtx.WaitStart()
	dbg.PrintForce("cloud.runETHWithdraw ---- START")

	InjectMasterWithdrawProcess(
		func(db mongo.DATABASE, nowAt mms.MMS) bool {
			return procMasterOutPending(db, nowAt)
		},
		func(db mongo.DATABASE, nowAt mms.MMS) bool {
			return procMasterOutTry(db, nowAt)
		},
	)

	InjectMasterWithdrawProcess(
		func(db mongo.DATABASE, nowAt mms.MMS) bool {
			return procETHWithdrawPending(db, nowAt)
		},
		func(db mongo.DATABASE, nowAt mms.MMS) bool {
			return procETHWithdarwTry(db, nowAt) //가상계좌 잔액 회수
		},
	)

	ticker := jticker.New(time.Second*5, time.Millisecond*200, true)

EXIT:
	for {
		select {
		case <-rtx.EndC():
			break EXIT
		default:
		} //select

		if ticker.IsWait() {
			continue
		}

		model.DB(func(db mongo.DATABASE) {
			nowAt := mms.Now()

			for _, pending_func := range master_withdraw_func.pending_funcs {
				if pending_func(db, nowAt) {
					return
				}
			} //for

			for _, try_func := range master_withdraw_func.try_funcs {
				if try_func(db, nowAt) {
					return
				}
			}

			// if procMasterOutPending(db, nowAt) {
			// 	return
			// }
			// if procMasterOutTry(db, nowAt) {
			// 	return
			// }

			// if procETHWithdrawPending(db, nowAt) {
			// 	return
			// }
			// if procETHWithdarwTry(db, nowAt) { //가상계좌 잔액 회수
			// 	return
			// }

		})
	} //for
}

const (
	MasterFail_InvalidSymbol         = "invalid_symbol"
	MasterFail_NeedPrice             = "need_price"
	MasterFail_ChainErrorBox         = "chain_error:box"
	MasterFail_ChainErrorNonce       = "chain_error:nonce"
	MasterFail_ChainErrorTx          = "chain_error:tx"
	MasterFail_ChainErrorSend        = "chain_error:send"
	MasterFail_ChianErrorPendingTime = "chain_error:pending_time"
	MasterFail_ChianErrorFail        = "chain_error:fail"
)

var (
	gasSpeed                = ecsx.GasFast
	pendingWaitHour float64 = 2
)

func procMasterOutTry(db mongo.DATABASE, nowAt mms.MMS) bool {
	list := model.TxETHMasterOutList{}

	selector := mongo.Bson{"state": model.TxStateNone}
	db.C(inf.TXETHMasterOutTry).Find(selector).Sort("timestamp").Limit(30).All(&list)
	if len(list) == 0 {
		return false
	}

	sender := get_sender_x()
	if sender == nil {
		model.LogError.WriteLog(
			db,
			model.ErrorFinderNull,
			"procMasterOutTry.sender",
		)
		return false
	}
	from := inf.Master()
	speed := gasSpeed

	senderWallet := model.NewCoinData()
	tokenlist := inf.TokenList()
	for _, token := range tokenlist {
		wei := sender.Balance2(from.Address, token.Contract)
		senderWallet[token.Symbol] = ecsx.WeiToToken(wei, token.Decimal)
	} //for

	sendingFail := func(item model.TxETHMasterOut, msg string) {
		item.State = model.TxStateFail
		item.FailMessage = msg
		item.InsertLog(db)
		item.RemoveTry(db)
	}
	_ = sendingFail

	type Fixed struct {
		item model.TxETHMasterOut
		box  ecsx.GasBoxData
	}
	fixlist := []Fixed{}
	sendinglist := []Fixed{}

	sumWallet := model.NewCoinData()
	for _, item := range list {
		info := tokenlist.GetSymbol(item.Symbol)
		if !info.Valid() {
			sendingFail(item, MasterFail_InvalidSymbol)
			continue
		}

		box := sender.GasBox(
			info.Contract,
			from.Address,
			item.ToAddress,
			item.Wei(),
			speed,
		)
		if box.Error != nil {
			sendingFail(item, MasterFail_ChainErrorBox)
			continue
		}

		if item.Symbol == model.ETH {
			sendETH := box.SpendETH()
			cmpValue := jmath.ADD(sumWallet.Price(info.Symbol), sendETH)
			if jmath.CMP(cmpValue, senderWallet.Price(info.Symbol)) > 0 {
				sendingFail(item, MasterFail_NeedPrice)
				continue
			}
			sumWallet.ADD(info.Symbol, cmpValue)
		} else {
			gas := box.GasETH()
			cmpETH := jmath.ADD(sumWallet.Price(model.ETH), gas)
			if jmath.CMP(cmpETH, senderWallet.Price(model.ETH)) > 0 {
				sendingFail(item, MasterFail_NeedPrice)
				continue
			}
			cmpToken := jmath.ADD(sumWallet.Price(info.Symbol), item.ToPrice)
			if jmath.CMP(cmpToken, senderWallet.Price(info.Symbol)) > 0 {
				sendingFail(item, MasterFail_NeedPrice)
				continue
			}
			sumWallet.ADD(model.ETH, cmpETH)
			sumWallet.ADD(info.Symbol, cmpToken)
		}
		fixlist = append(fixlist, Fixed{
			item: item,
			box:  box,
		})
	} //for

	if len(fixlist) == 0 {
		return true
	}

	nonce, err := sender.Nonce(from.PrivateKey)
	if err != nil {
		dbg.Red(MasterFail_ChainErrorNonce)
		return true
	}
	pending, _ := sender.XPendingNonceAt(from.Address)
	if nonce.NonceCount() != pending {
		dbg.RedItalic("procMasterOutTry.pendingNonce isDiffer (", nonce.NonceCount(), "/", pending, ")")
		return true
	}

	var nonceCount uint64 = 0
	for i, fix := range fixlist {
		ntx, err := nonce.BoxTx(fix.box, nonceCount)
		if err != nil {
			sendingFail(fix.item, MasterFail_ChainErrorTx)
			continue
		}
		stx, err := ntx.Tx()
		if err != nil {
			sendingFail(fix.item, MasterFail_ChainErrorTx)
			continue
		}
		if err := stx.Send(); err != nil {
			sendingFail(fix.item, MasterFail_ChainErrorSend)
			continue
		}

		fixlist[i].item.Hash = stx.Hash()
		fixlist[i].item.Gas = fix.box.GasETH()
		fixlist[i].item.Timestamp = nowAt
		fixlist[i].item.YMD = nowAt.YMD()
		fixlist[i].item.State = model.TxStatePending
		sendinglist = append(sendinglist, fixlist[i])
		nonceCount++
	} //for

	for _, sendingdata := range sendinglist {
		item := sendingdata.item
		item.UpdateTry(db)

		logWithdraw("MasterOut Send ", item.ToPrice, item.Symbol)
	} //for

	return true
}

// fail_hash_continue_check : true(continue) , false(fail)
func fail_hash_continue_check(checker *ecsx.Sender, fail_hash string) bool {
	time.Sleep(time.Second)

	r := checker.ReceiptByHash(fail_hash) //TX 실패일경우 한번더 확인
	if !r.IsAck() {
		return true //not found
	} else {
		if r.IsSuccess() {
			return true
		}
	}
	return false
}

func procMasterOutPending(db mongo.DATABASE, nowAt mms.MMS) bool {
	list := model.TxETHMasterOutList{}
	selector := mongo.Bson{"state": model.TxStatePending}
	db.C(inf.TXETHMasterOutTry).Find(selector).Sort("timestamp").All(&list)
	if len(list) == 0 {
		return false
	}

	for _, item := range list {
		checker := get_sender_x()
		if item.CancelTry {
			result, _, _, _ := checker.TransactionByHash(item.CancelHash)
			if !result.IsReceiptedByHash {
				continue
			}
			feeGas := result.GetTransactionFee()
			model.CoinDay{}.AddMasterGas(db, feeGas, mms.Now())
			item.State = model.TxStateFail
			item.FailMessage = MasterFail_ChianErrorPendingTime
			item.InsertLog(db)
			item.RemoveTry(db)

		} else {
			result, _, _, _ := checker.TransactionByHash(item.Hash)
			if !result.IsReceiptedByHash {
				duration := nowAt.Sub(item.Timestamp)
				if duration.Hours() >= pendingWaitHour { //2시간경과시 실패 처리 하자
					if !item.CancelTry {
						sender := get_sender_x()
						from := inf.Master()
						to := from.Address
						speed := ecsx.GasFast

						box := sender.GasBox(
							"eth",
							from.Address,
							to,
							model.ZERO,
							speed,
						)
						if box.Error != nil {
							continue
						}
						nonce, err := sender.Nonce(from.PrivateKey)
						if err != nil {
							continue
						}
						ntx, err := nonce.BoxTx(box)
						if err != nil {
							continue
						}
						stx, err := ntx.Tx()
						if err != nil {
							continue
						}
						if err := stx.Send(); err != nil {
							continue
						}
						dbg.Red("[ MasterOutTry Pending Over] hash:", item.Hash, " ---- Withdraw.Cancel.Try")
						item.CancelTry = true
						item.CancelHash = stx.Hash()
						item.CancelTime = nowAt
						item.CancelYMD = nowAt.YMD()
						item.UpdateTry(db)
					}
				}
				continue
			}

			item.Gas = result.GetTransactionFee()
			if result.IsError { //fail
				if fail_hash_continue_check(checker, item.Hash) { //TX 실패일경우 한번더 확인
					continue
				}

				logWithdraw("( MasterOutTry ) r-Fail ", item.ToPrice, item.Symbol)
				model.CoinDay{}.AddMasterGas(db, item.Gas, mms.Now())

				item.State = model.TxStateFail
				item.FailMessage = MasterFail_ChianErrorFail
				item.InsertLog(db)
				item.RemoveTry(db)

			} else { //success
				logWithdraw("( MasterOutTry ) r-Success ", item.ToPrice, item.Symbol)

				model.CoinDay{}.AddWithdraw(db, item.Symbol, item.ToPrice, item.Gas, mms.Now())

				item.State = model.TxStateSuccess
				item.InsertLog(db)
				item.RemoveTry(db)
			}
		}
	} //for

	pending_count, _ := db.C(inf.TXETHMasterOutTry).Find(selector).Count()
	return pending_count > 0
}

/////////////////
/////////////////
/////////////////
/////////////////
/////////////////
/////////////////

func procETHWithdarwTry(db mongo.DATABASE, nowAt mms.MMS) (isChargeAble bool) {
	isChargeAble = false

	list := model.TxETHWithdrawList{}

	selector := mongo.Bson{"state": model.TxStateNone}
	db.C(inf.TXETHWithdraw).Find(selector).Sort("timestamp").Limit(10).All(&list)

	if len(list) == 0 {
		return
	}

	sender := get_sender_x()
	if sender == nil {
		model.LogError.WriteLog(
			db,
			model.ErrorFinderNull,
			"procETHWithdarwTry.sender",
		)
		return false
	}

	from_master := inf.Master()
	speed := gasSpeed

	senderWallet := model.NewCoinData()
	tokenlist := inf.TokenList()
	for _, token := range tokenlist {
		wei := sender.Balance2(from_master.Address, token.Contract)
		senderWallet[token.Symbol] = ecsx.WeiToToken(wei, token.Decimal)
	} //for

	if jmath.CMP(senderWallet.Price(model.ETH), list.TotalPrice(model.ETH)) <= 0 {
		isChargeAble = true
	}

	sendingFail := func(item model.TxETHWithdraw) {
		model.LockMemberUID(db, item.UID, func(member model.Member) {
			log := item.MakeLogWithdraw(member, model.TxStateFail)
			log.InsertDB(db, mms.Now())

			db.C(inf.TXETHWithdraw).Remove(item.Selector())
		})
	}
	_ = sendingFail

	type Fixed struct {
		item model.TxETHWithdraw
		box  ecsx.GasBoxData
	}
	fixlist := []Fixed{}
	sendinglist := []Fixed{}

	sumWallet := model.NewCoinData()
	for _, item := range list {
		info := tokenlist.GetSymbol(item.Symbol)
		if !info.Valid() {
			sendingFail(item)
			continue
		}

		box := sender.GasBox(
			info.Contract,
			from_master.Address,
			item.ToAddress,
			item.Wei(),
			speed,
		)
		if box.Error != nil {
			sendingFail(item)
			continue
		}
		if item.Symbol == model.ETH {
			sendETH := box.SpendETH()
			cmpValue := jmath.ADD(sumWallet.Price(info.Symbol), sendETH)
			if jmath.CMP(cmpValue, senderWallet.Price(info.Symbol)) > 0 {
				isChargeAble = true
				break
			}
			sumWallet.ADD(info.Symbol, cmpValue)
		} else {
			gas := box.GasETH()
			cmpETH := jmath.ADD(sumWallet.Price(model.ETH), gas)
			if jmath.CMP(cmpETH, senderWallet.Price(model.ETH)) > 0 {
				isChargeAble = true
				break
			}
			cmpToken := jmath.ADD(sumWallet.Price(info.Symbol), item.ToPrice)
			if jmath.CMP(cmpToken, senderWallet.Price(info.Symbol)) > 0 {
				break
			}
			sumWallet.ADD(model.ETH, cmpETH)
			sumWallet.ADD(info.Symbol, cmpToken)
		}
		fixlist = append(fixlist, Fixed{
			item: item,
			box:  box,
		})
	} //for

	if len(fixlist) == 0 {
		return
	}

	nonce, err := sender.Nonce(from_master.PrivateKey)
	if err != nil {
		dbg.Red("eth.nonce", err)
		return
	}
	pending, _ := sender.XPendingNonceAt(from_master.Address)
	if nonce.NonceCount() != pending {
		dbg.RedItalic("procETHWithdarwTry.PendingNonce isDiffer (", nonce.NonceCount(), "/", pending, ")")
		return
	}

	var nonceCount uint64 = 0
	for i, fix := range fixlist {
		ntx, err := nonce.BoxTx(fix.box, nonceCount)
		if err != nil {
			sendingFail(fix.item)
			continue
		}
		stx, err := ntx.Tx()
		if err != nil {
			sendingFail(fix.item)
			continue
		}
		if err := stx.Send(); err != nil {
			sendingFail(fix.item)
			continue
		} else {
			isChargeAble = true
		}

		fixlist[i].item.Hash = stx.Hash()
		fixlist[i].item.Gas = fix.box.GasETH()
		fixlist[i].item.Timestamp = nowAt
		fixlist[i].item.YMD = nowAt.YMD()
		fixlist[i].item.State = model.TxStatePending
		sendinglist = append(sendinglist, fixlist[i])
		nonceCount++
	} //for

	for _, sendingdata := range sendinglist {
		item := sendingdata.item
		logWithdraw("(", item.UID, ") Send ", item.ToPrice, item.Symbol)
		item.UpdateDB(db)
	} //for

	return isChargeAble
}

func procETHWithdrawPending(db mongo.DATABASE, nowAt mms.MMS) bool {
	selector := mongo.Bson{"state": model.TxStatePending}
	list := model.TxETHWithdrawList{}
	db.C(inf.TXETHWithdraw).Find(selector).Sort("timestamp").All(&list)
	if len(list) == 0 {
		return false
	}

	for _, item := range list {
		checker := get_sender_x()
		if item.CancelTry {
			result, _, _, _ := checker.TransactionByHash(item.CancelHash)
			if !result.IsReceiptedByHash {
				continue
			}
			feeGas := result.GetTransactionFee()
			model.LockMemberUID(db, item.UID, func(member model.Member) {
				log := item.MakeLogWithdraw(member, model.TxStateFail)
				log.InsertDB(db, mms.Now())

				model.CoinDay{}.AddMasterGas(db, feeGas, mms.Now())

				db.C(inf.TXETHWithdraw).Remove(item.Selector())

				member.UpdateCoinDB_Legacy(db, checker)
			})

		} else {
			//result := checker.Receipt(item.Hash)
			result, _, _, _ := checker.TransactionByHash(item.Hash)
			if !result.IsReceiptedByHash {
				duration := nowAt.Sub(item.Timestamp)
				if duration.Hours() >= pendingWaitHour { //4시간경과시 실패 처리 하자
					if !item.CancelTry {
						sender := get_sender_x()
						from := inf.Master()
						to := from.Address
						speed := ecsx.GasFast

						box := sender.GasBox(
							"eth",
							from.Address,
							to,
							model.ZERO,
							speed,
						)
						if box.Error != nil {
							continue
						}
						nonce, err := sender.Nonce(from.PrivateKey)
						if err != nil {
							continue
						}
						ntx, err := nonce.BoxTx(box)
						if err != nil {
							continue
						}
						stx, err := ntx.Tx()
						if err != nil {
							continue
						}
						if err := stx.Send(); err != nil {
							continue
						}
						dbg.Red("[", item.UID, "] hash:", item.Hash, " ---- Withdraw.Cancel.Try")
						item.CancelTry = true
						item.CancelHash = stx.Hash()
						item.CancelTime = nowAt
						item.CancelYMD = nowAt.YMD()
						item.UpdateDB(db)
					}
				}
				continue
			}

			item.Gas = result.GetTransactionFee()
			if result.IsError { //fail
				model.LockMemberUID(db, item.UID, func(member model.Member) {
					logWithdraw("(", item.UID, ") r-Fail ", item.ToPrice, item.Symbol)
					log := item.MakeLogWithdraw(member, model.TxStateFail)
					log.InsertDB(db, mms.Now())

					model.CoinDay{}.AddMasterGas(db, item.Gas, mms.Now())

					db.C(inf.TXETHWithdraw).Remove(item.Selector())

					member.UpdateCoinDB_Legacy(db, checker)
				})

			} else { //success
				model.LockMemberUID(db, item.UID, func(member model.Member) {
					logWithdraw("(", item.UID, ") r-Success ", item.ToPrice, item.Symbol)
					member.Withdraw.ADD(item.Symbol, item.ToPrice)
					member.UpdateDB(db)

					model.CoinDay{}.AddWithdraw(db, item.Symbol, item.ToPrice, item.Gas, mms.Now())

					log := item.MakeLogWithdraw(member, model.TxStateSuccess)
					log.InsertDB(db, mms.Now())

					db.C(inf.TXETHWithdraw).Remove(item.Selector())

					member.UpdateCoinDB_Legacy(db, checker)
				})

			}

		}
	} //for

	return true
}
