package model

import (
	"strings"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/txm/inf"
)

const (
	ZERO = "0"

	V1  = "/v1"
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
