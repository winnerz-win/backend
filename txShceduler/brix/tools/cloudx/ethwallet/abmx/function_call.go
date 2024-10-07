package abmx

import (
	"encoding/hex"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"

	"github.com/ethereum/go-ethereum/accounts/abi"

	ebcmABI "txscheduler/brix/tools/cloudx/ebcmx/abix"
)

type ISender interface {
	CallContract(from, to string, data []byte) ([]byte, error)
}

// SenderCallContract : &ecsx.Sender.CallContract
type SenderCallContract func(from, to string, data []byte) ([]byte, error)

// Call :
func Call(finder ISender, contract string, method Method, caller string, f func(rs RESULT), isLogs ...bool) error {
	isLog := false
	if dbg.IsTrue2(isLogs...) {
		isLog = true
	}
	abiSafe := abi.ABI{
		Methods: method.getMethod(),
	}

	// if method.Returns.err != nil {
	// 	return method.Returns.err
	// }
	params := method.Params.getParames()
	inputBytes, err := abiSafe.Pack(method.Name, params...)
	if err != nil {
		dbg.Red("abi.Pack :", err)
		return err
	}
	receipt, err := finder.CallContract(
		caller,
		contract,
		inputBytes,
	)
	if err != nil {
		dbg.Red("abi.Call[", method.Name, "] :", err)
		return err
	}

	//isLog = true
	if isLog {
		dbg.Purple("abi.Call[", method.Name, "] RAW_DATA : ", len(receipt))
		v := hex.EncodeToString(receipt)
		dbg.Purple(v)
		loop := 0
		for len(v) > 0 {
			dbg.Purple("[", loop, "]", v[:64], "(", jmath.VALUE("0x"+v[:64]), ")")
			v = v[64:]
			loop++
		} //for

	}

	result := receiptDiv(receipt, method.Returns)

	if f == nil {
		dbg.Purple(result)
	} else {
		f(result)
	}

	return nil

}

// debugBufs :
func debugBufs(buf [][]byte) {
	defer dbg.Cyan("----------------------------------------------------------------------")
	dbg.Cyan("----------------------------------------------------------------------")

	for i, v := range buf {
		s := hex.EncodeToString(v)
		dbg.Cyan("[", i, "]", s)
	} //for

}

func callABI(
	finder ebcmABI.ISender,
	contract string,
	method ebcmABI.Method,
	caller string,
	f func(rs ebcmABI.RESULT),
	isLogs ...bool,
) error {

	return Call(
		finder,
		contract,
		Method{
			Name:    method.Name,
			Params:  EBCM_ABI_NewParams(method.Params...),
			Returns: EBCM_ABI_NewReturns(method.Returns...),
		},
		caller,
		func(rs RESULT) {
			f(rs.RESULT)
		},
		isLogs...,
	)
}

func GetEBCM(
	inputDataPure ebcmABI.DelegateInputDataPure,
) ebcmABI.Caller {
	item := ebcmABI.Caller{
		Call:          callABI,
		InputDataPure: inputDataPure,
	}
	return item
}
