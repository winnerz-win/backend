package model

import (
	"jtools/cloud/ebcm"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/txm/inf"
)

type NftTx struct {
	ebcm.TransactionBlock `bson:",inline" json:",inline"`
	Logs                  ebcm.TxLogList `bson:"logs" json:"logs"`

	Number int64 `bson:"number" json:"number"`
	//IsError              bool          `bson:"is_error" json:"is_error"`
}

func (my NftTx) String() string { return dbg.ToJSONString(my) }

func (my NftTx) Selector() mongo.Bson { return mongo.Bson{"hash": my.Hash} }

func (my NftTx) InsertDB(db mongo.DATABASE) error {
	return db.C(inf.NFTTX).Insert(my)
}

func (my NftTx) UpdateFuncNameDB(db mongo.DATABASE, mName string) {
	upQuery := mongo.Bson{"$set": mongo.Bson{
		"contract_method": mName,
	}}
	db.C(inf.NFTTX).Update(my.Selector(), upQuery)
}

func (NftTx) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.NFTTX)
		c.EnsureIndex(mongo.SingleIndex("hash", "1", true))

		c.EnsureIndex(mongo.SingleIndex("contract_method", "1", false))

		c.EnsureIndex(mongo.SingleIndex("number", "1", false))
		c.EnsureIndex(mongo.SingleIndex("is_error", "1", false))
		c.EnsureIndex(mongo.SingleIndex("tx_index", "1", false))

		c.EnsureIndex(mongo.SingleIndex("timestamp", "1", false))

		c.EnsureIndex(mongo.SingleIndex("logs.address", "1", false))
	})
}

//////////////////////////////////////////////////////////////////////////

type NftCache struct {
	Number int64                     `bson:"number" json:"number"`
	List   ebcm.TransactionBlockList `bson:"list" json:"list"`
}
type NftCacheList []NftCache

func (my NftCache) Selector() mongo.Bson { return mongo.Bson{"number": my.Number} }

func (my NftCache) InsertDB(db mongo.DATABASE) error {
	return db.C(inf.NFTCache).Insert(my)
}
func (my NftCache) RemoveDB(db mongo.DATABASE) error {
	return db.C(inf.NFTCache).Remove(my.Selector())
}

func (NftCache) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.NFTCache)
		c.EnsureIndex(mongo.SingleIndex("number", "1", true))
	})
}
