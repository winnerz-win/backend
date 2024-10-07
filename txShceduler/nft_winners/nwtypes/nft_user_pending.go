package nwtypes

import (
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/unix"
	"txscheduler/nft_winners/nwdb"
)

type NftUserPending struct {
	ReceiptCode RECEIPT_CODE `bson:"receipt_code" json:"receipt_code"`

	From string `bson:"from" json:"from"`

	DATA_SEQ DATA_FLOW `bson:"data_seq" json:"data_seq"`

	Hash string `bson:"hash" json:"hash"`

	TimeTryAt     unix.Time `bson:"time_try_at" json:"time_try_at"`
	TimePendingAt unix.Time `bson:"time_pending_at" json:"time_pending_at"`
}

func (my NftUserPending) Selector() mongo.Bson { return mongo.Bson{"receipt_code": my.ReceiptCode} }

func (NftUserPending) IndexingDB(db mongo.DATABASE) {
	c := db.C(nwdb.NftUserPending)
	c.EnsureIndex(mongo.SingleIndex("receipt_code", 1, true))

	c.EnsureIndex(mongo.SingleIndex("from", 1, false))
	c.EnsureIndex(mongo.SingleIndex("time_pending_at", 1, false))

}

func (my NftUserPending) InsertDB(db mongo.DATABASE) error {
	return db.C(nwdb.NftUserPending).Insert(my)
}

func (my NftUserPending) RemovePending(db mongo.DATABASE) error {
	return db.C(nwdb.NftUserPending).Remove(my.Selector())
}
