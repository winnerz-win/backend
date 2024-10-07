package kly

import (
	"jtools/cloud/ebcm"
)

const (
	ZERO       = "0"
	MainnetURL = "https://api.cypress.klaytn.net:8651"
	TestnetURL = "https://api.baobab.klaytn.net:8651"
)

func RPC_URL(mainnet bool) string {
	if mainnet {
		return MainnetURL
	}
	return TestnetURL
}

func New(host, key string, cacheId ...interface{}) *ebcm.Sender {

	client := NewClient(host, cacheId...)
	if client == nil {
		return nil
	}

	sender := ebcm.NewSender(
		client,
	)

	return sender
}
