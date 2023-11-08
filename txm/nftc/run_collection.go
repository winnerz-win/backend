package nftc

import (
	"sync"
	"time"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

var _startNumber string = model.ZERO

func runCollection() {
	defer dbg.PrintForce("nftc.runCollection ----------  END")
	dbg.PrintForce("nftc.runCollection ----------  START")

	aset := model.NftAset{}.FirstLoadDB(_startNumber)
	if aset.IsEnd {
		setGetAllow(200)
		return
	}

	txCnt := model.NewTxETHCounter(inf.Mainnet(), aset.Number)
	txCnt.LoadFromDB(true)

	const (
		sleepNormal = time.Second * 1
	)

	for {
		if jmath.CMP(aset.Number, txCnt.Number) >= 0 {
			aset.IsEnd = true
			aset.UpdateEndFlagDB()
			break
		}

		finder := Sender()
		finder.SetCustomMethods(cms...)
		number := aset.Number

		data := finder.BlockByNumber(number, false)
		if data == nil {
			time.Sleep(sleepNormal)
			continue
		}

		GetTxList(number, data.TxList, false)

		aset.Number = jmath.ADD(aset.Number, 1)
		aset.UpdateIncDB()

		time.Sleep(time.Millisecond * 100)
	} //for

	setCheckCache()
}

var (
	exGetState = 0
	getMu      sync.RWMutex
)

func getAllowState(db mongo.DATABASE, number string, txlist ecsx.TransactionBlockList) bool {
	defer getMu.RUnlock()
	getMu.RLock()
	if exGetState == 200 {
		return true
	}

	cYellow("Cache.Insert")
	item := model.NftCache{
		Number: jmath.Int64(number),
		List:   txlist,
	}
	item.InsertDB(db)

	return false
}

func setGetAllow(v int) {
	defer getMu.Unlock()
	getMu.Lock()
	exGetState = v
}

func GetTxList(number string, txlist ecsx.TransactionBlockList, isOuter bool) {
	model.DB(func(db mongo.DATABASE) {
		if isOuter {
			if !getAllowState(db, number, txlist) {
				return
			}
		}

		if isOuter {
			cYellow("parsingTxlist :", number, ", count :", len(txlist))
		}
		parsingTxlist(db, number, txlist)
	})

}

func setCheckCache() {
	defer dbg.YellowBG(">>>>>> NFT start CheckCache <<<<<<<  END")
	dbg.YellowBG(">>>>>> NFT start CheckCache <<<<<<<  START")

	defer getMu.Unlock()
	getMu.Lock()

	model.DB(func(db mongo.DATABASE) {
		iter := db.C(inf.NFTCache).Find(nil).Sort("number").Iter()
		v := model.NftCache{}
		for iter.Next(&v) {
			parsingTxlist(
				db,
				jmath.VALUE(v.Number),
				v.List,
			)
			v.RemoveDB(db)
		} //for
	})

	exGetState = 200

}
