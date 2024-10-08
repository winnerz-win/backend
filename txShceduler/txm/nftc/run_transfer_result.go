package nftc

import (
	"time"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

type NftTransferResultCallback func(db mongo.DATABASE, list model.NftTransferEndList)

/*
( 7. NFT 소유권 이전 결과 콜백  )
http://service.server.com:8080/api/v1/nft/transfer/callback

method : post
content-type : application/json

전송 데이터

	{
		"uid": string			// 유저 고유ID (전송자)
		"address": string		// 유저 발급 주소 (전송자)
		"name": string			// 유저 이름 (전송자)
		"token_id": string		// NFT 토큰 ID
		"to_address": string	// 수신자 주소
		"is_member_to": bool	// 수신자가 내부 회원인지여부
		"hash": string			// NFT 트랜젝션 해시값
		"status": int			// 전송 결과 ( 성공:200 , 실패:104 )
		"fail_code": int		// 실패 코드 ( 실패일경우 )
		"fail_message": string	// 실패 사유 ( 실패일경우 )
		"timestamp": long		// 결과 시간 (mms)
		"gas_fee_eth": string	// 트랜젝션의 가스비용 (ETH값)
	}

----------------------------------------------
실패코드 - fail_code(fail_message)
0	 ("")						: 없음
1002 (transaction_result_fail)	: 전송 트랜젝션 결과 실패
*/
func runTransferResult(transferCallback NftTransferResultCallback) {
	dbg.PrintForce("runTransferResult ------ START")
	go runTransferPending()
	for {
		model.DB(func(db mongo.DATABASE) {
			list := model.NftTransferEndList{}
			selector := mongo.Bson{"is_send": false}
			db.C(inf.NFTTransferEnd).Find(selector).All(&list)

			if len(list) > 0 {
				transferCallback(db, list)
			}

		})
		time.Sleep(time.Second)
	} //for
}

func runTransferPending() {
	defer dbg.PrintForce("nftc.runTransferPending ----------  END")
	dbg.PrintForce("nftc.runTransferPending ----------  START")

	for {
		model.DB(func(db mongo.DATABASE) {
			pendings := model.NftTransferTryList{}
			db.C(inf.NFTTransferTry).Find(nil).All(&pendings)
			for _, try := range pendings {
				r, _, _ := Sender().TransactionByHash(try.Hash)
				if !r.IsReceiptedByHash {
					continue
				}
				model.LockMember(db, try.Address, func(member model.Member) {
					member.UpdateCoinDB_Legacy(db)
				})

				endData := model.NftTransferEnd{
					NftTransferTry: try,
					GasFeeETH:      r.GetTransactionFee(),
				}
				if r.IsError {
					endData.SetFail(model.FailNftTxResult)
				} else {
					endData.Status = 200
					buyEnd := model.NftBuyEnd{}.GetTokenID(db, try.TokenId)
					if buyEnd.Valid() {
						buyEnd.LastOwner = try.ToAddress
						if try.IsMemberTo {
							toMember := model.LoadMemberAddress(db, try.ToAddress)
							if toMember.Valid() {
								buyEnd.ChangeInnerOwner(db, toMember.User)
							}

						} else {
							buyEnd.SetLastOwner(db, try.ToAddress)
						}
					}
				}
				if endData.InsertEndDB(db) == nil {
					try.RemoveTryDB(db)
				}

			} //for
		})
	} //for
}
