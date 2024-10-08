package cloud

import (
	"context"
	"jtools/cloud/ebcm"
	"jtools/cloud/ebcm/abi"
	"time"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func owner_task_relock(db mongo.DATABASE, owner_task model.OwnerTask) {
	task_data := owner_task.OwnerRelockTask()

	token_info := inf.FirstERC20()
	owner := inf.Owner()

	sender := inf.GetSender()

	TAG := dbg.Cat("[TASK_RELOCK](", owner_task.Key, ") ")
	TaskLog(TAG, "-------------- START")
	defer TaskLog(TAG, "-------------- END")

	make_unlock_tx := func(gas_price ebcm.GasPrice, nonce uint64) (model.TxTryData, ebcm.WrappedTransaction) {
		unlock_data := ebcm.MakePadBytesABI(
			"releaseLock",
			abi.TypeList{
				abi.NewAddress(task_data.Address),
				//abi.NewUint256(task_data.PositionIndex),
			},
		)

		limit, _ := sender.EstimateGas(
			context.Background(),
			ebcm.MakeCallMsg(
				owner.Address,
				token_info.Contract,
				"0",
				unlock_data,
			),
		)
		limit = ebcm.MMA_LimitBuffer_MasterOut(limit)

		return model.MakeTxTryData(
			sender,
			owner.PrivateKey,

			nonce,
			token_info.Contract,
			"0",
			limit,
			gas_price,
			unlock_data,
		)
	}
	make_lock_tx := func(gas_price ebcm.GasPrice, nonce uint64) (model.TxTryData, ebcm.WrappedTransaction) {
		lock_data := ebcm.MakePadBytesABI(
			"lock",
			abi.TypeList{
				abi.NewAddress(task_data.Address),
				abi.NewUint256(task_data.Amount),
				abi.NewUint256(task_data.ReleaseTime),
			},
		)
		limit := uint64(250000)

		return model.MakeTxTryData(
			sender,
			owner.PrivateKey,

			nonce,
			token_info.Contract,
			"0",
			limit,
			gas_price,
			lock_data,
		)
	}

	for {

		switch owner_task.State {
		case model.NONE: //first_action

			is_locked, ack_err := CheckLockState(task_data.Address)
			if ack_err != nil {
				TaskError(TAG, "CheckLockState ----- Call_Fail")
				return
			}

			if !is_locked { //언락시킬 데이터가 없음
				owner_task.State = model.FAIL
				task_data.StateMessage = "ALREADY_UNLOCKED_DATA"

				owner_task.UpdateDB(db, task_data)

				log := task_data.CbOwnerRelockLog(owner_task.Key)
				model.OwnerLog{}.InsertDB(db, log)
				return
			}

			gas_price, err := sender.SuggestGasPrice(context.Background())
			if err != nil {
				TaskError(TAG, "SuggestGasPrice :", err)
				return
			}
			//cc.Gray("base_gas :", gas_price.GET_GAS_GWEI())
			gas_price = CALC_GAS_PRICE(db, gas_price)
			//cc.Gray("add_gas :", gas_price.GET_GAS_GWEI())

			base_nonce, err := _getNonce(sender, owner.Address)
			if err != nil {
				TaskError(TAG, "_getNonce :", err)
				return
			}

			task_data.From = owner.Address

			/////////////////////////////////////////////////////

			unlock_nonce := base_nonce
			tx_try_data_1, unlock_stx := make_unlock_tx(gas_price, unlock_nonce)

			lock_nonce := base_nonce + 1
			tx_try_data_2, lock_stx := make_lock_tx(gas_price, lock_nonce)

			//////////////////////////////////////////////////////

			unlock_hash, err := sender.SendTransaction(context.Background(), unlock_stx)
			if err != nil {
				TaskError(TAG, "unlock.SendTx :", err)
				return
			}
			TaskLog(TAG, "unlock.SendTransaction :::", unlock_hash)
			task_data.UnlockHash = unlock_hash
			task_data.UnlockTryData = tx_try_data_1

			lock_hash, err := sender.SendTransaction(context.Background(), lock_stx)
			if err != nil {
				TaskError(TAG, "lock.SendTx :", err)
				continue
			}
			TaskLog(TAG, "lock.SendTransaction :::", lock_hash)
			task_data.LockHash = lock_hash
			task_data.LockTryData = tx_try_data_2

			///////////////////////////////////////////////////////
			owner_task.State = model.PENDING
			owner_task.UpdateDB(db, task_data)
			///////////////////////////////////////////////////////

		case model.PENDING:
			r1 := sender.ReceiptByHash(task_data.UnlockHash)
			r2 := sender.ReceiptByHash(task_data.LockHash)

			dt := model.GetGAS(db).GetLockUnLockPendingMin()
			if r1.IsNotFound {
				if task_data.UnlockTryData.IsTimeOver(dt) {
					stx_1 := task_data.UnlockTryData.STX(sender, owner.PrivateKey)
					hash_1, _ := sender.SendTransaction(context.Background(), stx_1)

					stx_2 := task_data.LockTryData.STX(sender, owner.PrivateKey)
					hash_2, _ := sender.SendTransaction(context.Background(), stx_2)

					task_data.UnlockHash = hash_1
					task_data.LockHash = hash_2
					owner_task.UpdateDB(db, task_data)

				}
				time.Sleep(time.Second)
				continue
			}

			if r2.IsNotFound {
				if task_data.LockTryData.IsTimeOver(dt) {
					stx_2 := task_data.LockTryData.STX(sender, owner.PrivateKey)
					hash_2, _ := sender.SendTransaction(context.Background(), stx_2)
					task_data.LockHash = hash_2
					owner_task.UpdateDB(db, task_data)
				}
				time.Sleep(time.Second)
				continue
			}

			task_data.UnlockTxFeeETH = r1.TxFeeETH()
			task_data.LockTxFeeETH = r2.TxFeeETH()

			is_s1 := !r1.IsError()
			is_s2 := !r2.IsError()

			success := is_s1 && is_s2
			if success {
				owner_task.State = model.SUCCESS
				task_data.State = model.SUCCESS
			} else {
				owner_task.State = model.FAIL
				task_data.State = model.FAIL
				task_data.StateMessage = dbg.Cat("Unlock(", is_s1, "), Lock(", is_s2, ")")
			}
			owner_task.UpdateDB(db, task_data)

			log := task_data.CbOwnerRelockLog(owner_task.Key)
			model.OwnerLog{}.InsertDB(db, log)

		default:
			TaskLog(TAG, owner_task.State)
			return

		} //switch
	} //for(0)

}
