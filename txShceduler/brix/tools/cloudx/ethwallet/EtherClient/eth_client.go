package EtherClient

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strings"
	"sync"

	"txscheduler/brix/tools/cloudx/ethwallet/EtherScanAPI"

	token "txscheduler/brix/tools/cloudx/ethwallet/EtherClient/ethdev/contracts_erc20"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// EClient :
type EClient struct {
	*ethclient.Client
	mainnet bool
	//mu      sync.Mutex
}

var (
	gMainClient *EClient
	gTestClient *EClient
	gMu         sync.Mutex

	gMainTokens = map[string]*EToken{}
	gTestTokens = map[string]*EToken{}
	gTu         sync.RWMutex
)

func gAddToken(mainnet bool, contract string, token *EToken) {
	defer gTu.Unlock()
	gTu.Lock()
	if mainnet {
		gMainTokens[contract] = token
	} else {
		gTestTokens[contract] = token
	}

}
func gGetToken(mainnet bool, contract string) *EToken {
	defer gTu.RUnlock()
	gTu.RLock()

	if mainnet {
		return gMainTokens[contract]
	} else {
		return gTestTokens[contract]
	}
}

// NewEClient :
func NewEClient(mainnet bool, forceWait ...bool) *EClient {
	defer gMu.Unlock()
	gMu.Lock()

	if mainnet == true {
		if gMainClient != nil {
			return gMainClient
		}
	} else {
		if gTestClient != nil {
			return gTestClient
		}
	}

	if ins := NewClient(mainnet, forceWait...); ins != nil {
		client := &EClient{
			Client:  ins,
			mainnet: mainnet,
		}
		if mainnet == true {
			gMainClient = client
		} else {
			gTestClient = client
		}
		return client
	}
	return nil
}

// Mainnet :
func (my EClient) Mainnet() bool {
	return my.mainnet
}

// GetToken : ERC-20
func (my *EClient) GetToken(contractHexAddress string, isLog ...bool) *EToken {
	defer gMu.Unlock()
	gMu.Lock()

	_log := false
	if len(isLog) > 0 && isLog[0] {
		_log = true
	}

	contractHexAddress = strings.TrimSpace(strings.ToLower(contractHexAddress))
	if cachedToken := gGetToken(my.mainnet, contractHexAddress); cachedToken != nil {
		if _log {
			fmt.Println("*** EClient.GetToken[CACHED] ***")
			fmt.Println(cachedToken.ToString())
		}
		return cachedToken
	} else if _log {
		fmt.Println("*** EClient.GetToken ***")
	}

	tokenAddress := common.HexToAddress(contractHexAddress)
	instance, err := token.NewToken(tokenAddress, my.Client)
	if err != nil {
		if _log {
			fmt.Println("token.NewToken :", err)
		}
		return nil
	}
	etoken := &EToken{
		Token:   instance,
		address: contractHexAddress,
		mainnet: my.mainnet,
		apiCfg:  nil,
	}
	fieldErr := etoken.setField()
	if fieldErr == nil {
		gAddToken(my.mainnet, contractHexAddress, etoken)
	}

	if _log {
		fmt.Println(etoken.ToString())
	}

	return etoken
}

// GetContractABI :
func (my *EClient) GetContractABI(contractHexAddress, abiString string, isLog ...bool) *EToken {
	defer gMu.Unlock()
	gMu.Lock()

	_log := false
	if len(isLog) > 0 && isLog[0] {
		_log = true
	}

	contractHexAddress = strings.TrimSpace(strings.ToLower(contractHexAddress))
	if cachedToken := gGetToken(my.mainnet, contractHexAddress); cachedToken != nil {
		if _log {
			fmt.Println("*** EClient.GetContractABI[CACHED] ***")
			fmt.Println(cachedToken.ToString())
		}
		return cachedToken
	} else if _log {
		fmt.Println("*** EClient.GetContractABI ***")
	}

	tokenAddress := common.HexToAddress(contractHexAddress)

	instance, err := token.NewContractABI(tokenAddress, my.Client, abiString)
	if err != nil {
		if _log {
			fmt.Println("token.NewContractABI :", err)
		}
		return nil
	}
	etoken := &EToken{
		Token:   instance,
		address: contractHexAddress,
		mainnet: my.mainnet,
		apiCfg:  nil,
	}
	etoken.isSet = true

	if _log {
		fmt.Println(etoken.ToString())
	}

	return etoken
}

