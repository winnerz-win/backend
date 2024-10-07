package ecsx

import (
	"context"
	"math/big"
	"time"
	"txscheduler/brix/tools/cloudx/ebcmx/abix"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsaa"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

/*
[ ATS ]
Assemble transaction steps for send.
*/
func (my *Sender) SetGasFuncDefault() *Sender {
	// my.gasFunc = func(Sender, GasSpeed) *big.Int {
	// 	return my.SuggestGasPrice()
	// }
	return my
}
func (my *Sender) AtsMakePadbytesETH() PADBYTES {
	zero := []byte{}
	return PadBytes{zero}
}
func (my *Sender) AtsMakePadbytesTOKEN(to string, value string) PADBYTES {
	return MakePadBytesABI(
		"transfer",
		abix.TypeList{
			abix.NewAddress(to),
			abix.NewUint256(value),
		},
	)
}

func (my *Sender) AtsMakePadbytesABI(functionName string, action_list abix.TypeList) PADBYTES {
	return MakePadBytesABI(functionName, action_list)
}

/////////////////////////////////////////////////////////////////////////////////

func _stob(wei string) *big.Int {
	var value *big.Int
	if jmath.CMP(wei, 0) > 0 {
		value = new(big.Int)
		value.SetString(jmath.VALUE(wei), 10)
	}
	return value
}

/////////////////////////////////////////////////////////////////////////////////

func (my *Sender) AtsGasLimit(padbytes PADBYTES, from, to string, wei string) (uint64, error) {
	to_address := common.HexToAddress(to)
	return my.client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  common.HexToAddress(from),
		To:    &to_address,
		Value: _stob(wei),
		Data:  padbytes.Bytes(),
	})
}

/////////////////////////////////////////////////////////////////////////////////

func (my *Sender) AtsNonceAt(from string) (uint64, error) {
	return my.client.NonceAt(context.Background(), common.HexToAddress(from), nil)
}
func (my *Sender) AtsPendingNonceAt(from string) (uint64, error) {
	return my.client.PendingNonceAt(context.Background(), common.HexToAddress(from))
}

/////////////////////////////////////////////////////////////////////////////////

type ATSGAS struct {
	tip string
	gas string
}

func (my *ATSGAS) AddGWEI(gwei GWEI) {
	wei := gwei.ToWEI().String()
	my.tip = jmath.VALUE(jmath.ADD(wei, my.tip))
	my.gas = jmath.VALUE(jmath.ADD(wei, my.gas))
}

func (my ATSGAS) String() string {
	m := map[string]interface{}{
		"tip": my.tip,
		"gas": my.gas,
	}
	return dbg.ToJSONString(m)
}

func (my ATSGAS) GAS_GWEI() string {
	return string(WEIToGWEI(WEI(my.gas)))
}
func (my ATSGAS) get_tip() *big.Int {
	return jmath.New(my.tip).ToBigInteger()
}
func (my ATSGAS) get_gas() *big.Int {
	return jmath.New(my.gas).ToBigInteger()
}
func (my ATSGAS) EstimateGasFeeWEI(limit interface{}) string {
	return jmath.MUL(my.gas, limit)
}
func (my ATSGAS) EstimateGasFeeETH(limit interface{}) string {
	wei := my.EstimateGasFeeWEI(limit)
	return WeiToETH(wei)
}

// CalcCostSubGas : 이더 전송일 경우 전송금액-수수료 (전액 전송시 사용)
func (my ATSGAS) CalcCostSubGas(limit interface{}, send_try_wei interface{}) (string, error) {
	if jmath.CMP(limit, 21000) != 0 {
		return "0", dbg.Error("Only ETH transfer. limit must is 21000")
	}
	fee_wei := my.EstimateGasFeeWEI(limit)
	if jmath.CMP(send_try_wei, fee_wei) < 0 {
		return "0", dbg.Error("Fee is Big than send_try_wei.")
	}
	send_wei := jmath.SUB(send_try_wei, fee_wei)
	return send_wei, nil
}

// AtsGasPrice : skipCap 은 이더 전송시 가스비를 반드시 예측해야 할경우에만 씀 ( ex: ATSGAS.CalcCostSubGas )
func (my *Sender) AtsGasPrice(skipCap ...bool) (ATSGAS, error) {
	r := ATSGAS{}
	ecsaa.SUGGEST_GAS_PRICE_2(my)
	value, err := my.client.SuggestGasPrice(context.Background())
	if err != nil {
		return r, err
	}
	r.gas = jmath.VALUE(ecsaa.GAS_ADD_MUL(value))

	if my.txnType == TXN_EIP_1559 {
		if dbg.IsTrue2(skipCap...) {
			r.tip = r.gas
		} else {
			tip, err := my.client.SuggestGasTipCap(context.Background())
			if err != nil {
				return r, err
			}
			r.tip = jmath.VALUE(tip)
		}
	}

	return r, nil
}

