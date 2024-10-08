package ebcm

import (
	"jtools/dbg"
	"jtools/unix"
)

type BlockByNumberData struct {
	BlockData  `bson:",inline" json:",inline"`
	TxCount    int                  `bson:"tx_count" json:"tx_count"`
	TxHashList []string             `bson:"tx_hash_list,omitempty" json:"tx_hash_list,omitempty"`
	TxList     TransactionBlockList `bson:"txlist" json:"txlist"`
}

func (my BlockByNumberData) GetTxCount() int { return len(my.TxList) }
func (my BlockByNumberData) String() string  { return dbg.ToJsonString(my) }

type BlockData struct {
	Number       uint64    `bson:"number" json:"number"`
	NumberString string    `bson:"numberString" json:"numberString"`
	Time         unix.Time `bson:"time" json:"time"`
	Hash         string    `bson:"hash" json:"hash"`
	PreHash      string    `bson:"pre_hash" json:"pre_hash"`
	RewardBase   string    `bson:"reward_base" json:"reward_base"` //klay
	GasUsed      uint64    `bson:"gas_used" json:"gas_used"`
	Nonce        uint64    `bson:"nonce" json:"nonce"`
	Extra        string    `bson:"extra" json:"extra"`
	ReceiptHash  string    `bson:"receipt_hash" json:"receipt_hash"`
	Root         string    `bson:"root" json:"root"`
	Size         string    `bson:"size" json:"size"`
	TxHash       string    `bson:"tx_hash" json:"tx_hash"`

	BlockScore int64 `bson:"block_score,omitempty" json:"block_score,omitempty"` //klay

	//ETH
	CoinBase   string `bson:"coin_base,omitempty" json:"coin_base,omitempty"`   //eth
	Difficulty string `bson:"difficulty,omitempty" json:"difficulty,omitempty"` //eth
	GasLimit   uint64 `bson:"gas_limit,omitempty" json:"gas_limit,omitempty"`   //eth
	BaseFee    string `bson:"baseFee,omitempty" json:"baseFee,omitempty"`
}

func (my BlockData) String() string { return dbg.ToJsonString(my) }
