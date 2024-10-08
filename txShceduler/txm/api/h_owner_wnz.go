package api

import (
	"jtools/cloud/ebcm"
	"jtools/cloud/jeth/ecs"
	"jtools/jmath"
	"jtools/mms"
	"jtools/unix"
	"net/http"
	"sync"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/database/mongo/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/jnet/doc"
	"txscheduler/txm/ack"
	"txscheduler/txm/cloud"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func applyOnwerHandles() {
	hOwner_ymd()

	hOwner_Request_LT()
	//hOwner_Request_Transfer()
	hOwner_State_Key()
	hOwner_State_Address()

	hOwner_Request_Lock()
	hOwner_Request_Unlock()
	hOwner_Request_ReLock()

}

var (
	//"[WNZ-LOCKUP] "
	OWNER_TAG = "[WNZ-LOCKUP] "
)

func hOwner_ymd() {

	type CDATA struct {
		Year  int `json:"year"`
		Month int `json:"month"`
		Day   int `json:"day"`
		Hour  int `json:"hour"`
		Min   int `json:"min"`
	}
	type Date struct {
		KST       string    `json:"kst"`
		UTC       string    `json:"utc"`
		Timestamp unix.Time `json:"timestamp"`
	}
	type RESULT struct {
		Today  Date `json:"today"`
		Target Date `json:"target"`

		Du string `json:"du"`
	}
	method := chttp.POST
	url := model.V1Owner + "/ymd"
	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := chttp.BindingStruct[CDATA](req)

			target_unix := unix.FromMMS(
				mms.FromYMDHM(
					cdata.Year,
					cdata.Month,
					cdata.Day,
					cdata.Hour,
					cdata.Min,
				),
			)

			target := Date{
				KST:       target_unix.KST(),
				UTC:       target_unix.String(),
				Timestamp: target_unix,
			}

			nowAt := unix.Now()
			today := Date{
				KST:       nowAt.KST(),
				UTC:       nowAt.String(),
				Timestamp: nowAt,
			}

			du := target.Timestamp.Sub(today.Timestamp)

			result := RESULT{
				Today:  today,
				Target: target,

				Du: du.String(),
			}

			chttp.OK(w, result)
		},
	)
}

// //////////////////////////////////////////////////////////////////////
func check_lock_ban(address string) chttp.CError {
	is_locked, ack_err := cloud.CheckLockState(address)
	if ack_err != nil {
		return ack_err
	}
	if is_locked {
		return ack.OWNER_ALREADY_LOCKED
	}
	return nil
}
func check_already_job(db mongo.DATABASE, address string) bool {

	selector := mongo.Bson{"$or": []mongo.Bson{
		{
			"task.recipient": address,
			"state":          mongo.Bson{"$lt": model.SUCCESS},
		},
		{
			"task.address": address,
			"state":        mongo.Bson{"$lt": model.SUCCESS},
		},
	}}
	cnt, err := db.C(inf.OwnerTask).Find(selector).Count()
	if err != nil {
		return true
	}

	is_already_job := cnt > 0

	return is_already_job
}

var owner_mu = sync.Mutex{}

func OWNER_DB(w http.ResponseWriter, address string, f func(db mongo.DATABASE)) {
	owner_mu.Lock()
	defer owner_mu.Unlock()

	model.DB(func(db mongo.DATABASE) {
		if check_already_job(db, address) {
			chttp.Fail(w, ack.OWNER_REQ_JOB_ONCE)
			return
		}
		f(db)
	})
}

