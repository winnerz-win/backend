package ebcm

import (
	"encoding/hex"
	"jtools/cc"
	"jtools/cloud/ebcm/abi"
	"jtools/dbg"
	"strings"

	"golang.org/x/crypto/sha3"
)

// MakeMethodIDHex : GetMethodIDHex
func MakeMethodIDHex(abiFuncString string, isView ...bool) cMethodIDData {
	return GetMethodIDHex(abiFuncString, isView...)
}

type Type interface {
	String() string
}

func GetMethodIDABM(name string, types ...abi.Type) cMethodIDData {
	sl := []string{}
	for _, v := range types {
		sl = append(sl, v.String())
	}
	paramsString := strings.Join(sl, ",")
	abiFuncName := name + "(" + paramsString + ")"

	item := GetMethodIDHex(abiFuncName)
	item.params = append(item.params, types...)
	return item
}
func MakeMethodIDDataABM(paramData cMethodIDData, parseCallback func(rs abi.RESULT, item *TransactionBlock)) cMethodIDData {
	paramData.parseInputABM = parseCallback
	return paramData
}

func SetCustomMethod(caller abi.Caller, item *TransactionBlock, customMethods ...cMethodIDData) {
	if !item.IsContract {
		return
	}
	if item.ContractMethod != MethodCustomFunction {
		return
	}
	for _, m := range customMethods {
		if m.ParseInput(caller, item) {
			break
		}
	} //for
}

// ///////////////////////////////////////////////////////////////////////////
func (my cMethodIDData) ParseInput(caller abi.Caller, item *TransactionBlock) bool {
	input := item.CustomInput
	if !strings.HasPrefix(input, my.MethodID) {
		return false
	}
	item.ContractMethod = my.FuncName

	if my.parseInputABM != nil {
		if len(input) > len(my.MethodID) {
			data := input[len(my.MethodID):]
			cdata := caller.InputDataPure(
				data,
				abi.NewReturns(
					my.params...,
				),
			)
			my.parseInputABM(cdata, item)
		} else {

		}
	}

	return true
}
