package inf

import (
	"jcloudnet/itype"
	"jtools/cloud/ebcm"
	"jtools/cloud/ebcm/abi"
	"jtools/cloud/jeth/ecs"
	"jtools/cloud/jeth/jwallet"
	"jtools/jmath"
	"strings"
	"sync"
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
func Version() string { return config.Version }
func Mainnet() bool   { return config.Mainnet }

func IsOnwerTaskMode() bool { return config.IsLockTransferByOwner }
func Owner() KeyPair        { return config.Owners[0] }

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

// /////////////////////////////////////////////////////////////
func GetFinder() *itype.IClient {
	return itype.New(
		ecs.RPC_URL(Mainnet()),
		false,
		InfuraKey(),
	)
}

func GetSender() *ebcm.Sender {
	s := GetFinder()
	sender := s.EBCMSender(ecs.TxSigner{})
	return sender
}

type ERC20 struct {
	contract string
	symbol   string
	decimals string
}

func (my ERC20) Valid() bool      { return my.contract != "" }
func (my ERC20) Address() string  { return my.contract }
func (my ERC20) Symbol() string   { return my.symbol }
func (my ERC20) Decimals() string { return my.decimals }
func (my ERC20) String() string {
	v := map[string]string{
		"address":  my.contract,
		"symbol":   my.symbol,
		"decimals": my.decimals,
	}
	return dbg.ToJSONString(v)
}

type ICaller interface {
	Call(contract string, method abi.Method, caller string, f func(rs abi.RESULT), isLogs ...bool) error
}

func GetERC20(
	caller ICaller,
	address string,
) ERC20 {
	re := ERC20{}
	if err := caller.Call(
		address,
		abi.Method{
			Name:   "symbol",
			Params: abi.NewParams(),
			Returns: abi.NewParams(
				abi.String,
			),
		},
		address,
		func(rs abi.RESULT) {
			re.symbol = rs.Text(0)
		},
	); err != nil {
		return re
	}
	if err := caller.Call(
		address,
		abi.Method{
			Name:   "decimals",
			Params: abi.NewParams(),
			Returns: abi.NewParams(
				abi.Uint8,
			),
		},
		address,
		func(rs abi.RESULT) {
			re.decimals = jmath.VALUE(rs.Uint8(0))
		},
	); err != nil {
		return re
	}

	re.contract = dbg.TrimToLower(address)
	return re
}

///////////////////////////////////////////////////////////////

// TokenList : all (ETH 포함)
func TokenList() TokenInfoList {
	defer config.mu.RUnlock()
	config.mu.RLock()
	return config.Tokens
}
func FirstERC20() TokenInfo {
	return TokenList().FirstERC20()
}

func AddToken(contract string) bool {

	contract = dbg.TrimToLower(contract)

	finder := GetFinder()
	token := GetERC20(finder, contract)
	if token.Valid() == false {
		return false
	}

	isAdd := true
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
func Wallet(uid int64) jwallet.IWallet {
	if config.Seed == "" {
		panic("model.Wallet -- secureSeed is not SET")
	}

	seed := config.Seed
	if config.Mainnet == false {
		dbg.Green(config.Seed)
	}
	//seed = ""
	return jwallet.NewSeedI(seed, uid)
}

func EtherScanAddressURL() string {
	if config.Mainnet {
		return "https://etherscan.io/address/address/"
	} else {
		return "https://sepolia.etherscan.io/address/"
	}
}

//////////////////////////////////////////////////////////////

var (
	file_log_writer func(v ...interface{}) = nil
)

func SetFileLogWriter(f func(v ...interface{})) {
	file_log_writer = f
}

func LogWrite(v ...interface{}) {
	if file_log_writer == nil {
		return
	}
	file_log_writer(v...)
}
