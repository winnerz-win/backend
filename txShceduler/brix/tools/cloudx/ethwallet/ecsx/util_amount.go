package ecsx

import (
	"strconv"
	"txscheduler/brix/tools/cloudx/ethwallet/EtherClient"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
	"txscheduler/brix/tools/jmath/decimal"
)

func init() {
	gwei := GWEI("1000")
	dbg.PrintInit("ecsx.util_amount : ", gwei.Desc())
	EtherClient.SetOverLimitWEI(gwei.ToWEI().String())
}

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
		dec = jmath.NewBigDecimal(decimal).ToString()
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
		dec = jmath.NewBigDecimal(decimal).ToString()
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
	zari := decimal.NewDecimal(ds)
	val := decimal.NewDecimal(wei)
	r := val.Mul(zari)
	return r.String()
}

func etherToWeiString(eth string, dec string) string {
	v, _ := strconv.ParseInt(dec, 10, 64)
	ds := "1"
	for i := 0; i < int(v); i++ {
		ds += "0"
	}
	zari := decimal.NewDecimal(ds)
	val := decimal.NewDecimal(eth)
	r := val.Mul(zari)
	return r.String()
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
func ETHToGWEI(eth ETH) GWEI {
	return GWEI(etherToWeiString(eth.String(), "9"))
}

// ETHToWEI :
func ETHToWEI(eth ETH) WEI {
	return WEI(etherToWeiString(eth.String(), "18"))
}

// WEIToETH :
func WEIToETH(wei WEI) ETH {
	return ETH(weiToEtherString(wei.String(), "18"))
}

// WEIToGWEI :
func WEIToGWEI(wei WEI) GWEI {
	return GWEI(weiToEtherString(wei.String(), "9"))
}

// GWEIToETH :
func GWEIToETH(gwei GWEI) ETH {
	return ETH(weiToEtherString(gwei.String(), "9"))
}

// GWEIToWEI :
func GWEIToWEI(gwei GWEI) WEI {
	return WEI(etherToWeiString(gwei.String(), "9"))
}