func hOwner_Request_LT() {
	type RESULT struct {
		Key model.TaskKey `json:"key"`
	}

	method := chttp.POST
	url := model.V1Owner + "/request/lock_transfer"

	Doc().Comment(OWNER_TAG+"Transfer(Master) & Lock(OWner) 요청").
		Method(method).URL(url).
		JParam(model.ReqOwner_LT_Try{}, model.ReqOwner_LT_Try{}.TagString()...).
		JResultOK(
			RESULT{},
			"key", "결과 조회/콜백 용 Key값",
		).
		JAckError(ack.BadParam).
		JAckError(ack.OWNER_ALREADY_LOCKED).
		JAckError(ack.OWNER_REQ_JOB_ONCE).
		JAckError(ack.InvalidAddress, "recipent 주소형식 오류").
		JAckError(ack.OWNER_Transfer_Price, "price 필드가 0이하").
		JAckError(ack.OWNER_Transfer_ReleaseTime, "release_time이 현재 시간 이하").
		JAckError(ack.OWNER_Transfer_MemberRecipient, "recipient가 회원 가상계좌 주소, 오너/마스터/가스비충전지갑 일경우.").
		Etc("", `
			<설명>
			recipient 주소가 다른 가상계좌 주소이거나 마스터/가스비충전지갑/오너 이면 안됩니다. ( 순수 외부 지갑 주소 여야만 합니다.)
			price 는 0 이하면 에러.
			release_time 은 요청시 서버시간 기준보다 이하이면 에러.

			<Request-Sample>
			{
				"recipient" : "0x8CE5bb2013887eD586e6a87211aa126453368b7A",
				"price" : "0.5",
				"release_time" : 1709168400				
			}

			: Master지갑 주소에서 recipient에게 0.5 WNZ를 전송한후, 해당 전송이 성공 하게 되면 
			 Owner가 recipient의 0.5 WNZ를 release_time기간 락을 건다.
		`).
		Apply(doc.Blue)

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {

			cdata := chttp.BindingStruct[model.ReqOwner_LT_Try](req)
			_ = cdata

			if !ebcm.IsAddressP(&cdata.Recipient) {
				chttp.Fail(w, ack.InvalidAddress)
				return
			}

			if jmath.CMP(cdata.Price, 0) <= 0 {
				chttp.Fail(w, ack.OWNER_Transfer_Price, "price is under zero.")
				return
			}

			nowAt := unix.Now()
			if cdata.ReleaseTime <= nowAt {
				chttp.Fail(w, ack.OWNER_Transfer_ReleaseTime, dbg.Cat("release_time less current server time."))
				return
			}

			OWNER_DB(w, cdata.Recipient, func(db mongo.DATABASE) {
				if cdata.Recipient == inf.Master().Address {
					chttp.Fail(w, ack.OWNER_Transfer_MemberRecipient, "recipient is master")
					return
				}
				if cdata.Recipient == inf.Charger().Address {
					chttp.Fail(w, ack.OWNER_Transfer_MemberRecipient, "recipient is charger")
					return
				}
				if cdata.Recipient == inf.Owner().Address {
					chttp.Fail(w, ack.OWNER_Transfer_MemberRecipient, "recipient is owner")
					return
				}

				other := model.Member{}
				db.C(inf.COLMember).Find(mongo.Bson{"address": cdata.Recipient}).One(&other)
				if other.Valid() {
					chttp.Fail(w, ack.OWNER_Transfer_MemberRecipient, "recipient is member")
					return
				}

				if ack_err := check_lock_ban(cdata.Recipient); ack_err != nil {
					chttp.Fail(w, ack_err)
					return
				}

				if check_already_job(db, cdata.Recipient) { //lock_transfer
					chttp.Fail(w, ack.OWNER_REQ_JOB_ONCE)
					return
				}

				decimals := inf.FirstERC20().Decimal
				task := cdata.GetOwnerTask(decimals)

				if err := db.C(inf.OwnerTask).Insert(task); err != nil {
					chttp.Fail(w, ack.DBJob)
					return
				}

				chttp.OK(w,
					RESULT{
						Key: task.Key,
					},
				)
			})

		},
	)
}

////////////////////////////////////////////////////////////////////////
//
//

// func hOwner_Request_Transfer() {

// 	type RESULT struct {
// 		Key model.TaskKey `json:"key"`
// 	}

// 	method := chttp.POST
// 	url := model.V1Owner + "/request/transfer"

// 	Doc().Comment(OWNER_TAG+"LockTransfer 요청").
// 		Method(method).URL(url).
// 		JParam(model.ReqOwnerTransferTry{}, model.ReqOwnerTransferTry{}.TagString()...).
// 		JResultOK(
// 			RESULT{},
// 			"key", "결과 조회/콜백 용 Key값",
// 		).
// 		ETCVAL(model.ReqOwnerTransferItem{}, model.ReqOwnerTransferItem{}.TagString()...).
// 		JAckError(ack.BadParam).
// 		JAckError(ack.InvalidAddress, "user / recipent 주소형식 오류").
// 		JAckError(ack.OWNER_Transfer_SameAddress, "user == recipient").
// 		JAckError(ack.OWNER_Transfer_EmptyDatas, "transfers 가 빈 배열 일경우.").
// 		JAckError(ack.OWNER_Transfer_Price, "price 필드가 0이하").
// 		JAckError(ack.OWNER_Transfer_ReleaseTime, "release_time이 현재 시간 이하").
// 		JAckError(ack.NotFoundAddress, "user가 회원주소가 아닐경우.").
// 		JAckError(ack.OWNER_Transfer_MemberRecipient, "recipient가 회원 가상 주소 일경우.").
// 		Etc("", `
// 			<설명>
// 			user 주소는 반드시 가상계좌(블록체인 서버에서 발급한) 주소이어야 합니다.
// 			user != recipient : 두 주소는 반드시 달라야 합니다.
// 			recipient 주소가 다른 가상계좌의 주소이면 안됩니다. ( 순수 외부 지갑 주소 여야만 합니다. )

