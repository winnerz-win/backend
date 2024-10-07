package api

func set_base_callback_api_doc() {
	Doc().Message("<---- 기본 입출금 기능 콜백 API ---->")
	Doc().Message(`
		<cc_blue>[ SERVICE CALLBACK ]</cc_blue> ( 1. 입금 콜백 ) 가입한 회원 주소로 코인이 입금되었을 경우 입금 내역 전송
		
		/v1/coins/deposit/callback

		method : post
		content-type : application/json

		data : 
		{
			"name" : string,			// 가입 회원ID
			"address" : string,			// 회원 입금 주소
			"symbol" : string,			// 들어온 토큰 심볼 (ETH , ...)
			"price" : string,			// 입금 수량 ("1.333")
			"hash" : string,			// 트랜젝셕 해시
			"deposit_result" : bool,	// 입금 성공/실패 여부
			"timestamp" : int			// 트랜젝션 확인된 시간 (mms)
		}
	`)

	Doc().Message(`
		<cc_blue>[ SERVICE CALLBACK ]</cc_blue> ( 2. 코인 출금 결과 콜백 ) 요청으로 코인 출금이 성공/실패 했을 경우 출금 결과 전송
		
		/v1/coins/withdraw/callback

		method : post
		content-type : application/json

		data : 
		{
			"receipt_code" : string,	// 출금신청시 발급한 영수증
			"name" : string,			// 가입 회원ID
			"from_address" : string,	// 회원 입금 주소
			"to_address" : string,		// 토큰을 받을 주소
			"symbol" : string,			// 출금 신청한 토큰 심볼 (ETH , ...)
			"price" : string,			// 출금 신청한 토큰 수량 ("1.333")
			"hash" : string,			// 출금 트랜젝셕 해시
			"withdraw_result" : int,	// 출금 요청 결과 ( 성공:200 , 실패:104 )
			"timestamp" : int			// 출금 확인된 시간 (mms)
		}
	`)

	Doc().Message(`
		<cc_blue>[ SERVICE CALLBACK ]</cc_blue> ( ex. 개인지갑 코인 출금신청 결과 ) 요청으로 코인 출금이 성공/실패 했을 경우 출금 결과 전송
		
		/v1/coins/withdraw_self/callback

		method : post
		content-type : application/json

		data : 
		{
			"receipt_code" : string,	// 출금신청시 발급한 영수증
			"name" : string,			// 가입 회원ID
			"from_address" : string,	// 회원 입금 주소
			"to_address" : string,		// 토큰을 받을 주소
			"symbol" : string,			// 출금 신청한 토큰 심볼 (ETH , ...)
			"price" : string,			// 출금 신청한 토큰 수량 ("1.333")
			"hash" : string,			// 출금 트랜젝셕 해시
			"withdraw_result" : int,	// 출금 요청 결과 ( 성공:200 , 실패:104 )
			"timestamp" : int			// 출금 확인된 시간 (mms)
		}
	`)

	Doc().Message(`
		<cc_blue>[ SERVICE CALLBACK ]</cc_blue> ( 3. 개인지갑 코인 마스터지갑으로 전송시 콜백 ) 회원 개인지갑의 코인이 마스터지갑으로 전송시 이동한 코인량과 현재 개인지갑의 코인 잔액		
		
		/v1/coins/tomaster/callback

		method : post
		content-type : application/json

		data : 
		{
			"name" : string,			// 가입 회원ID
			"uid" : int, 				// 가입 회원의 고유UID
			"from_address" : string,	// 회원 입금 주소 (개인주소 , 보낸이)
			"master_address" : string,	// 마스터 지갑주소 (받는이)
			"contract" : string,        // 이더전송이면 "eth" , 토큰 전송이면 컨트랙트 주소
			"symbol" : string,			// 마스터로 전송된 토큰 심볼 (ETH , GDG ...)
			"price" : string,			// 마스터로 전송된 토큰 수량 ("1.333")
			"gas_fee" : string,         // 마스터로 전송할때 사용된 Gas(ETH) 수수료
			"hash" : string,			// 전송 트랜젝션 해시
			"timestamp" : int			// 전송 확인된 시간 (mms)

			"remain_coin" : {			// 본 패킷을 보낼 당시의 개인지갑 코인 잔액 (전송 하고 난 이후의 개인지갑 블록체인 잔액)
				"ETH" : "0.11",
				"GDG" : "0",
			}
		}
	`)

	Doc().Message(`
		<cc_blue>[ SERVICE CALLBACK ]</cc_blue> ( 4. 개인지갑이 아닌 기타 외부에서 마스터 지갑으로 코인을 전송시 콜백 ) - 전송자의 주소와 코인종류 및 가격 , 트랜젝션 해시를 콜백함
		
		/v1/coins/exmaster/callback

		method : post
		content-type : application/json

		data : 
		{
			"hash" : string,			// 트랜젝션 해시
			"from" : string,			// 전송자의 지갑 주소 (0x...)
			"contract" : string,        // 이더전송이면 "eth" , 토큰 전송이면 컨트랙트 주소
			"symbol" : string,			// 마스터로 전송된 토큰 심볼 (ETH , GDG ...)
			"price" : string,			// 마스터로 전송된 토큰 수량 ("1.333")
			"timestamp" : int			// 전송 확인된 시간 (mms)
		}
	`)

	Doc().Message(`
		<cc_blue>[ SERVICE CALLBACK ]</cc_blue> ( 5. 마스터 계좌에서 외부 계좌로 출금 신청 결과 콜백  )
		
		/v1/coins/master_out/callback

		method : post
		content-type : application/json

		data : 
		{
			"receipt_code" : string,	// 주문번호
			"hash" : string,			// 트랜젝션 해시
			"from_address" : string,	// 마스터 지갑주소
			"to_address" : string		// 수신 지갑주소
			"symbol" : string,			// 토큰 심볼 (ETH , GDG ...)
			"price" : string,			// 토큰 수량 ("1.333")

			"withdraw_result" : int,	// 출금 요청 결과 ( 성공:200 , 실패:104 )
			"fail_message" : string 	// 실패 사유 ( 실패일경우 )
			"timestamp" : int			// 결과 시간 (mms)
		}
		----------------------------------------------
		실패사유 : fail_message
		invalid_symbol 				:	등록되어있지 않는 토큰 심볼 오류
		need_price 					:	마스터지갑의 잔액부족 (토큰/이더 수량 , 가스비 등등)
		chain_error:box 			:	메인 노드 에러
		chain_error:nonce 			:	메인 노드 에러
		chain_error:tx 				:	메인 노드 에러
		chain_error:send 			:	메인 노드 에러
		chain_error:pending_time 	:	전송한 트랜젝션이 시간 지연(4시간 체크)으로 인한 fallback 처리
		chain_error:fail 			:	실패한 트렉젝션
	`)

}
