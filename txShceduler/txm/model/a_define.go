package model

import (
	"jtools/cloud/ebcm"
	"strings"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/txm/inf"
)

const (
	ZERO = "0"

	V1      = "/v1"
	V1Owner = "/v1/owner" //inf.IsLockTransferMode()

	NFT = "/nft"

	ChanBuffers = 100
)

// DB :
func DB(f func(db mongo.DATABASE)) error {
	return inf.DB().Action(inf.DBName, func(db mongo.DATABASE) {
		f(db)
	})
}

// Trim :
func Trim(v *string) {
	*v = strings.TrimSpace(*v)
}

func CALC_GAS_PRICE(db mongo.DATABASE, gas_price ebcm.GasPrice) ebcm.GasPrice {
	mul_value := GetGAS(db).GetGasMulValue()
	gas_price.MUL_VALUE(mul_value)
	return gas_price
}
