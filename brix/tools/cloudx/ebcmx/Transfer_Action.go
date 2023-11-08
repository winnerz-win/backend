package ebcmx

import "txscheduler/brix/tools/dbg"

type XGasPrice interface {
	String() string
	Error() error

	WEI() string
	ETH() string

	AddWEI(weiVal interface{})
	AddGWEI(gwei string)
	FeeWEI(limit uint64) string
	FeeETH(limit uint64) string

	Valid() bool
}

type XSendResult struct {
	From        string `bson:"from" json:"from"`
	To          string `bson:"to" json:"to"`
	PadBytesHex string `bson:"padBytesHex" json:"padBytesHex"` // hexString
	Nonce       uint64 `bson:"nonce" json:"nonce"`
	GasLimit    uint64 `bson:"gasLimit" json:"gasLimit"`
	GasPrice    string `bson:"gasPrice" json:"gasPrice"` // wei
	Hash        string `bson:"hash" json:"hash"`
}

func (my XSendResult) String() string { return dbg.ToJSONString(my) }
