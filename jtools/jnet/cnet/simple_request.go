package cnet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"jtools/dbg"
	"net/http"
	"time"
)

const (
	default_timeout_value = 60
)

type SimpleNetAck struct {
	url  string
	code int
	buf  []byte
	err  error
}

func (my SimpleNetAck) Error() error {
	if my.err != nil {
		return fmt.Errorf("SimpleNetAck[%s] status_code:%v, i/o timeout : %v", my.url, my.code, my.err)
	}
	if my.code >= 200 && my.code < 300 {
		return nil
	}
	return fmt.Errorf("SimpleNetAck[%s] status_code:%v, msg:%v", my.url, my.code, string(my.buf))
}
func (my SimpleNetAck) StatusCode() int { return my.code }
func (my SimpleNetAck) Bytes() []byte   { return my.buf }
func (my SimpleNetAck) ParseJson(p interface{}) error {
	return json.Unmarshal(my.buf, p)
}
func (my SimpleNetAck) Text() string {
	if len(my.buf) == 0 {
		return ""
	}
	return string(my.buf)
}

/////////////////////////////////////////////////////////////////////

func GET(url string, time_out_sec ...int) (int, []byte, error) {
	time_sec := default_timeout_value * time.Second
	if len(time_out_sec) > 0 && time_out_sec[0] > 0 {
		time_sec = time.Duration(time_out_sec[0]) * time.Second
	}

	client := &http.Client{
		// Transport: &http.Transport{
		// 	Dial: (&net.Dialer{
		// 		Timeout: timeout_value * time.Second,
		// 	}).Dial,
		// 	TLSHandshakeTimeout: timeout_value * time.Second,
		// },
		Timeout: time_sec,
	}
	client.Transport = httpTransport()

	status_code := 0

	resp, err := client.Get(url)
	defer func() {
		if resp != nil {
			if resp.Body != nil {
				io.Copy(ioutil.Discard, resp.Body)
				resp.Body.Close()
			}
		}
	}()

	var buf []byte
	if err == nil {
		buf, err = io.ReadAll(resp.Body)
		if err != nil {
			return status_code, nil, err
		}
	}

	return status_code, buf, err
}

func GET_JSON(url string, header []string, time_out_sec ...int) (int, []byte, error) {
	const timeout_value = default_timeout_value
	time_sec := timeout_value * time.Second
	if len(time_out_sec) > 0 && time_out_sec[0] > 0 {
		time_sec = time.Duration(time_out_sec[0]) * time.Second
	}
	client := &http.Client{
		Timeout: time_sec,
	}
	client.Transport = httpTransport()

	status_code := 0

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return status_code, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	for i := 0; i < len(header); i += 2 {
		req.Header.Set(header[i], header[i+1])
	}

	resp, err := client.Do(req)
	defer func() {
		if resp != nil {
			if resp.Body != nil {
				io.Copy(ioutil.Discard, resp.Body)
				resp.Body.Close()
			}
		}
	}()

	var buf []byte
	if err == nil {

		status_code = resp.StatusCode

		buf, err = io.ReadAll(resp.Body)
		if err != nil {
			return status_code, nil, err
		}
	}

	return status_code, buf, err
}

func GET_JSON_F(
	url string,
	header []string,
	time_out_sec ...int,
) SimpleNetAck {
	status_code, buf, err := GET_JSON(url, header, time_out_sec...)
	ack := SimpleNetAck{
		url:  url,
		code: status_code,
		buf:  buf,
		err:  err,
	}
	return ack
}
func POST_JSON(url string, header []string, param interface{}, time_out_sec ...int) (int, []byte, error) {
	time_sec := default_timeout_value * time.Second
	if len(time_out_sec) > 0 && time_out_sec[0] > 0 {
		time_sec = time.Duration(time_out_sec[0]) * time.Second
	}

	client := &http.Client{
		Timeout: time_sec,
	}
	client.Transport = httpTransport()

	status_code := 0

	var reader io.Reader
	if param != nil {
		buf, err := json.Marshal(param)
		if err != nil {
			return status_code, nil, err
		}
		reader = bytes.NewReader(buf)
	} else {
		reader = bytes.NewReader(nil)
	}

	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return status_code, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	for i := 0; i < len(header); i += 2 {
		req.Header.Set(header[i], header[i+1])
	}

	resp, err := client.Do(req)
	defer func() {
		if resp != nil {
			if resp.Body != nil {
				io.Copy(ioutil.Discard, resp.Body)
				resp.Body.Close()
			}
		}
	}()

	var buf []byte
	if err == nil {

		status_code = resp.StatusCode

		buf, err = io.ReadAll(resp.Body)
		if err != nil {
			return status_code, nil, err
		}
	}

	return status_code, buf, err
}

func POST_JSON_F(
	url string,
	header []string,
	param interface{},
	time_out_sec ...int,
) SimpleNetAck {
	status_code, buf, err := POST_JSON(url, header, param, time_out_sec...)
	ack := SimpleNetAck{
		url:  url,
		code: status_code,
		buf:  buf,
		err:  err,
	}
	return ack
}

/////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////

type NetStructAck[T any] struct {
	url           string
	code          int
	item          T
	err           error
	is_fail_parse bool
	fail_buf_text string
}

func (my NetStructAck[T]) IsParseFail() bool {
	return my.is_fail_parse
}
func (my NetStructAck[T]) Error() error {
	if my.err != nil {
		return fmt.Errorf("NetStructAck[%s] status_code:%v, i/o timeout : %v", my.url, my.code, my.err)
	}
	if my.code < 200 || my.code >= 300 {
		return fmt.Errorf("NetStructAck[%s] status_code:%v, msg:%v", my.url, my.code, my.fail_buf_text)
	}
	if my.is_fail_parse {
		return fmt.Errorf("NetStructAck[%s] status_code:%v, msg:%v", my.url, my.code, my.fail_buf_text)
	}

	return nil
}
func (my NetStructAck[T]) StatusCode() int { return my.code }

func (my NetStructAck[T]) Item() T { return my.item }

func POST_STRUCT[T any](
	url string,
	header []string,
	param interface{},
	time_out_sec ...int,
) NetStructAck[T] {
	status_code, buf, err := POST_JSON(url, header, param, time_out_sec...)
	ack := NetStructAck[T]{
		url:  url,
		code: status_code,
		err:  err,
	}
	if err != nil {
		if len(buf) > 0 {
			var item T
			if err := dbg.ParseStruct(buf, &item); err != nil {
				ack.is_fail_parse = true
				ack.fail_buf_text = string(buf)
			}
		}
		return ack
	}

	//cc.White("BUF_SIZE:", len(buf))

	var item T
	if err := dbg.ParseStruct(buf, &item); err != nil {
		ack.is_fail_parse = true
		ack.fail_buf_text = string(buf)
		return ack
	}

	ack.item = item
	return ack
}

func GET_STRUCT[T any](
	url string,
	header []string,
	time_out_sec ...int,
) NetStructAck[T] {
	status_code, buf, err := GET_JSON(url, header, time_out_sec...)
	ack := NetStructAck[T]{
		url:  url,
		code: status_code,
		err:  err,
	}
	if err != nil {
		if len(buf) > 0 {
			var item T
			if err := dbg.ParseStruct(buf, &item); err != nil {
				ack.is_fail_parse = true
				ack.fail_buf_text = string(buf)
			}
		}
		return ack
	}
	var item T
	if err := dbg.ParseStruct(buf, &item); err != nil {
		ack.is_fail_parse = true
		ack.fail_buf_text = string(buf)
		return ack
	}

	ack.item = item
	return ack
}
