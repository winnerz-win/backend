package inf

import (
	"strings"
	"sync"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx/jwalletx"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
)

var (
	infuraMu sync.Mutex
	index1   int
	esMu     sync.Mutex
	index2   int
)

func InfuraKey() string {
	defer infuraMu.Unlock()
	infuraMu.Lock()
	key := config.InfuraKeys[index1]
	index1++
	if index1 >= len(config.InfuraKeys) {
		index1 = 0
	}
	return key
}
func ESKey() string {
	defer esMu.Unlock()
	esMu.Lock()
	key := config.ESKeys[index2]
	index2++
	if index2 >= len(config.ESKeys) {
		index2 = 0
	}
	return key
}
func Version() string      { return config.Version }
func Mainnet() bool        { return config.Mainnet }
func Master() KeyPair      { return config.Masters[0] }
func Charger() KeyPair     { return config.Chargers[0] }
func ClientHostIP() string { return config.ClientHost[config.Mainnet][0] }
func ClientAddress() string {
	ss := config.ClientHost[config.Mainnet]
	if ss[1] == "" {
		return ss[0]
	}
	return ss[0] + ":" + ss[1]
}
func SeedView() string {
	sv := "test_network_seed_key"
	if strings.HasPrefix(seed, "mainnet_") {
		sv = "main_network_seed_key"
	}
	return sv
}

func ValidSymbol(symbol string) bool {
	defer config.mu.RUnlock()
	config.mu.RLock()

	for _, token := range config.Tokens {
		if token.Symbol == symbol {
			return true
		}
	}
	return false
}
func SymbolList() []string {
	defer config.mu.RUnlock()
	config.mu.RLock()

	list := []string{}
	for _, token := range config.Tokens {
		list = append(list, token.Symbol)
	}
	return list
}

// TokenList : all (ETH 포함)
func TokenList() TokenInfoList {
	defer config.mu.RUnlock()
	config.mu.RLock()
	return config.Tokens
}

func AddToken(contract string) bool {

	contract = dbg.TrimToLower(contract)

	finder := ecsx.New(config.Mainnet, InfuraKey())
	token := finder.Token(contract)
	if token.Valid() == false {
		return false
	}

	isAdd := true
	_ = token.String()
	DB().Action(config.DB, func(db mongo.DATABASE) {
		newToken := TokenInfo{
			Mainnet:  config.Mainnet,
			Contract: token.Address(),
			Symbol:   token.Symbol(),
			Decimal:  token.Decimals(),
		}
		if newToken.insertDB(db) == nil {
			config.mu.Lock()
			for _, token := range config.Tokens {
				if token.Contract == newToken.Contract {
					isAdd = false
					break
				}
			}
			if isAdd {
				config.Tokens = append(config.Tokens, newToken)
			}
			config.mu.Unlock()
		}
	})

	return isAdd
}

func Confirms() int { return config.Confirms }

// Wallet :
func Wallet(uid int64) jwalletx.IWallet {
	if config.Seed == "" {
		panic("model.Wallet -- secureSeed is not SET")
	}

	seed := config.Seed
	if config.Mainnet == false {
		dbg.Green(config.Seed)
	}
	//seed = ""
	return jwalletx.NewSeedI(seed, uid)
}

func EtherScanAddressURL() string {
	if config.Mainnet {
		return "https://etherscan.io/address/address/"
	} else {
		return "https://goerli.etherscan.io/address/"
	}
}
