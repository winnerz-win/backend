package ecsx

import (
	"math/big"
	"txscheduler/brix/tools/cloudx/ebcmx"
	"txscheduler/brix/tools/cloudx/ethwallet/abmx"
	"txscheduler/brix/tools/dbg"
)

func (my TSender) Allowance(owner, spender string) string {
	wei := "0"
	abmx.Call(
		my,
		my.contractAddress,
		abmx.Method{
			Name: "allowance",
			Params: abmx.NewParams(
				abmx.NewAddress(owner),
				abmx.NewAddress(spender),
			),
			Returns: abmx.NewReturns(
				abmx.Uint256,
			),
		},
		owner,
		func(rs abmx.RESULT) {
			wei = rs.Uint256(0)
		},
	)
	return wei
}

func (my TSender) Approve(privateKey string, spender string, amount string) string {
	hash := ""
	h, n, e := my.TransferFunction(
		privateKey,
		MakePadBytes(
			"approve(address,uint256)",
			func(pad Appender) {
				pad.SetParamHeader(2)
				pad.SetAddress(0, spender)
				pad.SetAmount(1, amount)
			},
		),
		"0",
		GasFast,
	)
	_ = n
	if e != nil {
		dbg.Red(e)
	}
	hash = h
	return hash
}

func (my TSender) ApproveAll(privateKey string, spender string) string {
	maxValue := "115792089237316195423570985008687907853269984665640564039457584007913129639935"
	return my.Approve(privateKey, spender, maxValue)
}

func (my TSender) TransferFrom(privateKey string, owner, to string, amount string) string {
	hash := ""
	err := my.XPipe(
		privateKey,
		my.contractAddress,
		PadBytesTransferFrom(
			owner,
			to,
			amount,
		),
		"0",
		GasFast,
		nil, nil, nil,
		func(r XSendResult) {
			hash = r.Hash
		},
	)
	if err != nil {
		dbg.Red(err)
	}
	return hash
}

func (my TSender) ebcm_TransferToken(
	privateKey string,
	to string,
	tokenWEI string,
	speed ebcmx.GasSpeed,
	limitCB func(gasLimit uint64) uint64,
	nonceCB func(nonce uint64) uint64,
	gaspCB func(gasPrice ebcmx.XGasPrice) ebcmx.XGasPrice,
	resultCB func(r ebcmx.XSendResult),

) error {
	return my.ebcm_XPipe(
		privateKey,
		my.contractAddress,
		PadBytesTransfer(
			to,
			tokenWEI,
		),
		"0",
		speed,
		limitCB,
		nonceCB,
		gaspCB,
		resultCB,
	)
}

func (my TSender) ebcm_TransferTokenFixedGAS(
	privateKey string,
	to string,
	tokenWEI string,
	limitCB func(gasLimit uint64) uint64,
	nonceCB func(nonce uint64) uint64,
	gasPair []*big.Int,
	txFeeWeiAllow func(feeWEI string) bool,
	resultCB func(r ebcmx.XSendResult),

) error {
	return my.ebcm_XPipeFixedGAS(
		privateKey,
		my.contractAddress,
		PadBytesTransfer(
			to,
			tokenWEI,
		),
		"0",
		limitCB,
		nonceCB,
		gasPair,
		txFeeWeiAllow,
		resultCB,
	)
}

func (my TSender) ebcm_Write(
	privateKey string,
	padBytes ebcmx.PADBYTES,
	wei string,
	speed ebcmx.GasSpeed,
	limitCB func(gasLimit uint64) uint64,
	nonceCB func(nonce uint64) uint64,
	gaspCB func(gasPrice ebcmx.XGasPrice) ebcmx.XGasPrice,
	resultCB func(r ebcmx.XSendResult),
) error {
	return my.ebcm_XPipe(
		privateKey,
		my.contractAddress,
		padBytes,
		wei,
		speed,
		limitCB,
		nonceCB,
		gaspCB,
		resultCB,
	)
}

func (my TSender) ebcm_WriteFixedGAS(
	privateKey string,
	padBytes ebcmx.PADBYTES,
	wei string,
	limitCB func(gasLimit uint64) uint64,
	nonceCB func(nonce uint64) uint64,
	gasPair []*big.Int,
	txFeeWeiAllow func(feeWEI string) bool,
	resultCB func(r ebcmx.XSendResult),
) error {
	return my.ebcm_XPipeFixedGAS(
		privateKey,
		my.contractAddress,
		padBytes,
		wei,
		limitCB,
		nonceCB,
		gasPair,
		txFeeWeiAllow,
		resultCB,
	)
}
