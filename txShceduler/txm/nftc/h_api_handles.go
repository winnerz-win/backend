package nftc

import (
	"net/http"
	"sync"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/txm/ack"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

var handle = chttp.PContexts{}

func readyHandlers(classic *chttp.Classic) {
	defer func() {
		classic.SetContextHandles(handle)
	}()

	hNftVersion()
	hMasterInfo()
	hTxLog()
	hSaleBuyTry()
	hSalePendingReceipt()
	hTransferTry()
	hTransferCheckPending()
	hOwnerList()
}

func hNftVersion() {
	/*
		Comment : NFT 개발 버전확인
		Method : GET
		Response :
		{
			"success" : true,
			"data" : {
				"NFT" : "v1.0.0",
			}
		}
	*/
	method := chttp.GET
	url := model.NFT + "/version"
	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			chttp.OK(w, mongo.VOID{
				"NFT": "v1.0.0",
			})
		},
	)
}

type cMasterInfo struct {
	NFTContractAddress string `json:"nft_contract_address"`
	NFTOwnerAddress    string `json:"nft_owner_address"`
	NFTOwnerETHPrice   string `json:"nft_owner_eth"`

	NFTDepositMasterAddress string `json:"nft_deposit_master_address"`
}

func hMasterInfo() {
	/*
		Comment : NFT 마스터 주소 및 잔액 조회
		Method : GET

		Response :
		{
			"success" : true,
			"data" : {
				"nft_contract_address" : "0x03B9f58383dF6996C9445d983E56f68fbcC4046A",	// NFT컨트랙트 주소(테스트넷)
				"nft_owner_address" : "0x07e3595dd3662b3b6ce8f929c9c0c651edf921f4",		// NFT발행 오너 주소
				"nft_owner_eth" : "14.7",	// NFT발행 오너의 이더 잔액
				"nft_deposit_master_address" : "0xcc5168c6b85dc650b660b5b93f360a835e590cbc",	// 발행비용 수금 마스터 주소
			}
		}
	*/
	method := chttp.GET
	url := model.NFT + "/master_info"
	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			info := cMasterInfo{
				NFTContractAddress:      nftToken.Contract,
				NFTOwnerAddress:         nftToken.Address,
				NFTOwnerETHPrice:        Sender().CoinPrice(nftToken.Address),
				NFTDepositMasterAddress: depositAddress(),
			}

			chttp.OK(w, info)
		},
	)
}

