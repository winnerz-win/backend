package EtherClient

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
)

//logTransaction :
type logTransaction struct {
	*types.Transaction
}

//TransactionLog :
func TransactionLog(tx *types.Transaction, isView ...bool) string {
	logTx := logTransaction{Transaction: tx}
	if len(isView) > 0 && isView[0] {
		fmt.Println("**** EtherClient.TransactionLog ****")
		logTx.View()
	}
	return logTx.ToString()
}

//ToString :
func (my *logTransaction) ToString() string {

	return `{
	Hash     : ` + my.Hash().Hex() + `
	Cost     : ` + my.Cost().String() + `
	Value    : ` + my.Value().String() + `
	ChainId  : ` + my.ChainId().String() + `
	Gas      : ` + fmt.Sprintf("%v", my.Gas()) + `
	GasPrice : ` + my.GasPrice().String() + `
	Nonce    : ` + fmt.Sprintf("%v", my.Nonce()) + `
	Data     : ` + hex.EncodeToString(my.Data()) + `
	Size     : ` + my.Size().TerminalString() + `
}`
}

//View :
func (my logTransaction) View() {
	fmt.Println(my.ToString())
}
