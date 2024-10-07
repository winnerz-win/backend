package dbg

import (
	"encoding/json"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChangeStruct : void -> struct
func ChangeStruct(void interface{}, targetPointer interface{}) error {
	return ParseStruct(void, targetPointer)
}

func ParseStruct(src, dst interface{}) error {
	if src == nil {
		return Error("src is nil")
	}
	if dst == nil {
		return Error("dst is nil")
	}

	// rt := reflect.TypeOf(src)
	// switch rt.Kind() {
	// case reflect.Slice:
	// 	switch src.(type) {
	// 	case []byte:
	// 	default:
	// 		return copier.Copy(dst, src)
	// 	}

	// case reflect.Map:
	// 	return copier.Copy(dst, src)

	// case reflect.Struct:
	// 	return copier.Copy(dst, src)

	// case reflect.Ptr:
	// 	switch rt.Elem().Kind() {
	// 	case reflect.Slice:
	// 		switch src.(type) {
	// 		case []byte:
	// 		default:
	// 			return copier.Copy(dst, src)
	// 		}

	// 	case reflect.Map:
	// 		return copier.Copy(dst, src)

	// 	case reflect.Struct:
	// 		return copier.Copy(dst, src)
	// 	}
	// }

	_parse := func(src interface{}) error {
		if b, err := json.Marshal(src); err != nil {
			return Error("[Marshal]", err)
		} else {
			if err := json.Unmarshal(b, dst); err != nil {
				return err
			}
		}
		return nil
	}

	switch v := src.(type) {
	case primitive.D:
		m := v.Map()
		return _parse(m)

	case []byte:
		return json.Unmarshal(v, dst)
	case string:
		return json.Unmarshal([]byte(v), dst)
	}

	return _parse(src)
}

func DecodeStruct[T any](src interface{}) (T, error) {
	var r T
	if src == nil {
		return r, Error("src is nil")
	}

	_parse := func(src interface{}) error {
		if b, err := json.Marshal(src); err != nil {
			return err
		} else {
			if err := json.Unmarshal(b, &r); err != nil {
				return err
			}
		}
		return nil
	}

	// rt := reflect.TypeOf(src)
	// switch rt.Kind() {
	// case reflect.Slice:
	// 	switch src.(type) {
	// 	case []byte:
	// 	default:
	// 		err := copier.Copy(&r, src)
	// 		return r, err
	// 	}

	// case reflect.Map:
	// 	err := copier.Copy(&r, src)
	// 	return r, err

	// case reflect.Struct:
	// 	err := copier.Copy(&r, src)
	// 	return r, err

	// case reflect.Ptr:
	// 	switch rt.Elem().Kind() {
	// 	case reflect.Slice:
	// 		switch src.(type) {
	// 		case []byte:
	// 		default:
	// 			err := copier.Copy(&r, src)
	// 			return r, err
	// 		}

	// 	case reflect.Map:
	// 		err := copier.Copy(&r, src)
	// 		return r, err

	// 	case reflect.Struct:
	// 		err := copier.Copy(&r, src)
	// 		return r, err
	// 	}
	// }

	switch v := src.(type) {
	case primitive.D:
		m := v.Map()
		if err := _parse(m); err != nil {
			return r, err
		}

	case []byte:
		if err := json.Unmarshal(v, &r); err != nil {
			return r, err
		}
	case string:
		if err := json.Unmarshal([]byte(v), &r); err != nil {
			return r, err
		}
	default:
		if err := _parse(src); err != nil {
			return r, err
		}
	} //switch

	return r, nil
}

const (
	____gap___ = "    "
)

func jsonFactory(rt reflect.Type, space, tag string, write func(s string)) {
	lineGap := space + ____gap___
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		tagname, ok := f.Tag.Lookup(tag)
		if ok {
			if tagname == "-" || tagname == "" {
				continue
			}
			optionTag := ""
			if strings.Contains(tagname, ",") {
				ss := strings.Split(tagname, ",")
				tagname = strings.TrimSpace(ss[0])
				optionTag = ss[1]
			}

			if f.Type.Kind() == reflect.Struct {
				//Red(f.Name, f.Type.Name())
				if tagname != "" {
					write(lineGap + tagname + " : " + f.Type.String() + " {\n")
				} else if optionTag != "" {
					write(lineGap + f.Type.String() + " <" + optionTag + "> {\n")
				} else {
					write(lineGap + f.Type.String() + " {\n")
				}
				jsonFactory(f.Type, lineGap, tag, write)
				write(lineGap + "}\n")
				continue
			}

			if optionTag == "" {
				write(lineGap + tagname + " : " + f.Type.String() + "\n")
			} else {
				write(lineGap + tagname + " : " + f.Type.String() + "(" + optionTag + ")\n")
			}

		} else {
			if f.Type.Kind() == reflect.Struct {
				if f.Name == f.Type.Name() { //Embedding
					write(lineGap + f.Type.String() + " {\n")
					jsonFactory(f.Type, lineGap, tag, write)
					write(lineGap + "}\n")
					continue
				}

			}

			write(lineGap + strings.TrimSpace(f.Name) + " : " + f.Type.String() + "\n")
		}
	} //for
}
