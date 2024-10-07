package dec

import (
	"jtools/jmath"
	"strings"
)

type Type string

const (
	TYPE_NONE     = Type("none")
	TYPE_INT64    = Type("int64")
	TYPE_UINT256  = Type("uint256") //78
	TYPE_MINUS    = Type("minus")
	TYPE_OVERFLOW = Type("overflow")
)

///////////////////////////////////////////

type Number string

func (my Number) _size() int { return len(my) }

func (my Number) _string() string { return string(my) }
func (my Number) _v() string      { return jmath.VALUE(my) }

func (my Number) Type() Type {
	size := my._size()
	switch {
	case size > 78:
		return TYPE_OVERFLOW

	case size == 78:
		if jmath.CMP(my, UINT256MAX) > 0 {
			return TYPE_OVERFLOW
		}
		return TYPE_UINT256

	case size >= 1 && size < 78:
		if strings.HasPrefix(my._string(), "-") {
			return TYPE_MINUS
		}
		if jmath.CMP(my, UINT256MAX) > 0 {
			return TYPE_OVERFLOW
		}
		return TYPE_INT64
	default:
		return TYPE_NONE
	} //switch
}

// Int64Value : Int64string [0000...0009 -> string(9)]
func (my Number) Int64Value() string {
	switch my.Type() {
	case TYPE_INT64:
		return my._string()
	case TYPE_NONE:
		return my._string()
	} //switch

	return my.Int64()._string()
}

// Uint256Value : Uint256string [ string(9) -> string(000...0009) ]
func (my Number) Uint256Value() string {
	switch my.Type() {
	case TYPE_UINT256:
		return my._string()
	case TYPE_NONE:
		return my._string()
	} //switch

	return my.Uint256()._string()
}

func (my Number) Int64() Int64 {
	switch my.Type() {
	case TYPE_UINT256:
		return INT64(my)
	} //switch

	return Int64(my)
}

func (my Number) Uint256() Uint256 {
	switch my.Type() {
	case TYPE_INT64:
		return _UINT256(my)
	} //switch

	return Uint256(my)
}

func (my *Number) SetInt64() Int64 {
	if my.Type() != TYPE_INT64 {
		*my = INT64(*my).Number()
	}
	return Int64(*my)
}

func (my *Number) SetUint256() Uint256 {
	if my.Type() != TYPE_UINT256 {
		*my = _UINT256(*my).Number()
	}
	return Uint256(*my)
}

////////////////////////////////////////

func (my Number) Cmp(v any) int {
	return jmath.CMP(my.Int64()._string(), v)
}
