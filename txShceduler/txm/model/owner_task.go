package model

import (
	"jtools/mms"
	"jtools/unix"
	"strings"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/txm/inf"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type STATE int

const (
	NONE    = STATE(0)
	PENDING = STATE(1)

	MASTER_NONE    = STATE(10)
	MASTER_PENDING = STATE(11)

	FAIL    = STATE(104)
	SUCCESS = STATE(200)
)

type TaskKey string

func MakeOwnerTask_LT_Key() TaskKey {
	return TaskKey("OWNER" + strings.ToUpper(primitive.NewObjectID().Hex()) + "LT")
}
func MakeOwnerTask_TransferKey() TaskKey {
	return TaskKey("OWNER" + strings.ToUpper(primitive.NewObjectID().Hex()) + "TK")
}
func MakeOwnerTask_LockKey() TaskKey {
	return TaskKey("OWNER" + strings.ToUpper(primitive.NewObjectID().Hex()) + "LK")
}
func MakeOwnerTask_UnlockKey() TaskKey {
	return TaskKey("OWNER" + strings.ToUpper(primitive.NewObjectID().Hex()) + "UL")
}
func MakeOwnerTask_RelockKey() TaskKey {
	return TaskKey("OWNER" + strings.ToUpper(primitive.NewObjectID().Hex()) + "RL")
}

type OwnerTaskKind string

const (
	OwnerTaskKind_LT       = OwnerTaskKind("LT") //lock_transfer
	OwnerTaskKind_Transfer = OwnerTaskKind("transfer")
	OwnerTaskKind_Lock     = OwnerTaskKind("lock")
	OwnerTaskKind_Unlock   = OwnerTaskKind("unlock")
	OwnerTaskKind_Relock   = OwnerTaskKind("relock")
)

type OwnerTask struct {
	Key  TaskKey       `bson:"key" json:"key"`
	Kind OwnerTaskKind `bson:"kind" json:"kind"`

	State STATE `bson:"state" json:"state"` // 0:ready , 1:action, 200:end

	Task mongo.MAP `bson:"task" json:"task"`

	CreateAt unix.Time `bson:"create_at" json:"create_at"`
	UpdateAt unix.Time `bson:"update_at" json:"update_at"`

	ZZ_CREATE_KST string `bson:"zz_create_kst" json:"-"`
	ZZ_UPDATE_KST string `bson:"zz_update_kst" json:"-"`
}

func (my OwnerTask) Valid() bool          { return my.Key != "" }
func (my OwnerTask) String() string       { return dbg.ToJSONString(my) }
func (my OwnerTask) Selector() mongo.Bson { return mongo.Bson{"key": my.Key} }

func (OwnerTask) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.OwnerTask)

		c.EnsureIndex(mongo.SingleIndex("key", 1, true))
		c.EnsureIndex(mongo.SingleIndex("kind", 1, false))
		c.EnsureIndex(mongo.SingleIndex("state", 1, false))

		c.EnsureIndex(mongo.MultiIndexName(
			[]interface{}{
				"state", 1,
				"create_at", 1,
			},
			false,
			"m_index_1",
			0,
		))

		prefixDot := "task."
		Owner_LT_Task{}.indexingC(c, prefixDot)
		OwnerTransferTask{}.indexingC(c, prefixDot)
		OwnerLockTask{}.indexingC(c, prefixDot)
		OwnerUnlockTask{}.indexingC(c, prefixDot)
		OwnerRelockTask{}.indexingC(c, prefixDot)

	})
}

func (my OwnerTask) Owner_LT_Task() Owner_LT_Task {
	task, _ := dbg.DecodeStruct[Owner_LT_Task](my.Task)
	return task
}

func (my OwnerTask) OwnerTransferTask() OwnerTransferTask {
	task, _ := dbg.DecodeStruct[OwnerTransferTask](my.Task)
	return task
}

func (my OwnerTask) OwnerLockTask() OwnerLockTask {
	task, _ := dbg.DecodeStruct[OwnerLockTask](my.Task)
	return task
}

func (my OwnerTask) OwnerUnlockTask() OwnerUnlockTask {
	task, _ := dbg.DecodeStruct[OwnerUnlockTask](my.Task)
	return task
}

func (my OwnerTask) OwnerRelockTask() OwnerRelockTask {
	task, _ := dbg.DecodeStruct[OwnerRelockTask](my.Task)
	return task
}

