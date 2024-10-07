package ecsx

import (
	"strings"
	"txscheduler/brix/tools/cloudx/ebcmx"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// IsAddress :
func IsAddress(address string) bool {
	if !strings.HasPrefix(address, "0x") {
		return false
	}
	//return common.IsHexAddress(strings.ToLower(address))
	return ebcmx.EIP55(address) != ""
}

// ContractAddressNonce :
func ContractAddressNonce(from string, nonce uint64) string {
	if !IsAddress(from) {
		return ""
	}

	v := crypto.CreateAddress(
		common.HexToAddress(from),
		nonce,
	)
	return strings.ToLower(v.String())
}
