package itype

import (
	"context"
	"crypto/ecdsa"
	"jtools/cc"
	"jtools/cloud/ebcm"
	"jtools/cloud/ebcm/abi"
	"jtools/dbg"
	"jtools/jmath"
	"math/big"
	"strings"
	"time"
)

type IClient struct {
	isKlay bool
	srcURL string
	rpcURL string
	v3Key  string

	txnType ebcm.TXNTYPE
	isDebug bool

	IClientUtil

	//////////// MMA-CODE ////////
	tx_signer TxSigner
	_chain_id *big.Int
}

func (my IClient) String() string { return my.rpcURL }

func New(rpc_url string, is_klay bool, v3Key ...string) *IClient {
	if err := dbg.CutUrlSuffixSlashP(&rpc_url); err != nil {
		cc.RedBoldBG("(itype.New)", err)
		return nil
	}
	ic := IClient{
		isKlay: is_klay,
		srcURL: rpc_url,
		rpcURL: rpc_url,

		txnType: ebcm.TXN_LEGACY,
	}
	if len(v3Key) > 0 {
		if strings.HasSuffix(rpc_url, "/v3") {
			ic.v3Key = v3Key[0]
			if len(ic.v3Key) > 0 {
				if strings.HasSuffix(ic.rpcURL, "/") {
					ic.rpcURL += ic.v3Key
				} else {
					ic.rpcURL += "/" + ic.v3Key
				}
			}
		} //if

	}
	return &ic
}

//////////////////////////////////////////////////////////////////////////////////

func (my IClient) Client() interface{} {
	return my
}

func (my IClient) Host() string { return my.rpcURL }

func (my IClient) Info() ebcm.Info { return ebcm.Info{} }

func (my IClient) TXNTYPE() ebcm.TXNTYPE {
	return my.txnType
}

func (my *IClient) SetTXNTYPE(v ebcm.TXNTYPE) {
	my.txnType = v
}

func (my IClient) IsDebug() bool { return my.isDebug }

func (my *IClient) SetDebug(f ...bool) {
	my.isDebug = dbg.IsTrue(f)
}

//////////////////////////////////////////////////////////////////////////////////

type IClientUtil struct {
}

func (my IClientUtil) BytesToAddressHex(data32 []byte) string {
	return abi.ByteToAddressHexer().BytesToAddressHex(data32)
}

func (my IClientUtil) HexToAddress(address string) ebcm.WrappedAddress {
	return abi.HexToAddress(address)
}
func (my IClientUtil) ContractAddressNonce(from string, nonce uint64) string {
	return abi.ContractAddressNonce(from, nonce)
}

//////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////

func (my *IClient) SetSigner(tx_signer TxSigner) {
	my.tx_signer = tx_signer
}
func (my IClient) HasTxSigner() bool { return my.tx_signer != nil }

func (my *IClient) _Cache_ChainID() {
	if my._chain_id == nil {
		my._chain_id = my.ChainID()
	}
}

func (my *IClient) Sender() *ebcm.Sender {
	my._Cache_ChainID()

	ebcmSender := &wrapperSender{
		IClient: my,
		chainId: jmath.BigInt(my._chain_id),
	}

	if my.tx_signer != nil {
		ebcmSender.signer = my.tx_signer

	} else {
		cc.YellowItalic("[itype]IClient.Sender is not contain TxSigner.")

	}
	sender := &ebcm.Sender{
		IClient: ebcmSender,
	}
	return sender
}

// EBCMSender : only rpc.Call
func (my *IClient) EBCMSender(tx_signer_ ...TxSigner) *ebcm.Sender {
	my._Cache_ChainID()

	ebcmSender := &wrapperSender{
		IClient: my,
		chainId: jmath.BigInt(my._chain_id),
	}

	warn_log := func() {
		cc.YellowItalic("[itype]IClient.EBCMSender is not contain TxSigner.")
	}

	is_inject := false
	if len(tx_signer_) > 0 {
		if tx_signer_ != nil {
			ebcmSender.signer = tx_signer_[0]
			is_inject = true
		}
	}
	if !is_inject {
		if my.tx_signer != nil {
			ebcmSender.signer = my.tx_signer
		} else {
			warn_log()
		}
	}

	sender := &ebcm.Sender{
		IClient: ebcmSender,
	}
	return sender
}

type wrapperSender struct {
	*IClient

	signer  TxSigner
	chainId *big.Int
}

func _signer_empty_error_log() {
	cc.RedItalic("evSender.signer is nil.")
}

func (my wrapperSender) NewTransaction(
	nonce uint64,
	to interface{},
	amount interface{},
	gasLimit uint64,
	gasPrice ebcm.GasPrice,
	data []byte,
) ebcm.WrappedTransaction {
	if my.signer == nil {
		_signer_empty_error_log()
		return nil
	}
	return my.signer.NewTransaction(
		my.TXNTYPE(), //ebcm.TXN_LEGACY,
		nonce,
		to,
		amount,
		gasLimit,
		gasPrice,
		data,
	)
}

func (my *wrapperSender) SignTx(
	tx ebcm.WrappedTransaction,
	prv interface{},
) (ebcm.WrappedTransaction, error) {
	if my.signer == nil {
		_signer_empty_error_log()
		return nil, nil
	}

	if my.chainId == nil {
		my.chainId = my.ChainID()
	}
	chain_id := my.chainId

	return my.signer.SignTx(
		chain_id,
		tx,
		prv,
	)
}

func (my wrapperSender) UnmarshalBinary(buf []byte) ebcm.WrappedTransaction {
	if my.signer == nil {
		_signer_empty_error_log()
		return nil
	}
	return my.signer.UnmarshalBinary(buf)
}

