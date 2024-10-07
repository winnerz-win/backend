package model

import (
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/txm/inf"
)

type xMigrate struct {
	Key string `bson:"key" json:"key"`
}

func (my xMigrate) String() string { return dbg.ToJSONString(my) }

func (xMigrate) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		db.C(inf.XMigrate).EnsureIndex(mongo.SingleIndex("key", "1", true))
	})
}

func MigrationDB(key string, callback func(db mongo.DATABASE)) {
	if key == "" || callback == nil {
		return
	}
	DB(func(db mongo.DATABASE) {
		col := db.C(inf.XMigrate)
		if cnt, _ := col.Find(mongo.Bson{"key": key}).Count(); cnt == 0 {
			dbg.YellowBG("MIGRATION :", key)
			col.Insert(xMigrate{key})
			callback(db)
		}
	})
}
