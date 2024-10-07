/*
EhterClient 정리
*/
package ecsx

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"

	"txscheduler/brix/tools/cloudx/ethwallet/EtherClient"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsaa"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

const (
	hostMainnet = "https://mainnet.infura.io/v3"
	hostTestnet = "https://goerli.infura.io/v3"

	etherName1 = "eth"
	etherName2 = "ether"
	etherName3 = "ethereum"
)

func RPC_URL(mainnet bool) string {
	if mainnet {
		return hostMainnet
	}
	return hostTestnet
}

var ETHTAG = []string{
	etherName1,
	etherName2,
	etherName3,
}

func ETHTAGS() []string { return ETHTAG }

type ICore interface {
	CoreClient() *ethclient.Client
}

// Sender :
type Sender struct {
	client    *ethclient.Client
	mainnet   bool
	infuraKey string
	hostURL   string

	host_index int
	host_list  []string

	gasURL    string
	etherTags []string
	//gasFunc   IGasPrice
	minGasWei string // 10gwei -- 10000000000

	cacheChainID string
	chainID      *big.Int

	customMethods MethodIDDataList

	txnType TXNTYPE
}

// IGAS_SENDER : SUGGEST_GAS_PRICE
func (my Sender) GET_ETH_CLIENT() *ethclient.Client {
	_ = ecsaa.SUGGEST_GAS_PRICE
	return my.client
}

func (my Sender) CacheChainID() string { return my.cacheChainID }

func (my Sender) String() string {
	s := map[string]interface{}{
		"mainnet":   my.mainnet,
		"hostURL":   my.hostURL,
		"host_list": my.host_list,
		"key":       my.infuraKey,
		"gasURL":    my.gasURL,
		"etherTags": my.etherTags,
		"minGasWei": my.minGasWei,
	}
	return dbg.ToJSONString(s)
}

func (my *Sender) SetTxnLegacy() *Sender {
	if my != nil {
		my.txnType = TXN_LEGACY
	}
	return my
}
func (my *Sender) SetTxnEIP1559() *Sender {
	if my != nil {
		my.txnType = TXN_EIP_1559
	}
	return my
}

func (my *Sender) SetTxnType(txnType TXNTYPE) {
	if my != nil {
		my.txnType = txnType
	}
}

// CoreClient : ICore
func (my *Sender) CoreClient() *ethclient.Client { return my.client }

func (my *Sender) SetMainnet(v bool) { my.mainnet = v }

func (my Sender) isTagDo(tag string) bool {
	for _, v := range my.etherTags {
		if v == tag {
			return true
		}
	}
	return false
}
func (my Sender) isEtherTag(tokenAddress string) bool {
	for _, tag := range my.etherTags {
		if tokenAddress == tag {
			return true
		}
	}
	return false
}

// Client :
func (my Sender) Client() *ethclient.Client { return my.client }

// NetworkID :
func (my Sender) NetworkID(isFail ...*bool) string {
	id, err := my.client.NetworkID(context.Background())
	if id != nil {
		return id.String()
	}
	if err != nil {
		dbg.Red("NetworkID :", err)
		if len(isFail) > 0 {
			*isFail[0] = true
		}
	}
	return "0"
}
func (my Sender) ChainID(isFail ...*bool) string {
	return my.NetworkID(isFail...)
}

func chainIDCacheView(id string) {
	dbg.YellowBoldBG("★★★★★★★★★★★★★★★★★")
	dbg.YellowBoldBG("NewCACHE->ChainID is", id)
	dbg.YellowBoldBG("★★★★★★★★★★★★★★★★★")
}

