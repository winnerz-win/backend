package mongo

import (
	"fmt"
	"reflect"
	"strings"
	"txscheduler/brix/tools/database/mongo/tools/jmath"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Decimal128Zero() primitive.Decimal128 {
	zero, _ := primitive.ParseDecimal128("0")
	return zero
}

func Decimal128Max() primitive.Decimal128 {
	max, _ := primitive.ParseDecimal128("9999999999999999999999999999999999")
	return max
}

func Decimal128Min() primitive.Decimal128 {
	min, _ := primitive.ParseDecimal128("-9999999999999999999999999999999999")
	return min
}

func Decimal128(val any) primitive.Decimal128 {
	value := jmath.VALUE(val)
	dec, err := primitive.ParseDecimal128(
		value,
	)
	if err != nil {
		fmt.Println("mongo.Decimal128(", val, ") :", err)
		max := Decimal128Max()
		if jmath.CMP(value, max) > 0 {
			return max
		}
		min := Decimal128Min()
		if jmath.CMP(value, min) < 0 {
			return min
		}
		return Decimal128Zero()
	}

	return dec
}

/////////////////////////////////////////////////////////////////////////////////////////////////////

func is_string_decimal128_format(text string) (primitive.Decimal128, bool) {
	text = strings.TrimSpace(text)
	if !strings.HasPrefix(text, "Decimal128(") {
		return Decimal128Zero(), false
	}
	if !strings.HasSuffix(text, ")") {
		return Decimal128Zero(), false
	}
	text = strings.ReplaceAll(text, "Decimal128(", "")
	text = text[:len(text)-1]

	text = strings.TrimSpace(text)

	return Decimal128(text), true
}

func decimal128_to_string_format(dec primitive.Decimal128) string {
	return _Cat("Decimal128(", jmath.VALUE(dec), ")")
}

func _string_to_decimal128(item MAP) {
	for key, val := range item {
		if val == nil {
			continue
		}
		_, _ = key, val
		switch target := val.(type) {
		case string:
			if dt_format, ok := is_string_decimal128_format(target); ok {
				item[key] = dt_format
			}
		case []string:
			sl := []interface{}{}
			for _, v := range target {
				if dt_format, ok := is_string_decimal128_format(v); ok {
					sl = append(sl, dt_format)
				} else {
					sl = append(sl, v)
				}
			}
			item[key] = sl

		case []interface{}:
			sl := []interface{}{}
			for _, v := range target {
				if text, do := v.(string); do {
					if dt_format, ok := is_string_decimal128_format(text); ok {
						sl = append(sl, dt_format)
					} else {
						sl = append(sl, v)
					}
				} else {
					sl = append(sl, v)
				}
			}
			item[key] = sl

		default:
			switch reflect.TypeOf(val).Kind() {
			case reflect.Map:
				_string_to_decimal128(
					MakeMap(val),
				)
			}
			//cc.Cyan("[", key, "] ", reflect.TypeOf(val).Kind())
		}
	} //for
}

func _decimal128_to_string(item MAP) {
	for key, val := range item {
		if val == nil {
			continue
		}
		switch target := val.(type) {
		case primitive.Decimal128:
			item[key] = decimal128_to_string_format(target)

		case primitive.A:
			sl := []interface{}{}
			for _, v := range target {
				if datetime, do := v.(primitive.Decimal128); do {
					sl = append(sl, decimal128_to_string_format(datetime))
				} else {
					sl = append(sl, v)
				}
			}
			item[key] = sl

		default:
			switch reflect.TypeOf(val).Kind() {
			case reflect.Map:
				_decimal128_to_string(
					MakeMap(val),
				)
			case reflect.Slice:
				val_list := val.([]interface{})

				sl := []interface{}{}
				for _, v := range val_list {
					if datetime, do := v.(primitive.Decimal128); do {
						sl = append(sl, decimal128_to_string_format(datetime))
					} else {
						sl = append(sl, v)
					}
				}
				item[key] = sl

			} //switch
		} //switch
	} //for
}
