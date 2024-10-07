package ecsx

import (
	"context"
	"fmt"
	"strings"

	"txscheduler/brix/tools/cloudx/ebcmx"
	"txscheduler/brix/tools/dbg"

	"txscheduler/brix/tools/jmath"

	"github.com/ethereum/go-ethereum/common"
)

func checkEIP1559(txType string, callback func()) {
	if txType == "2" {
		callback()
	}
}

// TransactionByHash : (tx, number, pending, err) ecsx.TransactionBlock, blockNumber, isPending, err  [2 call]
func (my *Sender) TransactionByHash(hashString string, baseFee ...string) (TransactionBlock, string, bool, error) {
	hashString = strings.TrimSpace(hashString)

	//ETHKJS
	from := ""
	blockNumber := "0"
	dbg.D(from, blockNumber)
	tx, isPending, err := my.client.TransactionByHash(
		context.Background(),
		common.HexToHash(hashString),
		// &from,
		// &blockNumber,
	)
	dbg.White("isPending :", isPending, ", err :", err)
	//dbg.Red(dbg.ToJSONString(tx))

	isLog := false
	// if dbg.BoolsOne(isLogs...) {
	// 	isLog = true
	// }
	if err != nil {
		if isLog {
			dbg.Red("txhash.error :", err)
		}
		return TransactionBlock{}, blockNumber, false, err
	}
	if tx == nil {
		errMsg := "tx is nil"
		if isLog {
			dbg.Red(errMsg)
		}
		return TransactionBlock{}, blockNumber, false, fmt.Errorf("%v", errMsg)
	}
	//dbg.Red(tx)
	txitem := NewTxBlock(my, tx, jmath.NewBigDecimal(blockNumber).ToString())
	my.checkCustomMethod(&txitem)
	if isLog {
		//dbg.Purple("tx-low :", tx)
		dbg.Purple("is_pending :", isPending)
		dbg.Purple("blockNumber :", jmath.NewBigDecimal(blockNumber).ToString())
		dbg.Purple(txitem)
	}

	if !isPending {
		txitem.TxBlockReceipt(my, true)
	} else {
	}

	checkEIP1559(txitem.Type, func() {
		if txitem.GasTipCap != txitem.GasFeeCap {
			if len(baseFee) == 0 {
				if data := my.getBlockDataByNumber(txitem.BlockNumber); data != nil {
					txitem.BaseFee = data.BaseFee
				}
			} else {
				txitem.BaseFee = baseFee[0]

			}
			if txitem.BaseFee != "" {
				txitem.Gas = jmath.ADD(txitem.GasTipCap, txitem.BaseFee)
				txitem.GasPriceETH = jmath.MUL(WeiToETH(txitem.Gas), txitem.GasUsed)
			}
		}
	})

	return txitem, jmath.NewBigDecimal(blockNumber).ToString(), isPending, nil
}

func (my *Sender) ebcm_TransactionByHash(hashString string) (ebcmx.TransactionBlock, bool, error) {
	tx, _, pending, err := my.TransactionByHash(hashString)

	item := ebcmx.TransactionBlock{}
	dbg.ChangeStruct(tx, &item)
	//item.Finder = my

	return item, pending, err
}

// TxBlockReceipt :
func (my *TransactionBlock) TxBlockReceipt(sender *Sender, isSkipEIP1559 ...bool) TransactionBlock {
	r := sender.ReceiptByHash(my.Hash)

	if !dbg.IsTrue2(isSkipEIP1559...) {
		checkEIP1559(my.Type, func() {
			if my.GasTipCap != my.GasFeeCap {
				if my.BaseFee == "" {
					if data := sender.getBlockDataByNumber(r.BlockNumber); data != nil {
						my.BaseFee = data.BaseFee
					}
				}
				if my.BaseFee != "" {
					my.Gas = jmath.ADD(my.GasTipCap, my.BaseFee)
				}
			}
		})
	}

	my.InjectReceipt(r)
	return *my
}

func (my *TransactionBlock) InjectReceipt(r TxReceipt) {
	my.BlockNumber = r.BlockNumber
	my.Number = jmath.Uint64(r.BlockNumber)
	my.IsError = r.Status != 1
	my.TxIndex = r.TransactionIndex

	my.CumulativeGasUsed = r.CumulativeGasUsed

	my.Logs = r.Logs
	my.GasUsed = jmath.VALUE(r.GasUsed)
	my.IsReceiptedByHash = true

	my.GasPriceETH = jmath.MUL(WeiToETH(my.Gas), my.GasUsed)
}

func (my *Sender) ebcm_InjectReceipt(tx *ebcmx.TransactionBlock, r ebcmx.TxReceipt) {
	tx.BlockNumber = r.BlockNumber
	tx.Number = jmath.Uint64(r.BlockNumber)
	tx.IsError = r.Status != 1
	tx.TxIndex = r.TransactionIndex

	tx.CumulativeGasUsed = r.CumulativeGasUsed

	tx.Logs = r.Logs
	tx.GasUsed = jmath.VALUE(r.GasUsed)
	tx.IsReceiptedByHash = true

	tx.GasPriceETH = jmath.MUL(WeiToETH(tx.Gas), tx.GasUsed)
}

// GetTransactionFee : 가스 수수료 (ETH) : gas_limit * gasUsed
func (my *TransactionBlock) GetTransactionFee() string {
	if !my.IsReceiptedByHash {
		if my.finder != nil {
			my.TxBlockReceipt(my.finder)
		}
	}
	//return jmath.MUL(my.Gas, my.GasUsed)
	if my.GasPriceETH == "" {
		my.GasPriceETH = jmath.MUL(WeiToETH(my.Gas), my.GasUsed)
	}
	return my.GasPriceETH
}
