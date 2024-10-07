package kly

import (
	"crypto/ecdsa"
	"encoding/hex"
	"jtools/cc"
	"jtools/cloud/ebcm"
	"jtools/cloud/jklay/kwallet"
	"jtools/jmath"
	"math/big"
	"strings"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
)

////////////////////////////////////////////////////////////////////////

type ClientUtil struct{}

func (my ClientUtil) SignTooler(message_prefix ebcm.MessagePrefix) ebcm.SignTool {
	return getSignTool(message_prefix)
}

func (my ClientUtil) MakeWallet() ebcm.IWallet { return kwallet.New() }

func (my ClientUtil) MakeWalletFromSeed(text string, seq interface{}) ebcm.IWallet {
	return kwallet.EBCM_NewSeedI(text, seq)
}

func (my ClientUtil) Wallet(hexPrivate string) (ebcm.IWallet, error) {
	return kwallet.EBCM_Get(hexPrivate)
}

func (my ClientUtil) GetHash(wtx ebcm.WrappedTransaction) string {
	defer func() {
		if e := recover(); e != nil {
			cc.RedItalic("Client.Hash :", e)
		}
	}()
	tx := typesTransaction(wtx)
	return tx.Hash().Hex()
}

func (my ClientUtil) WrappedTransactionInfo(wtx ebcm.WrappedTransaction) ebcm.WrappedTxInfo {
	defer func() {
		if e := recover(); e != nil {
			cc.RedItalic("WrappedTransactionInfo :", e)
		}
	}()
	tx := typesTransaction(wtx)

	info := ebcm.WrappedTxInfo{
		IsSigned:    jmath.VALUE(tx.ChainId()) != "0",
		ChainID:     jmath.VALUE(tx.ChainId()),
		Hash:        tx.Hash().Hex(),
		Cost:        jmath.VALUE(tx.Cost()),
		GasPrice:    jmath.VALUE(tx.GasPrice()),
		GasTipCap:   jmath.VALUE(tx.GasTipCap()),
		GasFeeCap:   jmath.VALUE(tx.GasFeeCap()),
		Nonce:       tx.Nonce(),
		Gas:         tx.Gas(),
		TxnType:     ebcm.TXNTYPE(tx.Type()),
		TxnTypeText: ebcm.TXNTYPE(tx.Type()).String(),
		To:          tx.To().Hex(),
		Protected:   tx.Protected(),
		Data:        "0x" + hex.EncodeToString(tx.Data()),
		Value:       tx.Value().String(),
	}
	info.SetWrappedTransaction(wtx)
	return info
}

func (my ClientUtil) BytesToAddressHex(data32 []byte) string {
	return strings.ToLower(common.BytesToAddress(data32).Hex())
}

func (my ClientUtil) HexToAddress(address string) ebcm.WrappedAddress {
	return common.HexToAddress(strings.ToLower(address))
}
func wrappedAddress(iAddr interface{}) common.Address {
	switch v := iAddr.(type) {
	case string:
		return common.HexToAddress(strings.ToLower(v))
	}
	return iAddr.(common.Address)
}

func (my ClientUtil) HexToECDSA(private string) (ebcm.WrappedPrivateKey, error) {
	return crypto.HexToECDSA(private)
}
func (my ClientUtil) WrappedPrivateKey(iPrivate interface{}) *ecdsa.PrivateKey {
	switch v := iPrivate.(type) {
	case string:
		private, err := crypto.HexToECDSA(v)
		if err != nil {
			cc.RedItalic("wrappedPrivateKey :", err)
			return nil
		}
		return private
	}
	return iPrivate.(*ecdsa.PrivateKey)
}

func (my ClientUtil) NewLondonSigner(chainId *big.Int) ebcm.WrappedSigner {
	//return types.NewLondonSigner(chainId)
	return types.NewEIP155Signer(chainId)
}
func (my ClientUtil) WrappedLondonSigner(v ebcm.WrappedSigner) types.Signer {
	return v.(types.Signer)
}

func (my ClientUtil) NewEIP155Signer(chainId *big.Int) ebcm.WrappedSigner {
	return types.NewEIP155Signer(chainId)
}
func (my ClientUtil) WrappedEIP155Signer(v ebcm.WrappedSigner) types.EIP155Signer {
	return v.(types.EIP155Signer)
}

func typesTransaction(v ebcm.WrappedTransaction) *types.Transaction {
	return v.(*WrappedTx).Transaction
}

func (my ClientUtil) ContractAddressNonce(from string, nonce uint64) string {
	if !ebcm.IsAddress(from) {
		return ""
	}

	v := crypto.CreateAddress(
		common.HexToAddress(from),
		nonce,
	)
	return strings.ToLower(v.String())
}
