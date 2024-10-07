package nft_winners

import (
	"context"
	"time"
	"txscheduler/brix/tools/cloud/ebcm"
	"txscheduler/brix/tools/cloud/ebcm/abi"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
	"txscheduler/brix/tools/mms"
	"txscheduler/brix/tools/unix"
	"txscheduler/nft_winners/nwdb"
	"txscheduler/nft_winners/nwtypes"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func proc_master_send_try(db mongo.DATABASE, nowAt mms.MMS) bool {

	pendingCount, _ := db.C(nwdb.NftMasterPending).Find(nil).Count()
	if pendingCount > 0 {
		dbg.YellowItalic("[NFT] master_send_try is pending skip :", pendingCount)
		return false
	}

	master_try := nwtypes.NftMasterTry{}
	if db.C(nwdb.NftMasterTry).
		Find(nil).
		Sort("time_try_at").
		One(&master_try) != nil {
		return false
	}
	if !master_try.ReceiptCode.Valid() {
		return false
	}
	ctx := context.Background()

	sender := GetSender()
	if sender == nil {
		return false
	}

	from := inf.Master()

	nonce, err := sender.NonceAt(ctx, from.Address)
	if err != nil {
		return false
	}
	pending_nonce, _ := sender.PendingNonceAt(ctx, from.Address)
	if nonce != pending_nonce {
		dbg.RedItalic("master nonce differ (", nonce, "/", pending_nonce, ")")
		return false
	}
	time.Sleep(estimateWaitDu)

	mint_action := func() bool {
		mint_try_data, _ := dbg.DecodeStruct[nwtypes.DataMintTry](master_try.DATA_SEQ.Current().Item)

		pad_bytes := ebcm.MakePadBytesABI(
			"mint",
			abi.TypeList{
				abi.NewAddress(mint_try_data.Owner.Address),
				abi.NewUint(mint_try_data.NftInfo.TokenId),
			},
		)

		gas_limit, err := sender.EstimateGas(
			ctx,
			ebcm.MakeCallMsg(
				from.Address,
				nft_config.NftContract,
				"0",
				pad_bytes,
			),
		)
		if err != nil {
			return false
		}
		gas_limit = calc_limit(gas_limit, limit_real_pow, LIMIT_TAG_MINT)

		gas_price, err := sender.SuggestGasPrice(ctx, is_skip_tip_cap)
		if err != nil {
			return false
		}

		ntx := sender.NewTransaction(
			nonce,
			nft_config.NftContract,
			"0",
			gas_limit,
			gas_price,
			pad_bytes,
		)

		stx, err := sender.SignTx(
			ntx,
			from.PrivateKey,
		)
		if err != nil {
			return false
		}

		hash, err := sender.SendTransaction(ctx, stx)
		if err != nil {
			return false
		}

		//sender.CheckSendTxHashReceipt()
		master_try.RemoveTry(db)

		master_pending := nwtypes.NftMasterPending{
			ReceiptCode:   mint_try_data.ReceiptCode,
			DATA_SEQ:      master_try.DATA_SEQ,
			Hash:          hash,
			TimeTryAt:     master_try.TimeTryAt,
			TimePendingAt: unix.Now(),
		}
		master_pending.InsertDB(db)

		dbg.YellowItalic("MASTER_NFT_MINT_TRY[", mint_try_data.PayTokenInfoSymbol(), "] :", hash, ", limit:", gas_limit)
		return true
	} // mint_action()

	multi_transfer_action := func() bool {
		multi_transfer_try, _ := dbg.DecodeStruct[nwtypes.DataMultiTransferTry](master_try.DATA_SEQ.Current().Item)

		pay_token_info := multi_transfer_try.PayTokenInfo
		if pay_token_info.IsCoin {
			dbg.Red("MASTER_MULTI_TRANSFER_IS_COIN?????????????????????")
			return false
		}

		skip_count := 0
		receivers := []string{}
		values := []string{}
		for _, v := range multi_transfer_try.PriceToInfos {
			if jmath.CMP(v.Price, model.ZERO) <= 0 {
				skip_count++
				continue
			}
			receivers = append(receivers, v.To.Address)

			amount := pay_token_info.Wei(v.Price)
			values = append(values, amount)
		} //for

		if skip_count == len(multi_transfer_try.PriceToInfos) {
			master_try.RemoveTry(db)

			master_pending := nwtypes.NftMasterPending{
				ReceiptCode:   multi_transfer_try.ReceiptCode,
				DATA_SEQ:      master_try.DATA_SEQ,
				Hash:          "hash_skip",
				TimePendingAt: unix.Now(),
			}
			master_pending.InsertDB(db)
			dbg.YellowItalic("MASTER_MULTI_TRANSFER_TRY[", pay_token_info.Symbol, "] :", "hash_skip")
			return true
		}

		pad_bytes := ebcm.MakePadBytesABI(
			"multiTransferToken",
			abi.TypeList{
				abi.NewAddress(pay_token_info.Contract), // token
				abi.NewAddressArray(receivers...),
				abi.NewUint256Array(values),
			},
		)

		gas_limit, err := sender.EstimateGas(
			ctx,
			ebcm.MakeCallMsg(
				from.Address,
				nft_config.NftContract,
				"0",
				pad_bytes,
			),
		)
		if err != nil {
			return false
		}
		gas_limit = calc_limit(gas_limit, limit_real_pow, LIMIT_TAG_M_TRANS)

		gas_price, err := sender.SuggestGasPrice(ctx, is_skip_tip_cap)
		if err != nil {
			return false
		}

		ntx := sender.NewTransaction(
			nonce,
			nft_config.NftContract,
			"0",
			gas_limit,
			gas_price,
			pad_bytes,
		)

		stx, err := sender.SignTx(
			ntx,
			from.PrivateKey,
		)
		if err != nil {
			return false
		}

		hash, err := sender.SendTransaction(ctx, stx)
		if err != nil {
			return false
		}

		//sender.CheckSendTxHashReceipt()
		master_try.RemoveTry(db)

		master_pending := nwtypes.NftMasterPending{
			ReceiptCode:   multi_transfer_try.ReceiptCode,
			DATA_SEQ:      master_try.DATA_SEQ,
			Hash:          hash,
			TimeTryAt:     master_try.TimeTryAt,
			TimePendingAt: unix.Now(),
		}
		master_pending.InsertDB(db)
		dbg.YellowItalic("MASTER_MULTI_TRANSFER_TRY[", pay_token_info.Symbol, "] :", hash, ", limit:", gas_limit)
		return true
	} //multi_transfer_action

	set_base_uri_action := func() bool {
		set_base_uri_data, _ := dbg.DecodeStruct[nwtypes.DataSetBaseURI](master_try.DATA_SEQ.Current().Item)

		pad_bytes := ebcm.MakePadBytesABI(
			"setBaseURI",
			abi.MakeTypeList(
				abi.NewString(set_base_uri_data.NewURI),
			),
		)

		gas_limit, err := sender.EstimateGas(
			ctx,
			ebcm.MakeCallMsg(
				from.Address,
				nft_config.NftContract,
				"0",
				pad_bytes,
			),
		)
		if err != nil {
			return false
		}
		gas_limit = calc_limit(gas_limit, limit_real_pow, "")

		gas_price, err := sender.SuggestGasPrice(ctx, is_skip_tip_cap)
		if err != nil {
			return false
		}

		ntx := sender.NewTransaction(
			nonce,
			nft_config.NftContract,
			"0",
			gas_limit,
			gas_price,
			pad_bytes,
		)

		stx, err := sender.SignTx(
			ntx,
			from.PrivateKey,
		)
		if err != nil {
			return false
		}

		hash, err := sender.SendTransaction(ctx, stx)
		if err != nil {
			return false
		}

		//sender.CheckSendTxHashReceipt()
		master_try.RemoveTry(db)

		master_pending := nwtypes.NftMasterPending{
			ReceiptCode:   set_base_uri_data.ReceiptCode,
			DATA_SEQ:      master_try.DATA_SEQ,
			Hash:          hash,
			TimeTryAt:     master_try.TimeTryAt,
			TimePendingAt: unix.Now(),
		}
		master_pending.InsertDB(db)

		dbg.YellowItalic("MASTER_NFT_SET_BASE_URI[", set_base_uri_data.NewURI, "]", hash, ", limit:", gas_limit)

		return true
	}

	current := master_try.DATA_SEQ.Current()
	switch current.Kind {
	case nwtypes.NFT_MINT:
		return mint_action()

	case nwtypes.MULTI_TRANSFER:
		return multi_transfer_action()

	case nwtypes.NFT_BASE_URI:
		return set_base_uri_action()

	default:
		dbg.RedItalic("??????????????????????? ::::", current)
	}

	// switch master_try.DataKind {
	// case nwtypes.MASTER_TRY_MINT_TOKEN:
	// 	return mint_action()

	// case nwtypes.MASTER_TRY_MSEND_TOKEN:

	// case nwtypes.MASTER_TRY_MSEND_COIN:

	// case nwtypes.MASTER_TRY_MINT_COIN:
	// 	return mint_action()

	// } //switch

	return false
}
