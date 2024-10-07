package mongo

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"time"
	"txscheduler/brix/tools/database/mongo/tools/cc"
	"txscheduler/brix/tools/database/mongo/tools/dbg"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type dbInfo struct {
	url_list   []string
	isAuth     bool
	username   string
	pwd        string
	authSource string // "":admin
}

func (my dbInfo) GetURLList() []string    { return my.url_list }
func (my dbInfo) IsAuth() bool            { return my.isAuth }
func (my dbInfo) UserName() string        { return my.username }
func (my dbInfo) PWD() string             { return my.pwd }
func (my dbInfo) AuthSource() string      { return my.authSource }
func (my dbInfo) ClusterEndpoint() string { return "" }
func (my dbInfo) CaFilePath() string      { return "" }
func (my dbInfo) ReadPreference() string  { return "" }

type iInfo interface {
	GetURLList() []string
	IsAuth() bool
	UserName() string
	PWD() string
	AuthSource() string
	ClusterEndpoint() string
	CaFilePath() string
}

type dbSession struct {
	iInfo
	client *mongo.Client
	cctx   context.Context
}

func (my *dbSession) Close() error {
	return my.client.Disconnect(context.Background())
}

func (my dbSession) ctx() context.Context {
	if my.cctx != nil {
		return my.cctx
	}
	return context.Background()
}

type CDB struct {
	*dbSession
	sscnt      int64
	authSource string
}

func (my *CDB) GetAuthSource() string { return my.authSource }
func (my *CDB) IsAuthSource() bool    { return my.authSource != "" }

func (my *CDB) SessionCount() int64 {
	return atomic.LoadInt64(&my.sscnt)
}

func (my *CDB) copySession() *dbSession {
	atomic.AddInt64(&my.sscnt, 1)
	return my.dbSession
	//return newSession(my.dbInfo)
}

func (my *CDB) closeSession(session *dbSession) {
	//session.Close()
	atomic.AddInt64(&my.sscnt, -1)
}

func (my *CDB) GetSession() *dbSession { return my.dbSession }

// Ping : default waitDu : 3sec
func (my *CDB) Ping(waitDu ...time.Duration) error {
	var ctx context.Context
	if len(waitDu) == 0 {
		ctx, _ = context.WithTimeout(context.Background(), time.Second*3)
	} else {
		ctx, _ = context.WithTimeout(context.Background(), waitDu[0])
	}
	return my.client.Ping(ctx, nil)
}

func New(connectionAddress string, isAuth bool, auth_source string, id, pwd string) *CDB {
	return newCDB(
		[]string{connectionAddress},
		isAuth,
		auth_source,
		id, pwd,
	)
}

func NewList(address_list []string, isAuth bool, auth_source string, id, pwd string) *CDB {
	return newCDB(
		address_list,
		isAuth,
		auth_source,
		id, pwd,
	)
}

func newCDB(list []string, isAuth bool, auth_source string, id, pwd string) *CDB {
	for i := range list {
		list[i] = strings.TrimSpace(list[i])
		if strings.HasPrefix(list[i], "mongodb://") {
			list[i] = strings.Replace(list[i], "mongodb://", "", 1)
		}
	} //for

	info := dbInfo{
		url_list:   list,
		isAuth:     isAuth,
		authSource: auth_source,
	}
	if info.isAuth {
		info.username = id
		info.pwd = pwd
	}

	session := newSession(info)

	if session == nil {
		return nil
	}

	return &CDB{
		dbSession:  session,
		authSource: auth_source,
	}
}

