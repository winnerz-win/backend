package admin

import (
	"net/http"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/jnet/doc"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
	"txscheduler/txm/pwd"
)

var dc doc.Object

// Doc :
func Doc() doc.Object {
	if dc == nil {
		dc = doc.NewObjecter("TX_SCHEDULER.admin", "ADMIN API LIST", inf.CoreVersion)
		prefixDocHelp(dc)
	}
	return dc
}

// DocEnd :
func DocEnd(classic *chttp.Classic) {
	if dc == nil {
		return
	}
	// dc.Update()
	// dc = nil
	if !inf.Mainnet() {
		classic.SetHandler(
			chttp.GET, "/doc/admin",
			func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
				w.Write(dc.Bytes())
			},
		)
	}

}

func prefixDocHelp(dc doc.Object) {

	dc.Message(`
		---- ADMIN API ---
		헤더 토큰 : `+model.HeaderAdminToken+`

		비밀번호 Salt값은 프로젝트 별로 매번 다르므로 확인 요청 바람.
		Salt    : `+pwd.SaltVale()+`
		Iter    : `+dbg.Cat(pwd.IterCount)+`
		KeySize : `+dbg.Cat(pwd.KeySize)+`
		 
		최초 관리자 ID/PW
		ID : `+model.AdminDefaultName+`
		PW : `+model.AdminDefaultPWD+`

		----------------------------------------------------------------------------------------		
		모든 요청의 응답결과 형식은 아래와 같다.

		<성공응답>
		{
			"success" : true,
			"data" : {
				 format...
			},
		}

		<실패응답>
		{
			"success" : false,
			"data" : {
				"error_code" : int,
				"error_message" : string,
			}
		}
		----------------------------------------------------------------------------------------		
		`+chttp.ErrorStrings()+`

		----------------------------------------------------------------------------------------		
		 DB 테이블 설명
		 member : 가입 회원 테이블
		 {
			 data {} // 플랫폼서버가 회원가입시 추가적으로 저장하는 데이터 영역 (전화번호 등등)
			 coin {} // 회원 주소가 보유하고 있는 코인수량 (이더스캔과 동기화)
		 }

		 info_deposit : 입금계좌에서 마스터지갑으로 옮기기 위한 코인별 최소 수량. 없으면 base_value 값으로 대체

		 log_deposit : 회원의 코인 입금내역
		 {
			 deposit_result  bool  // 트랙젝션 성공/실패 여부
			 is_send  bool // 플랫폼 서버로 해당 내역 전송 완료 여부
		 }

		 log_withdraw : 회원의 코인 출금내역
		 {
			state int  // 출금성공:200 , 출금실패:104 , 출금대기:1
			is_send  bool // 플랫폼 서버로 해당 내역 전송 완료 여부
		 }

		 tx_eth_withdraw : 회원의 코인 출금신청 대기열

	`, doc.Blue)
}
