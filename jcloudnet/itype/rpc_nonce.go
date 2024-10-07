package itype

import (
	"context"
	"jtools/jmath"
)

func (my IClient) getTransactionCount(account interface{}, tail string) (uint64, error) {
	req := ReqJsonRpc{
		Method: _nmap["getTransactionCount"][my.isKlay],
		Params: []interface{}{
			ADDRESS(account),
			tail,
		},
	}
	r, err := req.Request(my.rpcURL, my.isDebug)
	if err != nil {
		return 0, err
	}
	return jmath.Uint64(r.Result), nil
}

func (my IClient) PendingNonceAt(ctx context.Context, account interface{}) (uint64, error) {
	return my.getTransactionCount(account, "pending")
}

func (my IClient) NonceAt(ctx context.Context, account interface{}) (uint64, error) {
	return my.getTransactionCount(account, "latest")

}
