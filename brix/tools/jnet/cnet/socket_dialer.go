package cnet

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"txscheduler/brix/tools/dbg"

	"github.com/gorilla/websocket"
)

//Dialer :
type Dialer struct {
	*websocket.Dialer
	address         string
	isHTTPS         bool
	isHTTPSCertKeys bool
	rootCAs         *x509.CertPool
	header          http.Header

	conmu sync.Mutex
	*dSession
}

//newSocket :
func newDialer(address string) *Dialer {
	address = strings.TrimSpace(address)
	checkHTTPS := false
	if strings.HasPrefix(address, "wss://") {
		checkHTTPS = true
	} else if !strings.HasPrefix(address, "ws://") {
		if strings.HasPrefix(address, "http://") {
			address = "ws" + address[4:]
		} else if strings.HasPrefix(address, "https://") {
			address = "wss" + address[5:]
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

//NewSocket :
func NewSocket(address string, sslCertPath ...string) *Dialer {
	dialer := newDialer(address)
	if len(sslCertPath) > 0 {
		dialer.setCertPath(sslCertPath[0])
	}
	return dialer
}

//HTTPSCertKey :  https://forfuncsake.github.io/post/2017/08/trust-extra-ca-cert-in-go-app/
func (my *Dialer) setCertPath(cert string) {
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
	dbg.Purple("cnet.Cert.Pem Apply")
	my.isHTTPSCertKeys = true
	my.rootCAs = rootCAs

}

//SetHeader :
func (my *Dialer) SetHeader(k, v string) *Dialer {
	my.header.Add(k, v)
	return my
}

//IsConnect :
func (my *Dialer) IsConnect() bool {
	if my.dSession == nil {
		return false
	}
	return my.dSession.IsConnect()
}

//SetTimeout :
func (my *Dialer) SetTimeout(sec int) {
	defer my.conmu.Unlock()
	my.conmu.Lock()
	if sec > 1 {
		my.HandshakeTimeout = time.Second * time.Duration(sec)
	}
}

//Connect :
func (my *Dialer) Connect(api string, callback ...func(res Responser)) error {
	defer my.conmu.Unlock()
	my.conmu.Lock()

	if my.IsConnect() {
		return errors.New("alive dSession")
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
		dbg.Red("[", address, "] Dial :", err)
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

	my.dSession = newSessionD(conn)
	return nil
}
