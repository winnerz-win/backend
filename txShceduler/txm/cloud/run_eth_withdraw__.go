package cloud

import (
	"context"
	"jtools/cc"
	"jtools/mms"
	"jtools/unix"
	"time"
	"txscheduler/brix/tools/console"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/runtext"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

var (
	master_tick_pause = false
)

func _cloud_master_console() {
	console.AppendCmd(
		"master_p",
		"",
		true,
		func(ps []string) {
			master_tick_pause = !master_tick_pause
			cc.YellowItalic("master_tick_pause :", master_tick_pause)
		},
	)

	type info_data struct {
		A_Nonce   uint64 `json:"a_nonce"`
		B_Pending uint64 `json:"b_pending"`
		Coin      string `json:"coin"`
		Token     string `json:"token"`
	}
	get_info_data := func(address string) info_data {
		finder := inf.GetFinder()
		wnz := inf.FirstERC20()
		ctx := context.Background()

		n, _ := finder.NonceAt(ctx, address)
		p, _ := finder.PendingNonceAt(ctx, address)
		a := finder.Price(address, "ETH", 18)
		b := finder.Price(address, wnz.Contract, wnz.Decimal)
		return info_data{
			A_Nonce:   n,
			B_Pending: p,
			Coin:      a,
			Token:     b,
		}
	}

	console.AppendCmd(
		"nv",
		"",
		true,
		func(ps []string) {
			{
				master := inf.Master()
				data := get_info_data(master.Address)
				cc.YellowItalic("MASTER(", master.Address, ") nonce(", data.A_Nonce, ")/pending(", data.B_Pending, ")")
				cc.YellowItalic("coin :", data.Coin)
				cc.YellowItalic("wnz  :", data.Token)
			}
			cc.White("------------------------------------------------")
			{
				owner := inf.Owner()
				data := get_info_data(owner.Address)
				cc.YellowItalic("OWNER(", owner.Address, ") nonce(", data.A_Nonce, ")/pending(", data.B_Pending, ")")
				cc.YellowItalic("coin :", data.Coin)
				cc.YellowItalic("wnz  :", data.Token)
			}
		},
	)
	go func() {
		model.DB(func(db mongo.DATABASE) {
			c := db.C("zz_nv")
			c.EnsureIndex(mongo.SingleIndex("key", 1, false))
		})
		for {
			model.DB(func(db mongo.DATABASE) {
				c := db.C("zz_nv")

				cnt, _ := c.Find(mongo.Bson{"key": mongo.Bson{"$exists": true}}).Count()
				if cnt <= 0 {
					return
				}
				c.RemoveAll(nil)

				data := mongo.MAP{
					"info":   dbg.Cat("key_", unix.Now().KST()),
					"owner":  get_info_data(inf.Owner().Address),
					"master": get_info_data(inf.Master().Address),
				}
				c.Insert(data)

			})
			time.Sleep(time.Second)
		} //for
	}()
}

func runETHWithdraw(rtx runtext.Runner, is_owner_master_same bool) {
	defer dbg.PrintForce("cloud.runETHWithdraw ---- END")
	<-rtx.WaitStart()
	dbg.PrintForce("cloud.runETHWithdraw ---- START")

	if inf.IsOnwerTaskMode() {
		InjectMasterWithdrawProcess(
			"<MASTER_LT_JOB>",
			owner_lt_master_pending,
			owner_lt_master_transfer,
		)
	}

	InjectMasterWithdrawProcess(
		"<MASTER_OUT_JOB>",
		func(db mongo.DATABASE, nowAt mms.MMS) bool {
			return procMasterOutPending(db, nowAt)
		},
		func(db mongo.DATABASE, nowAt mms.MMS) bool {
			return procMasterOutTry(db, nowAt)
		},
	)

	InjectMasterWithdrawProcess(
		"<MASTER_WITHDRAW_JOB>",
		func(db mongo.DATABASE, nowAt mms.MMS) bool {
			return procETHWithdrawPending(db, nowAt)
		},
		func(db mongo.DATABASE, nowAt mms.MMS) bool {
			return procETHWithdarwTry(db, nowAt) //가상계좌 잔액 회수
		},
	)

	if is_owner_master_same {
		go runOwnerTaskLogSender()
	}

	pause_du := time.Duration(0)
EXIT:
	for {
		select {
		case <-rtx.EndC():
			break EXIT
		default:
		} //select

		sleep_du := time.Millisecond * 100
		time.Sleep(sleep_du)

		if master_tick_pause {
			pause_du += sleep_du
			if pause_du > time.Second*3 {
				cc.Gray("MASTER_TICK_PAUSE")
				pause_du = 1
			}
			continue
		} else {
			if pause_du > 0 {
				pause_du = 0
				cc.Gray("MASTER_TICK_RESUME")
			}
		}

		model.DB(func(db mongo.DATABASE) {

			if is_owner_master_same {
				if owner_job_action(db) {
					return
				}
			}

			for _, job := range master_withdraw_func.jobs {
				if job.pending_func(db, mms.Now()) {
					cc.Yellow(job.tag, "pending_func .....")
					return
				}
			}

			for _, job := range master_withdraw_func.jobs {
				if job.try_func(db, mms.Now()) {
					cc.Yellow(job.tag, "try_func .....")
					return
				}
			}

		})
	} //for
}

const (
	MasterFail_InvalidSymbol         = "invalid_symbol"
	MasterFail_NeedPrice             = "need_price"
	MasterFail_ChainErrorBox         = "chain_error:box"
	MasterFail_ChainErrorNonce       = "chain_error:nonce"
	MasterFail_ChainErrorTx          = "chain_error:tx"
	MasterFail_ChainErrorSend        = "chain_error:send"
	MasterFail_ChianErrorPendingTime = "chain_error:pending_time"
	MasterFail_ChianErrorFail        = "chain_error:fail"
)

var (
	pendingWaitMinute float64 = 5
)

///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
