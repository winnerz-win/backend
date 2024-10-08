package nft_winners

import (
	"jtools/cloud/ebcm"
	"jtools/jmath"
	"jtools/unix"
	"net/http"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/database/mongo/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/jnet/doc"
	"txscheduler/nft_winners/nwtypes"
	"txscheduler/txm/ack"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func init() {
	hNftCalcETHPrice()
	hNftUserSaleETH()
}

type AckCalcETHPrice struct {
	IsTransferAllow bool   `json:"is_transfer_allow"`
	FailMessage     string `json:"fail_message,omitempty"`

	FromETHPrice string `json:"from_eth_price"`

	EstimateMinGasPrice string `json:"estimate_min_gas_price"`
	EstimateMinPayPrice string `json:"estimate_min_pay_price"`

	EstimateNeedGasPrice string `json:"estimate_need_gas_price"`
	EstimateNeedPayPrice string `json:"estimate_need_pay_price"`
}

func (AckCalcETHPrice) TagString() []string {
	return []string{
		"is_transfer_allow", "가능 여부 (false 이면 전송 불가 : (가스비+전송수량의합)이 전송자의 ETH보유량 보다 클경우)",
		"fail_message,omitempty", "실패 일경우 사유",

		"from_eth_price", "전송 주소의 실제 ETH 잔액",

		"estimate_min_gas_price", "[최소예측](is_transfer_allow = true일 경우)해당 트랜잭션을 실행하는데 필요한 가스비용",
		"estimate_min_pay_price", "[최소예측](is_transfer_allow = true일 경우)해당 트랜잭션을 실행하는데 필요한 가스비용 + 전송수량 총합",

		"estimate_need_gas_price", "[실제 필요 예측](is_transfer_allow = true일 경우)해당 트랜잭션을 실행하는데 필요한 가스비용",
		"estimate_need_pay_price", "[실제 필요 예측](is_transfer_allow = true일 경우)해당 트랜잭션을 실행하는데 필요한 가스비용 + 전송수량 총합",
	}
}

type RecipientData struct {
	Address string `json:"address"`
	Price   string `json:"price"`
}

func (RecipientData) TagString() []string {
	return []string{
		"address", "지갑 주소",
		"price", "가격",
	}
}

type ReqCalcETHPrice struct {
	FromAddress string `json:"from_address"`

	Recipients []RecipientData `json:"recipients"`

	pairs     []any
	total_wei string
}

func (ReqCalcETHPrice) TagString() []string {
	return []string{
		"from_address", "전송자 주소(회원 주소)",
		"recipients", "받는사람들 정보(주소,수량) 배열",

		"pairs", "안쓰는 필드입니다.(안씀)",
		"total_wei", "안쓰는 필드입니다.(안씀)",
	}
}

func (my *ReqCalcETHPrice) Valid() chttp.CError {
	if !ebcm.IsAddressP(&my.FromAddress) {
		return ack.InvalidAddress
	}
	size := len(my.Recipients)
	if size == 0 {
		return ack.BadParam
	}

	my.total_wei = model.ZERO

	for i := range my.Recipients {
		if !ebcm.IsAddressP(&my.Recipients[i].Address) {
			return ack.InvalidAddress
		}
		if jmath.CMP(my.Recipients[i].Price, 0) <= 0 {
			return ack.NFT_Param_price
		}

		wei := ebcm.ETHToWei(my.Recipients[i].Price)
		my.pairs = append(my.pairs, my.Recipients[i].Address, wei)
		my.total_wei = jmath.ADD(my.total_wei, wei)
	}

	return nil
}

func hNftCalcETHPrice() {
	method := chttp.POST
	url := model.V1 + "/nft/estimate/eth_price"
	Doc().Comment("[WINNERZ] ETH 멀티 전송시(ETH민팅 / 유저거래) 가스 비용 예측").
		Method(method).URL(url).
		JParam(ReqCalcETHPrice{}, ReqCalcETHPrice{}.TagString()...).
		JAckOK(AckCalcETHPrice{}, AckCalcETHPrice{}.TagString()...).
		ETCVAL(RecipientData{}, RecipientData{}.TagString()...).
		JAckError(ack.InvalidAddress, "요청 주소의 형식 오류").
		JAckError(ack.BadParam, "recipients 필드가 0개 인경우").
		JAckError(ack.NFT_Param_price, "").
		JAckError(ack.NotFoundAddress, "from_address가 회원 주소가 아닌경우").
		JAckError(ack.NFT_RPC_TIMEOUT, "블록체인 RPC노드의 일시적 장애").
		Apply()

	_help_estimate_gas()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := chttp.ParseRequestJson[ReqCalcETHPrice](req)
			if ack_err := cdata.Valid(); ack_err != nil {
				chttp.Fail(w, ack_err)
				return
			}

			model.DB(func(db mongo.DATABASE) {
				member := model.LoadMemberAddress(db, cdata.FromAddress)
				if !member.Valid() {
					chttp.Fail(w, ack.NotFoundAddress)
					return
				}

				finder := GetSender()
				if finder == nil {
					chttp.Fail(w, ack.NFT_RPC_TIMEOUT)
					return
				}

				from_wei := finder.Balance(cdata.FromAddress)
				if jmath.CMP(from_wei, cdata.total_wei) <= 0 {
					chttp.OK(w,
						AckCalcETHPrice{
							IsTransferAllow:     false,
							FailMessage:         "회원의 ETH잔고 부족",
							FromETHPrice:        ebcm.WeiToETH(from_wei),
							EstimateMinGasPrice: model.ZERO,
							EstimateMinPayPrice: model.ZERO,
						},
					)
					return
				}

				tx_min_gas_wei, tx_need_gas_wei, err := EstimateTxFee(
					db,
					"[M_TRANS_API]",
					finder,
					cdata.FromAddress,
					nft_config.NftContract,
					getPadBytes_MultiTransferETH(
						cdata.pairs...,
					),
					cdata.total_wei,
					LIMIT_TAG_M_TRANS,
				)

				EstimateMinGasPrice := model.ZERO
				EstimateMinPayPrice := model.ZERO

				EstimateNeedGasPrice := model.ZERO
				EstimateNeedPayPrice := model.ZERO

				is_transfer_allow := err == nil
				fail_message := ""
				if is_transfer_allow {
					EstimateMinGasPrice = ebcm.WeiToETH(tx_min_gas_wei)
					EstimateMinPayPrice = jmath.ADD(
						ebcm.WeiToETH(cdata.total_wei),
						EstimateMinGasPrice,
					)

					EstimateNeedGasPrice = ebcm.WeiToETH(tx_need_gas_wei)
					EstimateNeedPayPrice = jmath.ADD(
						ebcm.WeiToETH(cdata.total_wei),
						EstimateNeedGasPrice,
					)
				} else {
					fail_message = dbg.Cat("RPC_RESULT:", err)
				}

				r := AckCalcETHPrice{
					IsTransferAllow: is_transfer_allow,
					FailMessage:     fail_message,
					FromETHPrice:    ebcm.WeiToETH(from_wei),

					EstimateMinGasPrice: EstimateMinGasPrice,
					EstimateMinPayPrice: EstimateMinPayPrice,

					EstimateNeedGasPrice: EstimateNeedGasPrice,
					EstimateNeedPayPrice: EstimateNeedPayPrice,
				}
				if jmath.CMP(r.FromETHPrice, r.EstimateNeedPayPrice) < 0 {
					r.IsTransferAllow = false
				}

				chttp.OK(w, r)

			})
		},
	)
}

