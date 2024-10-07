package itype

import (
	"context"
	"jtools/jmath"
	"math/big"
)

func (my IClient) ChainID() *big.Int {
	req := ReqJsonRpc{
		Method: _nmap["chainId"][my.isKlay],
	}
	r, err := req.Request(my.rpcURL, my.isDebug)
	if err != nil {
		return nil
	}

	return jmath.BigInt(r.Result)
}

func (my IClient) NetworkID(ctx context.Context) (*big.Int, error) {
	req := ReqJsonRpc{
		Method: _nmap["chainId"][my.isKlay],
	}
	r, err := req.Request(my.rpcURL, my.isDebug)
	if err != nil {
		return nil, err
	}

	return jmath.BigInt(r.Result), nil
}
