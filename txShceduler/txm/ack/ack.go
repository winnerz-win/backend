package ack

import (
	"txscheduler/brix/tools/jnet/chttp"
)

var (
	HandlerPanic       = chttp.Error(1, "서버 핸들러 오류 (스케줄러서버 장애)")
	BadParam           = chttp.Error(1001, "요청 파라미터 오류(입력값 오류)")
	NotAllowSystem     = chttp.Error(1004, "시스템에서 지원하지 않는 기능")
	ExistedName        = chttp.Error(4001, "이미 존재하는 Name(유저)입니다.")
	NotFoundName       = chttp.Error(4101, "존재하지 않는 회원ID")
	NotFoundAddress    = chttp.Error(4102, "존재하지 않는 입금주소")
	NotFoundSymbol     = chttp.Error(4103, "지원하지 않는 코인 심볼(심볼은 대소문자 구분합니다.)")
	InvalidAddress     = chttp.Error(4104, "이더리움 주소 체계의 형식이 아님(0x.....)")
	NotAllowSameID     = chttp.Error(4105, "본인의 입금주소로 코인 출금 불가")
	UnderZERO          = chttp.Error(4107, "출금 요청 수량이 0이하일 경우")
	InvalidReceiptCode = chttp.Error(4108, "잘못된 출금영수증 코드")

	NotEnoughGasPrice    = chttp.Error(4201, "not enough gas price")
	NotEnoughTargetPrice = chttp.Error(4202, "not enough target price")

	DBJob = chttp.Error(5001, "서버 DB 작업 실패 (스케줄러서버 장애)")

	NotYetSendTx = chttp.Error(5002, "아직 트랜잭션 작업을 할 수 없습니다.(다른 프로세스에서 작업중)")

	ChainNonce    = chttp.Error(5101, "[Chain] nonce error")
	ChainGasLimit = chttp.Error(5102, "[Chain] gasLimit error")
	ChainGasPrice = chttp.Error(5102, "[Chain] gasPrice error")
	ChainNTX      = chttp.Error(5102, "[Chain] normal transaction error")
	ChainSTX      = chttp.Error(5102, "[Chain] signed transaction error")
	ChainSend     = chttp.Error(5102, "[Chain] send error")

	InvalidNick       = chttp.Error(9001, "잘못된 아이디")
	InvalidPassword   = chttp.Error(9002, "잘못된 비밀번호")
	InvalidToeken     = chttp.Error(9003, "잘못된 토큰")
	TokenExpired      = chttp.Error(9004, "토큰 기간만료")
	InvalidRootAdmin  = chttp.Error(9005, "Root권한 관리자만 요청 할 수 있다.")
	AlreadyProcessJob = chttp.Error(9006, "현재 진행중이 작업이 있다")

	NFTNotfoundTokenID = chttp.Error(3001, "tokenId를 찾을수 없음.")
	NFTZeroETH         = chttp.Error(3002, "ETH잔액이 0이하")
	NFTBuyPrice        = chttp.Error(3003, "NFT 구매/전송 비용 부족")
	NFTBuyTryInsert    = chttp.Error(3004, "NFT 구매 시도 실패 (DB 작업 실패)")
	NFTNotfoundData    = chttp.Error(3005, "NFT not found data")
	NFTExistTokenID    = chttp.Error(3006, "NFT existed token_id")
	NFTTransactionFail = chttp.Error(3007, "NFT transaction fail")
	NFTTTransferFail   = chttp.Error(3008, "NFT transfer fail")

	NFT_RPC_TIMEOUT       = chttp.Error(30303, "NFT 블록체인(RPC) 노트 통신 장애")
	NFT_InvalidOwner      = chttp.Error(3009, "NFT invalid owner")
	NFT_NeedETHPrice      = chttp.Error(3010, "NFT need more eth")
	NFT_InvalidPayPrice   = chttp.Error(3011, "NFT invalid pay price(estimate transaction fee error)")
	NFT_Param_price       = chttp.Error(3012, "NFT 요청 금액이 0이하 또는 오류값")
	NFT_TokenId_Format    = chttp.Error(3013, "NFT 토큰ID값이 잘못되었습니다.")
	NFT_NotfoundPayer     = chttp.Error(3014, "NFT Not found payer")
	NFT_NotfoundOwner     = chttp.Error(3015, "NFT Not found owner")
	NFT_ReceiptCodeFormat = chttp.Error(3016, "NFT invalid receiptCode format")
	NFT_SameBaseURI       = chttp.Error(3017, "NFT baseURI가 기존과 동일합니다.")

	OWNER_NotFoundData = chttp.Error(6000, "[OwnerTask] 요청한 데이터를 찾을수 없습니다.")
	OWNER_RpcFail      = chttp.Error(6001, "[OwnerTask] Blockchain rpc node 응답 실패")

	OWNER_ALREADY_LOCKED = chttp.Error(6002, "already locked address.")
	OWNER_REQ_JOB_ONCE   = chttp.Error(6003, "락업관련 예약중이 작업이 존재함.")

	OWNER_Transfer_SameAddress     = chttp.Error(6101, "[OwnerTransfer] recipient is same user.")
	OWNER_Transfer_EmptyDatas      = chttp.Error(6102, "[OwnerTransfer] transfers is empty.")
	OWNER_Transfer_Price           = chttp.Error(6103, "[OwnerTransfer] price is under zero.")
	OWNER_Transfer_ReleaseTime     = chttp.Error(6104, "[OwnerTransfer] release_time less current server time.")
	OWNER_Transfer_MemberRecipient = chttp.Error(6105, "[OwnerTransfer] recipient is member.")

	OWNER_Lock_OverPrice = chttp.Error(6201, "[OwnerLock] 락업 요청 금액이 보유잔액 보다 클경우")

	OWNER_Unlock_PositionIndex   = chttp.Error(6301, "[OwnerUnlock] position_index가 rpc-node의 정보와 불일치.")
	OWNER_Unlock_AlreadyTimeOver = chttp.Error(6302, "[OwnerUnlock] 이미 락업 해재 시간이 지났을경우.")
)