func hTxLog() {
	/*
		Comment : NFT 트랜젝션 로그 조회
		Method : POST (application/json)
		Param :
		{
			"kind" string			// Transfer(소유권이전), Mint(발행), Burn(소각), 미입력시 모두(ALL)
			"address" string 		// 검색할 주소 from, to --> (미입력시 모든 주소)
			"limit_count" int		// 검색할 최대 갯수 (미입력시 기본값 200개, 최대값 : 200개)
			"last_index" int		// 검색시작 인덱스값 (미입력 또는, 0 입력시 최근 데이타 부터)  --- 내림차순(최신 로그 데이타부터 역순으로 검색),
									// 응답 결과 리스트의 마지막 인덱스를 넣으면 index-1 부터 다음 limit_count 만큼 검색함.

			"token_id" string 		// 토큰 고유 ID (미입력시 모든 ID)
			"token_type" string		// 토큰 타입 (미입력시 모든 타입)
		}
		< 위의 요청 파라미터 값은 모두 생략 가능 합니다. >

		< EX >
		- token_id == 10, kind == Transfer, index < 2191 인 로그(200개) 조회
		{
			"kind" : "Transfer",
			"token_id" : "10"
			"last_index" : 2191,
		}
		-------------------------------------------------------------------
		Response :
		{
			"success" : true,
			"data" : [ log-data ]	//요청한 로그 데이타 배열 (내림차순 정렬)
		}
		< log-data 형식 >
		{
			"index" int				// 페이징 인덱스 값 ( 3345,3344,3343, ...)
			"number" long			// Tx 블럭번호
			"hash" string			// Tx 해시
			"tx_index" long			// Tx 인덱스
			"log_index" long		// Tx-log 인덱스
			"timestamp" long		// Tx 발생 시간 ( 13자리 mms)
			"name" string			// Transfer, Mint, Burn  (로그 종류)
			"from" string			// 전송자 주소 (Mint 일경우는 0x0000000000000000000000000000000000000000(0번지) 주소임)
			"to" string				// 수신자 주소 (Burn 일경우는 0x0000000000000000000000000000000000000000(0번지) 주소임)
			"token_id" string		// NFT 토큰 고유 ID
			"token_type" string		// NFT 토큰 타입
		}

	*/
	type CDATA struct {
		Kind      string `json:"kind,omitempty"` //Transfer, Mint, Burn, All
		Address   string `json:"address,omitempty"`
		TokenID   string `json:"token_id"` //KEY
		TokenType string `json:"token_type"`

		LimitCount int `json:"limit_count,omitempty"` //200
		LastIndex  int `json:"last_index,omitempty"`
	}
	method := chttp.POST
	url := model.NFT + "/txlog"
	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := CDATA{}
			chttp.BindingJSON(req, &cdata)

			if !jmath.IsNum(cdata.TokenID) || !jmath.IsNum(cdata.TokenType) {
				chttp.Fail(w, ack.BadParam)
				return
			}

			cdata.Address = dbg.TrimToLower(cdata.Address)

			model.DB(func(db mongo.DATABASE) {

				cdata.Address = dbg.TrimToLower(cdata.Address)
				cdata.Kind = dbg.TrimToLower(cdata.Kind)
				cdata.TokenID = dbg.TrimToLower(cdata.TokenID)
				cdata.TokenType = dbg.TrimToLower(cdata.TokenType)

				selector := mongo.Bson{}
				switch cdata.Kind {
				case "transfer":
					selector["kind"] = "Transfer"
				case "mint":
					selector["kind"] = "Mint"
				case "burn":
					selector["kind"] = "Burn"
				default:
				} //switch

				if cdata.Address != "" {
					selector["$or"] = []mongo.Bson{
						{"from": cdata.Address},
						{"to": cdata.Address},
					}
				}

				if cdata.TokenID != "" {
					selector["token_id"] = cdata.TokenID
				}
				if cdata.TokenType != "" {
					selector["token_type"] = cdata.TokenType
				}

				const _limitCount = 200

				if cdata.LimitCount > 0 {
					if cdata.LimitCount > _limitCount {
						cdata.LimitCount = _limitCount
					}
				} else {
					cdata.LimitCount = _limitCount
				}

				if cdata.LastIndex > 0 {
					selector["index"] = mongo.Bson{"$lt": cdata.LastIndex}
				}

				list := model.NftTxLogList{}
				db.C(inf.NFTTxLog).Find(selector).
					Sort("-number").
					Limit(cdata.LimitCount).
					All(&list)

				chttp.OK(w, list)

			})
		},
	)
}

func hSalePendingReceipt() {
	/*
		Comment : NFT ReceiptCode 값으로 현재 구매 진행상태 요청
		Method : POST
		-------------------------------------------------------------------
		Params :
		{
			"receipt_code" : string //구매신청때 받은 영수증 코드
		}
		Response :
		{ 	//case 1
			"success" : true,
			"data" : {
				"result_kind" : "deposit"
				"result_data" : {DEPOSTI-DATA}
			}
		}

		{ 	//case 2
			"success" : true,
			"data" : {
				"result_kind" : "try"
				"result_data" : {TRY-DATA}
			}
		}

		{	//case 3
			"success" : true,
			"data" : {
				"result_kind" : "result"
				"result_data" : {RESULT-DATA}
			}
		}

		Response : fail
		{
			"success" : false,
			"data" : {
				"error_code" : 3005,	// 해당 영수증의 데이타를 찾지 못하였음.
				"error_message" : "NFT not found data"
			}
		}

		case 1 > {DEPOSTI-DATA} Format
		{
			...
			"status" : 0 or 1 (0이면 구매대금 처리전 , 1이면 구매대금 pending상태)
			...
		}


		case 2 > {TRY-DATA} Format
		{
			...
			"status" : 0 or 1 (0이면 NFT 발행 처리전 , 1이면 NFT 발행 pending상태)
			...
		}


		case 3 > {RESULT-DATA} Format
		{
			...
			"status" : 200 or 104 (200이면 NFT 발행 성공 , 104이면 NFT 발행 실패)
			...
		}


	*/
	type CDATA struct {
		ReceiptCode string `json:"receipt_code"`
	}
	type RESULT struct {
		ResultKind string      `json:"result_kind"` // deposit, try , result
		ResultData interface{} `json:"result_data"` // NftBuyTry , NftBuyResult
	}
	method := chttp.POST
	url := model.NFT + "/sale/find/receipt_code"
	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := CDATA{}
			chttp.BindingJSON(req, &cdata)

			model.DB(func(db mongo.DATABASE) {
				deposit := model.NftDepositTry{}.GetByReceiptCode(db, cdata.ReceiptCode)
				if deposit.Valid() {
					chttp.OK(w, RESULT{
						ResultKind: "deposit",
						ResultData: deposit,
					})
					return
				}

				try := model.NftBuyTry{}.GetByReceiptCode(db, cdata.ReceiptCode)
				if try.Valid() {
					chttp.OK(w, RESULT{
						ResultKind: "try",
						ResultData: try,
					})
					return
				}

				result := model.NftBuyEnd{}.GetByReceiptCode(db, cdata.ReceiptCode)
				if result.Valid() {
					chttp.OK(w, RESULT{
						ResultKind: "result",
						ResultData: result,
					})
					return
				}

				chttp.Fail(w, ack.NFTNotfoundData)
			})
		},
	)
}

