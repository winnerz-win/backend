package ebcmx

import (
	"errors"
	"math/big"
	"txscheduler/brix/tools/cloudx/ebcmx/abix"
	"txscheduler/brix/tools/dbg"
)

/*

	Etheruem base Common Manager



*/

var (
	//TxCancelUserFee : tx_cancel_user_fee
	TxCancelUserFee = errors.New("tx_cancel_user_fee")
)

func (my *Sender) SetFuncPtr(
	ebcm_BlockByNumber func(number string) *BlockByNumberData,
	ebcm_TransactionByHash func(hashString string) (TransactionBlock, bool, error),
	ebcm_TSender func(contractAddress string) TSender,
) {
	my.ebcm_BlockByNumber = ebcm_BlockByNumber
	my.ebcm_TransactionByHash = ebcm_TransactionByHash
	my.ebcm_TSender = ebcm_TSender
}
func (my *Sender) BlockByNumber(number string) *BlockByNumberData {
	data := my.ebcm_BlockByNumber(number)
	if data != nil {
		for i := range data.TxList {
			data.TxList[i].Finder = my
		}
	}
	return data
}
func (my *Sender) TransactionByHash(hashTring string) (TransactionBlock, bool, error) {
	tx, pending, err := my.ebcm_TransactionByHash(hashTring)
	if err != nil {
		dbg.Red("ebcmx.TransactionByHash :", err, "[", hashTring, "]")
	}

	tx.Finder = my
	return tx, pending, err
}

type DataItemListParser struct {
	fn func(DataItemList, abix.TypeList, func(abix.RESULT)) bool
}

func (my DataItemListParser) ParseABI(list DataItemList, typelist abix.TypeList, callback func(rs abix.RESULT)) bool {
	return my.fn(list, typelist, callback)
}

func NewDataItemListParser(fn func(DataItemList, abix.TypeList, func(abix.RESULT)) bool) DataItemListParser {
	return DataItemListParser{
		fn: fn,
	}
}

const (
	MESSAGE_PREFIX_METAMASK = "\u0019Ethereum Signed Message:\n"
	MESSAGE_PREFIX_KLAYTN   = "\u0019Klaytn Signed Message:\n"
)

type GetSignTooler func(message_prefix string) SignTool

func (my *Sender) ToggleHost() {
	if my.toggleHostFunc == nil {
		return
	}
	my.toggleHostFunc()
}
func (my *Sender) SetToggleHost(f func()) {
	my.toggleHostFunc = f
}

type Sender struct {
	ISender interface{}
	abix.Caller
	SignTool   *SignTool
	SignTooler GetSignTooler

	toggleHostFunc func()

	Infomation func() string
	Mainnet    func() bool
	HostURL    func() string
	Key        func() string
	//Client() Client

	//Balance : coin-balance
	Balance      func(hexAddress string) string //Coin Balance
	CoinPrice    func(hexAddress string) string
	TokenBalance func(hexAddress string, contract string) string
	TokenPrice   func(hexAddress string, contract string, decimal interface{}) string
	ChainID      func(isFail ...*bool) string
	//BlockNumber() string
	BlockNumberTry func(defaultNumber string) string

	ebcm_BlockByNumber     func(number string) *BlockByNumberData                  //*Sender
	ebcm_TransactionByHash func(hashString string) (TransactionBlock, bool, error) //*Sender
	ReceiptByHash          func(hexHash string) TxReceipt
	InjectReceipt          func(tx *TransactionBlock, r TxReceipt)

	MakePadBytesABI func(pureName string, ebcmTypes abix.TypeList) PADBYTES
	MakePadBytes    func(method string, callback func(Appender)) PADBYTES
	MakePadBytes2   func(funcName string, types abix.TypeList, callback func(Appender)) PADBYTES

	ebcm_TSender func(contractAddress string) TSender
	Token        func(contractAddress string) Token

	PadBytesETH          func() PADBYTES
	TransferPadBytes     func(toAddress string, wei string) PADBYTES
	PadBytesTransfer     func(toAddress string, tokenWei string) PADBYTES
	PadBytesApprove      func(spender, amount string) PADBYTES
	PadBytesApproveAll   func(spender string) PADBYTES
	PadBytesTransferFrom func(from, to, amount string) PADBYTES

	XNonce   func(hexAddress string) (uint64, error)
	XNonceAt func(hexAddress string) (uint64, error)

	XEstimateFeeETH func(
		fromPrivate string,
		toAddress string,
		ipadBytes PADBYTES,
		wei string,
		speed GasSpeed,
	) string
	XGasLimit func(
		paddedData PADBYTES,
		fromAddress, toAddress string,
		wei string,
	) (uint64, error)

	XPipe func(
		fromPrivate string,
		toAddress string,
		ipadBytes PADBYTES,
		wei string,
		speed GasSpeed,
		limitCB func(gasLimit uint64) uint64,
		nonceCB func(nonce uint64) uint64,
		gaspCB func(gasPrice XGasPrice) XGasPrice,
		resultCB func(r XSendResult),
	) error

	XPipeFixedGAS func(
		fromPrivate string,
		toAddress string,
		ipadBytes PADBYTES,
		wei string,
		limitCB func(gasLimit uint64) uint64,
		nonceCB func(nonce uint64) uint64,
		gasPair []*big.Int,
		txFeeWeiAllow func(feeWEI string) bool,
		resultCB func(r XSendResult),
	) error

	TransferCoin func(
		fromPrivate string,
		toAddress string,
		wei string,
		speed GasSpeed,
		limitCB func(gasLimit uint64) uint64,
		nonceCB func(nonce uint64) uint64,
		gaspCB func(gasPrice XGasPrice) XGasPrice,
		resultCB func(r XSendResult),
	) error

	TransferCoinFixedGAS func(
		fromPrivate string,
		toAddress string,
		wei string,
		limitCB func(gasLimit uint64) uint64,
		nonceCB func(nonce uint64) uint64,
		gasPair []*big.Int,
		txFeeWeiAllow func(feeWEI string) bool,
		resultCB func(r XSendResult),
	) error

	DelegateCallContract func(from, to string, data []byte) ([]byte, error)

	///////////////////// public //////////////////////////
	IsAddress            func(address string) bool
	ContractAddressNonce func(from string, nonce uint64) string

	NewSeedI  func(text string, seq interface{}) IWallet
	GetWallet func(hexPrivate string) (IWallet, error)

	DataItemList_ParseABI func(DataItemList, abix.TypeList, func(abix.RESULT)) bool

	GetGasResult func() GasResult
}

