package jmath

import (
	"math/big"
)

const (
	//Int64Max :
	Int64Max = 9223372036854775807
	//Int32Max :
	Int32Max = 2147483647
	//Uint64Max :
	Uint64Max = ^uint64(0)

	//Uint32Max :
	Uint32Max = ^uint32(0)
)

type Number interface {
	String() string
	clone() *bigNumber
	BigInt() *big.Int
	Uint64() uint64
	Uint32() uint32
	Int64() int64
	Int() int
	Float64() float64
	MovePointRight(n interface{}) Number
	Pow(n interface{}) Number
	Mul(n interface{}) Number
	Div(n interface{}) Number
	Cmp(n interface{}) int
	Sub(n interface{}) Number
	Add(n interface{}) Number
	Round(place interface{}) Number
	Up(place interface{}) Number
	Floor() Number
}

func NEW(v interface{}) Number {
	n, _ := newBigNumber2(v)
	return n
}

func NEW2(v interface{}, isErr ...*error) Number {
	n, err := newBigNumber2(v)
	if len(isErr) > 0 {
		*isErr[0] = err
	}
	return n
}
