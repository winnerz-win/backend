package nwtypes

import "txscheduler/brix/tools/database/mongo"

type DataNftTransfer struct {
	ReceiptCode RECEIPT_CODE `bson:"receipt_code" json:"receipt_code"`

	From WalletInfo `bson:"from" json:"from"`
	To   WalletInfo `bson:"to" json:"to"`

	NftInfo NftInfo `bson:"nft_info" json:"nft_info"`

	ResultGasInfo `bson:",inline" json:",inline"`
}

func (my DataNftTransfer) Valid() bool { return my.ReceiptCode.Valid() }

func (my DataNftTransfer) Data() mongo.MAP {
	return mongo.MakeMap(my)
}

func (my DataNftTransfer) ResultNftTransfer() ResultNftTransfer {
	r := ResultNftTransfer{
		From:          my.From,
		To:            my.To,
		NftInfo:       my.NftInfo,
		ResultGasInfo: my.ResultGasInfo,
	}
	return r
}
