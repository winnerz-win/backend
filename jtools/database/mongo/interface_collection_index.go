package mongo

import (
	"jtools/cc"
	"jtools/dbg"
	"jtools/jmath"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IndexViewKey Dson
type IndexView struct {
	Key                IndexViewKey `json:"key"`
	Name               string       `json:"name"`
	NS                 string       `json:"ns"`
	V                  int          `json:"v"`
	Unique             bool         `json:"unique"`
	ExpireAfterSeconds *int32       `json:"expireAfterSeconds,omitempty"`
}
type IndexViewList []IndexView

func (my IndexView) String() string     { return toString(my) }
func (my IndexViewList) String() string { return toString(my) }

// NameKeyVal : key_1 --> key | 1
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
		if _Cat(v) == "hashed" {
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
		//cc.GreenItalic(dbg.ToJsonString(void))
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
			//return Eson{my.Key, bson.String(v)}
			return Eson{my.Key, v}

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
		cc.RedItalic("MULTI_INDEX_PAIRS ERROR : ", dbg.Stack())
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
func MultiIndexR(unique bool, pairs ...interface{}) IndexData {
	return MultiIndex(pairs, unique)
}

func MultiIndexNameR(name string, pairs ...interface{}) IndexData {
	data := MultiIndex(pairs, false)
	data.Name = name
	return data
}

func MultiIndexName(pairs []interface{}, unique bool, name string, expaireSeconds int32) IndexData {
	data := MultiIndex(pairs, unique)
	data.Name = name
	data.ExpireAfterSeconds = expaireSeconds
	return data
}

func GetSingleIndexName(key string, order interface{}) string {
	return _Cat(key, "_", order)
}
func GetMultiIndexName(dson []interface{}) string {
	name := ""
	for i := 0; i < len(dson); i += 2 {
		name += _Cat(dson[i], "_", dson[i+1])
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
		//cc.PurpleItalic(str)
	}
	return err
}
func (my *cCollection) EnsureIndexN(data IndexData) error {
	v := my.col.Indexes()
	index_model := data.getModel()
	if _, err := v.CreateOne(
		my.session.ctx(),
		index_model,
	); err == nil {
		return nil
	}

	if err := my.DropIndex(data.Name); err != nil {
		return err
	}
	_, err := v.CreateOne(
		my.session.ctx(),
		index_model,
	)
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

////////////////////////////////////////////////////////

type indexTagItem struct {
	Name   string
	Order  interface{}
	Unique bool
}

func (my indexTagItem) IndexName() string {
	return dbg.Cat(my.Name, "_", my.Order)
}

type indexTagItemList []indexTagItem

func (my indexTagItemList) String() string {
	sl := []string{}
	for _, v := range my {
		if v.Unique {
			sl = append(sl, _Cat(" ", v.Name, ",   ", v.Order, ",   ", v.Unique))
		} else {
			sl = append(sl, _Cat(" ", v.Name, ",   ", v.Order))
		}
	}
	return toString(sl)
}

func (my *indexTagItemList) Set(name, tag string) {

	_sort_order := func(s string) interface{} {
		if jmath.IsNum(s) {
			v := jmath.Int(s)
			switch v {
			case -1, 1:
				return v
			default:
				return 1
			}
		} else if s == "" || s == "false" {
			return 1
		}

		//"text" ....
		return s
	}

	var item indexTagItem
	if strings.Contains(tag, ",") {
		ss := strings.Split(tag, ",")
		item = indexTagItem{
			Name:   name,
			Order:  _sort_order(ss[0]),
			Unique: isTrue(ss[1]),
		}

	} else {
		if isTrue(tag) {
			item = indexTagItem{
				Name:   name,
				Order:  1,
				Unique: true,
			}
		} else {
			item = indexTagItem{
				Name:  name,
				Order: _sort_order(tag),
			}
		}
	}
	(*my) = append((*my), item)
}

func paretStructTag(array *indexTagItemList, v reflect.Type, prefix string, depth int) {
	fn := func(prefix, a string) string {
		if prefix == "" {
			return a
		}
		return prefix + "." + a
	}
	bson_look_up := func(f reflect.StructField) (string, bool) {
		bson_tag, ok := f.Tag.Lookup("bson")
		if !ok {
			return "", false
		}
		bson_tag = strings.ReplaceAll(bson_tag, ",omitempty", "")
		return bson_tag, ok
	}
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)

		// i tag
		if key_tag, ok := f.Tag.Lookup("i"); ok {
			if key_tag == "x" {
				//skip index
			} else {
				if depth > 0 {
					if si_tag, ok := f.Tag.Lookup("si"); ok {
						if si_tag == "x" {
							//skip index
						} else {
							if bson_tag, ok := bson_look_up(f); ok {
								array.Set(fn(prefix, bson_tag), si_tag)

							} else {
								array.Set(fn(prefix, f.Name), si_tag)
							}
						}
					}
				} else {
					if bson_tag, ok := bson_look_up(f); ok {
						array.Set(fn(prefix, bson_tag), key_tag)

					} else {
						array.Set(fn(prefix, f.Name), key_tag)
					}
				}

			}

		} else {
			// inline struct
			if bson_tag, ok := bson_look_up(f); ok {
				if bson_tag == ",inline" {
					if f.Type.Kind() == reflect.Struct {
						paretStructTag(array, f.Type, prefix, depth) //인라인은 0-depth로 한다.
					}
				} else {
					switch f.Type.Kind() {
					case reflect.Struct:
						paretStructTag(array, f.Type, fn(prefix, bson_tag), depth+1)

					case reflect.Ptr:
						elem_ptr := f.Type.Elem()
						if elem_ptr.Kind() == reflect.Struct {
							paretStructTag(array, elem_ptr, fn(prefix, bson_tag), depth+1)
						}

					default:
						if depth > 0 {
							if si_tag, ok := f.Tag.Lookup("si"); ok {
								if si_tag == "x" {
									//skip index
								} else {
									if bson_tag, ok := bson_look_up(f); ok {
										array.Set(fn(prefix, bson_tag), si_tag)

									} else {
										array.Set(fn(prefix, f.Name), si_tag)
									}
								}
							}
						}
					} //switch
				}
			}
		}

	} //for
}

/*
EnsureIndexStruct :
bson:"name" i:"" -> [indexing] key:name , sort:1 , unique: false
i:""  or i:"1"	-> sort 1
i:"-1"			-> sort -1
i:"true" 		-> sort 1 , unique
i:"1,true"		-> sort 1 , unique
i:"x"			-> [indexing] skip

	type Sub struct {
		X int `bson:"x" i:"true" si:""`
		Y int `bson:"y" i:"-1" si:"x"`
	}

	type Name struct {
		Sub `bson:",inline"`				// x,    1,  true
											// y,   -1
		A   int    `i:""`					// A,    1
		B   string `bson:"B" i:"-1,false"`	// B,   -1
		S   Sub    `bson:"s" i:"x"`			// <skip>
		S2  Sub    `bson:"s2"`				// s2.x, 1
	}

----------------------------
[

	" x,   1,   true",
	" y,   -1",
	" A,   1",
	" B,   -1",
	" s2.x,   1"

]
*/
func (my *cCollection) EnsureIndexStruct(i_struct interface{}, prefix ...string) error {
	r_type := reflect.TypeOf(i_struct)
	switch r_type.Kind() {
	case reflect.Struct:

	case reflect.Ptr:
		r_type = r_type.Elem()
		if r_type.Kind() != reflect.Struct {
			return _Error("i_struct is not struct type.")
		}
	default:
		return _Error("i_struct is not struct type.")
	} //switch
	return my._ensure_index_struct(r_type, prefix...)
}
func (my *cCollection) _ensure_index_struct(r_type reflect.Type, prefix ...string) error {
	prefix_tag := ""
	if len(prefix) > 0 {
		prefix_tag = prefix[0]
	}

	list := indexTagItemList{}
	paretStructTag(&list, r_type, prefix_tag, 0)

	for _, v := range list {
		my.EnsureIndex(
			SingleIndex(
				v.Name,
				v.Order,
				v.Unique,
			),
		)
	} //for

	return nil
}

func EnsureIndexStructView(i_struct interface{}) error {
	r_type := reflect.TypeOf(i_struct)
	switch r_type.Kind() {
	case reflect.Struct:
	case reflect.Ptr:
		r_type = r_type.Elem()
		if r_type.Kind() != reflect.Struct {
			return _Error("i_struct is not struct type.")
		}
	default:
		return _Error("i_struct is not struct type.")
	} //switch

	list := indexTagItemList{}
	paretStructTag(&list, r_type, "", 0)

	cc.Yellow(list)

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////

func (my *cCollection) EnsureIndexRenew(i_struct interface{}, datas ...IndexData) error {
	offset := time.Now()
	index_view_list, err := my.Indexes()
	if err != nil {
		return err
	}
	index_name_map := map[string]string{}
	for _, v := range index_view_list {
		index_name_map[v.Name] = v.Name
	} //for

	//create_index_single
	if i_struct != nil {
		r_type := reflect.TypeOf(i_struct)
		switch r_type.Kind() {
		case reflect.Struct:

		case reflect.Ptr:
			r_type = r_type.Elem()
			if r_type.Kind() != reflect.Struct {
				return _Error("i_struct is not struct type.")
			}
		default:
			return _Error("i_struct is not struct type.")
		} //switch

		list := indexTagItemList{}
		paretStructTag(&list, r_type, "", 0)

		for _, v := range list {
			data := SingleIndex(
				v.Name,
				v.Order,
				v.Unique,
			)
			my.EnsureIndex(
				data,
			)

			index_name := v.IndexName()
			delete(index_name_map, index_name)
		} //for
	} //if

	//create_index_multi
	for _, data := range datas {
		my.EnsureIndex(data)
		delete(index_name_map, data.Name)
	} //for

	//remove_index_etc
	if len(index_name_map) > 1 {
		rmv_index_names := []string{}
		for _, index_name := range index_name_map {
			rmv_index_names = append(rmv_index_names, index_name)
		}
		if len(rmv_index_names) > 0 {
			my.DropIndex(rmv_index_names...)
		}
	}
	du := time.Since(offset)
	cc.Green("EnsureIndexRenew :::::::::::::::::::::", du)
	return nil
}
