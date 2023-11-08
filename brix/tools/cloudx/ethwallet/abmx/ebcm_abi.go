package abmx

import (
	ebcmABI "txscheduler/brix/tools/cloudx/ebcmx/abix"
)

func swapABIResturnsType(t ebcmABI.Type) Type {
	switch t.String() {
	case ebcmABI.None.String():
		return None
	case ebcmABI.Address.String():
		return Address
	case ebcmABI.Uint256.String():
		return Uint256
	case ebcmABI.Uint.String():
		return Uint
	case ebcmABI.Uint128.String():
		return Uint128
	case ebcmABI.Uint112.String():
		return Uint112
	case ebcmABI.Uint64.String():
		return Uint64
	case ebcmABI.Uint32.String():
		return Uint32
	case ebcmABI.Uint16.String():
		return Uint16
	case ebcmABI.Uint8.String():
		return Uint8
	case ebcmABI.Bool.String():
		return Bool
	case ebcmABI.String.String():
		return String
	case ebcmABI.Bytes.String():
		return Bytes
	case ebcmABI.Bytes32.String():
		return Bytes32

	case ebcmABI.AddressArray.String():
		return AddressArray
	case ebcmABI.Uint256Array.String():
		return Uint256Array
	case ebcmABI.UintArray.String():
		return UintArray
	case ebcmABI.Uint128Array.String():
		return Uint128Array
	case ebcmABI.Uint112Array.String():
		return Uint112Array
	case ebcmABI.Uint64Array.String():
		return Uint64Array
	case ebcmABI.Uint32Array.String():
		return Uint32Array
	case ebcmABI.Uint16Array.String():
		return Uint16Array
	case ebcmABI.Uint8Array.String():
		return Uint8Array
	case ebcmABI.BoolArray.String():
		return BoolArray
	case ebcmABI.StringArray.String():
		return StringArray
	case ebcmABI.BytesArray.String():
		return BytesArray
	case ebcmABI.Bytes32Array.String():
		return Bytes32Array

	case ebcmABI.ITuple.String():
		sl := []*fixedType{}
		for _, data := range t.Datas {
			sl = append(sl,
				swapABIResturnsType(data.(ebcmABI.Type)).(*fixedType),
			)
		}
		return Tuple(sl...)

	case ebcmABI.ITupleFixedArray.String():
		sl := []*fixedType{}
		for _, data := range t.Datas {
			sl = append(sl,
				swapABIResturnsType(data.(ebcmABI.Type)).(*fixedType),
			)
		}
		return TupleArray(sl...)

	case ebcmABI.ITupleFlexArray.String():
		sl := []*fixedType{}
		for _, data := range t.Datas {
			sl = append(sl,
				swapABIResturnsType(data.(ebcmABI.Type)).(*fixedType),
			)
		}
		return TupleArray(sl...)
	}
	return nil
}
func EBCM_ABI_NewReturns(ps ...ebcmABI.Type) abiReturns { //interface{} {
	sl := TypeList{}
	for _, v := range ps {
		sl = append(sl, swapABIResturnsType(v))
	} //for
	return NewReturns(sl...)

}

// /////////////////////////////////////////////////////////////////////////
func swapParamsToABI(t Type) ebcmABI.Type {
	switch t.String() {
	case ebcmABI.None.String():
		return ebcmABI.None
	case ebcmABI.Address.String():
		return ebcmABI.Address
	case ebcmABI.Uint256.String():
		return ebcmABI.Uint256
	case ebcmABI.Uint.String():
		return ebcmABI.Uint
	case ebcmABI.Uint128.String():
		return ebcmABI.Uint128
	case ebcmABI.Uint112.String():
		return ebcmABI.Uint112
	case ebcmABI.Uint64.String():
		return ebcmABI.Uint64
	case ebcmABI.Uint32.String():
		return ebcmABI.Uint32
	case ebcmABI.Uint16.String():
		return ebcmABI.Uint16
	case ebcmABI.Uint8.String():
		return ebcmABI.Uint8
	case ebcmABI.Bool.String():
		return ebcmABI.Bool
	case ebcmABI.String.String():
		return ebcmABI.String
	case ebcmABI.Bytes.String():
		return ebcmABI.Bytes
	case ebcmABI.Bytes32.String():
		return ebcmABI.Bytes32

	case ebcmABI.AddressArray.String():
		return ebcmABI.AddressArray
	case ebcmABI.Uint256Array.String():
		return ebcmABI.Uint256Array
	case ebcmABI.UintArray.String():
		return ebcmABI.UintArray
	case ebcmABI.Uint128Array.String():
		return ebcmABI.Uint128Array
	case ebcmABI.Uint112Array.String():
		return ebcmABI.Uint112Array
	case ebcmABI.Uint64Array.String():
		return ebcmABI.Uint64Array
	case ebcmABI.Uint32Array.String():
		return ebcmABI.Uint32Array
	case ebcmABI.Uint16Array.String():
		return ebcmABI.Uint16Array
	case ebcmABI.Uint8Array.String():
		return ebcmABI.Uint8Array
	case ebcmABI.BoolArray.String():
		return ebcmABI.BoolArray
	case ebcmABI.StringArray.String():
		return ebcmABI.StringArray
	case ebcmABI.BytesArray.String():
		return ebcmABI.BytesArray
	case ebcmABI.Bytes32Array.String():
		return ebcmABI.Bytes32Array
	}
	return ebcmABI.None
}

