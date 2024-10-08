package nft_winners

import (
	"jtools/cloud/ebcm"
	"jtools/jmath"
	"jtools/unix"
	"net/http"
	"strings"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/jnet/doc"
	"txscheduler/nft_winners/nwtypes"
	"txscheduler/txm/ack"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func init() {
	hNftMintETH()
}

type AckNftMint struct {
	ReceiptCode nwtypes.RECEIPT_CODE `json:"receipt_code"`
	TokenId     string               `json:"token_id"`
}

func (AckNftMint) TagString() []string {
	return []string{
		"receipt_code", "영수증 코드 (nwtypes.RECEIPT_CODE == string)",
		"token_id", "요청한 토큰ID",
	}
}

///////////////////////////////////////////////////////////////////////////////////////

type ReqNftMintETH struct {
	PayerAddress string `json:"payer_address"` // ETH 로 민팅할 지갑 주소 (회원)
	OwnerAddress string `json:"owner_address"` // 소유자 지갑 주소 (회원)
	isSelf       bool

	TokenId string `json:"token_id"`

	PlatformAddress string `json:"platform_address"`
	PlatformPrice   string `json:"platform_price"`

	BenefitAddress string `json:"benefit_address"`
	BenefitPrice   string `json:"benefit_price"`
}

func (ReqNftMintETH) TagString() []string {
	return []string{
		//"pay_symbol", "구매 토큰 심볼 : ETH / WNZ(TEST) , 구매 토큰에 따라 처리방식이 분기.",

		"payer_address", "발급 요청자 지갑 주소 (회원 지갑 주소, ETH구매자)",
		"owner_address", "최초 소유자 지갑 주소 (회원 지갑 주소, payer_address == owner_address 자체 발급)",

		"token_id", "토큰ID",

		"platform_address", "수수료 받을 플랫폼 지갑 주소",
		"platform_price", "플랫폼 주소로 보낼 수량",

		"benefit_address", "수수료 받을 베네핏 지갑 주소",
		"benefit_price", "베네핏 주소로 보낼 수량",
	}
}

func (my ReqNftMintETH) BenefitWEI() string {
	return ebcm.ETHToWei(my.BenefitPrice)
}

func (my ReqNftMintETH) PlatformWEI() string {
	return ebcm.ETHToWei(my.PlatformPrice)
}

func (my ReqNftMintETH) PayPriceAllETH_WEI() string {
	sum := "0"
	sum = jmath.ADD(sum, my.PlatformWEI())
	sum = jmath.ADD(sum, my.BenefitWEI())
	return sum
}

func (my *ReqNftMintETH) Valid() chttp.CError {
	if !ebcm.IsAddressP(&my.PayerAddress) {
		return ack.InvalidAddress
	}
	if !ebcm.IsAddressP(&my.OwnerAddress) {
		return ack.InvalidAddress
	}
	if my.PayerAddress == my.OwnerAddress {
		my.isSelf = true
	}

	if !ebcm.IsAddressP(&my.PlatformAddress) {
		return ack.InvalidAddress
	}

	if !ebcm.IsAddressP(&my.BenefitAddress) {
		return ack.InvalidAddress
	}

	my.TokenId = strings.TrimSpace(my.TokenId)
	if !jmath.IsNum(my.TokenId) {
		return ack.NFT_TokenId_Format
	} else if jmath.CMP(my.TokenId, 0) < 0 {
		return ack.NFT_TokenId_Format
	}

	return nil
}
func (my ReqNftMintETH) MintKind() nwtypes.MintKind {
	if my.isSelf {
		return nwtypes.MintKindSelf
	}
	return nwtypes.MintKindGift
}

///////////////////////////////////////////////////////////////////////////////////////

func memberToWalletInfo(member model.Member) nwtypes.WalletInfo {
	return nwtypes.WalletInfo{
		IsMaster: false,
		Address:  member.Address,
		Name:     member.Name,
		UID:      member.UID,
	}
}

func hNftMintETH() {
	method := chttp.POST
	url := model.V1 + "/nft/mint/eth"
	Doc().Comment("[WINNERZ] ETH로 NFT 민팅 요청 ( 선물하기 )").
		Method(method).URL(url).
		JParam(ReqNftMintETH{}, ReqNftMintETH{}.TagString()...).
		JAckOK(AckNftMint{}, AckNftMint{}.TagString()...).
		//JAckError(ack.NotFoundSymbol, "pay_symbol 값이 유효하지 않음").
		JAckError(ack.InvalidAddress, "payer_address, owner_address, platform_address, benefit_address 형식 오류").
		JAckError(ack.NFT_NotfoundPayer, "payer_address가 회원 주소가 아닐경우").
		JAckError(ack.NFT_NotfoundOwner, "owner_address가 회원 주소가 아닐경우").
		JAckError(ack.NFT_InvalidPayPrice, "ETH 금액 부족").
		JAckError(ack.NFT_NeedETHPrice, "ETH 금액 부족").
		JAckError(ack.NFT_TokenId_Format, "토큰ID 형식 오류.(숫자가 아니거나 0미만일 경우)").
		JAckError(ack.NFTExistTokenID, "요청한 토큰ID가 이미 존재하거나 민팅 진행 중일 경우").
		Apply()

	_help_mint_ETH()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := chttp.ParseRequestJson[ReqNftMintETH](req)
			if ack_err := cdata.Valid(); ack_err != nil {
				chttp.Fail(w, ack_err)
				return
			}

			token_list := inf.TokenList()
			token_info := token_list.GetSymbol(model.ETH)

			model.DB(func(db mongo.DATABASE) {
				payer := model.LoadMemberAddress(db, cdata.PayerAddress)
				if !payer.Valid() {
					chttp.Fail(w, ack.NFT_NotfoundPayer)
					return
				}
				owner := model.LoadMemberAddress(db, cdata.OwnerAddress)
				if !owner.Valid() {
					chttp.Fail(w, ack.NFT_NotfoundOwner)
					return
				}

				is_tokenId_rev := nwtypes.NftTokenIDS{}.Reserve(db, cdata.TokenId, func() bool {
					sender := GetSender()
					on_pay_wei := sender.Balance(payer.Address)
					tx_pay_wei := cdata.PayPriceAllETH_WEI()
					_, tx_gas_wei, err := EstimateTxFee(
						db,
						"[MINT_ETH]",
						sender,
						payer.Address,
						nft_config.NftContract,
						getPadBytes_MultiTransferETH(
							cdata.PlatformAddress,
							cdata.PlatformWEI(),
							cdata.BenefitAddress,
							cdata.BenefitWEI(),
						),
						tx_pay_wei,
						LIMIT_TAG_M_TRANS,
					)
					if err != nil {
						chttp.Fail(w, ack.NFT_InvalidPayPrice, err)
						return false
					}
					estimate_wei := jmath.ADD(tx_pay_wei, tx_gas_wei)
					if jmath.CMP(on_pay_wei, estimate_wei) < 0 {
						chttp.Fail(w, ack.NFT_NeedETHPrice)
						return false
					}

					receipt_code := nwtypes.GetReceiptCode(nwtypes.RC_MINT_COIN, cdata.TokenId)

					user_try := nwtypes.NftUserTry{
						ReceiptCode: receipt_code,

						DATA_SEQ: nwtypes.MAKE_DATA_SEQ(
							nwtypes.MULTI_TRANSFER,
							nwtypes.DataMultiTransferTry{
								ReceiptCode:  receipt_code,
								PayTokenInfo: token_info,
								PayFrom:      memberToWalletInfo(payer),
								PriceToInfos: []nwtypes.PriceToInfo{
									{
										InfoKind: nwtypes.InfoKindPlatform,
										To: nwtypes.WalletInfo{
											Address: cdata.PlatformAddress,
										},
										Price: cdata.PlatformPrice,
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

							nwtypes.NFT_MINT,
							nwtypes.DataMintTry{
								ReceiptCode: receipt_code,

								Payer: memberToWalletInfo(payer),
								Owner: memberToWalletInfo(owner),
								Kind:  cdata.MintKind(),

								NftInfo:      nft_config.NftInfo(cdata.TokenId),
								PayTokenInfo: &token_info,
								FeeInfo: &nwtypes.FeeInfo{
									PlatformAddress: cdata.PlatformAddress,
									PlatformPrice:   cdata.PlatformPrice,
									BenefitAddress:  cdata.BenefitAddress,
									BenefitPrice:    cdata.BenefitPrice,
								},
							},
						),

						TimeTryAt: unix.Now(),
					}

					user_try.InsertDB(db)

					chttp.OK(w, AckNftMint{
						ReceiptCode: receipt_code,
						TokenId:     cdata.TokenId,
					})

					return true

				})
				if !is_tokenId_rev {
					chttp.Fail(w, ack.NFTExistTokenID)

				}

			})

		},
	)
}

func _help_mint_ETH() {
	Doc().Message(`
	<cc_purple>( ETH로 민팅할 경우 )</cc_purple>

		<요청>
		post `+doc_host_url+`/v1/nft/mint/eth

		{
			"payer_address" : "0x05b93b0feeb9f60a599ba4b4c76262c22e837579", // 요청자 회원 주소
			"owner_address" : "0x05b93b0feeb9f60a599ba4b4c76262c22e837579", // NFT 소유자 회원 주소
			"token_id" : "9999990000002", //토큰ID
			"platform_address" : "0x8CE5bb2013887eD586e6a87211aa126453368b7A", //플랫폼 주소
			"platform_price" : "0.1",	//플랫폼 수수료 (0.1 ETH)
			"benefit_address" : "0x8CE5bb2013887eD586e6a87211aa126453368b7A",	//베네핏 주소
			"benefit_price" : "0.1"		//베네핏 수수료 (0.1 ETH)
		}

		<성공 응답>
		{
			"success" : true,
			"data" : {
				"receipt_code" : "nft_mint_coin_[9999990000002]_ca26c518a6d12d4bf5e266d4d5fa2789193903a0afedb1197e0f97cd171443",
				"token_id" : "9999990000002"
			}
		}

		< 스케줄러 처리 로직 >
		1. [민팅할 회원 주소]가 플랫폼/베네핏 에게 요청한 수량의 ETH 전송 (가스비 + 수수료1 + 수수료2 부담 : 민팅할 회원 주소)
		2. 1번성공시 [마스터 지갑주소]가 [민팅할 회원 주소]에게 NFT(9999990000002) 민팅. (가스비 부담 : 마스터)


		< 모든 민팅과정이 끝나면 서비스 서버로 민팅 결과를 콜백 >
		<cc_purple>post {{서비스서버}}`+URL_NFTS_MINT_CALLBACK+`</cc_purple>

		{
			"data": {
				//민팅 정보
				"mint_info": {

					//민팅결과 HASH 
					"hash": "0x6e0d7accb869a1e1716f875c765a0cb5742014055ff6d29b62ff21542dfe21c5",
					"nft_info": {
						"contract": "0x1faa080c0e0c7b94d8571d236b182f15e0c1742a",
						"symbol": "WNZ",
						"token_id": "9999990000002"
					},

					//민팅 가스 비용 (마스터 부담)					
					"tx_gas_price": "0.061961012996683617",

					//민팅 요청 주소(회원)
					"payer": {
						"address": "0x05b93b0feeb9f60a599ba4b4c76262c22e837579",
						"is_master": false,
						"name": "MmaTestnetUser01",
						"uid": 1001
					},

					//민팅 NFT소유자 정보(회원) -- NFT소유자
					"owner": {
						"address": "0x05b93b0feeb9f60a599ba4b4c76262c22e837579",
						"is_master": false,
						"name": "MmaTestnetUser01",
						"uid": 1001
					}
					
					//민팅 종류 ( self , gift , free )
					"kind" : <cc_red>"self"</cc_red> // payer.address == owner.address ? "self" : "gift"
				},

				<cc_blue>//수수료 지불 정보
				"pay_info": {
					"benefit_address": "0x8ce5bb2013887ed586e6a87211aa126453368b7a",
					"benefit_fee_price": "0.1",
					"hash": "0x97184ac92bab67aebe04550b1a7c23e2a36a0f68b3146ad11fb12db2c4ce96fd",
					"pay_from": {	// ETH민팅 이므로 수수료는 회원 지갑주소가 지불
						"address": "0x05b93b0feeb9f60a599ba4b4c76262c22e837579",
						"is_master": false,
						"name": "MmaTestnetUser01",
						"uid": 1001
					},
					"pay_token_info": {
						"contract": "eth",
						"decimal": "18",
						"is_coin": true,
						"mainnet": false,
						"symbol": "ETH"
					},
					"platform_address": "0x8ce5bb2013887ed586e6a87211aa126453368b7a",
					"platform_fee_price": "0.1",
					"tx_gas_price": "0.03300828127931136"	//가스비용도 회원지갑 주소가 지불
				}</cc_blue>
			},
			"insert_at": 1683103273,
			"is_send": false,

			//영수증 코드
			"receipt_code": "nft_mint_coin_[9999990000002]_ca26c518a6d12d4bf5e266d4d5fa2789193903a0afedb1197e0f97cd171443",

			//민팅
			"result_type": "mint",
			"send_at": 0
		}
		
	`, doc.Blue)
}
