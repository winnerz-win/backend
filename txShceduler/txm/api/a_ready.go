package api

import (
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/runtext"
	"txscheduler/txm/aadev"
	"txscheduler/txm/inf"
)

var handle = chttp.PContexts{}

func init() {
	prefixDocHelp()

}

// Ready :
func Ready(classic *chttp.Classic) runtext.Starter {
	rtx := runtext.New("api")

	classic.SetHandlerFunc(handlerFunc(classic))

	is_onwer_task_mode := inf.IsOnwerTaskMode()
	if is_onwer_task_mode {
		applyOnwerHandles()
	}
	classic.SetContextHandles(handle)

	set_base_callback_api_doc()
	if is_onwer_task_mode {
		h_owner_callback_api_doc()
	}

	DocEnd(classic)

	return rtx
}

const (
	docUser    = `test@gmail.com`
	docID      = "1001"
	docAddress = `0xabcd...`
	docURL     = aadev.DEV_URL
)

func prefixDocHelp() {

	Doc().Message(`
		<cc_purple>--- 스케줄러 서버 API 설명서 ---</cc_purple>

	 모든 API 요청에 따른 응답 응답형식은 application/json 형식이며 아래와 같은 포멧으로 응답한다.
	 <cc_bold>
	 {
		 "success" : true / false,
		 "data" : {}
	 }</cc_bold>
	 
	 ----------------------------------------------------------------------------------------

	 <cc_blue><성공 예시></cc_blue>
	 <cc_bold>
	 {
		 "success" : true,
		 "data" : {
			"receipt_code" : "receipt_976f1ca6cfbd87549fc066cae60c5e3a67d24753"
		 }
	 }</cc_bold>	 
	 --> "success"가 true이면 요청결과 성공이며 "data"필드 안에는 각 요청에 맞는 응답 데이타 결과가 들어있다.
	 
	 ----------------------------------------------------------------------------------------

	 <cc_blue><실패 예시></cc_blue>
	 <cc_bold>
	 {
		"success" : false,
		"data" : {
			"error_code" : 4001,
			"error_message" : "존재하는 Name(유저)입니다."
		}
	 }</cc_bold>
	 --> "success"가 false이면 요청결과는 실패이며 "data"필드안에 "error_code"(코드)와 "error_message"(사유)가 들어있다.
	 
	 ----------------------------------------------------------------------------------------

	 ` + chttp.ErrorStrings() + `

	 ----------------------------------------------------------------------------------------
	 ` + docSign() + ` 

	 ----------------------------------------------------------------------------------------
	 ` + docMemberInfo() + `
	`)
}
