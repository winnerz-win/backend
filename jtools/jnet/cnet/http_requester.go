package cnet

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"jtools/jnet/chttp"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// TLSVersion :
type TLSVersion string

const (
	TLSDefault TLSVersion = ""
	TLS10      TLSVersion = "TLS 1.0"
	TLS11      TLSVersion = "TLS 1.1"
	TLS12      TLSVersion = "TLS 1.2"
	TLS13      TLSVersion = "TLS 1.3"
)

// Requester :
type Requester struct {
	isHTTPS    bool
	tlsVersion TLSVersion
	Address    string
	cookies    []http.Cookie
	TimeOut    time.Duration

	cheader map[string]string

	isHTTPSCertKeys bool
	rootCAs         *x509.CertPool

	reqURL  string
	reqBody string

	mu sync.Mutex
}

// SetTLS :
func (my *Requester) SetTLS(tls TLSVersion) { my.tlsVersion = tls }
func (my *Requester) getTLS() uint16 {
	switch my.tlsVersion {
	case TLS10:
		return tls.VersionTLS10
	case TLS11:
		return tls.VersionTLS11
	case TLS12:
		return tls.VersionTLS12
	case TLS13:
		return tls.VersionTLS13
	} //switch
	return 0
}

// GetCookies :
func (my *Requester) GetCookies() []http.Cookie {
	return my.cookies
}

func new(address string) *Requester {
	address = strings.TrimSuffix(address, "/")
	checkHTTPS := false
	if strings.HasPrefix(address, "https://") {
		checkHTTPS = true
	}
	requester := &Requester{
		isHTTPS:    checkHTTPS,
		tlsVersion: TLSDefault,
		Address:    address,
		cheader:    map[string]string{},
		TimeOut:    10 * time.Second,
	}
	return requester
}

// New :
func New(address string, sslCertPath ...string) *Requester {
	requester := new(address)
	if len(sslCertPath) > 0 {
		requester.setCertPath(sslCertPath[0], "")
	}
	return requester
}

// SetTimeout :
func (my *Requester) SetTimeout(d time.Duration) *Requester {
	my.TimeOut = d
	return my
}

// NewCertString :
func NewCertString(address string, certString string) *Requester {
	return NewCert(address, []byte(certString))
}

// NewCert :  https://forfuncsake.github.io/post/2017/08/trust-extra-ca-cert-in-go-app/
func NewCert(address string, cert []byte) *Requester {
	requester := new(address)
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	if ok := rootCAs.AppendCertsFromPEM(cert); ok == false {
		chttp.LogError("cnet.HTTPCertKey - No certs appended, using system certs only")
	} else {
		requester.isHTTPSCertKeys = true
		requester.rootCAs = rootCAs
	}
	return requester
}

// HTTPSCertKey :  https://forfuncsake.github.io/post/2017/08/trust-extra-ca-cert-in-go-app/
func (my *Requester) setCertPath(cert, key string) {
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	certs, err := ioutil.ReadFile(cert)
	if err != nil {
		chttp.LogError("cnet.HTTSCertKey :", err)
		return
	}
	if ok := rootCAs.AppendCertsFromPEM(certs); ok == false {
		chttp.LogError("cnet.HTTPCertKey - No certs appended, using system certs only")
		return
	}
	my.isHTTPSCertKeys = true
	my.rootCAs = rootCAs

}

// SetHeader :
func (my *Requester) SetHeader(key, val string) *Requester {
	my.cheader[key] = val
	return my
}

// ClearHeader :
func (my *Requester) ClearHeader() *Requester {
	my.cheader = map[string]string{}
	return my
}

// RequestURL :
func (my *Requester) RequestURL() string {
	return my.reqURL
}

// ViewFail :
func (my *Requester) ViewFail(f ...Responser) {
	chttp.LogError(my.reqURL, my.reqBody)
	if len(f) > 0 {
		f[0].ViewFail()
	}
}

