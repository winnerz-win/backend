package nft_winners

import (
	"net/http"
	"sort"
	"txscheduler/brix/tools/cloud/ebcm"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/jmath"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/jnet/doc"
	"txscheduler/nft_winners/nwdb"
	"txscheduler/nft_winners/nwtypes"
	"txscheduler/nft_winners/rpc"
	"txscheduler/txm/ack"
	"txscheduler/txm/model"
)

func init() {
	hNftInfo()
	hNftTokenIdOwner()
	NftOwnerTokenIds()
	hNftReceiptCode()
}

type NftInfoAck struct {
	NftAddress  string `json:"nft_address"`
	NftName     string `json:"nft_name"`
	NftSymbol   string `json:"nft_symbol"`
	BaseURI     string `json:"base_uri"`
	TotalSupply string `json:"total_supply"`
}

func (NftInfoAck) TagString() []string {
	return []string{
		"nft_address", "NFT 컨트랙트 주소",
		"nft_name", "NFT 이름",
		"nft_symbol", "NFT 심볼",
		"base_uri", "NFT BASE_URI",
		"total_supply", "NFT 총 발행량",
	}
}

func hNftInfo() {
	method := chttp.GET
	url := model.V1 + "/nft/info"

	Doc().Comment("[WINNERZ_INFO] NFT 컨트랙트 정보 요청").
		Method(method).URL(url).
		JAckOK(NftInfoAck{}, NftInfoAck{}.TagString()...).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			model.DB(func(db mongo.DATABASE) {
				r := NftInfoAck{
					NftAddress: nft_config.NftContract,
					NftName:    nft_config.NftName,
					NftSymbol:  nft_config.NftSymbol,
				}

				if finder := GetSender(); finder != nil {
					baseURI, _ := rpc.ERC721.BaseURI(
						finder,
						nft_config.Reader(),
					)
					r.BaseURI = baseURI

					rpc.ERC721.TotalSupply(
						finder,
						nft_config.Reader(),
						func(_total_supply string) {
							r.TotalSupply = _total_supply
						},
					)
				}

				chttp.OK(w, r)
			})
		},
	)
}

//////////////////////////////////////////////////////////////////////////////////////
//

type NftTokenIdOwnerAck struct {
	TokenId  string `json:"token_id"`
	Address  string `json:"address"`
	IsMember bool   `json:"is_member"`
	UID      int64  `json:"uid"`
	Name     string `json:"name"`
}

func (NftTokenIdOwnerAck) TagString() []string {
	return []string{
		"token_id", "토큰ID",
		"address", "소유자 주소",
		"is_member", "소유자가 회원인지여부",
		"uid", "회원 UID",
		"name", "회원 ID",
		"", "",
	}
}

func hNftTokenIdOwner() {
	method := chttp.GET
	url := model.V1 + "/nft/find_owner/:args"

	Doc().Comment("[WINNERZ_INFO] TokenID의 소유자 검색").
		Method(method).
		URLS(url,
			":args", "토큰ID",
		).
		JAckOK(NftTokenIdOwnerAck{}, NftTokenIdOwnerAck{}.TagString()...).
		JAckError(ack.NFTNotfoundTokenID).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			TokenId := ps.ByName("args")
			model.DB(func(db mongo.DATABASE) {

				data := nwtypes.NftTokenIDS{}
				db.C(nwdb.NftTokenIDS).Find(mongo.Bson{"token_id": TokenId}).One(&data)
				if !data.Valid() {
					chttp.Fail(w, ack.NFTNotfoundTokenID)
					return
				}

				if data.IsFixed {
					member := model.LoadMemberAddress(db, data.Owner)
					r := NftTokenIdOwnerAck{
						TokenId:  TokenId,
						Address:  data.Owner,
						IsMember: member.Valid(),
						UID:      member.UID,
						Name:     member.Name,
					}
					chttp.OK(w, r)
					return
				}

				address := ""
				err := rpc.ERC721.OwnerOf(
					GetSender(), nft_config.Reader(),
					TokenId,
					func(_owner string) {
						address = _owner
					},
				)
				if err == nil {
					member := model.LoadMemberAddress(db, address)
					r := NftTokenIdOwnerAck{
						TokenId:  TokenId,
						Address:  address,
						IsMember: member.Valid(),
						UID:      member.UID,
						Name:     member.Name,
					}
					chttp.OK(w, r)
					return
				}

				chttp.Fail(w, ack.NFTNotfoundTokenID)

			})
		},
	)

}

//////////////////////////////////////////////////////////////////////////////////////
//

type NftOwnerTokenIdsAck struct {
	Address  string `json:"address"`
	IsMember bool   `json:"is_member"`
	UID      int64  `json:"uid"`
	Name     string `json:"name"`

	TokenIds []string `json:"token_ids"`
}

func (NftOwnerTokenIdsAck) TagString() []string {
	return []string{
		"address", "소유자 주소",
		"is_member", "소유자가 회원인지여부",
		"uid", "회원 UID",
		"name", "회원 ID",
		"token_ids", "토큰ID 리스트",
		"", "",
	}
}

