package etherscanapi

import (
	"errors"
	"fmt"
	"jcloudnet/etherscan"
	"jcloudnet/etherscan/Response/BalanceData"
	"jcloudnet/etherscan/Response/TxExecuteStatus"
	"jcloudnet/etherscan/Response/TxReceiptStatus"
	"jtools/jmath"
	"strings"
)

// GetEtherTransactions : 주소 별 '정상적인'거래 목록을 가져옵니다. ( ETHER )
func GetEtherTransactions(cfg Config, depositAddress string, startBlock string, endBlock ...string) *etherscan.ResponeTransactional {
	fixedEndblock := endblockNumber
	if len(endBlock) > 0 {
		fixedEndblock = endBlock[0]
	}
	result := getService(cfg).getEtherTransactions(depositAddress, startBlock, fixedEndblock).Execute().Body()
	if result == nil {
		return nil
	}
	txs := result.(*etherscan.ResponeTransactional)
	for i := 0; i < len(txs.Result); i++ {
		txs.Result[i].ContractAddress = "eth"
		txs.Result[i].IsContract = false
	}
	txs.TxSort(cfg.SortOrder())

	return txs
}

// GetTokenTransactions : 주소별로 "ERC20-토큰 전송 이벤트"목록 가져 오기
func GetTokenTransactions(cfg Config, depositAddress string, startBlock string, endBlock ...string) *etherscan.ResponeTransactional {
	fixedEndblock := endblockNumber
	if len(endBlock) > 0 {
		fixedEndblock = endBlock[0]
	}
	result := getService(cfg).getTokenTransactions(depositAddress, startBlock, fixedEndblock).Execute().Body()
	if result == nil {
		return nil
	}
	txs := result.(*etherscan.ResponeTransactional)
	for i := 0; i < len(txs.Result); i++ {
		txs.Result[i].IsContract = true
	}
	txs.TxSort(cfg.SortOrder())
	return txs
}

// GetTransactionAll : 주소별 ETH, Token 모두 가져옴.
func GetTransactionAll(cfg Config, depositAddress string, startBlock string, endBlock ...string) *etherscan.ResponeTransactional {
	fixedEndblock := endblockNumber
	if len(endBlock) > 0 {
		fixedEndblock = endBlock[0]
	}

	var txs *etherscan.ResponeTransactional
	resultETH := getService(cfg).getEtherTransactions(depositAddress, startBlock, fixedEndblock).Execute().Body()
	if resultETH != nil {
		txs = resultETH.(*etherscan.ResponeTransactional)
		if txs.IsSuccess() {
			for i := 0; i < len(txs.Result); i++ {
				txs.Result[i].ContractAddress = "eth"
				txs.Result[i].IsContract = false
			}
		}

	}

	resultToken := getService(cfg).getTokenTransactions(depositAddress, startBlock, fixedEndblock).Execute().Body()
	if resultToken != nil {
		txsToken := resultToken.(*etherscan.ResponeTransactional)
		if txsToken.IsSuccess() {
			if txs != nil {
				if txs.IsSuccess() == false {
					txs.Status = txsToken.Status
					txs.Message = txsToken.Message
					txs.Result = etherscan.ResultTxDatas{}
				}
			}
			for i := 0; i < len(txsToken.Result); i++ {
				txsToken.Result[i].IsContract = true
			}
			txs.Result = append(txs.Result, txsToken.Result...)
		}
	}

	if txs != nil {
		txs.TxSort(cfg.SortOrder())
	}

	return txs
}

// GetInternalTransaction : 내부 거래 목록 조회
func GetInternalTransaction(cfg Config, address string, strarBlock string, endBlock ...string) *InternalTxData {
	fixedEndblock := endblockNumber
	if len(endBlock) > 0 {
		fixedEndblock = endBlock[0]
	}
	result := getService(cfg).getInternalTransactions(address, strarBlock, fixedEndblock).Execute().Body()
	if result == nil {
		return nil
	}
	txs := result.(*InternalTxData)
	return txs
}

func GetInternalTransactionAPI(cfg Config, target string, startblock string, endblock ...string) InternalTransactionList {
	data := GetInternalTransaction(cfg, target, startblock, endblock...)
	if data == nil || data.Status != "1" {
		return InternalTransactionList{}
	}
	return data.Result.ToList()
}

// GetContractTransactions : Contract Address의 transactions 목록 가져오기
func GetContractTransactions(cfg Config, contract string, page, offset string) *etherscan.ContractTransactional {
	result := getService(cfg).getContractTransactions(contract, page, offset).Execute().Body()
	if result == nil {
		return nil
	}
	txs := result.(*etherscan.ContractTransactional)
	txs.TxSort()
	return txs
}

// GetContractTxByBlock : 컨트랙 주소로 트랙젝션 조회 (블럭 단위)
func GetContractTxByBlock(cfg Config, contract string, startBlock string, endblock ...string) *etherscan.ContractTransactional {
	startBlock = strings.TrimSpace(startBlock)

	eb := "9999999999999"
	if len(endblock) > 0 {
		n := endblock[0]
		if jmath.IsNum(n) {
			if jmath.CMP(n, startBlock) > 0 {
				eb = jmath.VALUE(n)
			}
		}
	}
	result := getService(cfg).getContractTxsByBlock(contract, startBlock, eb).Execute().Body()
	if result == nil {
		return nil
	}
	txs := result.(*etherscan.ContractTransactional)
	txs.TxSort()
	return txs
}

// TokenBalance : ERC-20 토큰 잔액 조회
func TokenBalance(cfg Config, contract, address string) string {
	result := getService(cfg).getTokenBalance(contract, address).Execute().Body()
	if result == nil {
		return ""
	}
	r := result.(*BalanceData.ValueResult)
	return r.ResultValue()
}

// TokenBalanceString :
func TokenBalanceString(cfg Config, contract, address string) (string, error) {
	result := getService(cfg).getTokenBalance(contract, address).Execute().Body()
	if result == nil {
		return "0", errors.New("http.response.fail")
	}
	r := result.(*BalanceData.ValueResult)
	if r.IsStatusOk() == false {
		return "0", fmt.Errorf("%v", r.Message)
	}
	return r.ResultValue(), nil
}

// ETHBalance : ETH 잔액 조회
func ETHBalance(cfg Config, address string) string {
	result := getService(cfg).getEtherBalance(address).Execute().Body()
	if result == nil {
		return ""
	}
	r := result.(*BalanceData.ValueResult)
	return r.ResultValue()
}

// ETHBalanceString :
func ETHBalanceString(cfg Config, address string) (string, error) {
	result := getService(cfg).getEtherBalance(address).Execute().Body()
	if result == nil {
		return "0", errors.New("http.response.fail")
	}
	r := result.(*BalanceData.ValueResult)
	if r.IsStatusOk() == false {
		return "0", fmt.Errorf("%v", r.Message)
	}
	return r.ResultValue(), nil
}

// CheckTransactionExecutionStatus :
func CheckTransactionExecutionStatus(cfg Config, txHash string) *TxExecuteStatus.StatusResult {
	result := getService(cfg).checkTransactionExecuteStaus(txHash).Execute().Body()
	return result.(*TxExecuteStatus.StatusResult)
}

// CheckTransactionReceiptStaus : Result- "1":성공 , "0":실패, "":대기
func CheckTransactionReceiptStaus(cfg Config, txHash string) *TxReceiptStatus.StatusResult {
	result := getService(cfg).checkTransactionReceiptStaus(txHash).Execute().Body()
	return result.(*TxReceiptStatus.StatusResult)
}
