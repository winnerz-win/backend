package inf

import (
	"jtools/dbg"
	"txscheduler/brix/tools/database/mongo"
)

// Test4004DB :
func Test4004DB() *mongo.CDB {
	SetConfig(GetTestConfig())
	cdb = mongo.New("192.168.0.73:40004", true, "", "test_id", "test_pwd")
	DBName = config.DB
	TokenInfo{}.IndexingDB()
	return cdb
}

func GetTestConfig(is_mainnet ...bool) *IConfig {
	mainnet := false
	if dbg.IsTrue(is_mainnet) {
		mainnet = true
	}
	config := &IConfig{
		Mainnet: mainnet,
		Version: "2021.04.13",
		Seed:    "txscheduler/seed_text...",
		DB:      "txm_mma",
		IPCheck: false,
		ClientHost: map[bool][]string{
			false: []string{"http://127.0.0.1", "8080"},
		},
		AdminSalt: "txm1234PWD",

		Confirms: 1,

		Masters: KeyPairList{
			{
				Mainnet:    false,
				PrivateKey: "private_key",
				Address:    "address",
			},
		},
		Chargers: KeyPairList{
			{
				Mainnet:    false,
				PrivateKey: "private_key",
				Address:    "address",
			},
		},
		Tokens: TokenInfoList{
			{
				Mainnet:  false,
				Contract: "eth",
				Symbol:   "ETH",
				Decimal:  "18",
			},
			{
				Mainnet:  false,
				Contract: "0x55d820d980959fbcbcdcf7ebe361dae33a71d387",
				Symbol:   "ERCT",
				Decimal:  "18",
			},
			{
				Mainnet:  false,
				Contract: "0x014c406be26ec65d7a6fd38a2195111b32f74eef",
				Symbol:   "USDT",
				Decimal:  "6",
			},
		},
		InfuraKeys: []string{
			"key_1",
			"key_1",
			"key_1",
		},
		ESKeys: []string{
			"key_1",
			"key_1",
			"key_1",
			"key_1",
		},
	}

	return config
}
