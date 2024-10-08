package nft_winners

import (
	"context"
	"jtools/cloud/ebcm"
	"jtools/cloud/ebcm/abi"
	"jtools/jmath"
	"jtools/unix"
	"time"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/runtext"
	"txscheduler/nft_winners/nwdb"
	"txscheduler/nft_winners/nwtypes"
	"txscheduler/txm/model"
)

// EstimateTxFee : min , need , err
func EstimateTxFee(
	db mongo.DATABASE,
	TAG string,
	sender *ebcm.Sender,
	from string,
	to string,
	pad_bytes ebcm.PADBYTES,
	wei string,
	limit_tag string,
) (string, string, error) { //min , max(need)

	ctx := context.Background()
	gas_limit, err := sender.EstimateGas(
		ctx,
		ebcm.MakeCallMsg(
			from,
			to,
			wei,
			pad_bytes,
		),
	)
	if err != nil {
		return "0", "0", err
	}
	min_limit := gas_limit
	max_limit := calc_limit(gas_limit, limit_real_pow, limit_tag)

	gas_price, err := sender.SuggestGasPrice(ctx, is_skip_tip_cap)
	if err != nil {
		return "0", "0", err
	}
	if db != nil {
		gas_price = model.CALC_GAS_PRICE(db, gas_price)
	}

	dbg.Cyan("[", TAG, "]EstimateTxFee")
	dbg.Cyan("min_limit :", min_limit, ", max_limit :", max_limit)
	dbg.Cyan("gas :", gas_price.Gas, ", tip :", gas_price.Tip)
	dbg.Cyan("--------------------------------------------------")

	//gas_price.AddGWEI_ALL(ebcm.GWEI("1"))

	tx_min_gas_wei := gas_price.EstimateGasFeeWEI(min_limit)
	tx_need_gas_wei := gas_price.EstimateGasFeeWEI(max_limit)

	return tx_min_gas_wei, tx_need_gas_wei, nil
}

func getPadBytes_MultiTransferETH(pairs ...any) ebcm.PADBYTES {
	receivers := []string{}
	values := []string{}
	for i := 0; i < len(pairs); i += 2 {
		receiver := pairs[i].(string)
		value := pairs[i+1].(string)
		if jmath.CMP(value, model.ZERO) <= 0 {
			continue
		}
		receivers = append(receivers, receiver)
		values = append(values, value)
	} //for

	pad_bytes := ebcm.MakePadBytesABI(
		"multiTransferETH",
		abi.TypeList{
			abi.NewAddressArray(receivers...),
			abi.NewUint256Array(values),
		},
	)
	return pad_bytes
}

/////////////////////////////////////////////////////////////////

