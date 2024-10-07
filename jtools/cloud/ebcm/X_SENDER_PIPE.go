package ebcm

import (
	"context"
	"errors"
	"jtools/cc"
	"jtools/dbg"
	"jtools/jmath"
	"math/big"
)

func AmountZero() *big.Int {
	return jmath.BigInt(0)
}

type XSendResult struct {
	From     string `bson:"from" json:"from"`
	To       string `bson:"to" json:"to"`
	DataHex  string `bson:"data_hex" json:"data_hex"`
	Limit    uint64 `bson:"limit" json:"limit"`
	Nonce    uint64 `bson:"nonce" json:"nonce"`
	GasPrice string `bson:"gas_price" json:"gas_price"`
	Hash     string `bson:"hash" json:"hash"`
	Amount   string `bson:"amount,omitempty" json:"amount,omitempty"`
}

func (my XSendResult) String() string { return dbg.ToJsonString(my) }

type XNonce2 func(nonce, pending uint64) bool   //0
type XLimit func(gaslimit uint64) uint64        //1
type XGasPrice func(gasPrice GasPrice) GasPrice //2
const (
	XDebug     = "XPipeDebug"
	XWait10Min = 60 * 10
	XWait30Min = 60 * 30
)

func _check_xpipe2_callback(callbacks ...interface{}) ([]interface{}, bool) {
	if len(callbacks) == 0 {
		return []interface{}{
			XNonce2(nil),
			XLimit(nil),
			XGasPrice(nil),
		}, true
	}

	is_nonce := false
	is_limit := false
	is_gas := false
	list := []interface{}{}
	debug_index := 0
	for i := 0; i < len(callbacks); i++ {
		switch v := callbacks[i].(type) {
		case XNonce2:
			if !is_nonce {
				list = append(list, callbacks[i])
				is_nonce = true
			} else {
				return nil, false
			}
		case XLimit:
			if !is_limit {
				list = append(list, callbacks[i])
				is_limit = true
			} else {
				return nil, false
			}
		case XGasPrice:
			if !is_gas {
				list = append(list, callbacks[i])
				is_gas = true
			} else {
				return nil, false
			}

		case func(nonce, pending uint64) bool:
			if !is_nonce {
				list = append(list, XNonce2(callbacks[i].(func(nonce, pending uint64) bool)))
				is_nonce = true
			} else {
				return nil, false
			}
		case func(gaslimit uint64) uint64:
			if !is_limit {
				list = append(list, XLimit(callbacks[i].(func(gaslimit uint64) uint64)))
				is_limit = true
			} else {
				return nil, false
			}
		case func(gasPrice GasPrice) GasPrice:
			if !is_gas {
				list = append(list, XGasPrice(callbacks[i].(func(gasPrice GasPrice) GasPrice)))
				is_gas = true
			} else {
				return nil, false
			}

		case string:
			if v == XDebug {
				if debug_index == 0 {
					list = append(list, XDebug)
					debug_index = i
				}
			}

		default:
			return nil, false
		} //switch

	} //for
	if debug_index > 0 {
		for i := debug_index - 1; i >= 0; i-- {
			list[i+1] = list[i]
		}
		list[0] = XDebug
	}
	if !is_nonce {
		list = append(list, XNonce2(nil))
	}
	if !is_limit {
		list = append(list, XLimit(nil))
	}
	if !is_gas {
		list = append(list, XGasPrice(nil))
	}

	return list, true
}

