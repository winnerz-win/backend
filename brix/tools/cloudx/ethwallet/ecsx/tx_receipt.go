package ecsx

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"txscheduler/brix/tools/jmath"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TxResult :
type TxResult string

// String : pending , success , fail
func (my TxResult) String() string {
	return string(my)
}

const (
	TxResultNoData  = TxResult("nodata")
	TxResultPending = TxResult("pending")
	TxResultOK      = TxResult("success")
	TxResultFail    = TxResult("fail")
)

// Receipt :
type Receipt struct {
	client           *ethclient.Client
	data             *types.Receipt
	TransactionIndex uint
}

// String :
func (my Receipt) String() string {
	if my.data == nil {
		return `{ "status" : "nodata" }`
	}
	b, _ := json.MarshalIndent(my.data, "", "  ")
	return string(b)
}

// BlockNumber :
func (my Receipt) BlockNumber() string {
	if my.data != nil {
		return my.data.BlockNumber.String()
	}
	return "0"
}

// CumulativeGasUsed : payed limit
func (my Receipt) CumulativeGasUsed() uint64 {
	if my.data != nil {
		return jmath.Uint64(my.data.CumulativeGasUsed)
	}
	return 0
}

// GasUsed :""
func (my Receipt) GasUsed() uint64 {
	if my.data != nil {
		return jmath.Uint64(my.data.GasUsed)
	}
	return 0
}

// GasUsedString :
func (my Receipt) GasUsedString() string {
	if my.data != nil {
		return jmath.VALUE(my.data.GasUsed)
	}
	return "0"
}

// ParseString :
func (my Receipt) ParseString() string {
	if my.data == nil {
		return `{ "status" : "nodata" }`
	}
	c := my.data

	jdata := map[string]interface{}{}
	jdata["status"] = jmath.NewBigDecimal(c.Status).ToString()
	jdata["cumulativeGasUsed"] = jmath.NewBigDecimal(c.CumulativeGasUsed).ToString()
	jdata["gasUsed"] = jmath.NewBigDecimal(c.GasUsed).ToString()
	jdata["blockHash"] = strings.ToLower(c.BlockHash.Hex())
	jdata["transactionHash"] = strings.ToLower(c.TxHash.Hex())
	jdata["blockNumber"] = c.BlockNumber.String()
	jdata["contractAddress"] = strings.ToLower(c.ContractAddress.Hex())
	jdata["logs"] = c.Logs
	jdata["logsBloom"] = c.Bloom

	b, _ := json.MarshalIndent(jdata, "", "  ")
	return string(b)
}

// Result : TxResult(nodata, pending , success, fail)
func (my Receipt) Result(isContract bool) TxResult {
	if my.data == nil {
		return TxResultNoData
	}

	if isContract {
		if my.data.Status == 0 {
			if len(my.data.Logs) == 0 {
				return TxResultFail
			}
			return TxResultFail
			//return TxResultOK
		}
	}

	// if len(my.data.Logs) == 0 {
	// 	return TxResultOK
	// }
	if my.data.Status == 0 {
		return TxResultFail
	}
	return TxResultOK
}

// GetData :
func (my Receipt) GetData() *types.Receipt {
	return my.data
}

// TransactionInBlock :
func (my Receipt) TransactionInBlock() (*types.Transaction, error) {
	if my.data == nil {
		return nil, errors.New("Receipt data is nil")
	}
	tx, err := my.client.TransactionInBlock(
		context.Background(),
		my.data.BlockHash,
		my.data.TransactionIndex,
	)
	return tx, err
}
