package jrpc

import (
	"jtools/cloud/ebcm"
	"strings"
)

type IReader interface {
	Contract() string
	CallerAddress() string
}

type cReader struct {
	contract      string
	callerAddress string
}

func Reader(contract string, caller ...string) IReader {
	contract = strings.ToLower(strings.TrimSpace(contract))
	cr := cReader{
		contract: contract,
	}
	if len(caller) > 0 {
		cr.callerAddress = strings.ToLower(strings.TrimSpace(caller[0]))
	} else {
		cr.callerAddress = contract
	}
	return cr
}
func (my cReader) Contract() string      { return my.contract }
func (my cReader) CallerAddress() string { return my.callerAddress }

//////////////////////////////////////////////////////////////////////

type IWriter interface {
	getPrivateKey() string
	Contract() string
	CallerAddress() string
}

type cWriter struct {
	contract string
	wallet   ebcm.IWallet
}

func (my cWriter) getPrivateKey() string { return my.wallet.PrivateKey() }
func (my cWriter) Contract() string      { return my.contract }
func (my *cWriter) CallerAddress() string {
	return my.wallet.Address()
}

func Writer(contract string, wallet ebcm.IWallet) IWriter {
	cr := &cWriter{
		contract: contract,
		wallet:   wallet,
	}
	return cr
}

//////////////////////////////////////////////////////////////////////
