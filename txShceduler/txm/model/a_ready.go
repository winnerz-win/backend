package model

import (
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/txm/inf"
	"txscheduler/txm/scv"
)

// const (
// 	ToMasterCallback  = "ToMasterCallback"
// 	ExMasterCallback  = "ExMasterCallback"
// 	MasterOutCallback = "MasterOutCallback"
// )

var (
	cbFlag = map[scv.HKey]bool{
		scv.ToMasterCallback:  false,
		scv.ExMasterCallback:  false,
		scv.MasterOutCallback: false,
	}
)

func MasterCallback(key scv.HKey) bool {
	if v, do := cbFlag[key]; do {
		return v
	}
	return false
}

func Ready(
	scv_callback_list scv.CallbackList,
	// isToMaster bool,
	// isExMaster bool,
	// isMasterOut bool,
) {
	dbg.PrintForce("model.Ready")

	for _, hkey := range scv_callback_list.HKeys() {
		if _, do := cbFlag[hkey]; do {
			cbFlag[hkey] = true
		}
	} //for

	// cbFlag[ToMasterCallback] = isToMaster
	// cbFlag[ExMasterCallback] = isExMaster
	// cbFlag[MasterOutCallback] = isMasterOut

	index_list := []mongo.IndexingDBFunc{
		AASystemInfo{},
		AA_GAS_PRICE{},

		TxETHBlock{},
		TxETHDeposit{},
		TxETHInternalCnt{},
		TxETHWithdraw{},
		TxETHCharger{},
		TxETHCounter{},

		Member{},
		InfoMaster{},
		InfoDeposit{},
		LogDeposit{},
		LogWithdraw{},
		LogToMaster{},
		LogExMaster{},
		ConSum{},
		CoinDay{},
		SyncCoin{},

		Admin{},

		TxETHMasterOut{},

		LogWithdrawSELF{},

		XLog{},

		CTX_USER{},
	}
	if inf.IsOnwerTaskMode() {
		index_list = append(index_list,
			OwnerTask{},
			OwnerLog{},
			OwnerUnlockPool{},
		)
	}
	mongo.StartIndexingDB(index_list...)

	inf.Config()

}
