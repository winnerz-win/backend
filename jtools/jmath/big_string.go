package jmath

import (
	"encoding/hex"
	"fmt"
	"math"
	"strings"
)

// Int64MaxString :
func Int64MaxString() string {
	return fmt.Sprintf("%v", Int64Max)
}

// Int32MaxString :
func Int32MaxString() string {
	return fmt.Sprintf("%v", Int32Max)
}

// Uint64MaxString :
func Uint64MaxString() string {
	return fmt.Sprintf("%v", Uint64Max)
}

// Uint32MaxString :
func Uint32MaxString() string {
	return fmt.Sprintf("%v", Uint32Max)
}

func VALUE(i interface{}) string { return NEW(i).String() }

func IsFloat(i interface{}) bool { return strings.Contains(VALUE(i), ".") }

func BYTES(i interface{}) []byte { return NEW(i).BigInt().Bytes() }

func HEX(i interface{}, withdout0x ...bool) string {
	b := BYTES(i)
	v := hex.EncodeToString(b)
	if len(withdout0x) > 0 && withdout0x[0] {
		return v
	}
	v = strings.ToLower(v)

	return "0x" + v
}

func ADD(array ...interface{}) string {
	if len(array) == 0 {
		return "0"
	}
	sum := NEW(array[0])
	for i := 1; i < len(array); i++ {
		sum = sum.Add(array[i])
	}
	return sum.String()
}
func SUB(i, j interface{}) string {
	return NEW(i).Sub(j).String()
}

func MUL(i, j interface{}) string {
	return NEW(i).Mul(j).String()
}

func DIV(i, j interface{}) string {
	return NEW(i).Div(j).String()
}
func DIVDOT(i, j interface{}) string {
	val := DIV(i, j)
	return DOTCUT(val, 0)
}

func CMP(src, tar interface{}) int {
	return NEW(src).Cmp(tar)
}

func CMPSUB(src, tar interface{}) (int, string) {
	subVal := "0"
	cmp := CMP(src, tar)
	switch cmp {
	case -1:
		subVal = SUB(tar, src)
	case 1:
		subVal = SUB(src, tar)
	}
	return cmp, subVal

}

func ROUND(v, place interface{}) string {
	return NEW(v).Round(place).String()
}

func UP(v, place interface{}) string {
	return NEW(v).Up(place).String()
}

func FLOOR(v interface{}) string {
	return NEW(v).Floor().String()
}

func POW(x, n interface{}) string {
	return NEW(x).Pow(n).String()
}
func POWInt(x, n interface{}) string {
	if CMP(n, 0) == 0 {
		return "1"
	}
	r := VALUE(x)
	p := r
	loop := NEW(x).Int()
	for loop > 1 {
		r = MUL(r, p)
		loop--
	} //for
	return r
}

func DOTCUT(v, place interface{}) string {
	p := NEW(place).Int()
	val := NEW(v)
	if p <= 0 {
		return val.Floor().String()
	}

	a := val.String()
	ss := strings.Split(a, ".")
	if len(ss) == 1 {
		return a
	}
	dotValue := ss[1]
	if len(dotValue) <= p {
		return a
	}
	dot := dotValue[:p]
	return ss[0] + "." + dot
}

func MIN(x, y interface{}) string {
	if CMP(x, y) < 0 {
		return VALUE(x)
	}
	return VALUE(y)
}

func MAX(x, y interface{}) string {
	if CMP(x, y) < 0 {
		return VALUE(y)
	}
	return VALUE(x)
}

func ABS(v interface{}) string {
	if CMP(v, 0) < 0 {
		return MUL(v, -1)
	}
	return VALUE(v)
}

func MOD(a, b interface{}) string {
	v := DIV(a, b)
	c := DOTCUT(v, 0)
	r := SUB(a, MUL(c, b))
	return VALUE(r)
}

func SQRT_INT(v interface{}) string {
	// y := Uint64(v)
	// z := uint64(0)
	// if y > 3 {
	// 	z = y
	// 	x := y/2 + 1
	// 	for x < z {
	// 		z = x
	// 		x = (y/x + x) / 2
	// 	}
	// } else if y != 0 {
	// 	z = 1
	// }
	// return z
	y := Float64(v)
	r := math.Sqrt(y)
	return DOTCUT(r, 0)
}

func SQRT(v interface{}) string {
	y := Float64(v)
	r := math.Sqrt(y)
	return VALUE(r)
}

func Sqrt(v interface{}) string {
	y := VALUE(v)
	z := "0"
	if CMP(y, 3) > 0 {
		z = y
		x := ADD(DIV(y, 2), 1)
		for CMP(x, z) < 0 {
			z = x
			x = DIV(ADD(DIV(y, x), x), 2)
		}
	} else if CMP(y, 0) != 0 {
		z = "1"
	}
	return z
}

/////

// RATE100 : 100분율값 (i*rate)/100
func RATE100(i, rate interface{}) string {
	return DIV(MUL(i, rate), 100)
}

// RATESUB100 : i - 100분율값(i*rate/100)
func RATESUB100(i, rate interface{}) string {
	per := DIV(MUL(i, rate), 100)
	return SUB(i, per)
}

// RATEADD100 : i + 100분율값(i*rate/100)
func RATEADD100(i, rate interface{}) string {
	per := DIV(MUL(i, rate), 100)
	return ADD(i, per)
}

// RATE100SUBFEE  : i - 100분율값(i*rate/100) , feevalue
func RATE100SUBFEE(i, rate interface{}) (string, string) {
	fee := DIV(MUL(i, rate), 100)
	val := SUB(i, fee)
	return val, fee
}

// PER100 : val/total*100
func PER100(val, total interface{}) string {
	return MUL(DIV(val, total), 100)
}

func IsUnderZero(v any) bool {
	return CMP(v, 0) <= 0
}
