package nftc

import (
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx"
	"txscheduler/brix/tools/dbg"
)

var (
	/*
		GoldenGoal NFT
		GDGNFT
	*/
	nftToken = struct {
		Private  string
		Address  string
		Contract string
	}{
		Private:  "private_key",
		Address:  "eoa_address",
		Contract: "nft_contract_address",
	}

	nftDepositAddress string = "nft_deposit_address"

	isRun = false

	gasSpeed = ecsx.GasFast
)

func readyConfig(
	nftContract string,
	ownerAddress string,
	ownerPrivate string,
	depositAddress string,
) {
	isRun = true

	nftToken.Contract = dbg.TrimToLower(nftContract)
	nftToken.Address = dbg.TrimToLower(ownerAddress)
	nftToken.Private = ownerPrivate

	nftDepositAddress = dbg.TrimToLower(depositAddress)

}

// depositAddress :
func depositAddress() string {
	return nftDepositAddress
	//return nftToken.Address
}

func IsRun() bool { return isRun }

func debugMode() bool {
	return true
	//return !inf.Mainnet()
}

type WriteResult struct {
	Constract string `bson:"contract" json:"contract"`
	FuncName  string `bson:"func_name" json:"func_name"`
	Hash      string `bson:"hash" json:"hash"`
	Nonce     uint64 `bson:"nonce" json:"nonce"`
	Error     error  `bson:"error" json:"error"`
	From      string `bson:"from" json:"from"`
}

func (my *WriteResult) Set(from string, h string, n uint64, e error) {
	my.From = from
	my.Hash = h
	my.Nonce = n
	my.Error = e
}
