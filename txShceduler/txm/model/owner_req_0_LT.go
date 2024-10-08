package model

import (
	"jtools/cloud/ebcm"
	"jtools/unix"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/txm/inf"
)

type ReqOwner_LT_Try struct {
	Recipient   string    `bson:"recipient" json:"recipient"`
	Price       string    `bson:"price" json:"price"`
	ReleaseTime unix.Time `bson:"release_time" json:"release_time"`
}

func (ReqOwner_LT_Try) TagString() []string {
	return []string{
		"recipient", "외부 주소 (마스터가 실제로 토큰을 전달할 주소)",
		"price", "전송할 수량(가격)",
		"release_time", "언락될 시간(UTC 10자리 숫자), 미래의 시간이어야 함. (요청시 서버 시간과 비교함)",
	}
}

func (my ReqOwner_LT_Try) String() string { return dbg.ToJSONString(my) }

///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////

func (my ReqOwner_LT_Try) GetOwnerTask(decimals string) OwnerTask {
	nowAt := unix.Now()

	task := Owner_LT_Task{
		Recipient:   my.Recipient,
		Amount:      ebcm.TokenToWei(my.Price, decimals),
		ReleaseTime: my.ReleaseTime,

		TransferInfo: LT_State{
			From:         inf.Master().Address,
			FailHashList: []string{},
		},
		LockInfo: LT_State{
			From:         inf.Owner().Address,
			FailHashList: []string{},
		},

		Decimals: decimals,
	}
	taskMap, _ := dbg.DecodeStruct[mongo.MAP](task)

	owner_task := OwnerTask{
		Key:  MakeOwnerTask_LT_Key(),
		Kind: OwnerTaskKind_LT,

		State: MASTER_NONE,

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

type LT_State struct {
	From string `bson:"from" json:"from"` // master, owner

	State        STATE     `bson:"state" json:"state"` // 0:ready , 1:pending, 104:fail, 200:success
	StateMessage string    `bson:"state_message" json:"state_message"`
	TryData      TxTryData `bson:"try_data" json:"try_data"`

	FailHashList []string `bson:"fail_hash_list" json:"fail_hash_list"`
	Hash         string   `bson:"hash" json:"hash"`
	TxFeeETH     string   `bson:"tx_fee_eth" json:"tx_fee_eth"`
}

func (LT_State) TagString() []string {
	return []string{
		"from", "owner 또는 master 주소",
		"state", "0:ready , 1:pending, 104:fail, 200:success",
		"state_message", "상태 메시지",
		"hash", "트랜잭션 해시값",
		"tx_fee_eth", "트랜잭션 수수료",
	}
}

type Owner_LT_Task struct {
	Recipient   string    `bson:"recipient" json:"recipient"`
	Amount      string    `bson:"amount" json:"amount"`
	ReleaseTime unix.Time `bson:"release_time" json:"release_time"`

	TransferInfo LT_State `bson:"transfer_info" json:"transfer_info"`
	LockInfo     LT_State `bson:"lock_info" json:"lock_info"`

	Decimals string `bson:"decimals" json:"decimals"`
}

func (my Owner_LT_Task) String() string { return dbg.ToJsonString(my) }

func (Owner_LT_Task) indexingC(c mongo.Collection, prefixDot string) {
	c.EnsureIndex(mongo.SingleIndex(prefixDot+"recipient", 1, false))

}

func (my Owner_LT_Task) IsEndAll() bool {
	a := my.TransferInfo.State == FAIL || my.TransferInfo.State == SUCCESS
	b := my.LockInfo.State == FAIL || my.LockInfo.State == SUCCESS
	return a && b
}

func (my Owner_LT_Task) IsStateEnd() bool {
	//return my.State == FAIL || my.State == SUCCESS
	return false
}

///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////

func (my Owner_LT_Task) CbOwner_LT_Log(key TaskKey) CbOwner_LT_Log {
	timestamp := unix.Now()
	log := CbOwner_LT_Log{
		Key: key,

		Recipient:   my.Recipient,
		Price:       ebcm.WeiToToken(my.Amount, my.Decimals),
		ReleaseTime: my.ReleaseTime,

		TransferInfo: my.TransferInfo,
		LockInfo:     my.LockInfo,

		Timestamp: timestamp,
		YMD:       timestamp.YMD(),
	}
	return log
}

////////////////////////////////////////////////////////////////////////

type CbOwner_LT_Log struct {
	Key TaskKey `json:"key"`

	Recipient   string    `bson:"recipient" json:"recipient"`
	Price       string    `json:"price"`
	ReleaseTime unix.Time `json:"release_time"`

	TransferInfo LT_State `bson:"transfer_info" json:"transfer_info"`
	LockInfo     LT_State `bson:"lock_info" json:"lock_info"`

	Timestamp unix.Time `json:"timestamp"`
	YMD       int       `json:"ymd"`
}

func (my CbOwner_LT_Log) String() string { return dbg.ToJsonString(my) }

func (my CbOwner_LT_Log) TagString() []string {
	return []string{
		"key", "키값 (OWNER65D8036D8FBA6D349A35BB62LT)",
		"recipient", "지갑 주소",

		"price", "락시킨 금액 (10.5 WNZ)",
		"release_time", "언락시간(10자리UTC)",

		"transfer_info", "trasnfer정보 (Hash, 수수료, 상태값)",
		"lock_info", "lock정보 (Hash, 수수료, 상태값)",

		"timestamp", "로그시간 (10자리 UTC)",
		"ymd", "로그날짜",
	}
}

func (my CbOwner_LT_Log) GetKey() TaskKey         { return my.Key }
func (my CbOwner_LT_Log) Kind() OwnerTaskKind     { return OwnerTaskKind_LT }
func (my CbOwner_LT_Log) GetTimestamp() unix.Time { return my.Timestamp }
