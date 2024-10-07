package jmath

//GetN :
func GetN(n interface{}) int {
	var v int = 0
	switch n.(type) {
	case int:
		v = n.(int)
	case byte:
		v = int(n.(byte))
	case uint:
		v = int(n.(uint))
	case int8:
		v = int(n.(int8))
	case int16:
		v = int(n.(int16))
	case uint16:
		v = int(n.(uint16))
	case int32:
		v = int(n.(int32))
	case uint32:
		v = int(n.(uint32))
	case int64:
		v = int(n.(int64))
	case uint64:
		v = int(n.(uint64))
	default:
		panic("zava.GetN is ERROR VALUE")
	}

	return v
}
