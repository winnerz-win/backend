package jyml

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func LoadFile(file_name string, p interface{}) error {

	buf, err := ioutil.ReadFile(file_name)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(buf, p)

	return err
}

func Load(target interface{}, p interface{}) error {
	switch v := target.(type) {
	case string:
		return yaml.Unmarshal([]byte(v), p)

	case []byte:
		return yaml.Unmarshal(v, p)

	} //switch

	buf, err := yaml.Marshal(target)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(buf, p)
}

func ToBytes(p interface{}) []byte {

	return nil
}