///////////////////////////////////////////////////////////////////////////

type AckNftUserSale struct {
	ReceiptCode nwtypes.RECEIPT_CODE `json:"receipt_code"`
	TokenId     string               `json:"token_id"`
}

func (AckNftUserSale) TagString() []string {
	return []string{
		"receipt_code", "영수증 코드 (nwtypes.RECEIPT_CODE == string)",
		"token_id", "요청한 토큰ID",
	}
}

type ReqNftUserSale struct {
	//PaySymbol     string `json:"pay_symbol"` //ETH
	SellAddress string `json:"sell_address"` //판매자
	BuyAddress  string `json:"buy_address"`  //구매자

	TokenId string `json:"token_id"`

	SellPrice string `json:"sell_price"` //판매자 금액

	BenefitAddress string `json:"benefit_address"`
	BenefitPrice   string `json:"benefit_price"` //베네핏 금액
}

func (ReqNftUserSale) TagString() []string {
	return []string{
		"sell_address", "판매자 (회원)",
		"buy_address", "구매자 (회원)",

		"token_id", "판매자가 소유한 토큰ID",

		"sell_price", "(구매자가 판매자에게 지불할)판매금액",

		"benefit_address", "거래 수수료 받을 지갑",
		"benefit_price", "거래 수수료 금액",
	}
}

func (my ReqNftUserSale) BenefitWEI() string {
	return ebcm.ETHToWei(my.BenefitPrice)
}

func (my ReqNftUserSale) SellWEI() string {
	return ebcm.ETHToWei(my.SellPrice)
}

func (my ReqNftUserSale) PayPriceAllETH_WEI() string {
	sum := "0"
	sum = jmath.ADD(sum, my.SellWEI())
	sum = jmath.ADD(sum, my.BenefitWEI())
	return sum
}

