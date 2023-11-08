package mongo

import (
	"errors"
	"fmt"
	"strings"
	"txscheduler/brix/tools/database/mongo/tools/cc"
	"txscheduler/brix/tools/database/mongo/tools/dbg"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection interface {
	DBName() string
	Name() string
	Insert(docs ...interface{}) error

	UpdateID(id, update interface{}) error
	UpdateID2(id, update interface{}) UpdateResult

	Update(selector, update interface{}) error
	Update2(selector, update interface{}) UpdateResult

	UpdateAll(selector, update interface{}) (UpdateResult, error)
	UpdateAll2(selector, update interface{}) UpdateResult

	Upsert(selector, update interface{}) (UpdateResult, error)
	Upsert2(selector, update interface{}) UpdateResult

	Remove(selector interface{}) error
	Remove2(selector interface{}) DeleteResult

	RemoveAll(selector interface{}) (DeleteResult, error)
	RemoveAll2(selector interface{}) DeleteResult

	Count(selectors ...interface{}) (int, error)
	DropCollection() error

	Find(selector interface{}) Query
	FindProjection(selector interface{}, projection interface{}) Query

	FindOneAndUpdate(selector, update interface{}, result_before_document interface{}) error
	FindOneAndDelete(selector interface{}, result_remove_document interface{}) error
	FindOneAndReplace(selector, replace interface{}, result_before_document interface{}) error

	Aggregate(pipeline interface{}) Iterator

	Indexes() (IndexViewList, error)
	EnsureIndex(data IndexData) error
	DropIndex(names ...string) error
	DropIndexAll() error

	//Sharding [SCollection]
	SelectorInsert(v SelectorInterface) error
	ShardCollection(key string, val interface{}) error
	ShardCollectionHashed(key string) error
	ShardCollectionIDHashed() error
}

type cCollection struct {
	cdb     *CDB
	session *dbSession
	dbName  string
	col     *mongo.Collection
}

func (my *cCollection) DBName() string { return my.dbName }
func (my *cCollection) Name() string   { return my.col.Name() }

func (my *cCollection) Insert(docs ...interface{}) error {
	_, err := my.col.InsertMany(
		my.session.ctx(),
		docs,
	)
	return err
}

func getFilter(selector interface{}) interface{} {
	var filter interface{}
	if selector == nil {
		filter = Dson{}
	} else {
		filter = selector
	}
	return filter
}

func checkDalar(data interface{}) bool {
	query := VOID{}
	dbg.ChangeStruct(data, &query)
	for key, _ := range query {
		if strings.HasPrefix(key, "$") {
			return true
		}
	}
	return false
}

// /////////////////////////////////////////////////////////////////////////////////////
type UpdateResult struct {
	*mongo.UpdateResult `bson:",inline" json:",inline"`
	Error               error `bson:"error" json:"error"`
}

func makeUpdateResult() UpdateResult {
	return UpdateResult{
		UpdateResult: &mongo.UpdateResult{},
		Error:        nil,
	}
}
func (my UpdateResult) String() string { return dbg.ToJSONString(my) }
func (my UpdateResult) Valid() bool {
	if my.UpdateResult == nil || my.Error != nil {
		return false
	}
	return my.MatchedCount > 0 || my.ModifiedCount > 0 || my.UpsertedCount > 0
}

func (my UpdateResult) ValidError() error {
	if my.Error != nil {
		return my.Error
	}
	if my.UpdateResult == nil {
		return dbg.Error("updateResult is Nil")
	}

	if my.MatchedCount == 0 && my.ModifiedCount == 0 && my.UpsertedCount == 0 {
		return dbg.Error("UpdateResult 3 is ZERO")
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////

func (my *cCollection) Update(selector, update interface{}) error {
	r := my.Update2(selector, update)
	return r.Error
}

func (my *cCollection) Update2(selector, update interface{}) UpdateResult {

	result := VOID{}
	query := VOID{}
	dbg.ChangeStruct(update, &query)

	var r *mongo.UpdateResult
	var err error

	if checkDalar(update) {
		r, err = my.col.UpdateOne(
			my.session.ctx(),
			getFilter(selector),
			update,
		)
		dbg.ChangeStruct(r, &result)
		//cc.Green(result)
	} else {
		r, err = my.col.ReplaceOne(
			my.session.ctx(),
			getFilter(selector),
			update,
		)
	}

	if err == nil {
		result := makeUpdateResult()
		dbg.ChangeStruct(r, &result.UpdateResult)
		return result
	} else {
		cc.Red("collection.update____ :", err)
		cc.Red("selector :", selector)
		cc.Red("update :", update)
		cc.Red("-----------------------------")
	}

	return UpdateResult{
		Error: err,
	}
}

func (my *cCollection) UpdateAll(selector, update interface{}) (UpdateResult, error) {
	r, err := my.col.UpdateMany(
		my.session.ctx(),
		getFilter(selector),
		update,
	)
	if err == nil {
		result := makeUpdateResult()
		dbg.ChangeStruct(r, &result.UpdateResult)
		return result, err
	}
	return UpdateResult{
		Error: err,
	}, err
}
func (my *cCollection) UpdateAll2(selector, update interface{}) UpdateResult {
	r, err := my.col.UpdateMany(
		my.session.ctx(),
		getFilter(selector),
		update,
	)
	if err == nil {
		result := makeUpdateResult()
		dbg.ChangeStruct(r, &result.UpdateResult)
		return result
	}
	return UpdateResult{
		Error: err,
	}
}

func (my *cCollection) UpdateID(id, update interface{}) error {
	r := my.UpdateID2(id, update)
	return r.Error
}

func (my *cCollection) UpdateID2(id, update interface{}) UpdateResult {
	var keyID interface{}

	switch v := id.(type) {
	case string:
		v2, err1 := primitive.ObjectIDFromHex(v)
		if err1 != nil {
			return UpdateResult{
				Error: err1,
			}
		}
		keyID = v2

	case primitive.ObjectID:
		keyID = v
	}
	r, err := my.col.UpdateByID(
		my.session.ctx(),
		keyID,
		update,
	)
	if err == nil {
		result := makeUpdateResult()
		dbg.ChangeStruct(r, &result.UpdateResult)
		return result
	}
	return UpdateResult{
		Error: err,
	}
}

func (my *cCollection) Upsert(selector, update interface{}) (UpdateResult, error) {

	var r *mongo.UpdateResult
	var err error

	if checkDalar(update) {
		r, err = my.col.UpdateOne(
			my.session.ctx(),
			getFilter(selector),
			update,
			options.Update().SetUpsert(true),
		)

	} else {
		r, err = my.col.ReplaceOne(
			my.session.ctx(),
			getFilter(selector),
			update,
			options.Replace().SetUpsert(true),
		)
	}

	if err == nil {
		result := makeUpdateResult()
		dbg.ChangeStruct(r, &result.UpdateResult)
		return result, err
	}

	return UpdateResult{
		Error: err,
	}, err
}

func (my *cCollection) Upsert2(selector, update interface{}) UpdateResult {

	var r *mongo.UpdateResult
	var err error

	if checkDalar(update) {
		r, err = my.col.UpdateOne(
			my.session.ctx(),
			getFilter(selector),
			update,
			options.Update().SetUpsert(true),
		)

	} else {
		r, err = my.col.ReplaceOne(
			my.session.ctx(),
			getFilter(selector),
			update,
			options.Replace().SetUpsert(true),
		)
	}

	if err == nil {
		result := makeUpdateResult()
		dbg.ChangeStruct(r, &result.UpdateResult)
		return result
	}

	return UpdateResult{
		Error: err,
	}
}

// FindOneAndUpdate : 업데이트 하기전의 문서 반환
func (my *cCollection) FindOneAndUpdate(selector, update interface{}, result_before_document interface{}) error {
	r := my.col.FindOneAndUpdate(
		my.session.ctx(),
		getFilter(selector),
		update,
	)
	if r == nil {
		return errors.New("[FindOneAndUpdate] r is nil")
	}
	if err := r.Err(); err != nil {
		return err
	}
	if err := r.Decode(result_before_document); err != nil {
		return err
	}
	return nil
}

// FindOneAndDelete : 삭제하고 삭제된 문서 반환
func (my *cCollection) FindOneAndDelete(selector interface{}, result_remove_document interface{}) error {
	r := my.col.FindOneAndDelete(
		my.session.ctx(),
		getFilter(selector),
	)
	if r == nil {
		return errors.New("[FindOneAndDelete] r is nil")
	}
	if err := r.Err(); err != nil {
		return err
	}
	if err := r.Decode(result_remove_document); err != nil {
		return err
	}
	return nil
}

func (my *cCollection) FindOneAndReplace(selector, replace interface{}, result_before_document interface{}) error {
	r := my.col.FindOneAndReplace(
		my.session.ctx(),
		getFilter(selector),
		replace,
	)
	if r == nil {
		return errors.New("[FindOneAndReplace] r is nil")
	}
	if err := r.Err(); err != nil {
		return err
	}
	if err := r.Decode(result_before_document); err != nil {
		return err
	}
	return nil
}

// /////////////////////////////////////////////////////////////////////////////////////
type DeleteResult struct {
	*mongo.DeleteResult `bson:",inline" json:",inline"`
	Error               error `bson:"error" json:"error"`
}

func (my DeleteResult) String() string { return dbg.ToJSONString(my) }
func makeDeleteResult() DeleteResult {
	return DeleteResult{
		DeleteResult: &mongo.DeleteResult{},
	}
}
func (my DeleteResult) Valid() bool {
	if my.DeleteResult == nil || my.Error != nil {
		return false
	}
	return my.DeletedCount > 0
}

///////////////////////////////////////////////////////////////////////////////////////

func (my *cCollection) Remove(selector interface{}) error {
	r := my.Remove2(selector)
	return r.Error
}
func (my *cCollection) Remove2(selector interface{}) DeleteResult {
	r, err := my.col.DeleteOne(
		my.session.ctx(),
		getFilter(selector),
	)
	if err == nil {
		void := makeDeleteResult()
		dbg.ChangeStruct(r, &void.DeleteResult)
		return void
	}
	return DeleteResult{
		Error: err,
	}
}

func (my *cCollection) RemoveAll(selector interface{}) (DeleteResult, error) {
	r, err := my.col.DeleteMany(
		my.session.ctx(),
		getFilter(selector),
	)
	if err != nil {
		void := makeDeleteResult()
		dbg.ChangeStruct(r, &void.DeleteResult)
		return void, err
	}
	return DeleteResult{
		Error: err,
	}, err
}

func (my *cCollection) RemoveAll2(selector interface{}) DeleteResult {
	r, err := my.col.DeleteMany(
		my.session.ctx(),
		getFilter(selector),
	)
	if err != nil {
		void := makeDeleteResult()
		dbg.ChangeStruct(r, &void.DeleteResult)
		return void
	}
	return DeleteResult{
		Error: err,
	}
}

func (my *cCollection) Count(selectors ...interface{}) (int, error) {
	var filter interface{}
	if len(selectors) > 0 {
		filter = getFilter(selectors[0])
	} else {
		filter = getFilter(nil)
	}

	// cnt, err := my.col.CountDocuments(
	// 	my.session.ctx(),
	// 	filter,
	// )
	// return int(cnt), err

	return my.cdb.CollectionCount(
		my.dbName,
		my.Name(),
		filter,
		0, 0,
	)
}

func (my *cCollection) DropCollection() error {
	return my.col.Drop(my.session.ctx())
}

///////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////
//////////////////////////// [Shard] //////////////////////////////////
///////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////

type SelectorInterface interface {
	Selector() Bson
}

func (my *cCollection) SelectorInsert(v SelectorInterface) error {
	if v == nil {
		return errors.New("[SelectorInsert] nil pointer")
	}
	if cnt, _ := my.Find(v.Selector()).Count(); cnt == 0 {
		return my.Insert(v)
	}
	return fmt.Errorf("[SelectorInsert] Duplicatied : %v", v.Selector())
}

func (my *cCollection) ShardCollection(key string, val interface{}) error {
	return my.cdb.ShardCollection(
		my.dbName,
		my.Name(),
		key,
		val,
	)
}

func (my *cCollection) ShardCollectionHashed(key string) error {
	return my.cdb.ShardCollectionHashed(
		my.dbName,
		my.Name(),
		key,
	)
}
func (my *cCollection) ShardCollectionIDHashed() error {
	return my.cdb.ShardCollectionIDHashed(
		my.dbName,
		my.Name(),
	)
}
