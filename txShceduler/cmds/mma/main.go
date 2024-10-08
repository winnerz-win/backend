package main

import (
	"jtools/mms"
	"time"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/jnet/cnet"
	"txscheduler/nft_winners"
	"txscheduler/txm"
	"txscheduler/txm/cloud"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
	"txscheduler/txm/scv"
)

/*
http://192.168.0.9/git/brix.git
http://192.168.0.9/git/txShceduler.git
*/

func main() {
	mainnet := false
	config := &inf.IConfig{
		Mainnet: mainnet,
		Version: "2024.03.04" + nft_winners.Version,
		Seed:    "txscheduler/dummy_seed",
		DB:      "txm_mma",
		IPCheck: false,
		ClientHost: map[bool][]string{
			false: inf.NewHost("http://192.168.0.111", "18080"),
			true:  inf.NewHost("https://192.168.0.123", ""),
		},
		AdminSalt: "admin_salt",

		Confirms: 1,

		IsLockTransferByOwner: true,
		Owners: inf.KeyPairList{
			{
				Mainnet:    false,
				PrivateKey: "private_key",
				Address:    "address",
			},
			{
				Mainnet:    true,
				PrivateKey: "private_key",
				Address:    "address",
			},
		},

		Masters: inf.KeyPairList{
			{
				Mainnet:    false,
				PrivateKey: "private_key",
				Address:    "address",
			},
			{
				Mainnet:    true,
				PrivateKey: "private_key",
				Address:    "address",
			},
		},
		Chargers: inf.KeyPairList{
			{
				Mainnet:    false,
				PrivateKey: "private_key",
				Address:    "address",
			},
			{
				Mainnet:    true,
				PrivateKey: "private_key",
				Address:    "address",
			},
		},
		Tokens: inf.TokenInfoList{
			{
				Mainnet:  false,
				Contract: "eth",
				Symbol:   "ETH",
				Decimal:  "18",
			},
			{
				Mainnet:  false,
				Contract: "token_contract_address", //(sepolia)
				Symbol:   "WNZ",
				Decimal:  "18",
			},
			{
				Mainnet:  true,
				Contract: "eth",
				Symbol:   "ETH",
				Decimal:  "18",
			},
			{
				Mainnet:  true,
				Contract: "token_contract_address",
				Symbol:   "WNZ",
				Decimal:  "18",
			},
		},
		InfuraKeys: []string{
			"infura_key_1",
			"infura_key_2",
			"infura_key_3",
			"infura_key_4",
		},
		ESKeys: []string{
			"es_key_1",
			"es_key_2",
			"es_key_3",
		},
	}

	nft_config := nft_winners.NftConfig{}
	if !mainnet {
		nft_config.NftContract = "nft_contract_addresss" //sepolia
	} else {
		nft_config.NftContract = "nft_contract_address"
	}

	if err := nft_winners.InjectNftConfig(nft_config); err != nil {
		dbg.Exit("NFT_WINNERS_CONFIG : ", err)
	}

	cloud.ONLY_TOKEN_TO_MASTER = true
	txm.Start(
		config,
		makeCallback(),
		nil,
	)
}

func makeCallback() scv.CallbackList {
	return scv.CallbackList{
		{
			SubPackage: nft_winners.Ready,
		},
		{
			HKey:     scv.KeyNone,
			StartMsg: "DepositCallback ------ START",
			SleepDu:  time.Second,
			LooperFunc: func() {
				model.DB(func(db mongo.DATABASE) {
					list := model.LogDeposit{}.GetList(db)
					if len(list) == 0 {
						return
					}
					sendDeposit(db, list)
				})
			},
		},
		{
			HKey:     scv.KeyNone,
			StartMsg: "WithdrawCallback ------ START",
			SleepDu:  time.Second,
			LooperFunc: func() {
				model.DB(func(db mongo.DATABASE) {
					list := model.LogWithdraw{}.GetList(db)
					if len(list) == 0 {
						return
					}
					sendWithdraw(db, list)
				})
			},
		},
		{
			HKey:     scv.KeyNone,
			StartMsg: "WithdrawCallbackSELF ------ START",
			SleepDu:  time.Second,
			LooperFunc: func() {
				model.DB(func(db mongo.DATABASE) {
					list := model.LogWithdrawSELF{}.GetList(db)
					if len(list) == 0 {
						return
					}
					sendWithdrawSELF(db, list)
				})
			},
		},
		{
			HKey:     scv.ToMasterCallback,
			StartMsg: "ToMasterCallback ------ START",
			SleepDu:  time.Second,
			LooperFunc: func() {
				model.DB(func(db mongo.DATABASE) {
					list := model.LogToMaster{}.GetList(db)
					if len(list) == 0 {
						return
					}
					sendToMaster(db, list)
				})
			},
		},
		{
			HKey:     scv.ExMasterCallback,
			StartMsg: "ExMasterCallback ------ START",
			SleepDu:  time.Second,
			LooperFunc: func() {
				model.DB(func(db mongo.DATABASE) {
					list := model.LogExMaster{}.GetList(db)
					if len(list) == 0 {
						return
					}
					sendExMaster(db, list)
				})
			},
		},
		{
			HKey:     scv.MasterOutCallback,
			StartMsg: "MasterOutCallback ------ START",
			SleepDu:  time.Second,
			LooperFunc: func() {
				model.DB(func(db mongo.DATABASE) {
					list := model.TxETHMasterOut{}.GetCallbackList(db)
					if len(list) == 0 {
						return
					}
					sendMasterOut(db, list)
				})
			},
		},
	}
}

