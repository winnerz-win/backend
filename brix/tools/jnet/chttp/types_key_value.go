package chttp

import (
	"fmt"
	"sort"
	"strings"
)

///////////////////////////////////////////////////////////////////////
type keyValue struct {
	key   string
	value interface{}
}

//HTTPBody :
type HTTPBody interface {
	IsEmtpy() bool
	Map(m map[string]interface{}) HTTPBody
	Set(key string, val interface{}) HTTPBody
	Cat(key string, val interface{}) HTTPBody
	String() string
}

type httpBody struct {
	params []keyValue
	cats   []keyValue
}

//MakeHTTPBody :
func MakeHTTPBody(ms ...map[string]interface{}) HTTPBody {
	instance := &httpBody{}
	if len(ms) > 0 {
		instance.Map(ms[0])
	}
	return instance
}

func (my httpBody) IsEmtpy() bool {
	return len(my.params) == 0 && len(my.cats) == 0
}
func (my *httpBody) Map(m map[string]interface{}) HTTPBody {
	for key, val := range m {
		my.params = append(my.params, keyValue{key, val})
	}
	return my
}
func (my *httpBody) Set(key string, val interface{}) HTTPBody {
	my.params = append(my.params, keyValue{key, val})
	return my
}
func (my *httpBody) Cat(key string, val interface{}) HTTPBody {
	my.cats = append(my.cats, keyValue{key, val})
	return my
}
func (my httpBody) String() string {
	if my.IsEmtpy() {
		return ""
	}
	sort.Slice(my.params, func(i, j int) bool { return my.params[i].key < my.params[j].key })
	list := []string{}
	for _, v := range my.params {
		list = append(list, fmt.Sprintf("%v=%v", v.key, v.value))
	}
	for _, v := range my.cats {
		list = append(list, fmt.Sprintf("%v=%v", v.key, v.value))
	}
	return strings.Join(list, "&")
}

///////////////////////////////////////////////////////////////////////
