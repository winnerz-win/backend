package cloud

import (
	"time"
	"txscheduler/brix/tools/cloudx/ethwallet/EtherScanAPI"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
	"txscheduler/brix/tools/mms"
	"txscheduler/brix/tools/runtext"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func runETHInternal(rtx runtext.Runner, name, contract, startNumber string) {
	defer dbg.PrintForce("cloud.runETHInternal ----------  END")
	dbg.PrintForce("cloud.runETHInternal ----------  START")

	waitDuration := time.Minute * 2

	txCnt := model.NewETHTxInternalCnt(inf.Mainnet(), name, contract, startNumber)
	txCnt.LoadFromDB()
	startNumber = txCnt.Number

	tokenInfos := inf.Config().Tokens

EXIT:
	for {
		select {
		case <-rtx.EndC():
			break EXIT
		default:
		} //select

		list := ecsx.TransactionBlockList{}

		cfg := EtherScanAPI.NewConfig(inf.Mainnet(), inf.ESKey())

		lastNumber := startNumber
		txlist := ecsx.GetInternalTransaction(cfg, contract, startNumber)
		for _, tx := range txlist {
			if jmath.CMP(tx.BlockNumber, lastNumber) > 0 {
				lastNumber = tx.BlockNumber
			}
			if tx.IsError {
				continue
			}
			if tx.IsContract {
				info := tokenInfos.GetContract(tx.ContractAddress)
				if info.Valid() {
					item := ecsx.NewTxBlockInternal(tx)
					item.Symbol = info.Symbol
					item.Decimals = info.Decimal
					list = append(list, item)
				}
			} else {
				item := ecsx.NewTxBlockInternal(tx)
				item.Symbol = model.ETH
				item.ContractAddress = "eth"
				item.Decimals = "18"
				list = append(list, item)
			}
		} //for

		if len(list) > 0 {
			nowAt := mms.Now()
			processTxlist(list, nowAt, false)
		}

		txCnt.Update(lastNumber)
		startNumber = txCnt.Number

		time.Sleep(waitDuration)

	} //for

}
