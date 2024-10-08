package nwtypes

import (
	"jtools/unix"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/nft_winners/nwdb"
)

type DataSetBaseURI struct {
	ReceiptCode RECEIPT_CODE `bson:"receipt_code" json:"receipt_code"`
	NewURI      string       `bson:"new_uri" json:"new_uri"`
	IsCallback  bool         `bson:"is_callback" json:"is_callback"`

	ResultGasInfo `bson:",inline" json:",inline"`
}

func (my DataSetBaseURI) Valid() bool { return my.NewURI != "" }

func (my DataSetBaseURI) Data() mongo.MAP {
	return mongo.MakeMap(my)
}

type NftSetBaseURIResult struct {
	ReceiptCode RECEIPT_CODE `bson:"receipt_code" json:"receipt_code"`
	BaseURI     string       `bson:"base_uri" json:"base_uri"`

	Hash string `bson:"hash" json:"hash"`

	IsSucess bool `bson:"is_sucess" json:"is_sucess"`

	IsSend bool `bson:"is_send" json:"is_send"`

	Timestamp unix.Time `bson:"timestamp" json:"timestamp"`

	SendAt  unix.Time `bson:"send_at" json:"send_at"`
	SendMsg string    `bson:"send_msg,omitempty" json:"send_msg,omitempty"`
}

func (my NftSetBaseURIResult) SendOK(db mongo.DATABASE, nowAt unix.Time) {
	upQuery := mongo.Bson{"$set": mongo.Bson{
		"send_at": nowAt,
		"is_send": true,
	}}
	db.C(nwdb.NftSetBaseURIResult).Update(my.Selector(), upQuery)
}

func (my NftSetBaseURIResult) Selector() mongo.Bson {
	return mongo.Bson{"receipt_code": my.ReceiptCode}
}

func (NftSetBaseURIResult) IndexingDB(db mongo.DATABASE) {
	c := db.C(nwdb.NftSetBaseURIResult)
	c.EnsureIndex(mongo.SingleIndex("receipt_code", 1, true))
	c.EnsureIndex(mongo.SingleIndex("timestamp", 1, false))

	c.EnsureIndex(mongo.SingleIndex("is_send", 1, false))
}

func (my NftSetBaseURIResult) InsertDB(db mongo.DATABASE) {
	c := db.C(nwdb.NftSetBaseURIResult)
	c.Insert(my)
}
