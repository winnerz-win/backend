package jmath

import (
	"math/rand"
	"time"
)

//Random :
type Random struct {
	rnd *rand.Rand
}

//NewRandom :
func NewRandom() *Random {
	obj := &Random{
		rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return obj
}

//NewRandomSeed :
func NewRandomSeed(seed int64) *Random {
	obj := &Random{
		rnd: rand.New(rand.NewSource(seed)),
	}
	return obj
}

//GetInt : [이거대신 Int쓰셍] min(>=0) <= n < max  사이의 값 : 예외 -> 1
func (my *Random) GetInt(min int, max int) int {
	if min < 0 || max <= 0 {
		return 1
	}
	pa := max - min - 1
	if pa <= 1 {
		return 1
	}

	v := my.rnd.Intn(pa)
	v += min + 1

	return v
}

func abs(v int) int {
	if v < 0 {
		return v * -1
	}
	return v
}
func nav(v int) int {
	return v * -1
}

//Int : a ~ b-1 , 사이값
func (my *Random) Int(min int, max int) int {
	if min > max {
		return 0
	}
	gap := 0
	fv := 1
	if min < 0 && max >= 0 { // - , +
		gap = (max - min)
	} else if min >= 0 && max >= 0 { // + , +
		gap = max - min
	} else { // - , -
		_max := abs(min)
		min = abs(max)
		max = _max
		gap = max - min
		fv = -1
	}

	if gap <= 0 {
		return 0
	}

	v := my.rnd.Intn(gap)

	return (v + min) * fv
}

//Int100 : 1 <= 100 까지 값
func (my *Random) Int100() int {

	v := my.Int(0, 100)
	return v + 1
}
