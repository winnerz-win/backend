package abi

import (
	"fmt"
	"jtools/cc"
	"jtools/dbg"
)

type ISender interface {
	CallContract(from, to string, data []byte) ([]byte, error)
	BytesToAddressHex(data32 []byte) string
}

type Bytes32Stringer interface {
	BytesToAddressHex(data32 []byte) string
}

// InputParser :abi.byteToAddressHexer{} , wclie ebcm.IClient
type InputParser interface {
	BytesToAddressHex(data32 []byte) string
	ContractAddressNonce(from string, nonce uint64) string
}

func Call2(
	finder ISender,
	contract string,
	method Method,
	caller string,
	f func(rs RESULT),
	isLogs ...bool,
) (call_err error) {
	isLog := false
	if dbg.IsTrue(isLogs) {
		isLog = true
	}
	defer func() {
		if e := recover(); e != nil {
			call_err = fmt.Errorf("%v\n%v", e, dbg.Stack())
		}
	}()

	inputBytes := method.Params.GetBytes(method.Name)

	if isLog {
		debug_receipt(inputBytes, true)
	}

	receipt, err := finder.CallContract(
		caller,
		contract,
		inputBytes,
	)
	if err != nil {
		cc.RedItalic("abi.Call2[", method.Name, "] :", err)
		return err
	}

	//cc.White("BUF_SIZE:", len(receipt))

	if isLog {
		cc.PurpleItalic("abi.Call2[", method.Name, "] RAW_DATA : ", len(receipt))
		debug_receipt(receipt)
	}

	result := receiptDiv(finder, receipt, method.Returns)

	if f == nil {
		cc.PurpleItalic(result)
	} else {
		f(result)
	}

	return nil
}

func ReceiptDivDirect(receipt []byte, method_returns TypeList) (rs RESULT) {
	return receiptDiv(
		ByteToAddressHexer(),
		receipt,
		method_returns,
	)
}

func ReceiptDiv(stringer Bytes32Stringer, receipt []byte, rts TypeList) (rs RESULT) {
	defer func() {
		if e := recover(); e != nil {
			cc.PurpleItalic("abi.ReceiptDiv : ", e)
			rs = RESULT{}
			rs.IsError = true
		}
	}()
	if len(receipt)%32 != 0 {
		//cc.Red("ReceiptDiv__Size :", len(receipt))
		rs = RESULT{}
		rs.IsError = true
		return rs
	}

	if DEBUG_MODE {
		debug_receipt(receipt)
	}

	rs = receiptDiv(stringer, receipt, rts)
	return rs
}

func receiptDiv(stringer Bytes32Stringer, receipt []byte, rts TypeList) RESULT {
	if len(receipt) == 0 {
		rs := RESULT{}
		return rs
	}

	sl := [][]byte{}
	for len(receipt) > 0 {
		v := receipt[:32]
		sl = append(sl, v)
		receipt = receipt[32:]
	} //for

	/*
		hex - index
		20 - 1
		40 - 2
		60 - 3
		a0 - 5
		c0 - 6
		e0 - 7
		100 - 8
		120 - 9
		140 - 10
	*/

	result := RESULT{}

	headerIndex := 0
	for _, v := range rts {

		switch v.name {
		default:
			var val interface{}
			val, headerIndex = _getValue(stringer, v, sl, 0, headerIndex)
			result.EBCMaddItem(v.name, val)

		case ITupleFlex.name:
			pos := _getIndex(sl[headerIndex])

			tupleList := TupleList{}
			list := ResultItemList{}

			flex := v.FlexTypeList()

			start_pos := pos
			index_pos := start_pos
			for _, data := range flex {
				var val interface{}
				val, index_pos = _getValue(stringer, data, sl, start_pos, index_pos)
				list.EBCMadd(data.name, val)
			} //for

			tupleList = append(tupleList, list)
			result.EBCMaddItem(v.name, tupleList)
			headerIndex++

		case ITupleArray.name:
			pos := _getIndex(sl[headerIndex])
			count_idx := pos

			count := _getInt(sl[count_idx]) //2

			pos_offset := count_idx + 1

			flex := v.FlexTypeList()
			tupleList := TupleList{}

			for i := 0; i < count; i++ {
				start_index := pos_offset + (len(v.Datas) * i)
				start_pos := start_index
				index_pos := start_pos

				list := ResultItemList{}
				for _, data := range flex {
					var val interface{}
					val, index_pos = _getValue(stringer, data, sl, start_pos, index_pos)
					list.EBCMadd(data.name, val)
				} //for

				tupleList = append(tupleList, list)

			}
			result.EBCMaddItem(v.name, tupleList)
			headerIndex++

		case ITupleArrayFlex.name:
			pos := _getIndex(sl[headerIndex])
			count_idx := pos

			count := _getInt(sl[count_idx]) //2

			pos_offset := count_idx + 1

			flex := v.FlexTypeList()
			tupleList := TupleList{}

			for i := 0; i < count; i++ {
				pos_idx := _getIndex(sl[pos_offset+i])

				start_index := count_idx + pos_idx + 1
				start_pos := start_index
				index_pos := start_pos

				list := ResultItemList{}
				for _, data := range flex {
					var val interface{}
					val, index_pos = _getValue(stringer, data, sl, start_pos, index_pos)
					list.EBCMadd(data.name, val)
				} //for

				tupleList = append(tupleList, list)

			}
			result.EBCMaddItem(v.name, tupleList)
			headerIndex++
		}
	}

	return result
}
