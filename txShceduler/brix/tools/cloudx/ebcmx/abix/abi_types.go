package abix

import (
	"encoding/hex"
	"strings"
	"txscheduler/brix/tools/dbg/cc"
	"txscheduler/brix/tools/jmath"
)

type Type struct {
	name  string
	Data  interface{}
	Datas []interface{}
}
type TypeList []Type

func (my Type) String() string { return my.name }
func (my Type) clone(v interface{}) Type {
	return Type{
		name:  my.name,
		Data:  v,
		Datas: []interface{}{},
	}
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

	ITuple           = Type{name: "tuple"}
	ITupleFixedArray = Type{name: "tuple_fixed_array"} //고정길이 배열
	ITupleFlexArray  = Type{name: "tuple_flex_array"}  //가변길이 배열
)

func NewReturns(sl ...Type) TypeList {
	list := TypeList{}
	for _, v := range sl {
		list = append(list, v)
	}
	return list
}

///////////////////////////////////////////////////////////

func Tuple(params ...Type) Type {
	my := ITuple.clone(nil)
	for _, v := range params {
		my.Datas = append(my.Datas, v)
	}
	return my
}

func TupleArray(params ...Type) Type {
	my := ITupleFixedArray.clone(nil)
	for _, v := range params {
		my.Datas = append(my.Datas, v)
		switch v.String() {
		case String.String(), Bytes.String(),
			AddressArray.String(),
			Uint256Array.String(), Uint128Array.String(), Uint112Array.String(),
			Uint64Array.String(), Uint32Array.String(), Uint16Array.String(), Uint8Array.String(),
			BoolArray.String(),
			StringArray.String(), BytesArray.String(),
			Bytes32Array.String():
			my.name = ITupleFlexArray.name
		} //switch
	} //for
	return my
}

func NewAddress(data string) Type {
	data = strings.TrimSpace(data)
	return Address.clone(data)
}

func NewBytes32(data string) Type {
	data = strings.TrimSpace(data)
	return Bytes32.clone(data)
}

func HexToBytes(hexBytes string) []byte {
	hexBytes = strings.TrimPrefix(hexBytes, "0x")
	buf, err := hex.DecodeString(hexBytes)
	if err != nil {
		cc.PrintRed("abi.HexToBytes :", err)
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

func NewUint256Array(data ...interface{}) Type {
	array := []string{}
	for _, v := range data {
		array = append(array, jmath.VALUE(v))
	}
	return Uint256Array.clone(array)
}

func NewUint128Array(data ...interface{}) Type {
	array := []string{}
	for _, v := range data {
		array = append(array, jmath.VALUE(v))
	}
	return Uint128Array.clone(array)
}

func NewUint112Array(data ...interface{}) Type {
	array := []string{}
	for _, v := range data {
		array = append(array, jmath.VALUE(v))
	}
	return Uint112Array.clone(array)
}

func NewUint64Array(data ...uint64) Type {
	array := []uint64{}
	for _, v := range data {
		array = append(array, v)
	}
	return Uint64Array.clone(array)
}

func NewUint32Array(data ...uint32) Type {
	array := []uint32{}
	for _, v := range data {
		array = append(array, v)
	}
	return Uint32Array.clone(array)
}

func NewUint16Array(data ...uint16) Type {
	array := []uint16{}
	for _, v := range data {
		array = append(array, v)
	}
	return Uint16Array.clone(array)
}

func NewUint8Array(data ...uint8) Type {
	array := []uint8{}
	for _, v := range data {
		array = append(array, v)
	}
	return Uint8Array.clone(array)
}

func NewBoolArray(data ...bool) Type {
	array := []bool{}
	for _, v := range data {
		array = append(array, v)
	}
	return BoolArray.clone(array)
}

func NewStringArray(data ...string) Type {
	array := []string{}
	for _, v := range data {
		array = append(array, v)
	}
	return StringArray.clone(array)
}

func NewParams(sl ...Type) TypeList {
	list := TypeList{}
	for _, v := range sl {
		list = append(list, v)
	}
	return list
}
