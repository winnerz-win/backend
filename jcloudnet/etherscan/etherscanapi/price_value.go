package etherscanapi

import (
	"jtools/jmath"
	"strconv"
)

// WeiToEtherString :
func WeiToEtherString(wei string, dicimal int) string {
	ds := "0."
	for i := 0; i < dicimal-1; i++ {
		ds += "0"
	}
	ds += "1"
	zari := jmath.VALUE(ds)
	val := jmath.VALUE(wei)
	return jmath.MUL(val, zari)
}

// WeiToEtherString2 : EtherScanAPI.ChValueOutOfDecimals is obsoluted
func WeiToEtherString2(wei, decimal string) string {
	v, _ := strconv.ParseInt(decimal, 10, 64)
	return WeiToEtherString(wei, int(v))
}

// EtherToWeiString :
func EtherToWeiString(eth string, dicimal int) string {
	ds := "1"
	for i := 0; i < dicimal; i++ {
		ds += "0"
	}
	zari := jmath.VALUE(ds)
	val := jmath.VALUE(eth)
	return jmath.MUL(val, zari)
}

// EtherToWeiString2 :
func EtherToWeiString2(eth, decimal string) string {
	v, _ := strconv.ParseInt(decimal, 10, 64)
	return EtherToWeiString(eth, int(v))
}