func (my wrapperSender) SendTransaction(ctx context.Context, tx ebcm.WrappedTransaction) (string, error) {
	raw, _ := tx.MarshalBinary()
	return my.SendRawTransaction(raw)
}

func (my wrapperSender) CheckSendTxHashReceipt(
	tx ebcm.WrappedTransaction,
	limitSec int,
	is_debug ...bool,
) ebcm.CheckSendTxHashReceiptResult {
	hash := tx.HashHex()
	ack := ebcm.CheckSendTxHashReceiptResult{
		Hash: hash,
	}

	isDebug := dbg.IsTrue(is_debug)
	log := func(a ...interface{}) {
		if isDebug {
			cc.PurpleItalic(a...)
		}
	}

	for {
		if limitSec <= 0 {
			ack.FailMessage = "time_over"
			ack.IsTimeOver = true
			break
		}
		time.Sleep(time.Second)

		r, _, err := my.TransactionByHash(hash)
		if err != nil {
			limitSec--
			log("receipt wait -", limitSec)
			continue
		}

		if !r.IsReceiptedByHash {
			limitSec--
			log("receipt wait -", limitSec)
			continue
		}

		log("gas_fee :", r.TxFeeETH, " eth")
		ack.GasFeeETH = r.TxFeeETH
		if !r.IsError {
			ack.IsSuccess = true
		} else {
			ack.FailMessage = "tx_fail"
		}
		break
	} //for

	return ack
}

func (my wrapperSender) CheckSendTxHashToNonce(
	from interface{},
	hash string,
	limitSec int,
	is_debug ...bool,
) ebcm.CheckSendTxHashReceiptResult {
	ack := ebcm.CheckSendTxHashReceiptResult{
		Hash: hash,
	}

	ctx := context.Background()
	isDebug := dbg.IsTrue(is_debug)

	log := func(a ...interface{}) {
		if isDebug {
			cc.PurpleItalic(a...)
		}
	}

	pending, err := my.PendingNonceAt(ctx, from)
	if err != nil {
		ack.FailMessage = "pending fail"
		return ack
	}

	is_nonce_end := false
	for {
		if limitSec <= 0 {
			ack.FailMessage = "time_over"
			ack.IsTimeOver = true
			break
		}
		time.Sleep(time.Second)

		if !is_nonce_end {
			nonce, err := my.NonceAt(ctx, from)
			if err != nil {
				limitSec--
				log("nonce wait -", limitSec)
				continue
			}

			log("nonce :", nonce, " , pending :", pending)
			if nonce < pending {
				limitSec--
				log("nonce cmp wait -", limitSec)
				continue
			}

			is_nonce_end = true
		}

		r, _, _ := my.TransactionByHash(hash)
		if !r.IsReceiptedByHash {
			limitSec--
			log("receipt wait -", limitSec)
			continue
		}

		log("gas_fee :", r.TxFeeETH, " eth")
		ack.GasFeeETH = r.TxFeeETH
		if !r.IsError {
			ack.IsSuccess = true
		} else {
			ack.FailMessage = "tx_fail"
		}
		break
	} //for

	return ack
}

func (my wrapperSender) SignTooler(message_prefix ebcm.MessagePrefix) ebcm.SignTool {
	if my.signer == nil {
		_signer_empty_error_log()
		return ebcm.SignTool{}
	}
	return my.signer.SignTooler(message_prefix)
}

func (my wrapperSender) MakeWallet() ebcm.IWallet {
	if my.signer == nil {
		_signer_empty_error_log()
		return nil
	}
	return my.signer.MakeWallet()
}

func (my wrapperSender) MakeWalletFromSeed(text string, seq interface{}) ebcm.IWallet {
	if my.signer == nil {
		_signer_empty_error_log()
		return nil
	}
	return my.signer.MakeWalletFromSeed(text, seq)
}

func (my wrapperSender) Wallet(hexPrivate string) (ebcm.IWallet, error) {
	if my.signer == nil {
		_signer_empty_error_log()
		return nil, nil
	}
	return my.signer.Wallet(hexPrivate)
}

func (my wrapperSender) GetHash(wtx ebcm.WrappedTransaction) string {
	if my.signer == nil {
		_signer_empty_error_log()
		return ""
	}
	return my.signer.GetHash(wtx)
}

func (my wrapperSender) WrappedTransactionInfo(wtx ebcm.WrappedTransaction) ebcm.WrappedTxInfo {
	if my.signer == nil {
		_signer_empty_error_log()
		return ebcm.WrappedTxInfo{}
	}
	return my.signer.WrappedTransactionInfo(wtx)
}

func (my wrapperSender) HexToECDSA(private string) (ebcm.WrappedPrivateKey, error) {
	if my.signer == nil {
		_signer_empty_error_log()
		return nil, nil
	}
	return my.signer.HexToECDSA(private)
}

func (my wrapperSender) WrappedPrivateKey(iPrivate interface{}) *ecdsa.PrivateKey {
	if my.signer == nil {
		_signer_empty_error_log()
		return nil
	}
	return my.signer.WrappedPrivateKey(iPrivate)
}

func (my wrapperSender) NewLondonSigner(chainId *big.Int) ebcm.WrappedSigner {
	if my.signer == nil {
		_signer_empty_error_log()
		return nil
	}
	return my.signer.NewLondonSigner(chainId)
}

func (my wrapperSender) NewEIP155Signer(chainId *big.Int) ebcm.WrappedSigner {
	if my.signer == nil {
		_signer_empty_error_log()
		return nil
	}
	return my.signer.NewEIP155Signer(chainId)
}