// 			transfers 의 개수가 0이면 에러.
// 			transfers[i].price 의 수량이 0이하 이면 에러.
// 			transfers[i].release_time 의 시간 값이 0 미만 이면 에러.

// 			<Request-Sample>
// 			{
// 				"user" : "0xf811b879b9f4f24b411a92ebd10dfb7e79c4a200",
// 				"recipient" : "0x8CE5bb2013887eD586e6a87211aa126453368b7A",
// 				"transfers" : [
// 					{
// 						"price": "1",
// 						"release_time" : 1709168400
// 					},
// 					{
// 						"price" : "5",
// 						"release_time" : 0	//0이면 일반 Transfer로 전송
// 					},
// 					{
// 						"price": "0.6",
// 						"release_time" : 1709172000
// 					}
// 					,
// 					{
// 						"price": "0.4",
// 						"release_time" : 1709175600
// 					}
// 				]
// 			}

// 			: owner지갑 주소에서 user(회원지갑주소)가 지정한 recipient에게
// 				총 7개(1 + 5 + 0.6 + 0.4)의 WNZ를 전송합니다.
// 				여기서 2개는 Lock전송이고 5개는 일반 전송 입니다.
// 		`).
// 		Apply(doc.Blue)

// 	handle.Append(
// 		method, url,
// 		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {

// 			cdata := chttp.BindingStruct[model.ReqOwnerTransferTry](req)
// 			_ = cdata

// 			if !ebcm.IsAddressP(&cdata.User) {
// 				chttp.Fail(w, ack.InvalidAddress)
// 				return
// 			}
// 			if !ebcm.IsAddressP(&cdata.Recipient) {
// 				chttp.Fail(w, ack.InvalidAddress)
// 				return
// 			}

// 			if cdata.User == cdata.Recipient {
// 				chttp.Fail(w, ack.OWNER_Transfer_SameAddress, "user is same recipient")
// 				return
// 			}

// 			if cdata.Recipient == inf.Master().Address {
// 				chttp.Fail(w, ack.OWNER_Transfer_MemberRecipient, "recipient is master")
// 				return
// 			}
// 			if cdata.Recipient == inf.Charger().Address {
// 				chttp.Fail(w, ack.OWNER_Transfer_MemberRecipient, "recipient is charger")
// 				return
// 			}
// 			if cdata.Recipient == inf.Owner().Address {
// 				chttp.Fail(w, ack.OWNER_Transfer_MemberRecipient, "recipient is owner")
// 				return
// 			}

// 			if len(cdata.Transfers) == 0 {
// 				chttp.Fail(w, ack.OWNER_Transfer_EmptyDatas, "transfers is empty")
// 				return
// 			}

// 			nowAt := unix.Now()
// 			for i, v := range cdata.Transfers {
// 				if jmath.CMP(v.Price, 0) <= 0 {
// 					chttp.Fail(w, ack.OWNER_Transfer_Price, dbg.Cat("transfers[", i, "].price under zero"))
// 					return
// 				}
// 				if v.ReleaseTime <= nowAt {
// 					chttp.Fail(w, ack.OWNER_Transfer_ReleaseTime, dbg.Cat("transfers[", i, "].release_time less current server time."))
// 					return
// 				}
// 			} //for

// 			model.DB(func(db mongo.DATABASE) {
// 				member := model.Member{}
// 				db.C(inf.COLMember).Find(mongo.Bson{"address": cdata.User}).One(&member)
// 				if !member.Valid() {
// 					chttp.Fail(w, ack.NotFoundAddress, "not found member")
// 					return
// 				}
// 				other := model.Member{}
// 				db.C(inf.COLMember).Find(mongo.Bson{"address": cdata.Recipient}).One(&other)
// 				if other.Valid() {
// 					chttp.Fail(w, ack.OWNER_Transfer_MemberRecipient, "recipient is member")
// 					return
// 				}

// 				decimals := inf.FirstERC20().Decimal
// 				task := cdata.GetOwnerTask(decimals)

// 				if err := db.C(inf.OwnerTask).Insert(task); err != nil {
// 					chttp.Fail(w, ack.DBJob)
// 					return
// 				}

