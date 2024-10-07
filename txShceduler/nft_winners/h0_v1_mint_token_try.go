package nft_winners

import (
	"net/http"
	"strings"
	"txscheduler/brix/tools/cloud/ebcm"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/jmath"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/jnet/doc"
	"txscheduler/brix/tools/unix"
	"txscheduler/nft_winners/nwtypes"
	"txscheduler/txm/ack"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func init() {
	hNftMintTOKEN()
}

type ReqNftMintTOKEN struct {
	PayerAddress string `json:"payer_address"` // WNZ 로 민팅할 지갑 주소 (회원)
	OwnerAddress string `json:"owner_address"`
	isSelf       bool

	TokenId string `json:"token_id"`
}

func (ReqNftMintTOKEN) TagString() []string {
	return []string{
		//"pay_symbol", "구매 토큰 심볼 : ETH / WNZ(TEST) , 구매 토큰에 따라 처리방식이 분기.",

		"payer_address", "발급 요청자 지갑 주소 (회원 지갑 주소 , TOKEN 구매자)",
		"owner_address", "최초 소유자 지갑 주소 (회원 지갑 주소 , payer_address == owner_address 자체 발급)",

		"token_id", "토큰ID",
	}
}

func (my *ReqNftMintTOKEN) Valid() chttp.CError {
	if !ebcm.IsAddressP(&my.PayerAddress) {
		return ack.InvalidAddress
	}
	if !ebcm.IsAddressP(&my.OwnerAddress) {
		return ack.InvalidAddress
	}
	if my.PayerAddress == my.OwnerAddress {
		my.isSelf = true
	}

	my.TokenId = strings.TrimSpace(my.TokenId)
	if !jmath.IsNum(my.TokenId) {
		return ack.NFT_TokenId_Format
	} else if jmath.CMP(my.TokenId, 0) < 0 {
		return ack.NFT_TokenId_Format
	}

	return nil
}
func (my ReqNftMintTOKEN) MintKind() nwtypes.MintKind {
	if my.isSelf {
		return nwtypes.MintKindSelf
	}
	return nwtypes.MintKindGift
}

func hNftMintTOKEN() {
	method := chttp.POST
	url := model.V1 + "/nft/mint/token"
	Doc().Comment("[WINNERZ] TOKEN로 NFT 민팅 요청 ( 선물하기 )").
		Method(method).URL(url).
		JParam(ReqNftMintTOKEN{}, ReqNftMintTOKEN{}.TagString()...).
		JAckOK(AckNftMint{}, AckNftMint{}.TagString()...).
		//JAckError(ack.NotFoundSymbol, "pay_symbol 값이 유효하지 않음").
		JAckError(ack.InvalidAddress, "payer_address, owner_address 형식 오류").
		JAckError(ack.NFT_NotfoundPayer, "payer_address가 회원 주소가 아닐경우").
		JAckError(ack.NFT_NotfoundOwner, "owner_address가 회원 주소가 아닐경우").
		JAckError(ack.NFT_TokenId_Format, "토큰ID 형식 오류.(숫자가 아니거나 0미만일 경우)").
		JAckError(ack.NFTExistTokenID, "요청한 토큰ID가 이미 존재하거나 민팅 진행 중일 경우").
		Apply()

	_help_mint_token()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := chttp.ParseRequestJson[ReqNftMintTOKEN](req)
			if ack_err := cdata.Valid(); ack_err != nil {
				chttp.Fail(w, ack_err)
				return
			}

			token_list := inf.TokenList()
			token_info := token_list.FirstERC20() //토큰

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
					receipt_code := nwtypes.GetReceiptCode(nwtypes.RC_MINT_TOKEN, cdata.TokenId)

					master_try := nwtypes.NftMasterTry{
						ReceiptCode: receipt_code,

						DATA_SEQ: nwtypes.MAKE_DATA_SEQ(
							nwtypes.NFT_MINT,
							nwtypes.DataMintTry{
								ReceiptCode: receipt_code,

								Payer: memberToWalletInfo(payer),
								Owner: memberToWalletInfo(owner),
								Kind:  cdata.MintKind(),

								NftInfo:      nft_config.NftInfo(cdata.TokenId),
								PayTokenInfo: &token_info,
								FeeInfo:      nil,
								// FeeInfo: nwtypes.FeeInfo{
								// 	PlatformAddress: cdata.PlatformAddress,
								// 	PlatformPrice:   cdata.PlatformPrice,
								// 	BenefitAddress:  cdata.BenefitAddress,
								// 	BenefitPrice:    cdata.BenefitPrice,
								// },
							},

							// nwtypes.MULTI_TRANSFER,
							// nwtypes.DataMultiTransferTry{
							// 	ReceiptCode:  receipt_code,
							// 	PayTokenInfo: token_info,
							// 	PayFrom: nwtypes.WalletInfo{
							// 		IsMaster: true,
							// 		Address:  inf.Master().Address,
							// 	},
							// 	PriceToInfos: []nwtypes.PriceToInfo{
							// 		{
							// 			InfoKind: nwtypes.InfoKindPlatform,
							// 			To: nwtypes.WalletInfo{
							// 				Address: cdata.PlatformAddress,
							// 			},
							// 			Price: cdata.PlatformPrice,
							// 		},
							// 		{
							// 			InfoKind: nwtypes.InfoKindBenefit,
							// 			To: nwtypes.WalletInfo{
							// 				Address: cdata.BenefitAddress,
							// 			},
							// 			Price: cdata.BenefitPrice,
							// 		},
							// 	},
							// },
						),

						TimeTryAt: unix.Now(),
					}

					master_try.InsertDB(db)

					chttp.OK(w, AckNftMint{
						ReceiptCode: receipt_code,
						TokenId:     cdata.TokenId,
					})
					return true

				})
				if !is_tokenId_rev {
					chttp.Fail(w, ack.NFTExistTokenID)
					return
				}

			})

		},
	)
}

