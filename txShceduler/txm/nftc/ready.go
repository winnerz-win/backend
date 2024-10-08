package nftc

import (
	"jcloudnet/itype"
	"jtools/cloud/ebcm"
	"jtools/cloud/jeth/ecs"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/runtext"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

type NftCallback struct {
	NftTokenContract  string
	NftOwnerAddress   string
	NftOwnerPrivate   string
	NftDepositAddress string

	Buy      NftBuyResultCallback
	Transfer NftTransferResultCallback
}

func NewNftCallback(
	tokenContract string,
	owerAddress string,
	ownerPrivate string,
	depositAddress string,

	buy NftBuyResultCallback,
	transfer NftTransferResultCallback,

) *NftCallback {

	tokenContract = dbg.TrimToLower(tokenContract)
	owerAddress = dbg.TrimToLower(owerAddress)

	return &NftCallback{
		tokenContract,
		owerAddress,
		ownerPrivate,
		depositAddress,
		buy,
		transfer,
	}
}

func Ready(callback *NftCallback) {
	dbg.YellowBG("[ NFT READY ]")
	readyConfig(
		callback.NftTokenContract,
		callback.NftOwnerAddress,
		callback.NftOwnerPrivate,
		callback.NftDepositAddress,
	)
	readyCMD()

	//first ----
	model.NftAset{}.IndexingDB(firstSetDB)
	mongo.StartIndexingDB(
		//
		model.NftCache{},
		model.NftTx{},
		model.NftTxLog{},

		model.NftDepositTry{},
		model.NftBuyTry{},
		model.NftBuyEnd{},
		model.NftTransferTry{},
		model.NftTransferEnd{},
	)

	if inf.LOCALMODE {
		return
	}

	go runCollection()
	go runBuyTry()
	go runBuyResult(callback.Buy)
	go runTransferResult(callback.Transfer)
	go runDepositTry()
}

func Start(classic *chttp.Classic) runtext.Starter {
	rtx := runtext.New("nftc")

	readyHandlers(classic)

	//done := make(chan struct{})
	//go runCollection()
	//<-done

	return rtx
}

func Finder() *itype.IClient {
	return itype.New(ecs.RPC_URL(inf.Mainnet()), false, inf.InfuraKey())
}

func Sender() *ebcm.Sender {
	s := Finder()
	sender := s.EBCMSender(ecs.TxSigner{})
	return sender
}

func _logtag(tag string, a ...interface{}) []interface{} {
	logs := []interface{}{tag}
	logs = append(logs, a...)
	return logs
}

func cYellow(a ...interface{}) { dbg.Yellow(_logtag("[NFT]", a...)...) }
func cPurple(a ...interface{}) { dbg.Purple(_logtag("[NFT]", a...)...) }
func cGreen(a ...interface{})  { dbg.Green(_logtag("[NFT]", a...)...) }
func cCyan(a ...interface{})   { dbg.Cyan(_logtag("[NFT]", a...)...) }
func cRed(a ...interface{})    { dbg.Red(_logtag("[NFT]", a...)...) }
