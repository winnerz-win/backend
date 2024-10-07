package abix

import (
	"strings"
	"txscheduler/brix/tools/dbg"
)

type ResultItem struct {
	//kind Type
	Name string
	Data interface{}
}

func (my ResultItem) Address() string { return my.Data.(string) }
func (my ResultItem) Uint256() string { return my.Data.(string) }
func (my ResultItem) Uint() string    { return my.Uint256() }
func (my ResultItem) Uint128() string { return my.Data.(string) }
func (my ResultItem) Uint64() uint64  { return my.Data.(uint64) }
func (my ResultItem) Uint32() uint32  { return my.Data.(uint32) }
func (my ResultItem) Uint16() uint16  { return my.Data.(uint16) }
func (my ResultItem) Uint8() uint8    { return my.Data.(uint8) }
func (my ResultItem) Bool() bool      { return my.Data.(bool) }
func (my ResultItem) Text() string    { return my.Data.(string) }
func (my ResultItem) Bytes() string   { return my.Data.(string) }
func (my ResultItem) Bytes32() string { return my.Data.(string) }

func (my ResultItem) AddressArray() []string { return my.Data.([]string) }
func (my ResultItem) Uint256Array() []string { return my.Data.([]string) }
func (my ResultItem) UintArray() []string    { return my.Uint256Array() }
func (my ResultItem) Uint128Array() []string { return my.Data.([]string) }
func (my ResultItem) Uint64Array() []uint64  { return my.Data.([]uint64) }
func (my ResultItem) Uint32Array() []uint32  { return my.Data.([]uint32) }
func (my ResultItem) Uint16Array() []uint16  { return my.Data.([]uint16) }
func (my ResultItem) Uint8Array() []uint8    { return my.Data.([]uint8) }
func (my ResultItem) BoolArray() []bool      { return my.Data.([]bool) }
func (my ResultItem) TextArray() []string    { return my.Data.([]string) }
func (my ResultItem) BytesArray() []string   { return my.Data.([]string) }
func (my ResultItem) Bytes32Array() []string { return my.Data.([]string) }

type IRESULT interface {
	Address(i int) string
	Uint256(i int) string
	Uint(i int) string
	Uint128(i int) string
	Uint64(i int) uint64
	Uint32(i int) uint32
	Uint16(i int) uint16
	Uint8(i int) uint8
	Bool(i int) bool
	Text(i int) string
	Bytes(i int) string
	Bytes32(i int) string

	AddressArray(i int) []string
	Uint256Array(i int) []string
	UintArray(i int) []string
	Uint128Array(i int) []string
	Uint64Array(i int) []uint64
	Uint32Array(i int) []uint32
	Uint16Array(i int) []uint16
	Uint8Array(i int) []uint8
	BoolArray(i int) []bool
	TextArray(i int) []string
	BytesArray(i int) []string
	Bytes32Array(i int) []string
}

func (my ResultItemList) Address(i int) string { return my[i].Address() }
func (my ResultItemList) Uint256(i int) string { return my[i].Uint256() }
func (my ResultItemList) Uint(i int) string    { return my.Uint256(i) }
func (my ResultItemList) Uint128(i int) string { return my[i].Uint128() }
func (my ResultItemList) Uint64(i int) uint64  { return my[i].Uint64() }
func (my ResultItemList) Uint32(i int) uint32  { return my[i].Uint32() }
func (my ResultItemList) Uint16(i int) uint16  { return my[i].Uint16() }
func (my ResultItemList) Uint8(i int) uint8    { return my[i].Uint8() }
func (my ResultItemList) Bool(i int) bool      { return my[i].Bool() }
func (my ResultItemList) Text(i int) string    { return my[i].Text() }
func (my ResultItemList) Bytes(i int) string   { return my[i].Bytes() }
func (my ResultItemList) Bytes32(i int) string { return my[i].Bytes32() }

func (my ResultItemList) AddressArray(i int) []string { return my[i].AddressArray() }
func (my ResultItemList) Uint256Array(i int) []string { return my[i].Uint256Array() }
func (my ResultItemList) UintArray(i int) []string    { return my.Uint256Array(i) }
func (my ResultItemList) Uint128Array(i int) []string { return my[i].Uint128Array() }
func (my ResultItemList) Uint64Array(i int) []uint64  { return my[i].Uint64Array() }
func (my ResultItemList) Uint32Array(i int) []uint32  { return my[i].Uint32Array() }
func (my ResultItemList) Uint16Array(i int) []uint16  { return my[i].Uint16Array() }
func (my ResultItemList) Uint8Array(i int) []uint8    { return my[i].Uint8Array() }
func (my ResultItemList) BoolArray(i int) []bool      { return my[i].BoolArray() }
func (my ResultItemList) TextArray(i int) []string    { return my[i].TextArray() }
func (my ResultItemList) BytesArray(i int) []string   { return my[i].BytesArray() }
func (my ResultItemList) Bytes32Array(i int) []string { return my[i].Bytes32Array() }

const NAN = ""

func (my RESULT) checkSize(i int) bool { return i < len(my.items) }

