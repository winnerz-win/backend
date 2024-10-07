package ecsx

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type tokenInfo struct {
	Address  string `bson:"address" json:"address"`
	Symbol   string `bson:"symbol" json:"symbol"`
	Decimals string `bson:"decimals" json:"decimals"`
}
type tokenInfoList []tokenInfo

// TokenContainer :
type TokenContainer struct {
	symbolMap  map[string]*tokenInfo
	addressMap map[string]*tokenInfo
	mu         sync.RWMutex
}

func (my *TokenContainer) isAdd(token *Token) bool {
	defer my.mu.RUnlock()
	my.mu.RLock()
	_, do := my.addressMap[token.Address()]
	return do
}
func (my *TokenContainer) add(token *Token) error {
	defer my.mu.Unlock()
	my.mu.Lock()
	info := &tokenInfo{
		Address: token.Address(),
	}
	if v, err := token.ins.Symbol(&bind.CallOpts{}); err == nil {
		info.Symbol = v
	} else {
		return err
	}
	if v, err := token.ins.Decimals(&bind.CallOpts{}); err == nil {
		info.Decimals = fmt.Sprintf("%v", v)
	} else {
		return err
	}
	my.symbolMap[info.Symbol] = info
	my.addressMap[info.Address] = info
	return nil
}

var (
	tokenContainer = TokenContainer{
		symbolMap:  map[string]*tokenInfo{},
		addressMap: map[string]*tokenInfo{},
	}
)
