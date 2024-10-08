package model

import (
	"jtools/mms"
	"time"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/txm/inf"
)

type SyncCoin struct {
	Address    string  `bson:"address" json:"address"`
	IsToMaster bool    `bson:"is_tomaster" json:"is_tomaster"`
	At         mms.MMS `bson:"at" json:"at"`
}

type SyncCoinList []SyncCoin

func (SyncCoin) GetList(db mongo.DATABASE) SyncCoinList {
	list := SyncCoinList{}
	db.C(inf.COLSyncCoin).Find(nil).All(&list)
	return list
}

func (my SyncCoin) Selector() mongo.Bson { return mongo.Bson{"address": my.Address} }

// GetSyncMMS : 10 sec
func GetSyncMMS() mms.MMS {
	nowAt := mms.Now().Add(time.Second * 10)
	return nowAt
}

func (SyncCoin) InsertDB(db mongo.DATABASE, address string, at mms.MMS, isTomaster bool) bool {

	c := db.C(inf.COLSyncCoin)
	item := SyncCoin{
		Address:    address,
		At:         at,
		IsToMaster: isTomaster,
	}
	if c.Insert(item) != nil {
		if at == 0 || isTomaster == true {
			item.Address = address
			c.Upsert(
				item.Selector(),
				mongo.Bson{"$set": mongo.Bson{
					"at":          0,
					"is_tomaster": isTomaster,
				}},
			)
			return true
		}
		return false
	}
	return true
}

// RemoveDB : LogToMaster->IsReady = true
func (my SyncCoin) RemoveDB(db mongo.DATABASE) {
	selector := mongo.Bson{
		"address": my.Address,
	}
	upQuery := mongo.Bson{
		"$set": mongo.Bson{
			"is_ready": true,
		},
	}
	db.C(inf.COLLogToMaster).UpdateAll(selector, upQuery)

	db.C(inf.COLSyncCoin).Remove(my.Selector())
}

func (SyncCoin) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.COLSyncCoin)
		c.EnsureIndex(mongo.SingleIndex("address", "1", true))
		c.EnsureIndex(mongo.SingleIndex("at", "1", false))
	})

}
