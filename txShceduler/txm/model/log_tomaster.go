package model

import (
	"jtools/mms"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/txm/inf"
)

type LogToMaster struct {
	User `bson:",inline" json:",inline"`

	IsReady       bool   `bson:"is_ready" json:"is_ready"`
	Hash          string `bson:"hash" json:"hash"`
	Symbol        string `bson:"symbol" json:"symbol"`
	Contract      string `bson:"contract" json:"contract"`
	Decimal       string `bson:"decimal" json:"decimal"`
	Price         string `bson:"price" json:"price"`
	MasterAddress string `bson:"master_address" json:"master_address"`

	GasFee string `bson:"gas" json:"gas"`

	Timestamp mms.MMS `bson:"timestamp" json:"timestamp"`
	YMD       int     `bson:"ymd" json:"ymd"`

	SendAt  mms.MMS `bson:"send_at" json:"send_at"`
	SendYMD int     `bson:"send_ymd" json:"send_ymd"`
	IsSend  bool    `bson:"is_send" json:"is_send"`
}

func (my LogToMaster) AckPairs(coin CoinData) []interface{} {
	pairs := []interface{}{
		"name", my.Name,
		"uid", my.UID,
		"from_address", my.Address,
		"master_address", my.MasterAddress,
		"hash", my.Hash,
		"contract", my.Contract,
		"symbol", my.Symbol,
		"price", my.Price,
		"gas_fee", my.GasFee,
		"timestamp", my.Timestamp,
		"remain_coin", coin,
	}
	return pairs
}

/*
"name", my.Name,
"uid", my.UID,
"from_address", my.Address,
"master_address", my.MasterAddress,
"hash", my.Hash,
"contract", my.Contract,
"symbol", my.Symbol,
"price", my.Price,
"gas_fee", my.GasFee,
"timestamp", my.Timestamp,
"remain_coin", coin,
*/
func (my LogToMaster) AckJson(coin CoinData) chttp.JsonType {
	data := chttp.JsonType{}
	pairs := my.AckPairs(coin)
	for i := 0; i < len(pairs); i += 2 {
		key := pairs[i].(string)
		val := pairs[i+1]
		data[key] = val
	} //for
	return data
}

type LogToMasterList []LogToMaster

func (my LogToMaster) String() string     { return dbg.ToJSONString(my) }
func (my LogToMasterList) String() string { return dbg.ToJSONString(my) }

func (my LogToMaster) Valid() bool { return my.Hash != "" }

// GetList : is_ready == true && is_send == false
func (LogToMaster) GetList(db mongo.DATABASE) LogToMasterList {
	list := LogToMasterList{}
	selector := mongo.Bson{
		"is_ready": true,
		"is_send":  false,
	}
	db.C(inf.COLLogToMaster).
		Find(selector).
		Sort("timestamp").
		All(&list)

	return list
}

func (my LogToMaster) Selector() mongo.Bson { return mongo.Bson{"hash": my.Hash} }

func (LogToMaster) InsertDB(
	db mongo.DATABASE,
	masterAddress string,
	member Member,
	tx TxETHDeposit,
	gasFee string,
	nowAt mms.MMS) {

	log := LogToMaster{
		User:          member.User,
		Hash:          tx.Hash,
		Symbol:        tx.Symbol,
		Contract:      tx.Contract,
		Decimal:       tx.Decimal,
		Price:         tx.Price,
		MasterAddress: masterAddress,

		GasFee: gasFee,

		Timestamp: nowAt,
		YMD:       nowAt.YMD(),

		IsSend: false,
	}
	db.C(inf.COLLogToMaster).Insert(log)
}

func (my LogToMaster) SendOK(db mongo.DATABASE, nowAt mms.MMS) {
	upQuery := mongo.Bson{"$set": mongo.Bson{
		"send_at":  nowAt,
		"send_ymd": nowAt.YMD(),
		"is_send":  true,
	}}
	db.C(inf.COLLogToMaster).Update(my.Selector(), upQuery)
}

func (LogToMaster) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.COLLogToMaster)
		c.EnsureIndex(mongo.SingleIndex("hash", "1", true))
		c.EnsureIndex(mongo.SingleIndex("uid", "1", false))
		c.EnsureIndex(mongo.SingleIndex("address", "1", false))

		c.EnsureIndex(mongo.SingleIndex("is_ready", "1", false))
		c.EnsureIndex(mongo.SingleIndex("is_send", "1", false))

		c.EnsureIndex(mongo.MultiIndexName(
			[]interface{}{
				"is_send", 1,
				"timestamp", 1,
			},
			false,
			"m_index_1",
			0,
		))
	})
}
