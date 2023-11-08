package ecsx

import (
	"errors"
	"math/big"
	"txscheduler/brix/tools/cloudx/ebcmx"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx/jwalletx"
	"txscheduler/brix/tools/jmath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (my TSender) TransferFunction(
	fromPrivate string,
	data PADBYTES,
	wei string,
	speed GasSpeed,
	fixedNonce ...uint64,
) (
	string, //hash
	uint64, //nonce-value
	error,
) {
	from, err := jwalletx.Get(fromPrivate)
	if err != nil {
		return "", 0, err
	}
	nonce, err := my.Nonce(from.PrivateKey())
	if err != nil {
		return "", 0, err
	}

	ntx, err := nonce.NTX(
		data,
		wei,
		speed,
		fixedNonce...,
	)
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

func (my TSender) TransferFunctionFixedGAS(
	fromPrivate string,
	data PADBYTES,
	wei string,
	gasPair []*big.Int,
	txFeeWeiAllow func(feeWEI string) bool,
	fixedNonce ...uint64,
) (
	string, //hash
	uint64, //nonce-value
	error,
) {
	from, err := jwalletx.Get(fromPrivate)
	if err != nil {
		return "", 0, err
	}
	nonce, err := my.Nonce(from.PrivateKey())
	if err != nil {
		return "", 0, err
	}

	ntx, err := nonce.NTX_FixedGAS(
		data,
		wei,
		gasPair,
		fixedNonce...,
	)
	if err != nil {
		return "", 0, err
	}

	if txFeeWeiAllow != nil {
		txFeeWei := jmath.MUL(ntx.gasLimit, ntx.gasRight)
		isAllow := txFeeWeiAllow(txFeeWei)
		if !isAllow {
			return "", 0, ebcmx.TxCancelUserFee
		}
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

func (my TSender) ebcm_TransferFunction(
	fromPrivate string,
	padBytes ebcmx.PADBYTES,
	wei string,
	speed ebcmx.GasSpeed,
) (
	string, //hash
	uint64, //nonce-value
	error,
) {

	//padBytes := data.(PadBytes)

	return my.TransferFunction(
		fromPrivate,
		padBytes,
		wei,
		GasSpeed(speed),
	)

}

func (my TSender) ebcm_TransferFunctionFixedGAS(
	fromPrivate string,
	padBytes ebcmx.PADBYTES,
	wei string,
	gasPair []*big.Int,
	txFeeWeiAllow func(feeWEI string) bool, //isContinue
) (
	string, //hash
	uint64, //nonce-value
	error,
) {

	//padBytes := data.(PadBytes)

	return my.TransferFunctionFixedGAS(
		fromPrivate,
		padBytes,
		wei,
		gasPair,
		txFeeWeiAllow,
	)

}

func (my TSender) TransferFuncNTX(
	fromPrivate string,
	data PadBytes,
	wei string,
	speed GasSpeed,
	snap *GasSnapShot,
) (*NTX, error) {
	from, err := jwalletx.Get(fromPrivate)
	if err != nil {
		return &NTX{from: from.Address()}, err
	}
	nonce, err := my.Nonce(from.PrivateKey())
	if err != nil {
		return &NTX{from: from.Address()}, err
	}

	var ntx *NTX = nil
	var nErr error = nil
	snapNTX := func(paddedData PadBytes, wei string, speed GasSpeed) (*NTX, error) {
		contractHexAddress := common.HexToAddress(my.contractAddress)

		var ethValue *big.Int
		if !jmath.IsUnderZero(wei) {
			ethValue = new(big.Int)
			ethValue.SetString(jmath.VALUE(wei), 10)
		}

		gasLimit := snap.Limit

		nonce, _ := my.Nonce(from.PrivateKey())
		if snap.FixedNonce != 0 {
			nonce.nonceValue = snap.FixedNonce
		}

		gasPrice := snap.priceBig()

		if ethValue == nil {
			ethValue = big.NewInt(0) // in wei (0 eth)
		}
		tx := types.NewTransaction(
			nonce.nonceValue,
			contractHexAddress,
			ethValue,
			gasLimit,
			gasPrice,
			paddedData.cPadBytes,
		)

		newSnap := GasSnapShot{
			Limit:  gasLimit,
			Price:  gasPrice.String(),
			FeeWei: gasPrice.Mul(gasPrice, big.NewInt(int64(gasLimit))).String(),
		}
		ntx := &NTX{
			client:     my.Sender.client,
			tx:         tx,
			privatekey: nonce.privatekey,
			chainID:    nonce.chainID,

			nonceCount: nonce.nonceValue,
			from:       from.Address(),
			to:         my.contractAddress, //contract-Address
			wei:        wei,
			gasFeeWei:  newSnap.FeeWei,

			snap: newSnap,
		}
		return ntx, nil

	}
	if snap != nil {
		if snap.Limit == 0 && snap.FixedNonce != 0 {
			ntx, nErr = nonce.NTX(
				data,
				wei,
				speed,
				snap.FixedNonce,
			)
		} else {
			ntx, nErr = snapNTX(data, wei, speed)
		}
	} else {
		ntx, nErr = nonce.NTX(
			data,
			wei,
			speed,
		)
	}

	return ntx, nErr
}

func (my TSender) TransferFuncSEND(ntx *NTX) (
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
