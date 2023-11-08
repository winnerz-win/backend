package decimal

import (
	"fmt"
	"math/big"
	"strings"

	"txscheduler/brix/tools/jmath/decimal/idecimal"
)

//NewDecimal :
func NewDecimal(v interface{}, isErr ...*error) Decimal {
	value := "0"
	switch v.(type) {

	case idecimal.IDecimal:
		value = v.(idecimal.IDecimal).ToIDecimal()

	case []byte:
		vv := big.NewInt(0)
		vv = vv.SetBytes(v.([]byte))
		value = vv.String()

	case Decimal:
		value = v.(Decimal).String()
		value = strings.TrimSpace(value)

	default:
		value = fmt.Sprintf("%v", v)
	}
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "0x") {
		bi := big.NewInt(0)
		_, do := bi.SetString(value[2:], 16)
		if do {
			value = bi.String()
		}
	}
	dc, err := NewFromString(value)
	if len(isErr) > 0 {
		*isErr[0] = err
	}
	return dc
}
