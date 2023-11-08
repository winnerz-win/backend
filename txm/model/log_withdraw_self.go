package model

import (
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/mms"
	"txscheduler/txm/inf"
)

func (my TxETHWithdraw) MakeLogWithdrawSELF(member Member, state TxState) LogWithdrawSELF {
	log := LogWithdrawSELF{
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

type LogWithdrawSELF struct {
	User `bson:",inline" json:",inline"`

	ReceiptCode string `bson:"receipt_code" json:"receipt_code"`

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

type LogWithdrawSELFList []LogWithdrawSELF

func (my LogWithdrawSELF) Valid() bool          { return my.Hash != "" }
func (my LogWithdrawSELF) Selector() mongo.Bson { return mongo.Bson{"hash": my.Hash} }

func (my LogWithdrawSELF) String() string     { return dbg.ToJSONString(my) }
func (my LogWithdrawSELFList) String() string { return dbg.ToJSONString(my) }

func (LogWithdrawSELF) GetList(db mongo.DATABASE) LogWithdrawSELFList {
	list := LogWithdrawSELFList{}

	db.C(inf.COLLogWithdrawSELF).
		Find(mongo.Bson{"is_send": false}).
		Sort("timestamp").
		All(&list)

	return list
}

func (my LogWithdrawSELF) AckPairs() []interface{} {
	pairs := []interface{}{
		"receipt_code", my.ReceiptCode,

		"name", my.Name,
		"uid", my.UID,
		"from_address", my.Address,
		"to_address", my.ToAddress,
		"hash", my.Hash,
		"symbol", my.Symbol,
		"price", my.ToPrice,
		"gas", my.Gas,
		"withdraw_result", my.State, //0 , 1 , 104, 200
		"timestamp", my.Timestamp,
		//"is_send", my.IsSend,
		//"send_at", my.SendAt,
	}
	return pairs
}
func (my LogWithdrawSELF) AckJson() chttp.JsonType {
	data := chttp.JsonType{}
	pairs := my.AckPairs()
	for i := 0; i < len(pairs); i += 2 {
		key := pairs[i].(string)
		val := pairs[i+1]
		data[key] = val
	} //for
	return data
}

func (my LogWithdrawSELF) InsertDB(db mongo.DATABASE, nowAt mms.MMS) {
	my.Timestamp = nowAt
	my.YMD = nowAt.YMD()
	my.IsSend = false
	db.C(inf.COLLogWithdrawSELF).Insert(my)
}

func (my LogWithdrawSELF) SendOK(db mongo.DATABASE, nowAt mms.MMS) {
	upQuery := mongo.Bson{"$set": mongo.Bson{
		"send_at":  nowAt,
		"send_ymd": nowAt.YMD(),
		"is_send":  true,
	}}
	db.C(inf.COLLogWithdrawSELF).Update(my.Selector(), upQuery)
}

func (LogWithdrawSELF) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.COLLogWithdrawSELF)
		c.EnsureIndex(mongo.SingleIndex("hash", "1", true))
		c.EnsureIndex(mongo.SingleIndex("receipt_code", "1", false))

		c.EnsureIndex(mongo.SingleIndex("uid", "1", false))
		c.EnsureIndex(mongo.SingleIndex("address", "1", false))

		c.EnsureIndex(mongo.SingleIndex("is_send", "1", false))
	})
}

func (LogWithdrawSELF) GetData(db mongo.DATABASE, receiptCode string) LogWithdrawSELF {
	item := LogWithdrawSELF{}
	db.C(inf.COLLogWithdrawSELF).Find(mongo.Bson{"receipt_code": receiptCode}).One(&item)
	return item
}
