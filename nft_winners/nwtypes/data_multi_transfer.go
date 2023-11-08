package nwtypes

import (
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/txm/inf"
)

type InfoKind string

const (
	InfoKindMember   = InfoKind("member")
	InfoKindPlatform = InfoKind("platform")
	InfoKindBenefit  = InfoKind("benefit")
)

type PriceToInfo struct {
	InfoKind InfoKind   `bson:"info_kind" json:"info_kind"` //member, platform, benefit,
	To       WalletInfo `bson:"to" json:"to"`
	Price    string     `bson:"price" json:"price"`
}

type DataMultiTransferTry struct {
	ReceiptCode RECEIPT_CODE `bson:"receipt_code" json:"receipt_code"`

	PayTokenInfo inf.TokenInfo `bson:"pay_token_info" json:"pay_token_info"`
	PayFrom      WalletInfo    `bson:"pay_from" json:"pay_from"`

	PriceToInfos []PriceToInfo `bson:"price_to_infos" json:"price_to_infos"`

	ResultGasInfo `bson:",inline" json:",inline"`

	// PreData  mongo.MAP `bson:"pre_data" json:"pre_data"` //DataMintTry(mint) /
	// NextData mongo.MAP `bson:"next_data" json:"next_data"`
}

func (my DataMultiTransferTry) Data() mongo.MAP {
	return mongo.MakeMap(my)
}

func (my DataMultiTransferTry) ResultPayInfo() *ResultPayInfo {
	pay_info := &ResultPayInfo{
		PayTokenInfo:  my.PayTokenInfo,
		PayFrom:       my.PayFrom,
		ResultGasInfo: my.ResultGasInfo,
	}

	for _, v := range my.PriceToInfos {
		switch v.InfoKind {
		case InfoKindMember:
			pay_info.UserAddress = v.To.Address
			pay_info.UserUID = v.To.UID
			pay_info.UserName = v.To.Name
			pay_info.UserPayPrice = v.Price

		case InfoKindPlatform:
			pay_info.PlatformAddress = v.To.Address
			pay_info.PlatformFeePrice = v.Price

		case InfoKindBenefit:
			pay_info.BenefitAddress = v.To.Address
			pay_info.BenefitFeePrice = v.Price
		}
	} //for

	return pay_info
}
