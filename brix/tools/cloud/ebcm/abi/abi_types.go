package abi

import (
	"encoding/hex"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"

	"strings"
)

type Type struct {
	name  string
	Data  interface{}
	Datas []interface{}
}
type TypeList []Type

func (my Type) FlexTypeList() TypeList {
	list := TypeList{}
	for _, v := range my.Datas {
		list = append(list, v.(Type))
	}
	return list
}

func (my Type) String() string { return my.name }
func (my Type) clone(v interface{}) Type {
	return Type{
		name:  my.name,
		Data:  v,
		Datas: []interface{}{},
	}
}
func (my Type) makeTuple() Type {
	return Type{
		name:  my.name,
		Data:  nil,
		Datas: []interface{}{},
	}
}
func (my *Type) addTupleItem(v Type) {
	my.Datas = append(my.Datas, v)
}

// HasTuple : 0 , 1:tuple , 2:tuple-array
func (my Type) HasTuple() int {
	re := 0
	if strings.Contains(my.name, "tuple") {
		re = 1
		if strings.Contains(my.name, "array") {
			re = 2
		}
	}
	return re
}
func (my Type) TupleTypes() TypeList {
	list := TypeList{}
	for _, v := range my.Datas {
		list = append(list, v.(Type))
	}
	return list
}

func MakeTypeList(any_types ...interface{}) TypeList {
	list := TypeList{}
	list.Append(any_types...)
	return list
}
func (my *TypeList) Append(any_types ...interface{}) {
	for _, c := range any_types {
		switch v := c.(type) {
		case Type:
			*my = append(*my, v)
		case TypeList:
			*my = append(*my, v...)
		default:
			panic(c)
		}
	}
}

var (
	None    = Type{name: "none"}
	Address = Type{name: "address"}
	Uint256 = Type{name: "uint256"}
	Uint    = Uint256
	Uint128 = Type{name: "uint128"}
	Uint112 = Type{name: "uint112"}
	Uint64  = Type{name: "uint64"}
	Uint32  = Type{name: "uint32"}
	Uint16  = Type{name: "uint16"}
	Uint8   = Type{name: "uint8"}
	Bool    = Type{name: "bool"}
	String  = Type{name: "string"}
	Bytes   = Type{name: "bytes"}
	Bytes32 = Type{name: "bytes32"}

	AddressArray = Type{name: "address[]"}
	Uint256Array = Type{name: "uint256[]"}
	UintArray    = Uint256Array
	Uint128Array = Type{name: "uint128[]"}
	Uint112Array = Type{name: "uint112[]"}
	Uint64Array  = Type{name: "uint64[]"}
	Uint32Array  = Type{name: "uint32[]"}
	Uint16Array  = Type{name: "uint16[]"}
	Uint8Array   = Type{name: "uint8[]"}
	BoolArray    = Type{name: "bool[]"}
	StringArray  = Type{name: "string[]"}
	BytesArray   = Type{name: "bytes[]"}
	Bytes32Array = Type{name: "bytes32[]"}

	ITuple          = Type{name: "tuple"}
	ITupleFlex      = Type{name: "tuple_flex"}
	ITupleArray     = Type{name: "tuple_array"}
	ITupleArrayFlex = Type{name: "tuple_array_flex"}
)

func NewReturns(sl ...Type) TypeList {
	list := TypeList{}
	for _, v := range sl {
		list = append(list, v)
	}
	return list
}

///////////////////////////////////////////////////////////

func isFlexType(v Type) bool {
	switch v.String() {
	case String.String(),
		Bytes.String(),
		AddressArray.String(),
		Uint256Array.String(),
		Uint128Array.String(),
		Uint112Array.String(),
		Uint64Array.String(),
		Uint32Array.String(),
		Uint16Array.String(),
		Uint8Array.String(),
		BoolArray.String(),
		StringArray.String(),
		BytesArray.String(),
		Bytes32Array.String(),

		ITupleFlex.String(),
		ITupleArray.String(),
		ITupleArrayFlex.String():
		return true
	}
	return false
}

func Tuple(params ...Type) Type {
	my := ITuple.clone(nil)
	for _, v := range params {
		my.Datas = append(my.Datas, v)
		if isFlexType(v) {
			my.name = ITupleFlex.name
		}
	} //for
	return my
}

func TupleArray(params ...Type) Type {
	my := ITupleArray.clone(nil)
	for _, v := range params {
		my.Datas = append(my.Datas, v)
		if isFlexType(v) {
			my.name = ITupleArrayFlex.name
		}
	} //for
	return my
}

func NewTuple(fields ...Type) Type {
	tuple := ITuple.makeTuple()
	for _, v := range fields {
		tuple.addTupleItem(v)
		if isFlexType(v) {
			tuple.name = ITupleFlex.name
		}
	}
	tuple.Data = tuple.Datas
	return tuple
}

func NewTupleArray(params ...Type) Type {
	my := ITupleArray.clone(nil)
	for _, v := range params {
		my.Datas = append(my.Datas, v)
		if isFlexType(v) {
			my.name = ITupleArrayFlex.name
		}
	} //for
	return my
}

func NewAddress(data string) Type {
	data = strings.TrimSpace(data)
	return Address.clone(EIP55(data))
}

func NewBytes32(data string) Type {
	data = strings.TrimSpace(data)
	return Bytes32.clone(data)
}

