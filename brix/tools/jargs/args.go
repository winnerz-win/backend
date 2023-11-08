package jargs

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"txscheduler/brix/tools/dbg"
)

//ArgValue :
type ArgValue struct {
	val string
}

//String :
func (my ArgValue) String() string {
	return strings.TrimSpace(my.val)
}

//Int :
func (my ArgValue) Int() (int, error) {
	v, err := strconv.ParseInt(my.val, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

//SetInt :
func (my ArgValue) SetInt(pv *int) {
	if v, err := my.Int(); err == nil {
		*pv = v
	}
}

//ForceInt :
func (my ArgValue) ForceInt() int {
	v, _ := my.Int()
	return v
}

//Float :
func (my ArgValue) Float() (float64, error) {
	v, err := strconv.ParseFloat(my.val, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

//SetFloat :
func (my ArgValue) SetFloat(pv *float64) {
	if v, err := my.Float(); err == nil {
		*pv = v
	}
}

//ForceFloat :
func (my ArgValue) ForceFloat() float64 {
	v, _ := my.Float()
	return v
}

//Bool :
func (my ArgValue) Bool() bool {
	return strings.ToLower(my.val) == "true"
}

//ArgAction :
type ArgAction interface {
	String() string
	Action(key string, f func(v ArgValue)) error
	Word(word string, isLower ...bool) bool
	Contains(word string, isLower ...bool) bool
	ContainsCall(word string, f func(fullName string)) bool
	Next(word string, f func(next string)) bool
	Do(arg string) bool
	DoCall(arg string, call func())
}

//nArgAction :
type nArgAction struct {
	explictKey bool
	keyValue   map[string]ArgValue // key[=]value
	words      map[string]struct{} // key[=]value != > word
	argStrings []string            // all word
}

func (my nArgAction) String() string {
	m := map[string]interface{}{}
	m["explictKey"] = my.explictKey
	m["keyValue"] = my.keyValue
	m["words"] = my.words
	m["argStrings"] = my.argStrings
	return dbg.ToJSONString(m)
}

//Action :
func (my nArgAction) Action(key string, f func(v ArgValue)) error {
	if my.explictKey == false {
		key = strings.ToLower(key)
	}
	if v, isDo := my.keyValue[key]; isDo {
		f(v)
		return nil
	}
	return errors.New("not found key")
}

//Word : (word string, isLower ...bool) bool
func (my nArgAction) Word(word string, isLower ...bool) bool {
	if len(isLower) > 0 && isLower[0] == true {
		word = strings.ToLower(word)
		for key, _ := range my.words {
			key = strings.ToLower(key)
			if word == key {
				return true
			}
		} //for
	}
	_, isfind := my.words[word]
	return isfind
}

//Contains :
func (my nArgAction) Contains(word string, isLower ...bool) bool {
	lower := false
	if len(isLower) > 0 && isLower[0] == true {
		lower = true
		word = strings.ToLower(word)
	}
	for _, s := range my.argStrings {
		if lower {
			s = strings.ToLower(s)
		}
		if strings.Contains(s, word) {
			return true
		}
	} //for
	return false
}

func (my nArgAction) ContainsCall(word string, f func(fullName string)) bool {
	word = strings.ToLower(word)
	for _, s := range my.argStrings {
		cmp := strings.ToLower(s)
		if strings.Contains(cmp, word) {
			f(s)
			return true
		}
	} //for

	return false
}

func (my nArgAction) Next(word string, f func(next string)) bool {
	word = strings.ToLower(word)
	for i, s := range my.argStrings {
		cmp := strings.ToLower(s)
		if word == cmp {
			if i+1 >= len(my.argStrings) {
				return false
			}
			f(my.argStrings[i+1])
			return true
		}
	} //for
	return false
}

//Do : ARG , arg == true
func (my nArgAction) Do(arg string) bool {
	arg = dbg.TrimToLower(arg)
	for _, cmpArg := range my.argStrings {
		cmpArg = dbg.TrimToLower(cmpArg)
		if arg == cmpArg {
			return true
		}
	} //for
	return false
}

//DoCall :
func (my nArgAction) DoCall(arg string, call func()) {
	if my.Do(arg) == true {
		if call != nil {
			call()
		}
	}
}

/*
Args : div is "=" --> port=3000 ["port"]"3000"
explictKey = false : 대소문자 상관 없이 처리됨.
explictKey = true : 대소문자 구분
*/
func Args(div string, explictKey ...bool) ArgAction {
	args := nArgAction{
		keyValue: map[string]ArgValue{},
		words:    map[string]struct{}{},
	}
	if len(explictKey) > 0 && explictKey[0] == true {
		args.explictKey = true
	}

	if len(os.Args) > 1 {
		path := true
		for _, arg := range os.Args {
			if path {
				path = false
				continue
			}

			args.argStrings = append(args.argStrings, strings.Split(arg, " ")...)

			ss := strings.Split(arg, div)
			if len(ss) == 2 {
				if args.explictKey == false {
					ss[0] = strings.ToLower(ss[0])
				}
				args.keyValue[ss[0]] = ArgValue{val: ss[1]}
			} else if len(ss) == 1 {
				ss[0] = strings.TrimSpace(ss[0])
				if ss[0] == "" {
					continue
				}
				args.words[ss[0]] = struct{}{}
			}
		} //for
	}
	return args
}
