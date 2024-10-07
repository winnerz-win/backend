package chttp

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/mholt/binding"
)

type dummy struct{}

func (d *dummy) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{}
}
func BindingDummy(req *http.Request) error {
	d := &dummy{}
	return binding.Bind(req, d)
}

func BindingJSON(req *http.Request, mapper interface{}) error {

	bodybuf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return errors.New(`{"message":"Read Body error"}`)
	}
	defer req.Body.Close()
	if err := json.Unmarshal(bodybuf, mapper); err != nil {
		LogError("BindingJSON_ERR :", string(bodybuf))
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