func (my RESULT) Address(i int) string {
	if !my.checkSize(i) {
		return NAN
	}
	return my.items[i].Address()
}
func (my RESULT) Uint256(i int) string {
	if !my.checkSize(i) {
		return NAN
	}
	return my.items[i].Uint256()
}
func (my RESULT) Uint(i int) string {
	if !my.checkSize(i) {
		return NAN
	}
	return my.Uint256(i)
}
func (my RESULT) Uint128(i int) string {
	if !my.checkSize(i) {
		return NAN
	}
	return my.items[i].Uint128()
}
func (my RESULT) Uint64(i int) uint64 {
	if !my.checkSize(i) {
		return 0
	}
	return my.items[i].Uint64()
}
func (my RESULT) Uint32(i int) uint32 {
	if !my.checkSize(i) {
		return 0
	}
	return my.items[i].Uint32()
}
func (my RESULT) Uint16(i int) uint16 {
	if !my.checkSize(i) {
		return 0
	}
	return my.items[i].Uint16()
}
func (my RESULT) Uint8(i int) uint8 {
	if !my.checkSize(i) {
		return 0
	}
	return my.items[i].Uint8()
}
func (my RESULT) Bool(i int) bool {
	if !my.checkSize(i) {
		return false
	}
	return my.items[i].Bool()
}
func (my RESULT) Text(i int) string {
	if !my.checkSize(i) {
		return NAN
	}
	return my.items[i].Text()
}
func (my RESULT) Bytes(i int) string {
	if !my.checkSize(i) {
		return NAN
	}
	return my.items[i].Bytes()
}
func (my RESULT) Bytes32(i int) string {
	if !my.checkSize(i) {
		return NAN
	}
	return my.items[i].Bytes32()
}

func (my RESULT) AddressArray(i int) []string { return my.items[i].AddressArray() }
func (my RESULT) Uint256Array(i int) []string { return my.items[i].Uint256Array() }
func (my RESULT) UintArray(i int) []string    { return my.Uint256Array(i) }
func (my RESULT) Uint128Array(i int) []string { return my.items[i].Uint128Array() }
func (my RESULT) Uint64Array(i int) []uint64  { return my.items[i].Uint64Array() }
func (my RESULT) Uint32Array(i int) []uint32  { return my.items[i].Uint32Array() }
func (my RESULT) Uint16Array(i int) []uint16  { return my.items[i].Uint16Array() }
func (my RESULT) Uint8Array(i int) []uint8    { return my.items[i].Uint8Array() }
func (my RESULT) BoolArray(i int) []bool      { return my.items[i].BoolArray() }
func (my RESULT) TextArray(i int) []string    { return my.items[i].TextArray() }
func (my RESULT) BytesArray(i int) []string   { return my.items[i].BytesArray() }
func (my RESULT) Bytes32Array(i int) []string { return my.items[i].Bytes32Array() }

type TupleList []ResultItemList

func (my TupleList) String() string { return dbg.ToJSONString(my) }
func (my TupleList) Count() int     { return len(my) }

func (my RESULT) Tuple(i int) TupleList {
	cnt := i + 1
	if len(my.items) < cnt {
		return TupleList{}
	}
	target := my.items[i]
	if target.isTuple() {
		return target.Data.(TupleList)
	}
	return TupleList{}
}
func (my RESULT) TupleOne(i int) ResultItemList {
	tuple := my.Tuple(i)
	if len(tuple) > 0 {
		return tuple[0]
	}
	return ResultItemList{}
}

func (my ResultItem) isTuple() bool {
	tname := dbg.TrimToLower(my.Name)
	return strings.Contains(tname, "tuple")
}

type ResultItemList []ResultItem

func (my ResultItem) String() string     { return dbg.ToJSONString(my) }
func (my ResultItemList) String() string { return dbg.ToJSONString(my) }

// RESULT :
type RESULT struct {
	items   ResultItemList
	IsError bool
}

func (my RESULT) GetItems() ResultItemList {
	return my.items
}

func NewRESULT() RESULT {
	return RESULT{
		items:   ResultItemList{},
		IsError: true,
	}
}

func (my *RESULT) Inject(list ResultItemList) {
	my.items = list
}

func (my TupleList) helpString(i int, gap string) string {
	msg := dbg.Cat(gap, "[", i, "] TupleArray {", dbg.ENTER)

	ngap := gap + "   "
	for j, v := range my {
		val := dbg.Cat(ngap, "[", j, "] ", v.helpString(ngap+"   "))
		msg += val
		//dbg.Purple(val)
	} //for
	msg += dbg.Cat(gap, "}")
	return msg
}

func (my ResultItem) helpString(i int, gap string) string {
	if my.isTuple() == false {
		val := dbg.Cat(gap, "[", i, "] ", my.Name, " : ", my.Data)
		//dbg.Purple(val)
		return val
	}
	tuple := my.Data.(TupleList)
	return tuple.helpString(i, gap)
}

func (my ResultItemList) helpString(gap string) string {
	msg := dbg.Cat(gap, "[", dbg.ENTER)
	for i, v := range my {
		val := dbg.Cat(v.helpString(i, gap), dbg.ENTER)
		msg += val
	}
	msg += dbg.Cat(gap, "]", dbg.ENTER)
	return msg
}

func (my RESULT) String() string {
	msg := "<RESULT> " + dbg.ENTER
	if my.IsError {
		msg += "ERROR <FAIL>" + dbg.ENTER
	}
	for i, item := range my.items {
		msg += item.helpString(i, "") + dbg.ENTER

	}
	return msg
}

func (my *ResultItemList) EBCMadd(n string, val interface{}) {
	item := ResultItem{
		Name: n,
		Data: val,
	}
	(*my) = append((*my), item)
}

func (my *RESULT) EBCMaddItem(n string, val interface{}) {
	my.items.EBCMadd(n, val)
}