func checTryTokenID(db mongo.DATABASE, tokenId string) bool {
	deposit := model.NftDepositTry{}.GetTokenID(db, tokenId)
	if deposit.Valid() {
		return false
	}
	try := model.NftBuyTry{}.GetTokenID(db, tokenId)
	if try.Valid() {
		return false
	}
	result := model.NftBuyEnd{}.GetTokenID(db, tokenId)
	if result.Valid() {
		if result.IsBurn == false {
			return false
		}
	}
	return true
}

func checkNumber(v string) bool {
	if v == "" {
		return false
	}
	return jmath.IsNum(v)
}

func hSaleBuyTry() {
	/*
		Comment : NFT 상품 구매요청 ( 회원구매(무료/유료), 비회원구매(무료) )
		Method : POST
		-------------------------------------------------------------------
		Params :
		{
			"member_address" : string	// 가입 회원 입금 주소 (NFT발행 대상자 및 구매 비용을 지불할 주소)
			"token_id" : string			// 발행할 token_id (token_id는 고유키값임.)
			"token_type" : string		// 발행할 token_type
			"pay_symbol" : string		// 구매 토큰 심볼 (ETH , GDG(테스트넷은 ERCT))
			"pay_price" : string		// 구매 토큰 가격 (decimal을 제외한 값입니다.)
			"is_pay_free" : bool		// 무료 발행 여부 (true 일경우 pay_symbol,pay_price는 무시.)
			"is_external_address" : bool // member_address 가 회원주소가 아닌 외부주소 인지 여부
			( is_external_address == true 일경우,
				pay_symbol, pay_price, is_pay_free 값은 무시며 구매비용을 무료가 된다. )
		}
		Response :
		{
			"success" : true,
			"data" : {
				"receipt_code" : string // 구매 신청 완료 영수증 코드
			}
		}

		Fail:
		{
			"success" : false,
			"data" : {
				"error_code" : int
				"error_message" : string
			}
		}
		< error_code >
		1001 : 요청 파라미터 오류 ( 유료구매시 pay_price 값이 0이하일 경우 )
		4102 : 존재하지 않는 입금 주소
		4103 : 지원하지 않는 pay_symbol (ETH , GDG(ERCT) 만 지원)
		3003: NFT 구매 비용 부족 ( 구매비용 및 가스비 선 계산 )
				( ETH 일경우 pay_price(ETH) + gas(ETH) )
				( GDG(ERCT) 일경우 pay_price(GDG) + gas(ETH) )
		3006 : token_id 중복 ( 이미 구매 했거나 구매 진행중인 token_id )
		4104 : 이더리움 주소체계가 아님 (0x...) --> 외부주소로 구매신청 했을경우 검사.

	*/
	type CDATA struct {
		IsExternalAddress bool   `json:"is_external_address"`
		MemberAddress     string `json:"member_address"`
		TokenId           string `json:"token_id"`
		TokenType         string `json:"token_type"`
		PaySymbol         string `json:"pay_symbol"`
		PayPrice          string `json:"pay_price"`
		IsPayFree         bool   `json:"is_pay_free"`
	}
	type RESULT struct {
		ReceiptCode string `json:"receipt_code"`
	}
	method := chttp.POST
	url := model.NFT + "/sale/buy_try"
	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := CDATA{}
			chttp.BindingJSON(req, &cdata)
			cdata.MemberAddress = dbg.TrimToLower(cdata.MemberAddress)

			if !checkNumber(cdata.TokenId) || !checkNumber(cdata.TokenType) {
				chttp.Fail(w, ack.BadParam)
				return
			}

			model.DB(func(db mongo.DATABASE) {
				model.NFTLock(func() {
					if checTryTokenID(db, cdata.TokenId) == false {
						chttp.Fail(w, ack.NFTExistTokenID)
						return
					}

					if cdata.IsExternalAddress {
						if ecsx.IsAddress(cdata.MemberAddress) == false {
							chttp.Fail(w, ack.InvalidAddress)
							return
						}
						revData := model.NftRevData{
							User: model.User{
								UID:     0,
								Address: cdata.MemberAddress,
								Name:    "",
							},
							IsExternalAddress: true,
							PayAddress:        cdata.MemberAddress,
							PaySymbol:         "",
							PayPrice:          model.ZERO,
							IsPayFree:         true,

							TokenId:   cdata.TokenId,
							TokenType: cdata.TokenType,
						}

						revTry := model.NftDepositTry{
							NftRevData: revData,
							Snap:       ecsx.GasSnapShot{},
						}

						receiptCode, err := revTry.InsertTryDB(db)
						if err != nil {
							chttp.Fail(w, ack.NFTExistTokenID)
							return
						}

						chttp.OK(w, RESULT{ReceiptCode: receiptCode})
						return
					}

					member := model.LoadMemberAddress(db, cdata.MemberAddress)
					if member.Valid() == false {
						chttp.Fail(w, ack.NotFoundAddress)
						return
					}

					payAddress := cdata.MemberAddress
					payPrivate := member.PrivateKey()

					snap := ecsx.GasSnapShot{}
					if cdata.IsPayFree == false {
						token := inf.TokenList().GetSymbol(cdata.PaySymbol)
						if token.Valid() == false {
							chttp.Fail(w, ack.NotFoundSymbol)
							return
						}

						//GDG로 구매시는 입금마스터가 구매 대행을 한다.
						if token.Symbol != model.ETH {
							payAddress = inf.Master().Address
							payPrivate = inf.Master().PrivateKey
						}

						if jmath.IsUnderZero(cdata.PayPrice) {
							chttp.Fail(w, ack.BadParam)
							return
						}

						ethPrice := Sender().CoinPrice(payAddress)
						if token.Symbol == model.ETH {

							if jmath.CMP(ethPrice, cdata.PayPrice) <= 0 {
								chttp.Fail(w, ack.NFTBuyPrice)
								return
							}
							ntx, _ := Sender().EthTransferNTX(
								payPrivate,
								depositAddress(),
								ecsx.ETHToWei(cdata.PayPrice),
								gasSpeed,
								nil,
							)
							gasETH := ntx.GasFeeETH()
							needETH := jmath.ADD(gasETH, cdata.PayPrice)
							if jmath.CMP(ethPrice, needETH) < 0 {
								chttp.Fail(w, ack.NFTBuyPrice)
								return
							}

							snap = ntx.SnapShot()

						} else {
							if jmath.IsUnderZero(ethPrice) {
								chttp.Fail(w, ack.NFTBuyPrice)
								return
							}
							tkPrice := Sender().TokenPrice(payAddress, token.Contract, token.Decimal)
							if jmath.CMP(tkPrice, cdata.PayPrice) < 0 {
								chttp.Fail(w, ack.NFTBuyPrice)
								return
							}
							ts := Sender().TSender(token.Contract)
							ntx, _ := ts.TransferFuncNTX(
								payPrivate,
								ecsx.TransferPadBytes(
									depositAddress(),
									ecsx.TokenToWei(cdata.PayPrice, token.Decimal),
								),
								"0",
								gasSpeed,
								nil,
							)
							gasETH := ntx.GasFeeETH()
							if jmath.CMP(ethPrice, gasETH) < 0 {
								chttp.Fail(w, ack.NFTBuyPrice)
								return
							}

							snap = ntx.SnapShot()

						}
					} else {
						cdata.PaySymbol = ""
						cdata.PayPrice = model.ZERO
					}

					revData := model.NftRevData{
						User:       member.User,
						PayAddress: payAddress,
						PaySymbol:  cdata.PaySymbol,
						PayPrice:   cdata.PayPrice,
						IsPayFree:  cdata.IsPayFree,

						TokenId:   cdata.TokenId,
						TokenType: cdata.TokenType,
					}

					revTry := model.NftDepositTry{
						NftRevData: revData,
						Snap:       snap,
					}

					receiptCode, err := revTry.InsertTryDB(db)
					if err != nil {
						chttp.Fail(w, ack.NFTExistTokenID)
						return
					}

					chttp.OK(w, RESULT{ReceiptCode: receiptCode})
				}) //NFTLock

			})
		},
	)
}

