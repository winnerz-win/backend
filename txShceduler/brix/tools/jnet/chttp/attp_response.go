package chttp

import (
	"errors"
	"net/http"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jpath"

	"github.com/unrolled/render"
)

var (
	renderHTMLDebug  = false // SetHTMLRenderDebug(true)
	isResponseFormat = false // SetAckFormat()
)

//SetHTMLRenderDebug :
func SetHTMLRenderDebug(f ...bool) {
	if dbg.IsTrue2(f...) {
		renderHTMLDebug = true
	}
}

//SetAckFormat : 리스폰스 응답 200 룰 (신규)
func SetAckFormat(isErrView ...bool) {
	isResponseFormat = true
	if dbg.IsTrue2(isErrView...) {
		SetErrorView(true)
	}
}

type AckFormat struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func (AckFormat) TagString() []string {
	return []string{
		"success", "성공/실패",
		"data", "성공:요청에따른 응답, 실패:error_code , error_message",
	}
}

func (my AckFormat) Parse(p interface{}) error {
	if my.Success == false {
		return errors.New("ack_is_fail")
	}
	if my.Data == nil {
		return errors.New("ack_data_is_nil")
	}
	return dbg.ChangeStruct(my.Data, p)
}

func (my AckFormat) FailData() AckFailData {
	faildata := AckFailData{}
	if my.Success == true {
		return faildata
	}
	dbg.ChangeStruct(my.Data, &faildata)
	return faildata
}

//OK :
func OK(w http.ResponseWriter, v interface{}) {
	if isResponseFormat {
		ack := AckFormat{
			Success: true,
			Data:    v,
		}
		if ack.Data == nil {
			ack.Data = struct{}{}
		}
		Renderer.JSON(w, http.StatusOK, ack)
		return
	}
	////////////////////////////////////////////////////////////
	Renderer.JSON(w, http.StatusOK, v)
}
func JsonOK(w http.ResponseWriter, v interface{}) {
	Renderer.JSON(w, http.StatusOK, v)
}

type AckFailData struct {
	ErrorCode    int         `json:"error_code"`
	ErrorMessage string      `json:"error_message"`
	ETC          interface{} `json:"etc,omitempty"`
}

func (AckFailData) TagString() []string {
	return []string{
		"error_code", "실패 코드",
		"error_message", "실패 사유",
		"etc", "기타 데이타",
	}
}

//Fail :
func Fail(w http.ResponseWriter, e CError, etc ...interface{}) {
	if isResponseFormat {
		fail := AckFailData{
			ErrorCode:    e.Code(),
			ErrorMessage: e.Desc(),
		}
		if len(etc) > 0 {
			fail.ETC = etc[0]
		}
		ack := AckFormat{
			Success: false,
			Data:    fail,
		}
		if errorView {
			dbg.RedItalic("chttp.AckFail [", e.Code(), "]", e.Desc())
		}
		Renderer.JSON(w, http.StatusOK, ack)
		return
	}
	////////////////////////////////////////////////////////////
	v := JsonType{
		"code": e.Code(),
	}
	if len(etc) > 0 {
		if err, do := etc[0].(error); do {
			v["etc"] = err.Error()
		} else {
			v["etc"] = etc[0]
		}
	}
	if errorView {
		dbg.RedItalic("chttp.Fail ----- start")
		for key, val := range v {
			dbg.RedItalic(key, ":", val)
		} //for
		dbg.RedItalic("chttp.Fail ----- end")
	}
	Renderer.JSON(w, http.StatusBadRequest, v)
}

//Text :
func Text(w http.ResponseWriter, text string) {
	Renderer.Text(w, http.StatusOK, text)
}

func DATA(w http.ResponseWriter, buf []byte) {
	Renderer.Data(w, http.StatusOK, buf)
}

//HTML :
func HTML(w http.ResponseWriter, html string, binding interface{}) {
	if renderHTMLDebug {
		renderPath := jpath.NowPath() + "\\" + RenderRootPath
		Renderer = render.New(render.Options{
			Directory:  renderPath,
			Extensions: []string{".tmpl", ".html"},
		})
	}
	Renderer.HTML(w, http.StatusOK, html, binding)
}

//REDIRECT :
func REDIRECT(w http.ResponseWriter, req *http.Request, path string) {
	http.Redirect(w, req, path, http.StatusFound)
}

//ErrorHelp :
func ErrorHelp() string {
	help := ""
	if isResponseFormat {
		help = `
		-------- HTTP - ACK 응답 Format --------
	모든 API 요청에 따른 응답 응답형식은 application/json 형식이며 아래와 같은 포멧으로 응답한다.
	 {
		 "success" : true / false,
		 "data" : {}
	 }
	 
	 ----------------------------------------------------------------------------------------

	 <성공 예시>
	 {
		 "success" : true,
		 "data" : {
			"receipt_code" : "receipt_976f1ca6cfbd87549fc066cae60c5e3a67d24753"
		 }
	 }
	 --> "success"가 true이면 요청결과 성공이며 "data"필드 안에는 각 요청에 맞는 응답 데이타 결과가 들어있다.
	 
	 ----------------------------------------------------------------------------------------

	 <실패 예시>
	 {
		"success" : false,
		"data" : {
			"error_code" : 4001,
			"error_message" : "요청 데이터를 찾을수 없음"
		}
	 }
	 --> "success"가 false이면 요청결과는 실패이며 "data"필드안에 "error_code"(코드)와 "error_message"(사유)가 들어있다.
	 
	 ----------------------------------------------------------------------------------------
	` + ErrorStrings() + `
	`
	} else {
		help = `
		-------- HTTP - 응답 Format --------
		` + ErrorStrings() + `
		`
	}
	return help
}