func (my *ReqNftUserSale) Valid() chttp.CError {
	if !ebcm.IsAddressP(&my.SellAddress) {
		return ack.InvalidAddress
	}
	if !ebcm.IsAddressP(&my.BuyAddress) {
		return ack.InvalidAddress
	}
	if !ebcm.IsAddressP(&my.BenefitAddress) {
		return ack.InvalidAddress
	}
	return nil
}

func hNftUserSaleETH() {
	method := chttp.POST
	url := model.V1 + "/nft/sale/eth"
	Doc().Comment("[WINNERZ] NFT 유저간 거래 (ETH지불)").
		Method(method).URL(url).
		JParam(ReqNftUserSale{}, ReqNftUserSale{}.TagString()...).
		JAckOK(AckNftUserSale{}, AckNftUserSale{}.TagString()...).
		JAckError(ack.InvalidAddress, "요청 주소의 형식 오류").
		JAckError(ack.NotFoundAddress, "sell_address, buy_address가 회원 주소가 아닐경우").
		JAckError(ack.NFT_InvalidOwner, "토큰ID의 소유자가 sell_address가 아닌경우 ").
		JAckError(ack.NFT_InvalidPayPrice, "buy_address(구매자)의 ETH 수량 부족").
		JAckError(ack.NFT_NeedETHPrice, "buy_address(구매자)의 ETH 수량 부족").
		Apply()

	_help_userSale_ETH()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := chttp.ParseRequestJson[ReqNftUserSale](req)
			if ack_err := cdata.Valid(); ack_err != nil {
				chttp.Fail(w, ack_err)
				return
			}

			model.DB(func(db mongo.DATABASE) {
				buy_member := model.LoadMemberAddress(db, cdata.BuyAddress)
				if !buy_member.Valid() {
					chttp.Fail(w, ack.NotFoundAddress)
					return
				}
				sell_member := model.LoadMemberAddress(db, cdata.SellAddress)
				if !sell_member.Valid() {
					chttp.Fail(w, ack.NotFoundAddress)
					return
				}

				is_owner := nwtypes.NftTokenIDS{}.IsOwner(db, cdata.TokenId, cdata.SellAddress)
				if !is_owner {
					chttp.Fail(w, ack.NFT_InvalidOwner)
					return
				}

				sender := GetSender()
				on_pay_wei := sender.Balance(buy_member.Address)

				tx_pay_wei := cdata.PayPriceAllETH_WEI()
				_, tx_gas_wei, err := EstimateTxFee(
					db,
					"[M_TRANS_SALE]",
					sender,
					buy_member.Address,
					nft_config.NftContract,
					getPadBytes_MultiTransferETH(
						sell_member.Address,
						cdata.SellWEI(),
						cdata.BenefitAddress,
						cdata.BenefitWEI(),
					),
					tx_pay_wei,
					LIMIT_TAG_M_TRANS,
				)
				if err != nil {
					chttp.Fail(w, ack.NFT_InvalidPayPrice, err)
					return
				}

				estimate_wei := jmath.ADD(tx_pay_wei, tx_gas_wei)
				if jmath.CMP(on_pay_wei, estimate_wei) < 0 {
					chttp.Fail(w, ack.NFT_NeedETHPrice)
					return
				}

				nwtypes.NftTokenIDS{}.UpdatePending(db, cdata.TokenId, true)

				token_list := inf.TokenList()
				token_info := token_list.GetSymbol("ETH")

				receipt_code := nwtypes.GetReceiptCode(nwtypes.RC_SALE_COIN, cdata.TokenId)
				user_try := nwtypes.NftUserTry{
					ReceiptCode: receipt_code,

					DATA_SEQ: nwtypes.MAKE_DATA_SEQ(
						nwtypes.MULTI_TRANSFER,
						nwtypes.DataMultiTransferTry{
							ReceiptCode:  receipt_code,
							PayTokenInfo: token_info,
							PayFrom: nwtypes.WalletInfo{
								Address: buy_member.Address,
								UID:     buy_member.UID,
								Name:    buy_member.Name,
							},
							PriceToInfos: []nwtypes.PriceToInfo{
								{
									InfoKind: nwtypes.InfoKindMember,
									To: nwtypes.WalletInfo{
										Address: sell_member.Address,
										UID:     sell_member.UID,
										Name:    sell_member.Name,
									},
									Price: cdata.SellPrice,
								},
								{
									InfoKind: nwtypes.InfoKindBenefit,
									To: nwtypes.WalletInfo{
										Address: cdata.BenefitAddress,
									},
									Price: cdata.BenefitPrice,
								},
							},
						},

						nwtypes.NFT_TRANSFER,
						nwtypes.DataNftTransfer{
							ReceiptCode: receipt_code,
							From: nwtypes.WalletInfo{
								Address: sell_member.Address,
								UID:     sell_member.UID,
								Name:    sell_member.Name,
							},
							To: nwtypes.WalletInfo{
								Address: buy_member.Address,
								UID:     buy_member.UID,
								Name:    buy_member.Name,
							},
							NftInfo: nft_config.NftInfo(cdata.TokenId),
						},
					),

					TimeTryAt: unix.Now(),
				}

				user_try.InsertDB(db)

				chttp.OK(w, AckNftUserSale{
					ReceiptCode: receipt_code,
					TokenId:     cdata.TokenId,
				})

			})

		},
	)
}

