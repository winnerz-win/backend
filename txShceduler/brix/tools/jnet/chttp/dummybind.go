package chttp

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"txscheduler/brix/tools/dbg"

	"github.com/mholt/binding"
)

type dummyMapper struct {
}

func (d *dummyMapper) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{}
}

/*
BindingDummy : FieldMap 데이타를 수동으로 받기위해서 더미로 바인딩하는 함수..;;;;
단, [application/json] 으로 받아야 할 경우에는 본함수로 바인딩 하게 되면
req.Body 에서 json데이타를 읽을수 없게 된다.
*/
func BindingDummy(req *http.Request) error {
	d := &dummyMapper{}
	return binding.Bind(req, d)
}

type ParseRequestData[T any] struct {
	item T
	err  error
}

func (my ParseRequestData[T]) Error() error { return my.err }
func (my ParseRequestData[T]) Data() T      { return my.item }

func ParseRequestJson[T any](req *http.Request) T {
	var re T
	BindingJSON(req, &re)
	return re
}
func ParseRequestJson2[T any](req *http.Request) (T, error) {
	var re T
	if err := BindingJSON(req, &re); err != nil {
		return re, err
	}
	return re, nil
}

// BindingJSON :
func BindingJSON(req *http.Request, mapper interface{}) error {

	// contentType := req.Header.Get("Content-Type")
	// if strings.Contains(contentType, "json") {
	// 	fmt.Println("Content-Type_is_not_found_[application/json]....")
	// }

	bodybuf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return errors.New(`{"message":"Read Body error"}`)
	}
	defer req.Body.Close()
	// defer func() {
	// 	if e := recover(); e != nil {
	// 		dbg.Red(e)
	// 	}
	// }()
	if err := json.Unmarshal(bodybuf, mapper); err != nil {
		dbg.Red("BindingJSON_ERR :", string(bodybuf))
		return err
	}
	return nil
}

func BindingStruct[T any](req *http.Request) T {
	var re T
	BindingJSON(req, &re)
	return re
}

// BindingText :
func BindingText(req *http.Request, text *string) error {
	bodybuf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return errors.New(`{"message":"Read Body error"}`)
	}
	defer req.Body.Close()

	*text = string(bodybuf)
	return nil
}

// BindingBuffer :
func BindingBuffer(req *http.Request, f func(buf []byte)) error {
	bodybuf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return errors.New(`{"message":"Read Body error"}`)
	}
	defer req.Body.Close()

	f(bodybuf)

	return nil
}

// BindingReader : io.ReadCloser
func BindingReader(req *http.Request, f func(reader io.ReadCloser) error) error {
	defer req.Body.Close()
	return f(req.Body)
}

// func getErrorJSON(message string) error {
// 	type Er struct {
// 		Message string `json:"message"`
// 	}

// }

// Query : Request GET Query
func Query(req *http.Request) url.Values {
	return req.URL.Query()
}

// GetQuery : GET key-val
func GetQuery(req *http.Request, key string) string {
	return req.URL.Query().Get(key)
}

// GetRawQueryMap :
func GetRawQueryMap(req *http.Request) map[string]string {

	paramMap := make(map[string]string)

	seps := strings.Split(req.URL.RawQuery, "&")
	for _, div := range seps {
		vals := strings.Split(div, "=")
		if len(vals) > 1 {
			paramMap[vals[0]] = vals[1]
		} else {
			paramMap[vals[0]] = ""
		}
	} //for

	if len(paramMap) == 0 {
		return nil
	}
	return paramMap
}
