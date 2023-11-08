package nft_winners

import (
	"txscheduler/brix/tools/cloud/ebcm"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/mms"
	"txscheduler/brix/tools/unix"
	"txscheduler/nft_winners/nwdb"
	"txscheduler/nft_winners/nwtypes"
	"txscheduler/nft_winners/rpc"
	"txscheduler/txm/model"
)

func proc_master_pending_check(db mongo.DATABASE, nowAt mms.MMS) bool {

	is_mint_pending := false
	mongo.IterForeach(
		db.C(nwdb.NftMasterPending).
			Find(nil).
			Sort("time_pending_at").
			Iter(),
		func(cnt int, pending_data nwtypes.NftMasterPending) bool {
			is_mint_pending = true

			tx_timestamp := unix.ZERO
			result_gas_info := nwtypes.ResultGasInfo{}
			var finder *ebcm.Sender
			if pending_data.Hash != "hash_skip" {
				finder = GetSender()
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
				result_gas_info.TxGasPrice = tx.GasPriceETH
				result_gas_info.Limit = tx.Limit

				tx_timestamp = tx.Timestamp

			} else {
				result_gas_info.IsSuccess = true
				result_gas_info.Hash = ""
				result_gas_info.TxGasPrice = "0"
				result_gas_info.Limit = 0
			}

			if pending_data.RemovePending(db) != nil {
				return false
			}

			current := pending_data.DATA_SEQ.Current()

			switch current.Kind {
			case nwtypes.NFT_BASE_URI:
				set_base_uri_data, _ := dbg.DecodeStruct[nwtypes.DataSetBaseURI](current.Item)

				if result_gas_info.IsSuccess {

					nft_config.BaseURI, _ = rpc.ERC721.BaseURI(finder, nft_config.Reader())
					dbg.CyanBold(
						"NEW_BASE_URI : [", set_base_uri_data.NewURI, "/", nft_config.BaseURI, "]",
						nft_config.BaseURI == set_base_uri_data.NewURI,
					)

					db.C(nwdb.NftAInfo).Update(
						mongo.Bson{"nftcontract": nft_config.NftContract},
						mongo.Bson{"$set": mongo.Bson{
							"baseuri": nft_config.BaseURI,
						}},
					)
				}

				result_dat := nwtypes.NftSetBaseURIResult{
					ReceiptCode: set_base_uri_data.ReceiptCode,
					BaseURI:     set_base_uri_data.NewURI,
					IsSucess:    result_gas_info.IsSuccess,
					Hash:        result_gas_info.Hash,
					Timestamp:   tx_timestamp,
				}
				if !set_base_uri_data.IsCallback {
					result_dat.IsSend = true
					result_dat.SendAt = tx_timestamp
					result_dat.SendMsg = "CALLBACK_SKIP"
				}
				result_dat.InsertDB(db)

			case nwtypes.NFT_MINT:
				mint_data, _ := dbg.DecodeStruct[nwtypes.DataMintTry](current.Item)

				if result_gas_info.IsSuccess {
					nwtypes.NftTokenIDS{}.UpdateOwner(
						db,
						mint_data.NftInfo.TokenId,
						mint_data.Owner.Address,
					)

					mint_data.ResultGasInfo = result_gas_info
					pending_data.DATA_SEQ.UpdateCurrent(mint_data)

					next := pending_data.DATA_SEQ.SetNext()
					switch next.Kind {
					// case nwtypes.MULTI_TRANSFER: //WINZ
					// 	multi_transfer_try := nwtypes.NftMasterTry{
					// 		ReceiptCode: mint_data.ReceiptCode,
					// 		DATA_SEQ:    pending_data.DATA_SEQ,
					// 		TimeTryAt:   unix.Now(),
					// 	}
					// 	multi_transfer_try.InsertDB(db)

					case nwtypes.DATA_NONE:
						if mint_data.Kind == nwtypes.MintKindFree { //FREE 민팅
							result_mint_data := nwtypes.ResultMintData{
								MintInfo: mint_data.ResultMint(),
								PayInfo:  nil,
							}

							action_result := nwtypes.NftActionResult{
								ReceiptCode: pending_data.ReceiptCode,
								ResultType:  pending_data.ReceiptCode.ResultType(),
								Data:        result_mint_data.Data(),
								InsertAt:    unix.Now(),
								IsSend:      false,
								SendAt:      unix.ZERO,

								TimeTryAt: pending_data.TimeTryAt,
							}
							action_result.InsertDB(db)

						} else {
							if mint_data.PayTokenInfo.IsCoin { //ETH 민팅
								multi_data, _ := dbg.DecodeStruct[nwtypes.DataMultiTransferTry](pending_data.DATA_SEQ.Prev().Item)
								result_mint_data := nwtypes.ResultMintData{
									MintInfo: mint_data.ResultMint(),
									PayInfo:  multi_data.ResultPayInfo(),
								}

								action_result := nwtypes.NftActionResult{
									ReceiptCode: pending_data.ReceiptCode,
									ResultType:  pending_data.ReceiptCode.ResultType(),
									Data:        result_mint_data.Data(),
									InsertAt:    unix.Now(),
									IsSend:      false,
									SendAt:      unix.ZERO,

									TimeTryAt: pending_data.TimeTryAt,
								}
								action_result.InsertDB(db)

							} else { //WINNERZ 민팅
								result_mint_data := nwtypes.ResultMintData{
									MintInfo: mint_data.ResultMint(),
									PayInfo:  nil,
								}

								action_result := nwtypes.NftActionResult{
									ReceiptCode: pending_data.ReceiptCode,
									ResultType:  pending_data.ReceiptCode.ResultType(),
									Data:        result_mint_data.Data(),
									InsertAt:    unix.Now(),
									IsSend:      false,
									SendAt:      unix.ZERO,

									TimeTryAt: pending_data.TimeTryAt,
								}
								action_result.InsertDB(db)
							}
						}

					default:
						dbg.RedItalic("Master.P2 ---- ", next)

					} //switch
					dbg.CyanItalic("Master.Pending OK --- HASH[", pending_data.Hash, "]", ", limit:", result_gas_info.Limit)

				} else {
					//receipt-fail ....
					dbg.RedItalic("Master.Pending Fail --- Rollback Try [", pending_data.Hash, "]", ", limit:", result_gas_info.Limit)
					mater_try := nwtypes.NftMasterTry{
						ReceiptCode: pending_data.ReceiptCode,
						DATA_SEQ:    pending_data.DATA_SEQ,
						TimeTryAt:   pending_data.TimeTryAt,
					}
					mater_try.InsertDB(db)
				}

			case nwtypes.MULTI_TRANSFER:
				multi_data, _ := dbg.DecodeStruct[nwtypes.DataMultiTransferTry](current.Item)

				if result_gas_info.IsSuccess {
					finder := GetSender()
					for _, pti := range multi_data.PriceToInfos {
						model.SyncMemberCoin(db, pti.To.Address, finder)
					} //for

					multi_data.ResultGasInfo = result_gas_info
					pending_data.DATA_SEQ.UpdateCurrent(multi_data)

					next := pending_data.DATA_SEQ.SetNext()
					switch next.Kind {
					case nwtypes.DATA_NONE:

						data_mint_try, _ := dbg.DecodeStruct[nwtypes.DataMintTry](pending_data.DATA_SEQ.Prev().Item)
						result_mint_data := nwtypes.ResultMintData{
							MintInfo: data_mint_try.ResultMint(),
							PayInfo:  multi_data.ResultPayInfo(),
						}

						action_result := nwtypes.NftActionResult{
							ReceiptCode: pending_data.ReceiptCode,
							ResultType:  pending_data.ReceiptCode.ResultType(),
							Data:        result_mint_data.Data(),
							InsertAt:    unix.Now(),
							IsSend:      false,
							SendAt:      unix.ZERO,

							TimeTryAt: pending_data.TimeTryAt,
						}
						action_result.InsertDB(db)

					default:
						dbg.RedItalic("Master.P2 ---- ", next)
					} //switch

				} else {
					//receipt-fail ....
					dbg.RedItalic("Master.Pending Fail --- Rollback Try [", pending_data.Hash, "]")
					mater_try := nwtypes.NftMasterTry{
						ReceiptCode: pending_data.ReceiptCode,
						DATA_SEQ:    pending_data.DATA_SEQ,
						TimeTryAt:   pending_data.TimeTryAt,
					}
					mater_try.InsertDB(db)
				}

			default:
				dbg.RedItalic("Master.P1 ---- ", current)
			}

			return false
		},
	)

	return is_mint_pending
}
