package jmath

import (
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"
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

//Int64MaxString :
func Int64MaxString() string {
	return fmt.Sprintf("%v", Int64Max)
}

//Int32MaxString :
func Int32MaxString() string {
	return fmt.Sprintf("%v", Int32Max)
}

//Uint64MaxString :
func Uint64MaxString() string {
	return fmt.Sprintf("%v", Uint64Max)
}

//Uint32MaxString :
func Uint32MaxString() string {
	return fmt.Sprintf("%v", Uint32Max)
}

//VALUE :
func VALUE(i interface{}) string {
	return NewBigDecimal(i).ToString()
}

func IsFloat(i interface{}) bool {
	return strings.Contains(VALUE(i), ".")
}

//BYTES :
func BYTES(i interface{}) []byte {
	return NewBigDecimal(i).ToBigInteger().Bytes()
}

//HEX :
func HEX(i interface{}, without0x ...bool) string {
	b := BYTES(i)
	v := hex.EncodeToString(b)
	// if v == "" {
	// 	v = "00"
	// }
	if len(without0x) > 0 && without0x[0] {
		return v
	}
	return "0x" + v
}

//CMP : CompareTo
func CMP(src, tar interface{}) int {
	srcVal := NewBigDecimal(src)
	tarVal := NewBigDecimal(tar)
	return srcVal.CompareTo(tarVal)
}

//CMPSUB : ComapreTo , Val = bigger - smaller
func CMPSUB(src, tar interface{}) (int, string) {
	srcVal := NewBigDecimal(src)
	tarVal := NewBigDecimal(tar)
	cmp := srcVal.CompareTo(tarVal)
	subVal := "0"
	if cmp < 0 {
		subVal = SUB(tar, src)
	} else if cmp > 0 {
		subVal = SUB(src, tar)
	}
	return cmp, subVal
}

//RATE100 : 100분율값 (i*rate)/100
func RATE100(i, rate interface{}) string {
	return DIV(MUL(i, rate), 100)
}

//RATESUB100 : i - 100분율값(i*rate/100)
func RATESUB100(i, rate interface{}) string {
	per := DIV(MUL(i, rate), 100)
	return SUB(i, per)
}

//RATEADD100 : i + 100분율값(i*rate/100)
func RATEADD100(i, rate interface{}) string {
	per := DIV(MUL(i, rate), 100)
	return ADD(i, per)
}

//RATE100SUBFEE  : i - 100분율값(i*rate/100) , feevalue
func RATE100SUBFEE(i, rate interface{}) (string, string) {
	fee := DIV(MUL(i, rate), 100)
	val := SUB(i, fee)
	return val, fee
}

//PER100 : val/total*100
func PER100(val, total interface{}) string {
	return MUL(DIV(val, total), 100)
}

//ADD :
func ADD(i, j interface{}, etc ...interface{}) string {
	a := NewBigDecimal(i)
	b := NewBigDecimal(j)

	r := a.Add(b)
	for _, v := range etc {
		r = r.Add(NewBigDecimal(v))
	}
	return r.ToString()
}

//SUM :
func SUM(vals ...interface{}) string {
	sum := "0"
	for _, val := range vals {
		sum = ADD(sum, val)
	} //for
	return sum
}

//SUB :
func SUB(i, j interface{}, etc ...interface{}) string {
	a := NewBigDecimal(i)
	b := NewBigDecimal(j)
	r := a.Subtract(b)
	for _, v := range etc {
		r = r.Subtract(NewBigDecimal(v))
	}
	return r.ToString()
}

//MUL :
func MUL(i, j interface{}, etc ...interface{}) string {
	a := NewBigDecimal(i)
	b := NewBigDecimal(j)
	r := a.Multiply(b)
	for _, v := range etc {
		r = r.Multiply(NewBigDecimal(v))
	}
	return r.ToString()
}

//DIV :
func DIV(i, j interface{}, etc ...interface{}) string {
	a := NewBigDecimal(i)
	b := NewBigDecimal(j)
	r := a.Divide(b)
	for _, v := range etc {
		r = r.Divide(NewBigDecimal(v))
	}
	return r.ToString()
}

//DIVZERO : 4 / 3 = 1  cut>.33333333
func DIVZERO(i, j interface{}) string {
	return DOTCUT(DIV(i, j), 0)
}

//ROUND :
func ROUND(i, place interface{}) string {
	a := NewBigDecimal(i)
	p := int32(NewBigDecimal(place).ToBigInteger().Int64())
	return a.RoundC(p).ToString()
}

//UP :
func UP(i, place interface{}) string {
	a := NewBigDecimal(i)
	p := int32(NewBigDecimal(place).ToBigInteger().Int64())
	return a.UpC(p).ToString()
}

//DOTCUT :
func DOTCUT(i, place interface{}) string {
	p := int(NewBigDecimal(place).ToBigInteger().Int64())
	val := NewBigDecimal(i)
	if p <= 0 {
		return val.FloorC().ToString()
	}

	a := val.ToString()

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

func DotCount(v interface{}) int {
	v2 := VALUE(v)
	ss := strings.Split(v2, ".")
	if len(ss) <= 1 {
		return 0
	}
	return len(ss[1])
}

//Int :
func Int(v interface{}) int {
	return int(Int64(v))
}

//Int64 :
func Int64(v interface{}) int64 {
	val := NewBigDecimal(v)
	return val.ToBigInteger().Int64()
}

//Uint64 :
func Uint64(v interface{}) uint64 {
	return New(v).Uint64()
}

//Float64 :
func Float64(v interface{}) float64 {
	a := VALUE(v)
	f, _ := strconv.ParseFloat(a, 64)
	return f
}

func Pow(x, n interface{}) string {
	x1 := Float64(x)
	n1 := Float64(n)
	r := math.Pow(x1, n1)
	return VALUE(r)
}

func POW(x, n interface{}) string {
	if CMP(n, 0) == 0 {
		return "1"
	}
	//  else if CMP(n, 1) == 0 {
	// 	return VALUE(x)
	// }

	r := VALUE(x)
	p := r
	loop := Int(n)
	for loop > 1 {
		r = MUL(r, p)
		loop--
	} //for
	return r
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

func MIN(x, y interface{}) string {
	x1 := VALUE(x)
	y1 := VALUE(y)
	if CMP(x1, y1) < 0 {
		return x1
	}
	return y1
}
func MAX(x, y interface{}) string {
	x1 := VALUE(x)
	y1 := VALUE(y)
	if CMP(x1, y1) < 0 {
		return y1
	}
	return x1
}
func ABS(v interface{}) string {
	v1 := VALUE(v)
	if CMP(v1, 0) < 0 {
		return MUL(v1, -1)
	}
	return v1
}

//DivDot : 정수부 , 소수부 (0.xxxx)
func DivDot(v interface{}) (float64, float64) {
	val := VALUE(v)
	ss := strings.Split(val, ".")
	if len(ss) == 0 {
		return 0, 0
	} else if len(ss) == 1 {
		return Float64(ss[0]), 0
	}

	return Float64(ss[0]), Float64("0." + ss[1])
}

//MOD : 123 % 2 = 1
func MOD(a, b interface{}) string {
	div := DIV(a, b)
	c := DOTCUT(div, 0)
	r := SUB(a, MUL(c, b))
	return r
}
