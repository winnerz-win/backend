package model

import (
	"jtools/jmath"
	"strings"
	"txscheduler/brix/tools/dbg"
)

const (
	ETH = "ETH"
)

type CoinData map[string]string

func (my CoinData) String() string { return dbg.ToJSONString(my) }

// NewCoinData :
func NewCoinData() CoinData {
	return CoinData{
		ETH: ZERO,
	}
}

// NewCoinDataSymbol :
func NewCoinDataSymbol(symbols ...string) CoinData {
	my := CoinData{}
	for _, symbol := range symbols {
		symbol = strings.TrimSpace(symbol)
		if symbol == "" {
			continue
		}
		my[symbol] = ZERO
	}
	return my
}

// AddCoinData :
func (my CoinData) AddCoinData(coin CoinData) {
	for k, v := range coin {
		my.ADD(k, v)
	}
}

// Price :
func (my CoinData) Price(symbol string) string {
	if price, do := my[symbol]; do {
		return price
	}
	return ZERO
}

// SET :
func (my CoinData) SET(symbol, price string) string {
	my[symbol] = price
	return my[symbol]
}

// ZERO :
func (my CoinData) ZERO(symbol string) {
	my[symbol] = ZERO
}

// ADD :
func (my CoinData) ADD(symbol, price string) string {
	if v, do := my[symbol]; do {
		my[symbol] = jmath.ADD(v, price)
	} else {
		my[symbol] = price
	}
	return my[symbol]
}

// SUB :
func (my CoinData) SUB(symbol, price string) string {
	if v, do := my[symbol]; do {
		r := jmath.SUB(v, price)
		if jmath.CMP(r, 0) < 0 {
			r = ZERO
		}
		my[symbol] = r
	} else {
		if jmath.CMP(price, 0) < 0 {
			price = ZERO
		}
		my[symbol] = price
	}
	return my[symbol]
}

// Clone :
func (my CoinData) Clone() CoinData {
	clone := CoinData{}
	for k, v := range my {
		clone[k] = v
	}
	return clone
}

// Do :
func (my CoinData) Do() bool {
	for _, v := range my {
		if jmath.CMP(v, ZERO) > 0 {
			return true
		}
	}
	return false
}
