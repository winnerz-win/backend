package model

import (
	"jtools/cloud/ebcm"
	"jtools/unix"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
)

type ReqOwnerRelockTry struct {
	Address string `json:"address"`

	Price       string    `json:"price"`
	ReleaseTime unix.Time `json:"release_time"`
}

func (ReqOwnerRelockTry) TagString() []string {
	return []string{
		"address", "지갑 주소",
		"price", "락시킬 수량",
		"release_time", "락해제 시간(utc 10자리)",
		"", "",
	}
}

func (my ReqOwnerRelockTry) String() string { return dbg.ToJsonString(my) }

///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////

func (my ReqOwnerRelockTry) GetOwnerTask(decimals string) OwnerTask {
	nowAt := unix.Now()

	task := OwnerRelockTask{
		Address: my.Address,

		Amount:      ebcm.TokenToWei(my.Price, decimals),
		ReleaseTime: my.ReleaseTime,
		Decimals:    decimals,
	}

	taskMap, _ := dbg.DecodeStruct[mongo.MAP](task)

	owner_task := OwnerTask{
		Key:  MakeOwnerTask_RelockKey(),
		Kind: OwnerTaskKind_Relock,

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

type OwnerRelockTask struct {
	From    string `bson:"from" json:"from"`
	Address string `bson:"address" json:"address"`

	State        STATE  `bson:"state" json:"state"` // 0:ready , 1:pending, 104:fail, 200:success
	StateMessage string `bson:"state_message" json:"state_message"`

	UnlockHash     string    `bson:"unlock_hash" json:"unlock_hash"`
	UnlockTxFeeETH string    `bson:"unlock_tx_fee_eth" json:"unlock_tx_fee_eth"`
	UnlockTryData  TxTryData `bson:"unlock_try_data" json:"unlock_try_data"`

	LockHash     string    `bson:"lock_hash" json:"lock_hash"`
	LockTxFeeETH string    `bson:"lock_tx_fee_eth" json:"lock_tx_fee_eth"`
	LockTryData  TxTryData `bson:"lock_try_data" json:"lock_try_data"`

	Amount      string    `bson:"amount" json:"amount"`
	ReleaseTime unix.Time `bson:"release_time" json:"release_time"`
	Decimals    string    `bson:"decimals" json:"decimals"`
}

func (my OwnerRelockTask) String() string { return dbg.ToJSONString(my) }

func (OwnerRelockTask) indexingC(c mongo.Collection, prefixDot string) {
	c.EnsureIndex(mongo.SingleIndex(prefixDot+"address", 1, false))

}

func (my OwnerRelockTask) IsStateEnd() bool {
	return my.State == FAIL || my.State == SUCCESS
}

///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////

func (my OwnerRelockTask) CbOwnerRelockLog(key TaskKey) CbOwnerRelockLog {
	timestamp := unix.Now()
	log := CbOwnerRelockLog{
		Key: key,

		Address: my.Address,

		Price:       ebcm.WeiToToken(my.Amount, my.Decimals),
		ReleaseTime: my.ReleaseTime,

		UnlockHash:     my.UnlockHash,
		UnlockTxFeeETH: my.UnlockTxFeeETH,

		LockHash:     my.LockHash,
		LockTxFeeETH: my.LockTxFeeETH,

		State:        my.State,
		StateMessage: my.StateMessage,

		Timestamp: timestamp,
		YMD:       timestamp.YMD(),
	}

	return log
}

////////////////////////////////////////////////////////////////////////

type CbOwnerRelockLog struct {
	Key TaskKey `json:"key"`

	Address string `json:"address"`

	Price       string    `json:"price"`
	ReleaseTime unix.Time `json:"release_time"`

	UnlockHash     string `json:"unlock_hash"`
	UnlockTxFeeETH string `json:"unlock_tx_fee_eth"`

	LockHash     string `json:"lock_hash"`
	LockTxFeeETH string `json:"lock_tx_fee_eth"`

	State        STATE  `bson:"state" json:"state"` // 0:ready , 1:pending, 104:fail, 200:success
	StateMessage string `json:"state_message"`

	Timestamp unix.Time `json:"timestamp"`
	YMD       int       `json:"ymd"`
}

func (my CbOwnerRelockLog) String() string { return dbg.ToJsonString(my) }
func (my CbOwnerRelockLog) TagString() []string {
	return []string{
		"key", "키값 (OWNER65D8036D8FBA6D349A35BB62LK)",
		"address", "지갑 주소",

		"price", "락시킨 금액 (10.5 WNZ)",
		"release_time", "언락시간(10자리UTC)",

		"unlock_hash", "언락 트랜잭션 해시값",
		"unlock_tx_fee_eth", "언락 트랜잭션 수수료",

		"lock_hash", "락 트랜잭션 해시값",
		"lock_tx_fee_eth", "락 트랜잭션 수수료",

		"timestamp", "로그시간 (10자리 UTC)",
		"ymd", "로그날짜",
	}
}
func (my CbOwnerRelockLog) GetKey() TaskKey         { return my.Key }
func (my CbOwnerRelockLog) Kind() OwnerTaskKind     { return OwnerTaskKind_Relock }
func (my CbOwnerRelockLog) GetTimestamp() unix.Time { return my.Timestamp }
