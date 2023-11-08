package model

import (
	"sync"
	"txscheduler/brix/tools/cloud/ebcm"
	"txscheduler/brix/tools/cloud/ebcm/abi"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx/jwalletx"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/mms"
	"txscheduler/txm/inf"
)

// User :
type User struct {
	UID     int64  `bson:"uid" json:"uid"`
	Address string `bson:"address" json:"address"` //deposit address
	Name    string `bson:"name" json:"name"`
}

func (my User) String() string { return dbg.ToJSONString(my) }

func (my User) PrivateKey() string {
	seed := inf.Config().Seed
	wallet := jwalletx.NewSeed(seed, my.UID)
	return wallet.PrivateKey()
}

// Member :
type Member struct {
	User `bson:",inline" json:",inline"`
	Data map[string]interface{} `bson:"data" json:"data"`

	Coin CoinData `bson:"coin" json:"coin"`

	Deposit  CoinData `bson:"deposit" json:"deposit"`
	Withdraw CoinData `bson:"withdraw" json:"withdraw"`

	CreateAt  mms.MMS `bson:"create_at" json:"create_at"`
	CreateYMD int     `bson:"create_ymd" json:"create_ymd"`

	Timestamp mms.MMS `bson:"timestamp" json:"timestamp"`
	YMD       int     `bson:"ymd" json:"ymd"`
}
type MemberList []Member

func (my Member) String() string     { return dbg.ToJSONString(my) }
func (my MemberList) String() string { return dbg.ToJSONString(my) }

func (my Member) Valid() bool          { return my.UID != 0 }
func (my Member) Selector() mongo.Bson { return mongo.Bson{"uid": my.UID} }

// func (my Member) PrivateKey() string {
// 	seed := inf.Config().Seed
// 	wallet := jwallet.NewSeed(seed, my.UID)
// 	return wallet.PrivateKey()
// }

func (my Member) EtherScanURL() string { return inf.EtherScanAddressURL() + my.Address }

func (my Member) UpdateDB(db mongo.DATABASE) {
	my.Timestamp = mms.Now()
	my.YMD = my.Timestamp.YMD()
	db.C(inf.COLMember).Update(my.Selector(), my)
}

func (my Member) UpdateCoinDB_Legacy(db mongo.DATABASE, sender *ecsx.Sender) {
	if sender == nil {
		dbg.RedItalic("Member.UpdateCoinDB_Legacy.Sender is Nil :", dbg.Stack())
		return
	}
	ethPrice := sender.CoinPrice(my.Address)
	my.Coin.SET(ETH, ethPrice)
	for _, token := range inf.TokenList() {
		if token.Symbol == ETH {
			continue
		}
		tokenPrice := sender.TokenPrice(my.Address, token.Contract, token.Decimal)
		my.Coin.SET(token.Symbol, tokenPrice)
	} //for

	upQuery := mongo.Bson{"$set": mongo.Bson{
		"coin": my.Coin,
	}}
	db.C(inf.COLMember).Update(my.Selector(), upQuery)
}

func (my Member) UpdateCoinDB(db mongo.DATABASE, sender *ebcm.Sender) {
	if sender == nil {
		dbg.RedItalic("Member.UpdateCoinDB.Sender is Nil :", dbg.Stack())
		return
	}
	ethPrice := ebcm.WeiToETH(sender.Balance(my.Address))
	my.Coin.SET(ETH, ethPrice)
	for _, token := range inf.TokenList() {
		if token.Symbol == ETH {
			continue
		}
		sender.Call(
			token.Contract,
			abi.Method{
				Name: "balanceOf",
				Params: abi.NewParams(
					abi.NewAddress(my.Address),
				),
				Returns: abi.NewReturns(
					abi.Uint256,
				),
			},
			my.Address,
			func(rs abi.RESULT) {
				tokenPrice := ebcm.WeiToToken(
					rs.Uint256(0),
					token.Decimal,
				)
				my.Coin.SET(token.Symbol, tokenPrice)
			},
		)
	} //for

	upQuery := mongo.Bson{"$set": mongo.Bson{
		"coin": my.Coin,
	}}
	db.C(inf.COLMember).Update(my.Selector(), upQuery)
}

func LoadMemberAddress(db mongo.DATABASE, address string) Member {
	member := Member{}
	db.C(inf.COLMember).Find(mongo.Bson{"address": address}).One(&member)
	return member
}

func LoadMember(db mongo.DATABASE, uid int64) Member {
	member := Member{}
	db.C(inf.COLMember).Find(mongo.Bson{"uid": uid}).One(&member)
	return member
}

func LoadMemberName(db mongo.DATABASE, name string) Member {
	member := Member{}
	db.C(inf.COLMember).Find(mongo.Bson{"name": name}).One(&member)
	return member
}

// IsMemberAddress :
func IsMemberAddress(db mongo.DATABASE, address string) bool {
	cnt, _ := db.C(inf.COLMember).Find(mongo.Bson{"address": address}).Count()
	return cnt > 0
}

var (
	currentUID int64 = 1001
	muUID      sync.RWMutex
)

func CurrentUID() int64 {
	return currentUID
}
func IncUID() {
	currentUID++
}

// CreateLock :
func CreateLock(f func(db mongo.DATABASE)) {
	defer muUID.Unlock()
	muUID.Lock()
	DB(func(db mongo.DATABASE) {
		f(db)
	})
}

func (Member) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.COLMember)
		c.EnsureIndex(mongo.SingleIndex("uid", "1", true))
		c.EnsureIndex(mongo.SingleIndex("address", "1", true))
		c.EnsureIndex(mongo.SingleIndex("name", "1", true))

		last := Member{}
		if c.Find(nil).Sort("-uid").One(&last) == nil {
			currentUID = last.UID + 1
		}

		dbg.Yellow("==========================================")
		dbg.Yellow("  CurrentUID :", currentUID)
		dbg.Yellow("==========================================")
	})

	MigrationDB("member_cloud", func(db mongo.DATABASE) {
		iter := db.C(inf.COLMember).Find(nil).Iter()
		member := Member{}
		for iter.Next(&member) {
			sender := ecsx.New(inf.Mainnet(), inf.InfuraKey())
			member.UpdateCoinDB_Legacy(db, sender)
		} //for
	})
}

func (my *Member) Mig() {
	if my.Withdraw == nil {
		my.Withdraw = NewCoinData()
	}
	if my.Deposit == nil {
		my.Deposit = NewCoinData()
	}
}
