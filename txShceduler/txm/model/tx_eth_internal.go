package model

import (
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/jmath"
	"txscheduler/txm/inf"
)

/*
{
  "blockNumber": "11459686",
  "timeStamp": 1608062106000,
  "hash": "0x0c52ade5166c5003411e504f7de6622bc58f19b8afde9fa4df30e9eec37d34c5",
  "from": "0x6c0b51971650d28821ce30b15b02b9826a20b129",
  "to": "0xf57029c25af2c31e763ef643b8a999ae0ec9ee8a",
  "value": "3105640000000000000",
  "is_contract": false,
  "input": "",
  "type": "call",
  "gas": "202229",
  "gasUsed": "0",
  "traceId": "0",
  "isError": false,
  "errCode": "",
  "time_kst": "2020-12-16 04:55:06 +0000 KST"
}
*/

// TxETHInternalCnt :
type TxETHInternalCnt struct {
	Name     string `bson:"name" json:"name"`
	Contract string `bson:"contract" json:"contract"`
	Mainnet  bool   `bson:"mainnet" json:"mainnet"`
	Number   string `bson:"number" json:"number"`
}

// NewETHTxInternalCnt :
func NewETHTxInternalCnt(mainnet bool, name, contract, number string) *TxETHInternalCnt {
	return &TxETHInternalCnt{
		Name:     name,
		Contract: contract,
		Mainnet:  mainnet,
		Number:   number,
	}
}

// GetNumber :
func (my TxETHInternalCnt) GetNumber() string { return my.Number }

// Selector :
func (my TxETHInternalCnt) Selector() mongo.Bson {
	return mongo.Bson{
		"contract": my.Contract,
		"mainnet":  my.Mainnet,
	}
}

// IndexingDB :
func (my TxETHInternalCnt) IndexingDB() {
	inf.DBCollection(inf.TXETHInternalCnt, func(c mongo.Collection) {
		mongo.MultiIndex([]interface{}{
			"mainnet", "1",
			"contract", "1",
		}, false)
	})
}

// LoadFromDB :
func (my *TxETHInternalCnt) LoadFromDB() {
	startNumber := my.Number
	inf.DB().Run(inf.DBName, inf.TXETHInternalCnt, func(c mongo.Collection) {
		c.Find(my.Selector()).One(my)
	})
	if jmath.CMP(startNumber, my.Number) > 0 {
		my.Number = startNumber
	}
}

// Update :
func (my *TxETHInternalCnt) Update(lastNumber string) {
	if jmath.CMP(my.Number, lastNumber) >= 0 {
		return
	}
	my.Number = lastNumber
	inf.DB().Run(inf.DBName, inf.TXETHInternalCnt, func(c mongo.Collection) {
		c.Upsert(my.Selector(), my)
	})
}

// ResetDB :
func (my *TxETHInternalCnt) ResetDB() {
	my.Number = "0"
	inf.DB().Run(inf.DBName, inf.TXETHInternalCnt, func(c mongo.Collection) {
		c.Upsert(my.Selector(), my)
	})
}
