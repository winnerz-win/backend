package itype

import (
	"context"
	"jtools/cloud/ebcm"
	"jtools/jmath"
)

func (my IClient) EstimateGas(ctx context.Context, msg ebcm.CallMsg) (uint64, error) {
	param := map[string]interface{}{
		"from":  msg.From,
		"to":    msg.To,
		"value": param_hex_amount(msg.Value),
	}
	if len(msg.Data) > 0 {
		param["data"] = param_hex_amount(msg.Data)
	}

	req := ReqJsonRpc{
		Method: _nmap["estimateGas"][my.isKlay],
		Params: []any{
			param,
		},
	}

	r, err := req.Request(my.rpcURL, my.isDebug)
	if err != nil {
		return 0, err
	}

	return jmath.Uint64(r.Result), nil
}
