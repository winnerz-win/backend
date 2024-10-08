package cloud

import (
	"context"
	"jtools/cloud/ebcm"
	"jtools/cloud/ebcm/abi"
	"jtools/jmath"
	"jtools/mms"
	"time"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func owner_lt_master_pending(db mongo.DATABASE, _ mms.MMS) bool {
	owner_task := model.OwnerTask{}
	db.C(inf.OwnerTask).Find(
		mongo.Bson{"state": model.MASTER_PENDING}, //11
	).Sort("create_at").One(&owner_task)

	if !owner_task.Valid() {
		return false
	}

	task_data := owner_task.Owner_LT_Task()
	tx_info := &task_data.TransferInfo

	sender := inf.GetSender()

	receipt := sender.ReceiptByHash(tx_info.Hash)
	if !receipt.IsNotFound {
		if !receipt.IsError() {
			tx_info.State = model.SUCCESS
			tx_info.TxFeeETH = receipt.TxFeeETH()

			owner_task.State = model.NONE // >>toss>> lock try
			owner_task.UpdateDB(db, task_data)

		} else {
			tx_info.State = model.NONE
			tx_info.StateMessage = dbg.Cat(
				tx_info.StateMessage,
				dbg.Cat("fail_hash(", tx_info.Hash, "),"),
			)

			owner_task.State = model.MASTER_NONE //fail -> retry
			owner_task.UpdateDB(db, task_data)

		}

	} else {
		dt := model.GetGAS(db).GetLtTransferPendingMin()
		if tx_info.TryData.IsTimeOver(dt) { //lt_transer
			master := inf.Master()
			stx := tx_info.TryData.STX(
				sender,
				master.PrivateKey,
			)

			hash, err := sender.SendTransaction(context.Background(), stx)
			if err == nil {
				tx_info.Hash = hash
				owner_task.UpdateDB(db, task_data)
				TaskLog("LT(Master) RE_SendTransaction ::::", tx_info.Hash)
				return true
			}
		}

		time.Sleep(time.Second)
		TaskLog("LT(Master) PendingWait ::::", tx_info.Hash)
	}

	return true
}

func owner_lt_master_transfer(db mongo.DATABASE, _ mms.MMS) bool {
	owner_task := model.OwnerTask{}
	db.C(inf.OwnerTask).Find(
		mongo.Bson{"state": model.MASTER_NONE}, //10
	).Sort("create_at").One(&owner_task)

	if !owner_task.Valid() {
		return false
	}

	task_data := owner_task.Owner_LT_Task()
	_ = task_data.TransferInfo

	token_info := inf.FirstERC20()
	master := inf.Master()

	sender := inf.GetSender()

	TAG := dbg.Cat("[LT_TRANSFER](", owner_task.Key, ")")
	TaskLog(TAG, "-------------- START")
	defer TaskLog(TAG, "-------------- END")

	nonce, err := _getNonce(sender, master.Address)
	if err != nil {
		TaskError("_getNonce :", err)
		return false
	}

	master_token_value, err := model.Erc20Balance(
		sender,
		token_info.Contract,
		master.Address,
	)
	if err != nil {
		TaskError("LT(Master) Erc20Balance :", err)
		return false
	}

	if jmath.CMP(master_token_value, task_data.Amount) < 0 {
		TaskError(dbg.Cat(
			"LT(Master) Master Token Balance (",
			ebcm.WeiToToken(master_token_value, token_info.Decimal),
			"/",
			ebcm.WeiToToken(task_data.Amount, token_info.Decimal),
			")",
		))
		return false
	}

	master_coin_value := sender.Balance(master.Address)
	if jmath.CMP(master_coin_value, 0) <= 0 {
		TaskError("LT(Master) Master Coin Balance is 0.")
		return false
	}

	pad_bytes := ebcm.MakePadBytesABI(
		"transfer",
		abi.TypeList{
			abi.NewAddress(task_data.Recipient),
			abi.NewUint256(task_data.Amount),
		},
	)
	limit, err := sender.EstimateGas(
		context.Background(),
		ebcm.MakeCallMsg(
			master.Address,
			token_info.Contract,
			"0",
			pad_bytes,
		),
	)
	if err != nil {
		TaskError("LT(Master) EstimateGas :", err)
		return false
	}

	limit = ebcm.MMA_LimitBuffer(limit)

	gas_price, err := sender.SuggestGasPrice(context.Background())
	if err != nil {
		TaskError("LT(Master) SuggestGasPrice :", err)
		return false
	}
	gas_price = CALC_GAS_PRICE(db, gas_price) //transfer

	fee_wei := gas_price.EstimateGasFeeWEI(limit)
	if jmath.CMP(master_coin_value, fee_wei) < 0 {
		TaskError(TAG, "LT(Master) MasterrCoin(", ebcm.WeiToETH(master_coin_value), ") / Fee(", ebcm.WeiToETH(fee_wei), ") ")
		return false
	}

	tx_try_data, stx := model.MakeTxTryData(
		sender,
		master.PrivateKey,

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
		TaskError("LT(Master) SendTransaction :", err)
		return false
	}

	task_data.TransferInfo.From = master.Address
	task_data.TransferInfo.Hash = hash
	task_data.TransferInfo.State = model.PENDING
	task_data.TransferInfo.TryData = tx_try_data

	owner_task.State = model.MASTER_PENDING
	owner_task.UpdateDB(db, task_data)

	time.Sleep(time.Second)
	TaskLog("LT(Master) SendTransaction :::", hash)

	return true
}

/////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////

func owner_lt_lock(db mongo.DATABASE, owner_task model.OwnerTask) {
	task_data := owner_task.Owner_LT_Task()

	token_info := inf.FirstERC20()
	owner := inf.Owner()

	sender := inf.GetSender()

	TAG := dbg.Cat("[LT_LOCK](", owner_task.Key, ")")
	TaskLog(TAG, "-------------- START")
	defer TaskLog(TAG, "-------------- END")

	for {
		switch owner_task.State {
		case model.NONE:
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
					abi.NewAddress(task_data.Recipient),
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
				TaskError(TAG, "EstimateGas :", err)
				return
			}

			limit = ebcm.MMA_LimitBuffer_MasterOut(limit)

			gas_price, err := sender.SuggestGasPrice(context.Background())
			if err != nil {
				TaskError(TAG, "SuggestGasPrice :", err)
				return
			}
			gas_price = CALC_GAS_PRICE(db, gas_price) //lock

			fee_wei := gas_price.EstimateGasFeeWEI(limit)

			if jmath.CMP(owner_coin_value, fee_wei) < 0 {
				TaskError(TAG, "OwnerCoin(", owner_coin_value, ") / Fee(", fee_wei, ") ")
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

			task_data.LockInfo.From = owner.Address
			task_data.LockInfo.Hash = hash
			task_data.LockInfo.State = model.PENDING
			task_data.LockInfo.TryData = tx_try_data

			owner_task.State = model.PENDING

			owner_task.UpdateDB(db, task_data)
			TaskLog(TAG, "SendTransaction :::", hash)

		case model.PENDING:
			receipt := sender.ReceiptByHash(task_data.LockInfo.Hash)
			if !receipt.IsNotFound {

				if !receipt.IsError() {
					task_data.LockInfo.TxFeeETH = receipt.TxFeeETH()
					task_data.LockInfo.State = model.SUCCESS

					owner_task.State = model.SUCCESS
					owner_task.UpdateDB(db, task_data)

					log := task_data.CbOwner_LT_Log(owner_task.Key)
					model.OwnerLog{}.InsertDB(db, log)

				} else { //검증 실패면 재시도

					task_data.LockInfo.State = model.NONE
					task_data.LockInfo.StateMessage = dbg.Cat(
						task_data.LockInfo.StateMessage,
						dbg.Cat("fail_hash(", task_data.LockInfo.Hash, "),"),
					)

					owner_task.State = model.NONE //fail - retry
					owner_task.UpdateDB(db, task_data)
					time.Sleep(time.Second)

					model.LogError.InsertLog(
						"LT_LOCK_HASH_FAIL",
						dbg.Cat("Key :", owner_task.Key),
					)
				}

			} else {
				dt := model.GetGAS(db).GetLtLockPendingMin()
				if task_data.LockInfo.TryData.IsTimeOver(dt) { //lt_lock

					stx := task_data.LockInfo.TryData.STX(
						sender,
						owner.PrivateKey,
					)
					hash, err := sender.SendTransaction(
						context.Background(),
						stx,
					)
					if err == nil {
						task_data.LockInfo.Hash = hash
						owner_task.UpdateDB(db, task_data)
						TaskLog(TAG, "RE_SendTransaction :::", hash)
						continue
					}
				}

				time.Sleep(time.Second)
				TaskLog(TAG, "PendingWait :::", task_data.LockInfo.Hash)
			}

		default:
			TaskLog(TAG, owner_task.State)
			return

		} //switch
	} //for

}
