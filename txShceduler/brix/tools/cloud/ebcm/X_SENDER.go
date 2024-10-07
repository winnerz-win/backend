package ebcm

import (
	"context"
	"crypto"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"txscheduler/brix/tools/cloud/ebcm/abi"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
)

var (
	ZERO = "0"
	//TxCancelUserFee : tx_cancel_user_fee
	TxCancelUserFee = errors.New("tx_cancel_user_fee")
)

// IWallet :
type IWallet interface {
	Index() interface{}
	String() string
	PrivateKey() string
	Address() string
	CompareAddress(cmpAddress string) bool
}

type WrappedAddress interface {
	Bytes() []byte
	Hex() string
	String() string
}
type WrappedPrivateKey interface {
	Public() crypto.PublicKey
}
type WrappedSigner interface {
	ChainID() *big.Int
}

type WrappedTransaction interface {
	EncodeRLP(w io.Writer) error
	MarshalBinary() ([]byte, error)
	UnmarshalBinary(b []byte) error
	Cost() *big.Int
	GasPrice() *big.Int
	Nonce() uint64
	Gas() uint64 //limit

	ChainId() *big.Int
	Data() []byte
	Value() *big.Int //wei
}

type WrappedTxInfo struct {
	wrapped_tx  WrappedTransaction `json:"-"`
	IsSigned    bool               `json:"is_signed"`
	ChainID     string             `json:"chain_id"`
	Hash        string             `json:"hash"`
	Cost        string             `json:"cost"`
	GasPrice    string             `json:"gas_price"`
	GasTipCap   string             `json:"gas_tip_cap"`
	GasFeeCap   string             `json:"gas_fee_cap"`
	Nonce       uint64             `json:"nonce"`
	Gas         uint64             `json:"gas"`
	TxnType     TXNTYPE            `json:"txn_type"`
	TxnTypeText string             `json:"txn_type_text"`
	Value       string             `json:"value"`
	To          string             `json:"to"`
	Protected   bool               `json:"protected"`
	Data        string             `json:"data"`
}

func (my WrappedTxInfo) String() string                                { return dbg.ToJsonString(my) }
func (my *WrappedTxInfo) SetWrappedTransaction(wtx WrappedTransaction) { my.wrapped_tx = wtx }
func (my WrappedTxInfo) WrappedTransaction() WrappedTransaction        { return my.wrapped_tx }

type CheckSendTxHashReceiptResult struct {
	IsSuccess   bool   `json:"is_success"`
	Hash        string `json:"hash"`
	IsTimeOver  bool   `json:"is_time_over"`
	FailMessage string `json:"fail_message,omitempty"`
	GasFeeETH   string `json:"gas_fee_eth,omitempty"`
}

func (my CheckSendTxHashReceiptResult) String() string { return dbg.ToJsonString(my) }

//////////////////////////////////////////////////////////////////////////////////////////

type iHex interface {
	Hex() string
}
type CallMsg struct {
	From  interface{}
	To    interface{}
	Value *big.Int
	Data  []byte
}

func (my CallMsg) String() string {
	v := map[string]interface{}{
		"from":   my.FromAddress(),
		"to":     my.ToAddress(),
		"amount": my.Amount(),
		"data":   my.DataHex(),
	}
	return dbg.ToJsonString(v)
}

func MakeCallMsg(from, to interface{}, value interface{}, data []byte) CallMsg {
	return CallMsg{
		From:  from,
		To:    to,
		Value: jmath.BigInt(value),
		Data:  data,
	}
}

func (my CallMsg) FromAddress() string {
	address := ""
	switch v := my.From.(type) {
	case iHex:
		address = v.Hex()
	case string:
		address = v
	}
	return dbg.TrimToLower(address)
}
func (my CallMsg) ToAddress() string {
	address := ""
	switch v := my.To.(type) {
	case iHex:
		address = v.Hex()
	case string:
		address = v
	}
	return dbg.TrimToLower(address)
}
func (my CallMsg) Amount() string {
	return jmath.VALUE(my.Value)
}
func (my CallMsg) DataHex() string {
	return hex.EncodeToString(my.Data)
}

type Client interface {
	NetworkID(ctx context.Context) (*big.Int, error)
	Close()
}

type TXNTYPE int

func (my TXNTYPE) String() string {
	switch my {
	case TXN_LEGACY:
		return "Legacy"
	case TXN_EIP_1559:
		return "EIP-1559"
	}
	return ""
}
func (my TXNTYPE) Uint8() uint8   { return uint8(my) }
func (my TXNTYPE) Uint16() uint16 { return uint16(my) }

