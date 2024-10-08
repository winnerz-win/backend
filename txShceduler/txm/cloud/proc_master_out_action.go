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
	_MASTER_OUT_LIMIT = 10
)

func proc_master_send_action(
	db mongo.DATABASE,
	sender *ebcm.Sender,
	master inf.KeyPair,
	try_list model.TxETHMasterOutList,
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

	sendingFail := func(item model.TxETHMasterOut, msg string) {
		item.State = model.TxStateFail
		item.FailMessage = msg
		item.InsertLog(db)
		item.RemoveTry(db)
	}
	_ = sendingFail

	nonceCounter := base_nonce
	fixed_list := model.TxETHMasterOutList{}

	sumWallet := model.NewCoinData()
	for _, try := range try_list {
		info := tokenlist.GetSymbol(try.Symbol)
		if !info.Valid() {
			sendingFail(try, MasterFail_InvalidSymbol)
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

			coin_value = ebcm.ETHToWei(try.ToPrice) //buf_fix

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
			sendingFail(try, MasterFail_ChainErrorBox)
			continue
		}
		limit = ebcm.MMA_LimitBuffer_MasterOut(limit)

		if is_coin_transfer {
			fee_eth := gas_price.EstimateGasFeeETH(limit)
			sendETH := jmath.ADD(fee_eth, ebcm.WeiToETH(coin_value))

			cmpValue := jmath.ADD(sumWallet.Price(info.Symbol), sendETH)
			if jmath.CMP(cmpValue, senderWallet.Price(info.Symbol)) > 0 {
				sendingFail(try, MasterFail_NeedPrice)
				continue
			}
			sumWallet.ADD(info.Symbol, cmpValue)

		} else {
			fee_eth := gas_price.EstimateGasFeeETH(limit)
			cmpETH := jmath.ADD(sumWallet.Price(model.ETH), fee_eth)
			if jmath.CMP(cmpETH, senderWallet.Price(model.ETH)) > 0 {
				sendingFail(try, MasterFail_NeedPrice)
				continue
			}
			cmpToken := jmath.ADD(sumWallet.Price(info.Symbol), try.ToPrice)
			if jmath.CMP(cmpToken, senderWallet.Price(info.Symbol)) > 0 {
				sendingFail(try, MasterFail_NeedPrice)
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
		try.TryData = tx_try_data
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
		fixed.UpdateTry(db)

		logWithdraw("MasterOut Send ", fixed.ToPrice, fixed.Symbol)

		model.LogDebug.Set(
			"MASTER_OUT_TRY",
			"fixed", fixed,
		)
	} //for

	return false
}

///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

func procMasterOutTry(db mongo.DATABASE, nowAt mms.MMS) bool {
	TAG := "procMasterOutTry."
	list := model.TxETHMasterOutList{}

	selector := mongo.Bson{"state": model.TxStateNone}
	db.C(inf.TXETHMasterOutTry).Find(selector).Sort("timestamp").Limit(_MASTER_OUT_LIMIT).All(&list)
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
	master := inf.Master()

	base_nonce, err := ebcm.MMA_GetNonce(sender, master.Address, true)
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

	proc_master_send_action( //procMasterOutTry
		db,
		sender,
		master,
		list,
		base_nonce,
		gas_price,
		nowAt,
	)

	return true
}

func procMasterOutPending(db mongo.DATABASE, nowAt mms.MMS) bool {
	list := model.TxETHMasterOutList{}
	selector := mongo.Bson{"state": model.TxStatePending}
	db.C(inf.TXETHMasterOutTry).Find(selector).Sort("try_data.nonce").All(&list)
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
					retry.UpdateTry(db)
				} //for
			}
			return true
		} //if

		pending.Gas = receipt.TxFeeETH()
		if receipt.IsError() {
			logWithdraw("( MasterOutTry ) r-Fail ", pending.ToPrice, pending.Symbol)
			model.CoinDay{}.AddMasterGas(db, pending.Gas, mms.Now())

			pending.State = model.TxStateFail
			pending.FailMessage = MasterFail_ChianErrorFail
			pending.InsertLog(db)
			pending.RemoveTry(db)

		} else {
			logWithdraw("( MasterOutTry ) r-Success ", pending.ToPrice, pending.Symbol)

			model.CoinDay{}.AddWithdraw(db, pending.Symbol, pending.ToPrice, pending.Gas, mms.Now())

			pending.State = model.TxStateSuccess
			pending.InsertLog(db)
			pending.RemoveTry(db)
		}
	} //for

	pending_count, _ := db.C(inf.TXETHMasterOutTry).Find(selector).Count()
	return pending_count > 0
}