func (my *Requester) reqeust(api, method, contentType string, body io.Reader, callback func(res Responser)) error {
	defer my.mu.Unlock()
	my.mu.Lock()
	return my._reqeust(
		api,
		method,
		contentType,
		body,
		callback,
	)

}

func (my *Requester) _reqeust(api, method, contentType string, body io.Reader, callback func(res Responser)) error {
	my.reqURL = my.Address + api
	//cc.Cyan("url :", my.reqURL)
	req, err := http.NewRequest(method, my.reqURL, body)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for key, val := range my.cheader {
		req.Header.Set(key, val)
	}
	for _, cookie := range my.cookies {
		//cc.Purple("cookie:add")
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
		tr.MaxIdleConns = DefaultTransportMaxIdleConns
		tr.MaxIdleConnsPerHost = DefaultTransportMaxIdlePerHost
		client.Transport = tr
	} else {
		client.Transport = &http.Transport{
			Dial: (&net.Dialer{
				KeepAlive: 600 * time.Second,
			}).Dial,
			MaxIdleConns:        DefaultTransportMaxIdleConns,
			MaxIdleConnsPerHost: DefaultTransportMaxIdlePerHost,
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

// JSON : application/json
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
func (my *Requester) JSONX(api string, v interface{}, callback func(res Responser)) error {
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
	return my._reqeust(api, chttp.POST, "application/json", reader, callback)
}

// TEXT : text/plain
func (my *Requester) TEXT(api string, text string, callback func(res Responser)) error {
	my.reqBody = text
	return my.reqeust(api, chttp.POST, "text/plain", strings.NewReader(text), callback)
}
func (my *Requester) TEXTX(api string, text string, callback func(res Responser)) error {
	my.reqBody = text
	return my._reqeust(api, chttp.POST, "text/plain", strings.NewReader(text), callback)
}

//////////////////////////////////////////////////////////////////////////////////////////

type IParam interface {
	Body(prefix string) string
	Set(key string, val interface{}) IParam
}

type bodyParam struct {
	key []string
	val []string
}

func (my bodyParam) Valid() bool {
	return len(my.key) > 0
}
func (my bodyParam) Body(prefix string) string {
	if my.Valid() {
		u := url.Values{}

		cnt := len(my.key)
		for i := 0; i < cnt; i++ {
			k := my.key[i]
			v := my.val[i]
			u.Set(k, v)
		} //for
		return prefix + u.Encode()
	}
	return ""
}
func (my *bodyParam) Set(key string, val interface{}) IParam {
	my.key = append(my.key, key)
	my.val = append(my.val, fmt.Sprint(val))
	return my
}

func MakeParam() IParam {
	return &bodyParam{}
}

//////////////////////////////////////////////////////////////////////////////////////////

func (my *Requester) GET(api string, params IParam, callback func(res Responser)) error {
	body := ""
	if params != nil {
		body = params.Body("?")
	}
	return my.reqeust(api+body, chttp.GET, "", nil, callback)
}
func (my *Requester) GETX(api string, params IParam, callback func(res Responser)) error {
	body := ""
	if params != nil {
		body = params.Body("?")
	}
	return my._reqeust(api+body, chttp.GET, "", nil, callback)
}

func (my *Requester) WWWFORM(api string, params IParam, callback func(res Responser)) error {
	body := ""
	if params != nil {
		body = params.Body("?")
	}
	my.reqBody = body
	return my.reqeust(api, chttp.POST, "application/x-www-form-urlencoded", strings.NewReader(my.reqBody), callback)
}
func (my *Requester) WWWFORMX(api string, params IParam, callback func(res Responser)) error {
	body := ""
	if params != nil {
		body = params.Body("?")
	}
	my.reqBody = body
	return my._reqeust(api, chttp.POST, "application/x-www-form-urlencoded", strings.NewReader(my.reqBody), callback)
}
