package jcrypt

import "crypto/rand"

var (
	mcA0 = []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J",
		"K", "L", "M", "N", "O", "P", "Q", "R", "S", "T",
		"U", "V", "W", "X", "Y", "Z",
	}

	mcA1 = []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J",
		"K", "L", "M", "N", "O", "P", "Q", "R", "S", "T",
		"U", "V", "W", "X", "Y", "Z", "a", "b", "c", "d",
		"e", "f", "g", "h", "i", "j", "k", "m", "n", "p",
		"q", "r", "x", "t", "u", "v", "w", "x", "y", "z",
		//"l", "o",
	}
)

func _make_pin(max int, cs []string) string {
	size := max * 2
	rb := make([]byte, size)
	_, err := rand.Read(rb)
	if err != nil {
		panic(err)
	}
	mod := byte(36)

	code := ""
	for i := 0; i < size; i += 2 {
		code += cs[int((rb[i]+rb[i+1])%mod)]
	}
	return code
}

//MAKEPINCODE : [0~9,A~Z] default cnt : 6
func MAKEPINCODE(size ...int) string {
	max := 6
	if len(size) > 0 {
		max = size[0]
	}
	return _make_pin(max, mcA0)
}

//MakePinCode : [0~9,A~Z] X(l,o) default cnt : 6
func MakePinCode(size ...int) string {
	max := 6
	if len(size) > 0 {
		max = size[0]
	}
	return _make_pin(max, mcA1)
}
