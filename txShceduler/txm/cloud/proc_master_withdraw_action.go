package cloud

import (
	"context"
	"jtools/cc"
	"jtools/cloud/ebcm"
	"jtools/jmath"
	"jtools/mms"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

const (
	_MASTER_ETH_LIMIT = 10
)

func proc_eth_withdraw_action(
	db mongo.DATABASE,
	sender *ebcm.Sender,
	master inf.KeyPair,
	list model.TxETHWithdrawList,
	base_nonce uint64,
	gas_price ebcm.GasPrice,
	nowAt mms.MMS,
) bool {
	senderWallet := model.NewCoinData()
	tokenlist := inf.TokenList()
	for _, token := range tokenlist {
		wei := get_finder().TokenBalance(master.Address, token.Contract)
		senderWallet[token.Symbol] = ebcm.WeiToToken(wei, token.Decimal)
	} //for

	sendingFail := func(item model.TxETHWithdraw) {
		model.LockMemberUID(db, item.UID, func(member model.Member) {
			log := item.MakeLogWithdraw(member, model.TxStateFail)
			log.InsertDB(db, mms.Now())

			db.C(inf.TXETHWithdraw).Remove(item.Selector())
		})
	}
	_ = sendingFail

	nonceCounter := base_nonce
	fixed_list := model.TxETHWithdrawList{}

	sumWallet := model.NewCoinData()
	for _, try := range list {
		info := tokenlist.GetSymbol(try.Symbol)
		if !info.Valid() {
			sendingFail(try)
			continue
		}

		var pad_bytes ebcm.PADBYTES
		coin_value := model.ZERO
		send_to := ""
		is_coin_transfer := false
		if !ebcm.IsAddress(info.Contract) {
			is_coin_transfer = true
			pad_bytes = ebcm.PadByteETH()
			send_to = try.ToAddress
			coin_value = ebcm.ETHToWei(try.ToPrice)

		} else {
			pad_bytes = ebcm.PadByteTransfer(
				try.ToAddress,
				try.Wei(),
			)
			send_to = info.Contract
		}

		limit, err := sender.EstimateGas(
			context.Background(),
			ebcm.MakeCallMsg(
				master.Address,
				send_to,
				coin_value,
				pad_bytes,
			),
		)
		if err != nil {
			sendingFail(try)
			continue
		}
		limit = ebcm.MMA_LimitBuffer_MasterOut(limit)

		if is_coin_transfer {
			fee_eth := gas_price.EstimateGasFeeETH(limit)
			sendETH := jmath.ADD(fee_eth, ebcm.WeiToETH(coin_value))

			cmpValue := jmath.ADD(sumWallet.Price(info.Symbol), sendETH)
			if jmath.CMP(cmpValue, senderWallet.Price(info.Symbol)) > 0 {
				sendingFail(try)
				continue
			}
			sumWallet.ADD(info.Symbol, cmpValue)

		} else {
			fee_eth := gas_price.EstimateGasFeeETH(limit)
			cmpETH := jmath.ADD(sumWallet.Price(model.ETH), fee_eth)
			if jmath.CMP(cmpETH, senderWallet.Price(model.ETH)) > 0 {
				sendingFail(try)
				continue
			}
			cmpToken := jmath.ADD(sumWallet.Price(info.Symbol), try.ToPrice)
			if jmath.CMP(cmpToken, senderWallet.Price(info.Symbol)) > 0 {
				sendingFail(try)
				continue
			}
			sumWallet.ADD(model.ETH, cmpETH)
			sumWallet.ADD(info.Symbol, cmpToken)
		}

		tx_try_data, stx := model.MakeTxTryData(
			sender,
			master.PrivateKey,

			nonceCounter,
			send_to,
			coin_value,
			limit,
			gas_price,
			pad_bytes,
		)
		try.TryData = &tx_try_data
		try.TrySTX = stx

		fixed_list = append(fixed_list, try)

		nonceCounter++
	} //for

	if len(fixed_list) == 0 {
		return true
	}

	for _, fixed := range fixed_list {
		hash, _ := sender.SendTransaction(
			context.Background(),
			fixed.TrySTX,
		)
		fixed.Hash = hash
		fixed.Timestamp = nowAt
		fixed.YMD = nowAt.YMD()
		fixed.State = model.TxStatePending

		logWithdraw("(", fixed.UID, ") Send ", fixed.ToPrice, fixed.Symbol)
		fixed.UpdateDB(db)

		model.LogDebug.Set(
			"ETH_WITHDRAW_TRY",
			"fixed", fixed,
		)

	} //for
	return false
}

func procETHWithdarwTry(db mongo.DATABASE, nowAt mms.MMS) (isChargeAble bool) {
	TAG := "procETHWithdarwTry."

	isChargeAble = false

	list := model.TxETHWithdrawList{}

	selector := mongo.Bson{"state": model.TxStateNone}
	db.C(inf.TXETHWithdraw).Find(selector).Sort("timestamp").Limit(_MASTER_ETH_LIMIT).All(&list)

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

	master := inf.Master()

	nonce, err := ebcm.MMA_GetNonce(sender, master.Address, true)
	if err != nil {
		cc.RedItalic(TAG, "MMA_GetNonce :", err)
		return false
	}

	gas_price, err := sender.SuggestGasPrice(context.Background(), true)
	if err != nil {
		cc.RedItalic(TAG, "SuggestGasPrice :", err)
		return false
	}
	gas_price = CALC_GAS_PRICE(db, gas_price)

	proc_eth_withdraw_action(
		db,
		sender,
		master,
		list,
		nonce,
		gas_price,
		nowAt,
	)

	return true
}

func procETHWithdrawPending(db mongo.DATABASE, nowAt mms.MMS) bool {
	selector := mongo.Bson{"state": model.TxStatePending}
	list := model.TxETHWithdrawList{}
	db.C(inf.TXETHWithdraw).Find(selector).Sort("try_data.nonce").All(&list)
	if len(list) == 0 {
		return false
	}

	master := inf.Master()
	dt := model.GetGAS(db).GetLtTransferPendingMin()

	for i, pending := range list {
		sender := get_sender_x()
		receipt := sender.ReceiptByHash(pending.Hash)

		if receipt.IsNotFound {
			if pending.TryData.IsTimeOver(dt) {
				for j := i; j < len(list); j++ {
					retry := list[j]

					stx := retry.TryData.STX(
						sender,
						master.PrivateKey,
					)
					hash, _ := sender.SendTransaction(context.Background(), stx)
					retry.Hash = hash
					retry.Timestamp = nowAt
					retry.YMD = nowAt.YMD()
					retry.UpdateDB(db)
				} //for
			}
			return true
		} //if

		pending.Gas = receipt.TxFeeETH()
		if receipt.IsError() {
			model.LockMemberUID(db, pending.UID, func(member model.Member) {
				logWithdraw("(", pending.UID, ") r-Fail ", pending.ToPrice, pending.Symbol)
				log := pending.MakeLogWithdraw(member, model.TxStateFail)
				log.InsertDB(db, mms.Now())

				model.CoinDay{}.AddMasterGas(db, pending.Gas, mms.Now())

				db.C(inf.TXETHWithdraw).Remove(pending.Selector())

				member.UpdateCoinDB_Legacy(db)
			})

		} else {
			model.LockMemberUID(db, pending.UID, func(member model.Member) {
				logWithdraw("(", pending.UID, ") r-Success ", pending.ToPrice, pending.Symbol)
				member.Withdraw.ADD(pending.Symbol, pending.ToPrice)
				member.UpdateDB(db)

				model.CoinDay{}.AddWithdraw(db, pending.Symbol, pending.ToPrice, pending.Gas, mms.Now())

				log := pending.MakeLogWithdraw(member, model.TxStateSuccess)
				log.InsertDB(db, mms.Now())

				db.C(inf.TXETHWithdraw).Remove(pending.Selector())

				member.UpdateCoinDB_Legacy(db)
			})
		}

	} //for

	return true
}
