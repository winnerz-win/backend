package cloud

import (
	"jtools/cloud/ebcm"
	"jtools/mms"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/txm/model"
)

type MASTER_WITHDRAW_FUNC func(db mongo.DATABASE, nowAt mms.MMS) bool

type MasterTxJobData struct {
	tag          string
	pending_func MASTER_WITHDRAW_FUNC
	try_func     MASTER_WITHDRAW_FUNC
}

type _master_withdraw_data struct {
	jobs []MasterTxJobData
	// pending_funcs []MASTER_WITHDRAW_FUNC
	// try_funcs     []MASTER_WITHDRAW_FUNC
}

var (
	master_withdraw_func = _master_withdraw_data{}
)

func InjectMasterWithdrawProcess(
	tag string,
	pending_func MASTER_WITHDRAW_FUNC,
	try_func MASTER_WITHDRAW_FUNC,
) {
	job := MasterTxJobData{
		tag:          tag,
		pending_func: pending_func,
		try_func:     try_func,
	}
	master_withdraw_func.jobs = append(master_withdraw_func.jobs, job)
}

//////////////////////////////////////////////////////////////////////////////

func CALC_GAS_PRICE(db mongo.DATABASE, gas_price ebcm.GasPrice) ebcm.GasPrice {
	return model.CALC_GAS_PRICE(db, gas_price)
	// mul_value := model.GetGAS(db).GetGasMulValue()
	// gas_price.MUL_VALUE(mul_value)
	// return gas_price
}
