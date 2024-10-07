package cnet

import (
	"jtools/cc"
	"jtools/dbg"
	"net/http"
)

/*
https://syntaxsugar.tistory.com/entry/GoGolang-HTTP-%EC%84%B1%EB%8A%A5-%ED%8A%9C%EB%8B%9D
*/
const (
	DefaultTransportMaxIdleConns   = 1000
	DefaultTransportMaxIdlePerHost = 1000 // default-value : 2
)

var (
	default_transport = httpTransport()
)

func httpTransport() *http.Transport {
	transport := http.DefaultTransport
	tp := transport.(*http.Transport)
	tp.MaxIdleConns = DefaultTransportMaxIdleConns
	tp.MaxIdleConnsPerHost = DefaultTransportMaxIdlePerHost
	//tp.MaxConnsPerHost = 0
	return tp
}

func HttpTransport() *http.Transport {
	return default_transport
}

func SetHttpTransport(maxIdleConns, maxIdlePerHost int) {
	default_transport.MaxIdleConns = maxIdleConns
	default_transport.MaxIdleConnsPerHost = maxIdlePerHost
	cc.YellowItalicBG(
		"< cnet.SetHttpTransport >", dbg.ENTER,
		"MaxIdleConns :", maxIdleConns, dbg.ENTER,
		"MaxIdlePerHost :", maxIdlePerHost, dbg.ENTER,
	)
}
