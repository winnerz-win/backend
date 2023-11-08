package model

import (
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/jmath"
	"txscheduler/brix/tools/mms"
	"txscheduler/txm/inf"
)

//DepositPath : 유입경로.
type DepositPath string

const (
	DepositPathIncome = DepositPath("income")
	DepositPathAdmin  = DepositPath("admin`")
)

//TxETHDeposit :
type TxETHDeposit struct {
	Key      string      `bson:"key" json:"key"`
	UID      int64       `bson:"uid" json:"uid"`
	Path     DepositPath `bson:"path" json:"path"`
	Symbol   string      `bson:"symbol" json:"symbol"`
	Contract string      `bson:"contract" json:"contract"`
	Decimal  string      `bson:"decimal" json:"decimal"`
	Address  string      `bson:"address" json:"address"`
	Price    string      `bson:"price" json:"price"`
	State    TxState     `bson:"state" json:"state"`

	Hash       string `bson:"hash" json:"hash"`
	Gas        string `bson:"gas" json:"gas"`
	IsGasFixed bool   `bson:"is_gas_fixed" json:"is_gas_fixed"`
	GasLimit   string `bson:"gas_limit" json:"gas_limit"`
	GasPrice   string `bson:"gas_price" json:"gas_price"`

	CreateAt  mms.MMS `bson:"create_at" json:"create_at"`
	CreateYMD int     `bson:"create_ymd" json:"create_ymd"`
	Timestamp mms.MMS `bson:"timestamp" json:"timestamp"`
	YMD       int     `bson:"ymd" json:"ymd"`

	IsForce bool `bson:"is_force" json:"is_force"`
}

//Selector :
func (my TxETHDeposit) Selector() mongo.Bson { return mongo.Bson{"key": my.Key} }

//InsertDB :
func (my TxETHDeposit) InsertDB(db mongo.DATABASE, nowAt mms.MMS) {
	selector := mongo.Bson{
		"uid":      my.UID,
		"contract": my.Contract,
		"$or": []mongo.Bson{
			{"state": TxStateNone},
			{"state": TxStatePending},
		},
	}
	if cnt, _ := db.C(inf.TXETHDepositLog).Find(selector).Count(); cnt == 0 {
		my.State = TxStateNone
		my.CreateAt = nowAt
		my.CreateYMD = nowAt.YMD()
		my.Timestamp = nowAt
		my.YMD = nowAt.YMD()
		db.C(inf.TXETHDepositLog).Insert(my)
	}
}

//CalcGas :
func (my *TxETHDeposit) CalcGas() {
	my.Gas = jmath.MUL(my.GasLimit, my.GasPrice)
}

//UpdateDB :
func (my TxETHDeposit) UpdateDB(db mongo.DATABASE, nowAt mms.MMS) {
	my.Timestamp = nowAt
	my.YMD = nowAt.YMD()
	db.C(inf.TXETHDepositLog).Update(my.Selector(), my)
}

//RemoveDB :
func (my TxETHDeposit) RemoveDB(db mongo.DATABASE) {
	db.C(inf.TXETHDepositLog).Remove(my.Selector())
}

//IndexingDB :
func (TxETHDeposit) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		db.C(inf.TXETHDepositLog).EnsureIndex(mongo.SingleIndex("key", "1", true))
	})
}