func sendDeposit(db mongo.DATABASE, list model.LogDepositList) {
	/*
		황구택
		01076684789

		가입한 회원 주소로 코인이 입금되었을 경우 입금 내역 전송
		http://service.server.com:8080/v1/coins/deposit/callback

		method : post
		content-type : application/json

		전송 데이터
		{
			"name" : string,			// 가입 회원ID
			"address" : string,			// 회원 입금 주소
			"symbol" : string,			// 들어온 토큰 심볼 (ETH , ...)
			"price" : string,			// 입금 수량 ("1.333")
			"hash" : string,			// 트랜젝셕 해시
			"deposit_result" : bool,	// 입금 성공/실패 여부
			"timestamp" : int			// 트랜젝션 확인된 시간 (mms)
		}
	*/
	dbg.Yellow("DepositCallback --- Action", len(list))
	inf.ClientAddress()

	nowAt := mms.Now()
	for _, item := range list {
		data := chttp.JsonType{
			"name":           item.Name,
			"address":        item.Address,
			"symbol":         item.Symbol,
			"price":          item.Price,
			"hash":           item.Hash,
			"deposit_result": item.DepositResult,
			"timestamp":      item.Timestamp,
		}
		dbg.Yellow(data)

		ack := cnet.POST_JSON_F(
			inf.ClientAddress()+"/v1/coins/deposit/callback",
			nil,
			data,
		)
		if err := ack.Error(); err != nil {
			dbg.RedItalic("deposit.send_fail :", err)
			continue
		}
		item.SendOK(db, nowAt)

	}
}

func sendWithdraw(db mongo.DATABASE, list model.LogWithdrawList) {
	/*
		( 2. 코인 출금 신청 ) 요청으로 코인 출금이 성공/실패 했을 경우 출금 결과 전송
		http://service.server.com:8080/v1/coins/withdraw/callback

		method : post
		content-type : application/json

		전송 데이터
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
	*/
	dbg.Yellow("WithdrawCallback --- Action", len(list))
	inf.ClientAddress()

	nowAt := mms.Now()
	for _, item := range list {
		_ = chttp.JsonType{
			"receipt_code":    item.ReceiptCode,
			"name":            item.Name,
			"from_address":    item.Address,
			"to_address":      item.ToAddress,
			"symbol":          item.Symbol,
			"price":           item.ToPrice,
			"hash":            item.Hash,
			"withdraw_result": item.State,
			"timestamp":       item.Timestamp,
		}
		data := item.AckJson()
		dbg.Yellow(data)

		ack := cnet.POST_JSON_F(
			inf.ClientAddress()+"/v1/coins/withdraw/callback",
			nil,
			data,
		)
		if err := ack.Error(); err != nil {
			dbg.RedItalic("withdraw.send_fail :", err)
			continue
		}
		item.SendOK(db, nowAt)
	} //for
}

