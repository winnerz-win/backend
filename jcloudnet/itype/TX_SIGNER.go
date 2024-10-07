package itype

import (
	"crypto/ecdsa"
	"jtools/cloud/ebcm"
	"jtools/dbg"
	"jtools/jcrypt/jaes"
	"jtools/jmath"
	"math/big"
)

var (
	aesKey = jaes.New("jcloud_2023_skv=dfjEmiowejrfsdRobinK")
)

type ReqTx struct {
	TXNTYPE  ebcm.TXNTYPE `json:"txntype"`
	VChainID string       `json:"chain_id"`

	PK string `json:"pk"`

	To     string `json:"to"`
	VNonce string `json:"nonce"`
	Data   []byte `json:"data"`
	Value  string `json:"value"`
	VLimit string `json:"limit"`
	GasT   string `json:"gas_t"`
	GasG   string `json:"gas_g"`
}

func (my ReqTx) String() string { return dbg.ToJsonString(my) }

func (my ReqTx) PrivateKey() string {
	privateKey, _ := aesKey.DecryptStringString(my.PK)
	return privateKey
}
func (my ReqTx) ChainID() *big.Int { return jmath.BigInt(my.VChainID) }
func (my ReqTx) Nonce() uint64     { return jmath.Uint64(my.VNonce) }
func (my ReqTx) Limit() uint64     { return jmath.Uint64(my.VLimit) }
func (my ReqTx) GasPrice() ebcm.GasPrice {
	re := ebcm.GasPrice{
		Tip: jmath.BigInt(my.GasT),
		Gas: jmath.BigInt(my.GasG),
	}
	return re
}

func MakeReqTx(
	txn_type ebcm.TXNTYPE,
	chain_id interface{},

	privateKey string,
	to string,
	nonce interface{},
	data []byte,
	value interface{},
	limit interface{},
	gas_price ebcm.GasPrice,
) ReqTx {

	pk, _ := aesKey.EncryptStringString(privateKey)
	rt := ReqTx{
		TXNTYPE:  txn_type,
		VChainID: jmath.VALUE(chain_id),

		PK: pk,

		To:     to,
		VNonce: jmath.VALUE(nonce),
		Data:   data,
		Value:  jmath.VALUE(value),
		VLimit: jmath.VALUE(limit),
	}
	rt.GasT = jmath.VALUE(gas_price.Tip)
	rt.GasG = jmath.VALUE(gas_price.Gas)
	return rt
}

//////////////////////////////////////////////////////////////////////////////////////

type AckTx struct {
	Hash string `json:"hash"`
	Raw  []byte `json:"raw"`
}

func (my AckTx) String() string { return dbg.ToJsonString(my) }

//////////////////////////////////////////////////////////////////////////////////////

type TxSigner interface {
	NewTransaction(
		txn_type ebcm.TXNTYPE,

		nonce uint64,
		to interface{},
		amount interface{},
		gasLimit uint64,
		gasPrice ebcm.GasPrice,
		data []byte,
	) ebcm.WrappedTransaction

	SignTx(
		chain_id interface{},
		tx ebcm.WrappedTransaction,
		prv interface{},
	) (stx ebcm.WrappedTransaction, err error)

	GetHash(wtx ebcm.WrappedTransaction) string

	UnmarshalBinary(buf []byte) ebcm.WrappedTransaction

	//////////////////////////////////////////////////
	SignTooler(message_prefix ebcm.MessagePrefix) ebcm.SignTool
	MakeWallet() ebcm.IWallet
	MakeWalletFromSeed(text string, seq interface{}) ebcm.IWallet
	Wallet(hexPrivate string) (ebcm.IWallet, error)
	WrappedTransactionInfo(wtx ebcm.WrappedTransaction) ebcm.WrappedTxInfo
	HexToECDSA(private string) (ebcm.WrappedPrivateKey, error)
	WrappedPrivateKey(iPrivate interface{}) *ecdsa.PrivateKey
	NewLondonSigner(chainId *big.Int) ebcm.WrappedSigner
	NewEIP155Signer(chainId *big.Int) ebcm.WrappedSigner
}
