package cloud

import (
	"jtools/mms"
	"time"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/runtext"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func runSELFWithdraw(rtx runtext.Runner) {
	defer dbg.PrintForce("cloud.runSELFWithdraw ---- END")
	<-rtx.WaitStart()
	dbg.PrintForce("cloud.runSELFWithdraw ---- START")

EXIT:
	for {
		select {
		case <-rtx.EndC():
			break EXIT
		default:
		} //select

		model.DB(func(db mongo.DATABASE) {
			selector := mongo.Bson{"state": model.TxStatePendingSELF}
			list := model.TxETHWithdrawList{}
			db.C(inf.TXETHWithdraw).Find(selector).Sort("timestamp").All(&list)

			if len(list) <= 0 {
				return
			}

			for _, item := range list {
				checker := get_sender_x()
				if checker == nil {
					model.LogError.WriteLog(
						db,
						model.ErrorFinderNull,
						"runSELFWithdraw.checker",
					)
					return
				}

				result, _, _ := checker.TransactionByHash(item.Hash)
				if !result.IsReceiptedByHash {
					logWithdraw("(", item.UID, ") SELF PendingWait(", item.Hash, ")")
					continue
				}

				ctx_address := ""
				item.Gas = result.GetTransactionFee()
				if result.IsError { //false
					// if fail_hash_continue_check(checker, item.Hash) { //TX 실패일경우 한번더 확인
					// 	continue
					// }

					model.LockMemberUID(db, item.UID, func(member model.Member) {
						logWithdraw("(", item.UID, ") SELF r-Fail ", item.ToPrice, item.Symbol)

						member.UpdateCoinDB_Legacy(db)

						selfLOG := item.MakeLogWithdrawSELF(member, model.TxStateFail)
						selfLOG.InsertDB(db, mms.Now())

						db.C(inf.TXETHWithdraw).Remove(item.Selector())

						ctx_address = member.Address
					})

				} else { //success
					model.LockMemberUID(db, item.UID, func(member model.Member) {
						logWithdraw("(", item.UID, ") SELF r-Sucess ", item.ToPrice, item.Symbol)
						member.Withdraw.ADD(item.Symbol, item.ToPrice)
						member.UpdateDB(db)

						member.UpdateCoinDB_Legacy(db)

						model.CoinDay{}.AddWithdraw(db, item.Symbol, item.ToPrice, item.Gas, mms.Now())

						selfLOG := item.MakeLogWithdrawSELF(member, model.TxStateSuccess)
						selfLOG.InsertDB(db, mms.Now())

						db.C(inf.TXETHWithdraw).Remove(item.Selector())

						ctx_address = member.Address

					})
				}
				model.UserTransactionEnd(db, ctx_address)
			} //for
		})

		time.Sleep(time.Second * 5)
	} //for
}
