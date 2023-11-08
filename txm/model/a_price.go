package model

import (
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
)

//PRICE :
type PRICE string

func (my PRICE) val() string { return string(my) }

//Float64 :
func (my PRICE) Float64() float64 { return jmath.Float64(my) }

const (
	priceNONE = "none"
	//PRICEZERO :
	PRICEZERO = PRICE("0")
)

func parsePRICE(v interface{}) (string, bool) {
	switch v.(type) {
	case PRICE:
		return string(v.(PRICE)), true
	case *PRICE:
		return string(*(v.(*PRICE))), true
	default:
		if jmath.IsNum(v) {
			return jmath.VALUE(v), true
		} else {
			dbg.RedItalic(v, "is not PRICE format.")
		}
	}
	return priceNONE, false
}

//SET :
func (my *PRICE) SET(v interface{}) {
	if val, do := parsePRICE(v); do {
		*my = PRICE(val)
	}
}

//ADD :
func (my *PRICE) ADD(v interface{}) {
	if val, do := parsePRICE(v); do {
		*my = PRICE(jmath.ADD(my.val(), val))
	}
}

//SUB :
func (my *PRICE) SUB(v interface{}) {
	if val, do := parsePRICE(v); do {
		*my = PRICE(jmath.SUB(my.val(), val))
	}
}

//DIV :
func (my PRICE) DIV(a interface{}) PRICE {
	if val, do := parsePRICE(a); do {
		return PRICE(jmath.DIV(my.val(), val))
	}
	return PRICEZERO
}

//MUL :
func (my PRICE) MUL(a interface{}) PRICE {
	if val, do := parsePRICE(a); do {
		return PRICE(jmath.MUL(my.val(), val))
	}
	return PRICEZERO
}

//RATE100 :
func (my PRICE) RATE100(a interface{}) PRICE {
	if val, do := parsePRICE(a); do {
		return PRICE(jmath.RATE100(my.val(), val))
	}
	return PRICEZERO
}

//RATESUB100 :
func (my PRICE) RATESUB100(a interface{}) PRICE {
	if val, do := parsePRICE(a); do {
		return PRICE(jmath.RATESUB100(my.val(), val))
	}
	return PRICEZERO
}
