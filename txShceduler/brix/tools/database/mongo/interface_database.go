package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DATABASE interface {
	Name() string
	GetName() string

	CollectionNames() ([]string, error)
	DropDatabase() error
	C(collectionName string) Collection

	ShardCollection(colName string, key string, val interface{}) error
	ShardCollectionHashed(colName string, key string) error
	ShardCollectionIDHashed(colName string) error
}

type Database struct {
	cdb     *CDB
	session *dbSession
	*mongo.Database
}

func (my *Database) GetName() string { return my.Name() }

func (my *Database) C(collectionName string) Collection {
	return &cCollection{
		cdb:     my.cdb,
		session: my.session,

		dbName: my.Name(),
		col:    my.Collection(collectionName),
	}
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
