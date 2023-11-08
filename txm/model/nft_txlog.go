package model

import (
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/mms"
	"txscheduler/txm/inf"
)

type NftTxLog struct {
	Index int `bson:"index" json:"index"`

	Number    int64   `bson:"number" json:"number"`
	Hash      string  `bson:"hash" json:"hash"`
	TxIndex   uint    `bson:"tx_index" json:"tx_index"`
	LogIndex  uint    `bson:"log_index" json:"log_index"`
	Timestamp mms.MMS `bson:"timestamp" json:"timestamp"`

	Name      string `bson:"name" json:"name"`
	From      string `bson:"from" json:"from"`
	To        string `bson:"to" json:"to"`
	TokenID   string `bson:"token_id" json:"token_id"` //KEY
	TokenType string `bson:"token_type" json:"token_type"`
}

type NftTxLogList []NftTxLog

func (my *NftTxLog) Selector() mongo.Bson {
	return mongo.Bson{"index": my.Index}
}

func (my NftTxLog) InsertDB(db mongo.DATABASE) error {
	c := db.C(inf.NFTTxLog)
	cnt, _ := c.Count()
	my.Index = cnt + 1
	return c.Insert(my)
}

func (NftTxLog) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.NFTTxLog)

		c.EnsureIndex(mongo.SingleIndex("index", "1", true))

		c.EnsureIndex(mongo.SingleIndex("number", "1", false))
		c.EnsureIndex(mongo.SingleIndex("hash", "1", false))
		c.EnsureIndex(mongo.SingleIndex("tx_index", "1", false))
		c.EnsureIndex(mongo.SingleIndex("log_index", "1", false))

		c.EnsureIndex(mongo.SingleIndex("timestamp", "1", false))

		c.EnsureIndex(mongo.SingleIndex("name", "1", false))
		c.EnsureIndex(mongo.SingleIndex("from", "1", false))
		c.EnsureIndex(mongo.SingleIndex("to", "1", false))

		c.EnsureIndex(mongo.SingleIndex("token_id", "1", false))
		c.EnsureIndex(mongo.SingleIndex("token_type", "1", false))
	})
}
