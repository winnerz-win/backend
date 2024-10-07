package mongo

import (
	"fmt"
	"os"
	"reflect"
)

func StructNameToLower(stp interface{}) string {
	rt := reflect.TypeOf(stp)
	_, text := _struct_reflect_name(rt)
	return _struct_name_tolowter(text)
}

//////////////////////////////////////////////////////////////////////////////////

func _struct_reflect_name(rt reflect.Type) (reflect.Type, string) {
	var dt reflect.Type
	text := ""
	switch rt.Kind() {
	case reflect.Slice:
		dt, text = _struct_reflect_name(rt.Elem())

	case reflect.Struct:
		dt = rt
		text = rt.Name()

	case reflect.Pointer:
		dt, text = _struct_reflect_name(rt.Elem())

	default:
		fmt.Println("mongo.StructNameToLower[", rt, "] ERROR")
		fmt.Println(_Stack(9))
		os.Exit(400004)
	} //switch

	if text == "" { //struct{}{}
		fmt.Println("mongo.StructNameToLower[", rt, "] ERROR")
		fmt.Println(_Stack(9))
		os.Exit(400005)
	}

	return dt, text
}

func _struct_name_tolowter(text string) string {
	_gHL := func(v rune) int {
		if v == 95 { // _
			return 95
		}
		if v >= 97 && v <= 122 { //a~z
			return 0
		} else if v >= 65 && v <= 90 { //A~Z
			return 1
		}
		return 2 //0~9
	}
	_toLower := func(v rune) rune {
		if v >= 65 && v < 90 {
			return v + 32
		}
		return v
	}

	array := []rune(text)

	prv := _gHL(array[0]) //0:L , 1: H ,  2: number , 95:_
	re := []rune{_toLower(array[0])}

	for _, v := range array[1:] {
		cur := _gHL(v)
		switch cur {
		case 95: //_
			prv = cur
			re = append(re, _toLower(v))

		case 0: //a~z
			prv = cur
			re = append(re, _toLower(v))

		case 1: //A~Z
			if prv == 0 || prv == 2 {
				prv = cur
				re = append(re, 95, _toLower(v))
			} else {
				re = append(re, _toLower(v))
			}
		case 2: //0~9
			prv = cur
			re = append(re, _toLower(v))
		}
	} //for
	return string(re)
}
