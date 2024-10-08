package nftc

import (
	"txscheduler/brix/tools/console"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func readyCMD() {
	cmdInfo()
}

func cmdInfo() {
	console.AppendCmd(
		"nft.info",
		"nft.info",
		true,
		func(ps []string) {
			console.Yellow("NFT Contract :", nftToken.Contract)
			console.Yellow("NFT Owner    :", nftToken.Address)
			console.Yellow("NFT Deposit  :", depositAddress())
			console.Atap()

			depositETH := Finder().GetCoinPrice(depositAddress())

			console.Yellow("Deposit ETH :", depositETH)

			tokens := inf.TokenList()
			for _, token := range tokens {
				if token.Symbol == model.ETH {
					continue
				}

				tokenPrice := Finder().Price(depositAddress(), token.Contract, token.Decimal)
				console.Yellow("Deposit", token.Symbol, ":", tokenPrice)

			} //for

		},
	)
}
