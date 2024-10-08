package cloud

import (
	"context"
	"jtools/cloud/ebcm"
	"jtools/cloud/ebcm/abi"
	"jtools/jmath"
	"time"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func owner_task_lock(db mongo.DATABASE, owner_task model.OwnerTask) {
	task_data := owner_task.OwnerLockTask()

	set_state := func(state model.STATE) {
		owner_task.State = state
		task_data.State = state
	}

	token_info := inf.FirstERC20()
	owner := inf.Owner()

	sender := inf.GetSender()

	TAG := dbg.Cat("[TASK_LOCK](", owner_task.Key, ") ")
	TaskLog(TAG, "-------------- START")
	defer TaskLog(TAG, "-------------- END")

	for {
		switch owner_task.State {
		case model.NONE:

			availble_balance, err := model.LockTokenUtil{}.AvailableBalanceOf(
				sender,
				token_info.Contract,
				task_data.Address,
			)
			if err != nil {
				TaskError("AvailableBalanceOf :", err)
				return
			}

			is_locked, ack_err := CheckLockState(task_data.Address)
			if ack_err != nil {
				TaskError(TAG, "CheckLockState ----- Call_Fail")
				return
			}

			if is_locked { //이미 락이 되어있음...
				set_state(model.FAIL)
				task_data.StateMessage = "ALREADY_LOCKED_DATA"

				owner_task.UpdateDB(db, task_data)

				log := task_data.CbOwnerLockLog(owner_task.Key)
				model.OwnerLog{}.InsertDB(db, log)
				return
			}

			if jmath.CMP(availble_balance, task_data.Amount) < 0 {
				set_state(model.FAIL)
				task_data.StateMessage = dbg.Cat(
					"[FAIL_AvailableBalanceOf] ab_price:",
					ebcm.WeiToToken(availble_balance, token_info.Decimal),
				)

				owner_task.UpdateDB(db, task_data)

				log := task_data.CbOwnerLockLog(owner_task.Key)
				model.OwnerLog{}.InsertDB(db, log)
				return
			}

			owner_coin_value := sender.Balance(owner.Address)
			if jmath.CMP(owner_coin_value, 0) <= 0 {
				//show me the coin!
				TaskError(TAG, "Owner Coin Balance is 0.")
				return
			}

			nonce, err := _getNonce(sender, owner.Address)
			if err != nil {
				TaskError(TAG, "_getNonce :", err)
				return
			}

			pad_bytes := ebcm.MakePadBytesABI(
				"lock",
				abi.TypeList{
					abi.NewAddress(task_data.Address),
					abi.NewUint256(task_data.Amount),
					abi.NewUint256(task_data.ReleaseTime),
				},
			)

			limit, err := sender.EstimateGas(
				context.Background(),
				ebcm.MakeCallMsg(
					owner.Address,
					token_info.Contract,
					"0",
					pad_bytes,
				),
			)
			if err != nil {
				set_state(model.FAIL)
				task_data.StateMessage = dbg.Cat("[FAIL_EstimateGas] ", err.Error())

				owner_task.UpdateDB(db, task_data)

				log := task_data.CbOwnerLockLog(owner_task.Key)
				model.OwnerLog{}.InsertDB(db, log)

				return
			}

			limit = ebcm.MMA_LimitBuffer_MasterOut(limit)

			gas_price, err := sender.SuggestGasPrice(context.Background())
			if err != nil {
				TaskError(TAG, "SuggestGasPrice :", err)
				return
			}
			gas_price = CALC_GAS_PRICE(db, gas_price)

			fee_wei := gas_price.EstimateGasFeeWEI(limit)

			if jmath.CMP(owner_coin_value, fee_wei) < 0 {
				TaskError(TAG, "OwnerCoin(", ebcm.WeiToETH(owner_coin_value), ") / Fee(", ebcm.WeiToETH(fee_wei), ") ")
				return
			}

			tx_try_data, stx := model.MakeTxTryData(
				sender,
				owner.PrivateKey,

				nonce,
				token_info.Contract,
				"0",

				limit,
				gas_price,
				pad_bytes,
			)

			hash, err := sender.SendTransaction(
				context.Background(),
				stx,
			)
			if err != nil {
				TaskError(TAG, "SendTransaction :", err)
				return
			}

			task_data.TryData = tx_try_data
			task_data.Hash = hash
			set_state(model.PENDING)

			owner_task.UpdateDB(db, task_data)
			TaskLog(TAG, "SendTransaction :::", hash)

		case model.PENDING:
			receipt := sender.ReceiptByHash(task_data.Hash)
			if !receipt.IsNotFound {
				task_data.TxFeeETH = receipt.TxFeeETH()

				if !receipt.IsError() {
					set_state(model.SUCCESS)

					owner_task.UpdateDB(db, task_data)

				} else {
					set_state(model.FAIL)
					task_data.StateMessage = "[FAIL_ReceiptByHash]"

					owner_task.UpdateDB(db, task_data)
				}

				log := task_data.CbOwnerLockLog(owner_task.Key)
				model.OwnerLog{}.InsertDB(db, log)

			} else {

				dt := model.GetGAS(db).GetLockUnLockPendingMin()
				if task_data.TryData.IsTimeOver(dt) {
					stx := task_data.TryData.STX(
						sender,
						owner.PrivateKey,
					)
					hash, err := sender.SendTransaction(
						context.Background(),
						stx,
					)
					if err == nil {
						task_data.Hash = hash
						owner_task.UpdateDB(db, task_data)
						TaskLog(TAG, "RE_SendTransaction :::", hash)
						continue
					}
				}

				time.Sleep(time.Second)
				TaskLog(TAG, "PendingWait :::", task_data.Hash)

			}

		default:
			TaskLog(TAG, owner_task.State)
			return
		} //switch
	} //for

	//////////////////////////

}
