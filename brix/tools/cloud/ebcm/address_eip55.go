package ebcm

import (
	"txscheduler/brix/tools/cloud/ebcm/abi"
)

/*
	https://steemit.com/kr/@anpigon/ethereum-1
*/

func EIP55(address string) string {
	return abi.EIP55(address)
}
