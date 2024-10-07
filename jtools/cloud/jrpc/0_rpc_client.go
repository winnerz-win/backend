package jrpc

import (
	"jtools/cloud/ebcm"
	"jtools/dbg"
)

type Client struct {
	debug_call bool
	Sender     func() *ebcm.Sender
	SignTooler func(prefix_messge ebcm.MessagePrefix) ebcm.SignTool
	Wallet     func(private string) (ebcm.IWallet, error)
}

func New(sender *ebcm.Sender, debug_call ...bool) *Client {
	client := &Client{
		debug_call: dbg.IsTrue(debug_call),
		Sender: func() *ebcm.Sender {
			return sender
		},
		SignTooler: func(prefix_messge ebcm.MessagePrefix) ebcm.SignTool {
			return sender.SignTooler(prefix_messge)
		},
		Wallet: func(private string) (ebcm.IWallet, error) {
			return sender.Wallet(private)
		},
	}
	return client
}

func TrimAddress(address *string) {
	*address = dbg.TrimToLower(*address)
}

/////////////////////////////////////////////////////////////////////////////
