package model

import (
	"jtools/dbg"
	"jtools/mms"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/txm/inf"
)

const (
	ErrorFinderNull    = "finder_is_null"
	ErrorProcessTxlist = "process_tx_list"
)

type LogLevel int

const (
	LogError = LogLevel(0)
	LogWarn  = LogLevel(1)
	LogInfo  = LogLevel(2)
	LogDebug = LogLevel(3)
	LogTrace = LogLevel(4)
	LogCBC   = LogLevel(5)
)

type XLog struct {
	Kind      string    `bson:"kind" json:"kind"`
	Level     LogLevel  `bson:"level" json:"level"`
	Text      string    `bson:"text" json:"text"`
	Data      mongo.MAP `bson:"data" json:"data"`
	Timestamp mms.MMS   `bson:"timestamp" json:"timestamp"`
	YMD       int       `bson:"ymd" json:"ymd"`
}

func (XLog) IndexingDB() {
	DB(func(db mongo.DATABASE) {
		c := db.C(inf.XLog)
		c.EnsureIndex(mongo.SingleIndex("kind", "1", false))
		c.EnsureIndex(mongo.SingleIndex("level", "1", false))
		c.EnsureIndex(mongo.SingleIndex("timestamp", "-1", false))
		c.EnsureIndex(mongo.SingleIndex("ymd", "-1", false))

		LogInfo.WriteLog(db, "start", "indexingDB")
	})
}

func (my LogLevel) WriteLog(
	db mongo.DATABASE,
	kind string,
	text string,
) {
	at := mms.Now()
	x_log := XLog{
		Kind:      kind,
		Text:      text,
		Data:      mongo.MAP{},
		Timestamp: at,
		Level:     my,
	}
	x_log.YMD = at.YMD()

	db.C(inf.XLog).Insert(x_log)
}

func (my LogLevel) InsertLog(
	kind string,
	text string,
) {

	DB(func(db mongo.DATABASE) {
		my.WriteLog(db, kind, text)
	})
}

//////////////////////////////////////////////

func (my LogLevel) Set(
	kind string,
	data_pairs ...interface{},
) {
	DB(func(db mongo.DATABASE) {
		at := mms.Now()
		x_log := XLog{
			Kind:      kind,
			Text:      "data",
			Data:      mongo.MAP{},
			Timestamp: at,
			Level:     my,
		}
		for i := 0; i < len(data_pairs); i += 2 {
			key := dbg.Void(data_pairs[i])
			x_log.Data[key] = data_pairs[i+1]
		} //for

		x_log.YMD = at.YMD()
		db.C(inf.XLog).Insert(x_log)
	})
}
