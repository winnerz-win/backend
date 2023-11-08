package api

import (
	"net/http"
	"txscheduler/brix/tools/cloudx/ethwallet/ecsx"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jmath"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/mms"
	"txscheduler/txm/ack"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func init() {
	hWithdrawTry()
	hWithdrawWaitings()
	hWithdrawResult()
}

type WithdrawRequest struct {
	Name        string `json:"name"`
	FromAddress string `json:"from_address"`
	Symbol      string `json:"symbol"`
	Price       string `json:"price"`
	ToAddress   string `json:"to_address"`
}

func (WithdrawRequest) TagString() []string {
	return []string{
		"name", "회원 ID",
		"from_address", "가입시 발급한 입금주소(회원 주소)",
		"symbol", "코인 심볼",
		"price", "출금 가격 (20.12)",
		"to_address", "코인을 받을 주소 (받는이)",
	}
}

func hWithdrawTry() {
	type RESULT struct {
		ReceiptCode string `json:"receipt_code"`
	}

	method := chttp.POST
	url := model.V1 + "/withdraw/try"
	Doc().Comment("[ 코인 출금 요청 ] 코인 출금 요청").
		Method(method).URL(url).
		JParam(WithdrawRequest{}, WithdrawRequest{}.TagString()...).
		JResultOK(chttp.AckFormat{}).
		Etc("", `_
			<cc_blue>응답결과</cc_blue>
			{
				"success" : true,
				"data" :{
					"receipt_code" : "receipt_79f0b2da02de7f366a1125813aa3e32a146e8995" //출금 확인용 영수증(고유키값)
 				}
			}
		`).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := WithdrawRequest{}
			chttp.BindingJSON(req, &cdata)

			model.Trim(&cdata.Name)
			model.Trim(&cdata.Symbol)
			model.Trim(&cdata.Price)
			cdata.FromAddress = dbg.TrimToLower(cdata.FromAddress)
			cdata.ToAddress = dbg.TrimToLower(cdata.ToAddress)

			if ecsx.IsAddress(cdata.ToAddress) == false {
				chttp.Fail(w, ack.InvalidAddress)
				return
			}

			if jmath.IsUnderZero(cdata.Price) {
				chttp.Fail(w, ack.UnderZERO)
				return
			}

			if inf.ValidSymbol(cdata.Symbol) == false {
				chttp.Fail(w, ack.NotFoundSymbol)
				return
			}
			model.DB(func(db mongo.DATABASE) {
				member := model.LoadMemberAddress(db, cdata.FromAddress)
				if member.Valid() == false {
					chttp.Fail(w, ack.NotFoundAddress)
					return
				}
				if member.Name != cdata.Name {
					chttp.Fail(w, ack.NotFoundName)
					return
				}

				nowAt := mms.Now()
				token := inf.Config().Tokens.GetSymbol(cdata.Symbol)
				data := model.TxETHWithdraw{
					UID:       member.UID,
					ToAddress: cdata.ToAddress,
					ToPrice:   cdata.Price,

					Symbol:    cdata.Symbol,
					Decimal:   token.Decimal,
					Timestamp: nowAt,
					YMD:       nowAt.YMD(),

					State: model.TxStateNone,
				}
				receiptCode := data.InsertDB(db)

				chttp.OK(w, RESULT{
					ReceiptCode: receiptCode,
				})

			})

		},
	)
}

func hWithdrawWaitings() {
	type RESULT struct {
		Waiting int `json:"waitings"`
	}
	method := chttp.GET
	url := model.V1 + "/withdraw/waitings"
	Doc().Comment("[ 코인 출금 ] 현재 코인 출금을 위한 대기열 갯수").
		Method(method).URL(url).
		JResultOK(chttp.AckFormat{}).
		Etc(".", `_
			<cc_blue>응답결과</cc_blue>
			{
				"success": true,
				"data": {
					"waitings": 0	//출금 대기중인 요청의 갯수를 반환 합니다.
				}
			}
		`).
		Apply()
	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			model.DB(func(db mongo.DATABASE) {

				waitings, _ := db.C(inf.TXETHWithdraw).Count()
				chttp.OK(w, RESULT{
					Waiting: waitings,
				})

			})
		},
	)
}

