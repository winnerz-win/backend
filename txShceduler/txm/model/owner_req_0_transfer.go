package model

import (
	"jtools/cloud/ebcm"
	"jtools/unix"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
)

type ReqOwnerTransferItem struct {
	Price       string    `json:"price"`
	ReleaseTime unix.Time `json:"release_time"`
}

func (ReqOwnerTransferItem) TagString() []string {
	return []string{
		"price", "전송할 수량(가격)",
		"release_time", "언락될 시간(UTC 10자리 숫자), 0값이면 해당 수량은 lock_transfer가 아닌 일반 transfer로 전송됨.",
	}
}

type ReqOwnerTransferTry struct {
	User      string                 `json:"user"`
	Recipient string                 `json:"recipient"`
	Transfers []ReqOwnerTransferItem `json:"transfers"`
}

func (ReqOwnerTransferTry) TagString() []string {
	return []string{
		"user", "회원의 가상 계좌 주소",
		"recipient", "외부 주소 (오너가 실제로 토큰을 전달할 주소)",
		"transfers", "전송수량 및 언락시간 배열",
	}
}

func (my ReqOwnerTransferTry) String() string { return dbg.ToJSONString(my) }

///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////

func (my ReqOwnerTransferTry) GetOwnerTask(decimals string) OwnerTask {
	nowAt := unix.Now()

	task := OwnerTransferTask{
		User:      my.User,
		Recipient: my.Recipient,

		Seq:   0,
		Datas: make(OwnerTransferTryDataList, len(my.Transfers)),

		Decimals: decimals,
	}

	for i, v := range my.Transfers {
		data := OwnerTransferTryData{
			Amount:      ebcm.TokenToWei(v.Price, decimals),
			ReleaseTime: v.ReleaseTime,
			Hash:        "",
			State:       NONE,
			UpdateAt:    unix.ZERO,
		}
		task.Datas[i] = data
	} //for

	taskMap, _ := dbg.DecodeStruct[mongo.MAP](task)

	lockTask := OwnerTask{
		Key:  MakeOwnerTask_TransferKey(),
		Kind: OwnerTaskKind_Transfer,

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

type OwnerTransferTryData struct {
	Amount      string    `bson:"amount" json:"amount"`
	ReleaseTime unix.Time `bson:"release_time" json:"release_time"` //0이면 일반 transfer

	Hash     string `bson:"hash" json:"hash"`
	State    STATE  `bson:"state" json:"state"` // 0:ready , 1:pending, 104:fail, 200:success
	TxFeeETH string `bson:"tx_fee_eth" json:"tx_fee_eth"`

	UpdateAt unix.Time `bson:"update_at" json:"update_at"`
}
type OwnerTransferTryDataList []OwnerTransferTryData

func (my OwnerTransferTryData) String() string     { return dbg.ToJsonString(my) }
func (my OwnerTransferTryDataList) String() string { return dbg.ToJsonString(my) }

func (my *OwnerTransferTryDataList) IndexPointer(idx int) *OwnerTransferTryData {
	for i := range *my {
		if i == idx {
			return &(*my)[i]
		}
	} //for
	return nil
}

type OwnerTransferTask struct {
	User      string `bson:"user" json:"user"`
	Recipient string `bson:"recipient" json:"recipient"`

	Seq      int                      `bson:"seq" json:"seq"`
	Datas    OwnerTransferTryDataList `bson:"datas" json:"datas"`
	Decimals string                   `bson:"decimals" json:"decimals"`
}

func (my OwnerTransferTask) String() string { return dbg.ToJsonString(my) }

func (my *OwnerTransferTask) CurrentData() *OwnerTransferTryData {
	return my.Datas.IndexPointer(my.Seq)
}

func (OwnerTransferTask) indexingC(c mongo.Collection, prefixDot string) {

	c.EnsureIndex(mongo.SingleIndex(prefixDot+"user", 1, false))
	c.EnsureIndex(mongo.SingleIndex(prefixDot+"recipient", 1, false))

}

func (my *OwnerTransferTask) CheckStateEnd() bool {
	for _, v := range my.Datas {
		if v.State != SUCCESS {
			return false
		}
	} //for
	my.Seq = len(my.Datas)
	return true
}

func (my OwnerTransferTask) IsStateEnd() bool { return my.CheckStateEnd() }

///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////

func (my OwnerTransferTask) CbOwnerTransferLog(key TaskKey) CbOwnerTransferLog {

	timestamp := unix.Now()
	log := CbOwnerTransferLog{
		Key: key,

		User:      my.User,
		Recipient: my.Recipient,

		Transfers: make([]CbOwnerTransferTxLog, len(my.Datas)),

		Timestamp: timestamp,
		YMD:       timestamp.YMD(),
	}

	for i, v := range my.Datas {
		log.Transfers[i] = CbOwnerTransferTxLog{
			Price:       ebcm.WeiToToken(v.Amount, my.Decimals),
			ReleaseTime: v.ReleaseTime,
			Hash:        v.Hash,
			TxFeeETH:    v.TxFeeETH,
		}
	} //for

	return log
}

////////////////////////////////////////////////////////////////////////

type CbOwnerTransferTxLog struct {
	Price       string    `json:"price"`
	ReleaseTime unix.Time `json:"release_time"`
	Hash        string    `json:"hash"`
	TxFeeETH    string    `json:"tx_fee_eth"`
}

func (my CbOwnerTransferTxLog) String() string { return dbg.ToJsonString(my) }
func (my CbOwnerTransferTxLog) TagString() []string {
	return []string{
		"price", "전송금액 (10.5 WNZ)",
		"release_time", "언락시간(10자리UTC) , 0이면 일반 Transfer로 전송",
		"hash", "트랜잭션 해시값",
		"tx_fee_eth", "트랜잭션 수수료",
	}
}

type CbOwnerTransferLog struct {
	Key       TaskKey                `json:"key"`
	User      string                 `json:"user"`
	Recipient string                 `json:"recipient"`
	Transfers []CbOwnerTransferTxLog `json:"transfers"`
	Timestamp unix.Time              `json:"timestamp"`
	YMD       int                    `json:"ymd"`
}

func (my CbOwnerTransferLog) String() string { return dbg.ToJsonString(my) }
func (my CbOwnerTransferLog) TagString() []string {
	return []string{
		"key", "키값 (LOCK65D8036D8FBA6D349A35BB62TK)",
		"user", "가상 계좌 (회원주소)",
		"recipient", "외부주소 (실제 WNZ를 받을)",
		"transfers", "락업 트랜잭션 데이터들",
		"timestamp", "로그시간 (10자리 UTC)",
		"ymd", "로그날짜",
	}
}
func (my CbOwnerTransferLog) GetKey() TaskKey         { return my.Key }
func (my CbOwnerTransferLog) Kind() OwnerTaskKind     { return OwnerTaskKind_Transfer }
func (my CbOwnerTransferLog) GetTimestamp() unix.Time { return my.Timestamp }

/*
	Params : {
		"key" : "LOCK65D8036D8FBA6D349A35BB62TK",
		"user" : "가상 계좌 주소(회원 주소)",
		"recipient" : "외부주소(실제 WNZ를 전송할)",
		"transfers" : [	//락업 트랜잭션 전송 데이터 (총 15WNZ를 recipient에게 전송, 10.5 WNZ는 락업 전송, 4.5 WNZ는 일반 전송)
			{
				"price" : "10.5",
				"release_time" : 1234567890, // 10자리 UTC (언락 시간)
				"hash" : "0x..."			// 트랜잭션 해시값.
			},
			{
				"price" : "4.5",
				"release_time" : 0, 		// 0이면 일반 transfer로 전송
				"hash" : "0x..."			// 트랜잭션 해시값.
			}
		],
		"timestamp" : 1234567890, //10자리 UTC
		"ymd" : 20240101
	}
*/
