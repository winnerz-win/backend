package jmath

func Int(v interface{}) int {
	return NEW(v).Int()
}

func Int64(v interface{}) int64 {
	return NEW(v).Int64()
}
func Uint64(v interface{}) uint64 {
	return NEW(v).Uint64()
}

func Float64(v interface{}) float64 {
	return NEW(v).Float64()
}

//IsNum :
func IsNum(v interface{}) bool {
	_, err := newBigNumber2(v)
	return err == nil
}
