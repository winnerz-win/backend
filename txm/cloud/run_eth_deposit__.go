package cloud

import (
	"time"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx"
	"txscheduler/brix/tools/crypt"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
	"txscheduler/brix/tools/jticker"
	"txscheduler/brix/tools/mms"
	"txscheduler/brix/tools/runtext"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

type ETHDeposit struct {
	UID      int64
	Address  string
	Symbol   string
	Contract string
	Decimal  string
	IsForce  bool
}

type ETHDepositList []ETHDeposit

var (
	ETHDepositChan = make(chan ETHDepositList, model.ChanBuffers)

	ONLY_TOKEN_TO_MASTER = false //토큰만 마스터지갑으로 전송여부 플래그
)

func runETHDepositChn(rtx runtext.Runner) {
	defer dbg.PrintForce("cloud.runETHDepositChn ---- END")
	<-rtx.WaitStart()
	dbg.PrintForce("cloud.runETHDepositChn ---- START")

	for list := range ETHDepositChan {
		model.DB(func(db mongo.DATABASE) {
			nowAt := mms.Now()
			for _, item := range list {

				if ONLY_TOKEN_TO_MASTER { //토큰만 마스터로 전송한다.
					if item.Symbol == model.ETH {
						continue
					}
				}
				txDeposit := model.TxETHDeposit{
					Key:      crypt.MakeUID256(),
					UID:      item.UID,
					Path:     model.DepositPathIncome,
					Symbol:   item.Symbol,
					Contract: item.Contract,
					Decimal:  item.Decimal,
					Address:  item.Address,
					Price:    model.ZERO,
					IsForce:  item.IsForce,
				}
				txDeposit.InsertDB(db, nowAt)
			} //for
		})

	} //for
}

func runETHDepositToMaster(rtx runtext.Runner) {
	defer dbg.PrintForce("cloud.runETHDepositToMaster ----------  END")
	<-rtx.WaitStart()
	dbg.PrintForce("cloud.runETHDepositToMaster ----------  START")

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
			item := model.TxETHDeposit{}

			//pending-check
			{
				iter := db.C(inf.TXETHDepositLog).
					Find(mongo.Bson{"state": model.TxStatePending}).
					Sort("timestamp").
					Iter()
				defer iter.Close()

				for iter.Next(&item) {
					checker := get_sender_x()
					if checker == nil {
						model.LogError.WriteLog(
							db,
							model.ErrorFinderNull,
							"runETHDepositToMaster.pendingCheck",
						)
						continue
					}

					rblock, _, _, err := checker.TransactionByHash(item.Hash)
					if err != nil {
						continue
					}
					if !rblock.IsReceiptedByHash {
						continue
					}

					if rblock.IsError {
						model.SyncCoin{}.InsertDB(db, item.Address, model.GetSyncMMS(), false)

						txFee := rblock.GetTransactionFee()
						logDeposit("(", item.UID, ") deposit-fee : ", txFee)

						model.LockMemberUID(db, item.UID, func(member model.Member) {
							member.Coin.SUB(model.ETH, txFee)

							member.UpdateDB(db)
							model.CoinSumSub(db, model.ETH, txFee)
							model.CoinDay{}.AddMemberGas(db, txFee, nowAt)

							member.UpdateCoinDB_Legacy(db, checker)
						})
					} else {
						txFee := rblock.GetTransactionFee()

						logDeposit("(", item.UID, ") deposit-fee : ", txFee, ", price:", item.Price)

						model.LockMemberUID(db, item.UID, func(member model.Member) {
							member.Coin.SUB(model.ETH, txFee)
							member.Coin.SUB(item.Symbol, item.Price)

							member.UpdateDB(db)
							model.CoinSumAction(db, func(coin *model.ConSum) {
								coin.Coin.SUB(model.ETH, txFee)
								coin.Coin.SUB(item.Symbol, item.Price)
							})

							model.CoinDay{}.AddMemberGas(db, txFee, nowAt)

							model.SyncCoin{}.InsertDB(db, item.Address, model.GetSyncMMS(), true)
							model.LogToMaster{}.InsertDB(
								db,
								inf.Master().Address,
								member,
								item,
								txFee,
								nowAt,
							)

							member.UpdateCoinDB_Legacy(db, checker)
						})
					}

					model.UserTransactionEnd(db, item.Address)
					item.RemoveDB(db)
				} //for
			} //# pending-check

			{ //wallet -> master try
				idt := model.InfoDeposit{}.Get(db)
				iter := db.C(inf.TXETHDepositLog).
					Find(mongo.Bson{"state": model.TxStateNone}).
					Sort("timestamp").
					Iter()
				defer iter.Close()

				for iter.Next(&item) {
					pendingSelector := mongo.Bson{
						"uid":   item.UID,
						"state": model.TxStatePending,
					}
					cnt, _ := db.C(inf.TXETHDepositLog).Find(pendingSelector).Count()
					if cnt > 0 {
						continue
					}

					sender := get_sender_x()
					if sender == nil {
						model.LogError.WriteLog(
							db,
							model.ErrorFinderNull,
							"runETHDepositToMaster.sender",
						)
						continue
					}

					speed := ecsx.GasFast
					from := item.Address
					to := inf.Master().Address

					model.UserTransactionStart(
						db,
						from,
						"cloud.runETHDepositToMaster",
						func() bool {
							wei := sender.Balance2(from, item.Contract)
							price := ecsx.WeiToToken(wei, item.Decimal)

							if jmath.CMP(price, 0) <= 0 {
								//보낼 잔액 없음
								item.RemoveDB(db)
								model.SyncCoin{}.InsertDB(db, from, 0, false)
								return false
							}

							if !item.IsForce {
								if !idt.IsAllow(item.Symbol, price) {
									item.RemoveDB(db)
									model.SyncCoin{}.InsertDB(db, from, 0, false)
									return false
								}
							}

							item.Price = price

							box := sender.GasBox(item.Contract, from, to, wei, speed)
							if box.Error != nil {
								return false
							}

							if item.IsGasFixed {
								box.SetLimit(item.GasLimit)
								box.SetPrice(item.GasPrice)
								logDeposit("box.fixed :", box.GasETH())
								logDeposit("box.limit&price :", item.GasLimit, item.GasPrice)
							}

							sendGas := box.GasWei()
							cGasLimit := box.Limit()
							cGasPrice := box.Price()

							//check-gas
							is_less_gas := false
							eth_gas_wei := sender.Balance(from)
							if jmath.CMP(eth_gas_wei, sendGas) < 0 { //가스비 부족
								is_less_gas = true
							}
							// if !item.IsGasFixed {
							// 	ethWEI := sender.Balance(from)
							// 	if jmath.CMP(ethWEI, sendGas) < 0 { //가스비 부족
							// 		is_less_gas = true
							// 	}
							// }

							switch item.Symbol {
							case model.ETH:
								if is_less_gas {
									//ETH
									item.RemoveDB(db)
									model.SyncCoin{}.InsertDB(db, from, 0, false)
									return false
								}

								sendWei := jmath.SUB(wei, sendGas)
								box.SetWei(sendWei)

								item.Price = ecsx.WeiToETH(sendWei) //가스비를 제외한 잔액 모두

							default: // TOKEN
								if is_less_gas {
									gasETH := ecsx.WeiToETH(sendGas)

									charger := model.TxETHCharger{
										Key:       item.Key,
										Address:   item.Address,
										Price:     gasETH,
										Timestamp: nowAt,
										YMD:       nowAt.YMD(),
										State:     model.TxStateNone,
									}
									charger.InsertDB(db)

									logDeposit("(", item.UID, ") charger.send :", gasETH)
									item.IsGasFixed = true
									item.Gas = gasETH
									item.GasLimit = cGasLimit
									item.GasPrice = cGasPrice
									item.State = model.TxStateGas
									item.UpdateDB(db, item.Timestamp)
									return false
								}

							} //switch

							wallet := inf.Wallet(item.UID)
							nonce, err := sender.Nonce(wallet.PrivateKey())

							if err != nil {
								dbg.RedItalic("deposit.NONCE :", err)
								return false
							}

							pending, _ := sender.XNonceAt(wallet.Address())
							if nonce.NonceCount() != pending {
								dbg.RedItalic("runETHDepositToMaster.XNonceAt isDiffer (", nonce.NonceCount(), "/", pending, ")")
								return false
							}

							ntx, err := nonce.BoxTx(box)
							if err != nil {
								dbg.RedItalic("deposit.NTX :", err)
								return false
							}

							stx, err := ntx.Tx()
							if err != nil {
								dbg.RedItalic("deposit.STX :", err)
								return false
							}

							if err := stx.Send(); err != nil {
								dbg.RedItalic("deposit.SEND :", err)
								model.SyncCoin{}.InsertDB(db, from, 0, false)
								return false
							}

							item.State = model.TxStatePending
							item.Hash = stx.Hash()
							item.Gas = ecsx.WeiToETH(sendGas)
							item.GasLimit = cGasLimit
							item.GasPrice = cGasPrice
							item.UpdateDB(db, nowAt)
							logDeposit("(", item.UID, ") gas :", ecsx.WeiToETH(sendGas), ", hash :", item.Hash, ", price :", item.Price)
							return true
						},
					)

				} //for
			} //wallet -> master try

		})

	} //for
}
