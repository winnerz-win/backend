package rpc

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
	privatekey    string
	contract      string
	callerAddress string

	get_wallet_func func(hexPrivate string) (ebcm.IWallet, error)
}

func (my cWriter) getPrivateKey() string { return my.privatekey }
func (my cWriter) Contract() string      { return my.contract }
func (my *cWriter) CallerAddress() string {
	if my.callerAddress == "" {
		wallet, _ := my.get_wallet_func(my.privatekey)
		my.callerAddress = wallet.Address()
	}
	return my.callerAddress
}

func Writer(
	wallet_func func(hexPrivate string) (ebcm.IWallet, error),
	contract, privatekey string,
) IWriter {
	cr := &cWriter{
		contract:        contract,
		privatekey:      privatekey,
		get_wallet_func: wallet_func,
	}
	cr.CallerAddress()
	return cr
}
