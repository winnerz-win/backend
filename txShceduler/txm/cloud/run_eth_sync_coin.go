package cloud

import (
	"strings"
	"time"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/mms"
	"txscheduler/brix/tools/runtext"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func runSyncCoin(rtx runtext.Runner) {
	defer dbg.PrintForce("cloud.runSyncCoin ----------  END")
	<-rtx.WaitStart()
	dbg.PrintForce("cloud.runSyncCoin ----------  START")

EXIT:
	for {
		select {
		case <-rtx.EndC():
			break EXIT
		default:
		}

		model.DB(func(db mongo.DATABASE) {
			nowAt := mms.Now()

			selector := mongo.Bson{
				"at": mongo.Bson{"$lte": nowAt},
			}
			finder := db.C(inf.COLSyncCoin).Find(selector)
			cnt, _ := finder.Count()
			if cnt == 0 {
				return
			}

			tokenmap := map[string]*ecsx.Token{}
			tokenlist := inf.TokenList()
			for _, c := range tokenlist {
				if !strings.HasPrefix(c.Contract, "0x") {
					continue
				}

				checker := get_sender_x()
				if checker == nil {
					model.LogError.WriteLog(
						db,
						model.ErrorFinderNull,
						"runSyncCoin.checker.1",
					)
					return
				}
				token := checker.Token(c.Contract)
				if !token.Valid() {
					continue
				}
				tokenmap[c.Symbol] = token
			} //for

			iter := finder.Iter()

			sc := model.SyncCoin{}
			for iter.Next(&sc) {
				checker := get_sender_x()
				if checker == nil {
					model.LogError.WriteLog(
						db,
						model.ErrorFinderNull,
						"runSyncCoin.checker.2",
					)
					return
				}
				model.LockMember(db, sc.Address, func(member model.Member) {
					member.UpdateCoinDB_Legacy(db, checker)
				})

				sc.RemoveDB(db)
			} //for

		})

		time.Sleep(time.Second * 5)
	} //for
}
