package abmx

import (
	"errors"
	"math/big"
	"strings"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// //////////////////////////////////////////////////////
type fixedType struct {
	name string
}
type fixedTypeList []*fixedType

func (my *fixedType) String() string   { return my.name }
func (my *fixedType) Name() *fixedType { return my }

func (my *fixedType) GetArg() abi.Argument {
	re := abi.Argument{}
	tp := abi.Type{}
	tp, _ = abi.NewType(my.name, "", nil)
	re.Type = tp
	re.Indexed = false
	return re
}

// //////////////////////////////////////////////////////
type flexType struct {
	*fixedType
	datas fixedTypeList
}

////////////////////////////////////////////////////////

type Type interface {
	String() string
	Name() *fixedType
	GetArg() abi.Argument
}

type TypeList []Type

//func (my Type) String() string { return string(my) }

var (
	None = &fixedType{"none"}

	Address = &fixedType{"address"}
	Uint256 = &fixedType{"uint256"}
	Uint    = Uint256
	Uint128 = &fixedType{"uint128"}
	Uint112 = &fixedType{"uint112"}
	Uint64  = &fixedType{"uint64"}
	Uint32  = &fixedType{"uint32"}
	Uint16  = &fixedType{"uint16"}
	Uint8   = &fixedType{"uint8"}
	Bool    = &fixedType{"bool"}
	String  = &fixedType{"string"}
	Bytes   = &fixedType{"bytes"}
	Bytes32 = &fixedType{"bytes32"}

	AddressArray = &fixedType{"address[]"}
	Uint256Array = &fixedType{"uint256[]"}
	UintArray    = Uint256Array
	Uint128Array = &fixedType{"uint128[]"}
	Uint112Array = &fixedType{"uint112[]"}
	Uint64Array  = &fixedType{"uint64[]"}
	Uint32Array  = &fixedType{"uint32[]"}
	Uint16Array  = &fixedType{"uint16[]"}
	Uint8Array   = &fixedType{"uint8[]"}
	BoolArray    = &fixedType{"bool[]"}
	StringArray  = &fixedType{"string[]"}
	BytesArray   = &fixedType{"bytes[]"}
	Bytes32Array = &fixedType{"bytes32[]"}

	iTuple           = &fixedType{"tuple"}
	iTupleFlex       = &fixedType{"tuple_flex"}
	iTupleFixedArray = &fixedType{"tuple_fixed_array"} //고정길이 배열
	iTupleFlexArray  = &fixedType{"tuple_flex_array"}  //가변길이 배열

)

func Tuple(params ...*fixedType) *flexType {
	my := &flexType{
		fixedType: iTuple,
	}
	for _, v := range params {
		my.datas = append(my.datas, v)

		switch v {
		case Bytes,
			AddressArray,
			Uint256Array,
			UintArray,
			Uint128Array,
			Uint112Array,
			Uint64Array,
			Uint32Array,
			Uint16Array,
			Uint8Array,
			BoolArray,
			StringArray,
			BytesArray,
			Bytes32Array:
			my.fixedType = iTupleFlex
		}
	}
	return my
}

func TupleArray(params ...*fixedType) *flexType {
	my := &flexType{
		fixedType: iTupleFixedArray,
	}
	for _, v := range params {
		my.datas = append(my.datas, v)
		switch v {
		case String, Bytes,
			AddressArray,
			Uint256Array, Uint128Array, Uint112Array,
			Uint64Array, Uint32Array, Uint16Array, Uint8Array,
			BoolArray,
			StringArray, BytesArray,
			Bytes32Array:
			my.fixedType = iTupleFlexArray
		} //switch
	} //for
	return my
}

func (my abiReturns) IsTupleArray() {

}

/////////////////////////////////////////////////////////////

type AbiParam struct {
	p    Type
	data interface{}
}
type AbiParams []AbiParam

func (my AbiParam) getArg() abi.Argument { return my.p.GetArg() }
func (my AbiParam) getParam() interface{} {
	var re interface{}
	switch my.p.Name() {
	case Address:
		re = common.HexToAddress(my.data.(string))
	case Bytes32:
		re = common.HexToHash(my.data.(string))
		//re = jmath.New(my.data.(string)).ToBigInteger()

	case Uint256, Uint128, Uint112:
		re = jmath.New(my.data.(string)).ToBigInteger()
	case Uint64:
		re = my.data.(uint64)
	case Uint32:
		re = my.data.(uint32)
	case Uint16:
		re = my.data.(uint16)
	case Uint8:
		re = my.data.(uint8)

	case Bool:
		re = my.data.(bool)
	case String:
		re = my.data.(string)

	case AddressArray:
		list := []common.Address{}
		for _, v := range my.data.([]string) {
			list = append(list, common.HexToAddress(v))
		}
		re = list

	case Bytes32Array:
		list := []common.Hash{}
		for _, v := range my.data.([]string) {
			list = append(list, common.HexToHash(v))
		}
		re = list

	case Uint256Array, Uint128Array, Uint112Array:
		list := []*big.Int{}
		for _, v := range my.data.([]string) {
			list = append(list, jmath.New(v).ToBigInteger())
		}
		re = list
	case Bytes:
		re = my.data.([]byte)
	case BytesArray:
		re = my.data.([][]byte)

	case Uint64Array:
		re = my.data.([]uint64)
	case Uint32Array:
		re = my.data.([]uint32)
	case Uint16Array:
		re = my.data.([]uint16)
	case Uint8Array:
		re = my.data.([]uint8)

	case BoolArray:
		re = my.data.([]bool)
	case StringArray:
		re = my.data.([]string)
	}
	return re
}

func (my AbiParams) getArgument() []abi.Argument {
	args := []abi.Argument{}
	for _, v := range my {
		args = append(args, v.getArg())
	} //for
	return args
}
func (my AbiParams) getParames() []interface{} {
	var list []interface{}
	for _, v := range my {
		list = append(list, v.getParam())
	}
	return list
}

func NewAddress(data string) AbiParam {
	data = strings.TrimSpace(data)
	return AbiParam{Address, data}
}

func NewBytes32(data string) AbiParam {
	data = strings.TrimSpace(data)
	return AbiParam{Bytes32, data}
}
func NewBytes(data []byte) AbiParam {
	return AbiParam{Bytes, data}
}
func NewBytesArray(data [][]byte) AbiParam {
	return AbiParam{BytesArray, data}
}

func NewUint(data interface{}) AbiParam    { return NewUint256(data) }
func NewUint256(data interface{}) AbiParam { return AbiParam{Uint256, jmath.VALUE(data)} }
func NewUint128(data interface{}) AbiParam { return AbiParam{Uint128, jmath.VALUE(data)} }
func NewUint112(data interface{}) AbiParam { return AbiParam{Uint112, jmath.VALUE(data)} }
func NewUint64(data uint64) AbiParam       { return AbiParam{Uint64, data} }
func NewUint32(data uint32) AbiParam       { return AbiParam{Uint32, data} }
func NewUint16(data uint16) AbiParam       { return AbiParam{Uint16, data} }
func NewUint8(data uint8) AbiParam         { return AbiParam{Uint8, data} }
func NewBool(data bool) AbiParam           { return AbiParam{Bool, data} }
func NewString(data string) AbiParam       { return AbiParam{String, data} }

func NewAddressArray(data ...string) AbiParam {
	array := []string{}
	for _, v := range data {
		v = strings.TrimSpace(v)
		array = append(array, v)
	}
	return AbiParam{AddressArray, array}
}

func NewBytes32Array(data ...string) AbiParam {
	array := []string{}
	for _, v := range data {
		v = strings.TrimSpace(v)
		array = append(array, v)
	}
	return AbiParam{Bytes32Array, array}
}

func NewUint256Array(data ...interface{}) AbiParam {
	array := []string{}
	for _, v := range data {
		array = append(array, jmath.VALUE(v))
	}
	return AbiParam{Uint256Array, array}
}

func NewUint128Array(data ...interface{}) AbiParam {
	array := []string{}
	for _, v := range data {
		array = append(array, jmath.VALUE(v))
	}
	return AbiParam{Uint128Array, array}
}

func NewUint112Array(data ...interface{}) AbiParam {
	array := []string{}
	for _, v := range data {
		array = append(array, jmath.VALUE(v))
	}
	return AbiParam{Uint112Array, array}
}

func NewUint64Array(data ...uint64) AbiParam {
	array := []uint64{}
	for _, v := range data {
		array = append(array, v)
	}
	return AbiParam{Uint64Array, array}
}

func NewUint32Array(data ...uint32) AbiParam {
	array := []uint32{}
	for _, v := range data {
		array = append(array, v)
	}
	return AbiParam{Uint32Array, array}
}

func NewUint16Array(data ...uint16) AbiParam {
	array := []uint16{}
	for _, v := range data {
		array = append(array, v)
	}
	return AbiParam{Uint16Array, array}
}

func NewUint8Array(data ...uint8) AbiParam {
	array := []uint8{}
	for _, v := range data {
		array = append(array, v)
	}
	return AbiParam{Uint8Array, array}
}

func NewBoolArray(data ...bool) AbiParam {
	array := []bool{}
	for _, v := range data {
		array = append(array, v)
	}
	return AbiParam{BoolArray, array}
}

func NewStringArray(data ...string) AbiParam {
	array := []string{}
	for _, v := range data {
		array = append(array, v)
	}
	return AbiParam{StringArray, array}
}

func NewParams(ps ...AbiParam) AbiParams {
	list := AbiParams{}
	for _, p := range ps {
		list = append(list, p)
	}
	return list
}

/////////////////////////////////////////////////////////////

type abiReturns struct {
	args  AbiParams
	count int
	//realCount      int //rmv tuple
	// isTuple        bool
	// TupleArrayKind Type
	err error
}

func NewReturns(ps ...Type) abiReturns {
	r := abiReturns{
		//TupleArrayKind: None,
	}

	tupleCount := 0
	isTupleArray := false
	for _, v := range ps {
		r.args = append(r.args, AbiParam{p: v})

		switch v.Name() {
		case iTuple:
			tupleCount++
		case iTupleFlex, iTupleFixedArray, iTupleFlexArray:
			isTupleArray = true
		} //switch
	}

	if tupleCount > 1 || (tupleCount > 0 && isTupleArray) {
		errMsg := "TupleArray ---> dont used [Tuple] ==> do flat field."
		dbg.Red(errMsg)
		r.err = errors.New(errMsg)
	}

	r.count = len(r.args)
	// r.realCount = r.count
	// if isTuple {
	// 	r.realCount--
	// 	r.isTuple = true
	// }
	return r
}

func (my abiReturns) getArgument() []abi.Argument {
	return my.args.getArgument()
}

type Method struct {
	Name    string
	Params  AbiParams
	Returns abiReturns
}

// getMethod : inn Call
func (my Method) getMethod() map[string]abi.Method {
	// return map[string]abi.Method{
	// 	my.Name: {
	// 		my.Name,
	// 		my.Name,
	// 		false,
	// 		my.Params.getArgument(),
	// 		my.Returns.getArgument(),
	// 	},
	// }

	m := abi.NewMethod(
		my.Name,
		my.Name,
		abi.Function,
		"",
		false,
		false,
		my.Params.getArgument(),
		my.Returns.getArgument(),
	)
	return map[string]abi.Method{
		my.Name: m,
	}
}