func newSession(info dbInfo) *dbSession {
	//read_pref_opt, _ := readpref.New(readpref.SecondaryPreferredMode)
	opt := options.Client().
		SetHosts(
			info.url_list,
		).
		//SetReadPreference(read_pref_opt).
		SetAuth(
			options.Credential{
				//AuthMechanism: "SCRAM-SHA-256",
				//Username: url.QueryEscape(info.username),
				//Password: url.QueryEscape(info.pwd),
				Username:   info.username,
				Password:   info.pwd,
				AuthSource: info.authSource,
			},
		)
	is_direct := len(info.url_list) == 1
	opt.SetDirect(is_direct)

	// readpref.Primary()
	// readpref.PrimaryPreferred()
	// readpref.Secondary()
	// readpref.SecondaryPreferred()
	// readpref.Nearest()
	//opt.SetReadPreference(readpref.Primary())

	client, err := mongo.NewClient(opt)
	if err != nil {
		cc.RedItalicBG("[", info.url_list, "]", err)
		return nil
	}
	if err := client.Connect(context.Background()); err != nil {
		cc.Red("newSession :", err)
		return nil
	}
	return &dbSession{
		iInfo:  info,
		client: client,
	}
}

func (my *CDB) Session(callback func(*dbSession)) {
	session := my.copySession()
	defer my.closeSession(session)
	callback(session)
}

func (my *CDB) DatabaseNames() (names []string, err error) {
	session := my.copySession()
	defer my.closeSession(session)
	return session.client.ListDatabaseNames(session.ctx(), bson.M{})

}

func (my *CDB) CollectionNames(dbName string) (names []string, err error) {
	session := my.copySession()
	defer my.closeSession(session)

	db := session.client.Database(dbName)
	return db.ListCollectionNames(session.ctx(), bson.M{})
}

func (my *CDB) DropDatabase(dbName string) error {
	session := my.copySession()
	defer my.closeSession(session)

	db := session.client.Database(dbName)
	return db.Drop(session.ctx())
}

func (my *CDB) Action(dbName string, callback func(db DATABASE)) (callback_error error) {
	defer func() {
		if e := recover(); e != nil {
			dbg.StackError(e)
			callback_error = dbg.Error(e)
		}
	}()

	session := my.copySession()
	defer my.closeSession(session)

	database := &Database{
		cdb:      my,
		session:  session,
		Database: session.client.Database(dbName),
	}

	callback(database)
	return nil
}

func (my *CDB) Transaction(dbName string, callback func(db DATABASE) bool, limitDu ...time.Duration) error {
	session := my.copySession()
	defer my.closeSession(session)

	opt := options.SessionOptions{}
	limitDuration := time.Second * 60
	if len(limitDu) > 0 && limitDu[0] > 0 {
		limitDuration = limitDu[0]
	}
	opt.SetDefaultMaxCommitTime(&limitDuration)

	opt.SetDefaultReadPreference(readpref.Primary())

	err := session.client.UseSessionWithOptions(
		context.Background(),
		&opt,
		func(sc mongo.SessionContext) error {

			if err := sc.StartTransaction(); err != nil {
				return err
			}
			cc.Purple("StartTransaction")

			dbss := &dbSession{
				iInfo:  my.iInfo,
				client: sc.Client(),
				cctx:   sc,
			}
			database := &Database{
				cdb:      my,
				session:  dbss,
				Database: dbss.client.Database(dbName),
			}
			if ok := callback(database); !ok {
				err := sc.AbortTransaction(sc)
				cc.Purple("AbortTransaction :", err)

				return err

			} else {
				err := sc.CommitTransaction(sc)
				cc.Purple("CommitTransaction :", err)

				return err
			}

		},
	)

	return err
}

func (my *CDB) Run(dbName, colName string, callback func(c Collection)) {
	my.Action(dbName, func(db DATABASE) {
		callback(db.C(colName))
	})
}

func (my *CDB) CollectionCount(
	dbName, colName string,
	filter interface{},
	skip, limit int64,
) (int, error) {
	db := my.client.Database(dbName)

	eFilter := bson.E{}
	if filter != nil {
		eFilter = bson.E{
			"query", filter,
		}
	}

	opt := options.RunCmd()
	sr := db.RunCommand(
		context.Background(),
		Dson{
			{"count", colName},
			// bson.E{
			// 	"query",
			// 	//bson.M{"key": bson.M{"$gte": 7}},
			// },
			eFilter,
			bson.E{"skip", skip},
			bson.E{"limit", limit},
		},
		//Bson{"count": "aaa"},
		opt,
	)

	if sr.Err() != nil {
		return 0, sr.Err()
	}

	/*
		result {
			"shards"		[{"shard_name" : int32(count)} , ...]
			"n" 			int32
			"ok"			float64
			"$clusterTime"	[{clusterTime} , {signature}]
			"operationTime"	interface{}
		}
	*/
	result := VOID{}
	if err := sr.Decode(&result); err != nil {
		return 0, err
	}
	if v, do := result["ok"]; do {
		if ok, do := v.(float64); do && ok == 1 {
			if cnt, do := result["n"]; do {
				if total, do := cnt.(int32); do {
					return int(total), nil
				}
			}
		}
	}
	return 0, errors.New("not ok")
}

