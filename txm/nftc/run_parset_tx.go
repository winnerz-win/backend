package nftc

import (
	"txscheduler/brix/tools/cloudx/ethwallet/abmx"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
	"txscheduler/txm/model"
)

type cTopicMap map[string]string //key  , name

func (my cTopicMap) Do(n string) bool {
	_, do := my[n]
	return do
}

var (
	topicPair = cTopicMap{
		"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef": "Transfer",
		"0x2f00e3cdd69a77be7ed215ec7b2a36784dd158f921fca79ac29deffa353fe6ee": "Mint",
		"0x5d624aa9c148153ab3446c1b154f660ee7701e549fe9b62dab7171b1c80e6fa2": "Burn",
	}
)

var (
	cms = ecsx.MethodIDDataList{
		ecsx.MakeMethodIDDataM(
			ecsx.GetMethodIDHex("setBaseURI(string)"),
			func(data string, item *ecsx.TransactionBlock) {
				cdata := ecsx.InputDataPure(
					data, abmx.NewReturns(
						abmx.String,
					),
				)

				item.NewCustomInputParse()
				cip := item.CustomInputParse
				cip["uri"] = cdata.Text(0)

				cYellow("setBaseURI(", cip["uri"], ")")
			},
		),
	}
)

func CMS() ecsx.MethodIDDataList { return cms }

func isContract(ca string) bool {
	ca = dbg.TrimToLower(ca)
	return ca == nftToken.Contract
}

func parsingTxlist(db mongo.DATABASE, number string, txlist ecsx.TransactionBlockList) {
	number64 := jmath.Int64(number)

	for _, tx := range txlist {
		if !tx.IsContract {
			continue
		}
		if !isContract(tx.ContractAddress) {
			continue
		}
		cPurple("tx : ", tx.Hash)

		receipt := Sender().ReceiptByHash(tx.Hash)

		tx.IsError = receipt.IsError()
		item := model.NftTx{
			TransactionBlock: tx,
			Logs:             receipt.Logs,
			Number:           number64,
		}
		if item.InsertDB(db) != nil {
			continue
		}

		if item.IsError {
			continue
		}

		switch tx.ContractMethod {
		case "setBaseURI":
			baseURI := tx.CustomInputParse.GetString("url")
			model.NftAset{}.UpdateBaseURI(db, baseURI)
		} //switch

		for _, log := range item.Logs {
			topic := log.Topics

			if name, do := topicPair[topic.GetName()]; do {
				switch name {
				case "Transfer":
					item.UpdateFuncNameDB(db, name)

					from := topic[1].Address() // 0 -> Mint
					to := topic[2].Address()   // 0 -> Burn
					tokenId := topic[3].Number()

					if jmath.VALUE(from) == model.ZERO || jmath.VALUE(to) == model.ZERO {
						//Mint or Burn
					} else {
						item := model.NftTxLog{
							Number:    number64,
							Hash:      tx.Hash,
							TxIndex:   tx.TxIndex,
							LogIndex:  log.LogIndex,
							Timestamp: tx.Timestamp,

							Name:    name,
							From:    from,
							To:      to,
							TokenID: tokenId,
						}
						item.InsertDB(db)

						buyEnd := model.NftBuyEnd{}.GetTokenID(db, tokenId)
						if buyEnd.Valid() {
							if buyEnd.Address == to || buyEnd.Address == from {
								buyEnd.SetLastOwner(db, to)
							}
						}
					}

				case "Mint":
					item.UpdateFuncNameDB(db, name)

					from := topic[1].Address()
					to := topic[2].Address()
					tokenId := topic[3].Number()

					tokeyType := log.Data[0].Number()

					item := model.NftTxLog{
						Number:    number64,
						Hash:      tx.Hash,
						TxIndex:   tx.TxIndex,
						LogIndex:  log.LogIndex,
						Timestamp: tx.Timestamp,

						Name:      name,
						From:      from,
						To:        to,
						TokenID:   tokenId,
						TokenType: tokeyType,
					}
					item.InsertDB(db)

				case "Burn":
					item.UpdateFuncNameDB(db, name)

					from := topic[1].Address()
					to := topic[2].Address()
					tokenId := topic[3].Number()

					tokeyType := log.Data[0].Number()

					item := model.NftTxLog{
						Number:    number64,
						Hash:      tx.Hash,
						TxIndex:   tx.TxIndex,
						LogIndex:  log.LogIndex,
						Timestamp: tx.Timestamp,

						Name:      name,
						From:      from,
						To:        to,
						TokenID:   tokenId,
						TokenType: tokeyType,
					}
					item.InsertDB(db)

					buyEnd := model.NftBuyEnd{}.GetTokenID(db, tokenId)
					if buyEnd.Valid() {
						buyEnd.BurnDB(db, tx.Hash)
					}

				} //switch
			}

		} //for

	} //for

}
