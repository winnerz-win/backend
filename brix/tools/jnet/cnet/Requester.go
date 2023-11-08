package cnet

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"txscheduler/brix/tools/dbg"
)

//TLSVersion :
type TLSVersion string

const (
	TLSDefault TLSVersion = ""
	TLS10      TLSVersion = "TLS 1.0"
	TLS11      TLSVersion = "TLS 1.1"
	TLS12      TLSVersion = "TLS 1.2"
	TLS13      TLSVersion = "TLS 1.3"
)

//Requester :
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

//SetTLS :
func (my *Requester) SetTLS(tls TLSVersion) { my.tlsVersion = tls }
func (my Requester) getTLS() uint16 {
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

//GetCookies :
func (my *Requester) GetCookies() []http.Cookie {
	return my.cookies
}

func new(address string) *Requester {
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

//New :
func New(address string, sslCertPath ...string) *Requester {
	requester := new(address)
	if len(sslCertPath) > 0 {
		requester.setCertPath(sslCertPath[0], "")
	}
	return requester
}

//SetTimeout :
func (my *Requester) SetTimeout(d time.Duration) *Requester {
	my.TimeOut = d
	return my
}

//NewCertString :
func NewCertString(address string, certString string) *Requester {
	return NewCert(address, []byte(certString))
}

//NewCert :  https://forfuncsake.github.io/post/2017/08/trust-extra-ca-cert-in-go-app/
func NewCert(address string, cert []byte) *Requester {
	requester := new(address)
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	if ok := rootCAs.AppendCertsFromPEM(cert); ok == false {
		dbg.Red("cnet.HTTPCertKey - No certs appended, using system certs only")
	} else {
		requester.isHTTPSCertKeys = true
		requester.rootCAs = rootCAs
	}
	return requester
}

//HTTPSCertKey :  https://forfuncsake.github.io/post/2017/08/trust-extra-ca-cert-in-go-app/
func (my *Requester) setCertPath(cert, key string) {
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	certs, err := ioutil.ReadFile(cert)
	if err != nil {
		dbg.Red("cnet.HTTSCertKey :", err)
		return
	}
	if ok := rootCAs.AppendCertsFromPEM(certs); ok == false {
		dbg.Red("cnet.HTTPCertKey - No certs appended, using system certs only")
		return
	}
	my.isHTTPSCertKeys = true
	my.rootCAs = rootCAs

}

//SetHeader :
func (my *Requester) SetHeader(key, val string) *Requester {
	my.cheader[key] = val
	return my
}

//ClearHeader :
func (my *Requester) ClearHeader() *Requester {
	my.cheader = map[string]string{}
	return my
}

//RequestURL :
func (my *Requester) RequestURL() string {
	return my.reqURL
}

//ViewFail :
func (my *Requester) ViewFail(f ...Responser) {
	dbg.Red(my.reqURL, my.reqBody)
	if len(f) > 0 {
		f[0].ViewFail()
	}
}
