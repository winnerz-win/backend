package itype

import (
	"jtools/jmath"
)

func (my IClient) SendRawTransaction(raw []byte) (string, error) {

	req := ReqJsonRpc{
		Method: _nmap["sendRawTransaction"][my.isKlay],
		Params: []any{
			param_hex_amount(raw),
		},
	}
	ack, err := req.Request(my.rpcURL, my.isDebug)
	if err != nil {
		return "", err
	}

	hash := jmath.HEX(ack.Result)

	return hash, nil
}

func (my IClient) SendRawTransactionHex(hex string) (string, error) {

	req := ReqJsonRpc{
		Method: _nmap["sendRawTransaction"][my.isKlay],
		Params: []any{
			hex,
		},
	}
	ack, err := req.Request(my.rpcURL, my.isDebug)
	if err != nil {
		return "", err
	}

	hash := jmath.HEX(ack.Result)

	return hash, nil
}