func _help_mint_token() {
	Doc().Message(`
	<cc_purple>( WINNERZ 토큰으로 민팅할 경우 )</cc_purple>

		<요청>
		post `+doc_host_url+`/v1/nft/mint/token

		{
			"payer_address" : "0x05b93b0feeb9f60a599ba4b4c76262c22e837579", // 요청자 회원 주소
			"owner_address" : "0x05b93b0feeb9f60a599ba4b4c76262c22e837579", // NFT 소유자 회원 주소
			"token_id" : "9999990000007", //토큰ID
		}

		<성공 응답>
		{
			"success" : true,
			"data" : {
				"receipt_code" : "nft_mint_token_[9999990000007]_3c4fbf830f736899f79659127a11c2c0cf1d05281bccdc719fe48f80698979f", //영수증 코드 발급
				"token_id" : "9999990000007"
			}
		}

		< 스케줄러 처리 로직 >
		1. [마스터 지갑주소]로  [민팅할 회원 주소]에게 NFT(9999990000007) 민팅. (가스비 부담 : 마스터)

		< 모든 민팅과정이 끝나면 서비스 서버로 민팅 결과를 콜백 >
		<cc_purple>post {{서비스서버}}`+URL_NFTS_MINT_CALLBACK+`</cc_purple>


		{
			"data": {
				//민팅 정보
				"mint_info": {

					//민팅결과 HASH 
					"hash": "0x7a9ecb9cfb2c58fba1d7f999a692cef05604da5b3c39788924a2b5e5f67b00ec",
					"nft_info": {
						"contract": "0x1faa080c0e0c7b94d8571d236b182f15e0c1742a",
						"symbol": "WNZ",
						"token_id": "9999990000007"
					},

					//민팅 가스 비용 (마스터 부담)
					"tx_gas_price": "0.000000000075130779",

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

				<cc_blue>//수수료 지불 정보 (없음)
				"pay_info": null</cc_blue>
			},
			"insert_at": 1684138969,
			"is_send": false,

			//영수증 코드
			"receipt_code": "nft_mint_token_[9999990000007]_1f85e62168e09177cf3e02a7824f947a4234b3efc4b00eb9fd5423b8a974b63",
			
			//민팅
			"result_type": "mint",
			"send_at": 0,
			"send_ymd": 0
		}
		
	`, doc.Blue)
}
