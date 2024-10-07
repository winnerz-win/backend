package itype

import (
	"context"
	"jtools/cloud/ebcm"
	"jtools/cloud/ebcm/abi"
	"jtools/jmath"
	"math/big"
)

func ADDRESS(i interface{}) string {
	v := jmath.HEX(i, true)
	loop := 40 - len(v)
	for i := 0; i < loop; i++ {
		v = "0" + v
	}
	return "0x" + v
}

func (my IClient) BalanceAt(ctx context.Context, account interface{}, blockNumber *big.Int) (*big.Int, error) {

	tail := "latest"
	if blockNumber != nil {
		tail = param_hex_amount(blockNumber)
	}

	req := ReqJsonRpc{
		Method: _nmap["getBalance"][my.isKlay],
		Params: []interface{}{
			ADDRESS(account),
			tail,
		},
	}
	r, err := req.Request(my.rpcURL, my.isDebug)
	if err != nil {
		return nil, err
	}
	return jmath.BigInt(r.Result), nil
}

func (my IClient) GetCoinBalance(account interface{}) string {
	val, err := my.BalanceAt(
		context.Background(),
		account,
		nil,
	)
	if err != nil {
		return "0"
	}
	return jmath.VALUE(val)
}
func (my IClient) GetCoinPrice(account interface{}) string {
	return ebcm.WeiToETH(my.GetCoinBalance(account))
}

///////////////////////////////////////////////////////////////////////////////////////////

// TokenBalance : (account , contract/eth) wei
func (my IClient) TokenBalance(account string, contract string) string {

	//coin
	if !ebcm.IsAddress(contract) {
		v, err := my.BalanceAt(context.Background(), account, nil)
		if err != nil {
			v = jmath.BigInt(0)
		}
		return jmath.VALUE(v)
	}

	//token
	balance := "0"
	my.Call(
		contract,
		abi.Method{
			Name: "balanceOf",
			Params: abi.NewParams(
				abi.NewAddress(account),
			),
			Returns: abi.NewReturns(
				abi.Uint256,
			),
		},
		contract,
		func(rs abi.RESULT) {
			balance = rs.Uint(0)
		},
	)
	return balance
}

// Price : (account , contract/eth) price
func (my IClient) Price(account string, contract string, decimal interface{}) string {
	balance := my.TokenBalance(account, contract)

	return ebcm.WeiToToken(balance, decimal)
}
