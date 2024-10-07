package doc

import "fmt"

const (
	StateNone = 0
	StatusOK  = 200 // 성공

	StatusBadRequest    = 400 // 잘못된 요청 ( 주로 클라에서 보내주는 데이터가 잘못 되었을 경우.)
	StatusUnauthorized  = 401 // 인증 에러
	StatusForbidden     = 403 // 접근 금지 ( 블록 유저 )
	StatusNotFound      = 404 // 클라가 요청한 리소스가 서버에 없음. ( 유저정보 없음 등등.)
	StatusNotAcceptable = 406 // 요청한 콘텐츠를 찾지 못하였다.
	StatusConflict      = 409 // 이 응답은 요청이 현재 서버의 상태와 충돌될 때 보냅니다.

	StatusInternalServerError = 500 // 서버 DB 오류.
	StatusServiceUnavailable  = 503 // 정검중.
)

//getStatusString :
func getStatusString(status int) string {
	switch status {
	case StatusOK:
		return "StatusOK"
	case StatusBadRequest:
		return "StatusBadRequest"
	case StatusUnauthorized:
		return "StatusUnauthorized"
	case StatusForbidden:
		return "StatusForbidden"
	case StatusNotFound:
		return "StatusNotFound"
	case StatusNotAcceptable:
		return "StatusNotAcceptable"
	case StatusConflict:
		return "StatusConflict"
	case StatusInternalServerError:
		return "StatusInternalServerError"
	case StatusServiceUnavailable:
		return "StatusServiceUnavailable"
	}
	return fmt.Sprintf("%v", status)
}
