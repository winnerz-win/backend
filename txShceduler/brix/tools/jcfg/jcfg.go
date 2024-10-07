package jcfg

import (
	"fmt"
	"io/ioutil"
	"txscheduler/brix/tools/dbg"

	"gopkg.in/yaml.v2"
)

//LoadYAML :
func LoadYAML(path string, p interface{}, isSkipPanic ...bool) error {
	skipPanic := false
	if len(isSkipPanic) > 0 && isSkipPanic[0] {
		skipPanic = true
	}
	if !skipPanic {
		fmt.Println("jcfg.Load :", path)
	}
	if buf, err := ioutil.ReadFile(path); err != nil {
		if !skipPanic {
			panic(fmt.Sprintf("[ jcfg.Load ] ReadFileError : %v", err))
		} else {
			fmt.Println("jcfg.Load :", path)
			fmt.Println("[ jcfg.Load ] not found file :", path)
		}

	} else if err := yaml.Unmarshal(buf, p); err != nil {
		dbg.Red(string(buf))
		if !skipPanic {
			panic(fmt.Sprintf("[ jcfg.Load ] UnmarshalError : %v", err))
		} else {
			fmt.Println("jcfg.Load :", path)
			fmt.Println("[ jcfg.Load ] UnmarshalError :", path)
		}
	}
	return nil
}

//ReadBytes : yaml.Unmarshal([]byte , p)
func ReadBytes(buf []byte, p interface{}) error {
	return yaml.Unmarshal(buf, p)
}
