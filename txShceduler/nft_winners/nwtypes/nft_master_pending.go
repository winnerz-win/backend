package nwtypes

import (
	"jtools/unix"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/nft_winners/nwdb"
)

type NftMasterPending struct {
	ReceiptCode RECEIPT_CODE `bson:"receipt_code" json:"receipt_code"`

	DATA_SEQ DATA_FLOW `bson:"data_seq" json:"data_seq"`

	Hash string `bson:"hash" json:"hash"`

	TimeTryAt     unix.Time `bson:"time_try_at" json:"time_try_at"`
	TimePendingAt unix.Time `bson:"time_pending_at" json:"time_pending_at"`
}

func (my NftMasterPending) Selector() mongo.Bson { return mongo.Bson{"receipt_code": my.ReceiptCode} }

func (NftMasterPending) IndexingDB(db mongo.DATABASE) {
	c := db.C(nwdb.NftMasterPending)
	c.EnsureIndex(mongo.SingleIndex("receipt_code", 1, true))

	c.EnsureIndex(mongo.SingleIndex("time_pending_at", 1, false))

}

func (my NftMasterPending) InsertDB(db mongo.DATABASE) error {
	return db.C(nwdb.NftMasterPending).Insert(my)
}

func (my NftMasterPending) RemovePending(db mongo.DATABASE) error {
	return db.C(nwdb.NftMasterPending).Remove(my.Selector())
}

// events := rpc.WinnerzNFT.ParseEvent(
// 	tx,
// )
// for _, event := range events {
// 	if event.Name == "TransferNFT" {
// 		transfer_nft := event.Data.(rpc.EventTransferNFT)
// 		result_data.NftTokenId = transfer_nft.TokenId
// 		break
// 	}
// } //for
