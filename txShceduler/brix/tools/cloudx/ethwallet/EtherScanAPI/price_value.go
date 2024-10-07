package EtherScanAPI

import (
	"strconv"

	"txscheduler/brix/tools/jmath/decimal"
)

//WeiToEtherString :
func WeiToEtherString(wei string, dicimal int) string {
	ds := "0."
	for i := 0; i < dicimal-1; i++ {
		ds += "0"
	}
	ds += "1"
	zari := decimal.NewDecimal(ds)
	val := decimal.NewDecimal(wei)
	r := val.Mul(zari)
	return r.String()
}

//WeiToEtherString2 : EtherScanAPI.ChValueOutOfDecimals is obsoluted
func WeiToEtherString2(wei, decimal string) string {
	v, _ := strconv.ParseInt(decimal, 10, 64)
	return WeiToEtherString(wei, int(v))
}

//EtherToWeiString :
func EtherToWeiString(eth string, dicimal int) string {
	ds := "1"
	for i := 0; i < dicimal; i++ {
		ds += "0"
	}
	zari := decimal.NewDecimal(ds)
	val := decimal.NewDecimal(eth)
	r := val.Mul(zari)
	return r.String()
}

//EtherToWeiString2 :
func EtherToWeiString2(eth, decimal string) string {
	v, _ := strconv.ParseInt(decimal, 10, 64)
	return EtherToWeiString(eth, int(v))
}