/////////////////////////////////////////////////////////////////////////////////

type ATSTX struct {
	limit uint64

	tx *types.Transaction

	is_signed bool
	Error     error
}

func (my ATSTX) IsSigned() bool { return my.is_signed }
func (my ATSTX) Hash() string {
	if my.tx == nil {
		return ""
	}
	return my.tx.Hash().Hex()
}
func (my ATSTX) String() string {
	m := map[string]interface{}{
		"is_signed": my.is_signed,
		"error":     my.Error,
	}
	my.tx.Gas() //limit
	my.tx.ChainId()
	my.tx.Data()
	my.tx.Hash()
	if my.tx != nil {
		m["hash"] = my.tx.Hash().Hex()

		m["cost"] = jmath.VALUE(my.tx.Cost()) //Cost returns gas * gasPrice + value.
		m["gas_price"] = jmath.VALUE(my.tx.GasPrice())
		m["gas_tip_cap"] = my.tx.GasTipCap() //Max Priority
		m["gas_fee_cap"] = my.tx.GasFeeCap() //Max

		m["nonce"] = my.tx.Nonce()
		m["limit"] = my.limit
		m["limit_gas"] = my.tx.Gas()
	}
	return dbg.ToJSONString(m)
}
func (my ATSTX) Limit() uint64 { return my.limit }
func (my ATSTX) Nonce() uint64 {
	if my.tx == nil {
		return 0
	}
	return my.tx.Nonce()
}

func (my *Sender) AtsTx(data PADBYTES, to string, wei string, nonce uint64, limit uint64, gas ATSGAS) ATSTX {
	to_address := common.HexToAddress(to)
	value := _stob(wei)

	var tx *types.Transaction
	if my.txnType == TXN_EIP_1559 {
		tx = types.NewTx(&types.DynamicFeeTx{
			ChainID:   my.chainID,
			Nonce:     nonce,
			GasTipCap: gas.get_tip(),
			GasFeeCap: gas.get_gas(),
			Gas:       limit,
			To:        &to_address,
			Value:     value,
			Data:      data.Bytes(),
		})
	} else {
		tx = types.NewTransaction(
			nonce,
			to_address,
			value,
			limit,
			gas.get_gas(),
			data.Bytes(),
		)
	}
	return ATSTX{
		is_signed: false,
		tx:        tx,
		limit:     limit,
	}
}

/////////////////////////////////////////////////////////////////////////////////

func (my *Sender) AtsSignTx(privatekey_text string, ats_tx ATSTX) ATSTX {
	r := ATSTX{
		is_signed: true,
		Error:     nil,
		limit:     ats_tx.limit,
	}
	privatekey, err := crypto.HexToECDSA(privatekey_text)
	if err != nil {
		r.Error = dbg.Error("[AtsSignTx] HexToECDSA :", err)
		return r
	}

	chainID, err := my.client.NetworkID(context.Background())
	if err != nil {
		r.Error = dbg.Error("[AtsSignTx] NetworkID :", err)
		return r
	}

	var signer types.Signer
	if ats_tx.tx.Type() == 2 {
		signer = types.NewLondonSigner(chainID)
	} else {
		signer = types.NewEIP155Signer(chainID)
	}

	signedTx, err := types.SignTx(ats_tx.tx, signer, privatekey)
	if err != nil {
		r.Error = dbg.Error("[AtsSignTx] SignTx :", err)
		return r
	}

	r.tx = signedTx
	return r
}

/////////////////////////////////////////////////////////////////////////////////

func (my *Sender) AtsSend(ats_tx ATSTX) error {
	if !ats_tx.is_signed {
		return dbg.Error("[AtsSend] tx unsigned.")
	}
	return my.client.SendTransaction(context.Background(), ats_tx.tx)
}

/////////////////////////////////////////////////////////////////////////////////

type ats_receipt struct {
	IsSuccess   bool   `json:"is_success"`
	FailMessage string `json:"fail_message,omitempty"`
}

func (my ats_receipt) String() string { return dbg.ToJSONString(my) }

func (my *Sender) ATSCheckReceipt(hash string, limitSec int, is_debug ...bool) ats_receipt {
	isDebug := dbg.IsTrue2(is_debug...)
	ack := ats_receipt{}
	for {
		if limitSec <= 0 {
			ack.FailMessage = "time_over"
			break
		}
		time.Sleep(time.Second)
		r, _, _, err := my.TransactionByHash(hash)
		if err != nil {
			limitSec--
			if isDebug {
				dbg.PurpleItalic("receipt wait -", limitSec)
			}
			continue
		}
		if !r.IsReceiptedByHash {
			limitSec--
			if isDebug {
				dbg.PurpleItalic("receipt wait -", limitSec)
			}
			continue
		}
		if !r.IsError {
			ack.IsSuccess = true
		} else {
			ack.FailMessage = "tx_fail"
		}
		break
	} //for
	return ack
}