// 				chttp.OK(w,
// 					RESULT{
// 						Key: task.Key,
// 					},
// 				)

// 			})

// 		},
// 	)
// }

////////////////////////////////////////////////////////////////////////
//
//

func hOwner_State_Key() {

	type RESULT struct {
		State model.STATE `json:"state"`
		Log   interface{} `json:"log,omitempty"`
	}

	method := chttp.GET
	url := model.V1Owner + "/state/key/:args"

	Doc().Comment(OWNER_TAG+"Key값으로 현재 lock_transfer/lock/unlock/relock 결과 조회").
		Method(method).URLS(url, ":args", "... api로 받은 key값").
		JResultOK(
			RESULT{},
			"state", "상태 ( 0:대기중, 1:현재작업 진행중, 104:작업실패, 200:작업완료 )",
			"log", "state가 104/200일경우 lock transfer 콜백결과 log",
		).
		ETCVAL(model.CbOwnerTransferLog{}, model.CbOwnerTransferLog{}.TagString()...).
		ETCVAL(model.CbOwnerTransferTxLog{}, model.CbOwnerTransferTxLog{}.TagString()...).
		JAckError(ack.OWNER_NotFoundData).
		Etc("", `
			<Response-Sample>

			<트랜잭션 진행중 - 공통>
			{
				"success": true,
				"data": {
					"state": 1
				}
			}

			<트랜잭션 결과 - lock_transfer >
			{
				"success": true,
				"data": {
					"state": 200,
					"log": {
						"key": "OWNER65F2895065A75F54A65DB97ELT",
						"recipient": "0x8ce5bb2013887ed586e6a87211aa126453368b7a",
						"price": "10.5",
						"release_time": 1710892800,
						"transfer_info": {
							"from": "0x7d44dbb4a0fa180774611412edeacdd2e7a13ec8",
							"state": 200,
							"state_message": "",
							"fail_hash_list": [],
							"hash": "0xbbcfab5baaea91e5298289cc2b87ecaf1573f3b3ddfb593a17eedb1d633cb18a",
							"tx_fee_eth": "0.00017333168569541"
						},
						"lock_info": {
							"from": "0xff00bbd12f9ae64fe4dfa9f3dddfdc3d9ecc2378",
							"state": 200,
							"state_message": "",
							"fail_hash_list": [],
							"hash": "0x37ec9ce1915a1b4bf35cd3c38cffcbff4ba55b73888a9c56309159afb340e20d",
							"tx_fee_eth": "0.00031572909572282"
						},
						"timestamp": 1710394537,
						"ymd": 20240314
					}
				}
			}

			<트랜잭션 결과 - lock >
			{
				"success": true,
				"data": {
					"state": 200,
					"log": {
						"key": "OWNER65F28ABA65E79148D3FEF853LK",
						"address": "0x8ce5bb2013887ed586e6a87211aa126453368b7a",
						"price": "5",
						"release_time": 1711065600,
						"hash": "0x53e40992b4c8325255ee309e14a1cc1a959a34c664a6db71f63156d4f505266f",
						"tx_fee_eth": "0.000380373603145705",
						"state": 200,
						"state_message": "",
						"timestamp": 1710394517,
						"ymd": 20240314
					}
				}
			}


			<트랜잭션 결과 - unlock >
			{
				"success": true,
				"data": {
					"state": 200,
					"log": {
						"key": "OWNER66068223EA7BA4587AA1B03FUL",
						"address": "0x8ce5bb2013887ed586e6a87211aa126453368b7a",
						"hash": "0xfd37ffa65b84447f29b6fcc5c3f0c76db1ad262424ef1e6187a38dc9822654db",
						"tx_fee_eth": "0.00000875598444279",
						"state": 200,
						"state_message": "",
						"timestamp": 1711932116,
						"ymd": 20240401
					}
				}
			}

			<트랜잭션 결과 - relock >
			{
				"success": true,
				"data": {
					"state": 200,
					"log": {
						"key": "OWNER660683F1991D8F6B57B29C69RL",
						"address": "0x8ce5bb2013887ed586e6a87211aa126453368b7a",
						"price": "77",
						"release_time": 1791678251,
						"unlock_hash": "0x0d64ec5fc6b4f091336f5f83715947f7170165793bfdf8a2c5dbe33bc00d597b",
						"unlock_tx_fee_eth": "0.00000761663727922",
						"lock_hash": "0x2fc96450773f9430f652314a11bab29d7fc54d9d33898db9be56d24d58cc528e",
						"lock_tx_fee_eth": "0.000023016350184718",
						"state": 200,
						"state_message": "",
						"timestamp": 1711932267,
						"ymd": 20240401
					}
				}
			}
		`).
		Apply(doc.Blue)

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {

			cdata_key := ps.ByName("args")

			if cdata_key == "" {
				chttp.Fail(w, ack.BadParam)
				return
			}

			model.DB(func(db mongo.DATABASE) {

				selector := mongo.Bson{"key": cdata_key}

				owner_task := model.OwnerTask{}
				db.C(inf.OwnerTask).Find(selector).One(&owner_task)
				if !owner_task.Valid() {
					chttp.Fail(w, ack.OWNER_NotFoundData)
					return
				}

				_r_state := owner_task.State
				switch _r_state {
				case model.MASTER_NONE:
					_r_state = model.NONE

				case model.MASTER_PENDING:
					_r_state = model.PENDING
				} //switch

				result := RESULT{
					State: _r_state,
				}

				switch owner_task.State {
				case model.SUCCESS:

					switch owner_task.Kind {
					case model.OwnerTaskKind_LT:
						lt_task := owner_task.Owner_LT_Task()
						log := lt_task.CbOwner_LT_Log(owner_task.Key)
						result.Log = log

					case model.OwnerTaskKind_Transfer:
						transfer_task := owner_task.OwnerTransferTask()
						log := transfer_task.CbOwnerTransferLog(owner_task.Key)
						result.Log = log

					case model.OwnerTaskKind_Lock:
						lock_task := owner_task.OwnerLockTask()
						log := lock_task.CbOwnerLockLog(owner_task.Key)
						result.Log = log

					case model.OwnerTaskKind_Unlock:
						unlock_task := owner_task.OwnerUnlockTask()
						log := unlock_task.CbOwnerLockLog(owner_task.Key)
						result.Log = log

					case model.OwnerTaskKind_Relock:
						relock_task := owner_task.OwnerRelockTask()
						result.Log = relock_task.CbOwnerRelockLog(owner_task.Key)

					} //switch

				default:
				} //switch

				chttp.OK(w, result)

			})
		},
	)
}

