package ecsx

import (
	"strings"
	"time"

	"txscheduler/brix/tools/mms"

	"txscheduler/brix/tools/cloudx/ethwallet"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
)

// NewTxBlockEx :
func NewTxBlockEx(t ethwallet.ResultTxData) TransactionBlock {
	ut := time.Unix(jmath.NewBigDecimal(t.TimeStamp).ToBigInteger().Int64(), 0)

	isError := false
	if t.IsError == "" {
		isError = false
	} else {
		isError = t.IsError != "0"
	}

	tx := TransactionBlock{
		IsContract:     false,
		BlockNumber:    t.BlockNumber,
		Confirmations:  t.Confirmations,
		Timestamp:      mms.FromTime(ut),
		IsError:        isError,
		ContractMethod: "",
		Hash:           t.Hash,
		From:           t.From,
		Amount:         t.Value,
	}
	tx.To = t.To

	if strings.HasPrefix(t.Input, "0x") == true { //call contract
		if t.To == "" && t.ContractAddress != "" {
			//계약 생성임
			tx.IsContract = true
			tx.ContractMethod = "deploy"
			tx.ContractAddress = t.ContractAddress
			//Input : 0x60806040526000...
			return tx
		}
		//Contract
		if t.Input != "0x" {
			tx.IsContract = true
			tx.ContractAddress = tx.To
			tx.To = ""
			if strings.HasPrefix(t.Input, MethodTransfer.MethodID) {
				tx.ContractMethod = MethodTransfer.FuncName
				txInput := t.Input[2:] // cut 0x
				toHex := "0x" + txInput[32:32+40]
				tx.To = toHex
				weiHex := "0x" + txInput[32+40:]
				tx.Amount = jmath.NewBigDecimal(weiHex).ToString()
			}
		}
	} else { //call address (t.Input == deprecated)

		if t.ContractAddress == "" {
			tx.Symbol = "ETH"
			tx.Decimals = "18"
		} else {
			tx.IsContract = true
			tx.ContractAddress = t.ContractAddress
			tx.Symbol = t.TokenSymbol
			tx.Decimals = t.TokenDecimal
		}
	}

	return tx
}

// GetTxByHash :
func (my TransactionBlockList) GetTxByHash(hash string) (TransactionBlock, bool) {
	hash = dbg.TrimToLower(hash)
	for _, tx := range my {
		tx.Hash = dbg.TrimToLower(tx.Hash)
		if tx.Hash == hash {
			return tx, true
		}
	}
	return TransactionBlock{}, false
}
