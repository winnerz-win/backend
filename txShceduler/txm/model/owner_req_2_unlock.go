package model

import (
	"jtools/unix"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
)

type ReqOwnerUnlockTry struct {
	Address string `json:"address"`

	// PositionIndex int       `json:"position_index"`
	// ReleaseTime   unix.Time `json:"release_time"`
	// LockedPrice   string    `json:"locked_price"`
}

func (ReqOwnerUnlockTry) TagString() []string {
	return []string{
		"address", "지갑 주소",
		// "position_index", "락 인덱스",
		// "release_time", "락해제될 시간(utc 10자리)",
		// "locked_price", "락시켰던 수량",
		"", "",
	}
}

func (my ReqOwnerUnlockTry) String() string { return dbg.ToJsonString(my) }

///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////

func (my ReqOwnerUnlockTry) GetOwnerTask(decimals string) OwnerTask {
	nowAt := unix.Now()

	task := OwnerUnlockTask{
		Address: my.Address,

		// PositionIndex: my.PositionIndex,
		// Amount:        ebcm.TokenToWei(my.LockedPrice, decimals),
		// ReleaseTime:   my.ReleaseTime,

		State: NONE,
		Hash:  "",

		Decimals: decimals,
	}

	taskMap, _ := dbg.DecodeStruct[mongo.MAP](task)

	lockTask := OwnerTask{
		Key:  MakeOwnerTask_UnlockKey(),
		Kind: OwnerTaskKind_Unlock,

		State: NONE,

		Task: taskMap,

		CreateAt: nowAt,
		UpdateAt: nowAt,

		ZZ_CREATE_KST: nowAt.KST(),
		ZZ_UPDATE_KST: nowAt.KST(),
	}
	return lockTask
}

///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////

type OwnerUnlockTask struct {
	Address string `bson:"address" json:"address"`

	// PositionIndex int       `bson:"position_index" json:"position_index"`
	// Amount        string    `bson:"amount" json:"amount"`
	// ReleaseTime   unix.Time `bson:"release_time" json:"release_time"`

	State        STATE     `bson:"state" json:"state"` // 0:ready , 1:pending, 104:fail, 200:success
	StateMessage string    `bson:"state_message" json:"state_message"`
	TryData      TxTryData `bson:"try_data" json:"try_data"`

	Hash     string `bson:"hash" json:"hash"`
	TxFeeETH string `bson:"tx_fee_eth" json:"tx_fee_eth"`

	Decimals string `bson:"decimals" json:"decimals"`
}

func (my OwnerUnlockTask) String() string { return dbg.ToJsonString(my) }

func (OwnerUnlockTask) indexingC(c mongo.Collection, prefixDot string) {
	c.EnsureIndex(mongo.SingleIndex(prefixDot+"address", 1, false))

}

func (my OwnerUnlockTask) IsStateEnd() bool {
	return my.State == FAIL || my.State == SUCCESS
}

///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////

func (my OwnerUnlockTask) CbOwnerLockLog(key TaskKey) CbOwnerUnlockLog {

	timestamp := unix.Now()
	log := CbOwnerUnlockLog{
		Key: key,

		Address: my.Address,

		// PositionIndex: my.PositionIndex,
		// Price:         ebcm.WeiToToken(my.Amount, my.Decimals),
		// ReleaseTime:   my.ReleaseTime,

		Hash:     my.Hash,
		TxFeeETH: my.TxFeeETH,

		State:        my.State,
		StateMessage: my.StateMessage,

		Timestamp: timestamp,
		YMD:       timestamp.YMD(),
	}

	return log
}

////////////////////////////////////////////////////////////////////////

type CbOwnerUnlockLog struct {
	Key TaskKey `json:"key"`

	Address string `json:"address"`

	// PositionIndex int       `json:"position_index"`
	// Price         string    `json:"price"`
	// ReleaseTime   unix.Time `json:"release_time"`

	Hash     string `json:"hash"`
	TxFeeETH string `json:"tx_fee_eth"`

	State        STATE  `bson:"state" json:"state"` // 0:ready , 1:pending, 104:fail, 200:success
	StateMessage string `json:"state_message"`

	Timestamp unix.Time `json:"timestamp"`
	YMD       int       `json:"ymd"`
}

func (my CbOwnerUnlockLog) String() string { return dbg.ToJsonString(my) }
func (my CbOwnerUnlockLog) TagString() []string {
	return []string{
		"key", "키값 (OWNER65D8036D8FBA6D349A35BB62UL)",
		"address", "지갑 주소",

		"price", "락시킨 금액 (10.5 WNZ)",
		"release_time", "언락시간(10자리UTC)",
		"hash", "트랜잭션 해시값",
		"tx_fee_eth", "트랜잭션 수수료",

		"timestamp", "로그시간 (10자리 UTC)",
		"ymd", "로그날짜",
	}
}
func (my CbOwnerUnlockLog) GetKey() TaskKey         { return my.Key }
func (my CbOwnerUnlockLog) Kind() OwnerTaskKind     { return OwnerTaskKind_Unlock }
func (my CbOwnerUnlockLog) GetTimestamp() unix.Time { return my.Timestamp }

////////////////////////////////////////////////////////////////////////
