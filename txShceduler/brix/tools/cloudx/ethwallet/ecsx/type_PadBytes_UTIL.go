package ecsx

import "txscheduler/brix/tools/cloudx/ebcmx"

/*
	y_tsender.go
*/

//PadBytesETH: (ETH-transfer)
func PadBytesETH() PadBytes {
	return PadBytes{
		cPadBytes: []byte{},
	}
}

//TransferPadBytes : PadBytesTrasnfer (ERC20-transfer)
func TransferPadBytes(toAddress string, wei string) PadBytes {
	return MakePadBytes(
		"transfer(address,uint256)",
		func(pad Appender) {
			pad.SetParamHeader(2)
			pad.SetAddress(0, toAddress)
			pad.SetAmount(1, wei)
		},
	)
}

//PadBytesTrasnfer : TransferPadBytes (ERC20-transfer)
func PadBytesTransfer(toAddress string, tokenWei string) PadBytes {
	return TransferPadBytes(toAddress, tokenWei)
}

//PadBytesApprove : (ERC20-approve)
func PadBytesApprove(spender, amount string) PadBytes {
	return MakePadBytes(
		"approve(address,uint256)",
		func(pad Appender) {
			pad.SetParamHeader(2)
			pad.SetAddress(0, spender)
			pad.SetAmount(1, amount)
		},
	)
}

//PadBytesApproveAll : (ERC20-approve(MAX))
func PadBytesApproveAll(spender string) PadBytes {
	return PadBytesApprove(
		spender,
		"115792089237316195423570985008687907853269984665640564039457584007913129639935",
	)
}

//PadBytesTransferFrom : (ERC20-transferFrom)
func PadBytesTransferFrom(from, to, amount string) PadBytes {
	return MakePadBytes(
		"transferFrom(address,address,uint256)",
		func(pad Appender) {
			pad.SetParamHeader(3)
			pad.SetAddress(0, from)
			pad.SetAddress(1, to)
			pad.SetAmount(2, amount)
		},
	)
}

func ebcm_PadBytesETH() ebcmx.PADBYTES {
	return PadBytesETH()
}

func ebcm_TransferPadBytes(toAddress string, wei string) ebcmx.PADBYTES {
	return TransferPadBytes(toAddress, wei)
}

func ebcm_PadBytesTransfer(toAddress string, tokenWei string) ebcmx.PADBYTES {
	return PadBytesTransfer(toAddress, tokenWei)
}

func ebcm_PadBytesApprove(spender, amount string) ebcmx.PADBYTES {
	return PadBytesApprove(spender, amount)
}

func ebcm_PadBytesApproveAll(spender string) ebcmx.PADBYTES {
	return PadBytesApproveAll(spender)
}

func ebcm_PadBytesTransferFrom(from, to, amount string) ebcmx.PADBYTES {
	return PadBytesTransferFrom(from, to, amount)
}
