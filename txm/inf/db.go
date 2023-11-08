package inf

import (
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/jcfg"
	"txscheduler/brix/tools/jpath"
)

var (
	DBName = ""
)

const (
	LOCALMODE = false

	XMigrate = "x_migrate"
	XLog     = "x_log"

	TXETHCount       = "tx_eth_cnt"
	TXETHBlock       = "tx_eth_block"
	TXETHCharger     = "tx_eth_charger"     //가스차져
	TXETHDepositLog  = "tx_eth_deposit_log" //입금로그
	TXETHWithdraw    = "tx_eth_withdraw"    //출금대기열
	TXETHInternalCnt = "tx_eth_internal_cnt"
	TXETHTokenEx     = "tx_eth_token_ex" //기타토큰

	TXETHMasterOut    = "tx_eth_master_out"     //마스터계좌 출금
	TXETHMasterOutTry = "tx_eth_master_out_try" //

	//////////////////////////////////////////////////////
	COLConSum          = "coin_sum"
	COLCoinDay         = "coin_day"
	COLMember          = "member"
	COLInfoDeposit     = "info_deposit"
	COLInfoMaster      = "info_master"
	COLLogDeposit      = "log_deposit"
	COLLogWithdraw     = "log_withdraw"
	COLLogWithdrawSELF = "log_withdraw_self"
	COLLogToMaster     = "log_tomaster"
	COLLogExMaster     = "log_exmaster"

	COLSyncCoin = "sync_coin"

	CTX_USER = "ctx_user"

	//////////////////////////////////////////////////////
	COLAdmin = "admin"

	//////////////////////////////////////////////////////
	NFTASET        = "nft_aset"
	NFTTX          = "nft_tx"
	NFTCache       = "nft_cache"
	NFTTxLog       = "nft_txlog"
	NFTDepositTry  = "nft_deposit_try"
	NFTBuyTry      = "nft_buy_try"
	NFTBuyEnd      = "nft_buy_end"
	NFTTokenID     = "nft_tokenid"
	NFTTransferTry = "nft_transfer_try"
	NFTTransferEnd = "nft_transfer_end"
)

func NFTList() []string {
	return []string{
		NFTASET,
		NFTTX,
		NFTCache,
		NFTTxLog,
		NFTDepositTry,
		NFTBuyTry,
		NFTBuyEnd,
		NFTTokenID,
		NFTTransferTry,
		NFTTransferEnd,
	}
}

var cdb *mongo.CDB

// DB :
func DB() *mongo.CDB {
	return cdb
}

// InitMongo : TokenInfo{}.IndexingDB()
func InitMongo(filename string) {
	cfg := &mongo.Config{}
	jcfg.LoadYAML(jpath.NowPath()+"/"+filename, cfg)
	cdb = mongo.NewConfig(cfg)

	DBName = config.DB
	if DBName == "" {
		panic("inf.DBName is nil ...")
	}

	TokenInfo{}.IndexingDB()
}

// DBCollection :
func DBCollection(collection string, callback func(c mongo.Collection)) {
	DB().Run(DBName, collection, func(c mongo.Collection) {
		callback(c)
	})
}
