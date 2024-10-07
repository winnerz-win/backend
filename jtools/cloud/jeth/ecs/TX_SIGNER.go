package ecs

import (
	"crypto/ecdsa"
	"jtools/cc"
	"jtools/cloud/ebcm"
	"jtools/cloud/jeth/jwallet"
	"jtools/dbg"
	"jtools/jmath"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
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
	value := jmath.BigInt(amount)
	if txn_type == ebcm.TXN_EIP_1559 {
		toAddress := wrappedAddress(to)

		return WrappedTransaction(
			types.NewTx(
				&types.DynamicFeeTx{
					ChainID:   nil,
					Nonce:     nonce,
					GasTipCap: gasPrice.Tip,
					GasFeeCap: gasPrice.Gas,
					Gas:       gasLimit,
					To:        &toAddress,
					Value:     value,
					Data:      data,
				},
			),
		)
	}
	return WrappedTransaction(
		types.NewTransaction(
			nonce,
			wrappedAddress(to),
			value,
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

	ntx := typesTransaction(tx)
	var singer types.Signer
	if ntx.Type() == ebcm.TXN_EIP_1559.Uint8() {
		singer = types.NewLondonSigner(chainID)

	} else {
		singer = types.NewEIP155Signer(chainID)
	}

	defer func() {
		if e := recover(); e != nil {
			err = dbg.Error(e)
		}
	}()

	_stx, err := types.SignTx(
		typesTransaction(tx),
		singer,
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
		cc.RedItalic("[ECS]", err)
		return nil
	}
	return WrappedTransaction(&tx)
}

//////////////////////////////////////////////////

func (TxSigner) SignTooler(message_prefix ebcm.MessagePrefix) ebcm.SignTool {
	return getSignTool(message_prefix)
}

func (TxSigner) MakeWallet() ebcm.IWallet {
	return jwallet.New()
}

func (TxSigner) MakeWalletFromSeed(text string, seq interface{}) ebcm.IWallet {
	return jwallet.EBCM_NewSeedI(text, seq)
}

func (TxSigner) Wallet(hexPrivate string) (ebcm.IWallet, error) {
	return jwallet.EBCM_Get(hexPrivate)
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