func (my Sender) XPipe2(
	fromPrivate string,
	toAddress string,
	padBytes PADBYTES,
	wei string,
	callbacks ...interface{}, //XNonce,  XLimit, XGasPrice
) (XSendResult, error) {

	result := XSendResult{}

	callbacks, ok := _check_xpipe2_callback(callbacks...)
	if !ok {
		return result, dbg.Error("[XPipe2]invalid_callbacks_func")
	}

	fromWallet, err := my.Wallet(fromPrivate)
	if err != nil {
		return result, dbg.Error("[XPipe2]Wallet :", err)
	}

	privateKey, err := my.HexToECDSA(fromPrivate)
	if err != nil {
		return result, dbg.Error("[XPipe2]HexToECDSA :", err)
	}

	if jmath.CMP(wei, 0) <= 0 {
		wei = "0"
	}

	amount := jmath.BigInt(wei)

	from := my.HexToAddress(fromWallet.Address())
	to := my.HexToAddress(toAddress)

	var v_err error
	nonce := uint64(0)
	gasLimit := uint64(0)
	var gasPrice GasPrice

	isDebug := false
	debug := func(a ...interface{}) {
		if !isDebug {
			return
		}
		cc.PurpleItalic(a...)
	}

	for _, v := range callbacks {
		switch callback := v.(type) {
		case string:
			if callback == XDebug && !isDebug {
				isDebug = true
				debug("[XPipe2]============================= START")
				defer debug("[XPipe2]============================= END")
			}

		case XNonce2:
			nonce, v_err = my.NonceAt(
				context.Background(),
				from,
			)
			if v_err != nil {
				return result, dbg.Error("[XPipe2]NonceAt :", v_err)
			}
			debug("nonce :", nonce)
			if callback != nil {
				if pending, err := my.PendingNonceAt(
					context.Background(),
					from,
				); err != nil {
					return result, dbg.Error("[XPipe2]PendingNonceAt :", err)
				} else {
					debug("pending :", pending)
					if !callback(nonce, pending) {
						return result, dbg.Error("[XPipe2]nonce(", nonce, ") / pending(", pending, ")")
					}
				}
			}
		case XLimit:
			gasLimit, v_err = my.EstimateGas(
				context.Background(),
				CallMsg{
					from,
					to,
					amount,
					padBytes.Bytes(),
				},
			)
			if v_err != nil {
				return result, dbg.Error("[XPipe2]EstimateGas :", v_err)
			}
			debug("gas_limit :", gasLimit)
			if callback != nil {
				gasLimit = callback(gasLimit)
			}

		case XGasPrice:
			gasPrice, v_err = my.SuggestGasPrice(
				context.Background(),
			)
			if v_err != nil {
				return result, dbg.Error("[XPipe2]SuggestGasPrice :", v_err)
			}
			debug("gas_price :", gasPrice)
			if callback != nil {
				gasPrice = callback(gasPrice)
			}

		default:

		} //switch
	}

	//step.4
	tx := my.NewTransaction(
		nonce,
		to,
		amount,
		gasLimit,
		gasPrice,
		padBytes.Bytes(),
	)
	if tx == nil {
		return result, errors.New("[XPipe2]NewTransaction : tx is null")
	}

	stx, err := my.SignTx(
		tx,
		privateKey,
	)
	if err != nil {
		return result, dbg.Error("[XPipe2]SignTx :", err)
	}

	debug("stx :", my.WrappedTransactionInfo(stx))

	//step.5
	hash, err := my.SendTransaction(
		context.Background(),
		stx,
	)
	if err != nil {
		return result, dbg.Error("[XPipe2]SendTransaction :", v_err)
	}
	debug("send_hash:", hash)

	result = XSendResult{
		From:     fromWallet.Address(),
		To:       toAddress,
		DataHex:  padBytes.Hex(),
		Limit:    gasLimit,
		Nonce:    nonce,
		GasPrice: jmath.VALUE(gasPrice.Gas),
		Hash:     hash,
	}
	return result, nil
}
func (my Sender) XPipe(
	fromPrivate string,
	toAddress string,
	padBytes PADBYTES,
	wei string,
	limit_callback func(gasLimit uint64) uint64,
	nonce_callback func(nonce uint64) uint64,
	gasp_callback func(gasPrice GasPrice) GasPrice,
	result_callback func(r XSendResult),
) error {
	if result_callback == nil {
		return errors.New("result_callback is null")
	}

	fromWallet, err := my.Wallet(fromPrivate)
	if err != nil {
		return err
	}

	privateKey, err := my.HexToECDSA(fromPrivate)
	if err != nil {
		return err
	}

	if jmath.CMP(wei, 0) <= 0 {
		wei = "0"
	}

	amount := jmath.BigInt(wei)

	from := my.HexToAddress(fromWallet.Address())
	to := my.HexToAddress(toAddress)

	//step.1
	gasLimit, err := my.EstimateGas(
		context.Background(),
		CallMsg{
			from,
			to,
			amount,
			padBytes.Bytes(),
		},
	)
	if err != nil {
		return err
	}
	if limit_callback != nil {
		gasLimit = limit_callback(gasLimit)
	}

	//step.2
	nonce, err := my.NonceAt(
		context.Background(),
		from,
	)
	if err != nil {
		return err
	}
	if nonce_callback != nil {
		nonce = nonce_callback(nonce)
	}

	//step.3
	gasPrice, err := my.SuggestGasPrice(
		context.Background(),
	)
	if err != nil {
		return err
	}
	if gasp_callback != nil {
		gasPrice = gasp_callback(gasPrice)
	}

	//step.4
	tx := my.NewTransaction(
		nonce,
		to,
		amount,
		gasLimit,
		gasPrice,
		padBytes.Bytes(),
	)
	if tx == nil {
		return errors.New("tx is null")
	}

	stx, err := my.SignTx(
		tx,
		privateKey,
	)
	if err != nil {
		return err
	}

	//step.5
	hash, err := my.SendTransaction(
		context.Background(),
		stx,
	)
	_ = hash
	if err == nil {
		result_callback(
			XSendResult{
				From:     fromWallet.Address(),
				To:       toAddress,
				DataHex:  padBytes.Hex(),
				Limit:    gasLimit,
				Nonce:    nonce,
				GasPrice: jmath.VALUE(gasPrice.Gas),
				Hash:     hash,
			},
		)
	}

	return err
}

