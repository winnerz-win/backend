package admin

import (
	"net/http"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/jmath"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/jnet/doc"
	"txscheduler/txm/ack"
	"txscheduler/txm/cloud"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func init() {
	hDepositToMasterOne()
	hDepositToMasterAll()
}

func hDepositToMasterOne() {
	type CDATA struct {
		UID    int64  `json:"uid"`
		Symbol string `json:"symbol"`
	}

	method := chttp.POST
	url := model.V2 + "/deposit_to_master/one"
	Doc().Comment("[ 마스터송금 ] 입금계좌에서 마스터로 재화 전송 (ROOT 권한)").
		Method(method).URL(url).
		JParam(CDATA{},
			"symbol", "ETH , GDG",
			"min_price", "마스터로 보내기위한 최소 잔액량 (0이면 잔액 체크 없이 모두 보냄)",
		).
		JResultOK(chttp.AckFormat{}).
		ResultERRR(ack.InvalidRootAdmin).
		ResultERRR(ack.NotFoundSymbol).
		ResultERRR(ack.NotFoundName).
		ResultERRR(ack.UnderZERO).
		ResultERRR(ack.AlreadyProcessJob).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			admin := model.GetTokenAdmin(req)
			if admin.IsRoot == false {
				chttp.Fail(w, ack.InvalidRootAdmin)
				return
			}

			cdata := CDATA{}
			chttp.BindingJSON(req, &cdata)

			tokeninfo := inf.TokenList().GetSymbol(cdata.Symbol)
			if tokeninfo.Valid() == false {
				chttp.Fail(w, ack.NotFoundSymbol)
				return
			}

			list := cloud.ETHDepositList{}
			model.DB(func(db mongo.DATABASE) {

				member := model.LoadMember(db, cdata.UID)
				if member.Valid() == false {
					chttp.Fail(w, ack.NotFoundName)
					return
				}

				val := member.Coin.Price(cdata.Symbol)
				if jmath.CMP(val, 0) <= 0 {
					chttp.Fail(w, ack.UnderZERO)
					return
				}

				cnt, _ := db.C(inf.TXETHDepositLog).Find(member.Selector()).Count()
				if cnt > 0 {
					chttp.Fail(w, ack.AlreadyProcessJob)
					return
				}

				item := cloud.ETHDeposit{
					UID:      member.UID,
					Address:  member.Address,
					Symbol:   tokeninfo.Symbol,
					Contract: tokeninfo.Contract,
					Decimal:  tokeninfo.Decimal,
					IsForce:  true,
				}
				list = append(list, item)

			})

			if len(list) > 0 {
				cloud.ETHDepositChan <- list
			}

			chttp.OK(w, struct{}{})

		},
	)
}

func hDepositToMasterAll() {
	type CDATA struct {
		Symbol   string `json:"symbol"`
		MinPrice string `json:"min_price"`
	}
	type RESULT struct {
		TryCount int `json:"try_count"`
	}
	method := chttp.POST
	url := model.V2 + "/deposit_to_master/all"
	Doc().Comment("[ 마스터송금 ] 입금계좌에서 마스터로 전체 재화 전송 (ROOT 권한)").
		Method(method).URL(url).
		JParam(CDATA{},
			"symbol", "ETH , GDG",
			"min_price", "마스터로 보내기위한 최소 잔액량 (0이면 잔액 체크 없이 모두 보냄)",
		).
		JResultOK(chttp.AckFormat{}).
		ETC(doc.EV(RESULT{},
			"try_count", "적용된 계좌 갯수",
		)).
		ResultERRR(ack.InvalidRootAdmin).
		ResultERRR(ack.NotFoundSymbol).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			admin := model.GetTokenAdmin(req)
			if admin.IsRoot == false {
				chttp.Fail(w, ack.InvalidRootAdmin)
				return
			}

			cdata := CDATA{}
			chttp.BindingJSON(req, &cdata)

			if jmath.CMP(cdata.MinPrice, 0) < 0 {
				cdata.MinPrice = model.ZERO
			}

			tokeninfo := inf.TokenList().GetSymbol(cdata.Symbol)
			if tokeninfo.Valid() == false {
				chttp.Fail(w, ack.NotFoundSymbol)
				return
			}

			list := cloud.ETHDepositList{}
			model.DB(func(db mongo.DATABASE) {
				member := model.Member{}
				iter := db.C(inf.COLMember).Find(nil).Iter()
				for iter.Next(&member) {
					val := member.Coin.Price(cdata.Symbol)
					if jmath.CMP(val, 0) <= 0 {
						continue
					}
					cnt, _ := db.C(inf.TXETHDepositLog).Find(member.Selector()).Count()
					if cnt > 0 {
						continue
					}

					item := cloud.ETHDeposit{
						UID:      member.UID,
						Address:  member.Address,
						Symbol:   tokeninfo.Symbol,
						Contract: tokeninfo.Contract,
						Decimal:  tokeninfo.Decimal,
						IsForce:  true,
					}
					list = append(list, item)
				} //for

			})

			if len(list) > 0 {
				cloud.ETHDepositChan <- list
			}

			chttp.OK(w, RESULT{
				TryCount: len(list),
			})
		},
	)
}
