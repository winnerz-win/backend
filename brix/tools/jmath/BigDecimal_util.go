package jmath

///////////////////////////////////////////////////////////////////////////////////////

func BigDecimal_ZERO() *BigDecimal {
	return NewBigDecimal(0)
}

func BigDecimal_ONE() *BigDecimal {
	return NewBigDecimal(1)
}

func BigDecimal_TEN() *BigDecimal {
	return NewBigDecimal(10)
}

func BigDecimal_100() *BigDecimal {
	return NewBigDecimal(100)
}

//CompareTo :
func CompareTo(src, tar interface{}) int {
	srcVal := NewBigDecimal(src)
	tarVal := NewBigDecimal(tar)
	return srcVal.CompareTo(tarVal)
}

//Add :
func Add(a, b interface{}) *BigDecimal {
	a1 := NewBigDecimal(a)
	b1 := NewBigDecimal(b)
	return a1.Add(b1)
}

//Subtract :
func Subtract(a, b interface{}) *BigDecimal {
	a1 := NewBigDecimal(a)
	b1 := NewBigDecimal(b)
	return a1.Subtract(b1)
}

//Multiply :
func Multiply(a, b interface{}) *BigDecimal {
	a1 := NewBigDecimal(a)
	b1 := NewBigDecimal(b)
	return a1.Multiply(b1)
}

//Divide :
func Divide(a, b interface{}) *BigDecimal {
	a1 := NewBigDecimal(a)
	b1 := NewBigDecimal(b)
	return a1.Divide(b1)
}