func NewCACHE() func(mainnet bool, infuraKey string) *Sender {
	chainID := ""
	delegater := func(mainnet bool, infuraKey string) *Sender {
		sender := New(mainnet, infuraKey, chainID)
		if chainID == "" {
			if sender != nil {
				chainID = sender.CacheChainID()
				chainIDCacheView(chainID)
			}
		} else {
			// if sender != nil {
			// 	dbg.Yellow("cache.chainID :", chainID, sender.cacheChainID)
			// }
		}
		return sender
	}
	return delegater
}
func NewMainnetCACHE() func(host string, infuraKey string, tags []string) *Sender {
	chainID := ""
	delegater := func(host string, infuraKey string, tags []string) *Sender {
		sender := NewMainnet(host, infuraKey, tags)
		if chainID == "" {
			if sender != nil {
				chainID = sender.CacheChainID()
				chainIDCacheView(chainID)
			}
		} else {
			// if sender != nil {
			// 	dbg.Yellow("cache.chainID :", chainID, sender.cacheChainID)
			// }
		}
		return sender
	}
	return delegater
}

// New :
func New(mainnet bool, infuraKey string, cacheChainID ...string) *Sender {
	//"ded78ac6d48643c897c45048dd929df1"	-- younha
	//ff84cf3fca724ac1b461f966d7bafc08 		-- shera10004@brickstream
	host := hostMainnet
	if mainnet == false {
		host = hostTestnet
	}
	sender := newSender(host, mainnet, infuraKey, cacheChainID...)
	if sender == nil {
		return nil
	}
	sender.etherTags = append(sender.etherTags,
		etherName1,
		etherName2,
		etherName3,
	)
	sender.linkGasPriceFunc()
	return sender
}

// NewMainnet : 토모체인 등등 (이더리움 기반 메인넷)
func NewMainnet(host string, infuraKey string, tags []string, cacheChainID ...string) *Sender {
	// if len(tags) == 0 {
	// 	dbg.Red("tags param is zero ( eth , ether , ethereum )")
	// 	return nil
	// }
	sender := newSender(host, true, infuraKey, cacheChainID...)
	if sender == nil {
		return nil
	}
	if len(tags) > 0 {
		ts := []string{}
		for _, v := range tags {
			ts = append(ts, strings.ToLower(v))
		}
		sender.etherTags = append(sender.etherTags, ts...)
	} else {
		sender.etherTags = append(sender.etherTags,
			etherName1,
			etherName2,
			etherName3,
		)
	}

	//sender.linkGasPriceFunc()
	// sender.gasFunc = SuggestGasPrice
	// sender.SetLinkGasPriceFunc(SuggestGasPrice)
	return sender
}

func NewLegacy(host string, infuraKey string, tags []string, cacheChainID ...string) *Sender {
	sender := NewMainnet(host, infuraKey, tags, cacheChainID...)
	if sender != nil {
		sender = sender.SetTxnLegacy()
	}
	return sender
}

func newSender(host string, mainnet bool, infuraKey string, cacheChainID ...string) *Sender {
	hosturl := host
	url := ""
	infuraKey = strings.TrimSpace(infuraKey)
	if infuraKey != "" {
		url = hosturl + "/" + infuraKey
	} else {
		url = hosturl
	}

	client, err := ethclient.Dial(url)
	if err != nil {
		fmt.Println("Sender.New::Dial :", err)
		return nil
	}

	//client.NetworkID(context.Background())

	if len(cacheChainID) > 0 && jmath.CMP(cacheChainID, 0) > 0 {

	} else {

	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		fmt.Println("○○○○○○○○○○○○○○○○○○○○○○○○○○○○○○○")
		dbg.Red("newSender.ChainID :", err)
		dbg.Red("host :", host)
		fmt.Println("○○○○○○○○○○○○○○○○○○○○○○○○○○○○○○○")
		//chainID.SetString("0", 32)
		if len(cacheChainID) > 0 && jmath.IsNum(cacheChainID[0]) {
			if cc, do := big.NewInt(0).SetString(cacheChainID[0], 32); do {
				chainID = cc
			}
		} else {
			return nil
		}
	}

	return &Sender{
		client:    client,
		mainnet:   mainnet,
		infuraKey: infuraKey,
		hostURL:   hosturl,
		gasURL:    hosturl,
		minGasWei: "0",

		cacheChainID: jmath.VALUE(chainID.Int64()),
		chainID:      big.NewInt(chainID.Int64()),
		txnType:      TXN_EIP_1559,
	}
}

