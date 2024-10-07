package model

import (
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/txm/inf"
)

type NftAset struct {
	Key         string `bson:"key" json:"key"`
	Number      string `bson:"number" json:"number"`
	IsEnd       bool   `bson:"is_end" json:"is_end"`
	NFTContract string `bson:"nft_contract" json:"nft_contract"`
	BaseURI     string `bson:"baseURI" json:"baseURI"`
	NftName     string `bson:"nft_name" json:"nft_name"`
	NftSymbol   string `bson:"nft_symbol" json:"nft_symbol"`
}

func (my *NftAset) Selector() mongo.Bson {
	my.Key = inf.NFTASET
	return mongo.Bson{"key": my.Key}
}

func (my NftAset) UpdateIncDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.NFTASET)

		upquery := mongo.Bson{"$set": mongo.Bson{"number": my.Number}}
		c.Update(my.Selector(), upquery)
	})
}

func (my NftAset) UpdateEndFlagDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.NFTASET)

		upquery := mongo.Bson{"$set": mongo.Bson{"is_end": my.IsEnd}}
		c.Update(my.Selector(), upquery)
	})
}

func (my NftAset) UpdateBaseURI(db mongo.DATABASE, uri string) {
	c := db.C(inf.NFTASET)
	upquery := mongo.Bson{"$set": mongo.Bson{"baseURI": uri}}
	c.Update(my.Selector(), upquery)
}

func (my NftAset) UpdateQueryDB(db mongo.DATABASE, subQuery mongo.Bson) {
	c := db.C(inf.NFTASET)
	upQuery := mongo.Bson{"$set": subQuery}
	c.Update(my.Selector(), upQuery)
}

func (NftAset) FirstLoadDB(startNumber string) NftAset {
	item := NftAset{}
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.NFTASET)

		if c.Find(item.Selector()).One(&item) != nil {
			item.Number = startNumber
			item.IsEnd = false
			item.Selector()
			c.Insert(item)
		}
	})
	return item
}

func (my NftAset) IndexingDB(firstSetFunc func(db mongo.DATABASE)) {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.NFTASET)
		c.EnsureIndex(mongo.SingleIndex("key", "1", true))

		c.EnsureIndex(mongo.SingleIndex("is_end", "1", false))

		firstSetFunc(db)
	})
}