func swapABIParams(t ebcmABI.Type) AbiParam {
	stringsToPtrs := func(data interface{}) []interface{} {
		sl := []interface{}{}
		for _, v := range data.([]string) {
			sl = append(sl, v)
		}
		return sl
	}
	switch t.String() {
	case ebcmABI.None.String():
		return AbiParam{}
	case ebcmABI.Address.String():
		return NewAddress(t.Data.(string))
	case ebcmABI.Uint256.String():
		return NewUint256(t.Data.(string))
	case ebcmABI.Uint.String():
		return NewUint(t.Data.(string))
	case ebcmABI.Uint128.String():
		return NewUint128(t.Data.(string))
	case ebcmABI.Uint112.String():
		return NewUint112(t.Data.(string))
	case ebcmABI.Uint64.String():
		return NewUint64(t.Data.(uint64))
	case ebcmABI.Uint32.String():
		return NewUint32(t.Data.(uint32))
	case ebcmABI.Uint16.String():
		return NewUint16(t.Data.(uint16))
	case ebcmABI.Uint8.String():
		return NewUint8(t.Data.(uint8))
	case ebcmABI.Bool.String():
		return NewBool(t.Data.(bool))
	case ebcmABI.String.String():
		return NewString(t.Data.(string))

	case ebcmABI.Bytes.String():
		return NewBytes(t.Data.([]byte))
	case ebcmABI.BytesArray.String():
		return NewBytesArray(t.Data.([][]byte))

	case ebcmABI.Bytes32.String():
		return NewBytes32(t.Data.(string))

	case ebcmABI.AddressArray.String():
		return NewAddressArray(t.Data.([]string)...)
	case ebcmABI.Uint256Array.String():
		return NewUint256Array(stringsToPtrs(t.Data)...)
	case ebcmABI.UintArray.String():
		return NewUint256Array(stringsToPtrs(t.Data)...)

	case ebcmABI.Uint128Array.String():
		return NewUint128Array(stringsToPtrs(t.Data)...)

	case ebcmABI.Uint112Array.String():
		return NewUint112Array(stringsToPtrs(t.Data)...)

	case ebcmABI.Uint64Array.String():
		return NewUint64Array(t.Data.([]uint64)...)

	case ebcmABI.Uint32Array.String():
		return NewUint32Array(t.Data.([]uint32)...)
	case ebcmABI.Uint16Array.String():
		return NewUint16Array(t.Data.([]uint16)...)
	case ebcmABI.Uint8Array.String():
		return NewUint8Array(t.Data.([]uint8)...)
	case ebcmABI.BoolArray.String():
		return NewBoolArray(t.Data.([]bool)...)
	case ebcmABI.StringArray.String():
		return NewStringArray(t.Data.([]string)...)

	case ebcmABI.BytesArray.String():
		//return NewBytesArray(t.Data.([]uint64)...)

	case ebcmABI.Bytes32Array.String():
		return NewBytes32Array(t.Data.([]string)...)

		// case ebcmABI.ITuple.String():
		// 	for _, p := range t.Datas {
		// 		swapABIParams(p.(ebcmABI.Type))
		// 	}

		// 	return iTuple

		// case ebcmABI.ITupleFixedArray.String():
		// 	return iTupleFixedArray

		// case ebcmABI.ITupleFlexArray.String():
		// 	return iTupleFlexArray

	}
	return AbiParam{}
}

func EBCM_ABI_NewParams(ps ...ebcmABI.Type) AbiParams {
	list := AbiParams{}
	for _, p := range ps {
		// if elem, do := p.(ebcmABI.Type); do {
		// 	list = append(list, swapABIParams(elem))
		// } else {
		// 	list = append(list, p.(abiParam))
		// }
		list = append(list, swapABIParams(p))
	}
	return list
}
