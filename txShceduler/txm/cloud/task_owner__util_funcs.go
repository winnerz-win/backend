package cloud

import (
	"jtools/cc"
	"jtools/cloud/ebcm"
	"jtools/cloud/jeth/ecs"
	"jtools/jmath"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/txm/ack"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func CheckLockState(address string) (bool, chttp.CError) {
	erc20_info := inf.FirstERC20()
	caller := ecs.New(
		ecs.RPC_URL(inf.Mainnet()),
		inf.InfuraKey(),
	)
	r, err := model.LockTokenUtil{}.TotalLocked(
		caller, erc20_info.Contract, address,
	)
	if err != nil {
		cc.RedItalic("[check_lock_state] :", err)
		return false, ack.DBJob
	}
	is_locked := jmath.CMP(r.LockedAmount, 0) > 0

	return is_locked, nil
}

func makeStxTransaction(
	sender *ebcm.Sender,

	nonce uint64,
	to string,
	amount string,

	limit uint64,
	gas_price ebcm.GasPrice,
	data ebcm.PADBYTES,

	private_key string,

) ebcm.WrappedTransaction {
	stx, _ := sender.SignTx(
		sender.NewTransaction(
			nonce,
			to,
			amount,
			limit,
			gas_price,
			data,
		),
		private_key,
	)
	return stx
}
