package jmath

import (
	"math/big"

	"txscheduler/brix/tools/jmath/decimal"
)

// BigDecimal :
type BigDecimal struct {
	data  decimal.Decimal
	_sVal string
}

// NewBigDecimal :
func NewBigDecimal(v interface{}, isErr ...*error) *BigDecimal {
	ins := &BigDecimal{
		data: decimal.NewDecimal(v, isErr...),
	}
	ins.ToString()
	return ins
}

// New :
func New(v interface{}) *BigDecimal {
	return NewBigDecimal(v)
}

func NEW(v interface{}) *BigDecimal {
	return New(v)
}

func (my *BigDecimal) BigInt() *big.Int {
	return my.ToBigInteger()
}

// Uint64 :
func (my *BigDecimal) Uint64() uint64 {
	return my.ToBigInteger().Uint64()
}

// Int64 :
func (my *BigDecimal) Int64() int64 {
	return int64(my.Uint64())
}

// Int :
func (my *BigDecimal) Int() int {
	return int(my.Uint64())
}

func BigInt(v interface{}) *big.Int {
	return New(v).ToBigInteger()
}

// IsNum :
func IsNum(v interface{}) bool {
	var err error
	decimal.NewDecimal(v, &err)
	return err == nil
}

// IsUnderZero : 0이하 true
func IsUnderZero(v interface{}) bool {
	if IsNum(v) == false {
		return true
	}
	return CMP(v, 0) <= 0
}

// IsLessZero : 0미만 true
func IsLessZero(v interface{}) bool {
	if IsNum(v) == false {
		return true
	}
	return CMP(v, 0) < 0
}

// IsZERO :
func IsZERO(v interface{}) bool {
	return VALUE(v) == "0"
}

func (my *BigDecimal) ToIDecimal() string {
	return my.data.String()
}

func (my *BigDecimal) Clone() *BigDecimal {
	return NewBigDecimal(my.data)
}

func (my *BigDecimal) MovePointRight(n interface{}) *BigDecimal {
	v := GetN(n)
	my.data = my.data.Shift(int32(v))
	my.ToString()
	return my
}

// Pow : int , zava.INT ,
func (my *BigDecimal) Pow(nn interface{}) *BigDecimal {
	pv := decimal.NewDecimal(nn)
	my.data = my.data.Pow(pv)
	my.ToString()
	return my
}
func (my *BigDecimal) Multiply(v *BigDecimal) *BigDecimal {
	my.data = my.data.Mul(v.data)
	my.ToString()
	return my
}
func (my *BigDecimal) MultiplyC(v *BigDecimal) *BigDecimal {
	return my.Clone().Multiply(v)
}
func (my *BigDecimal) Multiplyv(v interface{}) *BigDecimal {
	my.data = my.data.Mul(NewBigDecimal(v).data)
	my.ToString()
	return my
}

func (my *BigDecimal) Divide(v *BigDecimal) *BigDecimal {
	my.data = my.data.Div(v.data)
	my.ToString()
	return my
}
func (my *BigDecimal) DivideC(v *BigDecimal) *BigDecimal {
	return my.Clone().Divide(v)
}
func (my *BigDecimal) Dividev(v interface{}) *BigDecimal {
	my.data = my.data.Div(NewBigDecimal(v).data)
	my.ToString()
	return my
}

func (my *BigDecimal) ToString() string {
	my._sVal = my.data.String()
	return my._sVal
}
func (my BigDecimal) String() string {
	return my.ToString()
}

func (my *BigDecimal) CompareTo(v *BigDecimal) int {
	return my.data.Cmp(v.data)
}
func (my *BigDecimal) CompareTov(v interface{}) int {
	return my.data.Cmp(NewBigDecimal(v).data)
}

func (my *BigDecimal) Subtract(v *BigDecimal) *BigDecimal {
	my.data = my.data.Sub(v.data)
	my.ToString()
	return my
}
func (my *BigDecimal) SubtractC(v *BigDecimal) *BigDecimal {
	return my.Clone().Subtract(v)
}
func (my *BigDecimal) Subtractv(v interface{}) *BigDecimal {
	my.data = my.data.Sub(NewBigDecimal(v).data)
	my.ToString()
	return my
}

func (my *BigDecimal) Add(v *BigDecimal) *BigDecimal {
	my.data = my.data.Add(v.data)
	my.ToString()
	return my
}
func (my *BigDecimal) AddC(v *BigDecimal) *BigDecimal {
	return my.Clone().Add(v)
}
func (my *BigDecimal) Addv(v interface{}) *BigDecimal {
	my.data = my.data.Add(NewBigDecimal(v).data)
	my.ToString()
	return my
}

func (my *BigDecimal) ToBigInteger() *big.Int {
	v := big.NewInt(0)
	v, _ = v.SetString(my.data.StringFixed(0), 10)
	return v
}

func (my *BigDecimal) LongValue() uint64 {
	v := my.ToBigInteger()
	return v.Uint64()
}

// Round :소숫점 자리 반올림
func (my *BigDecimal) Round(place int32) *BigDecimal {
	my.data = my.data.Round(place)
	my.ToString()
	return my
}

// RoundC :소숫점 자리 반올림
func (my *BigDecimal) RoundC(place int32) *BigDecimal {
	return my.Clone().Round(place)
}

// Up : 자릿수 무조껀 올림
func (my *BigDecimal) Up(place int32) *BigDecimal {
	my.data = my.data.Up(place)
	my.ToString()
	return my
}

// UpC : 자릿수 무조껀 올림
func (my *BigDecimal) UpC(place int32) *BigDecimal {
	return my.Clone().Up(place)
}

// Floor : 소수점 버림.
func (my *BigDecimal) Floor() *BigDecimal {
	my.data = my.data.Floor()
	my.ToString()
	return my
}

// FloorC : 소수점 버림.
func (my *BigDecimal) FloorC() *BigDecimal {
	return my.Clone().Floor()
}
