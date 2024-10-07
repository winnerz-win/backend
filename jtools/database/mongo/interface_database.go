package mongo

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DATABASE interface {
	Name() string
	GetName() string
	IsTransaction() bool

	CollectionNames() ([]string, error)
	DropDatabase() error
	C(collectionName string) Collection
	CS(param_struct interface{}) Collection
	EnsureIndexStruct(param_struct interface{}, prefix ...string) error

	ShardCollection(colName string, key string, val interface{}) error
	ShardCollectionHashed(colName string, key string) error
	ShardCollectionIDHashed(colName string) error
}

type Database struct {
	cdb     *CDB
	session *dbSession
	*mongo.Database

	is_tx_mode bool //transaction mode flag
}

func (my *Database) GetName() string { return my.Name() }

// IsTransaction  : IsTransaction
func (my Database) IsTransaction() bool { return my.is_tx_mode }

func (my *Database) C(collectionName string) Collection {

	return &cCollection{
		cdb:     my.cdb,
		session: my.session,

		dbName:    my.Name(),
		col:       my.Collection(collectionName),
		read_pref: my.ReadPreference(),

		is_tx_mode: my.is_tx_mode,
	}
}

func (my *Database) CS(param_struct interface{}) Collection {
	collection_name := StructNameToLower(param_struct)
	return my.C(collection_name)
}

func (my *Database) EnsureIndexStruct(param_struct interface{}, prefix ...string) error {
	r_type, collection_name := _struct_reflect_name(reflect.TypeOf(param_struct))
	collection_name = _struct_name_tolowter(collection_name)
	col := my.C(collection_name)
	return col.(*cCollection)._ensure_index_struct(r_type, prefix...)
}

func (my *Database) DropDatabase() error {
	return my.Drop(my.session.ctx())
}

func (my *Database) CollectionNames() ([]string, error) {
	return my.ListCollectionNames(my.session.ctx(), bson.M{})
}

///////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////
//////////////////////////// [Shard] //////////////////////////////////
///////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////

func (my *Database) ShardCollection(colName string, key string, val interface{}) error {
	return my.cdb.ShardCollection(
		my.GetName(),
		colName,
		key,
		val,
	)
}

func (my *Database) ShardCollectionHashed(colName string, key string) error {
	return my.cdb.ShardCollectionHashed(
		my.GetName(),
		colName,
		key,
	)
}

func (my Database) ShardCollectionIDHashed(colName string) error {
	return my.cdb.ShardCollectionIDHashed(
		my.GetName(),
		colName,
	)
}
