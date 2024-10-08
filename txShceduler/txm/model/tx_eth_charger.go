package model

import (
	"jtools/jmath"
	"jtools/mms"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/txm/inf"
)

// TxETHCharger :
type TxETHCharger struct {
	Key       string  `bson:"key" json:"key"`
	UID       int64   `bson:"uid" json:"uid"`
	Address   string  `bson:"address" json:"address"`
	Price     string  `bson:"price" json:"price"`
	Hash      string  `bson:"hash" json:"hash"`
	State     TxState `bson:"state" json:"state"`
	Timestamp mms.MMS `bson:"timestamp" json:"timestamp"`
	YMD       int     `bson:"ymd" json:"ymd"`
}

// TxETHChargerList :
type TxETHChargerList []TxETHCharger

// Selector :
func (my TxETHCharger) Selector() mongo.Bson { return mongo.Bson{"key": my.Key} }

// InsertDB :
func (my TxETHCharger) InsertDB(db mongo.DATABASE) {
	c := db.C(inf.TXETHCharger)
	c.Insert(my)
}

// TxChargeGroup :
type TxChargeGroup struct {
	Address string           `json:"address"`
	Price   string           `json:"price"`
	List    TxETHChargerList `json:"list"`
}

// RemoveForError :
func (my TxChargeGroup) RemoveForError(db mongo.DATABASE) {
	for _, v := range my.List {
		selector := v.Selector()
		db.C(inf.TXETHCharger).Remove(selector)
		db.C(inf.TXETHDepositLog).Remove(selector)
	}
}

// GroupBy : [address]
func (my TxETHChargerList) GroupBy() map[string]TxChargeGroup {
	group := map[string]TxChargeGroup{}
	for _, tx := range my {
		address := tx.Address
		if v, do := group[address]; do == false {
			group[address] = TxChargeGroup{
				Address: address,
				Price:   tx.Price,
				List:    TxETHChargerList{tx},
			}
		} else {
			v.List = append(v.List, tx)
			v.Price = jmath.ADD(v.Price, tx.Price)
			group[address] = v
		}
	}
	return group
}

// IndexingDB :
func (my TxETHCharger) IndexingDB() {
	inf.DBCollection(inf.TXETHCharger, func(c mongo.Collection) {
		c.EnsureIndex(mongo.SingleIndex("state", "1", false))
		c.EnsureIndex(mongo.SingleIndex("timestamp", "1", false))
		c.EnsureIndex(mongo.MultiIndex([]interface{}{
			"state", "1",
			"timestamp", "1",
		}, false))
	})
}
