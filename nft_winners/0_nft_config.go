package nft_winners

import (
	"errors"
	"txscheduler/brix/tools/cloudx/ebcmx"
	"txscheduler/brix/tools/dbg"
	"txscheduler/nft_winners/nwtypes"
	"txscheduler/nft_winners/rpc"
)

// nwdb.NftAInfo
type NftConfig struct {
	NftContract string `json:"nft_contract"` //0x1FAa080C0e0C7b94D8571D236b182f15e0C1742a
	NftSymbol   string `json:"nft_symbol"`
	NftName     string `json:"nft_name"`
	BaseURI     string `json:"base_uri"`
}

func (my NftConfig) String() string { return dbg.ToJSONString(my) }

func (my NftConfig) Reader() rpc.IReader { return rpc.Reader(my.NftContract) }

func (my *NftConfig) ValidAddress() error {
	if ok := ebcmx.IsAddressP(&my.NftContract); !ok {
		return errors.New("NFT_CONTRACT_ADDRESS_FORMAT_INVALID")
	}
	dbg.YellowItalic("nft_winners.nft_config.nft_contract :", my.NftContract)
	return nil
}

func (my NftConfig) NftInfo(token_id string) nwtypes.NftInfo {
	return nwtypes.NftInfo{
		Contract: my.NftContract,
		Symbol:   my.NftSymbol,
		TokenId:  token_id,
	}
}

var (
	nft_config NftConfig
)

func InjectNftConfig(config NftConfig) error {
	nft_config = config
	return nft_config.ValidAddress()
}
