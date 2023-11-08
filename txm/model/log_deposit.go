package model

import (
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/mms"
	"txscheduler/txm/inf"
)

type LogDeposit struct {
	User `bson:",inline" json:",inline"`

	Hash     string `bson:"hash" json:"hash"`
	Symbol   string `bson:"symbol" json:"symbol"`
	Contract string `bson:"contract" json:"contract"`
	Decimal  string `bson:"decimal" json:"decimal"`
	Price    string `bson:"price" json:"price"`
	From     string `bson:"from" json:"from"`

	Timestamp mms.MMS `bson:"timestamp" json:"timestamp"`
	YMD       int     `bson:"ymd" json:"ymd"`

	DepositResult bool `bson:"deposit_result" json:"deposit_result"`

	SendAt  mms.MMS `bson:"send_at" json:"send_at"`
	SendYMD int     `bson:"send_ymd" json:"send_ymd"`
	IsSend  bool    `bson:"is_send" json:"is_send"`
}

type LogDepositList []LogDeposit

func (my LogDeposit) String() string     { return dbg.ToJSONString(my) }
func (my LogDepositList) String() string { return dbg.ToJSONString(my) }

func (LogDeposit) GetList(db mongo.DATABASE) LogDepositList {
	list := LogDepositList{}

	db.C(inf.COLLogDeposit).
		Find(mongo.Bson{"is_send": false}).
		Sort("timestamp").
		All(&list)

	return list
}

func (my LogDeposit) InsertDB(db mongo.DATABASE) {
	db.C(inf.COLLogDeposit).Insert(my)
}

func (my LogDeposit) Selector() mongo.Bson { return mongo.Bson{"hash": my.Hash} }

func (my LogDeposit) SendOK(db mongo.DATABASE, nowAt mms.MMS) {
	upQuery := mongo.Bson{"$set": mongo.Bson{
		"send_at":  nowAt,
		"send_ymd": nowAt.YMD(),
		"is_send":  true,
	}}
	db.C(inf.COLLogDeposit).Update(my.Selector(), upQuery)
}

func (LogDeposit) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.COLLogDeposit)
		c.EnsureIndex(mongo.SingleIndex("hash", "1", true))

		c.EnsureIndex(mongo.SingleIndex("uid", "1", false))
		c.EnsureIndex(mongo.SingleIndex("address", "1", false))

		c.EnsureIndex(mongo.SingleIndex("is_send", "1", false))
	})
}
