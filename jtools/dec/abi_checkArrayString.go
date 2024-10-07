package dec

/*
	default:
		if array, ok := dec.CheckArrayString(void); ok {
			return array, ok
		}
*/

//CheckArrayString : ebcm/abi/abi_types.go [	func checkArrayString(data ...interface{}) ([]string, bool) {	]
func CheckArrayString(v any) ([]string, bool) {
	switch void := v.(type) {
	case []Int64:
		array := []string{}
		for _, v := range void {
			array = append(array, v.Value())
		}
		return array, true

	case Int64List:
		array := []string{}
		for _, v := range void {
			array = append(array, v.Value())
		}
		return array, true

	case []Uint256:
		array := []string{}
		for _, v := range void {
			array = append(array, v.Int64().Value())
		}
		return array, true

	case Uint256List:
		array := []string{}
		for _, v := range void {
			array = append(array, v.Int64().Value())
		}
		return array, true

	} //switch

	return nil, false
}
