package ecs

import "txscheduler/brix/tools/cloud/ebcm"

const (
	ZERO = "0"

	MainnetURL      = "https://mainnet.infura.io/v3" //id:1
	MainnetOptimism = "https://mainnet.optimism.io"  //id:10

	GoerliURL      = "https://goerli.infura.io/v3" // id:5
	GoerliOptimism = "https://goerli.optimism.io"  // id:420

	TestnetURL = GoerliURL //"https://ropsten.infura.io/v3" //Ropsten --id:3

	Kovan         = "https://kovan.infura.io/v3" //id:42
	KovanOptimism = "https://kovan.optimism.io"  //id:69

	/*
		Sepolia
		rpc-url : https://rpc.sepolia.org
		ChainID : 11155111
	*/
	SepoilaURL = "https://rpc.sepolia.org" //https://sepoliafaucet.net/		--id:11155111

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

)

func RPC_URL(mainnet bool) string {
	if mainnet {
		return MainnetURL
	}
	return GoerliURL
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
