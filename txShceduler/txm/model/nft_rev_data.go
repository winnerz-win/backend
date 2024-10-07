package model

import (
	"sync"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/mms"
	"txscheduler/txm/inf"
)

var nftMu sync.Mutex

func NFTLock(f func()) {
	defer nftMu.Unlock()
	nftMu.Lock()
	f()
}

///////////////////////////////////////////////////////////////////////////////////////

type NftRevData struct {
	User              `bson:",inline" json:",inline"`
	IsExternalAddress bool   `bson:"is_external_address" json:"is_external_address"`
	ReceiptCode       string `bson:"receipt_code" json:"receipt_code"`
	PayAddress        string `bson:"pay_address" json:"pay_address"`
	PaySymbol         string `bson:"pay_symbol" json:"pay_symbol"`
	PayPrice          string `bson:"pay_price" json:"pay_price"`
	IsPayFree         bool   `bson:"is_pay_free" json:"is_pay_free"`

	GasLimit  uint64 `bson:"gasLimit" json:"gasLimit"`
	GasPrice  string `bson:"gasPrice" json:"gasPrice"`
	GasFeeETH string `bson:"gasFeeETH" json:"gasFeeETH"`

	TokenId   string `bson:"token_id" json:"token_id"`
	TokenType string `bson:"token_type" json:"token_type"`

	DepositHash      string `bson:"deposit_hash" json:"deposit_hash"`
	DepositGasLimit  uint64 `bson:"deposit_gasLimit" json:"deposit_gasLimit"`
	DepositGasPrice  string `bson:"deposit_gasPrice" json:"deposit_gasPrice"`
	DepositGasFeeETH string `bson:"deposit_gasFeeETH" json:"deposit_gasFeeETH"`
}

// PayPrivateKey : 실제 구매비용을 지불할 계정
func (my NftRevData) PayPrivateKey() string {
	if my.PayAddress == inf.Master().Address {
		return inf.Master().PrivateKey
	}
	return my.PrivateKey()
}

// Valid : receipt_code
func (my NftRevData) Valid() bool { return my.ReceiptCode != "" }

// Selector : receipt_code
func (my NftRevData) Selector() mongo.Bson { return mongo.Bson{"receipt_code": my.ReceiptCode} }

func (NftRevData) _indexingDB(c mongo.Collection, tokenIdUnique bool) {
	c.EnsureIndex(mongo.SingleIndex("uid", "1", false))
	c.EnsureIndex(mongo.SingleIndex("address", "1", false))
	c.EnsureIndex(mongo.SingleIndex("name", "1", false))

	c.EnsureIndex(mongo.SingleIndex("is_external_address", "1", false))

	c.EnsureIndex(mongo.SingleIndex("receipt_code", "1", true))
	c.EnsureIndex(mongo.SingleIndex("token_id", "1", tokenIdUnique))
	c.EnsureIndex(mongo.SingleIndex("token_type", "1", false))

	c.EnsureIndex(mongo.SingleIndex("is_pay_free", "1", false))

	c.EnsureIndex(mongo.SingleIndex("deposit_hash", "1", false))
}

type NftDepositTry struct {
	NftRevData `bson:",inline" json:",inline"`

	Snap ecsx.GasSnapShot `bson:"snap" json:"-"`

	Status int `bson:"status" json:"status"`

	CreateTime mms.MMS `bson:"create_time" json:"create_time"`
	CreateYMD  int     `bson:"create_ymd" json:"create_ymd"`
}
type NftDepositTryList []NftDepositTry

func (my NftDepositTry) String() string     { return dbg.ToJSONString(my) }
func (my NftDepositTryList) String() string { return dbg.ToJSONString(my) }

func (my NftDepositTry) InsertTryDB(db mongo.DATABASE) (string, error) {
	receiptCode := NFTReceiptCode(my.PaySymbol)

	my.ReceiptCode = receiptCode
	my.CreateTime = mms.Now()
	my.CreateYMD = my.CreateTime.YMD()

	err := db.C(inf.NFTDepositTry).Insert(my)

	return receiptCode, err
}
func (my NftDepositTry) RemoveTryDB(db mongo.DATABASE) error {
	return db.C(inf.NFTDepositTry).Remove(my.Selector())
}

func (my NftDepositTry) UpdateTryDB(db mongo.DATABASE) error {
	return db.C(inf.NFTDepositTry).Update(my.Selector(), my)
}

func (NftDepositTry) GetByReceiptCode(db mongo.DATABASE, receiptCode string) NftDepositTry {
	item := NftDepositTry{}
	DB(func(db mongo.DATABASE) {
		selector := mongo.Bson{"receipt_code": receiptCode}
		db.C(inf.NFTDepositTry).Find(selector).One(&item)
	})
	return item
}
func (NftDepositTry) GetTokenID(db mongo.DATABASE, token_id string) NftDepositTry {
	item := NftDepositTry{}
	DB(func(db mongo.DATABASE) {
		selector := mongo.Bson{"token_id": token_id}
		db.C(inf.NFTDepositTry).Find(selector).One(&item)
	})
	return item
}

func (my NftDepositTry) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.NFTDepositTry)
		NftRevData{}._indexingDB(c, true)

		c.EnsureIndex(mongo.SingleIndex("status", "1", false))
		c.EnsureIndex(mongo.SingleIndex("create_time", "1", false))
		c.EnsureIndex(mongo.SingleIndex("create_ymd", "1", false))
	})
}
