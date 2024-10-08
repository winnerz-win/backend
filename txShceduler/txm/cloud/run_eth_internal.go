package cloud

import (
	"jcloudnet/etherscan/etherscanapi"
	"jtools/cloud/ebcm"
	"jtools/jmath"
	"jtools/mms"
	"time"
	"txscheduler/brix/tools/dbg"
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

		list := ebcm.TransactionBlockList{}

		cfg := etherscanapi.NewConfig(inf.Mainnet(), inf.ESKey())

		lastNumber := startNumber
		txlist := etherscanapi.GetInternalTransactionAPI(cfg, contract, startNumber)
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
					item := NewTxBlockInternal(tx)
					item.Symbol = info.Symbol
					item.Decimals = info.Decimal
					list = append(list, item)
				}
			} else {
				item := NewTxBlockInternal(tx)
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

// NewTxBlockInternal :
func NewTxBlockInternal(tx etherscanapi.InternalTransaction) ebcm.TransactionBlock {
	item := ebcm.TransactionBlock{
		IsInternal:      true,
		IsContract:      tx.IsContract,
		ContractAddress: tx.ContractAddress,
		Hash:            tx.Hash,
		From:            tx.From,
		To:              tx.To,
		Amount:          tx.Value,
		CustomInput:     tx.Input,
		IsError:         tx.IsError,
		Type:            ebcm.TxType(jmath.Int(tx.Type)),
		Gas:             tx.Gas,
		GasUsed:         tx.GasUsed,
		TradeID:         tx.TradeID,
		ErrCode:         tx.ErrCode,
	}
	return item
}
