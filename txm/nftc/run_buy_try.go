package nftc

import (
	"time"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func runBuyTry() {
	defer dbg.PrintForce("nftc.runBuyTry ----------  END")
	dbg.PrintForce("nftc.runBuyTry ----------  START")

	for {
		model.DB(func(db mongo.DATABASE) {

			{ //pending-check
				pendings := model.NftBuyTryList{}
				selector := mongo.Bson{
					"status": 1,
				}
				db.C(inf.NFTBuyTry).Find(selector).All(&pendings)
				for _, try := range pendings {
					r, _, _, _ := Sender().TransactionByHash(try.Hash)
					if !r.IsReceiptedByHash {
						continue
					}
					model.LockMember(db, try.Address, func(member model.Member) {
						member.UpdateCoinDB_Legacy(db, Sender())
					})

					feeGas := r.GetTransactionFee()

					try.GasLimit = r.Limit
					try.GasPrice = r.GasUsed
					try.GasFeeETH = feeGas

					result := model.NftBuyEnd{
						GasFeeETH: feeGas,
					}
					if r.IsError {
						try.SetFail(model.FailNftTxResult)

					} else {
						try.Status = 200
						NFT{}.TokenMetaData(try.TokenId, func(meta model.NftMetaData) {
							result.Meta = meta
						})
						NFT{}.TokenURI(try.TokenId, func(uri string) {
							result.TokenURI = uri
						})
					}
					result.NftBuyTry = try
					result.LastOwner = try.Address
					result.IsBurn = false
					if result.InsertEndDB(db) == nil {
						try.RemoveTryDB(db)
					}
				} //for
			}

			gasETH := model.ZERO
			if cnt, _ := db.C(inf.NFTBuyTry).Find(mongo.Bson{"status": 0}).Count(); cnt > 0 {
				gasETH = Sender().CoinPrice(nftToken.Address)
			}
			if jmath.IsUnderZero(gasETH) {
				return
			}

			iter := db.C(inf.NFTBuyTry).
				Find(mongo.Bson{"status": 0}).
				Sort("create_at").
				Iter()
			defer iter.Close()

			try := model.NftBuyTry{}
			for iter.Next(&try) {

				pendingSelector := mongo.Bson{
					"address": try.Address,
					"status":  1,
				}
				cnt, _ := db.C(inf.NFTBuyTry).Find(pendingSelector).Count()
				if cnt > 0 {
					dbg.Yellow("[NFT] pending Wait :", try.Address)
					continue
				}

				try.GasLimit = 0
				try.GasPrice = model.ZERO
				try.GasFeeETH = model.ZERO
				wr := NFT{}.Mint2(
					try.TokenId,
					try.Address,
					try.TokenType,
					&try,
				)
				if wr.Error != nil {
					break
				}

				try.Status = 1
				try.Hash = wr.Hash
				try.UpdateTryDB(db)

			} //for

		})
		time.Sleep(time.Second)
	} //for
}