// func (my Sender) String() string { return my.Infomation() }
func (my Sender) CallContract(from, to string, data []byte) ([]byte, error) {
	return my.DelegateCallContract(from, to, data)
}

func (my *Sender) TSender(contractAddress string) TSender {
	ts := my.ebcm_TSender(contractAddress)
	ts.Sender = my
	return ts
}

// IWallet :
type IWallet interface {
	Index() interface{}
	String() string
	PrivateKey() string
	Address() string
	CompareAddress(cmpAddress string) bool
}

type TSender struct {
	*Sender
	ContractAddress func() string

	Allowance    func(owner, spender string) string
	Approve      func(privateKey string, spender string, amount string) string
	ApproveAll   func(privateKey string, spender string) string
	TransferFrom func(privateKey string, owner, to string, amount string) string

	TransferFunction func(
		fromPrivate string,
		ipadBytes PADBYTES,
		wei string,
		speed GasSpeed,
	) (
		string, //hash
		uint64, //nonce-value
		error,
	)
	TransferFunctionFixedGAS func(
		fromPrivate string,
		ipadBytes PADBYTES,
		wei string,
		gasPair []*big.Int,
		txFeeWeiAllow func(feeWEI string) bool,
	) (
		string, //hash
		uint64, //nonce-value
		error,
	)

	TransferToken func(
		privateKey string,
		to string,
		tokenWEI string,
		speed GasSpeed,
		limitCB func(gasLimit uint64) uint64,
		nonceCB func(nonce uint64) uint64,
		gaspCB func(gasPrice XGasPrice) XGasPrice,
		resultCB func(r XSendResult),
	) error

	TransferTokenFixedGAS func(
		privateKey string,
		to string,
		tokenWEI string,
		limitCB func(gasLimit uint64) uint64,
		nonceCB func(nonce uint64) uint64,
		gasPair []*big.Int,
		txFeeWeiAllow func(feeWEI string) bool,
		resultCB func(r XSendResult),
	) error

	Write func(
		privateKey string,
		padBytes PADBYTES,
		wei string,
		speed GasSpeed,
		limitCB func(gasLimit uint64) uint64,
		nonceCB func(nonce uint64) uint64,
		gaspCB func(gasPrice XGasPrice) XGasPrice,
		resultCB func(r XSendResult),
	) error

	WriteFixedGAS func(
		privateKey string,
		padBytes PADBYTES,
		wei string,
		limitCB func(gasLimit uint64) uint64,
		nonceCB func(nonce uint64) uint64,
		gasPair []*big.Int,
		txFeeWeiAllow func(feeWEI string) bool,
		resultCB func(r XSendResult),
	) error
}

type Token interface {
	Valid() bool
	String() string

	Name(...bool) string
	Symbol() string
	Decimals(...bool) string
	TotalSupply() string

	Address() string
	Balance(hexAddress string) string
}