func (my Sender) XPipeFixedGAS(
	fromPrivate string,
	toAddress string,
	padBytes PADBYTES,
	fixed_gas_wei string,
	wei string,
	limit_callback func(gasLimit uint64) uint64,
	nonce_callback func(nonce uint64) uint64,
	result_callback func(r XSendResult),
) error {
	if result_callback == nil {
		return errors.New("result_callback is null")
	}

	fromWallet, err := my.Wallet(fromPrivate)
	if err != nil {
		return err
	}

	privateKey, err := my.HexToECDSA(fromPrivate)
	if err != nil {
		return err
	}

	if jmath.CMP(wei, 0) <= 0 {
		wei = "0"
	}

	amount := jmath.BigInt(wei)

	from := my.HexToAddress(fromWallet.Address())
	to := my.HexToAddress(toAddress)

	//step.1
	gasLimit, err := my.EstimateGas(
		context.Background(),
		CallMsg{
			from,
			to,
			amount,
			padBytes.Bytes(),
		},
	)
	if err != nil {
		return err
	}
	if limit_callback != nil {
		gasLimit = limit_callback(gasLimit)
	}

	//step.2
	nonce, err := my.NonceAt(
		context.Background(),
		from,
	)
	if err != nil {
		return err
	}
	if nonce_callback != nil {
		nonce = nonce_callback(nonce)
	}

	//step.3 --- Gas fixed
	gasPrice := GasPrice{
		Gas: jmath.BigInt(fixed_gas_wei),
		Tip: jmath.BigInt(fixed_gas_wei),
	}

	//step.4
	tx := my.NewTransaction(
		nonce,
		to,
		amount,
		gasLimit,
		gasPrice,
		padBytes.Bytes(),
	)
	if tx == nil {
		return errors.New("tx is null")
	}

	stx, err := my.SignTx(
		tx,
		privateKey,
	)
	if err != nil {
		return err
	}

	//step.5
	hash, err := my.SendTransaction(
		context.Background(),
		stx,
	)
	_ = hash
	if err == nil {
		result_callback(
			XSendResult{
				From:     fromWallet.Address(),
				To:       toAddress,
				DataHex:  padBytes.Hex(),
				Limit:    gasLimit,
				Nonce:    nonce,
				GasPrice: jmath.VALUE(gasPrice.Gas),
				Hash:     hash,
			},
		)
	}

	return err
}

func (my Sender) XSendCoinAll(
	fromPrivate string,
	toAddress string,
) (*XSendResult, error) {

	fromWallet, err := my.Wallet(fromPrivate)
	if err != nil {
		return nil, err
	}

	privateKey, err := my.HexToECDSA(fromPrivate)
	if err != nil {
		return nil, err
	}

	remain_amount := my.Balance(fromWallet.Address())
	if jmath.CMP(remain_amount, 0) <= 0 {
		return nil, errors.New("remain coin is zero.")
	}

	from := my.HexToAddress(fromWallet.Address())
	to := my.HexToAddress(toAddress)

	padBytes := PadByteETH()

	wei := jmath.BigInt(remain_amount)

	//step.1
	gasLimit, err := my.EstimateGas(
		context.Background(),
		CallMsg{
			from,
			to,
			wei,
			padBytes.Bytes(),
		},
	)
	if err != nil {
		return nil, err
	}

	//step.2
	nonce, err := my.NonceAt(
		context.Background(),
		from,
	)
	if err != nil {
		return nil, err
	}

	//step.3
	gasPrice, err := my.SuggestGasPrice(
		context.Background(),
		true,
	)
	if err != nil {
		return nil, err
	}

	send_amount, err := gasPrice.CalcCostSubGas(gasLimit, remain_amount)
	if err != nil {
		return nil, err
	}
	//step.4
	tx := my.NewTransaction(
		nonce,
		to,
		jmath.BigInt(send_amount),
		gasLimit,
		gasPrice,
		padBytes.Bytes(),
	)
	if tx == nil {
		return nil, errors.New("tx is null")
	}

	stx, err := my.SignTx(
		tx,
		privateKey,
	)
	if err != nil {
		return nil, err
	}

	//step.5
	hash, err := my.SendTransaction(
		context.Background(),
		stx,
	)
	if err != nil {
		return nil, err
	}

	result := &XSendResult{
		From:     fromWallet.Address(),
		To:       toAddress,
		DataHex:  padBytes.Hex(),
		Limit:    gasLimit,
		Nonce:    nonce,
		GasPrice: jmath.VALUE(gasPrice.Gas),
		Hash:     hash,
		Amount:   send_amount,
	}

	return result, nil
}