func _help_userSale_ETH() {
	Doc().Message(`
	<cc_purple>( 회원간의 ETH로 NFT 거래 )</cc_purple>

		<요청>
		post  `+doc_host_url+`/v1/nft/sale/eth
		{
			"sell_address" : "0x05b93b0feeb9f60a599ba4b4c76262c22e837579",	//판매자
			"sell_price" : "0.01",	//판매 금액

			"buy_address" : "0xe129243a027b25d813aced72fe34b22f5fc4bb20",	//구매자

			"token_id" : "9999990000002",

			"benefit_address" : "0x8CE5bb2013887eD586e6a87211aa126453368b7A",
			"benefit_price" : "0.01"	//거래 수수료
		}

		=> 구매자는 판매금액+거래수수료+가스비용의 ETH를 보유해야 한다.

		<성공 응답>
		{
			"success" : true,
			"data" : {
				"receipt_code" : "nft_sale_coin_[9999990000002]_d003e9fdee251276600f7dd67ebd37c373bdd8f8c151ef7d417f8ca4f84e16",
				"token_id" : "9999990000002"
			}
		}

		< 스케줄러 처리 로직 >
		1. [구매자]가 판매자+베네핏 에게 요청한 수량의 ETH 전송 (판매금액 + 거래 수수료 + 가스비 부담 : 구매자)
		2. 1번성공시 [판매자]가 [구매자]에게 NFT(9999990000002) 민팅. (가스비 부담 : 판매자)


		< 모든 거래 과정이 끝나면 서비스 서버로 거래 결과를 콜백 >
		post {{서비스서버}}`+URL_NFTS_SALE_CALLBACK+`

		{
			"data": {

				// NFT 소유권 정보 
				"nft_transfer_info": {
					"from": {	//판매자
						"address": "0x05b93b0feeb9f60a599ba4b4c76262c22e837579",
						"is_master": false,
						"name": "MmaTestnetUser01",
						"uid": 1001
					},
					"hash": "0xca0a2b516a40b333cf2653cd95bf8457d2b44eb5ef118fb3f431dd70f0186ee9",	//NFT 전송 HASH
					"nft_info": {
						"contract": "0x1faa080c0e0c7b94d8571d236b182f15e0c1742a",
						"symbol": "WNZ",
						"token_id": "9999990000002"
					},
					"to": {	//구매자
						"address": "0xe129243a027b25d813aced72fe34b22f5fc4bb20",
						"is_master": false,
						"name": "mma_testnet_user_02",
						"uid": 1002
					},
					"tx_gas_price": "0.00800025632308631" //가스비용 : 판매자 부담
				},

				//대금 지불 정보
				"pay_info": {
					"benefit_address": "0x8ce5bb2013887ed586e6a87211aa126453368b7a",	//베네핏주소
					"benefit_fee_price": "0.01",										//베네핏 수수료

					"hash": "0xbe2aaa319fbd1575689a79c133bf077dff8bf6073ee10380ee2295580dbbc892",  //TX-HASH
					"pay_from": {	//구매자 정보 ( 대금 납부자 )
						"address": "0xe129243a027b25d813aced72fe34b22f5fc4bb20",
						"is_master": false,
						"name": "mma_testnet_user_02",
						"uid": 1002
					},
					"pay_token_info": {
						"contract": "eth",	
						"decimal": "18",
						"is_coin": true,
						"mainnet": false,
						"symbol": "ETH"    //거래 대금은 ETH
					},
					"tx_gas_price": "0.006687872499028728", //가스비용 : 구매자 부담

					"user_address": "0x05b93b0feeb9f60a599ba4b4c76262c22e837579",	//판매자 주소
					"user_name": "MmaTestnetUser01",
					"user_pay_price": "0.01",										//판매 가격 
					"user_uid": 1001
				}
			},
			"insert_at": 1683114565,
			"is_send": false,
			"receipt_code": "nft_sale_coin_[9999990000002]_d003e9fdee251276600f7dd67ebd37c373bdd8f8c151ef7d417f8ca4f84e16",
			"result_type": "user_sale",
			"send_at": 0
		}
		
	`, doc.Blue)
}
