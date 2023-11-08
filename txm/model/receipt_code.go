package model

import (
	"txscheduler/brix/tools/crypt"
	"txscheduler/brix/tools/dbg"
)

//GetReceiptCode :
func GetReceiptCode() string {
	uuid := crypt.MakeUID256()
	prefix := "receipt_"

	cut := len(uuid)/2 + len(prefix)
	code := prefix + uuid[:cut]

	return code
}

//GetReceiptCodeSELF :
func GetReceiptCodeSELF() string {
	uuid := crypt.MakeUID256()
	prefix := "self_"

	cut := len(uuid)/2 + len(prefix)
	code := prefix + uuid[:cut]

	return code
}

func GetMasterCode() string {
	uuid := crypt.MakeUInt64()
	return dbg.Cat("master_", uuid)

}

func NFTReceiptCode(symbol string) string {
	symbol = dbg.TrimToLower(symbol)
	prefix := ""
	if symbol == "eth" {
		prefix += "eth_" + crypt.MakeUInt64String()
	} else {
		prefix += "token_" + crypt.MakeUInt64String()
	}

	prefix += "_nft"

	return prefix
}
