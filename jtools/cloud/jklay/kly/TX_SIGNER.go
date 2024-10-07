package kly

import (
	"crypto/ecdsa"
	"jtools/cc"
	"jtools/cloud/ebcm"
	"jtools/cloud/jklay/kwallet"
	"jtools/dbg"
	"jtools/jmath"
	"math/big"

	"github.com/klaytn/klaytn/blockchain/types"
)

type TxSigner struct{}

func (TxSigner) NewTransaction(
	txn_type ebcm.TXNTYPE,

	nonce uint64,
	to interface{},
	amount interface{},
	gasLimit uint64,
	gasPrice ebcm.GasPrice,
	data []byte,
) ebcm.WrappedTransaction {
	return WrappedTransaction(
		types.NewTransaction(
			nonce,
			wrappedAddress(to),
			jmath.BigInt(amount),
			gasLimit,
			gasPrice.Gas,
			data,
		),
	)
}

func (TxSigner) SignTx(
	chain_id interface{},

	tx ebcm.WrappedTransaction,
	prv interface{},
) (stx ebcm.WrappedTransaction, err error) {
	chainID := jmath.BigInt(chain_id)

	defer func() {
		if e := recover(); e != nil {
			err = dbg.Error(e)
		}
	}()

	_stx, err := types.SignTx(
		typesTransaction(tx),
		types.NewEIP155Signer(chainID),
		ClientUtil{}.WrappedPrivateKey(prv),
	)
	return WrappedTransaction(_stx), err
}

func (TxSigner) GetHash(wtx ebcm.WrappedTransaction) string {
	return ClientUtil{}.GetHash(wtx)
}

func (TxSigner) UnmarshalBinary(buf []byte) ebcm.WrappedTransaction {
	tx := types.Transaction{}
	if err := tx.UnmarshalBinary(buf); err != nil {
		cc.RedItalic("[KLY]", err)
		return nil
	}
	return WrappedTransaction(&tx)
}

//////////////////////////////////////////////////

func (TxSigner) SignTooler(message_prefix ebcm.MessagePrefix) ebcm.SignTool {
	return getSignTool(message_prefix)
}

func (TxSigner) MakeWallet() ebcm.IWallet {
	return kwallet.New()
}

func (TxSigner) MakeWalletFromSeed(text string, seq interface{}) ebcm.IWallet {
	return kwallet.EBCM_NewSeedI(text, seq)
}

func (TxSigner) Wallet(hexPrivate string) (ebcm.IWallet, error) {
	return kwallet.EBCM_Get(hexPrivate)
}

func (TxSigner) WrappedTransactionInfo(wtx ebcm.WrappedTransaction) ebcm.WrappedTxInfo {
	return ClientUtil{}.WrappedTransactionInfo(wtx)
}

func (TxSigner) HexToECDSA(private string) (ebcm.WrappedPrivateKey, error) {
	return ClientUtil{}.HexToECDSA(private)
}

func (TxSigner) WrappedPrivateKey(iPrivate interface{}) *ecdsa.PrivateKey {
	return ClientUtil{}.WrappedPrivateKey(iPrivate)
}

func (TxSigner) NewLondonSigner(chainId *big.Int) ebcm.WrappedSigner {
	return ClientUtil{}.NewLondonSigner(chainId)
}

func (TxSigner) NewEIP155Signer(chainId *big.Int) ebcm.WrappedSigner {
	return ClientUtil{}.NewEIP155Signer(chainId)
}
