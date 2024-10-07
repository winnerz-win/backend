package ecsx

/*
	RPC_KOVAN_URL:
	Mainnet[10] : https://mainnet.optimism.io
	Mainnet[10] : https://optimism-mainnet.gateway.pokt.network/v1/lb/62eb567f0fd618003965da18

	Testnet[69] : https://kovan.optimism.io
*/
func RPC_KOVAN_URL(mainnet bool) string {
	if mainnet {
		return "https://mainnet.optimism.io" //10
	}
	return "https://kovan.optimism.io" //69
}

func KovanLayer2(mainnet bool, infuraKey string) *Sender {
	tags := []string{"eth", "ETH", "ethereum"}
	sender := NewMainnet(RPC_KOVAN_URL(mainnet), "", tags)
	sender.mainnet = mainnet
	return sender
}

const (
	Layer2ETH = "0x4200000000000000000000000000000000000006" //WETH
)

/*
	https://kovan-explorer.optimism.io/

	https://gateway.optimism.io/
*/