func (my *EClient) _PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	defer gMu.Unlock()
	gMu.Lock()
	return my.Client.PendingNonceAt(ctx, account)
}

func (my *EClient) _CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	defer gMu.Unlock()
	gMu.Lock()

	return my.Client.CallContract(ctx, msg, blockNumber)
}

func (my *EClient) _SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	defer gMu.Unlock()
	gMu.Lock()
	return my.Client.SuggestGasPrice(ctx)
}

func (my *EClient) _EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	defer gMu.Unlock()
	gMu.Lock()
	return my.Client.EstimateGas(ctx, msg)
}

func (my *EClient) _NetworkID(ctx context.Context) (*big.Int, error) {
	defer gMu.Unlock()
	gMu.Lock()
	return my.Client.NetworkID(ctx)
}

func (my *EClient) _SendTransaction(ctx context.Context, tx *types.Transaction) error {
	defer gMu.Unlock()
	gMu.Lock()
	return my.Client.SendTransaction(ctx, tx)
}

//////////////////////////////////////////////////////////////////////////////////////////////////

// EToken :
type EToken struct {
	*token.Token
	address string //contract-address

	isSet       bool
	name        string
	symbol      string
	decimals    uint8
	totalSupply *big.Int

	mu      sync.Mutex
	mainnet bool
	apiCfg  EtherScanAPI.Config
}

// GetTokenCallerRaw :
func (my *EToken) GetTokenCallerRaw() token.TokenCallerRaw {
	raw := token.TokenCallerRaw{
		Contract: &my.TokenCaller,
	}
	return raw
}

func (my *EToken) setField() error {
	defer my.mu.Unlock()
	my.mu.Lock()

	if my.isSet == true {
		return nil
	}

	var _err error = nil
	_errIdx := ""
	instance := my.Token
	if name, err := instance.Name(&bind.CallOpts{}); err == nil {
		my.name = name
	} else {
		_err = err
		_errIdx = "Name"
	}
	if symbol, err := instance.Symbol(&bind.CallOpts{}); err == nil {
		my.symbol = symbol
	} else {
		_err = err
		_errIdx = "Symbol"
	}
	if decimals, err := instance.Decimals(&bind.CallOpts{}); err == nil {
		my.decimals = decimals
	} else {
		_err = err
		_errIdx = "Decimals"
	}
	if total, err := instance.TotalSupply(&bind.CallOpts{}); err == nil {
		my.totalSupply = total
	} else {
		my.totalSupply = big.NewInt(0)
		_err = err
		_errIdx = "TotalSupply"
	}

	if _err != nil {
		fmt.Println("EToken.setField("+_errIdx+") :", _err)
	} else {
		my.isSet = true
	}

	return _err

}

// GetBalance :
func (my *EToken) GetBalance(userAddr string, isAPI ...bool) *big.Int {
	defer my.mu.Unlock()
	my.mu.Lock()

	if len(isAPI) > 0 && isAPI[0] == true {
		if my.apiCfg == nil {
			my.apiCfg = EtherScanAPI.NewConfig(my.mainnet, strings.ReplaceAll(my.address, "0x", "ky"))
		}
		v := EtherScanAPI.TokenBalance(my.apiCfg, my.address, userAddr)
		vi := big.NewInt(0)
		re, _ := vi.SetString(v.ToString(), 10)
		return re
	}

	address := common.HexToAddress(userAddr)
	bal, err := my.Token.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		fmt.Println("EToken.GetBalance :", err)
	}
	return bal
}

// GetBalanceString : process-Decimals
func (my *EToken) GetBalanceString(userAddr string, isAPI ...bool) string {
	val := my.GetBalance(userAddr, isAPI...)
	if val == nil {
		return "0"
	}

	fVal := new(big.Float)
	fVal.SetString(val.String())

	dot := big.NewFloat(math.Pow10(my.Decimals()))
	tkVal := new(big.Float).Quo(fVal, dot)
	return tkVal.String()
}

