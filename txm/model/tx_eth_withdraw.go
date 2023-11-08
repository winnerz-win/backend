package model

import (
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
	"txscheduler/brix/tools/mms"
	"txscheduler/txm/inf"
)

// TxETHWithdraw :
type TxETHWithdraw struct {
	ReceiptCode string `bson:"receipt_code" json:"receipt_code"`
	UID         int64  `bson:"uid" json:"uid"`

	Symbol    string `bson:"symbol" json:"symbol"`
	Decimal   string `bson:"decimal" json:"decimal"`
	ToAddress string `bson:"to_address" json:"to_address"`
	ToPrice   string `bson:"to_price" json:"to_price"`

	Gas      string `bson:"gas" json:"gas"`
	GasLimit string `bson:"gas_limit,omitempty" json:"gas_limit,omitempty"`
	GasPrice string `bson:"gas_price,omitempty" json:"gas_price,omitempty"`

	Hash      string  `bson:"hash" json:"hash"`
	State     TxState `bson:"state" json:"state"`
	Timestamp mms.MMS `bson:"timestamp" json:"timestamp"`
	YMD       int     `bson:"ymd" json:"ymd"`

	CancelTry  bool    `bson:"cancel_try" json:"cancel_try"`
	CancelHash string  `bson:"cancel_hash,omitempty" json:"cancel_hash,omitempty"`
	CancelTime mms.MMS `bson:"cancel_time,omitempty" json:"cancel_time,omitempty"`
	CancelYMD  int     `bson:"cancel_ymd,omitempty" json:"cancel_ymd,omitempty"`
}

func (my TxETHWithdraw) Valid() bool { return my.ReceiptCode != "" }

// Wei : ecsx.TokenToWei(my.ToPrice, my.Decimal)
func (my TxETHWithdraw) Wei() string { return ecsx.TokenToWei(my.ToPrice, my.Decimal) }

// Selector : pair
func (my TxETHWithdraw) Selector() mongo.Bson { return mongo.Bson{"receipt_code": my.ReceiptCode} }

// TxETHWithdrawList :
type TxETHWithdrawList []TxETHWithdraw

// TotalPrice :
func (my TxETHWithdrawList) TotalPrice(symbol string) string {
	sum := ZERO
	for _, v := range my {
		if v.Symbol == symbol {
			sum = jmath.ADD(sum, v.ToPrice)
		}
	}
	return sum
}

func (my TxETHWithdraw) String() string     { return dbg.ToJSONString(my) }
func (my TxETHWithdrawList) String() string { return dbg.ToJSONString(my) }

// IndexingDB :
func (my TxETHWithdraw) IndexingDB() {
	DB(func(db mongo.DATABASE) {

		rc := db.C(inf.TXETHWithdraw)
		rc.EnsureIndex(mongo.SingleIndex("receipt_code", "1", true))

		rc.EnsureIndex(mongo.SingleIndex("hash", "1", false))
		rc.EnsureIndex(mongo.SingleIndex("uid", "1", false))
		rc.EnsureIndex(mongo.SingleIndex("state", "1", false))
		rc.EnsureIndex(mongo.SingleIndex("timestamp", "1", false))

		// lc := db.C(inf.TXWithdrawLog)
		// lc.EnsureIndex(mongo.SingleIndex("hash", "1", true))
		// lc.EnsureIndex(mongo.SingleIndex("uid", "1", false))
		// rc.EnsureIndex(mongo.SingleIndex("state", "1", false))
		// lc.EnsureIndex(mongo.SingleIndex("timestamp", "1", false))
	})
}

// InsertDB :
func (my TxETHWithdraw) InsertDB(db mongo.DATABASE) string {
	for {
		my.ReceiptCode = GetReceiptCode()
		if db.C(inf.TXETHWithdraw).Insert(my) == nil {
			break
		}
	} //for
	return my.ReceiptCode
}

// InsertDB :
func (my TxETHWithdraw) InsertSELF_DB(db mongo.DATABASE) string {
	for {
		my.ReceiptCode = GetReceiptCodeSELF()
		if db.C(inf.TXETHWithdraw).Insert(my) == nil {
			break
		}
	} //for
	return my.ReceiptCode
}

// UpdateDB :
func (my TxETHWithdraw) UpdateDB(db mongo.DATABASE) {
	db.C(inf.TXETHWithdraw).Update(my.Selector(), my)
}

func (my TxETHWithdraw) GetData(db mongo.DATABASE, receiptCode string) TxETHWithdraw {
	item := TxETHWithdraw{}
	db.C(inf.TXETHWithdraw).Find(mongo.Bson{"receipt_code": receiptCode}).One(&item)
	return item
}
