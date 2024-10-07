package cnet

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"jtools/jnet/chttp"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Dialer :
type Dialer struct {
	*websocket.Dialer
	address         string
	isHTTPS         bool
	isHTTPSCertKeys bool
	rootCAs         *x509.CertPool
	header          http.Header

	conmu sync.Mutex
	*dialerSession
}

// newSocket :
func newDialer(address string) *Dialer {
	checkHTTPS := false

	if strings.HasPrefix(address, "http://") {
		address = "ws://" + strings.TrimPrefix(address, "http://")

	} else if strings.HasPrefix(address, "https://") {
		address = "wss://" + strings.TrimPrefix(address, "https://")
		checkHTTPS = true
	} else {
		if strings.HasPrefix(address, "wss://") {
			checkHTTPS = true
		}
	}

	dialer := &Dialer{
		address: address,
		Dialer:  websocket.DefaultDialer,
		isHTTPS: checkHTTPS,
		header:  http.Header{},
	}
	return dialer
}

// NewSocket :
func NewSocket(address string, sslCertPath ...string) *Dialer {
	dialer := newDialer(address)
	if len(sslCertPath) > 0 {
		dialer.setCertPath(sslCertPath[0])
	}
	return dialer
}

// HTTPSCertKey :  https://forfuncsake.github.io/post/2017/08/trust-extra-ca-cert-in-go-app/
func (my *Dialer) setCertPath(cert string) {
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	certs, err := ioutil.ReadFile(cert)
	if err != nil {
		chttp.LogError("cnet.HTTSCertKey :", err)
		return
	}
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		chttp.LogError("cnet.HTTPCertKey - No certs appended, using system certs only")
		return
	}
	chttp.LogPurple("cnet.Cert.Pem Apply")
	my.isHTTPSCertKeys = true
	my.rootCAs = rootCAs

}

// SetHeader :
func (my *Dialer) SetHeader(k, v string) *Dialer {
	my.header.Add(k, v)
	return my
}

// IsConnect :
func (my *Dialer) IsConnect() bool {
	if my.dialerSession == nil {
		return false
	}
	return my.dialerSession.IsConnect()
}

// SetTimeout :
func (my *Dialer) SetTimeout(sec int) {
	defer my.conmu.Unlock()
	my.conmu.Lock()
	if sec > 1 {
		my.HandshakeTimeout = time.Second * time.Duration(sec)
	}
}

// Connect :
func (my *Dialer) Connect(api string, callback ...func(res Responser)) error {
	defer my.conmu.Unlock()
	my.conmu.Lock()

	if my.IsConnect() {
		return errors.New("alive session")
	}

	if my.isHTTPS {
		if my.isHTTPSCertKeys {
			my.TLSClientConfig = &tls.Config{RootCAs: my.rootCAs}
		} else {
			my.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}
	}
	address := my.address + api
	conn, response, err := my.Dial(address, my.header)
	if err != nil {
		chttp.LogError("[", address, "] Dial :", err)
		return err
	}
	if response != nil && response.Body != nil {
		func(response *http.Response) {
			defer response.Body.Close()
			nRes := newResponser(response)
			if len(callback) > 0 && callback[0] != nil {
				callback[0](nRes)
			}
		}(response)
	}

	my.dialerSession = newSessionD(conn)
	return nil
}
