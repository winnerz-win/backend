package model

import (
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/mms"
	"txscheduler/txm/inf"
)

type FailNftCode int

const (
	FailNftNone      = FailNftCode(0)
	FailNftDeposit   = FailNftCode(1001) //Pay계좌로 Deposit 실패
	FailNftTxResult  = FailNftCode(1002) //트랜젝션 결과 실패
	FailNftSendError = FailNftCode(1003) //구매 트랜젝션 실패 (토큰 또는 이더 잔액 부족, 구매 컨트랙트 approve 실패등등)
)

func NFTFailCodeMessage(code FailNftCode, exMsg ...string) string {
	msg := "transaction_error"
	switch code {
	case FailNftDeposit:
		msg = "pay_deposit_fail"
	case FailNftTxResult:
		msg = "transaction_result_fail"
	case FailNftSendError:
		msg = "transaction_try_fail"
	}
	if len(exMsg) > 0 {
		msg += "(" + exMsg[0] + ")"
	}
	return msg
}

type NftMetaData struct {
	TokenType string `bson:"token_type" json:"token_type"`
	Name      string `bson:"name" json:"name"`
	Content   string `bson:"content" json:"content"`
	Deployer  string `bson:"deployer" json:"deployer"`
	Desc      string `bson:"desc" json:"desc"`
}

func (my NftMetaData) String() string { return dbg.ToJSONString(my) }

type NftBuyEnd struct {
	NftBuyTry `bson:",inline" json:",inline"`

	LastOwner string `bson:"last_owner" json:"last_owner"`

	TokenURI string      `bson:"token_uri" json:"token_uri"`
	Meta     NftMetaData `bson:"meta,omitempty" json:"meta,omitempty"`

	Timestamp mms.MMS `bson:"timestamp" json:"timestamp"`
	YMD       int     `bson:"ymd" json:"ymd"`

	GasFeeETH string `bson:"gas_fee_eth" json:"gas_fee_eth"`

	IsBurn   bool   `bson:"is_burn" json:"is_burn"`
	BurnHash string `bson:"burn_hash,omitempty" json:"burn_hash,omitempty"`

	IsSend  bool    `bson:"is_send" json:"is_send"`
	SendAt  mms.MMS `bson:"send_at" json:"send_at"`
	SendYMD int     `bson:"send_ymd" json:"send_ymd"`
}

type NftBuyEndList []NftBuyEnd

func (my NftBuyEnd) CallbackData() chttp.JsonType {
	data := chttp.JsonType{
		"uid":          my.UID,
		"address":      my.Address,
		"name":         my.Name,
		"receipt_code": my.ReceiptCode,

		"pay_address": my.PayAddress,
		"pay_symbol":  my.PaySymbol,
		"pay_price":   my.PayPrice,
		"is_pay_free": my.IsPayFree,

		"last_owner": my.LastOwner,

		"token_id":   my.TokenId,
		"token_type": my.TokenType,

		"token_uri": my.TokenURI,
		"meta":      my.Meta,

		"deposit_hash":      my.DepositHash,
		"deposit_gasLimit":  my.DepositGasLimit,
		"deposit_gasPrice":  my.DepositGasPrice,
		"deposit_gasFeeETH": my.DepositGasFeeETH,

		"hash": my.Hash,

		"status":       my.Status,
		"fail_code":    my.FailCode,
		"fail_message": my.FailMessage,

		"timestamp": my.Timestamp,

		"gasLimit":  my.GasLimit,
		"gasPrice":  my.GasPrice,
		"gasFeeETH": my.NftBuyTry.GasFeeETH,

		"gas_fee_eth": my.GasFeeETH,

		"is_burn":   my.IsBurn,
		"burn_hash": my.BurnHash,
	}
	return data
}

func (NftBuyEnd) GetByReceiptCode(db mongo.DATABASE, receiptCode string) NftBuyEnd {
	item := NftBuyEnd{}
	DB(func(db mongo.DATABASE) {
		selector := mongo.Bson{"receipt_code": receiptCode}
		if err := db.C(inf.NFTBuyEnd).Find(selector).One(&item); err != nil {
			dbg.Red(err)
		}
	})
	return item
}
func (NftBuyEnd) GetTokenID(db mongo.DATABASE, token_id string) NftBuyEnd {
	item := NftBuyEnd{}
	DB(func(db mongo.DATABASE) {
		selector := mongo.Bson{"token_id": token_id, "is_burn": false}
		db.C(inf.NFTBuyEnd).Find(selector).One(&item)
	})
	return item
}

func (my NftBuyEnd) InsertEndDB(db mongo.DATABASE) error {
	my.Timestamp = mms.Now()
	my.YMD = my.Timestamp.YMD()
	if my.GasFeeETH == "" {
		my.GasFeeETH = ZERO
	}
	return db.C(inf.NFTBuyEnd).Insert(my)
}

func (my NftBuyEnd) UpdateTokenURICheck(db mongo.DATABASE) {
	if my.TokenURI == "" {

	}
}

func (my NftBuyEnd) BurnDB(db mongo.DATABASE, burnHash string) {
	upQuery := mongo.Bson{"$set": mongo.Bson{
		"is_burn":   true,
		"burn_hash": burnHash,
	}}
	db.C(inf.NFTBuyEnd).Update(my.Selector(), upQuery)
}
func (my NftBuyEnd) SetLastOwner(db mongo.DATABASE, owner string) {
	upQuery := mongo.Bson{"$set": mongo.Bson{
		"last_owner": owner,
	}}
	db.C(inf.NFTBuyEnd).Update(my.Selector(), upQuery)
}

func (my NftBuyEnd) ChangeInnerOwner(db mongo.DATABASE, newOwner User) {
	upQuery := mongo.Bson{"$set": mongo.Bson{
		"uid":        newOwner.UID,
		"address":    newOwner.Address,
		"name":       newOwner.Name,
		"last_owner": newOwner.Address,
	}}
	db.C(inf.NFTBuyEnd).Update(my.Selector(), upQuery)
}

func (my NftBuyEnd) SendOK(db mongo.DATABASE, nowAt mms.MMS) error {
	my.IsSend = true
	my.SendAt = nowAt
	my.SendYMD = nowAt.YMD()
	return db.C(inf.NFTBuyEnd).Update(my.Selector(), my)
}

func (my NftBuyEnd) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.NFTBuyEnd)
		NftRevData{}._indexingDB(c, false)

		c.EnsureIndex(mongo.SingleIndex("create_time", "1", false))
		c.EnsureIndex(mongo.SingleIndex("create_ymd", "1", false))

		c.EnsureIndex(mongo.SingleIndex("hash", "1", false))
		c.EnsureIndex(mongo.SingleIndex("timestamp", "1", false))
		c.EnsureIndex(mongo.SingleIndex("ymd", "1", false))

		c.EnsureIndex(mongo.SingleIndex("is_send", "1", false))
		c.EnsureIndex(mongo.SingleIndex("send_at", "1", false))
		c.EnsureIndex(mongo.SingleIndex("send_ymd", "1", false))

		c.EnsureIndex(mongo.SingleIndex("is_burn", "1", false))

	})
}
