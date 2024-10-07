package ethwallet

import (
	"encoding/json"
	"strconv"
)

//ContractTransactional :
type ContractTransactional struct {
	Status  string          `json:"status"`  //  1, 0
	Message string          `json:"message"` // OK , No transactions found
	Txlist  TokenTxDataList `json:"result"`

	_isLimitTimeError bool
	_limitTImeMessage string

	_isRefactData    bool
	_lastBlockNumber int64
}

//String
func (my ContractTransactional) String() string {
	var m map[string]interface{}

	if my._isLimitTimeError == false {
		sm := map[string]interface{}{
			"status":          my.Status,
			"message":         my.Message,
			"tx-count":        len(my.Txlist),
			"lastBlockNumber": my._lastBlockNumber,
		}
		m = sm
	} else {
		m = map[string]interface{}{
			"status":    my.Status,
			"message":   my.Message,
			"limitTime": my._limitTImeMessage,
		}
	}

	b, _ := json.MarshalIndent(m, "", "  ")
	return string(b)
}

//IsSuccess :
func (my ContractTransactional) IsSuccess() bool {
	return my.Status == "1"
}

//NewContractTransactional :
func NewContractTransactional() *ContractTransactional {
	return &ContractTransactional{}
}

//GETPTR :
func (my *ContractTransactional) GETPTR() interface{} {
	if my == nil {
		return nil
	}
	return my
}

//LastBlockNumber :
func (my ContractTransactional) LastBlockNumber() int64 {
	return my._lastBlockNumber
}

//TxSort :
func (my *ContractTransactional) TxSort() {
	if my._isRefactData == true {
		return
	}
	my._isRefactData = true

	count := len(my.Txlist)
	if count == 0 {
		return
	} else {
		my.Txlist.txSort()
		my._lastBlockNumber, _ = strconv.ParseInt(my.Txlist[count-1].BlockNumber, 10, 64)
	}
}

//TxCount :
func (my ContractTransactional) TxCount() int {
	return len(my.Txlist)
}
