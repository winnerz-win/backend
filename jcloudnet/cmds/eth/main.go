package main

import (
	"jcloudnet/txsigner"
	"jtools/cloud/jeth/ecs"
)

func main() {

	option := txsigner.Option{
		InfraTag: "ETH",
		Signer:   ecs.TxSigner{},
	}
	txsigner.StartServer(option)

}
