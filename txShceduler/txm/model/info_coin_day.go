package model

import (
	"jtools/mms"
	"sync"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/txm/inf"
)

// CoinDay :
type CoinDay struct {
	YMD        int      `bson:"ymd" json:"ymd"`
	Deposit    CoinData `bson:"deposit" json:"deposit"`
	Withdraw   CoinData `bson:"withdraw" json:"withdraw"`
	MasterGas  PRICE    `bson:"master_gas" json:"master_gas"`
	ChargerGas PRICE    `bson:"charger_gas" json:"charger_gas"`
	MemberGas  PRICE    `bson:"member_gas" json:"member_gas"`
}

// NewCoinDay :
func NewCoinDay(nowAt mms.MMS) *CoinDay {
	return &CoinDay{
		YMD:        nowAt.YMD(),
		Deposit:    NewCoinData(),
		Withdraw:   NewCoinData(),
		MasterGas:  ZERO,
		ChargerGas: ZERO,
		MemberGas:  ZERO,
	}
}

func (my CoinDay) Selector() mongo.Bson { return mongo.Bson{"ymd": my.YMD} }

type CoinDayList []CoinDay

func (my CoinDay) String() string     { return dbg.ToJSONString(my) }
func (my CoinDayList) String() string { return dbg.ToJSONString(my) }

var (
	_coin_day_mu sync.Mutex
)

func (CoinDay) AddDeposit(db mongo.DATABASE, symbol, price string, nowAt mms.MMS) {
	_coin_day_mu.Lock()
	defer _coin_day_mu.Unlock()

	item := getCoinDayDB(db, nowAt)
	item.Deposit.ADD(symbol, price)
	item.UpsertDB(db)
}
func (CoinDay) AddWithdraw(db mongo.DATABASE, symbol, price, gas string, nowAt mms.MMS) {
	_coin_day_mu.Lock()
	defer _coin_day_mu.Unlock()

	item := getCoinDayDB(db, nowAt)
	item.Withdraw.ADD(symbol, price)
	item.MasterGas.ADD(gas)
	item.UpsertDB(db)
}
func (CoinDay) AddChargerGas(db mongo.DATABASE, gas string, nowAt mms.MMS) {
	_coin_day_mu.Lock()
	defer _coin_day_mu.Unlock()

	item := getCoinDayDB(db, nowAt)
	item.ChargerGas.ADD(gas)
	item.UpsertDB(db)
}
func (CoinDay) AddMasterGas(db mongo.DATABASE, gas string, nowAt mms.MMS) {
	_coin_day_mu.Lock()
	defer _coin_day_mu.Unlock()

	item := getCoinDayDB(db, nowAt)
	item.MasterGas.ADD(gas)
	item.UpsertDB(db)
}
func (CoinDay) AddMemberGas(db mongo.DATABASE, gas string, nowAt mms.MMS) {
	_coin_day_mu.Lock()
	defer _coin_day_mu.Unlock()

	item := getCoinDayDB(db, nowAt)
	item.MemberGas.ADD(gas)
	item.UpsertDB(db)
}

func (CoinDay) Action(db mongo.DATABASE, nowAt mms.MMS, callback func(day *CoinDay)) {
	item := getCoinDayDB(db, nowAt)
	callback(item)
	item.UpsertDB(db)
}

func getCoinDayDB(db mongo.DATABASE, nowAt mms.MMS) *CoinDay {
	item := NewCoinDay(nowAt)
	if db.C(inf.COLCoinDay).Find(item.Selector()).One(item) != nil {
		item = NewCoinDay(nowAt)
	}
	return item
}

func (my CoinDay) UpsertDB(db mongo.DATABASE) {
	db.C(inf.COLCoinDay).Upsert(my.Selector(), my)
}

func (CoinDay) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.COLCoinDay)
		c.EnsureIndex(mongo.SingleIndex("ymd", "1", true))

		//migrate
		if cnt, _ := c.Count(); cnt == 0 {
			{
				iter := db.C(inf.COLMember).Find(nil).Iter()
				member := Member{}
				for iter.Next(&member) {
					member.Mig()
					member.UpdateDB(db)
					member = Member{}
				} //for
			}

			if cnt, _ := db.C(inf.COLLogDeposit).Count(); cnt > 0 {
				iter := db.C(inf.COLLogDeposit).Find(nil).Iter()
				log := LogDeposit{}
				for iter.Next(&log) {
					member := LoadMember(db, log.UID)
					member.Deposit.ADD(log.Symbol, log.Price)
					member.UpdateDB(db)

					CoinDay{}.AddDeposit(db, log.Symbol, log.Price, log.Timestamp)
				} //for
			}

			if cnt, _ := db.C(inf.COLLogWithdraw).Count(); cnt > 0 {
				iter := db.C(inf.COLLogWithdraw).Find(nil).Iter()
				log := LogWithdraw{}
				for iter.Next(&log) {
					member := LoadMember(db, log.UID)
					member.Withdraw.ADD(log.Symbol, log.ToPrice)
					member.UpdateDB(db)

					CoinDay{}.AddWithdraw(db, log.Symbol, log.ToPrice, log.Gas, log.Timestamp)
				}
			}
		}
	})
}
