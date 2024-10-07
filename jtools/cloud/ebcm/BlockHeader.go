package ebcm

import (
	"jtools/dbg"
	"jtools/jmath"
	"jtools/unix"
)

type BlockHeader struct {
	Number      uint64 `bson:"number" json:"number"`
	BlockNumber string `bson:"blockNumber,omitempty" json:"blockNumber,omitempty"`

	ParentHash       string `bson:"parent_hash" json:"parent_hash"`
	RewardCoinBase   string `bson:"reward_coin_base" json:"reward_coin_base"` //coinBase(ETH)
	StateRoot        string `bson:"state_root" json:"state_root"`
	TransactionsRoot string `bson:"transactions_root" json:"transactions_root"`
	ReceiptsRoot     string `bson:"receipts_root" json:"receipts_root"`

	//LogsBloom    string `bson:"logs_bloom" json:"logs_bloom"`

	GasUsed   string    `bson:"gas_used" json:"gas_used"`
	Timestamp unix.Time `bson:"timestamp" json:"timestamp"`

	Hash string `bson:"hash" json:"hash"`
}

func (my BlockHeader) String() string { return dbg.ToJsonString(my) }

func NewBlockHeader(data *BlockByNumberData, is_klay bool) *BlockHeader {
	get_reward_coin_base := func() string {
		if is_klay {
			return data.RewardBase
		}
		return data.CoinBase
	}
	header := &BlockHeader{
		Number:      data.Number,
		BlockNumber: data.NumberString,

		ParentHash:       data.PreHash,
		RewardCoinBase:   get_reward_coin_base(),
		StateRoot:        data.Root,
		TransactionsRoot: data.TxHash,
		ReceiptsRoot:     data.ReceiptHash,

		GasUsed:   jmath.VALUE(data.GasUsed),
		Timestamp: data.Time,

		Hash: data.Hash,
	}

	return header
}
