package jmath

import (
	"fmt"
	"testing"
)

func TestJmath(t *testing.T) {

	fmt.Println(
		VALUE(33),
		CMP(22, 33),
		ABS(-33),
		MUL(33, 2),
		DOTCUT(3.123456789, 3),
		MIN(39, 55),
		MAX(39, 55),
		MOD(2, 3),
	)

	fmt.Println(
		ADD(1, 2, 3, 4, 5, 6, 7, 8, 9, 10),
		1+2+3+4+5+6+7+8+9+10,
		ADD(1, 2),
		SUB(22, 1),
		DIV(33, 2),
	)

	fmt.Println(1, 2, 3)
	fmt.Printf("%v", 10)
	fmt.Printf("%v", 11)
	fmt.Println()
}
