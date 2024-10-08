package rcc

import (
	"jtools/cloud/ebcm"
	"jtools/jmath"
	"txscheduler/brix/tools/console"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func init() {
	cmdFinder()
	cmdSyncAddress()
}

func getMemberKey(db mongo.DATABASE, key string) model.Member {
	if ebcm.IsAddress(key) {
		return model.LoadMemberAddress(db, key)
	}
	if jmath.IsNum(key) == false {
		return model.LoadMemberName(db, key)
	}
	return model.LoadMember(db, jmath.Int64(key))
}

func cmdFinder() {
	console.AppendCmd(
		"find",
		"find [uid / name / address]",
		false,
		func(ps []string) {
			model.DB(func(db mongo.DATABASE) {
				member := getMemberKey(db, ps[0])
				if member.Valid() == false {
					console.Log("not found member :", ps[0])
					return
				}
				console.Log(member)
			})
		},
	)

	console.AppendCmd(
		"config",
		"config",
		true,
		func(ps []string) {
			v := inf.Config().View()
			console.Log(v)
		},
	)

	console.AppendCmd(
		"master.info",
		"master.info",
		true,
		func(ps []string) {
			mainnet := inf.Mainnet()
			masterAddress := inf.Master().Address
			chargerAddress := inf.Charger().Address

			masterPrice := model.NewCoinData()
			for _, token := range inf.TokenList() {
				finder := Finder(mainnet, inf.InfuraKey())
				wei := finder.TokenBalance(masterAddress, token.Contract)

				price := ebcm.WeiToToken(wei, token.Decimal)
				masterPrice.ADD(token.Symbol, price)
			}

			finder := Finder(mainnet, inf.InfuraKey())
			wei := finder.GetCoinBalance(chargerAddress)
			chargerPrice := model.NewCoinData()
			chargerPrice.ADD(model.ETH, ebcm.WeiToToken(wei, "18"))

			memberCount := 0
			model.DB(func(db mongo.DATABASE) {
				memberCount, _ = db.C(inf.COLMember).Count()
			})

			console.Log("mainnet         :", mainnet)
			console.Atap()
			console.Log("master_address  :", masterAddress)
			console.Log(masterPrice)
			console.Atap()
			console.Log("charger_address :", chargerAddress)
			console.Log(chargerPrice)
			console.Atap()
			console.Log("member_count    :", memberCount)
		},
	)
}

func cmdSyncAddress() {
	console.AppendCmd(
		"coinsync",
		"coinsync [uid / name / address | all]",
		false,
		func(ps []string) {
			isAll := false
			if ps[0] == "all" {
				isAll = true
			}

			model.DB(func(db mongo.DATABASE) {
				if isAll == false {
					member := getMemberKey(db, ps[0])
					if member.Valid() == false {
						console.Log("not found member :", ps[0])
						return
					}

					model.SyncCoin{}.InsertDB(db, member.Address, 0, false)
					console.Log("[coinsync] ", member.UID, " , ", member.Address)
					console.Log("reserved job success.")

				} else {
					iter := db.C(inf.COLMember).Find(nil).Iter()

					count := 0
					member := model.Member{}
					for iter.Next(&member) {
						model.SyncCoin{}.InsertDB(db, member.Address, 0, false)
						count++
					} //for

					console.Log("[coinsync] Count :", count)

				}
			})
		},
	)
}