////////////////////////////////////////////////////////////////////////
//
//

func hOwner_State_Address() {

	method := chttp.GET
	url := model.V1Owner + "/state/address/:args"

	Doc().Comment(OWNER_TAG+"요청 주소의 락업 상태를 WNZ컨트랙트에 조회").
		Method(method).URLS(url, ":args", "조회할 지갑 주소").
		JResultOK(model.LockAccountInfo{}, model.LockAccountInfo{}.TagString()...).
		ETCVAL(model.LockedStateInfo{}, model.LockedStateInfo{}.TagString()...).
		Etc("", `
			<Response-Sample>


			{
				"success": true,
				"data": {
					"account": "0x8ce5bb2013887ed586e6a87211aa126453368b7a",
					"total_price": "10000069.95",
					"locked_total_price": "77",		//locked_state_infos의 locked_price의 총합
					"locked_calc_price": "77",		//locked_state_infos의 is_time_over(false)의 총합 (실제 Lock된 금액)
					"locked_state_infos": [			//locked_state_infos의 데이터는 1개여야 합니다.
						{
							"position_index": 0,
							"release_time": 1791678251,
							"locked_price": "77",
							"is_time_over": false	//true이면 해당 데이터는 언락 되었음.(account가 토큰을 전송시 해당 데이터는 삭제됨.
						}
					],
					"response_server_time": 1711932889
				}
			}

			
		`).
		Apply(doc.Blue)

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {

			cdata_address := ps.ByName("args")

			if !ebcm.IsAddressP(&cdata_address) {
				chttp.Fail(w, ack.BadParam)
				return
			}

			erc20_info := inf.FirstERC20()

			caller := ecs.New(
				ecs.RPC_URL(inf.Mainnet()),
				inf.InfuraKey(),
			)

			rpc_lock_user_info, err := model.LockTokenUtil{}.GetUserLockInfoAll(
				caller,
				"유틸컨트랙트 안씀",
				erc20_info.Contract,
				cdata_address,
			)
			if err != nil {
				chttp.Fail(w, ack.DBJob, dbg.Cat("Blockchain node rpc error :", err))
				return
			}

			result := rpc_lock_user_info.LockUserInfo(erc20_info.Decimal)

			chttp.OK(w, result)
		},
	)
}

////////////////////////////////////////////////////////////////////////
//
//

