package itype

import (
	"jtools/cloud/ebcm"
	"jtools/jmath"
)

func (my IClient) InjectReceipt(tx *ebcm.TransactionBlock, r ebcm.TxReceipt) {
	tx.IsReceptFailByTxInject = !r.Valid()

	tx.BlockNumber = r.BlockNumber
	tx.Number = jmath.Uint64(r.BlockNumber)

	if !r.IsNotFound {
		tx.IsError = r.Status != 1
		tx.IsReceiptedByHash = true
	}

	//tx.TxIndex = r.TransactionIndex
	tx.Logs = r.Logs
	tx.GasUsed = jmath.VALUE(r.GasUsed)

	tx.EffectiveGasPrice = r.EffectiveGasPrice

	if my.isKlay {
		if tx.EffectiveGasPrice == "" {
			tx.TxFeeKLAY = jmath.MUL(ebcm.WeiToETH(tx.GasPrice), tx.GasUsed)
		} else {
			tx.TxFeeKLAY = jmath.MUL(ebcm.WeiToETH(tx.EffectiveGasPrice), tx.GasUsed)
		}

	} else {
		tx.TxFeeETH = jmath.MUL(ebcm.WeiToETH(tx.GasPrice), tx.GasUsed)
	}

}
