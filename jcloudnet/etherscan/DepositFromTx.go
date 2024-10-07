package etherscan

import "fmt"

//DepositFromTx :
type DepositFromTx struct {
	BlockNumber  int64  //blockNumber
	Hash         string //txHash
	FromAddress  string //fromAddress
	Value        string //Value
	ConfirmCount int64  //ConfirmCount
	TokenDecimal string

	ContractAddress string `json:"contractAddress"`
	TokenSymbol     string `json:"tokenSymbol"`
}

//ToString :
func (my DepositFromTx) ToString() string {
	return `{
	BlockNumber     : ` + fmt.Sprintf("%v", my.BlockNumber) + `
	Hash            : ` + fmt.Sprintf("%v", my.Hash) + `
	FromAddress     : ` + fmt.Sprintf("%v", my.FromAddress) + `
	Value           : ` + fmt.Sprintf("%v", my.Value) + `
	ConfirmCount    : ` + fmt.Sprintf("%v", my.ConfirmCount) + `
	ContractAddress : ` + fmt.Sprintf("%v", my.ContractAddress) + `
	TokenSymbol     : ` + fmt.Sprintf("%v", my.TokenSymbol) + `
	TokenDecimal    : ` + fmt.Sprintf("%v", my.TokenDecimal) + `
}`
}

//DepositTransaction :
type DepositTransaction struct {
	ContractAddress string
	DepositAddress  string
	LastBlockNumber int64
	Txlist          []DepositFromTx
}

//NewDepositTransaction :
func NewDepositTransaction(contractAddress string, depositAddress string, txlist []DepositFromTx) *DepositTransaction {
	ins := &DepositTransaction{
		ContractAddress: contractAddress,
		DepositAddress:  depositAddress,
		Txlist:          txlist,
	}
	ins.LastBlockNumber = 0
	for _, tx := range txlist {
		if ins.LastBlockNumber < tx.BlockNumber {
			ins.LastBlockNumber = tx.BlockNumber
		}
	} //for
	return ins
}

//DepositTxsMultiContract : 멀티 컨트랙 포함
type DepositTxsMultiContract struct {
	Contracts       []string
	DepositAddress  string
	LastBlockNumber int64
	Txlist          []DepositFromTx
}

//NewDepositTxsMultiContract :
func NewDepositTxsMultiContract(contracts []string, depositAddress string, txlist []DepositFromTx) *DepositTxsMultiContract {
	ins := &DepositTxsMultiContract{
		Contracts:      contracts,
		DepositAddress: depositAddress,
		Txlist:         txlist,
	}
	ins.LastBlockNumber = 0
	for _, tx := range txlist {
		if ins.LastBlockNumber < tx.BlockNumber {
			ins.LastBlockNumber = tx.BlockNumber
		}
	} //for
	return ins
}
