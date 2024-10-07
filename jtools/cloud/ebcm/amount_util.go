package ebcm

import (
	"jtools/jmath"
	"strconv"
	"strings"
)

const (
	UINT256MAX = "115792089237316195423570985008687907853269984665640564039457584007913129639935"
)

// ETHToWei :
func ETHToWei(eth string) string {
	return etherToWeiString(eth, "18")
}

// WeiToETH :
func WeiToETH(wei string) string {
	return weiToEtherString(wei, "18")
}

// TokenToWei :
func TokenToWei(tval string, decimal interface{}) string {
	dec := ""
	v, do := decimal.(string)
	if do {
		dec = v
	} else {
		dec = jmath.VALUE(decimal)
	}
	return etherToWeiString(tval, dec)
}

// WeiToToken :
func WeiToToken(wei string, decimal interface{}) string {
	dec := ""
	v, do := decimal.(string)
	if do {
		dec = v
	} else {
		dec = jmath.VALUE(decimal)
	}
	return weiToEtherString(wei, dec)
}

func weiToEtherString(wei string, dec string) string {
	v, _ := strconv.ParseInt(dec, 10, 64)
	ds := "0."
	for i := 0; i < int(v)-1; i++ {
		ds += "0"
	}
	ds += "1"
	return jmath.MUL(ds, wei)
}

func etherToWeiString(eth string, dec string) string {
	v, _ := strconv.ParseInt(dec, 10, 64)
	ds := "1"
	for i := 0; i < int(v); i++ {
		ds += "0"
	}
	wei := jmath.MUL(ds, eth)
	if strings.Contains(wei, ".") {
		return "0"
	}
	return wei
}

//////////////////////////////////////////////////////////////////////////////

// GWEI :
type GWEI string

// String :
func (my GWEI) String() string { return string(my) }

// Desc :
func (my GWEI) Desc() string { return string(my) + " gwei" }

// ToETH :
func (my GWEI) ToETH() ETH { return GWEIToETH(my) }

// ToWEI :
func (my GWEI) ToWEI() WEI { return GWEIToWEI(my) }

//////////////////////////////////////////////////////////////////////////////

// ETH :
type ETH string

// String :
func (my ETH) String() string { return string(my) }

// Desc :
func (my ETH) Desc() string { return string(my) + " eth" }

// ToGWEI :
func (my ETH) ToGWEI() GWEI { return ETHToGWEI(my) }

// ToWEI :
func (my ETH) ToWEI() WEI { return ETHToWEI(my) }

//////////////////////////////////////////////////////////////////////////////

// WEI :
type WEI string

// String :
func (my WEI) String() string { return string(my) }

// Desc :
func (my WEI) Desc() string { return string(my) + " wei" }

// ToGWEI :
func (my WEI) ToGWEI() GWEI { return WEIToGWEI(my) }

// ToETH :
func (my WEI) ToETH() ETH { return WEIToETH(my) }

// UInt64 :
func (my WEI) UInt64() uint64 { return uint64(jmath.Int64(my)) }

//////////////////////////////////////////////////////////////////////////////

// ETHToGWEI :
func ETHToGWEI(eth interface{}) GWEI {
	return GWEI(etherToWeiString(jmath.VALUE(eth), "9"))
}

// ETHToWEI :
func ETHToWEI(eth interface{}) WEI {
	return WEI(etherToWeiString(jmath.VALUE(eth), "18"))
}

// WEIToETH :
func WEIToETH(wei interface{}) ETH {
	return ETH(weiToEtherString(jmath.VALUE(wei), "18"))
}

// WEIToGWEI :
func WEIToGWEI(wei interface{}) GWEI {
	return GWEI(weiToEtherString(jmath.VALUE(wei), "9"))
}

// GWEIToETH :
func GWEIToETH(gwei interface{}) ETH {
	return ETH(weiToEtherString(jmath.VALUE(gwei), "9"))
}

// GWEIToWEI :
func GWEIToWEI(gwei interface{}) WEI {
	return WEI(etherToWeiString(jmath.VALUE(gwei), "9"))
}
