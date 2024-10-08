package api

import (
	"jtools/cloud/ebcm"
	"jtools/jmath"
	"jtools/mms"
	"net/http"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/txm/ack"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
	"txscheduler/txm/scv"
)

func init() {
	hMasterOutTry()
	hMasterOutResult()
}

func hMasterOutTry() {
	type CDATA struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
		To     string `json:"to"`
	}
	type RESULT struct {
		ReceiptCode string `json:"receipt_code"`
	}

	method := chttp.POST
	url := model.V1 + "/master/out/try"
	Doc().Comment("[ 마스터출금 ] 마스터 계좌에서 외부계좌로 출금신청").
		Method(method).URL(url).
		JParam(CDATA{},
			"symbol", "ETH , GDG",
			"price", "가격",
			"to", "받을 주소",
		).
		JAckOK(RESULT{},
			"receipt_code", "출금 영수증 (master_xxxxxxx)",
			"", "",
		).
		JAckError(ack.UnderZERO).
		JAckError(ack.NotFoundSymbol).
		JAckError(ack.InvalidAddress).
		JAckError(ack.DBJob).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {

			if model.MasterCallback(scv.MasterOutCallback) == false {
				chttp.Fail(w, ack.NotAllowSystem)
				return
			}

			cdata := CDATA{}
			chttp.BindingJSON(req, &cdata)

			if jmath.CMP(cdata.Price, 0) <= 0 {
				chttp.Fail(w, ack.UnderZERO)
				return
			}

			cdata.To = dbg.TrimToLower(cdata.To)
			if ebcm.IsAddress(cdata.To) == false {
				chttp.Fail(w, ack.InvalidAddress)
				return
			}

			tokeninfo := inf.TokenList().GetSymbol(cdata.Symbol)
			if tokeninfo.Valid() == false {
				chttp.Fail(w, ack.NotFoundSymbol)
				return
			}

			model.DB(func(db mongo.DATABASE) {

				nowAt := mms.Now()
				try := model.TxETHMasterOut{
					ReceiptCode: model.GetMasterCode(),
					Symbol:      cdata.Symbol,
					Decimal:     tokeninfo.Decimal,
					ToAddress:   cdata.To,
					ToPrice:     cdata.Price,

					State:     model.TxStateNone,
					Timestamp: nowAt,
					YMD:       nowAt.YMD(),
				}

				if try.InsertTry(db) != nil {
					chttp.Fail(w, ack.DBJob)
					return
				}

				chttp.OK(w, RESULT{
					ReceiptCode: try.ReceiptCode,
				})

			})
		},
	)

}

func hMasterOutResult() {
	type RESULT struct {
		ReceiptCode string `json:"receipt_code"`
		Hash        string `json:"hash"`
		FromAddress string `json:"from_address"`
		ToAddress   string `json:"to_address"`
		Symbol      string `json:"symbol"`
		Price       string `json:"price"`

		WithdrawResult model.TxState `json:"withdraw_result"`
		FailMessage    string        `json:"fail_message"`
		Timestamp      mms.MMS       `json:"timestamp"`
	}

	method := chttp.GET
	url := model.V1 + "/master/out/result/:args"
	Doc().Comment("[ 마스터출금 ] 영수증 결과 확인 ").
		Method(method).URLS(url, ":args", "출금영수증(receipt_code)").
		Etc(".", `_
			<cc_blue>응답결과</cc_blue>
			{
				"success": true,
				"data": {
					"receipt_code" : string,	//영수증
					"hash" : string,			//트랜젝션 해시
					"from_address" : string,	//마스터 지갑 주소
					"to_address" : string,		//출금 주소
					"symbol" : string,			//토큰 심볼
					"price" : string,			//출금 액수

					"withdraw_result" : int,	//전송 결과 ( 200 : 성공 , 104 : 실패 , 1 : 전송 대기중)
					"fail_message" : string,	//실패(104)일 경우 메시지
					"timestamp" : int64			//시간
				}
			}

			----------------------------------------------
			<cc_blue>실패사유 : fail_message </cc_blue>
			invalid_symbol 				:	등록되어있지 않는 토큰 심볼 오류
			need_price 					:	마스터지갑의 잔액부족 (토큰/이더 수량 , 가스비 등등)
			chain_error:box 			:	메인 노드 에러
			chain_error:nonce 			:	메인 노드 에러
			chain_error:tx 				:	메인 노드 에러
			chain_error:send 			:	메인 노드 에러
			chain_error:pending_time 	:	전송한 트렌젝션이 시간 지연(4시간 체크)으로 인한 fallback 처리
			chain_error:fail 			:	실패한 트렌젝션
		`).
		JAckError(ack.InvalidReceiptCode).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			if model.MasterCallback(scv.MasterOutCallback) == false {
				chttp.Fail(w, ack.NotAllowSystem)
				return
			}

			receiptCode := ps.ByName("args")
			if receiptCode == "" {
				chttp.Fail(w, ack.InvalidReceiptCode)
				return
			}

			model.DB(func(db mongo.DATABASE) {
				item := model.TxETHMasterOut{}.GetReceipt(db, receiptCode)

				if item.Valid() == false {
					chttp.Fail(w, ack.InvalidReceiptCode)
					return
				}

				chttp.OK(w, RESULT{
					ReceiptCode: item.ReceiptCode,
					Hash:        item.Hash,
					FromAddress: inf.Master().Address,
					ToAddress:   item.ToAddress,
					Symbol:      item.Symbol,
					Price:       item.ToPrice,

					WithdrawResult: item.State,
					FailMessage:    item.FailMessage,
					Timestamp:      item.Timestamp,
				})
			})
		},
	)

}
