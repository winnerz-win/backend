package admin

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
	"txscheduler/txm/scv"
)

func init() {
	// hMasterOutLog()
	// hMasterOutPending()
	// hMasterOutTry()
}

func hMasterOutLog() {

}
func hMasterOutPending() {

}

func hMasterOutTry() {
	type CDATA struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
		To     string `json:"to"`
	}

	method := chttp.POST
	url := model.V2 + "/master/out/try"
	Doc().Comment("[ 마스터출금 ] 마스터 계좌에서 외부계좌로 출금신청 (ROOT 권한)").
		Method(method).URL(url).
		JParam(CDATA{},
			"symbol", "ETH , GDG",
			"price", "가격",
			"to", "받을 주소",
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

			if model.MasterCallback(scv.MasterOutCallback) == false {
				chttp.Fail(w, ack.NotAllowSystem)
				return
			}

			admin := model.GetTokenAdmin(req)
			if admin.IsRoot == false {
				chttp.Fail(w, ack.InvalidRootAdmin)
				return
			}

			cdata := CDATA{}
			chttp.BindingJSON(req, &cdata)

			if jmath.CMP(cdata.Price, 0) <= 0 {
				chttp.Fail(w, ack.UnderZERO)
				return
			}

			cdata.To = dbg.TrimToLower(cdata.To)
			if ecsx.IsAddress(cdata.To) == false {
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

				chttp.OK(w, nil)

			})
		},
	)

}
