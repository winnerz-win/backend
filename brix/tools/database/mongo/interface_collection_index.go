package mongo

import (
	"txscheduler/brix/tools/database/mongo/tools/dbg"
	"txscheduler/brix/tools/database/mongo/tools/jmath"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type IndexViewKey map[string]interface{}
type IndexView struct {
	Key                IndexViewKey `json:"key"`
	Name               string       `json:"name"`
	NS                 string       `json:"ns"`
	V                  int          `json:"v"`
	Unique             bool         `json:"unique"`
	ExpireAfterSeconds *int32       `json:"expireAfterSeconds,omitempty"`
}
type IndexViewList []IndexView

func (my IndexView) String() string     { return dbg.ToJSONString(my) }
func (my IndexViewList) String() string { return dbg.ToJSONString(my) }

//NameKeyVal : key_1 --> key | 1
func (my IndexView) NameKeyVal() (string, string) {
	key := ""
	val := ""
	for i := len(my.Name) - 1; i >= 0; i-- {
		if my.Name[i] == '_' {
			key = my.Name[:i]
			val = my.Name[i+1:]
			break
		}
	} //for
	return key, val
}
func (my IndexView) IsHashed() bool {
	for _, v := range my.Key {
		if dbg.Cat(v) == "hashed" {
			return true
		}
	}
	return false
}

func (my *cCollection) Indexes() (IndexViewList, error) {
	view := my.col.Indexes()
	cur, err := view.List(
		my.session.ctx(),
	)
	if err != nil {
		return IndexViewList{}, err
	}
	defer cur.Close(my.session.ctx())

	list := IndexViewList{}
	for cur.Next(my.session.ctx()) {
		item := IndexView{}
		cur.Decode(&item)
		list = append(list, item)

		var void map[string]interface{}
		cur.Decode(&void)
		//dbg.Green(dbg.ToJSONString(void))
	} //for

	return list, nil
}

type IndexDataPair struct {
	Key   string
	Order interface{}
}

func (my IndexDataPair) getEson() Eson {
	if jmath.IsNum(my.Order) {
		return Eson{my.Key, jmath.Int(my.Order)}
	}
	if v, do := my.Order.(string); do {
		switch v {
		case "text", "2d", "2dsphere":
			return Eson{my.Key, bsonx.String(v)}

		} //switch
	}
	return Eson{my.Key, 1}
}

type IndexData struct {
	pairs              []IndexDataPair
	Unique             bool
	Name               string
	IsExpired          bool
	ExpireAfterSeconds int32
}

func (my *IndexData) SetExpired(flag bool, sec ...int32) {
	my.IsExpired = flag
	if flag && len(sec) > 0 {
		my.ExpireAfterSeconds = sec[0]
	}
}

type IndexDataList []IndexData

func (my IndexData) getModel() mongo.IndexModel {

	dsonKey := Dson{}
	for _, pair := range my.pairs {
		dsonKey = append(dsonKey, pair.getEson().E())
	}

	option := &options.IndexOptions{
		Unique: &my.Unique,
	}
	if my.Name != "" {
		option.SetName(my.Name)
	}
	if my.ExpireAfterSeconds > 0 {
		option.SetExpireAfterSeconds(my.ExpireAfterSeconds)
	}

	model := mongo.IndexModel{
		Keys:    dsonKey,
		Options: option,
	}

	return model
}

func SingleIndex(key string, order interface{}, unique bool) IndexData {
	data := IndexData{
		Unique: unique,
	}
	data.pairs = append(data.pairs,
		IndexDataPair{key, order},
	)
	return data
}
func SingleIndexName(key string, order interface{}, unique bool,
	name string, expaireSeconds int32) IndexData {

	data := SingleIndex(key, order, unique)
	data.Name = name
	data.ExpireAfterSeconds = expaireSeconds

	return data
}

func MultiIndex(pairs []interface{}, unique bool) IndexData {
	if len(pairs)%2 != 0 {
		return IndexData{}
	}
	data := IndexData{
		Unique: unique,
	}
	for i := 0; i < len(pairs); i += 2 {
		key, do := pairs[i].(string)
		if !do {
			return IndexData{}
		}
		order := pairs[i+1]
		data.pairs = append(data.pairs,
			IndexDataPair{key, order},
		)
	}
	return data
}
func MultiIndexName(pairs []interface{}, unique bool, name string, expaireSeconds int32) IndexData {
	data := MultiIndex(pairs, unique)
	data.Name = name
	data.ExpireAfterSeconds = expaireSeconds
	return data
}

func GetSingleIndexName(key string, order interface{}) string {
	return dbg.Cat(key, "_", order)
}
func GetMultiIndexName(dson []interface{}) string {
	name := ""
	for i := 0; i < len(dson); i += 2 {
		name += dbg.Cat(dson[i], "_", dson[i+1])
		if i < len(dson)-2 {
			name += "_"
		}
	}
	return name
}

func (my *cCollection) EnsureIndex(data IndexData) error {
	v := my.col.Indexes()
	str, err := v.CreateOne(
		my.session.ctx(),
		data.getModel(),
	)
	if err == nil {
		_ = str
		//dbg.Purple(str)
	}
	return err
}

func (my *cCollection) DropIndex(names ...string) error {
	view := my.col.Indexes()
	for _, name := range names {
		if name == "_id_" {
			continue
		}
		_, err := view.DropOne(
			my.session.ctx(),
			name,
		)
		if err != nil {
			return err
		}
	} //for

	return nil
}

func (my *cCollection) DropIndexAll() error {
	view := my.col.Indexes()
	_, err := view.DropAll(my.session.ctx())
	return err
}
