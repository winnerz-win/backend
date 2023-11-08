package jmath

//String : value-type to string
func String(v interface{}) string {
	return NewBigDecimal(v).ToString()
}
