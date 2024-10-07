package chttp

import (
	"jtools/dbg"
	"net/http"
)

type ResultFormat struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func (ResultFormat) TagString() []string {
	return []string{
		"success", "성공/실패",
		"data", "성공:요청에따른 응답, 실패:error_code , error_message",
	}
}

func (my ResultFormat) String() string { return dbg.ToJsonString(my) }

func OK(w ResponseWriter, v interface{}) {
	r := ResultFormat{
		Success: true,
		Data:    v,
	}
	if r.Data == nil {
		r.Data = struct{}{}
	}
	w.R().JSON(w.W(), http.StatusOK, r)
}
func JsonOK(w ResponseWriter, v interface{}) {
	w.R().JSON(w.W(), http.StatusOK, v)
}
func JsonFail(w ResponseWriter, v interface{}) {
	w.R().JSON(w.W(), http.StatusNotFound, v)
}
func JsonStatus(w ResponseWriter, status int, v interface{}) {
	w.R().JSON(w.W(), status, v)
}

type FailFormat struct {
	ErrorCode    int         `json:"error_code"`
	ErrorMessage string      `json:"error_message"`
	ETC          interface{} `json:"etc,omitempty"`
}

func (FailFormat) TagString() []string {
	return []string{
		"error_code", "실패 코드",
		"error_message", "실패 사유",
		"etc", "기타 데이타",
	}
}

func Fail(w ResponseWriter, e ERROR, etc ...interface{}) {
	data := FailFormat{
		ErrorCode:    e.Code(),
		ErrorMessage: e.Desc(),
	}
	if len(etc) > 0 {
		data.ETC = etc[0]
	}
	r := ResultFormat{
		Success: false,
		Data:    data,
	}
	if errorView {
		msg := dbg.Cat("chttp.AckFail [", e.Code(), "]", e.Desc())
		if data.ETC != nil {
			msg += " , ETC :" + dbg.ToJsonString(data.ETC)
		}
		LogError(msg)
	}
	w.R().JSON(w.W(), http.StatusOK, r)
}

func Text(w ResponseWriter, text string) {
	w.R().Text(w.W(), http.StatusOK, text)
}
func Bytes(w ResponseWriter, buf []byte) {
	w.R().Data(w.W(), http.StatusOK, buf)
}
func HTML(w ResponseWriter, html string, binding interface{}) {
	w.R().HTML(w.W(), http.StatusOK, html, binding)
}
func Redirect(w ResponseWriter, req *http.Request, path string) {
	http.Redirect(w.W(), req, path, http.StatusFound)
}

var (
	ERROR_ETC = Error(404, "data.etc(struct)")
)

func FailETC(w ResponseWriter, etc interface{}) {
	e := ERROR_ETC
	data := FailFormat{
		ErrorCode:    e.Code(),
		ErrorMessage: e.Desc(),
	}
	if etc != nil {
		data.ETC = etc
	}
	r := ResultFormat{
		Success: false,
		Data:    data,
	}
	LogError("chttp.AckFail [", e.Code(), "]", e.Desc(), ":", dbg.ToJsonString(etc))
	w.R().JSON(w.W(), http.StatusOK, r)
}