func sendWithdrawSELF(db mongo.DATABASE, list model.LogWithdrawSELFList) {
	/*
		( ex. 개인지갑 코인 출금신청 결과 ) 요청으로 코인 출금이 성공/실패 했을 경우 출금 결과 전송
		http://service.server.com:8080/v1/coins/withdraw_self/callback

		method : post
		content-type : application/json

		전송 데이터
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
	*/
	dbg.Yellow("WithdrawSELFCallback --- Action", len(list))
	inf.ClientAddress()

	nowAt := mms.Now()
	for _, item := range list {

		data := item.AckJson()
		dbg.Yellow(data)

		ack := cnet.POST_JSON_F(
			inf.ClientAddress()+"/v1/coins/withdraw_self/callback",
			nil,
			data,
		)
		if err := ack.Error(); err != nil {
			dbg.RedItalic("withdrawSELF.send_fail :", err)
			continue
		}

		item.SendOK(db, nowAt)
	} //for
}

func sendToMaster(db mongo.DATABASE, list model.LogToMasterList) {
	/*
		( 3. 개인지갑 코인 마스터지갑으로 전송시 콜백 ) 회원 개인지갑의 코인이 마스터지갑으로 전송시 이동한 코인량과 현재 개인지갑의 코인 잔액
		http://service.server.com:8080/v1/coins/tomaster/callback

		method : post
		content-type : application/json

		전송 데이터
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
	*/
	dbg.Yellow("ToMasterCallback --- Action")

	nowAt := mms.Now()
	for _, item := range list {
		member := model.LoadMember(db, item.UID)
		data := item.AckJson(member.Coin.Clone())
		dbg.Yellow(data)

		ack := cnet.POST_JSON_F(
			inf.ClientAddress()+"/v1/coins/tomaster/callback",
			nil,
			data,
		)
		if err := ack.Error(); err != nil {
			dbg.RedItalic("sendToMaster.send_fail :", err)
			continue
		}

		item.SendOK(db, nowAt)
	} //for
}

func sendExMaster(db mongo.DATABASE, list model.LogExMasterList) {
	/*
		( 4. 개인지갑이 아닌 기타 외부에서 마스터 지갑으로 코인을 전송시 콜백 ) - 전송자의 주소와 코인종류 및 가격 , 트랜젝션 해시를 콜백함
		http://service.server.com:8080/v1/coins/exmaster/callback

		method : post
		content-type : application/json

		전송 데이터
		{
			"hash" : string,			// 트랜젝션 해시
			"from" : string,			// 전송자의 지갑 주소 (0x...)
			"contract" : string,        // 이더전송이면 "eth" , 토큰 전송이면 컨트랙트 주소
			"symbol" : string,			// 마스터로 전송된 토큰 심볼 (ETH , GDG ...)
			"price" : string,			// 마스터로 전송된 토큰 수량 ("1.333")
			"timestamp" : int			// 전송 확인된 시간 (mms)
		}
	*/
	dbg.Yellow("ExMasterCallback --- Action")
	nowAt := mms.Now()
	for _, item := range list {
		data := item.AckJson()
		dbg.Yellow(data)

		ack := cnet.POST_JSON_F(
			inf.ClientAddress()+"/v1/coins/exmaster/callback",
			nil,
			data,
		)
		if err := ack.Error(); err != nil {
			dbg.RedItalic("sendExMaster.send_fail :", err)
			continue
		}

		item.SendOK(db, nowAt)
	} //for
}

func sendMasterOut(db mongo.DATABASE, list model.TxETHMasterOutList) {
	/*
		( 5. 마스터 계좌에서 외부 계좌로 출금 신청 결과 콜백  )
		http://service.server.com:8080/v1/coins/master_out/callback

		method : post
		content-type : application/json

		전송 데이터
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


	*/
	dbg.Yellow("MasterOutCallback --- Action", len(list))

	nowAt := mms.Now()
	for _, item := range list {
		data := chttp.JsonType{
			"receipt_code": item.ReceiptCode,
			"symbol":       item.Symbol,
			"decimal":      item.Decimal,
			"from_address": inf.Master().Address,
			"to_address":   item.ToAddress,
			"price":        item.ToPrice,
			"gas":          item.Gas,

			"hash":            item.Hash,
			"withdraw_result": item.State,
			"fail_message":    item.FailMessage,
			"timestamp":       item.Timestamp,
		}
		dbg.Yellow(data)

		ack := cnet.POST_JSON_F(
			inf.ClientAddress()+"/v1/coins/master_out/callback",
			nil,
			data,
		)
		if err := ack.Error(); err != nil {
			dbg.RedItalic("master_out.send_fail :", err)
			continue
		}

		item.SendOK(db, nowAt)
	} //for
}
