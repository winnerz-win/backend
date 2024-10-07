package kly

import (
	"jtools/cc"
	"jtools/cloud/ebcm"
	"jtools/jmath"
	"math/big"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
)

func ZMarshalJSON(stx *types.Transaction) ([]byte, error) {
	return stx.MarshalJSON()
}

func ZUnmarshalJSON(b []byte) *types.Transaction {
	stx := &types.Transaction{}
	// if err := tx.UnmarshalBinary(b); err != nil {
	// 	cc.PrintRed("a :", err)
	// }
	if err := stx.UnmarshalJSON(b); err != nil {
		cc.RedItalic("b :", err)
	}
	return stx
}

func ZCommonAddress(addressString string) common.Address {
	return common.HexToAddress(addressString)
}

/*
TxTypeFeeDelegatedValueTransfer
TxValueKeyFeePayer

TxTypeFeeDelegatedValueTransferWithRatio
TxValueKeyFeeRatioOfFeePayer

types.FeeRatio(30)
*/
func ZMakeTransactionMap(
	txtype types.TxType,
	values map[types.TxValueKeyType]interface{},
) (*types.Transaction, error) {
	return types.NewTransactionWithMap(txtype, values)
}

func ZMakeDelegateValueTransaction(
	nonce uint64,
	from string,
	to string,
	ipeb interface{},
	gasprice *big.Int,
	fee_payer string,
	fee_ratio uint8,
) (*types.Transaction, error) {
	amount := big.NewInt(0)

	if peb := jmath.VALUE(ipeb); jmath.CMP(peb, 0) > 0 {
		amount, _ = amount.SetString(peb, 10)
	}

	values := map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:  nonce,
		types.TxValueKeyFrom:   ZCommonAddress(from),
		types.TxValueKeyTo:     ZCommonAddress(to),
		types.TxValueKeyAmount: amount,
		// types.TxValueKeyGasLimit: uint64(2000000),
		types.TxValueKeyGasPrice: gasprice,
		types.TxValueKeyFeePayer: ZCommonAddress(fee_payer),
	}
	var tx_type types.TxType
	feeratio := types.FeeRatio(fee_ratio)
	if !feeratio.IsValid() {
		tx_type = types.TxTypeFeeDelegatedValueTransfer //9

		values[types.TxValueKeyGasLimit] = params.TxGasValueTransfer + params.TxGasFeeDelegated

	} else {
		tx_type = types.TxTypeFeeDelegatedValueTransferWithRatio //10

		values[types.TxValueKeyGasLimit] = uint64(36000)
		values[types.TxValueKeyFeeRatioOfFeePayer] = feeratio
	}

	return types.NewTransactionWithMap(tx_type, values)
}
func ZMakeDelegateContractTransaction(
	nonce uint64,
	from string,
	toContract string, //contract-address
	paddingdata interface{}, // []byte or PADBYTES
	ipeb interface{},
	limit uint64,
	gasprice *big.Int,
	fee_payer string,
	fee_ratio uint8,
) (*types.Transaction, error) {

	amount := big.NewInt(0)
	if peb := jmath.VALUE(ipeb); jmath.CMP(peb, 0) > 0 {
		amount, _ = amount.SetString(peb, 10)
	}
	values := map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyFrom:     ZCommonAddress(from),
		types.TxValueKeyTo:       ZCommonAddress(toContract),
		types.TxValueKeyAmount:   amount,
		types.TxValueKeyGasLimit: uint64(4000000),
		types.TxValueKeyGasPrice: gasprice,
		types.TxValueKeyFeePayer: ZCommonAddress(fee_payer),
	}

	switch pad_data := paddingdata.(type) {
	case []byte:
		values[types.TxValueKeyData] = pad_data
	case ebcm.PADBYTES:
		values[types.TxValueKeyData] = pad_data.Bytes()
	}

	var tx_type types.TxType
	feeratio := types.FeeRatio(fee_ratio)
	if !feeratio.IsValid() {
		tx_type = types.TxTypeFeeDelegatedSmartContractExecution //49
	} else {
		tx_type = types.TxTypeFeeDelegatedSmartContractExecutionWithRatio //50
		values[types.TxValueKeyFeeRatioOfFeePayer] = feeratio
	}

	return types.NewTransactionWithMap(tx_type, values)
}

func ZSignTransactionAsFeePayer(chainID *big.Int, privatekeyString string, tx *types.Transaction) (*types.Transaction, error) {
	privatekey, err := crypto.HexToECDSA(privatekeyString)
	if err != nil {
		return nil, err
	}

	signedTx, err := types.SignTxAsFeePayer(tx, types.NewEIP155Signer(chainID), privatekey)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func ZGasLimitGasFeeCancel(
	paddedData ebcm.PADBYTES,
	fromAddress, toAddress string,
	ethWEIs ...string,
) (uint64, error) {
	return uint64(8000000), nil
	// limit, err := my.XGasLimit(
	// 	paddedData,
	// 	fromAddress,
	// 	toAddress,
	// 	ethWEIs...,
	// )
	// if err == nil {
	// 	limit += params.TxGasCancel // 21000
	// }
	// return limit, err
}
