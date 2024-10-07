package mongo

import (
	"jtools/jmath"
	"jtools/unix"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func is_string_datetime_format(text string) (primitive.DateTime, bool) {
	text = strings.TrimSpace(text)
	if !strings.HasPrefix(text, "DateTime(") {
		return 0, false
	}
	if !strings.HasSuffix(text, ")") {
		return 0, false
	}
	text = strings.ReplaceAll(text, "DateTime(", "")
	text = text[:len(text)-1]

	text = strings.TrimSpace(text)

	//cc.Red(text)
	if jmath.IsNum(text) { //unix.Time
		//cc.Red(unix.Time(jmath.Int64(text)).DateTimeString())
		return unix.Time(jmath.Int64(text)).DateTime(), true
	}
	//cc.Red(unix.FromRFC3339(text).DateTimeString())
	return unix.FromRFC3339(text).DateTime(), true
}
func datetime_to_string_format(datetime primitive.DateTime) string {
	return _Cat("DateTime(", unix.FromTime(datetime.Time()).DateTimeString(), ")")
}

func _string_to_datetime(item MAP) {
	for key, val := range item {
		if val == nil {
			continue
		}
		_, _ = key, val
		switch target := val.(type) {
		case string:
			if dt_format, ok := is_string_datetime_format(target); ok {
				item[key] = dt_format
			}
		case []string:
			sl := []interface{}{}
			for _, v := range target {
				if dt_format, ok := is_string_datetime_format(v); ok {
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
					if dt_format, ok := is_string_datetime_format(text); ok {
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
				_string_to_datetime(
					MakeMap(val),
				)
			}
			//cc.Cyan("[", key, "] ", reflect.TypeOf(val).Kind())
		}
	} //for
}

func _datetime_to_string(item MAP) {
	for key, val := range item {
		if val == nil {
			continue
		}
		switch target := val.(type) {
		case primitive.DateTime:
			item[key] = datetime_to_string_format(target)

		case primitive.A:
			sl := []interface{}{}
			for _, v := range target {
				if datetime, do := v.(primitive.DateTime); do {
					sl = append(sl, datetime_to_string_format(datetime))
				} else {
					sl = append(sl, v)
				}
			}
			item[key] = sl

		default:
			switch reflect.TypeOf(val).Kind() {
			case reflect.Map:
				_datetime_to_string(
					MakeMap(val),
				)
			case reflect.Slice:
				val_list := val.([]interface{})

				sl := []interface{}{}
				for _, v := range val_list {
					if datetime, do := v.(primitive.DateTime); do {
						sl = append(sl, datetime_to_string_format(datetime))
					} else {
						sl = append(sl, v)
					}
				}
				item[key] = sl

			} //switch
		} //switch
	} //for
}

////////////////////////////////////////////////////////////////////////

/*
DAteTimeFormat :

	"DateTime(1627601407)"
	"DateTime(2023-05-16T22:00:00+09:00)"
	"DateTime(2023-05-16T22:00:00+09:00)"
*/
func DateTimeFormat(text string) (primitive.DateTime, bool) {
	return is_string_datetime_format(text)
}
