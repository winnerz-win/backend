package jbsc

import (
	"jtools/cloud/ebcm"
	"jtools/cloud/jeth/ecs"
)

const (
	ZERO = "0"
	/*
		"https://bsc-dataseed.binance.org"
		"https://bsc-dataseed1.defibit.io"
		"https://bsc-dataseed1.ninicoin.io"

		gasURL:   "https://bsc-dataseed.binance.org",
		Explorer: "https://bscscan.com/",
	*/
	MainnetURL = "https://bsc-dataseed.binance.org"

	/*
		"https://data-seed-prebsc-1-s1.binance.org:8545"
		"https://data-seed-prebsc-2-s1.binance.org:8545"
		"https://data-seed-prebsc-1-s2.binance.org:8545"
		"https://data-seed-prebsc-2-s2.binance.org:8545"
		"https://data-seed-prebsc-1-s3.binance.org:8545"
		"https://data-seed-prebsc-2-s3.binance.org:8545"

		gasURL:   "https://bsc-dataseed.binance.org",
		Explorer: "https://testnet.bscscan.com/",

		Faucet : https://testnet.binance.org/faucet-smart
	*/
	TestnetURL = "https://data-seed-prebsc-2-s2.binance.org:8545"
)

func RPC_URL(mainnet bool) string {
	if mainnet {
		return MainnetURL
	}
	return TestnetURL
}

func New(host, key string, cacheId ...interface{}) *ebcm.Sender {
	sender := ecs.New(host, key, cacheId...)
	if sender != nil {
		sender.SetTXNTYPE(ebcm.TXN_LEGACY)
	}
	return sender
}
