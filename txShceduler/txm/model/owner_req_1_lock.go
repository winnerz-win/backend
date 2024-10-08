package model

import (
	"jtools/cloud/ebcm"
	"jtools/dbg"
	"jtools/unix"
	"txscheduler/brix/tools/database/mongo"
)

// func (ReqOwnerLockItem) TagString() []string {
// 	return []string{
// 		"kind", "lock , unlock (락 , 언락)",
// 		"price", "kind가 lock일경우 락시킬 수량",
// 		"position_index", "kind가 unlock일경우 포지션 인덱스 값",
// 		"release_time", "kind가 lock일경우 락해제 시간(utc 10자리) , kind가 unlock일경우 해당 포지션의 락해제 시간",
// 	}
// }

type ReqOwnerLockTry struct {
	Address string `json:"address"`

	Price       string    `json:"price"`
	ReleaseTime unix.Time `json:"release_time"`
}

func (ReqOwnerLockTry) TagString() []string {
	return []string{
		"address", "지갑 주소",
		"price", "락시킬 수량",
		"release_time", "락해제 시간(utc 10자리)",
		"", "",
	}
}

func (my ReqOwnerLockTry) String() string { return dbg.ToJsonString(my) }

///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////

func (my ReqOwnerLockTry) GetOwnerTask(decimals string) OwnerTask {
	nowAt := unix.Now()

	task := OwnerLockTask{
		Address: my.Address,

		Amount:      ebcm.TokenToWei(my.Price, decimals),
		ReleaseTime: my.ReleaseTime,

		State: NONE,
		Hash:  "",

		Decimals: decimals,
	}

	taskMap, _ := dbg.DecodeStruct[mongo.MAP](task)

	owner_task := OwnerTask{
		Key:  MakeOwnerTask_LockKey(),
		Kind: OwnerTaskKind_Lock,

		State: NONE,

		Task: taskMap,

		CreateAt: nowAt,
		UpdateAt: nowAt,

		ZZ_CREATE_KST: nowAt.KST(),
		ZZ_UPDATE_KST: nowAt.KST(),
	}
	return owner_task
}

///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////

type OwnerLockTask struct {
	Address string `bson:"address" json:"address"`

	Amount      string    `bson:"amount" json:"amount"`
	ReleaseTime unix.Time `bson:"release_time" json:"release_time"`

	State        STATE     `bson:"state" json:"state"` // 0:ready , 1:pending, 104:fail, 200:success
	StateMessage string    `bson:"state_message" json:"state_message"`
	TryData      TxTryData `bson:"try_data" json:"try_data"`

	Hash     string `bson:"hash" json:"hash"`
	TxFeeETH string `bson:"tx_fee_eth" json:"tx_fee_eth"`

	Decimals string `bson:"decimals" json:"decimals"`
}

func (my OwnerLockTask) String() string { return dbg.ToJsonString(my) }

func (OwnerLockTask) indexingC(c mongo.Collection, prefixDot string) {
	c.EnsureIndex(mongo.SingleIndex(prefixDot+"address", 1, false))

}

func (my OwnerLockTask) IsStateEnd() bool {
	return my.State == FAIL || my.State == SUCCESS
}

///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////

func (my OwnerLockTask) CbOwnerLockLog(key TaskKey) CbOwnerLockLog {

	timestamp := unix.Now()
	log := CbOwnerLockLog{
		Key: key,

		Address: my.Address,

		Price:       ebcm.WeiToToken(my.Amount, my.Decimals),
		ReleaseTime: my.ReleaseTime,

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

type CbOwnerLockLog struct {
	Key TaskKey `json:"key"`

	Address string `json:"address"`

	Price       string    `json:"price"`
	ReleaseTime unix.Time `json:"release_time"`

	Hash     string `json:"hash"`
	TxFeeETH string `json:"tx_fee_eth"`

	State        STATE  `bson:"state" json:"state"` // 0:ready , 1:pending, 104:fail, 200:success
	StateMessage string `json:"state_message"`

	Timestamp unix.Time `json:"timestamp"`
	YMD       int       `json:"ymd"`
}

func (my CbOwnerLockLog) String() string { return dbg.ToJsonString(my) }
func (my CbOwnerLockLog) TagString() []string {
	return []string{
		"key", "키값 (OWNER65D8036D8FBA6D349A35BB62LK)",
		"address", "지갑 주소",

		"price", "락시킨 금액 (10.5 WNZ)",
		"release_time", "언락시간(10자리UTC)",
		"hash", "트랜잭션 해시값",
		"tx_fee_eth", "트랜잭션 수수료",

		"timestamp", "로그시간 (10자리 UTC)",
		"ymd", "로그날짜",
	}
}
func (my CbOwnerLockLog) GetKey() TaskKey         { return my.Key }
func (my CbOwnerLockLog) Kind() OwnerTaskKind     { return OwnerTaskKind_Lock }
func (my CbOwnerLockLog) GetTimestamp() unix.Time { return my.Timestamp }
