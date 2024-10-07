package model

import (
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/mms"
	"txscheduler/txm/inf"
)

type LogExMaster struct {
	Hash     string `bson:"hash" json:"hash"`
	From     string `bson:"from" json:"from"`
	Symbol   string `bson:"symbol" json:"symbol"`
	Contract string `bson:"contract" json:"contract"`
	Decimal  string `bson:"decimal" json:"decimal"`
	Price    string `bson:"price" json:"price"`

	Timestamp mms.MMS `bson:"timestamp" json:"timestamp"`
	YMD       int     `bson:"ymd" json:"ymd"`

	SendAt  mms.MMS `bson:"send_at" json:"send_at"`
	SendYMD int     `bson:"send_ymd" json:"send_ymd"`
	IsSend  bool    `bson:"is_send" json:"is_send"`
}

type LogExMasterList []LogExMaster

/*
	"hash":      my.Hash,
	"from":      my.From,
	"symbol":    my.Symbol,
	"contract":  my.Contract,
	"price":     my.Price,
	"timestamp": my.Timestamp,
*/
func (my LogExMaster) AckJson() chttp.JsonType {
	data := chttp.JsonType{
		"hash":      my.Hash,
		"from":      my.From,
		"symbol":    my.Symbol,
		"contract":  my.Contract,
		"price":     my.Price,
		"timestamp": my.Timestamp,
	}
	return data
}

func (my LogExMaster) String() string     { return dbg.ToJSONString(my) }
func (my LogExMasterList) String() string { return dbg.ToJSONString(my) }

func (my LogExMaster) Valid() bool          { return my.Hash != "" }
func (my LogExMaster) Selector() mongo.Bson { return mongo.Bson{"hash": my.Hash} }

func (my LogExMaster) InsertDB(db mongo.DATABASE) {
	db.C(inf.COLLogExMaster).Insert(my)
}

func (LogExMaster) GetList(db mongo.DATABASE) LogExMasterList {
	list := LogExMasterList{}
	selector := mongo.Bson{
		"is_send": false,
	}
	db.C(inf.COLLogExMaster).
		Find(selector).
		Sort("timestamp").
		All(&list)

	return list
}

func (my LogExMaster) SendOK(db mongo.DATABASE, nowAt mms.MMS) {
	upQuery := mongo.Bson{"$set": mongo.Bson{
		"send_at":  nowAt,
		"send_ymd": nowAt.YMD(),
		"is_send":  true,
	}}
	db.C(inf.COLLogExMaster).Update(my.Selector(), upQuery)
}

func (LogExMaster) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.COLLogExMaster)
		c.EnsureIndex(mongo.SingleIndex("hash", "1", true))

		c.EnsureIndex(mongo.SingleIndex("timestamp", "1", false))
		c.EnsureIndex(mongo.SingleIndex("ymd", "1", false))
		c.EnsureIndex(mongo.SingleIndex("from", "1", false))
		c.EnsureIndex(mongo.SingleIndex("is_send", "1", false))
	})
}
