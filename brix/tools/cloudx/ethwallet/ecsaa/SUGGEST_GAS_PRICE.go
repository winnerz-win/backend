package ecsaa

import (
	"context"
	"math/big"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"

	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	MUL_VAL = 1.3
)

type IGAS_SENDER interface {
	GET_ETH_CLIENT() *ethclient.Client
}

func GAS_ADD_MUL(gas interface{}) *big.Int {
	return jmath.BigInt(
		jmath.MUL(gas, MUL_VAL),
	)
}

func SUGGEST_GAS_PRICE_2(sender IGAS_SENDER) (*big.Int, error) {
	v, e := sender.GET_ETH_CLIENT().SuggestGasPrice(context.Background())
	if e != nil {
		dbg.Red("ecsx.SUGGEST_GAS_PRICE_2 :", e)
		return v, e
	}
	v = GAS_ADD_MUL(v)
	return v, e
}

func SUGGEST_GAS_PRICE(sender IGAS_SENDER) *big.Int {
	v, _ := SUGGEST_GAS_PRICE_2(sender)
	return v
}

///////////////////////////////////////////////////////////////////

func SUGGEST_TIP_PRICE_2(sender IGAS_SENDER) (*big.Int, error) {
	v, e := sender.GET_ETH_CLIENT().SuggestGasTipCap(context.Background())
	if e != nil {
		dbg.Red("ecsx.SUGGEST_TIP_PRICE :", e)
		return v, e
	}
	v = GAS_ADD_MUL(v)
	return v, e
}

func SUGGEST_TIP_PRICE(sender IGAS_SENDER) *big.Int {
	v, _ := SUGGEST_TIP_PRICE_2(sender)

	return v
}
