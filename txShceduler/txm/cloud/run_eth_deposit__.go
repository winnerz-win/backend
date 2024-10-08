package cloud

import (
	"context"
	"jtools/cloud/ebcm"
	"jtools/jmath"
	"jtools/mms"
	"time"
	"txscheduler/brix/tools/crypt"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jticker"
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
					checker := get_finder()
					if checker == nil {
						model.LogError.WriteLog(
							db,
							model.ErrorFinderNull,
							"runETHDepositToMaster.pendingCheck",
						)
						continue
					}

					rblock, _, err := checker.TransactionByHash(item.Hash)
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

							member.UpdateCoinDB_Legacy(db)
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

							member.UpdateCoinDB_Legacy(db)
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

					// speed := ecsx.GasFast
					from := item.Address
					master_address := inf.Master().Address

					model.UserTransactionStart(
						db,
						from,
						"cloud.runETHDepositToMaster",
						func() bool {
							send_wei := get_finder().TokenBalance(from, item.Contract)
							send_price := ebcm.WeiToToken(send_wei, item.Decimal)

							if jmath.CMP(send_price, 0) <= 0 {
								//보낼 잔액 없음
								item.RemoveDB(db)
								model.SyncCoin{}.InsertDB(db, from, 0, false)
								return false
							}

							if !item.IsForce {
								if !idt.IsAllow(item.Symbol, send_price) {
									item.RemoveDB(db)
									model.SyncCoin{}.InsertDB(db, from, 0, false)
									return false
								}
							}

							nonce, err := ebcm.MMA_GetNonce(
								sender,
								from,
								true,
							)
							if err != nil {
								item.RemoveDB(db)
								return false
							}

							item.Price = send_price

							var pad_bytes ebcm.PADBYTES
							to := ""
							coin_value := model.ZERO
							is_coin_transfer := false

							switch item.Symbol {
							case model.ETH:
								pad_bytes = ebcm.PadByteETH()
								to = master_address
								coin_value = send_wei
								is_coin_transfer = true

							default: //TOKEN
								pad_bytes = ebcm.PadByteTransfer(
									master_address,
									send_wei,
								)
								to = item.Contract
							}

							try_limit := uint64(0)
							try_gas_price := ebcm.GasPrice{}
							if item.IsGasFixed {
								try_limit = jmath.Uint64(item.GasLimit)
								try_gas_price = ebcm.MakeGasPrice(item.GasPrice)

							} else {
								limit, err := sender.EstimateGas(
									context.Background(),
									ebcm.MakeCallMsg(
										from,
										to,
										coin_value,
										pad_bytes,
									),
								)
								if err != nil {
									return false
								}
								try_limit = limit

								gas_price, err := sender.SuggestGasPrice(
									context.Background(),
									true,
								)
								if err != nil {
									return false
								}
								try_gas_price = gas_price
							}

							//check-gas
							fee_wei := try_gas_price.EstimateGasFeeWEI(try_limit)
							fee_eth := try_gas_price.EstimateGasFeeETH(try_limit)
							if is_coin_transfer {
								if jmath.CMP(send_price, fee_eth) <= 0 {

									//ETH
									item.RemoveDB(db)
									model.SyncCoin{}.InsertDB(db, from, 0, false)
									return false
								}
							} else {
								from_coin_price := sender.Price(from)
								if jmath.CMP(from_coin_price, fee_eth) < 0 {

									gasETH := jmath.SUB(fee_eth, from_coin_price)

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
									item.GasLimit = jmath.VALUE(try_limit)
									item.GasPrice = try_gas_price.GET_GAS_WEI()
									item.State = model.TxStateGas
									item.UpdateDB(db, item.Timestamp)
									return false
								}
							}

							wallet := inf.Wallet(item.UID)

							ntx := sender.NewTransaction(
								nonce,
								to,
								coin_value,
								try_limit,
								try_gas_price,
								pad_bytes,
							)
							stx, err := sender.SignTx(ntx, wallet.PrivateKey())
							if err != nil {
								dbg.RedItalic("deposit.SignTx ::", err)
								return false
							}

							hash, err := sender.SendTransaction(
								context.Background(),
								stx,
							)
							if err != nil {
								dbg.RedItalic("deposit.SendTransaction :", err)
								model.SyncCoin{}.InsertDB(db, from, 0, false)
								return false
							}

							item.State = model.TxStatePending
							item.Hash = hash
							item.Gas = try_gas_price.GET_GAS_ETH()
							item.GasLimit = jmath.VALUE(try_limit)
							item.GasPrice = fee_wei
							item.UpdateDB(db, nowAt)
							logDeposit("(", item.UID, ") fee_eth :", fee_eth, ", hash :", item.Hash, ", price :", item.Price)
							return true
						},
					)

				} //for
			} //wallet -> master try

		})

	} //for
}
