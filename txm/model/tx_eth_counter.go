package model

import (
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
	"txscheduler/txm/inf"
)

//TxETHCounter :
type TxETHCounter struct {
	Mainnet     bool   `bson:"mainnet"  json:"mainnet"`
	Number      string `bson:"number" json:"number"`
	ChainNumber string `bson:"chain_number" json:"chain_number"`
}

//NewTxETHCounter :
func NewTxETHCounter(mainnet bool, startNumber string) *TxETHCounter {
	return &TxETHCounter{
		Mainnet: mainnet,
		Number:  startNumber,
	}
}

//GetNumber :
func (my TxETHCounter) GetNumber() string { return my.Number }

//Selector :
func (my TxETHCounter) Selector() mongo.Bson { return mongo.Bson{"mainnet": my.Mainnet} }

//IndexingDB :
func (my TxETHCounter) IndexingDB() {
	inf.DB().Run(inf.DBName, inf.TXETHCount, func(c mongo.Collection) {
		c.EnsureIndex(mongo.SingleIndex("mainnet", "1", true))
	})
}

//LoadFromDB :
func (my *TxETHCounter) LoadFromDB(notOverWrite ...bool) {
	startNumber := my.Number
	inf.DB().Run(inf.DBName, inf.TXETHCount, func(c mongo.Collection) {
		if c.Find(my.Selector()).One(my) != nil {
			my.Number = ZERO
		}
	})
	if dbg.IsTrue2(notOverWrite...) == false {
		if jmath.CMP(startNumber, my.Number) > 0 {
			my.Number = startNumber
		}
	}
}

//Update :
func (my *TxETHCounter) Update(lastNumber string) {
	if jmath.CMP(my.Number, lastNumber) >= 0 {
		return
	}
	my.Number = lastNumber
	inf.DB().Run(inf.DBName, inf.TXETHCount, func(c mongo.Collection) {
		c.Upsert(my.Selector(), my)
	})
}

//Inc :
func (my *TxETHCounter) Inc(chainNumber string) {
	my.Number = jmath.ADD(my.Number, 1)
	my.ChainNumber = chainNumber
	inf.DB().Run(inf.DBName, inf.TXETHCount, func(c mongo.Collection) {
		c.Upsert(my.Selector(), my)
	})
}

//ResetDB :
func (my *TxETHCounter) ResetDB() {
	my.Number = "0"
	inf.DB().Run(inf.DBName, inf.TXETHCount, func(c mongo.Collection) {
		c.Upsert(my.Selector(), my)
	})
}
