package itype

import (
	"jtools/cc"
	"strings"
)

var (
	_nmap = map[string]map[bool]string{
		"call": {
			false: "eth_call",
			true:  "klay_call",
		},
		"getBalance": {
			false: "eth_getBalance",
			true:  "klay_getBalance",
		},
		"blockNumber": {
			false: "eth_blockNumber",
			true:  "klay_blockNumber",
		},
		"chainId": {
			false: "eth_chainId",
			true:  "klay_chainId",
		},
		"estimateGas": {
			false: "eth_estimateGas",
			true:  "klay_estimateGas",
		},
		"getTransactionReceipt": {
			false: "eth_getTransactionReceipt",
			true:  "klay_getTransactionReceipt",
		},
		"getTransactionByHash": {
			false: "eth_getTransactionByHash",
			true:  "klay_getTransactionByHash",
		},
		"getBlockByNumber": {
			false: "eth_getBlockByNumber",
			true:  "klay_getBlockByNumber",
		},
		"getTransactionCount": {
			false: "eth_getTransactionCount",
			true:  "klay_getTransactionCount",
		},
		"gasPrice": {
			false: "eth_gasPrice",
			true:  "klay_gasPrice",
		},
		"sendRawTransaction": {
			false: "eth_sendRawTransaction",
			true:  "klay_sendRawTransaction",
		},
	}
)

func _IsTrue(p interface{}) bool {
	switch v := p.(type) {
	case bool:
		return v
	case []bool:
		if len(v) > 0 {
			return v[0]
		}
	case string:
		return strings.ToLower(strings.TrimSpace(v)) == "true"
	case []interface{}:
		if len(v) > 0 {
			return _IsTrue(v[0])
		}
	}
	return false
}

func _debug_log(a ...interface{}) {
	cc.Gray(a...)
}
