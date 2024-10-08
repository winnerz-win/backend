package nftc

import (
	"jtools/cloud/ebcm"
	"jtools/jmath"
	"time"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func runDepositTry() {
	func() {
		for i := 0; i < 20; i++ {
			dbg.YellowBG("TEST_TEST_CASCH_BACKKKKKK")
		}
	}()

	defer dbg.PrintForce("nftc.runDepositTry ----------  END")
	dbg.PrintForce("nftc.runDepositTry ----------  START")

	for {
		model.DB(func(db mongo.DATABASE) {
			depositFail := func(depositTry model.NftDepositTry, gasFeeETH string) {
				buyTry := model.NftBuyTry{
					NftRevData: depositTry.NftRevData,
				}
				buyTry.SetFail(model.FailNftDeposit)
				buyTry.CreateTime = depositTry.CreateTime
				buyTry.CreateYMD = depositTry.CreateYMD
				result := model.NftBuyEnd{
					NftBuyTry: buyTry,
					GasFeeETH: gasFeeETH,
				}
				if result.InsertEndDB(db) == nil {
					depositTry.RemoveTryDB(db)
				}
			}

			{ //pending-check
				pendings := model.NftDepositTryList{}
				db.C(inf.NFTDepositTry).Find(mongo.Bson{"status": 1}).All(&pendings)
				for _, depositTry := range pendings {
					r, _, _ := Sender().TransactionByHash(depositTry.DepositHash)
					if !r.IsReceiptedByHash {
						continue
					}

					depositTry.GasLimit = r.Limit
					depositTry.GasPrice = r.GasUsed
					depositTry.GasFeeETH = r.GetTransactionFee()

					depositTry.DepositGasLimit = depositTry.GasLimit
					depositTry.DepositGasPrice = depositTry.GasPrice
					depositTry.DepositGasFeeETH = depositTry.GasFeeETH

					if r.IsError {
						depositFail(depositTry, r.GetTransactionFee())
						if depositTry.Address == depositTry.PayAddress {
							model.LockMemberUID(db, depositTry.UID, func(member model.Member) {
								member.UpdateCoinDB_Legacy(db)
							})
						}

					} else {
						//testCachBack(db, depositTry)

						buyTry := model.NftBuyTry{
							NftRevData: depositTry.NftRevData,
						}
						if buyTry.InsertTryDB(db) == nil {
							depositTry.RemoveTryDB(db)
						}
						if depositTry.Address == depositTry.PayAddress {
							model.LockMemberUID(db, depositTry.UID, func(member model.Member) {
								member.Withdraw.ADD(depositTry.PaySymbol, depositTry.PayPrice)
								member.UpdateDB(db)

								member.UpdateCoinDB_Legacy(db)
							})
						}
					}
				} //for
			} //pending-check

			tokenlist := inf.TokenList()

			iter := db.C(inf.NFTDepositTry).
				Find(mongo.Bson{"status": 0}).
				Sort("create_at").
				Iter()
			defer iter.Close()

			depositTry := model.NftDepositTry{}
			for iter.Next(&depositTry) {
				if depositTry.IsPayFree {
					buyTry := model.NftBuyTry{
						NftRevData: depositTry.NftRevData,
					}
					if buyTry.InsertTryDB(db) == nil {
						depositTry.RemoveTryDB(db)
					}
					continue
				}

				ethPrice := Finder().GetCoinPrice(depositTry.Address)

				token := tokenlist.GetSymbol(depositTry.PaySymbol)
				if token.Symbol == model.ETH {
					// ntx, err := Sender().EthTransferNTX(
					// 	depositTry.PrivateKey(),
					// 	depositAddress(),
					// 	ebcm.ETHToWei(depositTry.PayPrice),
					// 	gasSpeed,
					// 	&depositTry.Snap,
					// )
					ntx, err := TransferEtherNTX(
						depositTry.PrivateKey(),
						depositAddress(),
						ebcm.ETHToWei(depositTry.PayPrice),
					)

					if err != nil {
						depositFail(depositTry, model.ZERO)
						continue
					}
					gasETH := ntx.GasFeeETH()
					needETH := jmath.ADD(gasETH, depositTry.PayPrice)
					if jmath.CMP(ethPrice, needETH) < 0 {
						depositFail(depositTry, model.ZERO)
						continue
					}

					h, err := TransferNTX_Send(ntx)
					if err != nil {
						depositFail(depositTry, model.ZERO)
						continue
					}
					depositTry.DepositHash = h
					depositTry.Status = 1
					depositTry.UpdateTryDB(db)

				} else {
					if jmath.IsUnderZero(ethPrice) {
						depositFail(depositTry, model.ZERO)
						continue
					}
					tkPrice := Finder().Price(depositTry.Address, token.Contract, token.Decimal)
					if jmath.CMP(tkPrice, depositTry.PayPrice) < 0 {
						depositFail(depositTry, model.ZERO)
						continue
					}

					// ts := Sender().TSender(token.Contract)
					// ntx, err := ts.TransferFuncNTX(
					// 	depositTry.PayPrivateKey(),
					// 	ebcm.TransferPadBytes(
					// 		depositAddress(),
					// 		ebcm.TokenToWei(depositTry.PayPrice, token.Decimal),
					// 	),
					// 	"0",
					// 	gasSpeed,
					// 	&depositTry.Snap,
					// )

					ntx, err := TransferTokenNTX(
						token.Contract,
						depositTry.PayPrivateKey(),
						depositAddress(),
						ebcm.TokenToWei(depositTry.PayPrice, token.Decimal),
					)

					if err != nil {
						depositFail(depositTry, model.ZERO)
						continue
					}
					depositTry.Snap = ntx.SnapShot()

					gasETH := ntx.GasFeeETH()
					if jmath.CMP(ethPrice, gasETH) < 0 {
						depositFail(depositTry, model.ZERO)
						continue
					}

					//h, _, err := ts.TransferFuncSEND(ntx)
					h, err := TransferNTX_Send(ntx)
					if err != nil {
						depositFail(depositTry, model.ZERO)
						continue
					}
					depositTry.DepositHash = h
					depositTry.Status = 1
					depositTry.UpdateTryDB(db)

				}

			} //for
		})

		time.Sleep(time.Second)
	} //for
}