func HexToBytes(hexBytes string) []byte {
	hexBytes = strings.TrimPrefix(hexBytes, "0x")
	buf, err := hex.DecodeString(hexBytes)
	if err != nil {
		dbg.RedItalic("abi.HexToBytes :", err)
		return nil
	}
	return buf
}

func NewBytes(data []byte) Type {
	return Bytes.clone(data)
}

func HexArrayToBytes(array ...string) [][]byte {
	if len(array) == 0 {
		return nil
	}
	bufs := make([][]byte, len(array))
	for i, hexBytes := range array {
		bufs[i] = HexToBytes(hexBytes)
	}
	return bufs
}

func NewBytesArray(data [][]byte) Type {
	return BytesArray.clone(data)
}

func NewUint(data interface{}) Type    { return NewUint256(data) }
func NewUint256(data interface{}) Type { return Uint256.clone(jmath.VALUE(data)) }
func NewUint128(data interface{}) Type { return Uint128.clone(jmath.VALUE(data)) }
func NewUint112(data interface{}) Type { return Uint112.clone(jmath.VALUE(data)) }
func NewUint64(data uint64) Type       { return Uint64.clone(data) }
func NewUint32(data uint32) Type       { return Uint32.clone(data) }
func NewUint16(data uint16) Type       { return Uint16.clone(data) }
func NewUint8(data uint8) Type         { return Uint8.clone(data) }
func NewBool(data bool) Type           { return Bool.clone(data) }
func NewString(data string) Type       { return String.clone(data) }

func NewAddressArray(data ...string) Type {
	array := []string{}
	for _, v := range data {
		v = strings.TrimSpace(v)
		array = append(array, v)
	}
	return AddressArray.clone(array)
}

func NewBytes32Array(data ...string) Type {
	array := []string{}
	for _, v := range data {
		v = strings.TrimSpace(v)
		array = append(array, v)
	}
	return Bytes32Array.clone(array)
}

func checkArrayString(data ...interface{}) ([]string, bool) {
	if len(data) == 1 {
		switch void := data[0].(type) {
		case []int8:
			array := []string{}
			for _, v := range void {
				array = append(array, jmath.VALUE(v))
			}
			return array, true
		case []uint8:
			array := []string{}
			for _, v := range void {
				array = append(array, jmath.VALUE(v))
			}
			return array, true
		case []int:
			array := []string{}
			for _, v := range void {
				array = append(array, jmath.VALUE(v))
			}
			return array, true
		case []uint:
			array := []string{}
			for _, v := range void {
				array = append(array, jmath.VALUE(v))
			}
			return array, true
		case []int16:
			array := []string{}
			for _, v := range void {
				array = append(array, jmath.VALUE(v))
			}
			return array, true
		case []uint16:
			array := []string{}
			for _, v := range void {
				array = append(array, jmath.VALUE(v))
			}
			return array, true
		case []int32:
			array := []string{}
			for _, v := range void {
				array = append(array, jmath.VALUE(v))
			}
			return array, true
		case []uint32:
			array := []string{}
			for _, v := range void {
				array = append(array, jmath.VALUE(v))
			}
			return array, true
		case []int64:
			array := []string{}
			for _, v := range void {
				array = append(array, jmath.VALUE(v))
			}
			return array, true
		case []uint64:
			array := []string{}
			for _, v := range void {
				array = append(array, jmath.VALUE(v))
			}
			return array, true

		case []string:
			array := []string{}
			for _, v := range void {
				array = append(array, jmath.VALUE(v))
			}
			return array, true
		}
	}
	return nil, false
}

func NewUint256Array(data ...interface{}) Type {
	if array, ok := checkArrayString(data...); ok {
		return Uint256Array.clone(array)
	}

	array := []string{}

	for _, v := range data {
		array = append(array, jmath.VALUE(v))
	}
	return Uint256Array.clone(array)
}

func NewUint128Array(data ...interface{}) Type {
	if array, ok := checkArrayString(data...); ok {
		return Uint128Array.clone(array)
	}

	array := []string{}
	for _, v := range data {
		array = append(array, jmath.VALUE(v))
	}
	return Uint128Array.clone(array)
}

func NewUint112Array(data ...interface{}) Type {
	if array, ok := checkArrayString(data...); ok {
		return Uint112Array.clone(array)
	}

	array := []string{}
	for _, v := range data {
		array = append(array, jmath.VALUE(v))
	}
	return Uint112Array.clone(array)
}

func NewUint64Array(data ...uint64) Type {
	array := []uint64{}
	array = append(array, data...)
	return Uint64Array.clone(array)
}

func NewUint32Array(data ...uint32) Type {
	array := []uint32{}
	array = append(array, data...)
	return Uint32Array.clone(array)
}

func NewUint16Array(data ...uint16) Type {
	array := []uint16{}
	array = append(array, data...)
	return Uint16Array.clone(array)
}

func NewUint8Array(data ...uint8) Type {
	array := []uint8{}
	array = append(array, data...)
	return Uint8Array.clone(array)
}

func NewBoolArray(data ...bool) Type {
	array := []bool{}
	array = append(array, data...)
	return BoolArray.clone(array)
}

func NewStringArray(data ...string) Type {
	array := []string{}
	array = append(array, data...)
	return StringArray.clone(array)
}

func NewParams(sl ...Type) TypeList {
	list := TypeList{}
	for _, v := range sl {
		list = append(list, v)
	}
	return list
}
