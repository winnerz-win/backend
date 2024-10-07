package ethwallet

import (
	"strings"

	"txscheduler/brix/tools/jmath"
)

//ResponeTransactional :
type ResponeTransactional struct {
	Status  string        `json:"status"`  //  1, 0
	Message string        `json:"message"` // OK , No transactions found
	Result  ResultTxDatas `json:"result"`

	_isRefactData    bool
	_depositAddr     string
	_lastBlockNumber int64
	_isToken         bool
}

//IsToken : do not used by GetTransactionAll
func (my ResponeTransactional) IsToken() bool {
	return my._isToken
}

//IsSuccess :
func (my ResponeTransactional) IsSuccess() bool {
	return my.Status == "1"
}

//NewResponeTransactional :
func NewResponeTransactional(depositAddr string, isToken bool) *ResponeTransactional {
	return &ResponeTransactional{
		_depositAddr: depositAddr,
		_isToken:     isToken,
	}
}

//GETPTR :
func (my *ResponeTransactional) GETPTR() interface{} {
	if my == nil {
		return nil
	}
	return my
}

//ToString :
func (my ResponeTransactional) ToString() string {
	return `{
status  :` + my.Status + `,
message :` + my.Message + `,
result  :` + my.Result.ToString() + `
}`

	//return fmt.Sprintf("{status:%v, message:%v, result:%v}", my.Status, my.Message, my.Result.ToString())
}

//LastBlockNumber :
func (my ResponeTransactional) LastBlockNumber() int64 {
	return my._lastBlockNumber
}

//TxSort :
func (my *ResponeTransactional) TxSort(sortOrder ...string) {
	if my._isRefactData == true {
		return
	}
	my._isRefactData = true

	count := len(my.Result)
	if count == 0 {
		return
	} else {
		my._lastBlockNumber = my.Result.rtxSort(sortOrder...)
		//my._lastBlockNumber, _ = strconv.ParseInt(my.Result[count-1].BlockNumber, 10, 64)
	}
}

//TxCount :
func (my ResponeTransactional) TxCount() int {
	return len(my.Result)
}

//GetContractData :
func (my *ResponeTransactional) GetContractData(contractAddress string, confirmCount int, isOutDecimal bool) *DepositTransaction {
	if my.TxCount() == 0 {
		return nil
	}

	if my._isRefactData == false {
		my.TxSort()
	}

	contractAddress = strings.ToLower(contractAddress)
	depositAddress := strings.ToLower(my._depositAddr)

	depositTxlist := []DepositFromTx{}
	for _, tx := range my.Result {
		tx.ContractAddress = strings.ToLower(tx.ContractAddress)
		if tx.ContractAddress != contractAddress {
			continue
		}
		tx.To = strings.ToLower(tx.To)
		if tx.To != depositAddress {
			continue
		}
		//if tx._confirmCount < Define.ETHEREUM_TRANSACTION_CONFIRM_COUNT {
		if tx._confirmCount < int64(confirmCount) {
			continue
		}

		if isOutDecimal == true {
			outValue := "0"
			bigVal := jmath.NewBigDecimal(tx.Value)
			if bigVal.CompareTo(jmath.BigDecimal_ZERO()) > 0 {
				//outValue = bigVal.Divide(jmath.BigDecimal_TEN().Pow(Define.DECIMAL_ERC20)).ToString()
				outValue = bigVal.Divide(jmath.BigDecimal_TEN().Pow(tx.TokenDecimal)).ToString()
			}
			tx.Value = outValue
		}

		data := DepositFromTx{
			BlockNumber:  tx._blockNumber,
			Hash:         tx.Hash,
			FromAddress:  tx.From,
			Value:        tx.Value,
			ConfirmCount: tx._confirmCount,
			TokenDecimal: tx.TokenDecimal,

			ContractAddress: strings.ToLower(tx.ContractAddress),
			TokenSymbol:     tx.TokenSymbol,
		}
		depositTxlist = append(depositTxlist, data)

	} //for

	dt := NewDepositTransaction(contractAddress, depositAddress, depositTxlist)
	return dt
}

//GetContractDataList :
func (my *ResponeTransactional) GetContractDataList(contracts []string, confirmCount int, isOutDecimal bool) *DepositTxsMultiContract {
	if my.TxCount() == 0 {
		return nil
	}

	if my._isRefactData == false {
		my.TxSort()
	}

	for i := 0; i < len(contracts); i++ {
		contracts[i] = strings.ToLower(contracts[i])
	}
	depositAddress := strings.ToLower(my._depositAddr)

	depositTxlist := []DepositFromTx{}
	for _, tx := range my.Result {
		tx.To = strings.ToLower(tx.To)
		if tx.To != depositAddress {
			continue
		}
		if tx._confirmCount < int64(confirmCount) {
			continue
		}

		isContinue := true
		tx.ContractAddress = strings.ToLower(tx.ContractAddress)
		for _, contractAddr := range contracts {
			if tx.ContractAddress == contractAddr {
				isContinue = false
				break
			}
		} //for

		if isContinue == true {
			continue
		}

		if isOutDecimal == true {
			outValue := "0"
			bigVal := jmath.NewBigDecimal(tx.Value)
			if bigVal.CompareTo(jmath.BigDecimal_ZERO()) > 0 {
				//outValue = bigVal.Divide(jmath.BigDecimal_TEN().Pow(Define.DECIMAL_ERC20)).ToString()
				outValue = bigVal.Divide(jmath.BigDecimal_TEN().Pow(tx.TokenDecimal)).ToString()
			}
			tx.Value = outValue
		}

		data := DepositFromTx{
			BlockNumber:  tx._blockNumber,
			Hash:         tx.Hash,
			FromAddress:  tx.From,
			Value:        tx.Value,
			ConfirmCount: tx._confirmCount,
			TokenDecimal: tx.TokenDecimal,

			ContractAddress: strings.ToLower(tx.ContractAddress),
			TokenSymbol:     tx.TokenSymbol,
		}
		depositTxlist = append(depositTxlist, data)

	} //for

	dt := NewDepositTxsMultiContract(contracts, depositAddress, depositTxlist)
	return dt
}
