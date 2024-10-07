package ebcmx

import (
	"strings"
	"txscheduler/brix/tools/cloudx/ebcmx/abix"
	"txscheduler/brix/tools/dbg"

	"golang.org/x/crypto/sha3"
)

const (
	MethodCustomFunction = "custom"
	DeployInputPrefix    = "" //"6080604052"
)

type MethodIDData struct {
	FuncName  string `json:"func_name"`  // transfer
	ABIString string `json:"abi_string"` // transfer(address,uint256)
	MethodID  string `json:"method_id"`  // 0xa9059cbb

	params        abix.TypeList
	parseInputABM func(rs abix.RESULT, item *TransactionBlock)
}
type MethodIDDataList []MethodIDData

func (my MethodIDData) String() string     { return dbg.ToJSONString(my) }
func (my MethodIDDataList) String() string { return dbg.ToJSONString(my) }

// MakeMethodIDHex : GetMethodIDHex
func MakeMethodIDHex(abiFuncString string, isView ...bool) MethodIDData {
	return GetMethodIDHex(abiFuncString, isView...)
}

// GetMethodIDHex : MakeMethodIDHex (Warpper)
func GetMethodIDHex(abiFuncString string, isView ...bool) MethodIDData {
	abiFuncString = strings.ReplaceAll(abiFuncString, " ", "")
	ss := strings.Split(abiFuncString, "(")
	data := MethodIDData{
		FuncName:  ss[0],
		ABIString: abiFuncString,
	}
	fnSignature := []byte(data.ABIString)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(fnSignature)
	methodID := hash.Sum(nil)[:4]
	data.MethodID = Hexutil_Encode(methodID)

	if dbg.IsTrue2(isView...) {
		dbg.Purple(data)
	}
	return data
}

type Type interface {
	String() string
}

func GetMethodIDABM(name string, types ...abix.Type) MethodIDData {
	sl := []string{}
	for _, v := range types {
		// typeSTring := v.(Type).String()
		// sl = append(sl, typeSTring)
		sl = append(sl, v.String())
	}
	paramsString := strings.Join(sl, ",")
	abiFuncName := name + "(" + paramsString + ")"

	item := GetMethodIDHex(abiFuncName)
	item.params = append(item.params, types...)
	return item
}
func MakeMethodIDDataABM(paramData MethodIDData, f func(rs abix.RESULT, item *TransactionBlock)) MethodIDData {
	paramData.parseInputABM = f
	return paramData
}

func SetCustomMethod(caller abix.Caller, item *TransactionBlock, customMethods ...MethodIDData) {
	if item.IsContract == false {
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
func (my MethodIDData) ParseInput(caller abix.Caller, item *TransactionBlock) bool {
	input := item.CustomInput
	if strings.HasPrefix(input, my.MethodID) == false {
		return false
	}
	item.ContractMethod = my.FuncName

	if my.parseInputABM != nil {
		if len(input) > len(my.MethodID) {
			data := input[len(my.MethodID):]
			cdata := caller.InputDataPure(
				data,
				abix.NewReturns(
					my.params...,
				),
			)
			my.parseInputABM(cdata, item)
		} else {

		}
	}

	return true
}
