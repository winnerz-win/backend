package api

import (
	"net/http"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/mms"
	"txscheduler/txm/ack"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func init() {
	hSign()
}

func docSign() string {
	return `
	<cc_blue>< 회원 가입 예제 ></cc_blue> 
	플랫폼 서버에서 회원가입을 하면 가입한 회원에게 코인을 입금할 수 있는 주소를 발급 하기 위해서
	블록체인 연동서버에 회워정보를 주어서 입금주소를 발급 받는다. ( 회원ID == 입금주소 )

	<cc_bold>URL></cc_bold> ` + docURL + `/v1/sign

	<cc_bold>Method</cc_bold> : POST

	<cc_purple>요청파라미터</cc_purple> : 
	{
		"name" : "` + docUser + `",		
	}

	<cc_purple>성공응답</cc_purple> :
	{
		"success" : true,
		"data" : {
			"uid" : ` + docID + `,
			"address" : "` + docAddress + `",
			"create_at" : 1618303423126
		}
	}
	
	<cc_bold>uid</cc_bold> : 회원가입시 발금되는 고유ID
	<cc_bold>address</cc_bold> : 회원이 ETH나 플랫폼서버와 약속된 ERC20토큰을 입금하는 계좌 주소
	<cc_bold>create_at</cc_bold> : 입금주소 발글 시간 (mms) 
	`
}

func hSign() {
	type CDATA struct {
		Name string                 `json:"name"`
		Data map[string]interface{} `json:"data"`
	}
	type RESULT struct {
		UID      int64   `json:"uid"`
		Address  string  `json:"address"`
		CreateAt mms.MMS `json:"create_at"`
	}

	method := chttp.POST
	url := model.V1 + "/sign"

	Doc().Comment("[ 회원가입(입금주소 생성) ] 가입 요청").
		Method(method).URL(url).
		JParam(CDATA{},
			"name", "계정이름 (이메일, UID등등)",
			"data", "추가로 저장할 부가 데이타(전화번호 등등)",
		).
		JResultOK(chttp.AckFormat{}).
		Etc(".", `_
			<cc_blue>성공응답</cc_blue> :
			{
				"success" : true,
				"data" : {
					"uid": 1001,					// 회원 키값
					"address" : "0x....."			// 회원의 가상계좌 ( 코인 입금 주소 )
					"create_at" : 1617349692884		// 가입시간(mms)
				}
			}

			<cc_blue>실패응답</cc_blue> :
			{
				"success" : false,
				"error_code" : int,
				"error_message" : string,
			}
		`).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := CDATA{}
			chttp.BindingJSON(req, &cdata)
			model.Trim(&cdata.Name)

			if cdata.Name == "" {
				chttp.Fail(w, ack.BadParam)
				return
			}

			if cdata.Data == nil || len(cdata.Data) == 0 {
				cdata.Data = map[string]interface{}{}
			}

			model.CreateLock(func(db mongo.DATABASE) {
				c := db.C(inf.COLMember)
				if cnt, _ := c.Find(mongo.Bson{"name": cdata.Name}).Count(); cnt > 0 {
					chttp.Fail(w, ack.ExistedName)
					return
				}

				nowAt := mms.Now()
				uid := model.CurrentUID()
				wallet := inf.Wallet(uid)
				member := model.Member{
					User: model.User{
						UID:     uid,
						Address: wallet.Address(),
						Name:    cdata.Name,
					},
					Data:      cdata.Data,
					Coin:      model.NewCoinDataSymbol(inf.SymbolList()...),
					Deposit:   model.NewCoinDataSymbol(inf.SymbolList()...),
					Withdraw:  model.NewCoinDataSymbol(inf.SymbolList()...),
					CreateAt:  nowAt,
					CreateYMD: nowAt.YMD(),
					Timestamp: nowAt,
					YMD:       nowAt.YMD(),
				}

				if c.Insert(member) != nil {
					chttp.Fail(w, ack.DBJob)
					return
				}

				model.IncUID()

				chttp.OK(w, RESULT{
					UID:      member.UID,
					Address:  member.Address,
					CreateAt: member.CreateAt,
				})
			})

		},
	)
}
