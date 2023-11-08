package ebcmx

import (
	"strings"
	"txscheduler/brix/tools/cloudx/ebcmx/abix"
)

type PMethodIDDataMap map[string]MethodIDData

func MethodID(funcName string, args abix.TypeList, parseCallback func(rs abix.RESULT, item *TransactionBlock)) MethodIDData {
	params := ""
	if args != nil || len(args) > 0 {
		sl := []string{}
		for _, arg := range args {
			sl = append(sl, arg.String())
		}
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

func MethodIDDataMap(items ...MethodIDData) PMethodIDDataMap {
	ddm := PMethodIDDataMap{}
	for _, item := range items {
		ddm[item.MethodID] = item
	}
	return ddm
}
func (my PMethodIDDataMap) AddMethodIDData(item MethodIDData) {
	my[item.MethodID] = item
}

func (my PMethodIDDataMap) Merge(target PMethodIDDataMap) {
	for k, v := range target {
		my[k] = v
	}
}

func SetCustomMethodMap(caller abix.Caller, item *TransactionBlock, ddm PMethodIDDataMap) {
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
		cdata := caller.InputDataPure(
			data,
			abix.NewReturns(
				v.params...,
			),
		)
		v.parseInputABM(cdata, item)
	}
}
