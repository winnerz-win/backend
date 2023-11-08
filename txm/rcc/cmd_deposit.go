package rcc

import (
	"txscheduler/brix/tools/console"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/jmath"
	"txscheduler/txm/cloud"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func init() {
	cmdDepositMaster()
}

func cmdDepositMaster() {
	console.AppendCmd(
		"deposit_master",
		"deposit_master [symbol] [uid , name , address] / [all]",
		false,
		func(ps []string) {
			symbol := ps[0]

			if inf.ValidSymbol(symbol) == false {
				console.Log("invalid symbol :", symbol)
				return
			}
			token := inf.TokenList().GetSymbol(symbol)

			minValue := "0"
			if symbol == model.ETH {
				minValue = "0.002"
			}

			isAll := ps[1] == "all"
			if isAll == true {
				model.DB(func(db mongo.DATABASE) {
					depositList := cloud.ETHDepositList{}

					iter := db.C(inf.COLMember).Find(nil).Iter()
					member := model.Member{}
					for iter.Next(&member) {
						if jmath.CMP(member.Coin.Price(symbol), minValue) >= 0 {
							depositList = append(depositList, cloud.ETHDeposit{
								UID:      member.UID,
								Address:  member.Address,
								Symbol:   token.Symbol,
								Contract: token.Contract,
								Decimal:  token.Decimal,
								IsForce:  true,
							})
						}
					} //for

					console.Log("depositToMaster :", len(depositList))
					if len(depositList) > 0 {
						cloud.ETHDepositChan <- depositList
					}
				})
			} else {
				model.DB(func(db mongo.DATABASE) {
					member := getMemberKey(db, ps[1])
					if member.Valid() == false {
						console.Log("not found member :", ps[1])
						return
					}

					depositList := cloud.ETHDepositList{}
					if jmath.CMP(member.Coin.Price(symbol), minValue) >= 0 {
						depositList = append(depositList, cloud.ETHDeposit{
							UID:      member.UID,
							Address:  member.Address,
							Symbol:   token.Symbol,
							Contract: token.Contract,
							Decimal:  token.Decimal,
							IsForce:  true,
						})
					}

					console.Log("depositToMaster :", len(depositList))
					if len(depositList) > 0 {
						cloud.ETHDepositChan <- depositList
					}
				})
			}
		},
	)
}
