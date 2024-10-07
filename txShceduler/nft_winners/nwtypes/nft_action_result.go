package nwtypes

import (
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/unix"
	"txscheduler/nft_winners/nwdb"
	"txscheduler/txm/inf"
)

type ResultGasInfo struct {
	IsSuccess  bool   `bson:"-" json:"-"`
	Hash       string `bson:"hash" json:"hash"`
	TxGasPrice string `bson:"tx_gas_price" json:"tx_gas_price"`

	Limit uint64 `bson:"-" json:"-"`
}

type ResultPayInfo struct {
	PayTokenInfo inf.TokenInfo `bson:"pay_token_info" json:"pay_token_info"`
	PayFrom      WalletInfo    `bson:"pay_from" json:"pay_from"`

	PlatformAddress  string `bson:"platform_address,omitempty" json:"platform_address,omitempty"`
	PlatformFeePrice string `bson:"platform_fee_price,omitempty" json:"platform_fee_price,omitempty"`

	BenefitAddress  string `bson:"benefit_address,omitempty" json:"benefit_address,omitempty"`
	BenefitFeePrice string `bson:"benefit_fee_price,omitempty" json:"benefit_fee_price,omitempty"`

	UserAddress  string `bson:"user_address,omitempty" json:"user_address,omitempty"`
	UserUID      int64  `bson:"user_uid,omitempty" json:"user_uid,omitempty"`
	UserName     string `bson:"user_name,omitempty" json:"user_name,omitempty"`
	UserPayPrice string `bson:"user_pay_price,omitempty" json:"user_pay_price,omitempty"`

	ResultGasInfo `bson:",inline" json:",inline"`
}

type MintKind string

const (
	MintKindSelf = MintKind("self")
	MintKindGift = MintKind("gift")
	MintKindFree = MintKind("free")
)

type ResultMint struct {
	Payer WalletInfo `bson:"payer" json:"payer"`
	Owner WalletInfo `bson:"owner" json:"owner"`
	Kind  MintKind   `bson:"kind" json:"kind"`

	NftInfo NftInfo `bson:"nft_info" json:"nft_info"`

	ResultGasInfo `bson:",inline" json:",inline"`
}

type ResultMintData struct {
	MintInfo ResultMint     `bson:"mint_info" json:"mint_info"`
	PayInfo  *ResultPayInfo `bson:"pay_info" json:"pay_info"`
}

func (my ResultMintData) Data() mongo.MAP {
	return mongo.MakeMap(my)
}

type ResultNftTransfer struct {
	From WalletInfo `bson:"from" json:"from"`
	To   WalletInfo `bson:"to" json:"to"`

	NftInfo NftInfo `bson:"nft_info" json:"nft_info"`

	ResultGasInfo `bson:",inline" json:",inline"`
}

type ResultSaleData struct {
	NftTransferInfo ResultNftTransfer `bson:"nft_transfer_info" json:"nft_transfer_info"`
	PayInfo         *ResultPayInfo    `bson:"pay_info" json:"pay_info"`
}

func (my ResultSaleData) Data() mongo.MAP {
	return mongo.MakeMap(my)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////

type NftActionResult struct {
	ReceiptCode RECEIPT_CODE `bson:"receipt_code" json:"receipt_code"`

	ResultType RESULT_TYPE `bson:"result_type" json:"result_type"` //mint

	Data mongo.MAP `bson:"data" json:"data"`

	TimeTryAt unix.Time `bson:"time_try_at" json:"time_try_at"`
	InsertAt  unix.Time `bson:"insert_at" json:"insert_at"`

	IsSend  bool      `bson:"is_send" json:"is_send"`
	SendAt  unix.Time `bson:"send_at" json:"send_at"`
	SendYMD int       `bson:"send_ymd" json:"send_ymd"`
}

func (my NftActionResult) Selector() mongo.Bson { return mongo.Bson{"receipt_code": my.ReceiptCode} }

func (NftActionResult) IndexingDB(db mongo.DATABASE) {
	c := db.C(nwdb.NftActionResult)
	c.EnsureIndex(mongo.SingleIndex("receipt_code", 1, true))
	c.EnsureIndex(mongo.SingleIndex("insert_at", 1, false))

	c.EnsureIndex(mongo.SingleIndex("is_send", 1, false))
}

func (my NftActionResult) InsertDB(db mongo.DATABASE) error {
	return db.C(nwdb.NftActionResult).Insert(my)
}

func (my NftActionResult) SendOK(db mongo.DATABASE, nowAt unix.Time) {
	upQuery := mongo.Bson{"$set": mongo.Bson{
		"send_at":  nowAt,
		"send_ymd": nowAt.YMD(),
		"is_send":  true,
	}}
	db.C(nwdb.NftActionResult).Update(my.Selector(), upQuery)
}
