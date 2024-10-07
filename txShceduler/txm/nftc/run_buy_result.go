package nftc

import (
	"time"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

/*
	( 6. NFT 구매 결과 콜백  )

	method : post
	content-type : application/json

	전송 데이터
	{
		"uid" : long				// 유저 고유ID (스케줄러 UID)
		"address" : string			// 유저 발급 주소 (유저 입금 계좌)
		"name" : string				// 유저 이름
		"receipt_code" : string		// 구매신청시 발금한 영수증 코드

		"pay_address": string		// 실제 구매비용 지불한 주소 ( address == pay_address : 회원납부(ETH) ,  pay_address == 입금마스터주소 : 회사납부(GDG) )

		"pay_symbol" : string		// ETH , ERCT (구매 토큰 심볼 , 상용은 GDG)
		"pay_price" : string		// NFT의 구매 가격(decimal제외한 값)
		"is_pay_free" : string 		// 무료발행 여부 (무료발행시 pay_symbol, pay_price 는 무시됨.)

		"last_owner" : string		// 최종 소유자 주소 ( address == last_owner , 추후에 유저간 거래시에 last_owner가 바뀔수 있다.)

		"token_id" : string			// NFT 토큰 ID
		"token_type" : string		// NFT 토큰 타입

		"token_uri" : string 		// NFT 토큰 URI
		"meta" : {					// NFT 메타데이타 (Meta)
			"token_type" : string	// [Meta] NFT 토큰 타입
			"name" : string			// [Meta] 이름
			"content" : string		// [Meta] 컨텐츠 텍스트
			"deployer" : string		// [Meta] 발행자 텍스트
			"desc" : string			// [Meta] 설명 텍스트
		}

		"deposit_hash"				// 구매대금("pay_price") 트랜젝션 해시값
		"hash" : string				// NFT 트랜젝션 해시값

		"status" : int				// 구매 결과 ( 성공:200 , 실패:104 )
		"fail_code" : int	 		// 실패 코드 ( 실패일경우 )
		"fail_message" : string 	// 실패 사유 ( 실패일경우 )

		"timestamp" : int			// 결과 시간 (mms)

		"gas_fee_eth" : string		// 트랜젝션의 가스비용 (ETH값)

		"is_burn" : bool 			// 기본값 false (해당 NFT아이템이 소각 되었을 경우 true)
		"burn_hash" : string		// is_burn == true 일때의 소각한 트랜젝션 해시 값.
	}
	----------------------------------------------
	실패코드 - fail_code(fail_message)
	0	 ("")						: 없음
	1001 (pay_deposit_fail)			: 구매대금 지불 실패 (ETH == 구매가(ETH)+가스비(ETH) , GDG == 구매가(GDG) + 가스비(ETH) )
	1002 (transaction_result_fail)	: 구매 트랜젝션 결과 실패
	1003 (transaction_try_fail)		: NFT구매 트랜젝션 시도 실패 ( 토큰 잔액, 이더 잔액 , Approve상태에 따른 경우의수가 존재)

*/

type NftBuyResultCallback func(db mongo.DATABASE, list model.NftBuyEndList)

func runBuyResult(resultCallback NftBuyResultCallback) {
	dbg.PrintForce("NftBuyResultCallback ------ START")
	for {
		model.DB(func(db mongo.DATABASE) {
			list := model.NftBuyEndList{}
			selector := mongo.Bson{"is_send": false}
			db.C(inf.NFTBuyEnd).Find(selector).All(&list)

			if len(list) > 0 {
				resultCallback(db, list)
			}

		})
		time.Sleep(time.Second)
	} //for
}
