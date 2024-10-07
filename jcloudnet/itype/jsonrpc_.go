package itype

import (
	"encoding/hex"
	"errors"
	"fmt"
	"jtools/cc"
	"jtools/cloud/ebcm/abi"
	"jtools/dbg"
	"jtools/jmath"
	"jtools/jnet/cnet"
	"strings"
)

type RpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (my RpcError) String() string { return dbg.ToJsonString(my) }
func (my RpcError) Error(method string) error {
	return fmt.Errorf(`{"method":"%v","code":%v,"message":"%v"}`, method, my.Code, my.Message)
}

type AckJsonRpc struct {
	Method string `json:"method,omitempty"`
	/////////////////////////////////////////////////

	Jsonrpc string      `json:"jsonrpc"` //2.0
	Id      interface{} `json:"id"`      //1
	Result  interface{} `json:"result"`
	/*
		"error": {
			"code": -32602,
			"message": "invalid argument 0: json: cannot unmarshal hex string of odd length into Go struct field CallArgs.to of type common.Address"
		}
	*/
	Error *RpcError `json:"error,omitempty"`
}

func (my AckJsonRpc) String() string { return dbg.ToJsonString(my) }
func (my AckJsonRpc) ResultBytes() []byte {
	if hexString, do := my.Result.(string); do {
		return abi.HexToBytes(hexString)
	}
	return []byte{}
}

/////////////////////////////////////////////////////////////////////////////////////////

type ReqJsonRpc struct {
	Jsonrpc string `json:"jsonrpc"` //2.0
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	Id      int    `json:"id"` //1
}

func (my ReqJsonRpc) String() string { return dbg.ToJsonString(my) }

const (
	INVALID_EOC_ADDRESS_FORMAT = "<jcloudnet> invalid EOC address format."
)

func (my ReqJsonRpc) Request(rpc_url string, is_debugs ...bool) (*AckJsonRpc, error) {
	my.Jsonrpc = "2.0"
	my.Id = 1

	is_debug := _IsTrue(is_debugs)
	if is_debug {
		_debug_log(my)
	}

	ack := cnet.POST_STRUCT[AckJsonRpc](
		rpc_url,
		nil,
		my,
	)
	if err := ack.Error(); err != nil {
		if is_debug {
			cc.Red(err)
		}
		return nil,
			RpcError{
				Code:    0,
				Message: err.Error(),
			}.Error(my.Method)
	}

	result := ack.Item()
	result.Method = my.Method //
	if v, do := result.Result.(string); do {
		if v == "0x" {
			return nil, errors.New(INVALID_EOC_ADDRESS_FORMAT)
		}
	}

	if is_debug {
		_debug_log(result)
	}
	if result.Error != nil {
		return nil, result.Error.Error(my.Method)
	}
	return &result, nil
}

// ///////////////////////////////////////////////////////////////////////////////////////

// Encode encodes b as a hex string with 0x prefix.
func _Encode(b []byte) string {
	enc := make([]byte, len(b)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], b)
	return string(enc)
}

func param_hex_amount(v any) string {
	if b, do := v.([]byte); do {
		return _Encode(b)
	}
	hv := jmath.HEX(v)
	if len(hv) >= 4 { //0x01, 0x0f
		if strings.HasPrefix(hv, "0x0") {
			hv = "0x" + hv[3:]
		}
	}
	if hv == "0x" {
		hv = "0x0"
	}
	return hv
	//return fmt.Sprintf("0x%x", v)
}

func _0xToLower(v any) string {
	return strings.ToLower(fmt.Sprintf("%v", v))
}
