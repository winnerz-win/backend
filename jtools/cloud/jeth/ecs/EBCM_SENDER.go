package ecs

import (
	"jtools/cloud/ebcm"
)

const (
	ZERO = "0"

	MainnetURL      = "https://mainnet.infura.io/v3" //id:1
	MainnetOptimism = "https://mainnet.optimism.io"  //id:10

	Kovan         = "https://kovan.infura.io/v3" //id:42
	KovanOptimism = "https://kovan.optimism.io"  //id:69

	/*
		Sepolia
		rpc-url : https://rpc.sepolia.org
		ChainID : 11155111
		Explorer : https://sepolia.etherscan.io/
	*/
	SepoilaURL = "https://rpc.ankr.com/eth_sepolia/6b6ea5eaea29cc09d1b895224a62c829b90fb14b07ecb5ab7a0e488039d244ca" //"https://rpc.ankr.com/eth_sepolia" //"https://rpc.sepolia.org" //https://sepolia-faucet.pk910.de/		--id:11155111

	TestnetURL = SepoilaURL //id:11155111   ---- //"https://goerli.infura.io/v3" //id:5

	/*
		< Goerli >
		rpc-url : https://goerli.infura.io/v3
		ChainID : 5
		Explorer : https://goerli.etherscan.io/
		Pow-faucet : https://goerli-faucet.pk910.de/

		< Optimistic-Goerili (Layer2) >
		rpc-url : https://goerli.optimism.io
		ChainID : 420
		Explorer : https://blockscout.com/optimism/goerli/

		< Bridge-V2 >
		https://app.optimism.io/bridge

	*/
	GoerliURL      = "https://goerli.infura.io/v3" // id:5
	GoerliOptimism = "https://goerli.optimism.io"  // id:420
)

func RPC_URL(mainnet bool) string {
	if mainnet {
		return MainnetURL
	}
	return TestnetURL
}

func New(host, key string, cacheId ...interface{}) *ebcm.Sender {

	client := NewClient(host, key, cacheId...)
	if client == nil {
		return nil
	}

	sender := ebcm.NewSender(
		client,
	)
	sender.SetTXNTYPE(ebcm.TXN_EIP_1559)

	return sender

}
