package nwtypes

// import (
// 	"jtools/unix"
// )

// type NftMintPendingData struct {
// 	DataMintTry `bson:",inline" json:",inline"`

// 	MintHash string `bson:"mint_hash" json:"mint_hash"`

// 	TimePendingAt unix.Time `bson:"time_pending_at" json:"time_pending_at"`
// }

// func (NftMintPending) IndexingDB(db mongo.DATABASE) {
// 	c := db.C(nwdb.NftMintPending)
// 	c.EnsureIndex(mongo.SingleIndex("receipt_code", 1, true))
// 	c.EnsureIndex(mongo.SingleIndex("token_info.is_coin", 1, false))

// 	c.EnsureIndex(mongo.SingleIndex("time_pending_at", 1, false))

// }

// func (my *NftMintPending) InsertDB(db mongo.DATABASE, nowAt unix.Time) error {
// 	my.TimePendingAt = nowAt
// 	return db.C(nwdb.NftMintPending).Insert(my)
// }

// func (my NftMintPending) RemovePending(db mongo.DATABASE) error {
// 	return db.C(nwdb.NftMintTry).Remove(my.Selector())
// }