func hOwner_Request_Lock() {

	type RESULT struct {
		Key model.TaskKey `json:"key"`
	}

	method := chttp.POST
	url := model.V1Owner + "/request/lock"

	Doc().Comment(OWNER_TAG+"Lock 요청").
		Method(method).URL(url).
		JParam(model.ReqOwnerLockTry{}, model.ReqOwnerLockTry{}.TagString()...).
		JResultOK(
			RESULT{},
			"key", "결과 조회/콜백 용 Key값",
		).
		//ETCVAL(model.ReqOwnerTransferItem{}, model.ReqOwnerTransferItem{}.TagString()...).
		JAckError(ack.InvalidAddress).
		JAckError(ack.OWNER_ALREADY_LOCKED).
		JAckError(ack.OWNER_REQ_JOB_ONCE).
		JAckError(ack.OWNER_RpcFail).
		JAckError(ack.OWNER_Lock_OverPrice).
		JAckError(ack.DBJob).
		Apply(doc.Blue)

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {

			cdata := chttp.BindingStruct[model.ReqOwnerLockTry](req)
			_ = cdata

			if !ebcm.IsAddressP(&cdata.Address) {
				chttp.Fail(w, ack.InvalidAddress)
				return
			}

			if jmath.CMP(cdata.Price, 0) <= 0 {
				chttp.Fail(w, ack.BadParam, "price under zero.")
				return
			}

			if ack_err := check_lock_ban(cdata.Address); ack_err != nil {
				chttp.Fail(w, ack_err)
				return
			}

			token_info := inf.FirstERC20()
			caller := inf.GetFinder()

			ab_price, err := model.LockTokenUtil{}.AvailablePriceOf(
				caller,
				token_info.Contract,
				cdata.Address,
				token_info.Decimal,
			)
			if err != nil {
				chttp.Fail(w, ack.OWNER_RpcFail)
				return
			}

			if jmath.CMP(cdata.Price, ab_price) > 0 {
				chttp.Fail(w, ack.OWNER_Lock_OverPrice, "request over price (", ab_price, ")")
				return
			}

			OWNER_DB(w, cdata.Address, func(db mongo.DATABASE) {
				task := cdata.GetOwnerTask(token_info.Decimal)

				if err := db.C(inf.OwnerTask).Insert(task); err != nil {
					chttp.Fail(w, ack.DBJob)
					return
				}

				chttp.OK(w,
					RESULT{
						Key: task.Key,
					},
				)
			})

		},
	)
}

////////////////////////////////////////////////////////////////////////
//
//

func hOwner_Request_Unlock() {

	type RESULT struct {
		Key model.TaskKey `json:"key"`
	}

	method := chttp.POST
	url := model.V1Owner + "/request/unlock"

	Doc().Comment(OWNER_TAG+"Unlock 요청").
		Method(method).URL(url).
		JParam(model.ReqOwnerUnlockTry{}, model.ReqOwnerUnlockTry{}.TagString()...).
		JResultOK(
			RESULT{},
			"key", "결과 조회/콜백 용 Key값",
		).
		//ETCVAL(model.ReqOwnerTransferItem{}, model.ReqOwnerTransferItem{}.TagString()...).
		JAckError(ack.InvalidAddress).
		JAckError(ack.OWNER_REQ_JOB_ONCE).
		JAckError(ack.DBJob).
		Apply(doc.Blue)

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {

			cdata := chttp.BindingStruct[model.ReqOwnerUnlockTry](req)
			_ = cdata

			if !ebcm.IsAddressP(&cdata.Address) {
				chttp.Fail(w, ack.InvalidAddress)
				return
			}

			OWNER_DB(w, cdata.Address, func(db mongo.DATABASE) {
				// err := model.OwnerUnlockPool{}.InsertDB(db, cdata.Address)
				// if err != nil {
				// 	chttp.Fail(w, ack.AlreadyProcessJob)
				// 	return
				// }

				chttpFail := func(w http.ResponseWriter, e chttp.CError, etc ...interface{}) {
					//model.OwnerUnlockPool{}.RemoveDB(db, cdata.Address)
					chttp.Fail(w, e, etc...)
				}

				is_locked, ack_err := cloud.CheckLockState(cdata.Address)
				if ack_err != nil {
					chttp.Fail(w, ack_err)
					return
				}
				if !is_locked {
					chttp.Fail(w, ack.BadParam, "락업 정보가 없습니다.")
					return
				}

				token_info := inf.FirstERC20()
				task := cdata.GetOwnerTask(token_info.Decimal)

				if err := db.C(inf.OwnerTask).Insert(task); err != nil {
					chttpFail(w, ack.DBJob)
					return
				}

				chttp.OK(w,
					RESULT{
						Key: task.Key,
					},
				)
			})

		},
	)
}

