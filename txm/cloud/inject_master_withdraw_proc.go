package cloud

import (
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/mms"
)

type MASTER_WITHDRAW_FUNC func(db mongo.DATABASE, nowAt mms.MMS) bool

type _master_withdraw_data struct {
	pending_funcs []MASTER_WITHDRAW_FUNC
	try_funcs     []MASTER_WITHDRAW_FUNC
}

var (
	master_withdraw_func = _master_withdraw_data{}
)

func InjectMasterWithdrawProcess(
	pending_func MASTER_WITHDRAW_FUNC,
	try_func MASTER_WITHDRAW_FUNC,
) {
	master_withdraw_func.pending_funcs = append(master_withdraw_func.pending_funcs, pending_func)
	master_withdraw_func.try_funcs = append(master_withdraw_func.try_funcs, try_func)
}

//////////////////////////////////////////////////////////////////////////////
