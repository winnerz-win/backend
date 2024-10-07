package dec

import (
	"fmt"
	"jtools/dbg"
	"jtools/jmath"
	"strings"
)

const (
	EmptyInt64   = Int64("")
	EmptyUint256 = Uint256("")
)

func NoneUint256(v any) Uint256 {
	return Uint256(strings.TrimSpace(dbg.Void(v)))
}
func NoneInt64(v any) Int64 {
	return Int64(strings.TrimSpace(dbg.Void(v)))
}

///////////////////////////////////////////

const (
	UINT256MAX = "115792089237316195423570985008687907853269984665640564039457584007913129639935"
	UINT128MAX = "340282366920938463463374607431768211455"
	UINT112MAX = "5192296858534827628530496329220095"
	UINT64MAX  = "18446744073709551615"
	UINT32MAX  = "4294967295"
	UINT16MAX  = "65535"
	UINT8MAX   = "255"
)

func UINT256(v any, ignore_uint256_max ...bool) Uint256 {
	num := jmath.VALUE(v)
	if !dbg.IsTrue(ignore_uint256_max) {
		if jmath.CMP(num, UINT256MAX) > 0 {
			num = UINT256MAX
		}
	}

	num = fmt.Sprintf("%078v", num)
	return Uint256(Number(num))
}

func _UINT256(v any) Uint256 {
	num := jmath.VALUE(v)
	num = fmt.Sprintf("%078v", num)
	return Uint256(Number(num))
}

func INT64(v any) Int64 {
	num := jmath.VALUE(v)
	return Int64(Number(num))
}
func INT64Value(v any) string {
	return INT64(v).Value()
}
func UINT256Value(v any) string {
	return _UINT256(v).Value()
}

func INT64_LIST(ns ...string) Int64List {
	list := make(Int64List, len(ns))
	for i, v := range ns {
		list[i] = INT64(v)
	}
	return list
}
func INT64_LIST2(vs ...int64) Int64List {
	list := make(Int64List, len(vs))
	for i, v := range vs {
		list[i] = INT64(v)
	}
	return list
}

func UINT256_LIST(ns ...string) Uint256List {
	list := make(Uint256List, len(ns))
	for i, v := range ns {
		list[i] = _UINT256(v)
	}
	return list
}

func UINT256_LIST2(vs ...int64) Uint256List {
	list := make(Uint256List, len(vs))
	for i, v := range vs {
		list[i] = _UINT256(v)
	}
	return list
}