// ShardingDBList : 샤딩된 DB 리스트
func (my *CDB) ShardingDBList() []string {
	list := []string{}
	my.Action("config", func(db DATABASE) {
		voids := VOIDS{}
		db.C("databases").Find(nil).All(&voids)
		for _, v := range voids {
			list = append(list, v.IDString())
		} //for
	})
	return list
}

func (my *CDB) IsShardingDBCheck() bool {
	list := my.ShardingDBList()
	return len(list) > 0
}

func (my *CDB) RunAdmin(cmd interface{}) *mongo.SingleResult {
	session := my.copySession()
	defer my.closeSession(session)

	opt := &options.RunCmdOptions{}
	db := session.client.Database("admin")
	sr := db.RunCommand(
		session.ctx(),
		cmd,
		opt,
	)

	return sr
}

func (my *CDB) EnableSharding(dbName string) error {
	list := my.ShardingDBList()
	for _, v := range list {
		if v == dbName {
			cc.Yellow("Already Sharding DB :", dbName)
			return nil
		}
	} //for

	sr := my.RunAdmin(
		Dson{
			{"enableSharding", dbName},
		},
	)

	if sr.Err() == nil {
		void := VOID{}
		sr.Decode(&void)
		//cc.Purple(void)
		cc.Cyan(dbg.Cat("EnableSharding[ ", dbName, " ]"))
	}

	return sr.Err()
}

// ShardingCollectionList : 샤딩된 컬렉션 리스트 ( []{"dbName.colName",..." 형식 )
func (my *CDB) ShardingCollectionList() []string {
	list := []string{}
	my.Action("config", func(db DATABASE) {
		voids := VOIDS{}
		db.C("collections").Find(nil).All(&voids)
		for _, v := range voids {
			if elem, do := v["dropped"]; do {
				if flag, do := elem.(bool); do {
					if flag {
						continue
					}
				}
			}
			list = append(list, v.IDString())
		}
	})
	return list
}

func (my *CDB) ShardCollection(dbName, colName string, key string, val interface{}) error {
	compareName := dbName + "." + colName
	list := my.ShardingCollectionList()
	for _, v := range list {
		if v == compareName {
			cc.Yellow("Already Sharding Collection :", compareName)
			return nil
		}
	} //for

	sr := my.RunAdmin(
		Dson{
			{"shardCollection", dbName + "." + colName},
			{"key", Bson{key: val}},
		},
	)

	if sr.Err() == nil {
		void := VOID{}
		sr.Decode(&void)
		//cc.Purple(void)
		cc.Cyan(dbg.Cat("ShardCollection[", dbName, ".", colName, "] ", key, " ", val))
	}

	return sr.Err()
}

func (my *CDB) ShardCollectionHashed(dbName, colName string, key string) error {
	return my.ShardCollection(dbName, colName, key, "hashed")
}

func (my *CDB) ShardCollectionIDHashed(dbName, colName string) error {
	return my.ShardCollection(dbName, colName, "_id", "hashed")
}

type ShardNode struct {
	ID    string `bson:"_id" json:"_id"`
	Host  string `bson:"host" json:"host"`
	State int    `bson:"state" json:"state"`
}
type ShardNodeList []ShardNode

func (my ShardNode) String() string     { return dbg.ToJSONString(my) }
func (my ShardNodeList) String() string { return dbg.ToJSONString(my) }

func (my CDB) ShardNodeList() ShardNodeList {
	list := ShardNodeList{}
	my.Action("config", func(db DATABASE) {
		db.C("shards").Find(nil).All(&list)
	})
	return list
}
