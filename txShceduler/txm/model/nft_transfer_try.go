package model

import (
	"jtools/mms"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/txm/inf"
)

type NftTransferTry struct {
	User    `bson:",inline" json:",inline"`
	TokenId string `bson:"token_id" json:"token_id"`

	ToAddress  string `bson:"to_address" json:"to_address"`
	IsMemberTo bool   `bson:"is_member_to" json:"is_member_to"`

	Hash string `bson:"hash" json:"hash"`

	CreateTime mms.MMS `bson:"create_time" json:"create_time"`
	CreateYMD  int     `bson:"create_ymd" json:"create_ymd"`
}

func (my NftTransferTry) Valid() bool { return my.TokenId != "" }

type NftTransferTryList []NftTransferTry

func (my NftTransferTry) SelectorTry() mongo.Bson { return mongo.Bson{"token_id": my.TokenId} }

func (my NftTransferTry) InsertDB(db mongo.DATABASE) error {
	nowAt := mms.Now()
	my.CreateTime = nowAt
	my.CreateYMD = nowAt.YMD()
	return db.C(inf.NFTTransferTry).Insert(my)
}
func (my NftTransferTry) RemoveTryDB(db mongo.DATABASE) error {
	return db.C(inf.NFTTransferTry).Remove(my.SelectorTry())
}

func (my NftTransferTry) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.NFTTransferTry)
		c.EnsureIndex(mongo.SingleIndex("uid", "1", false))
		c.EnsureIndex(mongo.SingleIndex("address", "1", false))
		c.EnsureIndex(mongo.SingleIndex("name", "1", false))

		c.EnsureIndex(mongo.SingleIndex("token_id", "1", true))

		c.EnsureIndex(mongo.SingleIndex("Hash", "1", false))

		c.EnsureIndex(mongo.SingleIndex("create_time", "1", false))
		c.EnsureIndex(mongo.SingleIndex("create_ymd", "1", false))
	})
}

type NftTransferEnd struct {
	NftTransferTry `bson:",inline" json:",inline"`

	Key string `bson:"key" json:"-"`

	Status      int         `bson:"status" json:"status"`
	FailCode    FailNftCode `bson:"fail_code" json:"fail_code"`
	FailMessage string      `bson:"fail_message" json:"fail_message"`

	Timestamp mms.MMS `bson:"timestamp" json:"timestamp"`
	YMD       int     `bson:"ymd" json:"ymd"`

	GasFeeETH string `bson:"gas_fee_eth" json:"gas_fee_eth"`

	IsSend  bool    `bson:"is_send" json:"is_send"`
	SendAt  mms.MMS `bson:"send_at" json:"send_at"`
	SendYMD int     `bson:"send_ymd" json:"send_ymd"`
}

type NftTransferEndList []NftTransferEnd

func (my NftTransferEnd) CallbackData() chttp.JsonType {
	data := chttp.JsonType{
		"uid":      my.UID,
		"address":  my.Address,
		"name":     my.Name,
		"token_id": my.TokenId,

		"to_address":   my.ToAddress,
		"is_member_to": my.IsMemberTo,

		"hash": my.Hash,

		"status":       my.Status,
		"fail_code":    my.FailCode,
		"fail_message": my.FailMessage,

		"timestamp": my.Timestamp,

		"gas_fee_eth": my.GasFeeETH,
	}
	return data
}
func (my *NftTransferEnd) SetFail(code FailNftCode, exMsg ...string) {
	my.Status = 104
	my.FailCode = code
	my.FailMessage = NFTFailCodeMessage(code, exMsg...)
}
func (my NftTransferEnd) InsertEndDB(db mongo.DATABASE) error {
	my.Timestamp = mms.Now()
	my.YMD = my.Timestamp.YMD()
	return db.C(inf.NFTTransferEnd).Insert(my)
}

func (my NftTransferEnd) SendOK(db mongo.DATABASE, nowAt mms.MMS) error {
	my.IsSend = true
	my.SendAt = nowAt
	my.SendYMD = nowAt.YMD()
	return db.C(inf.NFTTransferEnd).Update(my.SelectorEnd(), my)
}

func (my NftTransferEnd) SelectorEnd() mongo.Bson { return mongo.Bson{"key": my.Key} }

func (my NftTransferEnd) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.NFTTransferEnd)
		c.EnsureIndex(mongo.SingleIndex("key", "1", true))

		c.EnsureIndex(mongo.SingleIndex("uid", "1", false))
		c.EnsureIndex(mongo.SingleIndex("address", "1", false))
		c.EnsureIndex(mongo.SingleIndex("name", "1", false))

		c.EnsureIndex(mongo.SingleIndex("token_id", "1", false))

		c.EnsureIndex(mongo.SingleIndex("Hash", "1", false))

		c.EnsureIndex(mongo.SingleIndex("timestamp", "1", false))
		c.EnsureIndex(mongo.SingleIndex("ymd", "1", false))

		c.EnsureIndex(mongo.SingleIndex("is_send", "1", false))
		c.EnsureIndex(mongo.SingleIndex("send_at", "1", false))
		c.EnsureIndex(mongo.SingleIndex("send_ymd", "1", false))
	})
}
