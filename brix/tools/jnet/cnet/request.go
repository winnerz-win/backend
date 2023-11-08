package cnet

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"txscheduler/brix/tools/dbg"

	"txscheduler/brix/tools/jnet/chttp"
)

//not used :
func (my *Requester) reqeust222(api, method, contentType string, body io.Reader, callback func(res Responser)) error {
	my.reqURL = my.Address + api
	req, err := http.NewRequest(method, my.reqURL, body)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for key, val := range my.cheader {
		req.Header.Set(key, val)
	}
	for _, cookie := range my.cookies {
		req.AddCookie(&cookie)
	} //for

	if err != nil {
		return err
	}
	client := http.DefaultClient
	client.Timeout = my.TimeOut

	if my.isHTTPS {
		var tr *http.Transport
		if my.isHTTPSCertKeys {
			tr = &http.Transport{
				TLSClientConfig: &tls.Config{
					//InsecureSkipVerify: true,
					RootCAs: my.rootCAs,
				},
				DisableKeepAlives: true,
			}
		} else {
			tr = &http.Transport{
				TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
				DisableKeepAlives: true,
			}
		}
		tr.Dial = (&net.Dialer{
			KeepAlive: 600 * time.Second,
		}).Dial
		tr.MaxIdleConns = 100
		tr.MaxIdleConnsPerHost = 100
		client.Transport = tr
		client.Transport = tr
	} else {
		client.Transport = &http.Transport{
			Dial: (&net.Dialer{
				KeepAlive: 600 * time.Second,
			}).Dial,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			DisableKeepAlives:   true,
		}
	}

	response, err := client.Do(req)
	if response != nil {
		if response.Body != nil {
			defer response.Body.Close()
		}
	}
	if err != nil {
		return err
	}

	if response != nil && response.Body != nil {
		my.cookies = my.cookies[:0]
		for _, cookie := range response.Cookies() {
			my.cookies = append(my.cookies, *cookie)
		}
		nResponser := newResponser(response, my.cookies...)
		callback(nResponser)
	}
	return nil
}

func (my *Requester) reqeust(api, method, contentType string, body io.Reader, callback func(res Responser)) error {
	defer my.mu.Unlock()
	my.mu.Lock()

	my.reqURL = my.Address + api
	//dbg.Cyan("url :", my.reqURL)
	req, err := http.NewRequest(method, my.reqURL, body)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for key, val := range my.cheader {
		req.Header.Set(key, val)
	}
	for _, cookie := range my.cookies {
		dbg.Purple("cookie:add")
		req.AddCookie(&cookie)
	} //for

	if err != nil {
		return err
	}
	client := &http.Client{}
	client.Timeout = my.TimeOut

	if my.isHTTPS {
		tr := &http.Transport{
			DisableKeepAlives: true,
		}
		var tlsConfig *tls.Config

		if my.isHTTPSCertKeys {
			tlsConfig = &tls.Config{
				//InsecureSkipVerify: true,
				RootCAs: my.rootCAs,
			}
		} else {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		if my.tlsVersion != TLSDefault {
			tlsConfig.MinVersion = my.getTLS()
		}
		tr.DisableKeepAlives = true
		tr.TLSClientConfig = tlsConfig

		tr.Dial = (&net.Dialer{
			KeepAlive: 600 * time.Second,
		}).Dial
		tr.MaxIdleConns = 100
		tr.MaxIdleConnsPerHost = 100
		client.Transport = tr
	} else {
		client.Transport = &http.Transport{
			Dial: (&net.Dialer{
				KeepAlive: 600 * time.Second,
			}).Dial,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		}
	}

	response, err := client.Do(req)
	if response != nil {
		if response.Body != nil {
			defer func() {
				response.Body.Close()
				response = nil
			}()
		}
	}
	if err != nil {
		return err
	}

	if response != nil && response.Body != nil {
		my.cookies = my.cookies[:0]
		for _, cookie := range response.Cookies() {
			my.cookies = append(my.cookies, *cookie)
		}
		nResponser := newResponser(response, my.cookies...)
		callback(nResponser)
	}
	return nil
}

//JSON : POST - json
func (my *Requester) JSON(api string, v interface{}, callback func(res Responser)) error {
	var reader io.Reader
	if v != nil {
		buffer, err := json.Marshal(v)
		if err != nil {
			return err
		}
		my.reqBody = string(buffer)
		reader = bytes.NewReader(buffer)
	} else {
		reader = bytes.NewReader(nil)
	}
	return my.reqeust(api, chttp.POST, "application/json", reader, callback)
}

//PTEXT : POST = text
func (my *Requester) PTEXT(api string, text string, callback func(res Responser)) error {
	my.reqBody = text
	return my.reqeust(api, chttp.POST, "text/plain", strings.NewReader(text), callback)
}

//FORM : POST - form [ JsonType.ReqeustFORMBody ]
func (my *Requester) FORM(api string, params chttp.JsonType, callback func(res Responser)) error {
	my.reqBody = params.RequestBody()
	return my.reqeust(api, chttp.POST, "application/x-www-form-urlencoded", strings.NewReader(my.reqBody), callback)
}

//FORM2 : POST - form [ JsonType.ReqeustFORMBody ]
func (my *Requester) FORM2(api string, params chttp.HTTPBody, callback func(res Responser)) error {
	my.reqBody = params.String()
	return my.reqeust(api, chttp.POST, "application/x-www-form-urlencoded", strings.NewReader(my.reqBody), callback)
}

//FORMString : POST - form [ string ]
func (my *Requester) FORMString(api string, params string, callback func(res Responser)) error {
	return my.reqeust(api, chttp.POST, "application/x-www-form-urlencoded", strings.NewReader(params), callback)
}

//GET : JsonType.RequestGETBody
func (my *Requester) GET(api string, params chttp.JsonType, callback func(res Responser)) error {
	body := ""
	if params != nil {
		body = "?" + params.RequestBody()
	}
	return my.reqeust(api+body, chttp.GET, "", nil, callback)
}

//GET2 : GET - [ string ]
func (my *Requester) GET2(api string, params chttp.HTTPBody, callback func(res Responser)) error {
	body := ""
	if params.IsEmtpy() == false {
		body = "?" + params.String()
	}
	return my.reqeust(api+body, chttp.GET, "", nil, callback)
}

//DELETE : JsonType.RequestGETBody
func (my *Requester) DELETE(api string, params chttp.JsonType, callback func(res Responser)) error {
	if params != nil {
		my.reqBody = params.RequestBody()
	}
	return my.reqeust(api, chttp.DELETE, "application/x-www-form-urlencoded", strings.NewReader(my.reqBody), callback)

	// body := ""
	// if params != nil {
	// 	body = "?" + params.RequestGETBody()
	// }
	// return my.reqeust(api+body, chttp.DELETE, "", nil, callback)
}

//DELETE2 : JsonType.RequestGETBody
func (my *Requester) DELETE2(api string, params chttp.HTTPBody, callback func(res Responser)) error {
	my.reqBody = params.String()
	return my.reqeust(api, chttp.DELETE, "application/x-www-form-urlencoded", strings.NewReader(my.reqBody), callback)

	// body := ""
	// if params != nil {
	// 	body = "?" + params.RequestGETBody()
	// }
	// return my.reqeust(api+body, chttp.DELETE, "", nil, callback)
}
