package nft_winners

import (
	"jtools/unix"
	"time"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/runtext"
	"txscheduler/nft_winners/nwdb"
	"txscheduler/nft_winners/nwtypes"
	"txscheduler/txm/model"
)

func run_user_tx_pending(rtx runtext.Runner) {
	defer dbg.PrintForce("nft_winners.run_user_tx_pending ----------  END")
	<-rtx.WaitStart()
	dbg.PrintForce("nft_winners.run_user_tx_pending ----------  START")

EXIT:
	for {
		select {
		case <-rtx.EndC():
			break EXIT
		default:
		} //select
		time.Sleep(time.Second)
		model.DB(func(db mongo.DATABASE) {
			mongo.IterForeach(
				db.C(nwdb.NftUserPending).Find(nil).Sort("time_pending_at").Iter(),
				func(cnt int, pending_data nwtypes.NftUserPending) bool {
					result_gas_info := nwtypes.ResultGasInfo{}
					finder := GetSender()
					if pending_data.Hash != "hash_skip" {
						if finder == nil {
							return true
						}
						tx, _, err := finder.TransactionByHash(pending_data.Hash)
						if err != nil {
							return false
						}
						if !tx.IsReceiptedByHash {
							return false
						}

						result_gas_info.IsSuccess = !tx.IsError
						result_gas_info.Hash = pending_data.Hash
						result_gas_info.TxGasPrice = tx.GetTransactionFee()
						result_gas_info.Limit = tx.Limit

					} else {
						result_gas_info.IsSuccess = true
						result_gas_info.Hash = ""
						result_gas_info.TxGasPrice = "0"
						result_gas_info.Limit = 0
					}

					if pending_data.RemovePending(db) != nil {
						return false
					}

					model.SyncMemberCoin(db, pending_data.From, finder)

					model.UserTransactionEnd(db, pending_data.From)

					current := pending_data.DATA_SEQ.Current()
					switch current.Kind {
					case nwtypes.MULTI_TRANSFER:
						multi_data, _ := dbg.DecodeStruct[nwtypes.DataMultiTransferTry](current.Item)
						if result_gas_info.IsSuccess {
							for _, pti := range multi_data.PriceToInfos {
								model.SyncMemberCoin(db, pti.To.Address, finder)
							} //for

							multi_data.ResultGasInfo = result_gas_info
							pending_data.DATA_SEQ.UpdateCurrent(multi_data)

							next := pending_data.DATA_SEQ.SetNext()
							switch next.Kind {
							case nwtypes.NFT_MINT:

								mint_try := nwtypes.NftMasterTry{
									ReceiptCode: pending_data.ReceiptCode,
									DATA_SEQ:    pending_data.DATA_SEQ,
									TimeTryAt:   unix.Now(),
								}
								mint_try.InsertDB(db)

							case nwtypes.NFT_TRANSFER:

								nft_transfer_try := nwtypes.NftUserTry{
									ReceiptCode: pending_data.ReceiptCode,
									DATA_SEQ:    pending_data.DATA_SEQ,
									TimeTryAt:   unix.Now(),
								}
								nft_transfer_try.InsertDB(db)

							default:
								dbg.RedItalic("User.P2 ---- ", next)

							} //switch
							dbg.CyanItalic("User.Pending OK --- HASH[", pending_data.Hash, "]", ", limit:", result_gas_info.Limit)

						} else {
							//receipt-fail ....
							dbg.RedItalic("User.Pending Fail --- Rollback Try [", pending_data.Hash, "]")
							user_try := nwtypes.NftUserTry{
								ReceiptCode: pending_data.ReceiptCode,
								DATA_SEQ:    pending_data.DATA_SEQ,
								TimeTryAt:   pending_data.TimeTryAt,
							}
							user_try.InsertDB(db)
						}

					case nwtypes.NFT_TRANSFER:
						nft_transfer_data, _ := dbg.DecodeStruct[nwtypes.DataNftTransfer](current.Item)
						if result_gas_info.IsSuccess {

							nwtypes.NftTokenIDS{}.UpdateOwner(
								db,
								nft_transfer_data.NftInfo.TokenId,
								nft_transfer_data.To.Address,
							)

							nft_transfer_data.ResultGasInfo = result_gas_info
							pending_data.DATA_SEQ.UpdateCurrent(nft_transfer_data)

							next := pending_data.DATA_SEQ.SetNext()
							switch next.Kind {
							case nwtypes.DATA_NONE:

								multi_data, _ := dbg.DecodeStruct[nwtypes.DataMultiTransferTry](pending_data.DATA_SEQ.Prev().Item)
								result_sale_data := nwtypes.ResultSaleData{
									NftTransferInfo: nft_transfer_data.ResultNftTransfer(),
									PayInfo:         multi_data.ResultPayInfo(),
								}

								action_result := nwtypes.NftActionResult{
									ReceiptCode: pending_data.ReceiptCode,
									ResultType:  pending_data.ReceiptCode.ResultType(),
									Data:        result_sale_data.Data(),
									InsertAt:    unix.Now(),
									IsSend:      false,
									SendAt:      unix.ZERO,

									TimeTryAt: pending_data.TimeTryAt,
								}
								action_result.InsertDB(db)

							default:
								dbg.RedItalic("User.P2 ---- ", next)

							} //switch
							dbg.CyanItalic("User.Pending OK --- HASH[", pending_data.Hash, "]", ", limit:", result_gas_info.Limit)

						} else {
							//receipt-fail ....
							dbg.RedItalic("User.Pending Fail --- Rollback Try [", pending_data.Hash, "]")
							user_try := nwtypes.NftUserTry{
								ReceiptCode: pending_data.ReceiptCode,
								DATA_SEQ:    pending_data.DATA_SEQ,
								TimeTryAt:   pending_data.TimeTryAt,
							}
							user_try.InsertDB(db)
						}

					default:
						dbg.RedItalic("User.P1 ---- ", current)
					} //switch

					return false
				},
			)
		})

	} //for
}