func NftOwnerTokenIds() {
	method := chttp.GET
	url := model.V1 + "/nft/find_tokenids/:owner"

	Doc().Comment("[WINNERZ_INFO] (회원)주소가 보유한 TokenID 리스트 요청").
		Method(method).
		URLS(url,
			":owner", "(회원)주소",
		).
		JAckOK(NftOwnerTokenIdsAck{}, NftOwnerTokenIdsAck{}.TagString()...).
		JAckError(ack.InvalidAddress).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			owner := ps.ByName("owner")
			if !ebcm.IsAddressP(&owner) {
				chttp.Fail(w, ack.InvalidAddress)
				return
			}
			model.DB(func(db mongo.DATABASE) {
				r := NftOwnerTokenIdsAck{
					Address:  owner,
					TokenIds: []string{},
				}
				member := model.LoadMemberAddress(db, owner)
				if member.Valid() {
					r.IsMember = true
					r.UID = member.UID
					r.Name = member.Name

					//rpc.ERC721.BalanceOf()

					mongo.IterForeach(
						db.C(nwdb.NftTokenIDS).Find(mongo.Bson{"owner": owner}).Iter(),
						func(cnt int, item nwtypes.NftTokenIDS) bool {
							if !item.IsFixed {
								return false
							}
							r.TokenIds = append(r.TokenIds, item.TokenId)
							return false
						},
					)

				} else {
					list, _ := rpc.ERC721.OwnerTokenAll(
						GetSender(), nft_config.Reader(),
						owner,
					)
					r.TokenIds = list
				}

				sort.Slice(r.TokenIds, func(i, j int) bool {
					return jmath.CMP(r.TokenIds[i], r.TokenIds[j]) < 0
				})

				chttp.OK(w, r)
			})
		},
	)
}

//////////////////////////////////////////////////////////////////////////////////////
//

type AckReceiptCode struct {
	ActionState  int                      `json:"action_state"`
	ActionResult *nwtypes.NftActionResult `json:"action_result,omitempty"`
}

func (AckReceiptCode) TagString() []string {
	return []string{
		"action_state", "0 : 존재하지 않음, 1: 작업 대기/진행중, 200 : 작업 완료",
		"action_result,omitempty", "action_state == 200 일경우 결과 데이터",
	}
}

func hNftReceiptCode() {
	method := chttp.GET
	url := model.V1 + "/nft/receipt_code/:code"

	Doc().Comment("[WINNTERZ_CODE_INFO] 발급한 영수증 코드 정보 보기").
		Method(method).
		URLS(url,
			":code", "발급한 영수증 코드",
		).
		JAckOK(AckReceiptCode{}, AckReceiptCode{}.TagString()...).
		JAckError(ack.NFT_ReceiptCodeFormat).
		Apply()

	_help_receiptCode()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			code := ps.ByName("code")

			receipt_code := nwtypes.RECEIPT_CODE(code)

			if !receipt_code.Valid() {
				chttp.Fail(w, ack.NFT_ReceiptCodeFormat)
				return
			}

			model.DB(func(db mongo.DATABASE) {
				selector := mongo.Bson{"receipt_code": receipt_code}

				action_result := nwtypes.NftActionResult{}
				if err := db.C(nwdb.NftActionResult).
					Find(selector).
					One(&action_result); err == nil {
					r := AckReceiptCode{
						ActionState:  200,
						ActionResult: &action_result,
					}
					chttp.OK(w, r)
					return
				}

				if cnt, _ := db.C(nwdb.NftMasterTry).Find(selector).Count(); cnt > 0 {
					chttp.OK(w, AckReceiptCode{ActionState: 1})
					return
				}
				if cnt, _ := db.C(nwdb.NftUserTry).Find(selector).Count(); cnt > 0 {
					chttp.OK(w, AckReceiptCode{ActionState: 1})
					return
				}

				if cnt, _ := db.C(nwdb.NftMasterPending).Find(selector).Count(); cnt > 0 {
					chttp.OK(w, AckReceiptCode{ActionState: 1})
					return
				}
				if cnt, _ := db.C(nwdb.NftUserPending).Find(selector).Count(); cnt > 0 {
					chttp.OK(w, AckReceiptCode{ActionState: 1})
					return
				}

				chttp.OK(w, AckReceiptCode{ActionState: 0})
				return

			})
		},
	)
}

func _help_receiptCode() {
	Doc().Message(`
	<cc_purple>( 영수증 처리 결과 보기 )</cc_purple>

		<요청>
		get `+doc_host_url+`/v1/nft/receipt_code/nft_mint_token_[9999990000008]_45de2d2096d27960f665fb733fc8ad5e89c4636aa16790ef62d9f24812882d9
		
		<성공 응답>
		{
			"success": true,
			"data": {
				"action_state": 200,
				"action_result": {
				"receipt_code": "nft_mint_token_[9999990000008]_45de2d2096d27960f665fb733fc8ad5e89c4636aa16790ef62d9f24812882d9",
				"result_type": "mint",
				"data": {
					"mint_info": {
					"hash": "0xc8c26b1f1062665c57552d71e2d36bc18bc0b0f8c888abfe63a6c3fe9da224db",
					<cc_red>"kind": "gift"</cc_red>, //민팅 선물하기 payer != owner 정보가 다름
					"nft_info": {
						"contract": "0x1faa080c0e0c7b94d8571d236b182f15e0c1742a",
						"symbol": "WNZ",
						"token_id": "9999990000008"
					},
					"owner": {
						"address": "0x05b93b0feeb9f60a599ba4b4c76262c22e837579",
						"is_master": false,
						"name": "MmaTestnetUser01",
						"uid": 1001
					},
					"payer": {
						"address": "0xe129243a027b25d813aced72fe34b22f5fc4bb20",
						"is_master": false,
						"name": "mma_testnet_user_02",
						"uid": 1002
					},
					"tx_gas_price": "0.000000558247470114"
					},
					"pay_info": null
				},
				"insert_at": 1685432715,
				"is_send": true,
				"send_at": 0,
				"send_ymd": 0
				}
			}
		}
		
	`, doc.Blue)
}
