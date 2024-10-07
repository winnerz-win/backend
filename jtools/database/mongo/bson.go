package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
)

/*
Bson : map[string]interface{}
*/
type Bson bson.M

type Bsons []Bson

////////////////////////////////////////////////////////////////////

/*
Dson : []Eson
*/
type Dson bson.D

// Set : vs is Eson or Bson or Pair("key" : value)
func (my *Dson) Set(vs ...interface{}) {

	_eson := func(items ...bson.E) {
		for _, item := range items {
			*my = append(*my, item)
		} //for
	}
	_bson := func(items ...bson.M) {
		for _, item := range items {
			for key, val := range item {
				e := Eson{key, val}
				*my = append(*my, e.E())
			} //for
		}
	}

	skip := false
	for i, v := range vs {
		if skip {
			skip = false
			continue
		}

		switch item := v.(type) {
		case string:
			*my = append(*my,
				bson.E{
					Key:   item,
					Value: vs[i+1],
				},
			)
			skip = true

		case Eson:
			_eson(bson.E(item))
		case bson.E:
			_eson(item)

		case Dson:
			_eson(item...)
		case bson.D:
			_eson(item...)

		case Bson:
			_bson(bson.M(item))
		case bson.M:
			_bson(item)

		case Bsons:
			for _, v := range item {
				_bson(bson.M(v))
			}
		case []bson.M:
			_bson(item...)

		} //switch
	} //for

}

func (my *Dson) Add(key string, val interface{}) {
	*my = append(*my,
		bson.E{
			Key:   key,
			Value: val,
		},
	)
}

// DSON : vs is Eson or Bson or Pair("key" : value)
func DSON(vs ...interface{}) Dson {
	dson := Dson{}
	dson.Set(vs...)
	return dson
}

func (my Dson) Map() Bson {
	return Bson(bson.D(my).Map())
}

////////////////////////////////////////////////////////////////////

/*
	Eson {
		Key stirng
		Value interface{}
	}
*/
type Eson bson.E

func (my Eson) E() bson.E { return bson.E(my) }