type TaskEndChecker interface {
	IsStateEnd() bool
}

func (my *OwnerTask) UpdateDB(db mongo.DATABASE, task_any TaskEndChecker) {
	if task_any.IsStateEnd() {
		my.State = SUCCESS
	}

	task, _ := dbg.DecodeStruct[mongo.MAP](task_any)
	my.Task = task

	my.UpdateAt = unix.Now()
	my.ZZ_UPDATE_KST = my.UpdateAt.KST()

	db.C(inf.OwnerTask).Update(
		my.Selector(),
		my,
	)
}

///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

type OwnerLog struct {
	Key  TaskKey       `bson:"key" json:"key"`
	Kind OwnerTaskKind `bson:"kind" json:"kind"`

	Log mongo.MAP `bson:"log" json:"log"`

	Timestamp unix.Time `bson:"timestamp" json:"timestamp"`
	YMD       int       `bson:"ymd" json:"ymd"`

	SendAt  mms.MMS `bson:"send_at" json:"send_at"`
	SendYMD int     `bson:"send_ymd" json:"send_ymd"`
	IsSend  bool    `bson:"is_send" json:"is_send"`
}

type LockLogList []OwnerLog

func (my OwnerLog) String() string    { return dbg.ToJSONString(my) }
func (my LockLogList) String() string { return dbg.ToJSONString(my) }

func (my OwnerLog) Valid() bool          { return my.Key != "" }
func (my OwnerLog) Selector() mongo.Bson { return mongo.Bson{"key": my.Key} }

func (OwnerLog) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.OwnerLog)
		c.EnsureIndex(mongo.SingleIndex("key", "1", true))
		c.EnsureIndex(mongo.SingleIndex("kind", 1, false))

		c.EnsureIndex(mongo.SingleIndex("is_send", "1", false))

		c.EnsureIndex(mongo.MultiIndexName(
			[]interface{}{
				"is_send", 1,
				"timestamp", 1,
			},
			false,
			"m_index_1",
			0,
		))

	})
}

func (OwnerLog) GetList(db mongo.DATABASE) LockLogList {
	list := LockLogList{}
	db.C(inf.OwnerLog).
		Find(mongo.Bson{"is_send": false}).
		Sort("timestamp").
		All(&list)
	return list
}

func (my OwnerLog) SendOK(db mongo.DATABASE, nowAt mms.MMS) {
	upQuery := mongo.Bson{"$set": mongo.Bson{
		"send_at":  nowAt,
		"send_ymd": nowAt.YMD(),
		"is_send":  true,
	}}
	db.C(inf.OwnerLog).Update(my.Selector(), upQuery)
}

type LockLoger interface {
	GetKey() TaskKey
	Kind() OwnerTaskKind
	GetTimestamp() unix.Time
}

func (my OwnerLog) InsertDB(db mongo.DATABASE, log_any LockLoger) {

	my.Key = log_any.GetKey()
	my.Kind = log_any.Kind()

	my.Timestamp = log_any.GetTimestamp()
	my.YMD = my.Timestamp.YMD()

	log, _ := dbg.DecodeStruct[mongo.MAP](log_any)
	my.Log = log

	my.SendAt = mms.Zero()
	my.SendYMD = 0
	my.IsSend = false

	db.C(inf.OwnerLog).Insert(my)
}

////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////

type OwnerUnlockPool struct {
	Address  string    `bson:"address" json:"address"`
	InsertAt unix.Time `bson:"insert_at" json:"insert_at"`
}

func (my OwnerUnlockPool) Valid() bool { return my.Address != "" }

func (OwnerUnlockPool) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.OwnerUnlockPool)
		c.EnsureIndex(mongo.SingleIndex("address", 1, true))
		c.EnsureIndex(mongo.SingleIndex("insert_at", 1, false))
	})
}

func (my OwnerUnlockPool) InsertDB(db mongo.DATABASE, address string) error {
	my.Address = address
	my.InsertAt = unix.Now()

	return db.C(inf.OwnerUnlockPool).Insert(my)
}
func (OwnerUnlockPool) RemoveDB(db mongo.DATABASE, address string) {
	db.C(inf.OwnerUnlockPool).Remove(mongo.Bson{"address": address})
}
