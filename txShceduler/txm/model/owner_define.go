package model

import (
	"jtools/cloud/ebcm"
	"jtools/dec"
	"jtools/unix"
	"time"
	"txscheduler/brix/tools/database/mongo/tools/jmath"
)

const (
	OwnerCallbackApi_LT       = "/v1/owner/lock_transfer/callback"
	OwnerCallbackApi_Transfer = "/v1/owner/transfer/callback"
	OwnerCallbackApi_Lock     = "/v1/owner/lock/callback"
	OwnerCallbackApi_Unlock   = "/v1/owner/unlock/callback"
	OwnerCallbackApi_Relock   = "/v1/owner/relock/callback"
)

type TxTryData struct {
	Nonce dec.Uint256 `bson:"nonce" json:"nonce"`
	To    string      `bson:"to" json:"to"`
	Value string      `bson:"value" json:"value"`

	Limit string `bson:"limit" json:"limit"`
	Gwei  string `bson:"gwei" json:"gwei"`
	Hex   string `bson:"hex" json:"hex"`

	Time    unix.Time `bson:"time" json:"time"`
	Counter int       `bson:"counter" json:"counter"`
}

func MakeTxTryData(
	sender *ebcm.Sender,
	private_key string,

	nonce any,
	to string,
	value string,
	limit any,
	gas_price ebcm.GasPrice,
	data ebcm.PADBYTES,

) (TxTryData, ebcm.WrappedTransaction) {
	tx_try_data := TxTryData{
		Nonce: dec.UINT256(jmath.VALUE(nonce)),
		To:    to,
		Value: value,

		Limit: jmath.VALUE(limit),
		Gwei:  gas_price.GET_GAS_GWEI(),
		Hex:   data.Hex(),

		Time:    unix.Now(),
		Counter: 1,
	}

	stx, _ := sender.SignTx(
		sender.NewTransaction(
			jmath.Uint64(nonce),
			to,
			value,
			jmath.Uint64(limit),
			gas_price,
			data,
		),
		private_key,
	)
	return tx_try_data, stx
}

func (my TxTryData) IsTimeOver(dst time.Duration) bool {
	du := unix.Now().Sub(my.Time)
	return du >= dst
}

func (my *TxTryData) STX(
	sender *ebcm.Sender,
	private_key string,
	//mul_val any,
) ebcm.WrappedTransaction {
	my.Gwei = jmath.DOTCUT(
		jmath.MUL(my.Gwei, 1.5),
		0,
	)
	my.Counter++

	nonce := jmath.Uint64(my.Nonce.Int64())
	limit := jmath.Uint64(my.Limit)
	gas_price := ebcm.MakeGasPriceGWEI(my.Gwei)
	data := ebcm.PADBYTESFromHex(my.Hex)

	stx, _ := sender.SignTx(
		sender.NewTransaction(
			nonce,
			my.To,
			my.Value,
			limit,
			gas_price,
			data,
		),
		private_key,
	)

	return stx
}
