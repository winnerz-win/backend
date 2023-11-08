package model

import (
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/jmath"
	"txscheduler/txm/inf"
)

//InfoDeposit :
type InfoDeposit struct {
	Key       string   `bson:"key" json:"key"`
	Coin      CoinData `bson:"coin" json:"coin"`
	BaseValue string   `bson:"base_value" json:"base_value"`
}

//TagString :
func (InfoDeposit) TagString() []string {
	return []string{
		"coin", "[Symbol]마스터로 보내기위한 기준 수량",
		"base_value", "기본 수량",
	}
}

//Collector :
func (my *InfoDeposit) Collector() mongo.Bson {
	my.Key = "info_deposit"
	return mongo.Bson{"key": my.Key}
}

//Get :
func (InfoDeposit) Get(db mongo.DATABASE) InfoDeposit {
	data := InfoDeposit{}
	db.C(inf.COLInfoDeposit).Find(data.Collector()).One(&data)
	return data
}
func (my InfoDeposit) UpdateDB(db mongo.DATABASE) {
	db.C(inf.COLInfoDeposit).Update(my.Collector(), my)
}

//IsAllow :
func (my InfoDeposit) IsAllow(symbol string, price string) bool {
	if jmath.CMP(price, 0) <= 0 {
		return false
	}

	if v, do := my.Coin[symbol]; do {
		return jmath.CMP(price, v) >= 0
	}

	return jmath.CMP(price, my.BaseValue) >= 0
}

//IndexingDB :
func (InfoDeposit) IndexingDB() {
	inf.DBCollection(inf.COLInfoDeposit, func(c mongo.Collection) {
		if cnt, _ := c.Find(nil).Count(); cnt == 0 {
			data := InfoDeposit{
				Coin: CoinData{
					ETH: "1000000000000000.0",
				},
				BaseValue: "100000000000000000000",
			}
			data.Collector()
			c.Insert(data)
		}
	})
}
