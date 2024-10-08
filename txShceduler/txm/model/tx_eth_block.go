package model

import (
	"jtools/cloud/ebcm"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"

	"txscheduler/txm/inf"
)

// TxState :
type TxState int

const (
	TxStateNone        = TxState(0)   //기본값 - 출금 대기
	TxStatePending     = TxState(1)   //출금 실행
	TxStatePendingSELF = TxState(22)  //출금 실행 (개인지갑 출금)
	TxStateGas         = TxState(33)  //가스비 대기중
	TxStateFail        = TxState(104) //결과-실패
	TxStateSuccess     = TxState(200) //결과-성공
)

// TxETHBlock :
type TxETHBlock struct {
	ebcm.TransactionBlock `bson:",inline" json:",inline"`
	Price                 string  `bson:"price" json:"price"`
	UID                   int64   `bson:"uid" json:"uid"`
	Order                 int64   `bson:"order" json:"order"` //blocknumber
	TxState               TxState `bson:"txstate" json:"txstate"`
}

func (my TxETHBlock) String() string { return dbg.ToJSONString(my) }

func TxErrorState(isError bool) TxState {
	if isError {
		return TxStateFail
	}
	return TxStateSuccess
}

// IndexingDB : .
func (my TxETHBlock) IndexingDB() {
	inf.DB().Run(inf.DBName, inf.TXETHBlock, func(c mongo.Collection) {
		c.EnsureIndex(mongo.SingleIndex("hash", "1", true))
		c.EnsureIndex(mongo.SingleIndex("to", "1", false))
		c.EnsureIndex(mongo.SingleIndex("is_contract", "1", false))
		c.EnsureIndex(mongo.SingleIndex("contractAddress", "1", false))

		c.EnsureIndex(mongo.SingleIndex("uid", "1", false))
		c.EnsureIndex(mongo.SingleIndex("order", "1", false))
		c.EnsureIndex(mongo.SingleIndex("txstate", "1", false))
	})
}

// IsInert :
func (my *TxETHBlock) IsInert(db mongo.DATABASE) error {
	price := ebcm.WeiToToken(my.Amount, my.Decimals)
	my.Price = price
	return db.C(inf.TXETHBlock).Insert(my)
}

// Load :
func (my *TxETHBlock) Load(hash string) error {
	var err error
	inf.DB().Run(inf.DBName, inf.TXETHBlock, func(c mongo.Collection) {
		err = c.Find(mongo.Bson{"hash": hash}).One(my)
	})
	return err
}

// TestLoad :
func (my *TxETHBlock) TestLoad(hash string) error {
	var err error
	inf.Test4004DB().Run(inf.DBName, inf.TXETHBlock, func(c mongo.Collection) {
		err = c.Find(mongo.Bson{"hash": hash}).One(my)
	})
	return err
}
