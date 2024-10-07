package mongo

import (
	"jtools/dbg"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewInsertOneModel() *mongo.InsertOneModel {
	return mongo.NewInsertOneModel()
}

func NewDeleteOneModel() *mongo.DeleteOneModel {
	return mongo.NewDeleteOneModel()
}

func NewDeleteManyModel() *mongo.DeleteManyModel {
	return mongo.NewDeleteManyModel()
}

// NewReplaceOneModel creates a new ReplaceOneModel.
func NewReplaceOneModel() *mongo.ReplaceOneModel {
	return mongo.NewReplaceOneModel()
}

// NewUpdateOneModel creates a new UpdateOneModel.
func NewUpdateOneModel() *mongo.UpdateOneModel {
	return mongo.NewUpdateOneModel()
}

// NewUpdateManyModel creates a new UpdateManyModel.
func NewUpdateManyModel() *mongo.UpdateManyModel {
	return mongo.NewUpdateManyModel()
}

//////////////////////////////////////////////////////////////////////////

type _WriteCounter struct {
	insert_one_count int64
	delete_one_count int64
	update_one_count int64
}

func (my _WriteCounter) InsertCount() int64 { return my.insert_one_count }
func (my _WriteCounter) DeleteCount() int64 { return my.delete_one_count }
func (my _WriteCounter) UpdateCount() int64 { return my.update_one_count }

func (my _WriteCounter) getMap() map[string]any {
	return map[string]any{
		"insert_one_count": my.insert_one_count,
		"delete_one_count": my.delete_one_count,
		"update_one_count": my.update_one_count,
	}
}

type WriteModeler struct {
	list []mongo.WriteModel
	wc   _WriteCounter
}

func (my WriteModeler) Count() int { return len(my.list) }

func NewWriteModels() *WriteModeler {
	return &WriteModeler{
		list: []mongo.WriteModel{},
	}
}
func (my *WriteModeler) Append(wms ...mongo.WriteModel) {
	for _, v := range wms {
		switch item := v.(type) {
		case *mongo.InsertOneModel:
			my.wc.insert_one_count++

		case *mongo.DeleteOneModel:
			my.wc.delete_one_count++

		case *mongo.UpdateOneModel:
			if item.Upsert != nil {
				if !*item.Upsert {
					my.wc.update_one_count++
				}
			} else {
				my.wc.update_one_count++
			}
		} //switch
	} //for
	my.list = append(my.list, wms...)
}

// Ex : Update to WriteModelEx
func (my *WriteModeler) Ex(collectionName string) *WriteModelEx {
	return &WriteModelEx{
		WriteModeler:   my,
		CollectionName: collectionName,
	}
}

type WriteModelEx struct {
	*WriteModeler
	CollectionName string
}

func (my *WriteModelEx) Insert(doc interface{}) {
	my.Append(NewInsertOneModel().SetDocument(doc))
}

func (my *WriteModelEx) Update(filter, update interface{}) {
	item := NewUpdateOneModel()
	item.SetFilter(filter).SetUpdate(update)
	my.Append(item)
}

func (my *WriteModelEx) Upsert(filter, update interface{}) {
	item := NewUpdateOneModel()
	item.SetUpsert(true)
	item.SetFilter(filter).SetUpdate(update)
	my.Append(item)
}

func (my *WriteModelEx) UpdateAll(filter, update interface{}) {
	my.Append(
		NewUpdateManyModel().
			SetFilter(filter).
			SetUpdate(update),
	)
}
func (my *WriteModelEx) Delete(filter interface{}) {
	my.Append(
		NewDeleteOneModel().SetFilter(filter),
	)
	NewDeleteManyModel().SetFilter(filter)
}
func (my *WriteModelEx) DeleteAll(filter interface{}) {
	my.Append(
		NewDeleteManyModel().SetFilter(filter),
	)
}

func (my WriteModelEx) BulkWriteDB(db DATABASE, ordered ...bool) BulkWriteResult {
	if my.Count() <= 0 {
		return BulkWriteResult{}
	}
	return db.C(my.CollectionName).BulkWrite(my.WriteModeler, ordered...)
}

func (my WriteModelEx) BulkWriteSliceDB(db DATABASE, cutCnt int, ordered ...bool) BulkWriteResult {
	if len(my.WriteModeler.list) <= cutCnt {
		return my.BulkWriteDB(db, ordered...)
	}

	if my.Count() <= 0 {
		return BulkWriteResult{}
	}

	ssl := dbg.CutSliceSlice[mongo.WriteModel](my.WriteModeler.list, cutCnt)

	result := BulkWriteResult{
		BulkWriteResult: &mongo.BulkWriteResult{
			UpsertedIDs: map[int64]interface{}{},
		},
		req_filter_count: int64(len(my.WriteModeler.list)),
		req_wc:           my.wc,
	}
	c := db.C(my.CollectionName)
	for i, write_model_list := range ssl {
		_ = i
		//cc.Gray("[", i, "] :", len(write_model_list))

		wm := &WriteModeler{
			list: write_model_list,
		}

		r := c.BulkWrite(wm, ordered...)
		if r.err != nil {
			result.err = r.err
			return result
		}

		result.InsertedCount += r.InsertedCount
		result.MatchedCount += r.MatchedCount
		result.ModifiedCount += r.ModifiedCount
		result.UpsertedCount += r.UpsertedCount
		if r.UpsertedIDs != nil {
			for k, v := range r.UpsertedIDs {
				result.UpsertedIDs[k] = v
			}
		}

	} //for

	return result
}

//////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////

func (my *cCollection) BulkWrite(
	write_modeler *WriteModeler,
	ordered ...bool,
) BulkWriteResult {
	if write_modeler == nil {
		return BulkWriteResult{
			err: dbg.Error("write_modeler is nil"),
		}
	}

	is_ordered := isTrue(ordered)

	opt := options.BulkWrite()
	if !is_ordered {
		opt = opt.SetOrdered(false)
	}

	r := makeBulkWriteResult(
		write_modeler,
		is_ordered,
	)

	bwr, err := my.col.BulkWrite(
		my.session.ctx(),
		write_modeler.list,
		opt,
	)
	r.BulkWriteResult = bwr
	r.err = err

	return r
}

type BulkWriteResult struct {
	*mongo.BulkWriteResult `bson:",inline" json:",inline"`

	req_filter_count int64
	req_wc           _WriteCounter

	is_ordered bool
	err        error
}

func (my BulkWriteResult) String() string {
	if my.err != nil {
		return my.err.Error()
	}
	data := map[string]any{
		"bulk_result":      my.BulkWriteResult,
		"req_filter_count": my.req_filter_count,
		"req_wc":           my.req_wc.getMap(),
		"is_ordered":       my.is_ordered,
	}
	return dbg.ToJsonString(data)
}

func makeBulkWriteResult(
	write_modeler *WriteModeler,
	is_ordered bool,
) BulkWriteResult {
	return BulkWriteResult{
		req_filter_count: int64(len(write_modeler.list)),
		req_wc:           write_modeler.wc,
		is_ordered:       is_ordered,
	}
}
func (my BulkWriteResult) Error() error {
	return my.err
}
func (my BulkWriteResult) IsOrdered() bool  { return my.is_ordered }
func (my BulkWriteResult) FilterCount() int { return int(my.req_filter_count) }

func (my BulkWriteResult) ReqData() _WriteCounter {
	return my.req_wc
}

func (my BulkWriteResult) ValidInsertCount() bool {
	if my.err != nil {
		return false
	}
	return my.InsertedCount == my.req_wc.insert_one_count
}

func (my BulkWriteResult) ValidDeleteCount() bool {
	if my.err != nil {
		return false
	}
	return my.DeletedCount == my.req_wc.delete_one_count
}
func (my BulkWriteResult) ValidUpdateCount() bool {
	if my.err != nil {
		return false
	}
	return my.ModifiedCount >= my.req_wc.update_one_count
}

func (my BulkWriteResult) TestValidCount() bool {
	if my.err != nil {
		return false
	}

	if !my.ValidDeleteCount() {
		return false
	}

	if !my.ValidInsertCount() {
		return false
	}

	if !my.ValidUpdateCount() {
		return false
	}

	return true

	/*
		필터가 적중해야 하는것들
		update , delete

		update시 db의 문서와 갱신할 데이터가 동일하다면 modify가 일어나지 않는다.



	*/

	// my.req_filter_count -= (my.req_delete_one_count + my.req_insert_one_count)
	// tot_cnt := my.MatchedCount + //The number of documents matched by filters in update and replace operations.
	// 	//my.ModifiedCount + //The number of documents modified by update and replace operations.
	// 	my.UpsertedCount //The number of documents upserted by update and replace operations.

	// return tot_cnt == my.req_filter_count
}
