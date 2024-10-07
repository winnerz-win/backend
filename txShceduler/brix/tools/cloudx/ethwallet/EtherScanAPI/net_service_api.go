package EtherScanAPI

import (
	"txscheduler/brix/tools/cloudx/ethwallet"
	"txscheduler/brix/tools/cloudx/ethwallet/Response/BalanceData"
	"txscheduler/brix/tools/cloudx/ethwallet/Response/TxExecuteStatus"
	"txscheduler/brix/tools/cloudx/ethwallet/Response/TxReceiptStatus"
)

/*
getTokenBalance : erc-20 토큰 잔액 조회
https://api.etherscan.io
/api?module=account&action=tokenbalance&
contractaddress=0x8d97C127236D3aEf539171394212F2e43ad701C4&
address=0xd1432c8fcaf299f184401dca99f2b4b77595025d&
tag=latest&
apikey=YourApiKeyToken
*/
func (my *ETHScanAPI) getTokenBalance(contract, address string) *ETHScanAPI {
	my.method = get
	my.url = "/api?module=account&action=tokenbalance&"
	my.makeGetParams(
		"contractaddress", contract,
		"address", address,
		"tag", "latest",
		"apikey", my.apikey)

	my.resultCallbackData = BalanceData.New()

	return my
}

/*
getEtherBalance : ETH 코인 잔액 조회
/api?module=account&action=balance&
address=0xd1432c8fcaf299f184401dca99f2b4b77595025d&
tag=latest&
apikey=YourApiKeyToken
*/
func (my *ETHScanAPI) getEtherBalance(address string) *ETHScanAPI {
	my.method = get
	my.url = "/api?module=account&action=balance&"
	my.makeGetParams(
		"address", address,
		"tag", "latest",
		"apikey", my.apikey)
	my.resultCallbackData = BalanceData.New()
	return my
}

// getEtherTransactions :
// 주소 별 '정상적인'거래 목록을 가져옵니다.
// [선택적 매개 변수] startblock : 결과를 검색하려면 시작 blockNo, endblock : 결과를 검색하려면 종료 blockNo
func (my *ETHScanAPI) getEtherTransactions(address, startblock string, endblock string) *ETHScanAPI {
	my.method = get
	my.url = "/api?module=account&action=txlist&"
	my.makeGetParams(
		"address", address,
		"startblock", startblock,
		"endblock", endblock, //"999999999",
		"sort", my.sortOrder,
		"apikey", my.apikey)

	my.resultCallbackData = ethwallet.NewResponeTransactional(address, false)

	return my
}

// getTokenTarnsactions :
// 주소별로 "ERC20-토큰 전송 이벤트"목록 가져 오기
// [선택적 매개 변수] startblock : 결과를 검색하려면 시작 blockNo, endblock : 결과를 검색하려면 종료 blockNo
func (my *ETHScanAPI) getTokenTransactions(address, startblock string, endblock string) *ETHScanAPI {

	my.method = get
	my.url = "/api?module=account&action=tokentx&"
	my.makeGetParams(
		"address", address,
		"startblock", startblock,
		"endblock", endblock, //"999999999",
		"sort", my.sortOrder,
		"apikey", my.apikey)

	my.resultCallbackData = ethwallet.NewResponeTransactional(address, true)

	return my
}

// getInternalTransactions :
// 내부 거래 목록 조회
func (my *ETHScanAPI) getInternalTransactions(address string, startblock, endblock string) *ETHScanAPI {
	my.method = get
	my.url = "/api?module=account&action=txlistinternal&"
	my.makeGetParams(
		"address", address,
		"startblock", startblock,
		"endblock", endblock, //"999999999",
		"sort", my.sortOrder,
		"apikey", my.apikey,
	)
	my.resultCallbackData = newInternalTxData()
	return my
}

// 컨트랙주소로 트랜젝션 조회 (페이지 단위)
func (my *ETHScanAPI) getContractTransactions(constract string, page string, offset string) *ETHScanAPI {

	my.method = get
	my.url = "/api?module=account&action=tokentx&"
	my.makeGetParams(
		"contractaddress", constract,
		"page", page,
		"offset", offset,
		"apikey", my.apikey)

	my.resultCallbackData = ethwallet.NewContractTransactional()
	my.failCallbackData = ethwallet.NewCallTimeLimitError("ContractTransactional")
	return my
}

// 컨트랙 주소로 트랙젝션 조회 (블럭 단위)
func (my *ETHScanAPI) getContractTxsByBlock(constract string, startblock, endblock string) *ETHScanAPI {

	my.method = get
	my.url = "/api?module=account&action=tokentx&"
	my.makeGetParams(
		"contractaddress", constract,
		"startblock", startblock,
		"endblock", endblock, //"999999999",
		"sort", my.sortOrder,
		"apikey", my.apikey)

	my.resultCallbackData = ethwallet.NewContractTransactional()
	my.failCallbackData = ethwallet.NewCallTimeLimitError("ContractTransactional")
	return my
}

/*
https://api.etherscan.io/api?module=transaction&action=getstatus&txhash=0x1834c9c372b6a991f4f6a9303022e89caa993ec27c228be467b1bf11e0ae4ffd&apikey=TJV15UE6DFRE7UHYZ1F1SJ6IWQHFKN9YUJ
*/
//checkTransactionExecuteStaus : txHash Excute 검증
func (my *ETHScanAPI) checkTransactionExecuteStaus(txHash string) *ETHScanAPI {
	my.method = get
	my.url = "/api?module=transaction&action=getstatus&"
	my.makeGetParams(
		"txhash", txHash,
		"apikey", my.apikey)
	my.resultCallbackData = TxExecuteStatus.New()
	return my
}

/*
tps://api.etherscan.io/api?module=transaction&action=gettxreceiptstatus&txhash=0x513c1ba0bebf66436b5fed86ab668452b7805593c05073eb2d51d3a52f480a76&apikey=YourApiKeyToken
*/
//checkTransactionReceiptStaus : txHash 검증
func (my *ETHScanAPI) checkTransactionReceiptStaus(txHash string) *ETHScanAPI {
	my.method = get
	my.url = "/api?module=transaction&action=gettxreceiptstatus&"
	my.makeGetParams(
		"txhash", txHash,
		"apikey", my.apikey)
	my.resultCallbackData = TxReceiptStatus.New()
	return my
}
