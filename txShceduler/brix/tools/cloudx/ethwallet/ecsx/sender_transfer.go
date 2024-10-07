package ecsx

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsaa"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx/jwalletx"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (my Sender) Transfer(
	tokenAddress string,
	fromPrivate string,
	toAddress string,
	wei string,
	speed GasSpeed,
	isIncludeGas ...bool,
) (
	string,
	uint64, //nonce-value
	error,
) {
	from, err := jwalletx.Get(fromPrivate)
	if err != nil {
		return "", 0, err
	}

	box := my.GasBox(
		tokenAddress,
		from.Address(),
		toAddress,
		wei,
		speed,
		isIncludeGas...,
	)
	if box.Error != nil {
		return "", 0, box.Error
	}
	nonce, err := my.Nonce(from.PrivateKey())
	if err != nil {
		return "", 0, err
	}
	ntx, err := nonce.BoxTx(box)
	if err != nil {
		return "", 0, err
	}
	stx, err := ntx.Tx()
	if err != nil {
		return "", 0, err
	}
	if err := stx.Send(); err != nil {
		return "", 0, err
	}

	return stx.Hash(), nonce.NonceCount(), nil
}

func (my Sender) EthTransferNTX(
	fromPrivate string,
	toAddress string,
	wei string,
	speed GasSpeed,
	snap *GasSnapShot,
) (*NTX, error) {
	from, err := jwalletx.Get(fromPrivate)
	if err != nil {
		return &NTX{from: from.Address()}, err
	}

	fromHexAddress := common.HexToAddress(strings.ToLower(from.Address()))
	toHexAddress := common.HexToAddress(strings.ToLower(toAddress))

	value := new(big.Int)
	value.SetString(wei, 10)
	gasLimit, err := my.client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  fromHexAddress,
		To:    &toHexAddress,
		Value: value,
		Data:  []byte{},
	})
	if err != nil {
		return &NTX{from: from.Address()}, err
	}
	if snap != nil {
		gasLimit = snap.Limit
	}

	gasPrice := ecsaa.SUGGEST_GAS_PRICE(my)
	if snap != nil {
		gasPrice = snap.priceBig()
	}

	nonce, err := my.Nonce(from.PrivateKey())
	if err != nil {
		return &NTX{from: from.Address()}, err
	}
	if snap != nil {
		if snap.FixedNonce != 0 {
			nonce.nonceValue = snap.FixedNonce
		}
	}

	var data []byte

	var tx *types.Transaction
	switch my.txnType {
	case TXN_EIP_1559:
		tipPrice := ecsaa.SUGGEST_TIP_PRICE(my)

		tx = types.NewTx(&types.DynamicFeeTx{
			ChainID:   my.chainID,
			Nonce:     nonce.nonceValue,
			GasTipCap: tipPrice,
			GasFeeCap: gasPrice,
			Gas:       gasLimit,
			To:        &toHexAddress,
			Value:     value,
			Data:      data,
		})
	default:
		tx = types.NewTransaction(nonce.nonceValue, toHexAddress, value, gasLimit, gasPrice, data)
	}

	_price := big.NewInt(gasPrice.Int64())
	_limit := big.NewInt(int64(gasLimit))
	fee := _price.Mul(_price, _limit)

	newSnap := GasSnapShot{
		Limit:  gasLimit,
		Price:  gasPrice.String(),
		FeeWei: fee.String(),
	}
	ntx := &NTX{
		client:     my.client,
		privatekey: nonce.privatekey,
		chainID:    nonce.chainID,
		tx:         tx,

		nonceCount: nonce.nonceValue,

		from: from.Address(),
		to:   toAddress,
		wei:  wei,

		gasFeeWei: fee.String(),

		UserData: newTxUserData(),

		snap: newSnap,
	}

	return ntx, nil
}

func (my Sender) EthTransferSEND(ntx *NTX) (
	string, //hash
	uint64, //nonce-value
	error,
) {
	if ntx == nil {
		return "", 0, errors.New("ntx is nil")
	}
	stx, err := ntx.Tx()
	if err != nil {
		return "", 0, err
	}
	if err := stx.Send(); err != nil {
		return "", 0, err
	}

	return stx.Hash(), ntx.nonceCount, nil
}