// SetKey :
func (my *Sender) SetKey(infuraKey string) {
	my.infuraKey = infuraKey
	my.client.Close()
	my.client, _ = ethclient.Dial(my.hostURL + "/" + infuraKey)

}

// Mainnet :
func (my Sender) Mainnet() bool { return my.mainnet }

// HostURL :
func (my Sender) HostURL() string { return my.hostURL }

// InfuraKey :
func (my Sender) InfuraKey() string { return my.infuraKey }

// Key :
func (my Sender) Key() string { return my.InfuraKey() }

// Balance : ETH - wei
func (my Sender) Balance(hexAddress string) string {
	hexAddress = strings.TrimSpace(hexAddress)
	if hexAddress == "" {
		dbg.Red("hexAddress is empty")
		return "0"
	}
	address := common.HexToAddress(strings.ToLower(hexAddress))
	v, err := my.client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		dbg.Red("Sender.Balance :", err)
		return "0"
	}
	return v.String()
}

// CoinPrice : ETH - value
func (my Sender) CoinPrice(hexAddress string) string {
	wei := my.Balance(hexAddress)
	return WeiToETH(wei)
}

func (my Sender) TokenBalance(hexAddress string, contract string) string {
	wei := my.Balance2(hexAddress, contract)
	return wei
}

// TokenPRice : Token - value
func (my Sender) TokenPrice(hexAddress string, contract string, decimal interface{}) string {
	wei := my.TokenBalance(hexAddress, contract)
	return WeiToToken(wei, jmath.VALUE(decimal))
}

// Balance2 : address , contract (eth , ETH , ether , ETHER , ethereum)
func (my Sender) Balance2(hexAddress string, contract string) string {
	lowerstr := strings.TrimSpace(strings.ToLower(contract))
	wei := "0"

	if my.isEtherTag(lowerstr) {
		wei = my.Balance(hexAddress)
	} else {
		if strings.HasPrefix(contract, "0x") == false {
			fmt.Println("[ecsx.Sender] contract name is not allow :", contract)
			return wei
		}
		token := my.Token(contract)
		if token == nil {
			dbg.Red("token is nil")
			return wei
		}
		wei = token.Balance(hexAddress)
	}
	return wei
}

// GasPriceETH :
func (my Sender) GasPriceETH(fromAddress, toAddress string, wei string, speed GasSpeed, isLog ...bool) (string, error) {
	fromHexAddress := common.HexToAddress(strings.ToLower(fromAddress))
	toHexAddress := common.HexToAddress(strings.ToLower(toAddress))

	value := new(big.Int)
	value.SetString(wei, 10)
	gasLimit, err := my.client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  fromHexAddress,
		To:    &toHexAddress,
		Value: value,
		Data:  []byte{},
	})
	if err != nil {
		return "0", fmt.Errorf("Sender.GasPriceETH : %v", err)
	}
	//gasLimit := uint64(21000) // in units

	limit := big.NewInt(int64(gasLimit))
	gasPrice := NewGasStation().Price(speed.Value())

	if len(isLog) > 0 && isLog[0] {
		fmt.Println("limit    :", gasLimit, "(", WeiToETH(fmt.Sprintf("%v", gasLimit)), ")")
		fmt.Println("gasPrice :", gasPrice, "(", WeiToETH(gasPrice.String()), ")")
	}

	fee := gasPrice.Mul(gasPrice, limit)
	return fee.String(), nil
}

