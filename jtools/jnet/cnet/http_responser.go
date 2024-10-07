package cnet

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"jtools/jnet/chttp"
	"net/http"
)

// Responser :
type Responser struct {
	StatusCode int
	cookies    []http.Cookie
	isBuffer   bool
	buffer     []byte
	header     Header
}

func (my Responser) Header() Header { return my.header }

func newResponser(res *http.Response, cookies ...http.Cookie) Responser {
	my := Responser{
		StatusCode: res.StatusCode,
		header:     Header{res.Header},
	}
	if len(cookies) > 0 {
		my.cookies = append(my.cookies, cookies...)
	}
	if buffer, err := ioutil.ReadAll(res.Body); err == nil {
		my.buffer = buffer
		my.isBuffer = true
	}
	if n, err := io.Copy(ioutil.Discard, res.Body); err != nil || n != 0 {
		chttp.LogError("cnet.newResponser[", res.StatusCode, "] io.Copy[", n, "]", err)
	}
	return my
}

// Bytes :
func (my *Responser) Bytes() []byte {
	return my.buffer
}

// ToJSON :
func (my *Responser) ToJSON(v interface{}) error {
	buffer := my.Bytes()
	if buffer == nil {
		return errors.New("body buffer is nil")
	}
	return json.Unmarshal(buffer, v)
}

// Text :
func (my *Responser) Text() string {
	return my.String()
}

// String :
func (my *Responser) String() string {
	buffer := my.Bytes()
	if buffer == nil {
		return ""
	}
	return string(buffer)
}

// Code :
func (my Responser) Code() int {
	return my.StatusCode
}

// Success :
func (my Responser) Success() bool {
	return my.StatusCode >= 200 && my.StatusCode < 300
}

// GetCookies :
func (my Responser) GetCookies() []http.Cookie {
	return my.cookies
}

// ViewFail :
func (my Responser) ViewFail() error {
	failErr := chttp.Cat("[", my.StatusCode, "]", my.String())
	chttp.LogError(failErr)
	return errors.New(failErr)
}
