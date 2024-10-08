package ebcm

import (
	"context"
	"crypto/ecdsa"
	"jtools/cloud/ebcm/abi"
	"math/big"
)

type IClient interface {
	Client() interface{}
	Host() string
	Info() Info
	ChainID() *big.Int
	TXNTYPE() TXNTYPE
	SetTXNTYPE(v TXNTYPE)
	IsDebug() bool
	SetDebug(f ...bool)
	////////////////////////////////////////////////////////////////////////

	NewTransaction(
		nonce uint64,
		to interface{},
		amount interface{},
		gasLimit uint64,
		gasPrice GasPrice,
		data []byte,
	) WrappedTransaction

	SignTx(
		tx WrappedTransaction,
		prv interface{},
	) (WrappedTransaction, error)

	UnmarshalBinary(buf []byte) WrappedTransaction

	SendTransaction(ctx context.Context, tx WrappedTransaction) (string, error)

	CheckSendTxHashReceipt(
		tx WrappedTransaction,
		limitSec int,
		is_debug ...bool,
	) CheckSendTxHashReceiptResult

	CheckSendTxHashToNonce(
		from interface{},
		hash string,
		limitSec int,
		is_debug ...bool,
	) CheckSendTxHashReceiptResult

	NetworkID(ctx context.Context) (*big.Int, error)
	BlockNumber(ctx context.Context) (*big.Int, error)
	BalanceAt(ctx context.Context, account interface{}, blockNumber *big.Int) (*big.Int, error)

	SuggestGasPrice(ctx context.Context, is_skip_tip_cap ...bool) (GasPrice, error)

	EstimateGas(ctx context.Context, msg CallMsg) (uint64, error)
	PendingNonceAt(ctx context.Context, account interface{}) (uint64, error)
	NonceAt(ctx context.Context, account interface{}) (uint64, error)
	CallContract(from, to string, data []byte) ([]byte, error)

	BlockByNumberSimple(number interface{}) *BlockByNumberData
	BlockByNumber(number interface{}) *BlockByNumberData
	TransactionByHash(hashString string) (TransactionBlock, bool, error)
	ReceiptByHash(hexHash string) TxReceipt
	InjectReceipt(tx *TransactionBlock, r TxReceipt)

	HeaderByNumber(number any) *BlockHeader

	Call(
		contract string,
		method abi.Method,
		caller string,
		f func(rs abi.RESULT),
		isLogs ...bool,
	) error

	IClientUtil
}

type CALLER interface {
	Call(
		contract string,
		method abi.Method,
		caller string,
		f func(rs abi.RESULT),
		isLogs ...bool,
	) error
}

// IClientUtil : Client Util Func
type IClientUtil interface {
	SignTooler(message_prefix MessagePrefix) SignTool
	MakeWallet() IWallet
	MakeWalletFromSeed(text string, seq interface{}) IWallet
	Wallet(hexPrivate string) (IWallet, error)

	GetHash(wtx WrappedTransaction) string
	WrappedTransactionInfo(wtx WrappedTransaction) WrappedTxInfo

	BytesToAddressHex(data32 []byte) string
	HexToAddress(address string) WrappedAddress
	HexToECDSA(private string) (WrappedPrivateKey, error)

	WrappedPrivateKey(iPrivate interface{}) *ecdsa.PrivateKey

	NewLondonSigner(chainId *big.Int) WrappedSigner
	NewEIP155Signer(chainId *big.Int) WrappedSigner

	ContractAddressNonce(from string, nonce uint64) string
}
