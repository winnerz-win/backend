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
	"txscheduler/txm/model"
)

func init() {
	hNftMintFree()
}

type ReqNftMintFree struct {
	OwnerAddress string `json:"owner_address"`
	TokenId      string `json:"token_id"`
}

func (ReqNftMintFree) TagString() []string {
	return []string{
		"owner_address", "최초 소유자 지갑 주소 (회원 지갑 주소)",
		"token_id", "토큰ID",
	}
}

func (my *ReqNftMintFree) Valid() chttp.CError {

	if !ebcm.IsAddressP(&my.OwnerAddress) {
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

func hNftMintFree() {
	method := chttp.POST
	url := model.V1 + "/nft/mint/free"
	Doc().Comment("[WINNERZ] 운영자 무료 발급(이벤트용)").
		Method(method).URL(url).
		JParam(ReqNftMintFree{}, ReqNftMintFree{}.TagString()...).
		JAckOK(AckNftMint{}, AckNftMint{}.TagString()...).
		JAckError(ack.InvalidAddress, "owner_address 형식 오류").
		JAckError(ack.NFT_NotfoundOwner, "owner_address가 회원 주소가 아닐경우").
		JAckError(ack.NFT_TokenId_Format, "토큰ID 형식 오류.(숫자가 아니거나 0미만일 경우)").
		JAckError(ack.NFTExistTokenID, "요청한 토큰ID가 이미 존재하거나 민팅 진행 중일 경우").
		Apply()

	_help_mint_free()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := chttp.ParseRequestJson[ReqNftMintFree](req)
			if ack_err := cdata.Valid(); ack_err != nil {
				chttp.Fail(w, ack_err)
				return
			}

			model.DB(func(db mongo.DATABASE) {
				owner := model.LoadMemberAddress(db, cdata.OwnerAddress)
				if !owner.Valid() {
					chttp.Fail(w, ack.NFT_NotfoundOwner)
					return
				}

				is_tokenId_rev := nwtypes.NftTokenIDS{}.Reserve(db, cdata.TokenId, func() bool {

					receipt_code := nwtypes.GetReceiptCode(nwtypes.RC_MINT_FREE, cdata.TokenId)

					master_try := nwtypes.NftMasterTry{
						ReceiptCode: receipt_code,

						DATA_SEQ: nwtypes.MAKE_DATA_SEQ(
							nwtypes.NFT_MINT,
							nwtypes.DataMintTry{
								ReceiptCode: receipt_code,

								Payer: nwtypes.MasterWalletInfo(),
								Owner: memberToWalletInfo(owner),
								Kind:  nwtypes.MintKindFree,

								NftInfo:      nft_config.NftInfo(cdata.TokenId),
								PayTokenInfo: nil, //FREE
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

func _help_mint_free() {
	Doc().Message(`
	<cc_purple>( 운영자 무료(FREE) 민팅할 경우 )</cc_purple>

		<요청>
		post `+doc_host_url+`/v1/nft/mint/free

		{
			"owner_address" : "0xe129243a027b25d813aced72fe34b22f5fc4bb20", // NFT 소유자 회원 주소
			"token_id" : "9999990000009", //토큰ID
		}

		<성공 응답>
		{
			"success" : true,
			"data" : {
				"receipt_code" : "nft_mint_free_[9999990000009]_7176c8db2e09e797bf572518afee533d0ccead2b8688c4136c7b643e6a909b", //영수증 코드 발급
				"token_id" : "9999990000009"
			}
		}

		< 스케줄러 처리 로직 >
		1. [마스터 지갑주소]로  [민팅할 회원 주소]에게 NFT(9999990000009) 민팅. (가스비 부담 : 마스터)

		< 모든 민팅과정이 끝나면 서비스 서버로 민팅 결과를 콜백 >
		<cc_purple>post {{서비스서버}}`+URL_NFTS_MINT_CALLBACK+`</cc_purple>


		{
			"data": {
				//민팅 정보
				"mint_info": {

					//민팅결과 HASH 
					"hash": "0x47a39609df76046bfa1e6ef79a6ead098da4db3370831c5c85c01d2667378020",
					"nft_info": {
						"contract": "0x1faa080c0e0c7b94d8571d236b182f15e0c1742a",
						"symbol": "WNZ",
						"token_id": "9999990000009"
					},

					//민팅 가스 비용 (마스터 부담)
					"tx_gas_price": "0.000000558284330523",

					//민팅 요청 주소(마스터)
					"payer": {
						"address": "",
						<cc_blue>"is_master": true,</cc_blue>
						"uid": 0
					},

					//민팅 NFT소유자 정보(회원) -- NFT소유자
					"owner": {
						"address": "0xe129243a027b25d813aced72fe34b22f5fc4bb20",
						"is_master": false,
						"name": "mma_testnet_user_02",
						"uid": 1002
					}

					//민팅 종류 ( self , gift , free )
					"kind" : <cc_red>"free"</cc_red> // payer.is_master == true
				},

				<cc_blue>//수수료 지불 정보 (없음)
				"pay_info": null</cc_blue>
			},
			"insert_at": 1685434766,
			"is_send": false,

			//영수증 코드
			"receipt_code": "nft_mint_free_[9999990000009]_7176c8db2e09e797bf572518afee533d0ccead2b8688c4136c7b643e6a909b",
			
			//민팅
			"result_type": "mint",
			"send_at": 0,
			"send_ymd": 0
		}
		
	`, doc.Blue)
}
