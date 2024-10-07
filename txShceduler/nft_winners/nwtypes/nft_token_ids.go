package nwtypes

import (
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/jmath"
	"txscheduler/nft_winners/nwdb"
)

type NftTokenIDS struct {
	TokenId string `bson:"token_id" json:"token_id"`
	Owner   string `bson:"owner" json:"owner"`
	IsFixed bool   `bson:"is_fixed" json:"is_fixed"`
}

func (my NftTokenIDS) Valid() bool {
	if !my.IsFixed {
		return false
	}
	return my.TokenId != ""
}
func (my NftTokenIDS) Selector() mongo.Bson { return mongo.Bson{"token_id": my.TokenId} }

func (NftTokenIDS) IndexingDB(db mongo.DATABASE) {
	c := db.C(nwdb.NftTokenIDS)
	c.EnsureIndex(mongo.SingleIndex("token_id", 1, true))
	c.EnsureIndex(mongo.SingleIndex("owner", 1, false))
	c.EnsureIndex(mongo.SingleIndex("is_fixed", 1, false))
}

func (NftTokenIDS) UpdateOwner(db mongo.DATABASE, token_id any, owner string) {
	TOKEN_ID := jmath.VALUE(token_id)

	db.C(nwdb.NftTokenIDS).Upsert(
		mongo.Bson{"token_id": TOKEN_ID},
		mongo.Bson{
			"token_id": TOKEN_ID,
			"owner":    owner,
			"is_fixed": true,
		},
	)

}

func (NftTokenIDS) UpdatePending(db mongo.DATABASE, token_id any, is_pending bool) {
	TOKEN_ID := jmath.VALUE(token_id)
	db.C(nwdb.NftTokenIDS).Upsert(
		mongo.Bson{"token_id": TOKEN_ID},
		mongo.Bson{"$set": mongo.Bson{
			"is_fixed": !is_pending,
		}},
	)
}

func (my NftTokenIDS) IsOwner(db mongo.DATABASE, token_id any, owner string) bool {
	my.TokenId = jmath.VALUE(token_id)

	db.C(nwdb.NftTokenIDS).Find(my.Selector()).One(&my)

	ok := my.Owner == owner && my.IsFixed

	return ok
}

func (my NftTokenIDS) Reserve(
	db mongo.DATABASE,
	token_id string,
	f func() bool,
) bool {
	my.TokenId = token_id

	if err := db.C(nwdb.NftTokenIDS).Insert(my); err == nil {
		if !f() {
			db.C(nwdb.NftTokenIDS).Remove(my.Selector())
		}
		return true
	}
	return false
}
