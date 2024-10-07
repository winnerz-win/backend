package ecs

import (
	"jtools/cloud/ebcm"

	"github.com/ethereum/go-ethereum/core/types"
)

type WrappedTx struct {
	*types.Transaction
}

func WrappedTransaction(tx *types.Transaction) ebcm.WrappedTransaction {
	if tx == nil {
		return nil
	}
	return &WrappedTx{tx}
}

func (my WrappedTx) HashHex() string {
	tx := my.Transaction
	if tx == nil {
		return ""
	}
	return my.Hash().Hex()
}
func (my WrappedTx) To() string {
	tx := my.Transaction
	return tx.To().Hex()
}

func (my WrappedTx) Size() string {
	tx := my.Transaction
	return tx.Size().TerminalString()
}
