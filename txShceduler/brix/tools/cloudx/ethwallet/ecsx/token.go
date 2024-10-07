package ecsx

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"txscheduler/brix/tools/cloudx/ebcmx"
	token "txscheduler/brix/tools/cloudx/ethwallet/EtherClient/ethdev/contracts_erc20"
	"txscheduler/brix/tools/cloudx/ethwallet/abmx"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"

	"github.com/ethereum/go-ethereum/common"
)

const nan = "-"

// Token :
type Token struct {
	ins     *token.Token
	address string
	mainnet bool

	name        string
	symbol      string
	decimal     string
	totalsupply string

	sender       *Sender
	callContract func(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
}

func (my *Token) String() string {
	if my.name == nan {
		my.Name()
		my.Symbol()
		my.Decimals()
		my.TotalSupply()
	}
	mapper := map[string]interface{}{
		"address":     my.address,
		"mainnet":     my.mainnet,
		"name":        my.name,
		"symbol":      my.symbol,
		"decimal":     my.decimal,
		"totalsupply": my.totalsupply,
	}
	return dbg.ToJSONString(mapper)
}

func (my *Token) Valid() bool {
	if my.ins == nil {
		return false
	}
	_ = my.String()
	return jmath.IsNum(my.decimal)
}

// Token :
func (my Sender) Token(contractHexAddress string) *Token {
	contractHexAddress = strings.TrimSpace(strings.ToLower(contractHexAddress))
	tokenAddress := common.HexToAddress(contractHexAddress)
	tc, err := token.NewToken(tokenAddress, my.client)
	if err != nil {
		fmt.Println("Sender.Token :", err)
		return nil
	}
	token := &Token{
		sender:      &my,
		ins:         tc,
		address:     contractHexAddress,
		mainnet:     my.mainnet,
		name:        nan,
		symbol:      nan,
		decimal:     nan,
		totalsupply: nan,
	}
	//return my.ec.GetToken(contractHexAddress, isLog...)
	return token
}

func (my Sender) ebcm_Token(contractAddress string) ebcmx.Token {
	return my.Token(contractAddress)
}

func (my Sender) TokenABI(contractHexAddress string, abiString string) *Token {
	contractHexAddress = strings.TrimSpace(strings.ToLower(contractHexAddress))
	tokenAddress := common.HexToAddress(contractHexAddress)
	tc, err := token.NewContractABI(tokenAddress, my.client, abiString)
	if err != nil {
		fmt.Println("Sender.TokenABI :", err)
		return nil
	}
	token := &Token{
		ins:          tc,
		address:      contractHexAddress,
		mainnet:      my.mainnet,
		name:         nan,
		symbol:       nan,
		decimal:      nan,
		totalsupply:  nan,
		callContract: my.client.CallContract,
	}
	//return my.ec.GetToken(contractHexAddress, isLog...)
	return token
}

// Address :
func (my Token) Address() string { return my.address }

// Symbol :
func (my *Token) Symbol() string {
	if strings.HasPrefix(my.symbol, nan) {
		if v, err := my.ins.Symbol(&bind.CallOpts{}); err != nil {
			my.symbol = fmt.Errorf("%v %v", nan, err).Error()
		} else {
			my.symbol = v
		}
	}
	return my.symbol
}

// Decimals :
func (my *Token) Decimals(failThenDefault ...bool) string {
	if strings.HasPrefix(my.decimal, nan) {
		if v, err := my.ins.Decimals(&bind.CallOpts{}); err != nil {
			my.decimal = fmt.Errorf("%v %v", nan, err).Error()
			if len(failThenDefault) > 0 && failThenDefault[0] {
				return "18"
			}
			return "0"
		} else {
			my.decimal = fmt.Sprintf("%v", v)
		}
	}
	return my.decimal
}

// Name :
func (my *Token) Name(failThenDefault ...bool) string {
	if strings.HasPrefix(my.name, nan) {
		if v, err := my.ins.Name(&bind.CallOpts{}); err != nil {
			my.name = ""
		} else {
			my.name = v
		}
	}
	return my.name
}

// TotalSupply :
func (my *Token) TotalSupply() string {
	if strings.HasPrefix(my.totalsupply, nan) {
		if v, err := my.ins.TotalSupply(&bind.CallOpts{}); err != nil {
			my.totalsupply = fmt.Errorf("%v %v", nan, err).Error()
			return "0"
		} else {
			my.totalsupply = v.String()
		}
	}
	return my.totalsupply
}

// Balance : erc-20 wei
func (my Token) Balance(hexAddress string) string {
	address := common.HexToAddress(hexAddress)
	bal, err := my.ins.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		fmt.Println("Token.Balance :", err)
		return "0"
	}
	return bal.String()
}

// func (my Token) CallbyABI(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
// 	return my.ins.TokenCaller.KJSCall(opts, result, method, params...)
// }

func (my Token) Call(method abmx.Method, callerAddress string, f func(rs abmx.RESULT)) error {
	return abmx.Call(
		my.sender,
		my.address,
		method,
		callerAddress,
		func(rs abmx.RESULT) {
			f(rs)
		},
	)
}

func (my Token) Action(privateKey string, padBytes PadBytes, wei string, speed GasSpeed) string {
	hash := ""
	err := my.sender.XPipe(
		privateKey,
		my.address,
		padBytes,
		wei,
		speed,
		nil, nil, nil,
		func(r XSendResult) {
			hash = r.Hash
		},
	)
	if err != nil {
		dbg.Red(err)
		hash = ""
	}
	return hash
}
