package zzp

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type fStack []interface{}

func (my *fStack) Push(v interface{}) {
	*my = append(*my, v)
}
func (my *fStack) Pop() interface{} {
	idx := len(*my) - 1
	if idx < 0 {
		return nil
	}
	item := (*my)[idx]
	*my = (*my)[:idx]
	return item
}
func (my fStack) Count() int    { return len(my) }
func (my fStack) IsEmpty() bool { return len(my) == 0 }

func (my fStack) First() interface{} {
	if len(my) == 0 {
		return nil
	}
	return my[0]
}

//IsFirstOne : count == 1 && [0] == v
func (my fStack) IsFirstOne(v interface{}) bool {
	if my.Count() > 1 {
		return false
	}
	if my.IsEmpty() {
		return false
	}
	fv := my.First()
	return fv == v
}

func (my fStack) Cmp(v interface{}) bool {
	if len(my) == 0 {
		return false
	}
	return my[len(my)-1] == v
}
func (my fStack) ContainCount(vs ...interface{}) int {
	if len(my) == 0 {
		return 0
	}
	cnt := 0
	for _, item := range my {
		for _, v := range vs {
			if item == v {
				cnt++
			}
		}
	} //for
	return cnt
}

func toJsonString(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

func Error(a ...interface{}) error {
	msg := []string{}
	for _, v := range a {
		msg = append(msg, fmt.Sprintf("%v", v))
	} //for
	return errors.New(strings.Join(msg, " "))
}
