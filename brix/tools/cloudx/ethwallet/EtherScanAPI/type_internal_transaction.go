package EtherScanAPI

import (
	"strconv"

	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/mms"
)

//InternalTransaction :
type InternalTransaction struct {
	BlockNumber     string  `json:"blockNumber"`
	Timestamp       mms.MMS `json:"timeStamp"`
	Hash            string  `json:"hash"`
	From            string  `json:"from"`
	To              string  `json:"to"`
	Value           string  `json:"value"`
	IsContract      bool    `json:"is_contract"`
	ContractAddress string  `json:"contractAddress,omitempty"`
	Input           string  `json:"input"`
	Type            string  `json:"type"`
	Gas             string  `json:"gas"`
	GasUsed         string  `json:"gasUsed"`
	TradeID         string  `json:"traceId"`
	IsError         bool    `json:"isError"`
	ErrCode         string  `json:"errCode"`
	TimeKST         string  `json:"time_kst"`
}

//InternalTransactionList :
type InternalTransactionList []InternalTransaction

type internalTx struct {
	BlockNumber     string `json:"blockNumber"`
	Timestamp       string `json:"timeStamp"`
	Hash            string `json:"hash"`
	From            string `json:"from"`
	To              string `json:"to"`
	Value           string `json:"value"`
	ContractAddress string `json:"contractAddress"`
	Input           string `json:"input"`
	Type            string `json:"type"`
	Gas             string `json:"gas"`
	GasUsed         string `json:"gasUsed"`
	TradeID         string `json:"traceId"`
	IsError         string `json:"isError"`
	ErrCode         string `json:"errCode"`
}

//GetTimestamp :
func (my internalTx) getTimestamp() mms.MMS {
	v, _ := strconv.ParseInt(my.Timestamp, 10, 64)
	return mms.MMS(v * 1000)
}
func (my internalTx) Tx() InternalTransaction {
	tx := InternalTransaction{
		BlockNumber: my.BlockNumber,
		Timestamp:   my.getTimestamp(),
		Hash:        dbg.TrimToLower(my.Hash),
		From:        dbg.TrimToLower(my.From),
		To:          dbg.TrimToLower(my.To),
		Value:       my.Value,
		Input:       my.Input,
		Type:        my.Type,
		Gas:         my.Gas,
		GasUsed:     my.GasUsed,
		TradeID:     my.TradeID,
		IsError:     my.IsError != "0",
		ErrCode:     my.ErrCode,
	}
	contract := dbg.TrimToLower(my.ContractAddress)
	if contract != "" {
		tx.IsContract = true
		tx.ContractAddress = contract
	}
	tx.TimeKST = tx.Timestamp.KST()

	return tx
}

type internalTxList []internalTx

func (my internalTxList) ToList() InternalTransactionList {
	list := InternalTransactionList{}
	for _, v := range my {
		list = append(list, v.Tx())
	} //for
	return list
}

//InternalTxData :
type InternalTxData struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Result  internalTxList `json:"result"`
}

func newInternalTxData() *InternalTxData {
	return &InternalTxData{
		Result: internalTxList{},
	}
}

func (my InternalTxData) String() string          { return dbg.ToJSONString(my) }
func (my InternalTransaction) String() string     { return dbg.ToJSONString(my) }
func (my InternalTransactionList) String() string { return dbg.ToJSONString(my) }
