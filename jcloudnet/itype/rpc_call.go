package itype

import (
	"jtools/cc"
	"jtools/cloud/ebcm/abi"
	"jtools/dbg"
)

func (my IClient) Call(
	contract string,
	method abi.Method,
	caller string,
	f func(rs abi.RESULT),
	isLogs ...bool,
) error {

	is_debug := dbg.IsTrue(isLogs)

	inputBytes := method.Params.GetBytes(method.Name)
	if is_debug {
		cc.Purple(inputBytes)
	}

	view_func_data := param_hex_amount(inputBytes)
	if is_debug {
		cc.Purple(view_func_data)
	}
	// cc.Purple(view_func_data)
	// cc.Purple("SIZE :", len(view_func_data))

	req := ReqJsonRpc{
		Method: _nmap["call"][my.isKlay],
		Params: []any{
			map[string]interface{}{
				"from": caller,
				"to":   contract,
				"data": view_func_data,
			},
			"latest",
		},
		Id:      1,
		Jsonrpc: "2.0",
	}
	ack, err := req.Request(my.rpcURL, my.isDebug)
	if err != nil {
		return err
	}

	receipt := ack.ResultBytes()

	result := abi.ReceiptDivDirect(
		receipt,
		method.Returns,
	)

	if is_debug {
		cc.Purple(result)
	}
	f(result)

	return nil
}
