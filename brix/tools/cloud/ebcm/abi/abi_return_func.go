package abi

import (
	"encoding/hex"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"

	"math/big"
)

func _toAddress(stringer Bytes32Stringer, data32 []byte) string {
	return stringer.BytesToAddressHex(data32)
}

func _amountRe(data32 []byte) string {
	v := big.NewInt(0)
	v = v.SetBytes(data32)
	return v.String()
}

func _amountValue(data32 []byte) string {
	return jmath.VALUE(_amountRe(data32))
}

func _boolValue(data32 []byte) bool {
	if data32[len(data32)-1] == 0 {
		return false
	} else {
		return true
	}
}

func _bytesHex(data32 []byte) string {
	return dbg.TrimToLower("0x" + hex.EncodeToString(data32))
}
func _getInt(data32 []byte) int {
	return jmath.Int("0x" + hex.EncodeToString(data32))
}

func _getIndex(data32 []byte) int {
	return _getInt(data32) / 32
}

func _set_flex_value(start_idx int, data32 []byte, sl [][]byte, f func(data []byte)) {
	buf := []byte{}
	pos_idx := _getIndex(data32)

	size_idx := pos_idx + start_idx

	size := _getInt(sl[size_idx])
	size_idx++

	loop := size / 32
	for j := 0; j < loop; j++ {
		data32 = sl[size_idx]
		buf = append(buf, data32...)
		size_idx++
	}
	if v := size % 32; v > 0 {
		data32 = sl[size_idx]
		buf = append(buf, data32[:v]...)
	}

	f(buf)
}

func _set_flex_array_value(start_idx int, data32 []byte, sl [][]byte, f func(data []byte)) {
	pos := _getIndex(data32)

	count_idx := pos + start_idx
	count := _getInt(sl[count_idx])

	pos_offset := count_idx + 1

	for i := 0; i < count; i++ {
		pos_idx := _getIndex(sl[pos_offset+i])

		size_idx := pos_offset + pos_idx
		size := _getInt(sl[size_idx])
		size_idx++

		buf := []byte{}
		loop := size / 32
		for j := 0; j < loop; j++ {
			data32 = sl[size_idx]
			buf = append(buf, data32...)
			size_idx++
		}
		if v := size % 32; v > 0 {
			data32 = sl[size_idx]
			buf = append(buf, data32[:v]...)
		}
		f(buf)
	}

}

func _set_fixed_array_value(start_idx int, data32 []byte, sl [][]byte, f func(data []byte)) {
	pos := _getIndex(data32)
	count_idx := pos + start_idx

	count := _getInt(sl[count_idx])

	data_idx := count_idx + 1
	for i := 0; i < count; i++ {
		data32 = sl[data_idx+i]
		f(data32)
	}
}

