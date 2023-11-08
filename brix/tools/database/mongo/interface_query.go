package mongo

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

///////////////////////////////////////////////////////////////////

type Query interface {
	Skip(n int) Query
	Limit(n int) Query
	Count() (int, error)
	Sort(sorts ...string) Query

	One(r interface{}) error
	All(r interface{}) error

	Iter(ctxs ...context.Context) Iterator
}

type cQuery struct {
	c          *cCollection
	selector   interface{}
	projection interface{} //bson.D
	skip       int64
	limit      int64
	sorts      bson.D
}

func (my *cCollection) Find(selector interface{}) Query {
	query := &cQuery{
		c:        my,
		selector: getFilter(selector),
	}
	query.projection = nil
	return query
}

func get_projection_field(dson Dson, rt reflect.Type) Dson {
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if !field.IsExported() {
			continue
		}

		if tag, do := field.Tag.Lookup("bson"); do {
			if tag == "-" {
				continue
			}
			if tag == ",inline" {
				switch field.Type.Kind() {
				case reflect.Struct:
					dson = get_projection_field(dson, field.Type)

				case reflect.Ptr:
					if field.Type.Elem().Kind() == reflect.Struct {
						dson = get_projection_field(dson, field.Type.Elem())
					}
				}
			} else {
				item := primitive.E{
					Key:   tag,
					Value: 1,
				}
				dson = append(dson, item)
			}
		} else {
			item := primitive.E{
				Key:   strings.ToLower(field.Name),
				Value: 1,
			}
			dson = append(dson, item)
		}
	}
	return dson
}

func (my *cCollection) FindProjection(selector interface{}, projection interface{}) Query {
	query := &cQuery{
		c:        my,
		selector: getFilter(selector),
	}

	var _projection_data interface{}
	ptype := reflect.TypeOf(projection)
	switch ptype.Kind() {
	case reflect.Ptr:
		ptype = ptype.Elem()
		if ptype.Kind() == reflect.Struct {
			query.projection = get_projection_field(Dson{}, ptype)
		} else {
			_projection_data = projection
		}

	case reflect.Struct:
		query.projection = get_projection_field(Dson{}, ptype)

	default:
		_projection_data = projection
	} //switch

	query.projection = _projection_data
	return query
}

func (my *cQuery) Skip(n int) Query {
	my.skip = int64(n)
	return my
}
func (my *cQuery) Limit(n int) Query {
	if n <= 0 {
		n = 0
	}
	my.limit = int64(n)
	return my
}
func (my *cQuery) Sort(sorts ...string) Query {
	getE := func(s string) bson.E {
		if strings.HasPrefix(s, "-") {
			return bson.E{s[1:], -1}
		} else {
			return bson.E{s, 1}
		}
	}
	for _, v := range sorts {
		my.sorts = append(my.sorts, getE(v))
	}
	return my
}
func (my *cQuery) getSort() bson.D {
	if len(my.sorts) > 0 {
		return my.sorts
	}
	return bson.D{bson.E{"_id", 1}}
}

func (my *cQuery) getOption() *options.FindOptions {
	opt := options.Find()
	if my.projection != nil {
		opt = opt.SetProjection(my.projection)
	}
	opt.SetSkip(my.skip)
	opt.SetLimit(my.limit)
	// if len(my.sorts) > 0 {
	// 	opt.SetSort(my.sorts)
	// }
	opt.Sort = my.getSort()

	return opt
}

func (my *cQuery) getCursor() (*mongo.Cursor, error) {
	return my.c.col.Find(
		my.c.session.ctx(),
		my.selector,
		my.getOption(),
	)
}

func (my *cQuery) Count() (int, error) {
	// cursor, err := my.getCursor()
	// if err != nil {
	// 	return 0, err
	// }
	// return cursor.RemainingBatchLength(), err

	// opt := options.Count()
	// if my.skip > 0 {
	// 	opt.SetSkip(my.skip)
	// }
	// if my.limit > 0 {
	// 	opt.SetLimit(my.limit)
	// }

	// cnt, err := my.c.col.CountDocuments(
	// 	my.c.session.ctx(),
	// 	my.selector,
	// 	opt,
	// )

	// return int(cnt), err

	return my.c.cdb.CollectionCount(
		my.c.dbName,
		my.c.col.Name(),
		my.selector,
		my.skip,
		my.limit,
	)

}

func (my *cQuery) One(r interface{}) error {
	opt := &options.FindOneOptions{
		Skip: &my.skip,
	}
	// if len(my.sorts) > 0 {
	// 	opt.Sort = my.sorts
	// }
	opt.Sort = my.getSort()

	if my.projection != nil {
		opt = opt.SetProjection(my.projection)
	}

	sr := my.c.col.FindOne(
		my.c.session.ctx(),
		my.selector,
		opt,
	)
	if sr == nil {
		return errors.New("One is nil")
	}
	return sr.Decode(r)
}

func (my *cQuery) All(slicePointer interface{}) error {
	cursor, err := my.getCursor()
	if err != nil {
		return err
	}

	return cursor.All(
		my.c.session.ctx(),
		slicePointer,
	)
}

///////////////////////////////////////////////////////////////////

func (my *cQuery) Iter(ctxs ...context.Context) Iterator {
	cursor, err := my.getCursor()
	if err != nil {
		return &cIterator{
			err: err,
		}
	}
	iter := &cIterator{
		ctx:    context.Background(),
		cursor: cursor,
	}
	if len(ctxs) > 0 {
		iter.ctx = ctxs[0]
	}
	return iter
}

func (my *cQuery) Iterator(f func(iter Iterator), ctxs ...context.Context) (g_err error) {
	cursor, err := my.getCursor()
	if err != nil {
		return err
	}
	iter := &cIterator{
		ctx:    my.c.session.ctx(),
		cursor: cursor,
	}
	if len(ctxs) > 0 {
		iter.ctx = ctxs[0]
	}

	defer func() {
		if e := recover(); e != nil {
			g_err = _Error("Query.Iterator Panic[", e, "]", _Stack())
		}
		iter.Close()
	}()

	f(iter)

	return g_err
}

// IterForeach : [return true]일 경우 loop 종료!, cnt는 1부터 ~
func IterForeach[T any](iter Iterator, f func(cnt int, item T) bool) {
	if iter.Error() != nil {
		return
	}

	defer iter.Close()

	counter := 1

	for {
		var v T
		if iter.Next(&v) {
			if f(counter, v) {
				break
			}
			counter++
		} else {
			break
		}
	} //for
}

// IterAll :
func IterAll[T any](iter Iterator) []T {
	list := []T{}
	if iter.Error() != nil {
		return list
	}
	defer iter.Close()

	for {
		var v T
		if iter.Next(&v) {
			list = append(list, v)
		} else {
			break
		}
	} //for

	return list
}
