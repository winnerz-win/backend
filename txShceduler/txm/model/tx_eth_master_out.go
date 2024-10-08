package model

import (
	"jtools/cloud/ebcm"
	"jtools/mms"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/database/mongo/tools/jmath"
	"txscheduler/brix/tools/dbg"
	"txscheduler/txm/inf"
)

type TxETHMasterOut struct {
	ReceiptCode string `bson:"receipt_code" json:"receipt_code"`
	Symbol      string `bson:"symbol" json:"symbol"`
	Decimal     string `bson:"decimal" json:"decimal"`
	ToAddress   string `bson:"to_address" json:"to_address"`
	ToPrice     string `bson:"to_price" json:"to_price"`

	Gas      string `bson:"gas" json:"gas"`
	GasLimit string `bson:"gas_limit,omitempty" json:"gas_limit,omitempty"`
	GasPrice string `bson:"gas_price,omitempty" json:"gas_price,omitempty"`

	Hash        string  `bson:"hash" json:"hash"`
	Nonce       string  `bson:"nonce" json:"nonce"`
	State       TxState `bson:"state" json:"state"`
	FailMessage string  `bson:"fail_message,omitempty" json:"fail_message,omitempty"`
	Timestamp   mms.MMS `bson:"timestamp" json:"timestamp"`
	YMD         int     `bson:"ymd" json:"ymd"`

	//
	//
	TryData TxTryData               `bson:"try_data" json:"try_data"`
	TrySTX  ebcm.WrappedTransaction `bson:"-" json:"-"`
	// RetryCnt    int     `bson:"retry_cnt" json:"retry_cnt"`
	// Mv          string  `bson:"mv" json:"mv"`

	// CancelTry  bool    `bson:"cancel_try" json:"cancel_try"`
	// CancelHash string  `bson:"cancel_hash,omitempty" json:"cancel_hash,omitempty"`
	// CancelTime mms.MMS `bson:"cancel_time,omitempty" json:"cancel_time,omitempty"`
	// CancelYMD  int     `bson:"cancel_ymd,omitempty" json:"cancel_ymd,omitempty"`

	SendAt  mms.MMS `bson:"send_at" json:"send_at"`
	SendYMD int     `bson:"send_ymd" json:"send_ymd"`
	IsSend  bool    `bson:"is_send" json:"is_send"`
}

type TxETHMasterOutList []TxETHMasterOut

func (my *TxETHMasterOut) SetNonce(n any)  { my.Nonce = jmath.VALUE(n) }
func (my TxETHMasterOut) GetNonce() uint64 { return jmath.Uint64(my.Nonce) }

// Wei : ebcm.TokenToWei(my.ToPrice, my.Decimal)
func (my TxETHMasterOut) Wei() string { return ebcm.TokenToWei(my.ToPrice, my.Decimal) }

func (my TxETHMasterOut) Selector() mongo.Bson { return mongo.Bson{"receipt_code": my.ReceiptCode} }
func (my TxETHMasterOut) Valid() bool          { return my.ReceiptCode != "" }

func (my TxETHMasterOut) String() string     { return dbg.ToJSONString(my) }
func (my TxETHMasterOutList) String() string { return dbg.ToJSONString(my) }

func (my TxETHMasterOut) IndexingDB() {
	DB(func(db mongo.DATABASE) {

		runIndex := func(c mongo.Collection) {
			c.EnsureIndex(mongo.SingleIndex("receipt_code", "1", true))
			c.EnsureIndex(mongo.SingleIndex("hash", "1", false))

			c.EnsureIndex(mongo.SingleIndex("try_data.nonce", "1", false))

			c.EnsureIndex(mongo.SingleIndex("state", "1", false))
			c.EnsureIndex(mongo.SingleIndex("timestamp", "1", false))

			c.EnsureIndex(mongo.SingleIndex("is_send", "1", false))
		}

		runIndex(db.C(inf.TXETHMasterOut))
		runIndex(db.C(inf.TXETHMasterOutTry))
	})
}

func (TxETHMasterOut) GetReceipt(db mongo.DATABASE, receipt string) TxETHMasterOut {
	item := TxETHMasterOut{}

	selector := mongo.Bson{"receipt_code": receipt}
	if db.C(inf.TXETHMasterOutTry).Find(selector).One(&item) == nil {
		return item
	} else {
		db.C(inf.TXETHMasterOut).Find(selector).One(&item)
	}

	return item
}

//============================================\

func (my TxETHMasterOut) UpdateTry(db mongo.DATABASE) error {
	return db.C(inf.TXETHMasterOutTry).Update(my.Selector(), my)
}

func (my TxETHMasterOut) InsertTry(db mongo.DATABASE) error {
	return db.C(inf.TXETHMasterOutTry).Insert(my)
}

func (my TxETHMasterOut) RemoveTry(db mongo.DATABASE) error {
	return db.C(inf.TXETHMasterOutTry).Remove(my.Selector())
}

//============================================\

func (my TxETHMasterOut) UpdateLog(db mongo.DATABASE) error {
	return db.C(inf.TXETHMasterOut).Update(my.Selector(), my)
}

func (my TxETHMasterOut) InsertLog(db mongo.DATABASE) error {
	return db.C(inf.TXETHMasterOut).Insert(my)
}

func (my TxETHMasterOut) RemoveLog(db mongo.DATABASE) error {
	return db.C(inf.TXETHMasterOut).Remove(my.Selector())
}

func (my TxETHMasterOut) GetCallbackList(db mongo.DATABASE) TxETHMasterOutList {
	list := TxETHMasterOutList{}

	selector := mongo.Bson{
		"is_send": false,
	}
	db.C(inf.TXETHMasterOut).
		Find(selector).
		Sort("timestamp").
		All(&list)
	return list
}

func (my TxETHMasterOut) SendOK(db mongo.DATABASE, nowAt mms.MMS) {
	upQuery := mongo.Bson{"$set": mongo.Bson{
		"send_at":  nowAt,
		"send_ymd": nowAt.YMD(),
		"is_send":  true,
	}}
	db.C(inf.TXETHMasterOut).Update(my.Selector(), upQuery)
}
