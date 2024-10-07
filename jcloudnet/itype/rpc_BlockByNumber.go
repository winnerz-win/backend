package itype

import (
	"jtools/cloud/ebcm"
	"jtools/jmath"
	"jtools/unix"
)

func checkEIP1559(txType ebcm.TxType, callback func()) bool {
	if txType.Uint16() == ebcm.TXN_EIP_1559.Uint16() {
		callback()
		return true
	}
	return false
}

func (my IClient) _blockByNumber(number interface{}, is_full bool) *ebcm.BlockByNumberData {
	req := ReqJsonRpc{
		Method: _nmap["getBlockByNumber"][my.isKlay],
		Params: []any{
			param_hex_amount(number),
			is_full,
		},
	}
	ack, err := req.Request(my.rpcURL, my.isDebug)
	if err != nil {
		return nil
	}
	v, do := ack.Result.(map[string]interface{})
	if !do {
		return nil
	}

	data := ebcm.BlockByNumberData{}

	if !my.isKlay {
		data.BlockData = ebcm.BlockData{
			Number:       jmath.Uint64(v["number"]),
			NumberString: jmath.VALUE(v["number"]),
			Time:         unix.Time(jmath.Int64(v["timestamp"])),
			Hash:         _0xToLower(v["hash"]),
			PreHash:      _0xToLower(v["parentHash"]),
			CoinBase:     _0xToLower(v["miner"]),
			Difficulty:   _0xToLower(v["difficulty"]),
			GasLimit:     jmath.Uint64(v["gasLimit"]),

			GasUsed: jmath.Uint64(v["gasUsed"]),
			Nonce:   jmath.Uint64(v["nonce"]),

			Extra:       _0xToLower(v["extraData"]),
			ReceiptHash: _0xToLower(v["receiptsRoot"]),
			Root:        _0xToLower(v["stateRoot"]),
			Size:        jmath.VALUE(v["size"]),
			TxHash:      _0xToLower(v["transactionsRoot"]),

			BaseFee: jmath.VALUE(v["baseFeePerGas"]),
			//BlockScore: jmath.Int64(v["blockScore"]),
		}
	} else {
		data.BlockData = ebcm.BlockData{
			Number:       jmath.Uint64(v["number"]),
			NumberString: jmath.VALUE(v["number"]),
			Time:         unix.Time(jmath.Int64(v["timestamp"])),
			Hash:         _0xToLower(v["hash"]),
			PreHash:      _0xToLower(v["parentHash"]),
			RewardBase:   _0xToLower(v["reward"]),

			GasUsed: jmath.Uint64(v["gasUsed"]),

			Extra:       _0xToLower(v["extraData"]),
			ReceiptHash: _0xToLower(v["receiptsRoot"]),
			Root:        _0xToLower(v["stateRoot"]),
			Size:        jmath.VALUE(v["size"]),
			TxHash:      _0xToLower(v["transactionsRoot"]),

			BlockScore: jmath.Int64(v["blockScore"]),
		}
	}

	if c, do := v["transactions"]; do {
		list := c.([]interface{})
		data.TxCount = len(list)

		if is_full {
			for _, tv := range list {
				txv := tv.(map[string]interface{})

				tx_hash := ""
				// if !my.isKlay {
				// 	tx_hash = _0xToLower(txv["hash"])
				// } else {
				// 	tx_hash = _0xToLower(txv["senderTxHash"])
				// }
				tx_hash = _0xToLower(txv["hash"])

				txitem := my._makeTxData(
					tx_hash,
					txv,
				)

				txitem.Number = data.BlockData.Number
				txitem.BlockNumber = data.BlockData.NumberString
				txitem.Timestamp = data.Time

				if !my.isKlay {
					checkEIP1559(txitem.Type, func() {
						if txitem.GasTipCap != txitem.GasFeeCap {
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
					})
				}

				data.TxList = append(data.TxList, txitem)
			}
		} else {
			for _, hash := range list {
				data.TxHashList = append(data.TxHashList, jmath.HEX(hash))
			}
		}
	}

	return &data
}

func (my IClient) BlockByNumberSimple(number interface{}) *ebcm.BlockByNumberData {
	return my._blockByNumber(number, false)
}

func (my IClient) BlockByNumber(number interface{}) *ebcm.BlockByNumberData {
	return my._blockByNumber(number, true)
}
