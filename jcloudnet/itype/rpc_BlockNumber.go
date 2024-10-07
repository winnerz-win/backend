package itype

import (
	"context"
	"jtools/jmath"
	"math/big"
)

func (my IClient) BlockNumber(ctx context.Context) (*big.Int, error) {
	req := ReqJsonRpc{
		Method: _nmap["blockNumber"][my.isKlay],
	}
	r, err := req.Request(my.rpcURL, my.isDebug)
	if err != nil {
		return nil, err
	}

	return jmath.BigInt(r.Result), nil
}
