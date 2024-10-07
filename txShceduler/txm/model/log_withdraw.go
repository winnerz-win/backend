package model

import (
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/mms"
	"txscheduler/txm/inf"
)

type LogWithdraw struct {
	User        `bson:",inline" json:",inline"`
	ReceiptCode string `bson:"receipt_code" json:"receipt_code"` //UID

	Symbol    string `bson:"symbol" json:"symbol"`
	Decimal   string `bson:"decimal" json:"decimal"`
	ToAddress string `bson:"to_address" json:"to_address"`
	ToPrice   string `bson:"to_price" json:"to_price"`

	Hash      string  `bson:"hash" json:"hash"`
	Gas       string  `bson:"gas" json:"gas"`
	State     TxState `bson:"state" json:"state"`
	Timestamp mms.MMS `bson:"timestamp" json:"timestamp"`
	YMD       int     `bson:"ymd" json:"ymd"`

	SendAt  mms.MMS `bson:"send_at" json:"send_at"`
	SendYMD int     `bson:"send_ymd" json:"send_ymd"`
	IsSend  bool    `bson:"is_send" json:"is_send"`
}

func (my LogWithdraw) AckPairs() []interface{} {
	pairs := []interface{}{
		"receipt_code", my.ReceiptCode,
		"name", my.Name,
		"uid", my.UID,
		"from_address", my.Address,
		"to_address", my.ToAddress,
		"hash", my.Hash,
		"gas", my.Gas,
		"symbol", my.Symbol,
		"price", my.ToPrice,
		"withdraw_result", my.State, //0 , 1 , 104, 200
		"timestamp", my.Timestamp,
		//"is_send", my.IsSend,
		//"send_at", my.SendAt,
	}
	return pairs
}
func (my LogWithdraw) AckJson() chttp.JsonType {
	data := chttp.JsonType{}
	pairs := my.AckPairs()
	for i := 0; i < len(pairs); i += 2 {
		key := pairs[i].(string)
		val := pairs[i+1]
		data[key] = val
	} //for
	return data
}

type LogWithdrawList []LogWithdraw

func (my LogWithdraw) String() string     { return dbg.ToJSONString(my) }
func (my LogWithdrawList) String() string { return dbg.ToJSONString(my) }

func (my LogWithdraw) Valid() bool { return my.ReceiptCode != "" }

func (LogWithdraw) GetList(db mongo.DATABASE) LogWithdrawList {
	list := LogWithdrawList{}

	db.C(inf.COLLogWithdraw).
		Find(mongo.Bson{"is_send": false}).
		Sort("timestamp").
		All(&list)

	return list
}

func (my TxETHWithdraw) MakeLogWithdraw(member Member, state TxState) LogWithdraw {
	log := LogWithdraw{
		ReceiptCode: my.ReceiptCode,
		User:        member.User,

		ToAddress: my.ToAddress,
		Symbol:    my.Symbol,
		ToPrice:   my.ToPrice,
		Decimal:   my.Decimal,

		Hash:  my.Hash,
		Gas:   my.Gas,
		State: state,
	}

	return log
}

func (my LogWithdraw) Selector() mongo.Bson { return mongo.Bson{"receipt_code": my.ReceiptCode} }

func (my LogWithdraw) InsertDB(db mongo.DATABASE, nowAt mms.MMS) {
	my.Timestamp = nowAt
	my.YMD = nowAt.YMD()
	my.IsSend = false
	db.C(inf.COLLogWithdraw).Insert(my)
}

func (my LogWithdraw) SendOK(db mongo.DATABASE, nowAt mms.MMS) {
	upQuery := mongo.Bson{"$set": mongo.Bson{
		"send_at":  nowAt,
		"send_ymd": nowAt.YMD(),
		"is_send":  true,
	}}
	db.C(inf.COLLogWithdraw).Update(my.Selector(), upQuery)
}

func (LogWithdraw) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.COLLogWithdraw)
		c.EnsureIndex(mongo.SingleIndex("hash", "1", true))
		c.EnsureIndex(mongo.SingleIndex("receipt_code", "1", true))

		c.EnsureIndex(mongo.SingleIndex("uid", "1", false))
		c.EnsureIndex(mongo.SingleIndex("address", "1", false))

		c.EnsureIndex(mongo.SingleIndex("is_send", "1", false))
	})
}

func (LogWithdraw) GetData(db mongo.DATABASE, receiptCode string) LogWithdraw {
	item := LogWithdraw{}
	db.C(inf.COLLogWithdraw).Find(mongo.Bson{"receipt_code": receiptCode}).One(&item)
	return item
}