func _getValue(stringer Bytes32Stringer, t Type, sl [][]byte, start_idx, offset int) (interface{}, int) {
	var val interface{}
	data32 := sl[offset]
	switch t.name {
	case Address.name:
		val = _toAddress(stringer, data32)

	case Uint256.name,
		Uint128.name,
		Uint112.name:
		val = _amountValue(data32)

	case Uint64.name:
		val = jmath.Uint64(_amountRe(data32))

	case Uint32.name:
		val = uint32(jmath.Uint64(_amountRe(data32)))

	case Uint16.name:
		val = uint16(jmath.Uint64(_amountRe(data32)))

	case Uint8.name:
		val = uint8(jmath.Uint64(_amountRe(data32)))

	case Bool.name:
		val = _boolValue(data32)

	case Bytes32.name:
		val = _bytesHex(data32)

	case String.name:
		_set_flex_value(
			start_idx, data32, sl,
			func(data []byte) {
				val = string(data)
			},
		)

	case Bytes.name:
		_set_flex_value(
			start_idx, data32, sl,
			func(data []byte) {
				val = _bytesHex(data)
			},
		)

	case StringArray.name:
		ss := []string{}
		_set_flex_array_value(
			start_idx, data32, sl,
			func(data []byte) {
				ss = append(ss, string(data))
			},
		)
		val = ss

	case BytesArray.name:
		ss := []string{}
		_set_flex_array_value(
			start_idx, data32, sl,
			func(data []byte) {
				ss = append(ss, _bytesHex(data))
			},
		)
		val = ss

	case AddressArray.name:
		obj := []string{}
		_set_fixed_array_value(
			start_idx, data32, sl,
			func(data []byte) {
				val := _toAddress(stringer, data)
				obj = append(obj, val)
			},
		)
		val = obj

	case Uint256Array.name,
		Uint128Array.name,
		Uint112Array.name:
		obj := []string{}
		_set_fixed_array_value(
			start_idx, data32, sl,
			func(data []byte) {
				val := _amountValue(data)
				obj = append(obj, val)
			},
		)
		val = obj

	case Uint64Array.name:
		obj := []uint64{}
		_set_fixed_array_value(
			start_idx, data32, sl,
			func(data []byte) {
				val := jmath.Uint64(_amountRe(data))
				obj = append(obj, val)
			},
		)
		val = obj

	case Uint32Array.name:
		obj := []uint32{}
		_set_fixed_array_value(
			start_idx, data32, sl,
			func(data []byte) {
				val := uint32(jmath.Uint64(_amountRe(data)))
				obj = append(obj, val)
			},
		)
		val = obj

	case Uint16Array.name:
		obj := []uint16{}
		_set_fixed_array_value(
			start_idx, data32, sl,
			func(data []byte) {
				val := uint16(jmath.Uint64(_amountRe(data)))
				obj = append(obj, val)
			},
		)
		val = obj

	case Uint8Array.name:
		obj := []uint8{}
		_set_fixed_array_value(
			start_idx, data32, sl,
			func(data []byte) {
				val := uint8(jmath.Uint64(_amountRe(data)))
				obj = append(obj, val)
			},
		)
		val = obj

	case BoolArray.name:
		obj := []bool{}
		_set_fixed_array_value(
			start_idx, data32, sl,
			func(data []byte) {
				val := _boolValue(data)
				obj = append(obj, val)
			},
		)
		val = obj

	case Bytes32Array.name:
		obj := []string{}
		_set_fixed_array_value(
			start_idx, data32, sl,
			func(data []byte) {
				val := _bytesHex(data)
				obj = append(obj, val)
			},
		)
		val = obj

	case ITuple.name:

		flex := t.FlexTypeList()
		tupleList := TupleList{}

		idx := offset
		list := ResultItemList{}
		for _, data := range flex {
			var val interface{}
			val, idx = _getValue(stringer, data, sl, start_idx+offset, idx)
			list.EBCMadd(data.name, val)
		}
		tupleList = append(tupleList, list)
		val = tupleList
		offset = idx - 1

	case ITupleFlex.name:
		pos := _getIndex(sl[offset])

		flex := t.FlexTypeList()
		tupleList := TupleList{}
		list := ResultItemList{}

		start_pos := pos + start_idx
		index_pos := start_pos
		for _, data := range flex {
			var val interface{}
			val, index_pos = _getValue(stringer, data, sl, start_pos, index_pos)
			list.EBCMadd(data.name, val)
		}
		tupleList = append(tupleList, list)
		val = tupleList

	case ITupleArray.name:

		pos := _getIndex(sl[offset])
		count_idx := pos + start_idx

		count := _getInt(sl[count_idx])

		pos_offset := count_idx + 1

		flex := t.FlexTypeList()
		tupleList := TupleList{}

		for j := 0; j < count; j++ {
			start_index := pos_offset + (len(t.Datas) * j)
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
		val = tupleList

	case ITupleArrayFlex.name:
		pos := _getIndex(sl[offset])
		count_idx := pos + start_idx

		count := _getInt(sl[count_idx])

		pos_offset := count_idx + 1

		flex := t.FlexTypeList()
		tupleList := TupleList{}

		for j := 0; j < count; j++ {
			pos_idx := _getIndex(sl[pos_offset+j])

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
		val = tupleList

	} //t

	offset++
	return val, offset
}