var (
	transferMu sync.Mutex
)

func hTransferTry() {
	/*
		Comment : 회원의 NFT토큰을 다른주소(다른회원/비회원)로 소유권 이전
		Method : POST
		-------------------------------------------------------------------
		Param :
		{
			"member_address" : string	// 회원의 주소 (전송자)
			"to_address" : string		// 받을 주소 (수신자 , 회원/비회원)
			"token_id" : string			// NFT토큰의 ID
		}
		Response : success
		{
			"success" : true,
			"data" : {
				"token_id" : string		// 요청한 NFT토큰의 ID
				"hash" : string			// NFT토큰을 발송한 트랜젝션 해시값
			}
		}
		Response : fail
		{
			"success" : false,
			"data" : {
				"error_code" : int,
				"error_message" : string
			}
		}
		< error_code (실패 응답 코드) >
		1001 : 요청파라미터 오류 (member_addres == to_address 가 같은 값일경우.)
		4102 : 존재하지 않는 입금주소 (member_address 가 회원 주소가 아님)
		4104 : to_address가 이더리움 주소체계가 아님 (0x...)
		3001 : tokenId를 찾을수 없음.
		3002 : ETH잔액이 0이하
		3003 : NFT 전송 가스비용 부족
		3005 : tokenId를 찾았으나 소유자(last_owner)가 member_address가 아닐경우
			또는 발급된 tokenId가 소각된 경우.
		3007 : 전송 transaction 실패
		3008 : 전송 transfer 실패

	*/
	type CDATA struct {
		MemberAddress string `json:"member_address"`
		ToAddress     string `json:"to_address"`
		TokenId       string `json:"token_id"`
	}
	type RESULT struct {
		TokenId string `json:"token_id"`
		Hash    string `json:"hash"`
	}

	method := chttp.POST
	url := model.NFT + "/transfer/change_owner"
	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := CDATA{}
			chttp.BindingJSON(req, &cdata)

			if !checkNumber(cdata.TokenId) {
				chttp.Fail(w, ack.BadParam)
				return
			}

			cdata.MemberAddress = dbg.TrimToLower(cdata.MemberAddress)
			cdata.ToAddress = dbg.TrimToLower(cdata.ToAddress)

			if cdata.MemberAddress == cdata.ToAddress {
				chttp.Fail(w, ack.BadParam)
				return
			}

			if !ecsx.IsAddress(cdata.ToAddress) {
				chttp.Fail(w, ack.InvalidAddress)
				return
			}

			model.DB(func(db mongo.DATABASE) {
				defer transferMu.Unlock()
				transferMu.Lock()

				member := model.LoadMemberAddress(db, cdata.MemberAddress)
				if member.Valid() == false {
					chttp.Fail(w, ack.NotFoundAddress)
					return
				}

				buyEnd := model.NftBuyEnd{}.GetTokenID(db, cdata.TokenId)
				if !buyEnd.Valid() {
					chttp.Fail(w, ack.NFTNotfoundTokenID)
					return
				}
				if buyEnd.IsBurn {
					chttp.Fail(w, ack.NFTNotfoundData)
					return
				}
				if buyEnd.LastOwner != member.Address {
					chttp.Fail(w, ack.NFTNotfoundData)
					return
				}

				ntx, err := NFT{}.TransferFromNTX(
					member.PrivateKey(),
					cdata.ToAddress,
					cdata.TokenId,
				)
				if err != nil {
					chttp.Fail(w, ack.NFTTransactionFail)
					return
				}

				memberETH := Sender().CoinPrice(member.Address)
				if jmath.CMP(memberETH, ntx.GasFeeETH()) < 0 {
					chttp.Fail(w, ack.NFTBuyPrice)
					return
				}

				hash, err := NFT{}.TransferFromSEND(ntx)
				if err != nil {
					chttp.Fail(w, ack.NFTTTransferFail)
					return
				}

				isMemberTo := false
				toMember := model.LoadMemberAddress(db, cdata.ToAddress)
				if toMember.Valid() {
					isMemberTo = true
				}

				try := model.NftTransferTry{
					User:       member.User,
					TokenId:    cdata.TokenId,
					ToAddress:  cdata.ToAddress,
					IsMemberTo: isMemberTo,
					Hash:       hash,
				}
				try.InsertDB(db)

				chttp.OK(w, RESULT{
					TokenId: cdata.TokenId,
					Hash:    hash,
				})

			})
		},
	)
}

