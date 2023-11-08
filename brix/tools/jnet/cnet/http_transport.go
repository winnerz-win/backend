package cnet

import (
	"net/http"
	"txscheduler/brix/tools/dbg"
)

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
	return tp
}

func HttpTransport() *http.Transport {
	return default_transport
}

func SetHttpTransport(maxIdleConns, maxIdlePerHost int) {
	default_transport.MaxIdleConns = maxIdleConns
	default_transport.MaxIdleConnsPerHost = maxIdlePerHost
	dbg.YellowItalicBG(
		"< cnet.SetHttpTransport >", dbg.ENTER,
		"MaxIdleConns :", maxIdleConns, dbg.ENTER,
		"MaxIdlePerHost :", maxIdlePerHost, dbg.ENTER,
	)
}
