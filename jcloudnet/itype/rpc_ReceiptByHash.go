package itype

import (
	"jtools/cloud/ebcm"
	"jtools/dbg"
	"jtools/jmath"
	"time"
)

func (my IClient) ReceiptByHash(hexHash string) ebcm.TxReceipt {

	req := ReqJsonRpc{
		Method: _nmap["getTransactionReceipt"][my.isKlay],
		Params: []any{
			hexHash,
		},
	}
	ack, err := req.Request(my.rpcURL, my.isDebug)
	if err != nil {
		return ebcm.TxReceipt{
			IsNotFound: true,
		}
	}

	re := ebcm.TxReceipt{}
	v, do := ack.Result.(map[string]interface{})
	if !do {
		re.IsNotFound = true
		return re
	}

	re.TransactionHash = _0xToLower(hexHash)
	re.BlockNumber = jmath.VALUE(v["blockNumber"])

	re.BlockHash = _0xToLower(v["blockHash"])

	if gas, do := v["gas"]; do {
		re.Gas = jmath.VALUE(gas)
	}

	if gasUsed, do := v["gasUsed"]; do {
		re.GasUsed = jmath.Uint64(gasUsed)
	} else {
		re.GasUsed = jmath.Uint64(re.Gas)
	}

	if gasPrice, do := v["gasPrice"]; do {
		re.GasPrice = jmath.VALUE(gasPrice)
	}

	if effectiveGasPrice, do := v["effectiveGasPrice"]; do {
		re.EffectiveGasPrice = jmath.VALUE(effectiveGasPrice)
	}

	re.From = _0xToLower(v["from"])
	re.Bloom = _0xToLower(v["logsBloom"])
	re.Nonce = jmath.Uint64(v["nonce"])

	re.SenderTxHash = _0xToLower(v["senderTxHash"])
	re.Status = jmath.Uint64(v["status"])

	toAddress := ""
	if to, do := v["to"]; do {
		if to != nil {
			switch v := to.(type) {
			case string:
				toAddress = v
			default:
				toAddress = _0xToLower(v)
			}
		}
	}
	re.To = dbg.TrimToLower(toAddress)
	tx_type := ebcm.TxType(uint16(jmath.Int(v["typeInt"])))
	re.Type = &tx_type

	re.Amount = jmath.VALUE(v["value"])

	if v["contractAddress"] != nil {
		if ca, do := v["contractAddress"].(string); do {
			re.ContractAddress = _0xToLower(ca)
		} else {
			re.ContractAddress = _0xToLower(v["contractAddress"])
		}
	}

	re.Logs = ebcm.TxLogList{}
	if log_field, do := v["logs"]; do {
		logs := log_field.([]interface{})

		for _, a := range logs {
			l := a.(map[string]interface{})
			log := ebcm.TxLog{
				Address:     _0xToLower(l["address"]),
				Data:        ebcm.MakeDataItemList(_0xToLower(l["data"])),
				BlockNumber: jmath.Uint64(l["blockNumber"]),
				BlockHash:   _0xToLower(l["blockHash"]),
				LogIndex:    uint(jmath.Int64(l["logIndex"])),
				TxHash:      _0xToLower(l["transactionHash"]),
				TxIndex:     uint(jmath.Int64(l["transactionIndex"])),
				Removed:     dbg.IsTrue(l["removed"]),
			}
			if t_field, do := l["topics"]; do {
				t_list := t_field.([]interface{})
				for _, topic := range t_list {
					c := _0xToLower(topic)
					log.Topics = append(log.Topics, ebcm.Topic(c))
				} //for
			}
			re.Logs = append(re.Logs, log)
		}

	}

	return re
}

func (my IClient) WaitConfirmTxHash(hexHash string, pending_wait_limit ...time.Duration) (ebcm.TxReceipt, bool) {

	pending_wait_time := time.Duration(time.Minute * 30)
	if len(pending_wait_limit) > 0 {
		if pending_wait_limit[0] > time.Second {
			pending_wait_time = pending_wait_limit[0]
		}
	}

	elipsed_time := time.Duration(0)
	for {
		time.Sleep(time.Second)

		r := my.ReceiptByHash(hexHash)

		if r.IsNotFound {
			elipsed_time += time.Second
			if elipsed_time >= pending_wait_time {
				break
			}
			continue
		}

		return r, true

	} //for

	return ebcm.TxReceipt{IsNotFound: true}, false
}