func hTransferCheckPending() {
	/*
		Comment : 소유권 이전중인 토큰의 전송여부 체크
		Method : POST
		-------------------------------------------------------------------
		Param :
		{
			"hash" : string			// NFT토큰을 전송한 트랜잭션 해시값
		}
		Response : success
		{ 	//case 1
			"success" : true,
			"data" : {
				"result_kind" : "pending"	//현재 처리중임.
				"result_data" : {PENDING-DATA}
			}
		}

		{ 	//case 2
			"success" : true,
			"data" : {
				"result_kind" : "result"
				"result_data" : {RESULT-DATA}
			}
		}
		< {RESULT-DATA} format >
		{
			...
			"status" : 200 or 104 (200이면 전송 성공 , 104이면 전송 실패)
			...
		}

		Response : fail
		{
			"success" : false,
			"data" : {
				"error_code" : int,
				"error_message" : string
			}
		}
		< error_code (실패 응답 코드) >
		3005 : 요청한 hash로 전송중인 token_id를 찾지 못함


	*/
	type CDATA struct {
		Hash string `json:"hash"`
	}
	method := chttp.POST
	url := model.NFT + "/transfer/check_pending"
	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := CDATA{}
			chttp.BindingJSON(req, &cdata)

			cdata.Hash = dbg.TrimToLower(cdata.Hash)

			model.DB(func(db mongo.DATABASE) {
				selector := mongo.Bson{"hash": cdata.Hash}

				try := model.NftTransferTry{}
				db.C(inf.NFTTransferTry).Find(selector).One(&try)
				if try.Valid() {
					chttp.OK(w, chttp.JsonType{
						"result_kind": "pending",
						"result_data": try,
					})
					return
				}

				result := model.NftTransferEnd{}
				db.C(inf.NFTTransferEnd).Find(selector).One(&result)
				if result.Valid() {
					chttp.OK(w, chttp.JsonType{
						"result_kind": "result",
						"result_data": result,
					})
					return
				}

				chttp.Fail(w, ack.NFTNotfoundData)

			})
		},
	)

}

