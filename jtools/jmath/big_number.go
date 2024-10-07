package jmath

import (
	"errors"
	"fmt"
	"jtools/jmath/decimal"
	"jtools/unix"
	"math/big"
	"strconv"
	"strings"
)

type bigNumber struct {
	data decimal.Decimal
}

func newBigNumber2(v interface{}) (*bigNumber, error) {

	value := "0"
	if v == nil {
		data, _ := decimal.NewFromString(value)
		return &bigNumber{data}, errors.New("v is null")
	}

	switch n := v.(type) {
	case *bigNumber:
		return &bigNumber{
			data: n.data,
		}, nil

	case bigNumber:
		return &bigNumber{
			data: n.data,
		}, nil

	case [32]byte:
		sl := []byte{}
		sl = append(sl, n[:]...)
		value = big.NewInt(0).SetBytes(sl).String()

	case []byte:
		value = big.NewInt(0).SetBytes(n).String()

	case decimal.Decimal:
		value = strings.TrimSpace(n.String())

	case unix.Time:
		value = n.Int64String()

	default:
		value = fmt.Sprintf("%v", v)

	} //switch

	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "0x") {
		cv, do := big.NewInt(0).SetString(value[2:], 16)
		if do {
			value = cv.String()
		}
	}

	data, err := decimal.NewFromString(value)
	number := &bigNumber{
		data: data,
	}
	return number, err
}

func newBigNumber(v interface{}) *bigNumber {
	n, _ := newBigNumber2(v)
	return n
}
func iNumber(v decimal.Decimal) *bigNumber {
	return &bigNumber{v}
}

func (my bigNumber) String() string { return my.data.String() }

func (my bigNumber) clone() *bigNumber {
	return newBigNumber(my)
}

func (my bigNumber) BigInt() *big.Int {
	v, _ := big.NewInt(0).SetString(my.data.StringFixed(0), 10)
	return v
}

func (my bigNumber) Uint64() uint64 { return my.BigInt().Uint64() }
func (my bigNumber) Uint32() uint32 { return uint32(my.BigInt().Uint64()) }
func (my bigNumber) Int64() int64   { return int64(my.BigInt().Int64()) }
func (my bigNumber) Int() int       { return int(my.BigInt().Int64()) }

func (my bigNumber) Float64() float64 {
	f, _ := strconv.ParseFloat(my.String(), 64)
	return f
}

func (my bigNumber) MovePointRight(n interface{}) Number {
	return iNumber(my.data.Shift(int32(newBigNumber(n).Int())))
}

func (my bigNumber) Pow(n interface{}) Number {
	return iNumber(my.data.Pow(newBigNumber(n).data))
}

func (my bigNumber) Mul(n interface{}) Number {
	return iNumber(my.data.Mul(newBigNumber(n).data))
}

func (my bigNumber) Div(n interface{}) Number {
	return iNumber(my.data.Div(newBigNumber(n).data))
}

func (my bigNumber) Cmp(n interface{}) int {
	return my.data.Cmp(newBigNumber(n).data)
}

func (my bigNumber) Sub(n interface{}) Number {
	return iNumber(my.data.Sub(newBigNumber(n).data))
}

func (my bigNumber) Add(n interface{}) Number {
	return iNumber(my.data.Add(newBigNumber(n).data))
}

// Round :소숫점 자리 반올림
func (my bigNumber) Round(place interface{}) Number {
	return iNumber(my.data.Round(int32(newBigNumber(place).Uint32())))
}

// Up : 자릿수 무조껀 올림
func (my bigNumber) Up(place interface{}) Number {
	return iNumber(my.data.Up(int32(newBigNumber(place).Uint32())))
}

// Floor : 소수점 버림.
func (my bigNumber) Floor() Number {
	return iNumber(my.data.Floor())
}
