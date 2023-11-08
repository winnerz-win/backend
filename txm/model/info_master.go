package model

import (
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/txm/inf"
)

type InfoMaster struct {
	Mainnet        bool     `bson:"mainnet" json:"mainnet"`
	MasterAddress  string   `bson:"master_address" json:"master_address"`
	ChargerAddress string   `bson:"charger_address" json:"charget_address"`
	Symbols        []string `bson:"symbols" json:"symbols"`
}

func (my InfoMaster) String() string { return dbg.ToJSONString(my) }

func (my InfoMaster) Selector() mongo.Bson { return mongo.Bson{"mainnet": my.Mainnet} }

func (InfoMaster) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.COLInfoMaster)
		c.EnsureIndex(mongo.SingleIndex("mainnet", "1", true))

		item := InfoMaster{
			Mainnet:        inf.Mainnet(),
			MasterAddress:  inf.Master().Address,
			ChargerAddress: inf.Charger().Address,
			Symbols:        inf.SymbolList(),
		}

		c.Upsert(item.Selector(), item)
	})
}