////////////////////////////////////////////////////////////////////////
//
//

func hOwner_Request_ReLock() {
	type RESULT struct {
		Key model.TaskKey `json:"key"`
	}

	method := chttp.POST
	url := model.V1Owner + "/request/relock"

	Doc().Comment(OWNER_TAG+"이미 Lock되어있는 계정의 ReLock 설정").
		Method(method).URL(url).
		JParam(model.ReqOwnerRelockTry{}, model.ReqOwnerRelockTry{}.TagString()...).
		JResultOK(
			RESULT{},
			"key", "결과 조회/콜백 용 Key값",
		).
		//ETCVAL(model.ReqOwnerTransferItem{}, model.ReqOwnerTransferItem{}.TagString()...).
		JAckError(ack.InvalidAddress).
		JAckError(ack.OWNER_REQ_JOB_ONCE).
		JAckError(ack.DBJob).
		Apply(doc.Blue)

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {

			cdata := chttp.BindingStruct[model.ReqOwnerRelockTry](req)
			_ = cdata

			if !ebcm.IsAddressP(&cdata.Address) {
				chttp.Fail(w, ack.InvalidAddress)
				return
			}

			if jmath.CMP(cdata.Price, 0) <= 0 {
				chttp.Fail(w, ack.BadParam, "price under zero.")
				return
			}

			OWNER_DB(w, cdata.Address, func(db mongo.DATABASE) {
				is_locked, ack_err := cloud.CheckLockState(cdata.Address)
				if ack_err != nil {
					chttp.Fail(w, ack_err)
					return
				}
				if !is_locked {
					chttp.Fail(w, ack.BadParam, "락업 정보가 없습니다.")
					return
				}

				caller := inf.GetFinder()
				token_info := inf.FirstERC20()
				balance, _ := model.Erc20Balance(
					caller, token_info.Contract,
					cdata.Address,
				)
				price := ebcm.WeiToToken(balance, token_info.Decimal)
				if jmath.CMP(price, cdata.Price) < 0 {
					chttp.Fail(w, ack.OWNER_Lock_OverPrice)
					return
				}

				task := cdata.GetOwnerTask(token_info.Decimal)
				if err := db.C(inf.OwnerTask).Insert(task); err != nil {
					chttp.Fail(w, ack.DBJob)
					return
				}

				chttp.OK(w,
					RESULT{
						Key: task.Key,
					},
				)

			})

		},
	)
}

// ////////////////////////////////////////////////////////////////////////
// /
// /
// /