func run_user_tx_try(rtx runtext.Runner) {
	defer dbg.PrintForce("nft_winners.run_user_tx_try ----------  END")
	<-rtx.WaitStart()
	dbg.PrintForce("nft_winners.run_user_tx_try ----------  START")

EXIT:
	for {
		select {
		case <-rtx.EndC():
			break EXIT
		default:
		} //select
		time.Sleep(time.Millisecond * 100)

		model.DB(func(db mongo.DATABASE) {
			c := db.C(nwdb.NftUserTry)

			mongo.IterForeach(
				c.Find(nil).Sort("time_try_at").Iter(),
				func(cnt int, user_try nwtypes.NftUserTry) bool {

					current := user_try.DATA_SEQ.Current()

					from := nwtypes.WalletInfo{}
					tx_to := ""
					eth_wei := "0"

					var void_data interface{}
					switch current.Kind {
					case nwtypes.MULTI_TRANSFER:
						multi_transfer_try, _ := dbg.DecodeStruct[nwtypes.DataMultiTransferTry](current.Item)
						from = multi_transfer_try.PayFrom
						tx_to = nft_config.NftContract

						void_data = multi_transfer_try

					case nwtypes.NFT_TRANSFER:
						nft_transfer, _ := dbg.DecodeStruct[nwtypes.DataNftTransfer](current.Item)
						from = nft_transfer.From
						tx_to = nft_config.NftContract

						void_data = nft_transfer

					default:
						dbg.Red("user_tx_try : ", current.Kind, "??????????????")
						return false
					} //switch

					if !model.UserTransactionStart(db, from.Address, "nft_winners.run_user_tx_try") {
						return false
					}

					if pendingCount, _ := db.C(nwdb.NftUserPending).Find(
						mongo.Bson{"from": from.Address},
					).Count(); pendingCount > 0 {
						dbg.YellowItalic("[NFT] user(", current.Kind, ") is pending skip :", pendingCount)
						return false
					}

					sender := GetSender()
					if sender == nil {
						return false
					}
					ctx := context.Background()

					nonce, err := sender.NonceAt(ctx, from.Address)
					if err != nil {
						return false
					}
					pending_nonce, _ := sender.PendingNonceAt(ctx, from.Address)
					if nonce != pending_nonce {
						dbg.RedItalic("user[", from.UID, "] nonce differ (", nonce, "/", pending_nonce, ")")
						return false
					}
					time.Sleep(estimateWaitDu)

					limit_tag := ""
					sending_msg := ""
					pad_bytes := ebcm.PADBYTES{}
					switch current.Kind {
					case nwtypes.MULTI_TRANSFER:
						limit_tag = LIMIT_TAG_M_TRANS

						multi_transfer_try := void_data.(nwtypes.DataMultiTransferTry)

						pay_token_info := multi_transfer_try.PayTokenInfo
						if !pay_token_info.IsCoin {
							dbg.Red("user_tx_try : is not coin [", pay_token_info.Symbol, "]")
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
							eth_wei = jmath.ADD(eth_wei, amount)
						} //for

						if skip_count == len(multi_transfer_try.PriceToInfos) {
							user_try.RemoveTry(db)

							user_pending := nwtypes.NftUserPending{
								ReceiptCode: user_try.ReceiptCode,
								DATA_SEQ:    user_try.DATA_SEQ,
							}
							user_pending.InsertDB(db)
							dbg.YellowItalic("USER_MULTI_TRANSFER_TRY[", pay_token_info.Symbol, "] :", "hash_skip")
							return false
						}

						pad_bytes = ebcm.MakePadBytesABI(
							"multiTransferETH",
							abi.TypeList{
								abi.NewAddressArray(receivers...),
								abi.NewUint256Array(values),
							},
						)
						sending_msg = dbg.Cat("USER_MULTI_TRANSFER_TRY[", pay_token_info.Symbol, "]")

					case nwtypes.NFT_TRANSFER:
						nft_transfer := void_data.(nwtypes.DataNftTransfer)

						pad_bytes = ebcm.MakePadBytesABI(
							"transferFrom",
							abi.TypeList{
								abi.NewAddress(from.Address),
								abi.NewAddress(nft_transfer.To.Address),
								abi.NewUint256(nft_transfer.NftInfo.TokenId),
							},
						)
						eth_wei = "0"
						sending_msg = dbg.Cat("USER_NFT_TRANSFER_TRY[", nft_transfer.NftInfo.TokenId, "]")

					} //switch

					gas_limit, err := sender.EstimateGas(
						ctx,
						ebcm.MakeCallMsg(
							from.Address,
							tx_to,
							eth_wei,
							pad_bytes,
						),
					)
					if err != nil {
						return false
					}
					gas_limit = calc_limit(gas_limit, limit_real_pow, limit_tag)

					gas_price, err := sender.SuggestGasPrice(ctx, is_skip_tip_cap)
					if err != nil {
						return false
					}
					gas_price = model.CALC_GAS_PRICE(db, gas_price)

					ntx := sender.NewTransaction(
						nonce,
						tx_to,
						eth_wei,
						gas_limit,
						gas_price,
						pad_bytes,
					)

					stx, err := sender.SignTx(
						ntx,
						from.UserPrivateKey(),
					)
					if err != nil {
						return false
					}

					hash, err := sender.SendTransaction(ctx, stx)
					if err != nil {
						return false
					}

					user_try.RemoveTry(db)

					user_pending := nwtypes.NftUserPending{
						ReceiptCode: user_try.ReceiptCode,
						From:        from.Address,

						DATA_SEQ: user_try.DATA_SEQ,
						Hash:     hash,

						TimeTryAt:     user_try.TimeTryAt,
						TimePendingAt: unix.Now(),
					}
					user_pending.InsertDB(db)
					dbg.YellowItalic(sending_msg, ":", hash, ", limit:", gas_limit)

					return false
				},
			)

			// switch try.DataKind {
			// case nwtypes.MASTER_TRY_MINT_COIN:
			// } //switch

		})
	} //for
}
