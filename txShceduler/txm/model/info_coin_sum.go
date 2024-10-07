package model

import (
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/txm/inf"
)

type ConSum struct {
	Key  string   `bson:"key" json:"key"`
	Coin CoinData `bson:"coin" json:"coin"`
}

func (my *ConSum) Collector() mongo.Bson {
	my.Key = "info_coin_sum"
	return mongo.Bson{"key": my.Key}
}

func (ConSum) Get(db mongo.DATABASE) ConSum {
	data := ConSum{}
	db.C(inf.COLConSum).Find(data.Collector()).One(&data)
	return data
}
func (my ConSum) UpdateDB(db mongo.DATABASE) {
	db.C(inf.COLConSum).Update(my.Collector(), my)
}

func CoinSumAdd(db mongo.DATABASE, symbol, price string) {
	coin := ConSum{}.Get(db)
	coin.Coin.ADD(symbol, price)
	coin.UpdateDB(db)
}
func CoinSumSub(db mongo.DATABASE, symbol, price string) {
	coin := ConSum{}.Get(db)
	coin.Coin.SUB(symbol, price)
	coin.UpdateDB(db)
}

func CoinSumAction(db mongo.DATABASE, f func(coin *ConSum)) {
	coin := ConSum{}.Get(db)
	f(&coin)
	coin.UpdateDB(db)
}

func (ConSum) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		db.C(inf.COLConSum).EnsureIndex(mongo.SingleIndex("key", "1", true))

		if cnt, _ := db.C(inf.COLConSum).Count(); cnt == 0 {
			data := ConSum{
				Coin: NewCoinData(),
			}
			data.Collector()

			if cnt, _ := db.C(inf.COLMember).Count(); cnt > 0 {
				member := Member{}
				iter := db.C(inf.COLMember).Find(nil).Iter()
				for iter.Next(&member) {
					data.Coin.AddCoinData(member.Coin)
				} //for
			}
			db.C(inf.COLConSum).Insert(data)
		}
	})
}