func hWithdrawResult() {
	type RESULT struct {
		ReceiptCode   string        `json:"receipt_code"`
		Name          string        `json:"name"`
		UID           int64         `json:"uid"`
		FromAddress   string        `json:"from_address"`
		ToAddress     string        `json:"to_address"`
		Hash          string        `json:"hash"`
		Symbol        string        `json:"symbol"`
		Price         string        `json:"price"`
		WithdrawState model.TxState `json:"withdraw_state"`
		Timestamp     mms.MMS       `json:"timestamp"`
		IsSend        bool          `json:"is_send"`
		SendAt        mms.MMS       `json:"send_at"`
	}

	method := chttp.GET
	url := model.V1 + "/withdraw/result/:args"
	Doc().Comment("[ 코인 출금 결과 ] 코인 출금 결과").
		Method(method).URLS(url, ":args", "출금영수증(receipt_code)").
		Etc(".", `_
			<cc_blue>응답결과</cc_blue>
			{
				"success": true,
				"data": {
					"address": "0xef14301c9530d52f20f7acad0049780a3927a1fe",	//회원 지갑주소
					"hash": "0xb495dd89e0d72ae00d459dce8a42922aed40b942fcc4cee8f39153efcb8e5646",	//트랜잭션 해시
					"is_send": true,	// 어플리케이션 서버에 응답 완료여부
					"name": "test2",	// 회원ID
					"receipt_code": "receipt_b4f8f9dcdca00c355d64e53ee69a087bdd991fb2",	//영수증 코드
					"send_at": 1617519519833,	// 어플리케이션 서버에 응답 시간
					"withdraw_state": 200,		// 전송 결과 ( 200 : 성공 , 104 : 실패 , 1 : 전송 대기중)
					"symbol": "ETH",	// 코인 심볼
					"timestamp": 1617519519113,	// 코인 전송 시간
					"to_address": "0x8ce5bb2013887ed586e6a87211aa126453368b7a",		// 받는이 주소
					"to_price": "3.2",	// 코인 전송 수량
					"uid": 1002			// 회원 UID
				}
			}
		`).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {

			receiptCode := ps.ByName("args")

			if receiptCode == "" {
				chttp.Fail(w, ack.InvalidReceiptCode)
				return
			}

			model.DB(func(db mongo.DATABASE) {
				witem := model.TxETHWithdraw{}.GetData(db, receiptCode)
				if witem.Valid() {
					member := model.LoadMember(db, witem.UID)
					log := witem.MakeLogWithdraw(member, model.TxStatePending)
					chttp.OK(w, RESULT{
						ReceiptCode:   log.ReceiptCode,
						Name:          log.Name,
						UID:           log.UID,
						FromAddress:   log.Address,
						ToAddress:     log.ToAddress,
						Hash:          log.Hash,
						Symbol:        log.Symbol,
						Price:         log.ToPrice,
						WithdrawState: log.State,
						Timestamp:     log.Timestamp,
						IsSend:        log.IsSend,
						SendAt:        log.SendAt,
					})
					return
				}

				log := model.LogWithdraw{}.GetData(db, receiptCode)
				if log.Valid() {
					chttp.OK(w, RESULT{
						ReceiptCode:   log.ReceiptCode,
						Name:          log.Name,
						UID:           log.UID,
						FromAddress:   log.Address,
						ToAddress:     log.ToAddress,
						Hash:          log.Hash,
						Symbol:        log.Symbol,
						Price:         log.ToPrice,
						WithdrawState: log.State,
						Timestamp:     log.Timestamp,
						IsSend:        log.IsSend,
						SendAt:        log.SendAt,
					})
					return
				}

				chttp.Fail(w, ack.InvalidReceiptCode)
			})
		},
	)
}
