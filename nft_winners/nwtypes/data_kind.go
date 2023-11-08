package nwtypes

import (
	"txscheduler/brix/tools/database/mongo"
)

type DATA_KIND string

func (my DATA_KIND) String() string { return string(my) }

const (
	DATA_NONE      = DATA_KIND("")
	NFT_MINT       = DATA_KIND("nft_mint")
	MULTI_TRANSFER = DATA_KIND("multi_transfer")
	NFT_TRANSFER   = DATA_KIND("nft_transfer")

	NFT_BASE_URI = DATA_KIND("set_base_uri")
)

func MAKE_DATA_SEQ(pair ...interface{}) DATA_FLOW {
	my := DATA_FLOW{
		Index: 0,
	}
	for i := 0; i < len(pair); i += 2 {
		kind := pair[i].(DATA_KIND)
		item := mongo.MakeMap(pair[i+1])
		my.List = append(my.List,
			DATA{
				Kind: kind,
				Item: item,
			},
		)
	} //for
	return my
}

type DATA struct {
	Kind DATA_KIND `bson:"kind" json:"kind"`
	Item mongo.MAP `bson:"item" json:"item"`
}

type DATA_FLOW struct {
	List  []DATA `bson:"list" json:"list"`
	Index int    `bson:"index" json:"index"`
}

func (my DATA_FLOW) Current() DATA {
	return my.List[my.Index]
}
func (my *DATA_FLOW) UpdateCurrent(data any) {
	my.List[my.Index].Item = mongo.MakeMap(data)
}

func (my DATA_FLOW) Prev() DATA {
	return my.List[my.Index-1]
}
func (my DATA_FLOW) Next() DATA {
	if my.Index+1 >= len(my.List) {
		return DATA{Kind: DATA_NONE}
	}
	return my.List[my.Index+1]
}

func (my *DATA_FLOW) SetNext() DATA {
	if my.Index+1 >= len(my.List) {
		return DATA{Kind: DATA_NONE}
	}
	my.Index++
	return my.List[my.Index]
}