func hOwnerList() {
	/*
		Comment : 회원/비회원(회사에서 발급) 이 소유하고 있는 NFT상품 리스트 요청
		Method : POST
		-------------------------------------------------------------------
		Param :
		{
			"uid" : long				// 회원의 고유ID로 검색
			"member_address" : string	//
			"name" : string
		}
		Response : success
		{
			"success" : true,
			"data" : {
				"receipt_code" : string,	// NFT 구매 신청 영수증( "token_7573871901640588701_nft" , "eth_5593420227436337062_nft" )
			}
		}
		Response : fail
		{
			"success" : false,
			"data" : {
				"error_code" : int,
				"error_message" : string
			}
		}
		< error_code (실패 응답 코드) >
		4102 : 존재하지 않는 입금주소
		3001 : tokenId를 찾을수 없음.
		3002 : ETH잔액이 0이하
		3003 : NFT 구매 비용 부족 (ETH,GDG의 구매 비용 부족)
		3004 : NFT 구매 시도 실패 (DB 작업 실패) ---> tokenId를 다른 회원이 먼저 선점하였을 경우.

	*/
	type CDATA struct {
		UID           int64  `json:"uid"`
		MemberAddress string `json:"member_address"`
		Name          string `json:"name"`
	}

	method := chttp.POST
	url := model.NFT + "/owner/list"
	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := CDATA{}
			chttp.BindingJSON(req, &cdata)

			cdata.MemberAddress = dbg.TrimToLower(cdata.MemberAddress)

			model.DB(func(db mongo.DATABASE) {

				selector := mongo.Bson{}
				if cdata.UID != 0 {
					selector["uid"] = cdata.UID
				} else if cdata.MemberAddress != "" {
					selector["address"] = cdata.MemberAddress
				} else if cdata.Name != "" {
					selector["name"] = cdata.Name
				} else {
					chttp.Fail(w, ack.BadParam)
					return
				}

				list := model.NftBuyEndList{}
				db.C(inf.NFTBuyEnd).Find(selector).All(&list)
				chttp.OK(w, list)
			})
		},
	)

}
