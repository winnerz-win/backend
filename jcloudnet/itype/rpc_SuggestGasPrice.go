package itype

import (
	"context"
	"errors"
	"jtools/cloud/ebcm"
	"jtools/jmath"
	"math/big"
)

func (my IClient) eth_maxPriorityFeePerGas() (*big.Int, error) {
	if my.isKlay {
		return nil, errors.New("isKlay(true)")
	}
	req := ReqJsonRpc{
		Method: "eth_maxPriorityFeePerGas",
	}
	r, err := req.Request(my.rpcURL, my.isDebug)
	if err != nil {
		return nil, err
	}
	return jmath.BigInt(r.Result), nil
}

func (my IClient) SuggestGasPrice(ctx context.Context, is_skip_tip_cap ...bool) (ebcm.GasPrice, error) {
	gas_price := ebcm.GasPrice{}

	req := ReqJsonRpc{
		Method: _nmap["gasPrice"][my.isKlay],
	}
	r, err := req.Request(my.rpcURL, my.isDebug)
	if err != nil {
		return gas_price, err
	}

	gas_price = ebcm.GasPrice{
		Tip: jmath.BigInt(r.Result),
		Gas: jmath.BigInt(r.Result),
	}
	if my.isKlay {
		return gas_price, nil
	}

	if my.TXNTYPE() == ebcm.TXN_EIP_1559 {
		if _IsTrue(is_skip_tip_cap) {
			// gas_price = ebcm.GasPrice{
			// 	Tip: jmath.BigInt(r.Result),
			// 	Gas: jmath.BigInt(r.Result),
			// }

		} else {
			//eth_maxPriorityFeePerGas

			tip, err := my.eth_maxPriorityFeePerGas()
			if err != nil {
				return gas_price, err
			}

			gas_price.Gas = jmath.BigInt(
				jmath.DOTCUT(
					jmath.ADD(jmath.MUL(gas_price.Gas, 1.3), tip),
					0,
				),
			)
			gas_price.Tip = tip
		}
	}

	return gas_price, nil
}
