package etherscan

import (
	"encoding/json"
	"jtools/jmath"
	"sort"
	"strconv"
	"strings"
)

// TokenTxData :
type TokenTxData struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	From              string `json:"from"`
	ContractAddress   string `json:"contractAddress"`
	To                string `json:"to"`
	Value             string `json:"value"`
	TokenName         string `json:"tokenName"`
	TokenSymbol       string `json:"tokenSymbol"`
	TokenDecimal      string `json:"tokenDecimal"`
	TransactionIndex  string `json:"transactionIndex"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	GasUsed           string `json:"gasUsed"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	Input             string `json:"input"`
	Confirmations     string `json:"confirmations"`

	_isRefresh bool
	_timestamp int64
}

func (my *TokenTxData) refresh() {
	if my._isRefresh {
		return
	}
	my._isRefresh = true
	my.Hash = strings.ToLower(my.Hash)
	my.ContractAddress = strings.ToLower(my.ContractAddress)
	my.From = strings.ToLower(my.From)
	my.To = strings.ToLower(my.To)

	tn, _ := strconv.ParseInt(my.TimeStamp, 10, 64)
	my._timestamp = tn
}

// TimeStampInt64 :
func (my TokenTxData) TimeStampInt64() int64 {
	return my._timestamp
}

// ToString :
func (my TokenTxData) ToString() string {
	b, _ := json.MarshalIndent(my, "", "  ")
	return string(b)
}

// TokenTxDataList :
type TokenTxDataList []TokenTxData

func (my TokenTxDataList) txSort() {
	for i := 0; i < len(my); i++ {
		my[i].refresh()
	} //for
	sort.Slice(my, func(i, j int) bool {
		return jmath.CMP(my[i].BlockNumber, my[j].BlockNumber) < 0
	})
}

// ToString :
func (my TokenTxDataList) ToString() string {
	b, _ := json.MarshalIndent(my, "", "  ")
	return string(b)
}

// ValidConfirms :
func (my TokenTxDataList) ValidConfirms(hash string, count interface{}) bool {
	for _, tx := range my {
		if tx.Hash == hash {
			if jmath.CMP(tx.Confirmations, count) >= 0 {
				return true
			}
		}
	} //for
	return false
}
