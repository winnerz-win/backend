package itype

import (
	"jtools/cc"
	"jtools/cloud/ebcm"
	"jtools/cloud/ebcm/abi"
	"jtools/dbg"
	"jtools/jmath"
	"strings"
)

func (my IClient) _makeTxData(hashString string, v map[string]interface{}) ebcm.TransactionBlock {
	if my.isDebug {
		cc.Gray(dbg.ToJsonString(v))
	}

	tx := ebcm.TransactionBlock{}

	tx.Hash = dbg.TrimToLower(hashString)
	tx.Number = jmath.Uint64(v["blockNumber"])
	tx.BlockNumber = jmath.VALUE(tx.Number)

	tx.From = ADDRESS(v["from"])
	tx.To = ADDRESS(v["to"])

	tx.Gas = jmath.VALUE(v["gas"])
	tx.Cost = tx.Gas
	tx.Limit = uint64(jmath.Int64(tx.Gas))
	tx.GasUsed = tx.Gas

	tx.Nonce = jmath.Uint64(v["nonce"])
	tx.GasPrice = jmath.VALUE(v["gasPrice"])
	tx.CustomInput = jmath.HEX(v["input"])

	tx.TxIndex = uint(jmath.Uint64(v["transactionIndex"]))

	tx.Amount = jmath.VALUE(v["value"])

	if my.isKlay {
		tx.Type = ebcm.TxType(jmath.Int(v["typeInt"]))

		tx.GasUsed = tx.Gas
		tx.TxFeeKLAY = jmath.MUL(ebcm.WeiToETH(tx.GasPrice), tx.GasUsed)

		tx.TypeString = dbg.Cat(v["type"])
		if fee_payer, do := v["feePayer"]; !do {
			tx.FeePayer = "0x0000000000000000000000000000000000000000"
		} else {
			tx.FeePayer = ADDRESS(fee_payer)
		}
		switch tx.Type {
		case 49: //TxTypeFeeDelegatedSmartContractExecution
			ratio := ebcm.Klay_FeeRatio(100)
			tx.FeeRatio = &ratio

		case 30722: //TxTypeEthereumDynamicFee
			if maxPriorityFeePerGas, do := v["maxPriorityFeePerGas"]; do {
				tx.GasTipCap = jmath.VALUE(maxPriorityFeePerGas)
			}
			if maxFeePerGas, do := v["maxFeePerGas"]; do {
				tx.GasFeeCap = jmath.VALUE(maxFeePerGas)
			}
		}

	} else {
		tx.Type = ebcm.TxType(jmath.Int(v["type"]))
		if tx.Type == ebcm.TxType(ebcm.TXN_EIP_1559) {
			tx.GasTipCap = jmath.VALUE(v["maxPriorityFeePerGas"])
			tx.GasFeeCap = jmath.VALUE(v["maxFeePerGas"])
		}
	}

	tx.Logs = ebcm.TxLogList{}

	if !strings.HasPrefix(tx.To, "0x") {
		tx.ContractMethod = "deploy"
		tx.ContractAddress = abi.ContractAddressNonce(tx.From, tx.Nonce)
	}

	ebcm.CheckMethodERC20(
		abi.GetInputParser(),
		tx.CustomInput,
		&tx,
	)

	return tx
}

func (my IClient) TransactionByHash(hashString string) (ebcm.TransactionBlock, bool, error) {

	req := ReqJsonRpc{
		Method: _nmap["getTransactionByHash"][my.isKlay],
		Params: []any{
			hashString,
		},
	}
	ack, err := req.Request(my.rpcURL, my.isDebug)
	if err != nil {
		return ebcm.TransactionBlock{}, false, err
	}

	v, do := ack.Result.(map[string]interface{})
	if !do {
		return ebcm.TransactionBlock{}, true, nil
	}

	txitem := my._makeTxData(hashString, v)

	receipt := my.ReceiptByHash(hashString)

	if !receipt.IsNotFound {
		my.InjectReceipt(&txitem, receipt)

		if !my.isKlay {
			if !checkEIP1559(txitem.Type, func() {
				if txitem.GasTipCap != txitem.GasFeeCap {
					if block := my.BlockByNumberSimple(txitem.BlockNumber); block != nil {
						txitem.BaseFee = block.BaseFee
						txitem.Timestamp = block.Time
					}

					if txitem.BaseFee != "" {
						/*
							gas = block.Base + tx.GasTipCap(Max Priority)
							gas 가 FeeCap (MAX) 보다 크면 FeeCap이 적용됨.
						*/
						gas := jmath.ADD(txitem.GasTipCap, txitem.BaseFee)
						txitem.Gas = gas
						if jmath.CMP(gas, txitem.GasFeeCap) >= 0 {
							txitem.Gas = txitem.GasFeeCap
						}
						txitem.TxFeeETH = jmath.MUL(ebcm.WeiToETH(txitem.Gas), txitem.GasUsed)
					}
				} else {
					if block := my.BlockByNumberSimple(txitem.BlockNumber); block == nil {
						txitem.BaseFee = block.BaseFee
						txitem.Timestamp = block.Time
					}
				}
			}) {
				if block := my.BlockByNumberSimple(txitem.BlockNumber); block == nil {
					txitem.Timestamp = block.Time
				}
			}
		} else {

		}
	}

	return txitem, receipt.IsNotFound, nil
}
