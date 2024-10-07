package dbg

import (
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/yaml.v2"
)

func ToJsonString(v interface{}) string {
	if v == nil {
		return ""
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err.Error()
	} else if !json.Valid(b) {
		return fmt.Sprintln("Invalid ToJsonString :", v)
	}
	return fmt.Sprintln(string(b))
}
func ToJsonStringFlat(v interface{}) string {
	if v == nil {
		return ""
	}
	b, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	} else if !json.Valid(b) {
		return fmt.Sprintln("Invalid ToJsonString :", v)
	}
	return fmt.Sprintln(string(b))
}

func ToYamlString(v interface{}) string {
	b, _ := yaml.Marshal(v)
	return string(b)
}

func ParseStruct(src, dst interface{}) error {
	if src == nil {
		return Error("src is nil")
	}
	if dst == nil {
		return Error("dst is nil")
	}

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
			//fmt.Println(string(b))
			if err := json.Unmarshal(b, &r); err != nil {
				return err
			}
		}
		return nil
	}

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

func JsonToBytes(j interface{}) []byte {
	b, err := json.Marshal(j)
	if err != nil {
		return nil
	}
	return b
}

func BoolsOne(vals ...bool) bool {
	return len(vals) > 0 && vals[0]
}