// Address :
func (my *EToken) Address() string {
	return my.address
}

// Name :
func (my *EToken) Name() string {
	my.setField()
	return my.name
}

// Symbol :
func (my *EToken) Symbol() string {
	my.setField()
	return strings.TrimSpace(my.symbol)
}

// Decimals :
func (my *EToken) Decimals() int {
	my.setField()
	return int(my.decimals)
}

// TotalSupply :
func (my *EToken) TotalSupply() string {
	my.setField()
	return my.totalSupply.String()
}

// ToString :
func (my *EToken) ToString() string {
	my.setField()
	return `{
  address     : ` + my.address + `
  name        : ` + my.name + `
  symbol      : ` + my.symbol + `
  decimals    : ` + fmt.Sprintf("%v", my.decimals) + `
  totalSupply : ` + my.totalSupply.String() + `
}
`
}

// NewClient : EtherClient.NewClient
func NewClient(mainnet bool, forceWait ...bool) *ethclient.Client {
	apikey := "ded78ac6d48643c897c45048dd929df1"
	url := "https://mainnet.infura.io/v3/" + apikey
	if mainnet == false {
		url = "https://goerli.infura.io/v3/" + apikey
	}

	network_id := "1"
	fw := false
	if len(forceWait) > 0 && forceWait[0] == true {
		fw = true
	}

	loopLimit := 1
	var client *ethclient.Client
	var err error
	for {
		client, err = ethclient.Dial(url)
		if err != nil {
			return nil
		}

		var id *big.Int
		id, err = client.NetworkID(context.Background())
		if err == nil {
			fmt.Println("mainnet", mainnet, ", network_id", id)
		} else {
			fmt.Println("NewClient", err)
		}

		if id.String() == network_id {
			break
		}

		if fw == false {
			loopLimit--
			if loopLimit == 0 {
				break
			}
		}
	} //for

	return client
}

// HenaContractAddress :
func HenaContractAddress() string {
	return strings.ToLower("0x8d97C127236D3aEf539171394212F2e43ad701C4")
}

// NewToken :
func NewToken(mainnet bool, contractAddress string, isLog bool) *token.Token {
	if isLog {
		fmt.Println("*** EtherClient.NewToken New ***")
	}

	client := NewClient(mainnet)
	if client == nil {
		return nil
	}

	// HENA Address
	tokenAddress := common.HexToAddress(contractAddress)
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		if isLog {
			fmt.Println("EtherClient.NewToken :", err)
		}
		return nil
	}

	if isLog {
		fmt.Println("contract_address :", contractAddress)
	}
	if name, err := instance.Name(&bind.CallOpts{}); err == nil {
		if isLog {
			fmt.Println("token_name :", name)
		}
	}
	if symbol, err := instance.Symbol(&bind.CallOpts{}); err == nil {
		if isLog {
			fmt.Println("token_symbol :", symbol)
		}
	}
	if decimals, err := instance.Decimals(&bind.CallOpts{}); err == nil {
		if isLog {
			fmt.Println("token_decimals :", decimals)
		}
	}
	if total, err := instance.TotalSupply(&bind.CallOpts{}); err == nil {
		if isLog {
			fmt.Println("token_total :", total)
		}
	}

	if isLog {
		fmt.Println("----------------------------")
	}

	return instance
}

// GetTokenBalance :	EtherScanAPI.TokenBalance()
func GetTokenBalance(tk *token.Token, userAddr string) *big.Int {
	address := common.HexToAddress(userAddr)
	bal, err := tk.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		//fmt.Println("GetHenaBalance :", err)
		return big.NewInt(0)
	}
	return bal
}

// GetTokenBalance2 : EtherScanAPI.TokenBalance()
func GetTokenBalance2(tk *token.Token, userAddr string) string {
	val := GetTokenBalance(tk, userAddr)
	decimal, _ := tk.Decimals(&bind.CallOpts{})

	fVal := new(big.Float)
	fVal.SetString(val.String())
	dot := big.NewFloat(math.Pow10(int(decimal)))
	tkVal := new(big.Float).Quo(fVal, dot)
	return tkVal.String()
}
