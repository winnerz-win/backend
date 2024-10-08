package cloud

import (
	"jtools/cc"
	"jtools/cloud/ebcm"
	"jtools/mms"
	"time"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jnet/cnet"
	"txscheduler/brix/tools/runtext"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func TaskLog(a ...any) {
	cc.CyanItalic(a...)

}
func TaskError(a ...any) {
	cc.RedItalic(a...)
}

func runOwnerTaskAction(rtx runtext.Runner) {
	dbg.PrintForce("cloud.runOwnerTaskAction ---- END")
	<-rtx.WaitStart()
	dbg.PrintForce("cloud.runOwnerTaskAction ---- START")

EXIT:
	for {
		select {
		case <-rtx.EndC():
			break EXIT
		default:
		} //select

		model.DB(func(db mongo.DATABASE) {
			owner_job_action(db)
		})

		time.Sleep(time.Second)

	} //for

}

func owner_job_action(db mongo.DATABASE) bool {

	owner_task := model.OwnerTask{}
	db.C(inf.OwnerTask).Find(
		mongo.Bson{"state": mongo.Bson{"$lte": model.PENDING}},
	).Sort("create_at").One(&owner_task)

	if !owner_task.Valid() {
		return false
	}

	switch owner_task.Kind {
	case model.OwnerTaskKind_LT: //lockup
		owner_lt_lock(db, owner_task)

	case model.OwnerTaskKind_Transfer:
		owner_task_transfer(db, owner_task)

	case model.OwnerTaskKind_Lock:
		owner_task_lock(db, owner_task)

	case model.OwnerTaskKind_Unlock:
		owner_task_unlock(db, owner_task)

	case model.OwnerTaskKind_Relock:
		owner_task_relock(db, owner_task)
	} //switch

	return true
}

func _getNonce(caller *ebcm.Sender, account string) (uint64, error) {
	return ebcm.MMA_GetNonce(
		caller,
		account,
	)
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func runOwnerTaskLogSender() {
	dbg.PrintForce("cloud.runOwnerTaskLogSender ---- START")

	for {
		model.DB(func(db mongo.DATABASE) {

			list := model.OwnerLog{}.GetList(db)
			if len(list) == 0 {
				return
			}

			client_url := inf.ClientAddress()
			for _, item := range list {

				callback_url := ""
				var params any

				switch item.Kind {
				case model.OwnerTaskKind_LT:
					params, _ = dbg.DecodeStruct[model.CbOwner_LT_Log](item.Log)
					callback_url = model.OwnerCallbackApi_LT

				case model.OwnerTaskKind_Transfer:
					params, _ = dbg.DecodeStruct[model.CbOwnerTransferLog](item.Log)
					callback_url = model.OwnerCallbackApi_Transfer

				case model.OwnerTaskKind_Lock:
					params, _ = dbg.DecodeStruct[model.CbOwnerLockLog](item.Log)
					callback_url = model.OwnerCallbackApi_Lock

				case model.OwnerTaskKind_Unlock:
					params, _ = dbg.DecodeStruct[model.CbOwnerUnlockLog](item.Log)
					callback_url = model.OwnerCallbackApi_Unlock

				case model.OwnerTaskKind_Relock:
					params, _ = dbg.DecodeStruct[model.CbOwnerRelockLog](item.Log)
					callback_url = model.OwnerCallbackApi_Relock

				} //switch

				ack := cnet.POST_JSON_F(
					client_url+callback_url,
					nil,
					params,
				)
				if err := ack.Error(); err != nil {
					dbg.RedItalic("owner_log.send_fail :", err)
					break
				}
				item.SendOK(db, mms.Now())
			} //for
		})

		time.Sleep(time.Second)
	} //for

}

/*

function transferWithLock(
	address recipient,
	uint256 amount,
	uint256 releaseTime
) external onlyOwner returns (bool success)


function transfer(address recipient, uint256 amount) external virtual override returns (bool)


function lock(address recipient, uint256 amount, uint256 releaseTime) external onlyOwner returns (bool success) {
	_lock(recipient, amount, releaseTime);
	success = true;
}

function unlock(address from, uint256 idx) external onlyOwner returns (bool success) {
	require(_locks[from].length > idx, "There is not lock info.");
	_unlock(from, idx);
	success = true;
}



*/
