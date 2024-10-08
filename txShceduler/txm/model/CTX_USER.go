package model

import (
	"jtools/mms"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/txm/inf"
)

type CTX_USER struct {
	Address   string  `bson:"address" json:"address"`
	Tag       string  `bson:"tag" json:"tag"`
	Timestamp mms.MMS `bson:"timestamp" json:"timestamp"`
	YMD       int     `bson:"ymd" json:"ymd"`
}

func (CTX_USER) IndexingDB() {
	inf.DBCollection(inf.CTX_USER, func(c mongo.Collection) {
		c.EnsureIndex(mongo.SingleIndex("address", 1, true))
	})
}

func ctx_insert_db(db mongo.DATABASE, address, tag string) error {
	my := CTX_USER{}
	my.Address = address
	my.Tag = tag
	my.Timestamp = mms.Now()
	my.YMD = my.Timestamp.YMD()
	return db.C(inf.CTX_USER).Insert(my)
}

func ctx_remove_db(db mongo.DATABASE, address string) {
	db.C(inf.CTX_USER).Remove(mongo.Bson{"address": address})
}

func UserTransactionEnd(db mongo.DATABASE, user_address string) {
	ctx_remove_db(db, user_address)
}

func UserTransactionStart(
	db mongo.DATABASE,
	user_address string,
	tag string,
	process_func ...func() bool,
) bool {
	if ctx_insert_db(db, user_address, tag) == nil {
		if len(process_func) > 0 {
			if !process_func[0]() {
				ctx_remove_db(db, user_address)
				//return false
			}
		}
		return true
	}
	return false
}
