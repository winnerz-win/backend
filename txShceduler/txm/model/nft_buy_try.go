package model

import (
	"jtools/mms"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/txm/inf"
)

type NftBuyTry struct {
	NftRevData `bson:",inline" json:",inline"`

	Hash        string      `bson:"hash" json:"hash"`
	Status      int         `bson:"status" json:"status"`
	FailCode    FailNftCode `bson:"fail_code" json:"fail_code"`
	FailMessage string      `bson:"fail_message" json:"fail_message"`

	// GasLimit  uint64 `bson:"gasLimit" json:"gasLimit"`
	// GasPrice  string `bson:"gasPrice" json:"gasPrice"`
	// GasFeeETH string `bson:"gasFeeETH" json:"gasFeeETH"`

	CreateTime mms.MMS `bson:"create_time" json:"create_time"`
	CreateYMD  int     `bson:"create_ymd" json:"create_ymd"`
}
type NftBuyTryList []NftBuyTry

func (my NftBuyTry) String() string     { return dbg.ToJSONString(my) }
func (my NftBuyTryList) String() string { return dbg.ToJSONString(my) }

func (my *NftBuyTry) SetFail(code FailNftCode, exMsg ...string) {
	my.Status = 104
	my.FailCode = code
	my.FailMessage = NFTFailCodeMessage(code, exMsg...)
}

func (NftBuyTry) GetByReceiptCode(db mongo.DATABASE, receiptCode string) NftBuyTry {
	item := NftBuyTry{}
	DB(func(db mongo.DATABASE) {
		selector := mongo.Bson{"receipt_code": receiptCode}
		db.C(inf.NFTBuyTry).Find(selector).One(&item)
	})
	return item
}
func (NftBuyTry) GetTokenID(db mongo.DATABASE, token_id string) NftBuyTry {
	item := NftBuyTry{}
	DB(func(db mongo.DATABASE) {
		selector := mongo.Bson{"token_id": token_id}
		db.C(inf.NFTBuyTry).Find(selector).One(&item)
	})
	return item
}

func (my NftBuyTry) InsertTryDB(db mongo.DATABASE) error {

	my.CreateTime = mms.Now()
	my.CreateYMD = my.CreateTime.YMD()

	return db.C(inf.NFTBuyTry).Insert(my)

}
func (my NftBuyTry) RemoveTryDB(db mongo.DATABASE) error {
	return db.C(inf.NFTBuyTry).Remove(my.Selector())
}

func (my NftBuyTry) UpdateTryDB(db mongo.DATABASE) {
	db.C(inf.NFTBuyTry).Update(my.Selector(), my)
}

func (my NftBuyTry) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.NFTBuyTry)
		NftRevData{}._indexingDB(c, true)

		c.EnsureIndex(mongo.SingleIndex("create_time", "1", false))
		c.EnsureIndex(mongo.SingleIndex("create_ymd", "1", false))

		c.EnsureIndex(mongo.SingleIndex("hash", "1", false))

		c.EnsureIndex(mongo.SingleIndex("status", "1", false))

	})
}
