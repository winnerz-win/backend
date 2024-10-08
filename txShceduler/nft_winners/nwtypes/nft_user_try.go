package nwtypes

import (
	"jtools/unix"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/nft_winners/nwdb"
)

type NftUserTry struct {
	ReceiptCode RECEIPT_CODE `bson:"receipt_code" json:"receipt_code"`

	DATA_SEQ DATA_FLOW `bson:"data_seq" json:"data_seq"`

	TimeTryAt unix.Time `bson:"time_try_at" json:"time_try_at"`
}

func (my NftUserTry) Selector() mongo.Bson { return mongo.Bson{"receipt_code": my.ReceiptCode} }

func (NftUserTry) IndexingDB(db mongo.DATABASE) {
	c := db.C(nwdb.NftUserTry)
	c.EnsureIndex(mongo.SingleIndex("receipt_code", 1, true))

	c.EnsureIndex(mongo.SingleIndex("time_try_at", 1, false))

}

func (my NftUserTry) InsertDB(db mongo.DATABASE) error {
	return db.C(nwdb.NftUserTry).Insert(my)
}

func (my NftUserTry) RemoveTry(db mongo.DATABASE) error {
	return db.C(nwdb.NftUserTry).Remove(my.Selector())
}