func h_owner_callback_api_doc() {
	Doc().Message("---- " + OWNER_TAG + " 콜백 API ----------------------")

	Doc().Message(`
	<cc_blue>` + OWNER_TAG + `[ SERVICE CALLBACK ]</cc_blue> ( 0. 오너 토큰 락업전송 콜백 ) WNZ Transfer&Lock

	<cc_red>` + model.OwnerCallbackApi_LT + `</cc_red>

	method : post
	content-type : application/json
	
	data :
	{
		"key": "OWNER65F286CE49305605EC991FD4LT",		// 발급한 KEY 값
		"lock_info": {	//<step2> 락업정보
			"fail_hash_list": [],
			"from": "0xff00bbd12f9ae64fe4dfa9f3dddfdc3d9ecc2378",	//오너주소
			"hash": "0xd0f41c9d99a4f898d6d80d1c87e548604005728d5e65029ceada97dfa8d92e9a",
			"state": 200,
			"state_message": "",
			"tx_fee_eth": "0.000270179305333445"
		},
		"price": "1.5",												// 전송 금액
		"recipient": "0x8ce5bb2013887ed586e6a87211aa126453368b7a",	// 실제 토큰을 받은 지갑주소
		"release_time": 1710504000,									// 언락 시간
		"timestamp": 1710393292,
		"transfer_info": {	//<step1> 전송정보
			"fail_hash_list": [],
			"from": "0x7d44dbb4a0fa180774611412edeacdd2e7a13ec8",	//마스터주소
			"hash": "0xb4732c00d530e4b0cfe20d1d6d3d910e08de0842523019574e94269d4200d171",
			"state": 200,
			"state_message": "",
			"tx_fee_eth": "0.00012485058381433"
		},
		"ymd": 20240314
	}
	`)

	Doc().Message(`
	<cc_blue>` + OWNER_TAG + `[ SERVICE CALLBACK ]</cc_blue> ( 2. Lock ) WNZ 홀더 주소의 잔액을 Lock

	<cc_red>` + model.OwnerCallbackApi_Lock + `</cc_red>

	method : post
	content-type : application/json
	
	data :
	{
		"address": "0x8ce5bb2013887ed586e6a87211aa126453368b7a",						// WNZ 홀더 주소
		"hash": "0xbc00abc356735c2919c8f2cf8edb0348cc7665874dc2c04e7586d64e83fe42bc",	// Tx-hash
		"key": "OWNER65DF27C5708544C184CF4189LK",										// 발급 키값
		"price": "0.1",																	// Lock 시킨 금액
		"release_time": 1711929600,														// Unlock 될 시간
		"state": 200,																	// 200(성공) , 104(실패)
		"state_message": "",															// 104(실패)이면 무조껀 상태 메시지 제공
		"timestamp": 1709123545,														// 서버 시간
		"tx_fee_eth": "0.00008456499983087",											// Tx-Fee-ETH
		"ymd": 20240228																	// 날짜
	}

	> 104(실패) 일때 "state_message" 
	- "[FAIL_AvailableBalanceOf] ab_price: 0.01" : 홀더 주소의 사용 가능한 잔액이 요청한 Lock시킬 금액보다 적을 경우 ( 0.01은 홀더의 사용 가능한 금액 )
	- "[FAIL_EstimateGas] : error_message"		 : Owner의 트랜잭션 가스 예측 실패.
	- "[FAIL_ReceiptByHash]"					 : 트랜잭션이 실패 했을경우.

	`)

	Doc().Message(`
	<cc_blue>` + OWNER_TAG + `[ SERVICE CALLBACK ]</cc_blue> ( 3. Unlock ) WNZ 홀더 주소의 잔액을 UnLock (해제).

	<cc_red>` + model.OwnerCallbackApi_Unlock + `</cc_red>

	method : post
	content-type : application/json
	
	data :
	{
		"address": "0x8ce5bb2013887ed586e6a87211aa126453368b7a",			// 홀더 주소
		"hash": "0xfd37ffa65b84447f29b6fcc5c3f0c76db1ad262424ef1e6187a38dc9822654db",
		"key": "OWNER66068223EA7BA4587AA1B03FUL",
		"state": 200,
		"state_message": "",
		"timestamp": 1711702580,
		"tx_fee_eth": "0.00000875598444279",
		"ymd": 20240329
	}

	> 104(실패) 일때 "state_message" 
	- "[FAIL_EstimateGas] : error_message"	: Owner의 트랜잭션 가스 예측 실패.
	- "[FAIL_ReceiptByHash]"				: 트랜잭션이 실패 했을경우.


	> 200(성공)일때 "state_message" 
	- "[lock_time_over]"	: Unlock을 하기 전에 데이터 검증 과정에서 이미 락업해제 시간이 지나서 unlock을 할 필요가 없을경우.
							( 위의 경우에는 트랜잭션을 쏘지 않기 때문에 hash, tx_fee_eth 값이 존재하지 않는다. )

	

	`)

	Doc().Message(`
	<cc_blue>` + OWNER_TAG + `[ SERVICE CALLBACK ]</cc_blue> ( 3. Relock ) Lock 되어있는 WNZ 홀더 주소의 잔액/기간을 새로운 데이터로 다시 락.

	<cc_red>` + model.OwnerCallbackApi_Relock + `</cc_red>

	method : post
	content-type : application/json
	
	data :
	{
		"address": "0x8ce5bb2013887ed586e6a87211aa126453368b7a",							// WNZ홀더 주소
		"key": "OWNER660682B8EA7BA4587AA1B045RL",											// 발급한 키값
		"lock_hash": "0x21d58fd707e7459779d62edfeaa6e12f3cc73428b37c987987ad86b9c92b2c38",	// 새로운 lock Hash
		"lock_tx_fee_eth": "0.00002385916596274",
		"price": "11",																		// 새로운 lock 수량
		"release_time": 1791678251,															// 새로운 unlock 시간
		"state": 200,
		"state_message": "",
		"timestamp": 1711702739,
		"unlock_hash": "0x26f94221af3c16a2877a67b89f10708fbc934f06c76f456796dbc9adf2993fd0",	//이전 lock정보 unlock hash
		"unlock_tx_fee_eth": "0.0000078963420022",
		"ymd": 20240329
	}

	> 104(실패) 일때 "state_message" 
	- "[FAIL_EstimateGas] : error_message"	: Owner의 트랜잭션 가스 예측 실패.
	- "[FAIL_ReceiptByHash]"				: 트랜잭션이 실패 했을경우.

	

	`)

}
