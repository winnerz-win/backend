package ecsx

import (
	"encoding/hex"
	"strings"

	"txscheduler/brix/tools/cloudx/ethwallet/abmx"
	"txscheduler/brix/tools/dbg"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/crypto/sha3"

	ebcmABI "txscheduler/brix/tools/cloudx/ebcmx/abix"
)

const (
	MethodCustomFunction = "custom"
	DeployInputPrefix    = "" //"6080604052"
)

type cMethodIDData struct {
	FuncName   string `json:"func_name"`  // transfer
	ABIString  string `json:"abi_string"` // transfer(address,uint256)
	MethodID   string `json:"method_id"`  // 0xa9059cbb
	parseInput func(data string, item *TransactionBlock)

	isAbm         bool
	params        abmx.TypeList
	parseInputABM func(rs abmx.RESULT, item *TransactionBlock)
}
type MethodIDDataList []cMethodIDData

func (my cMethodIDData) String() string    { return dbg.ToJSONString(my) }
func (my MethodIDDataList) String() string { return dbg.ToJSONString(my) }

// MakeMethodIDData : ecsx.Sender.SetCustomMethods(...) ---> MakeMethodIDDataABM (GetMethodIDABM ...)
func MakeMethodIDData(funcName, abifundName, methodID string, inputfunc func(data string, item *TransactionBlock)) cMethodIDData {
	return cMethodIDData{
		FuncName:   funcName,
		ABIString:  abifundName,
		MethodID:   methodID,
		parseInput: inputfunc,
	}
}

// MakeMethodIDDataM : ecsx.Sender.SetCustomMethods(...) ---> MakeMethodIDDataABM (GetMethodIDABM ...)
func MakeMethodIDDataM(m cMethodIDData, inputfunc func(data string, item *TransactionBlock)) cMethodIDData {
	m.isAbm = false
	m.parseInput = inputfunc
	return m
}

// //////////////////////////////////////////////////////////////////////////
func GetMethodIDABM(name string, types ...abmx.Type) cMethodIDData {
	sl := []string{}
	for _, v := range types {
		sl = append(sl, v.String())
	}
	paramsString := strings.Join(sl, ",")
	abiFuncName := name + "(" + paramsString + ")"

	item := GetMethodIDHex(abiFuncName)
	item.isAbm = true
	item.params = append(item.params, types...)
	return item
}
func MakeMethodIDDataABM(paramData cMethodIDData, f func(rs abmx.RESULT, item *TransactionBlock)) cMethodIDData {
	paramData.isAbm = true
	paramData.parseInputABM = f
	return paramData
}

/////////////////////////////////////////////////////////////////////////

func (my *Sender) SetCustomMethods(customMethods ...cMethodIDData) {
	if len(customMethods) > 0 {
		my.customMethods = append(my.customMethods, customMethods...)
	}
}

func (my Sender) checkCustomMethod(item *TransactionBlock) {
	if item.IsContract == false {
		return
	}
	if item.ContractMethod != MethodCustomFunction {
		return
	}
	for _, m := range my.customMethods {
		if m.ParseInput(item) {
			break
		}
	} //for
}

// SetCustomMethod :
func SetCustomMethod(item *TransactionBlock, customMethods ...cMethodIDData) {
	if item.IsContract == false {
		return
	}
	if item.ContractMethod != MethodCustomFunction {
		return
	}
	for _, m := range customMethods {
		if m.ParseInput(item) {
			break
		}
	} //for
}

// InputPrefix : "0x"를 제외한 hex value
func (my cMethodIDData) InputPrefix() string {
	return my.MethodID[2:]
}

// InputDataPure
func InputDataPure(data string, abmReturns interface{}) abmx.RESULT {
	v, e := hex.DecodeString(data)
	if e != nil {
		r := abmx.RESULT{}
		r.IsError = true
		return r
	}
	return abmx.ReceiptDiv(v, abmReturns)
}

func ebcm_InputDataPure(data string, typelist ebcmABI.TypeList) ebcmABI.RESULT {
	r := InputDataPure(
		data,
		abmx.EBCM_ABI_NewReturns(typelist...),
	)
	return r.RESULT
}

func getAddress64Bytes(data string) string {
	if len(data) < 64 {
		return ""
	}
	return data[24:64]
}

func (my cMethodIDData) ParseInput(item *TransactionBlock) bool {
	input := item.CustomInput
	if strings.HasPrefix(input, my.MethodID) == false {
		return false
	}
	item.ContractMethod = my.FuncName

	if my.isAbm == false {
		if my.parseInput != nil {
			if len(input) > len(my.MethodID) {
				data := input[len(my.MethodID):]
				my.parseInput(data, item)
			} else {

			}
		}
	} else {
		if my.parseInputABM != nil {
			if len(input) > len(my.MethodID) {
				data := input[len(my.MethodID):]
				cdata := InputDataPure(
					data,
					abmx.NewReturns(
						my.params...,
					),
				)
				my.parseInputABM(cdata, item)
			} else {

			}
		}
	}

	return true
}

// CheckMethodERC20 :
func CheckMethodERC20(input string, item *TransactionBlock) {
	limitSize := 8 //method name (without 0x)
	getInputSize := len(input)
	isContractCheck := false
	if input != "" && getInputSize >= limitSize {
		isContractCheck = true
	}

	if isContractCheck == false {
		input = strings.ToLower(input)
		item.CustomInput = input

		item.IsContract = false
		return
	}

	input = strings.TrimSpace(input)
	if input != "" {
		item.IsContract = true
		item.ContractAddress = item.To
		item.To = ""
	}

	if item.IsContract {
		if strings.HasPrefix(input, "0x") == false {
			input = "0x" + input
		}
	}
	input = strings.ToLower(input)
	item.CustomInput = input

	isChecked := false
	for _, erc := range methodERC20s {
		isChecked = erc.ParseInput(item)
		if isChecked {
			break
		}
	}

	if isChecked == false {
		item.ContractMethod = MethodCustomFunction
	}

}

// MakeMethodIDHex : GetMethodIDHex
func MakeMethodIDHex(abiFuncString string, isView ...bool) cMethodIDData {
	return GetMethodIDHex(abiFuncString, isView...)
}

// GetMethodIDHex : MakeMethodIDHex (Warpper)
func GetMethodIDHex(abiFuncString string, isView ...bool) cMethodIDData {
	abiFuncString = strings.ReplaceAll(abiFuncString, " ", "")
	ss := strings.Split(abiFuncString, "(")
	data := cMethodIDData{
		FuncName:  ss[0],
		ABIString: abiFuncString,
	}
	fnSignature := []byte(data.ABIString)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(fnSignature)
	methodID := hash.Sum(nil)[:4]
	data.MethodID = hexutil.Encode(methodID)

	if dbg.IsTrue2(isView...) {
		dbg.Purple(data)
	}
	return data
}
