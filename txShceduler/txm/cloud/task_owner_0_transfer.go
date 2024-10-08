package cloud

import (
	"context"
	"jtools/cloud/ebcm"
	"jtools/cloud/ebcm/abi"
	"jtools/jmath"
	"jtools/unix"
	"time"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func owner_task_transfer(db mongo.DATABASE, owner_task model.OwnerTask) {
	task_data := owner_task.OwnerTransferTask()

	token_info := inf.FirstERC20()
	owner := inf.Owner()

	sender := inf.GetSender()

	TAG := dbg.Cat("[TASK_TRANSFER](", owner_task.Key, ")")
	TaskLog(TAG, "-------------- START")
	defer TaskLog(TAG, "-------------- END")

	for {
		try_data := task_data.CurrentData()
		if try_data == nil {
			if task_data.CheckStateEnd() {
				owner_task.UpdateDB(db, task_data)
			}
			TaskLog(TAG, "try_data is nil.")
			return
		}

		SEQ := TAG + dbg.Cat("[", task_data.Seq, "] ")

		switch try_data.State {
		case model.NONE:
			owner_coin_value := sender.Balance(owner.Address)
			if jmath.CMP(owner_coin_value, 0) <= 0 {
				//show me the coin!
				TaskError(SEQ, "Owner Coin Balance is 0.")
				return
			}

			owner_token_value, err := model.Erc20Balance(
				sender,
				token_info.Contract,
				owner.Address,
			)
			if err != nil {
				TaskError(SEQ, "Erc20Balance :", err)
				return
			}
			if jmath.CMP(owner_token_value, try_data.Amount) < 0 {
				//show me the token!
				TaskError(
					SEQ,
					dbg.Cat(
						"Owner Token Balance (",
						ebcm.WeiToToken(owner_token_value, token_info.Decimal),
						"/",
						ebcm.WeiToToken(try_data.Amount, token_info.Decimal),
						")",
					),
				)
				return
			}

			nonce, err := _getNonce(sender, owner.Address)
			if err != nil {
				TaskError(SEQ, "_getNonce :", err)
				return
			}

			var pad_bytes ebcm.PADBYTES
			if try_data.ReleaseTime > unix.ZERO {
				pad_bytes = ebcm.MakePadBytesABI(
					"transferWithLock",
					abi.TypeList{
						abi.NewAddress(task_data.Recipient),
						abi.NewUint256(try_data.Amount),
						abi.NewUint256(try_data.ReleaseTime),
					},
				)
			} else {
				pad_bytes = ebcm.MakePadBytesABI(
					"transfer",
					abi.TypeList{
						abi.NewAddress(task_data.Recipient),
						abi.NewUint256(try_data.Amount),
					},
				)
			}

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
				TaskError(SEQ, "EstimateGas :", err)
				return
			}

			limit = ebcm.MMA_LimitBuffer_MasterOut(limit)

			gas_price, err := sender.SuggestGasPrice(context.Background())
			if err != nil {
				TaskError(SEQ, "SuggestGasPrice :", err)
				return
			}
			gas_price = CALC_GAS_PRICE(db, gas_price)

			stx := makeStxTransaction(
				sender,

				nonce,
				token_info.Contract,
				"0",

				limit,
				gas_price,
				pad_bytes,

				owner.PrivateKey,
			)

			if err != nil {
				TaskError(SEQ, "SignTx :", err)
				return
			}

			hash, err := sender.SendTransaction(
				context.Background(),
				stx,
			)
			if err != nil {
				TaskError(SEQ, "SendTransaction :", err)
				return
			}

			try_data.Hash = hash
			try_data.State = model.PENDING

			owner_task.State = model.PENDING
			owner_task.UpdateDB(db, task_data)

			time.Sleep(time.Second)
			TaskLog(SEQ, "SendTransaction :::", hash)

		case model.PENDING:

			receipt := sender.ReceiptByHash(try_data.Hash)
			if !receipt.IsNotFound {
				if !receipt.IsError() {
					try_data.State = model.SUCCESS
					try_data.TxFeeETH = receipt.TxFeeETH()
					task_data.Seq++

					isTaskEnd := task_data.CheckStateEnd()
					owner_task.UpdateDB(db, task_data)

					if isTaskEnd {
						log := task_data.CbOwnerTransferLog(owner_task.Key)
						model.OwnerLog{}.InsertDB(db, log)
						return
					}
				} else {
					try_data.State = model.NONE

					owner_task.UpdateDB(db, task_data)
					TaskLog(SEQ, "PendingWait :::", try_data.Hash)
				}

			} else {
				time.Sleep(time.Second)
				TaskLog(SEQ, "PendingWait :::", try_data.Hash)
			}

		} //switch
	} //for

}