// GasPriceToken : [ tokenAddress == "eth", "ether", "ethereum" -> GasPriceETH ]
func (my Sender) GasPriceToken(tokenAddress, fromAddress, toAddress string, value string, speed GasSpeed, isLog ...bool) (string, error) {
	ethName := strings.ToLower(tokenAddress)

	if my.isEtherTag(ethName) {
		return my.GasPriceETH(fromAddress, toAddress, value, speed, isLog...)
	}

	tokenHexAddress := common.HexToAddress(strings.ToLower(tokenAddress))
	fromHexAddress := common.HexToAddress(strings.ToLower(fromAddress))
	toHexAddress := common.HexToAddress(strings.ToLower(toAddress))

	gasPrice := NewGasStation().Price(speed.Value())

	paddedAddress := common.LeftPadBytes(toHexAddress.Bytes(), 32)

	amount := new(big.Int)
	amount.SetString(value, 10)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gasLimit, err := my.client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: fromHexAddress,
		To:   &tokenHexAddress,
		//Value: amount,	//이더value만 들어간다.
		Data: data,
	})
	if err != nil {
		return "0", fmt.Errorf("Sender.GasPriceToken : %v", err)
	}
	gasLimit += EtherClient.TXLimitTokenAddValue

	if len(isLog) > 0 && isLog[0] {
		fmt.Println("limit    :", gasLimit, "(", WeiToToken(fmt.Sprintf("%v", gasLimit), 18), ")")
		fmt.Println("gasPrice :", gasPrice, "(", WeiToToken(gasPrice.String(), 18), ")")
	}

	limit := big.NewInt(int64(gasLimit))
	fee := gasPrice.Mul(gasPrice, limit)
	return fee.String(), nil
}

// Nonce :
func (my *Sender) Nonce(privateKeyString string) (*NonceData, error) {
	privatekey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		return nil, fmt.Errorf("Nonce.HexToECDSA : %v", err)
	}
	publicKey := privatekey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("error casting public key to ECDSA")
	}
	fromHexAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := my.client.NonceAt(context.Background(), fromHexAddress, nil)
	if err != nil {
		return nil, fmt.Errorf("Nonce.NonceAt : %v", err)
	}

	chainID, err := my.client.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Nonce.NetworkID : %v", err)
	}

	nd := &NonceData{
		Sender:           my,
		client:           my.client,
		privateKeyString: privateKeyString,
		privatekey:       privatekey,
		fromHexAddress:   fromHexAddress,
		nonceValue:       nonce,
		chainID:          chainID,
	}
	for _, tag := range my.etherTags {
		nd.etherTags = append(nd.etherTags, tag)
	}
	return nd, nil
}

// TxByHash : Convert to transaction data from Hash
func (my Sender) TxByHash(hexHash string) *STX {
	hash := common.HexToHash(hexHash)

	//ETHKJS
	from := ""
	blockNumber := ""
	dbg.D(from, blockNumber)
	//tx, isPending, err := my.client.TransactionByHashJ(context.Background(), hash, &from, &blockNumber)
	tx, isPending, err := my.client.TransactionByHash(context.Background(), hash)
	if err != nil {
		return nil
	}

	if blockNumber != "" {
		blockNumber = jmath.NewBigDecimal(blockNumber).ToString()
	}
	stx := &STX{
		client:      my.client,
		tx:          tx,
		isSend:      true,
		isPending:   isPending,
		from:        from,
		blockNumber: blockNumber,
	}
	return stx
}

// Receipt :
func (my Sender) Receipt(hexHash string) *Receipt {
	hexHash = strings.TrimSpace(hexHash)
	hash := common.HexToHash(hexHash)
	receipt, err := my.client.TransactionReceipt(context.Background(), hash)

	//b, _ := receipt.MarshalJSON()
	// dbg.Yellow("receipt :", string(b))
	// dbg.Yellow("receipt.Status :", receipt.Status)

	if err != nil {
		return &Receipt{
			client: my.client,
			data:   nil,
		}
	}
	return &Receipt{
		client:           my.client,
		data:             receipt,
		TransactionIndex: receipt.TransactionIndex,
	}
}

func (my Sender) CallContract(from, to string, data []byte) ([]byte, error) {
	fromHexAddress := common.HexToAddress(strings.ToLower(from))
	toHexAddress := common.HexToAddress(strings.ToLower(to))
	param := ethereum.CallMsg{
		From: fromHexAddress,
		To:   &toHexAddress,
		Data: data,
	}
	// blockNumber := big.NewInt(1)
	return my.client.CallContract(context.Background(), param, nil)
}

func (my Sender) CallBy(callerAddress string) error {

	return nil
}
