package nwtypes

import (
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/unix"
	"txscheduler/nft_winners/nwdb"
)

type NftMasterTry struct {
	ReceiptCode RECEIPT_CODE `bson:"receipt_code" json:"receipt_code"`

	DATA_SEQ DATA_FLOW `bson:"data_seq" json:"data_seq"`

	TimeTryAt unix.Time `bson:"time_try_at" json:"time_try_at"`
}

func (my NftMasterTry) Selector() mongo.Bson { return mongo.Bson{"receipt_code": my.ReceiptCode} }

func (NftMasterTry) IndexingDB(db mongo.DATABASE) {
	c := db.C(nwdb.NftMasterTry)
	c.EnsureIndex(mongo.SingleIndex("receipt_code", 1, true))

	c.EnsureIndex(mongo.SingleIndex("time_try_at", 1, false))

}

func (my NftMasterTry) InsertDB(db mongo.DATABASE) error {
	return db.C(nwdb.NftMasterTry).Insert(my)
}

func (my NftMasterTry) RemoveTry(db mongo.DATABASE) error {
	return db.C(nwdb.NftMasterTry).Remove(my.Selector())
}
