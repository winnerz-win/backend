package dec

import "jtools/jmath"

///////////////////////////////////////////

type dec_number interface {
	_v() string
	Int64() Int64
	Uint256() Uint256
}

func CMP(a, b dec_number) int {
	return jmath.CMP(a._v(), b._v())
}

func ADD(a, b dec_number) dec_number {
	return Number(jmath.ADD(a._v(), b._v()))
}

func SUB(a, b dec_number) dec_number {
	return Number(jmath.SUB(a._v(), b._v()))
}

func MUL(a, b dec_number) dec_number {
	return Number(jmath.MUL(a._v(), b._v()))
}

func DIVInt64(a, b dec_number) (Int64, Int64) {
	v := jmath.DIV(a._v(), b._v())
	c := jmath.DOTCUT(v, 0)
	r := jmath.SUB(a._v(), jmath.MUL(c, b._v()))
	return INT64(c), INT64(r)
}

func DIVUint256(a, b dec_number) (Uint256, Uint256) {
	v := jmath.DIV(a._v(), b._v())
	c := jmath.DOTCUT(v, 0)
	r := jmath.SUB(a._v(), jmath.MUL(c, b._v()))
	return _UINT256(c), _UINT256(r)
}