const (
	TXN_LEGACY   = TXNTYPE(0) //"Legacy"
	TXN_EIP_1559 = TXNTYPE(2) //"EIP-1559"
)

//////////////////////////////////////////////////////////////////////////////////////////

/*
[EIP-1559]
Gas == Tip 일경우 Fee계산 정확.
Tip을 따로 설정할경우  block의 Base값에 Tx의 Max Priority값을 더한것이 gasPrice가 된다.  (block.Base + tx.TipCap)
*/
type GasPrice struct {
	Tip *big.Int //Max Priority
	Gas *big.Int //Max
}

func (my *GasPrice) AddGWEI_ALL(gwei GWEI) {
	wei := gwei.ToWEI().String()
	my.Tip = jmath.BigInt(jmath.ADD(wei, my.Tip))
	my.Gas = jmath.BigInt(jmath.ADD(wei, my.Gas))
}
func (my *GasPrice) AddGWEI_TIP(gwei GWEI) {
	wei := gwei.ToWEI().String()

	sum := jmath.ADD(wei, my.Tip)
	if jmath.CMP(sum, my.Gas) >= 0 {
		my.Tip = jmath.BigInt(my.Gas)
	} else {
		my.Tip = jmath.BigInt(sum)
	}
}

func (my *GasPrice) AddGWEI_EACH(_gas GWEI, _tip GWEI) {
	my.Gas = jmath.BigInt(jmath.ADD(_gas, my.Gas))
	my.Tip = jmath.BigInt(jmath.ADD(_tip, my.Tip))
}

func (my GasPrice) String() string {
	m := map[string]interface{}{
		"tip": jmath.VALUE(my.Tip),
		"gas": jmath.VALUE(my.Gas),
	}
	return dbg.ToJsonString(m)
}
func (my GasPrice) EstimateGasFeeWEI(limit interface{}) string {
	return jmath.MUL(my.Gas, limit)
}
func (my GasPrice) EstimateGasFeeETH(limit interface{}) string {
	wei := my.EstimateGasFeeWEI(limit)
	return WeiToETH(wei)
}

// CalcCostSubGas : 이더 전송일 경우 전송금액-수수료 (전액 전송시 사용) ( ex: suggestGasPrice( is_skip_tip_cat == true ) )
func (my GasPrice) CalcCostSubGas(limit interface{}, send_try_wei interface{}) (string, error) {
	if jmath.CMP(limit, 21000) != 0 {
		return "0", dbg.Error("Only ETH transfer. limit must is 21000")
	}
	fee_wei := my.EstimateGasFeeWEI(limit)
	if jmath.CMP(send_try_wei, fee_wei) < 0 {
		return "0", fmt.Errorf("Fee is Big than send_try_wei. (try_amount:%v)", send_try_wei)
	}
	send_wei := jmath.SUB(send_try_wei, fee_wei)
	return send_wei, nil
}

//////////////////////////////////////////////////////////////////////////////////////////
// Sender
//////////////////////////////////////////////////////////////////////////////////////////

type Info struct {
	Host      string
	Key       string
	NetworkID *big.Int
	IsDebug   bool
}

type Sender struct {
	IClient
	///////////////////////////
}

func NewSender(
	client IClient,
) *Sender {
	sender := &Sender{
		IClient: client,
	}
	return sender

}

func (my Sender) RpcClient() interface{} { return my.IClient.Client() }

func (my Sender) Balance(address string) string {
	v, err := my.IClient.BalanceAt(
		context.Background(),
		my.HexToAddress(address),
		nil,
	)
	if err != nil {
		dbg.RedItalic("Balance :", err)
	}
	return jmath.VALUE(v)
}

func (my Sender) Price(address string) string {
	return WeiToETH(
		my.Balance(address),
	)
}

func (my Sender) BlockNumber(defaultNumber ...string) string {
	v, err := my.IClient.BlockNumber(context.Background())
	if err != nil {
		if len(defaultNumber) > 0 {
			return defaultNumber[0]
		}
	}
	return jmath.VALUE(v)
}

func MakePadBytesABI(pureName string, type_list abi.TypeList) PADBYTES {
	return type_list.GetBytes(pureName, true)
}

func (my Sender) MakePadBytesABI(pureName string, type_list abi.TypeList) PADBYTES {
	return type_list.GetBytes(pureName, true)
}

//////////////////////////////////////////////////////////////////////////////////////
