package ethwallet

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"

	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
	"txscheduler/brix/tools/mms"
)

//ResultTxData :
type ResultTxData struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	From              string `json:"from"`
	ContractAddress   string `json:"contractAddress"` //ETH is ""
	To                string `json:"to"`
	Value             string `json:"value"`
	TokenName         string `json:"tokenName,omitempty"`    //ETH is ""	- tokentx
	TokenSymbol       string `json:"tokenSymbol,omitempty"`  //ETH is ""	- tokentx
	TokenDecimal      string `json:"tokenDecimal,omitempty"` //ETH is ""	- tokentx
	TransactionIndex  string `json:"transactionIndex"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	GasUsed           string `json:"gasUsed"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	Input             string `json:"input"`
	Confirmations     string `json:"confirmations"`

	IsError         string `json:"isError"`          //ETH only : "0"
	TxReceiptStatus string `json:"txreceipt_status"` //ETH & TOKEN : "1"
	IsContract      bool   `json:"is_contract"`

	_isRefresh    bool
	_blockNumber  int64
	_txIndex      int64
	_confirmCount int64
	_isToken      bool
}

//ResultTxDatas :
type ResultTxDatas []ResultTxData

func (my ResultTxData) String() string  { return dbg.ToJSONString(my) }
func (my ResultTxDatas) String() string { return dbg.ToJSONString(my) }

func (my ResultTxData) TimeMMS() mms.MMS {
	return mms.MMS(jmath.Int64(my.TimeStamp) * 1000)
}

//IsReceiptError : my.TxReceiptStatus != "1"
func (my ResultTxData) IsReceiptError() bool { return my.TxReceiptStatus != "1" }

//IsTxSuccess : ETHER TX가 정상 처리되었는지 확인.
func (my *ResultTxData) IsTxSuccess() bool {
	//my.Refresh()
	if my._isToken == true {
		return my.TxReceiptStatus == "1"
	}
	if my.IsError != "0" || my.TxReceiptStatus != "1" {
		return false
	}
	return true
}

//IsTxError :
func (my *ResultTxData) IsTxError() bool {
	if my.IsTxSuccess() {
		return false
	}
	return true
}

//GetBlockNumber :
func (my *ResultTxData) GetBlockNumber() int64 {
	return my._blockNumber
}

//GetConfirmCount :
func (my *ResultTxData) GetConfirmCount() int64 {
	return my._confirmCount
}

//ToString :
func (my ResultTxData) ToString() string {
	b, _ := json.MarshalIndent(my, "", "  ")
	return string(b)
}

//Refresh : TxData가 Token인지 Ehter인지 확인..
func (my *ResultTxData) Refresh() {
	if my._isRefresh {
		return
	}
	//dbg.Red(dbg.ToJSONString(my))

	my._isRefresh = true
	my.Hash = strings.ToLower(my.Hash)
	my.ContractAddress = strings.ToLower(my.ContractAddress)
	my.From = strings.ToLower(my.From)
	my.To = strings.ToLower(my.To)

	no, _ := strconv.ParseInt(my.BlockNumber, 10, 64)
	my._blockNumber = no

	txindex, _ := strconv.ParseInt(my.TransactionIndex, 10, 64)
	my._txIndex = txindex

	cn, _ := strconv.ParseInt(my.Confirmations, 10, 64)
	my._confirmCount = cn

	if my.ContractAddress == "" {
		my._isToken = false
		my.TokenDecimal = "18"
	} else {
		my._isToken = true
		_ = my.Input
	}
}

//IsTokenTransfer : token - transfer 함수 인지 체크
func (my *ResultTxData) IsTokenTransfer() bool {
	my.Refresh()
	if my._isToken == false {
		return false
	}
	name := my.Input[:34]
	if name == "0xa9059cbb000000000000000000000000" {
		return true
	}
	return false
}

//TokenTo : token to user
func (my *ResultTxData) TokenTo() string {
	my.Refresh()
	if my._isToken == false {
		return ""
	}
	buffer := my.Input[34:]
	buffer = buffer[:40]
	return "0x" + strings.ToLower(buffer)
}

//TokenValue : token transfer value(wei)
func (my *ResultTxData) TokenValue() string {
	my.Refresh()
	if my._isToken == false {
		return "0"
	}
	buffer := "0x" + my.Input[34+40+46:]
	return jmath.NewBigDecimal(buffer).ToString()
}

//ToString :
func (my ResultTxDatas) ToString() string {
	b, _ := json.MarshalIndent(my, "", "  ")
	return string(b)
}

func (my ResultTxDatas) rtxSort(sortOrder ...string) int64 {
	var lastBlockNumber int64 = 0
	for i := 0; i < len(my); i++ {
		my[i].Refresh()
		if lastBlockNumber < my[i]._blockNumber {
			lastBlockNumber = my[i]._blockNumber
		}
	} //for

	orderString := "asc"
	if len(sortOrder) > 0 {
		orderString = sortOrder[0]
	}
	if orderString == "asc" {
		sort.Slice(my, func(i, j int) bool {
			if my[i]._blockNumber < my[j]._blockNumber {
				return true
			} else if my[i]._blockNumber == my[j]._blockNumber {
				return my[i]._txIndex < my[j]._txIndex
			} else {
				return false
			}
		})
	} else {
		sort.Slice(my, func(i, j int) bool {
			if my[i]._blockNumber < my[j]._blockNumber {
				return false
			} else if my[i]._blockNumber == my[j]._blockNumber {
				return my[i]._txIndex > my[j]._txIndex
			} else {
				return true
			}
		})
	}

	return lastBlockNumber
}
