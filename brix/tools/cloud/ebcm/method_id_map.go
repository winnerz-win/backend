package ebcm

import (
	"encoding/hex"
	"strings"
	"txscheduler/brix/tools/cloud/ebcm/abi"
	"txscheduler/brix/tools/dbg"

	"golang.org/x/crypto/sha3"
)

const (
	MethodCustomFunction = "custom"
	DeployInputPrefix    = "" //"6080604052"
)

type cMethodIDData struct {
	FuncName  string `json:"func_name"`  // transfer
	ABIString string `json:"abi_string"` // transfer(address,uint256)
	MethodID  string `json:"method_id"`  // 0xa9059cbb

	params        abi.TypeList
	parseInputABM func(rs abi.RESULT, item *TransactionBlock)
	parseDeploy   func(wclient IClient, rs abi.RESULT, item *TransactionBlock)
}
type MethodIDDataList []cMethodIDData

func (my cMethodIDData) String() string    { return dbg.ToJsonString(my) }
func (my MethodIDDataList) String() string { return dbg.ToJsonString(my) }

type PMethodIDDataMap map[string]cMethodIDData

func MethodID(
	funcName string,
	args abi.TypeList,
	parseCallback func(rs abi.RESULT, item *TransactionBlock),
) cMethodIDData {
	params := ""
	if args != nil || len(args) > 0 {
		sl := []string{}
		abi.MakeInputMethodParam(&sl, args, false)
		params = strings.Join(sl, ",")
	}
	abiFuncName := funcName + "(" + params + ")"

	item := GetMethodIDHex(abiFuncName)
	if len(params) > 0 {
		item.params = append(item.params, args...)
	}

	item.parseInputABM = parseCallback
	return item
}

func MethodIDDataMap(items ...cMethodIDData) PMethodIDDataMap {
	ddm := PMethodIDDataMap{}
	for _, item := range items {
		ddm[item.MethodID] = item
	}
	return ddm
}

func (my PMethodIDDataMap) Merge(target PMethodIDDataMap) {
	for k, v := range target {
		my[k] = v
	}
}

func SetCustomMethodMap(stringer abi.Bytes32Stringer, item *TransactionBlock, ddm PMethodIDDataMap) {
	if !item.IsContract {
		return
	}
	const methodIdSize = 10
	if len(item.CustomInput) < methodIdSize { //0x2db9a59c
		return
	}

	methodID := item.CustomInput[:methodIdSize]
	if v, do := ddm[methodID]; do {
		item.ContractMethod = v.FuncName

		if v.parseInputABM == nil {
			return
		}
		data := item.CustomInput[methodIdSize:]
		cdata := InputDataPure(
			stringer,
			data,
			abi.NewReturns(
				v.params...,
			),
		)
		v.parseInputABM(cdata, item)
	}
}

// Encode encodes b as a hex string with 0x prefix.
func Encode(b []byte) string {
	enc := make([]byte, len(b)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], b)
	return string(enc)
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
	data.MethodID = Encode(methodID)

	if dbg.IsTrue(isView) {
		dbg.PurpleItalic(data)
	}
	return data
}
