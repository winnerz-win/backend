package nwtypes

import (
	"jtools/cloud/jeth/jwallet"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/txm/inf"
)

type WalletInfo struct {
	IsMaster bool   `bson:"is_master" json:"is_master"`
	Address  string `bson:"address" json:"address"`
	UID      int64  `bson:"uid" json:"uid"`
	Name     string `bson:"name,omitempty" json:"name,omitempty"`
}

func MasterWalletInfo() WalletInfo {
	return WalletInfo{
		IsMaster: true,
		Address:  "",
		UID:      0,
	}
}

func (my WalletInfo) UserPrivateKey() string {
	if my.UID == 0 {
		dbg.RedItalic("nwtypes.WalletInfo.UserPrivateKey :", dbg.Stack())
		return ""
	}
	seed := inf.Config().Seed
	wallet := jwallet.NewSeed(seed, my.UID)
	return wallet.PrivateKey()
}

type FeeInfo struct {
	PlatformAddress string `json:"platform_address"`
	PlatformPrice   string `json:"platform_price"`

	BenefitAddress string `json:"benefit_address"`
	BenefitPrice   string `json:"benefit_price"`
}

type NftInfo struct {
	Contract string `bson:"contract" json:"contract"`
	Symbol   string `bson:"symbol" json:"symbol"`
	TokenId  string `bson:"token_id" json:"token_id"`
}

type DataMintTry struct {
	ReceiptCode RECEIPT_CODE `bson:"receipt_code" json:"receipt_code"`

	Payer WalletInfo `bson:"payer" json:"payer"`
	Owner WalletInfo `bson:"owner" json:"owner"`
	Kind  MintKind   `bson:"kind" json:"kind"`

	NftInfo NftInfo `bson:"nft_info" json:"nft_info"`

	PayTokenInfo *inf.TokenInfo `bson:"pay_token_info" json:"pay_token_info"` // nil == free

	FeeInfo *FeeInfo `bson:"fee_info" json:"fee_info"`

	ResultGasInfo `bson:",inline" json:",inline"`

	//TimeTryAt unix.Time `bson:"time_try_at" json:"time_try_at"`

	//PreData mongo.MAP `bson:"pre_data" json:"pre_data"` //DataMultiTransferTry(mint)
}

func (my DataMintTry) Valid() bool { return my.ReceiptCode.Valid() }

func (my DataMintTry) PayTokenInfoSymbol() string {
	if my.PayTokenInfo != nil {
		return my.PayTokenInfo.Symbol
	}
	return "FREE"
}

func (my DataMintTry) IsCoinTry() bool {
	if my.PayTokenInfo != nil {
		return my.PayTokenInfo.IsCoin
	}
	return false
}

func (my DataMintTry) Data() mongo.MAP {
	return mongo.MakeMap(my)
}

func (my DataMintTry) ResultMint() ResultMint {
	r := ResultMint{
		Payer: my.Payer,
		Owner: my.Owner,
		Kind:  my.Kind,

		NftInfo:       my.NftInfo,
		ResultGasInfo: my.ResultGasInfo,
	}
	return r
}
