package model

import (
	"context"
	"net/http"
	"time"

	"jtools/mms"
	"txscheduler/brix/tools/jtoken"
)

const verifyAdminTokenKey = "TXM/admin,8v8_NIMIchuiPAlomA@v@"

//ValidTokenAdmin :
func ValidTokenAdmin(str string) (jtoken.AuthToken, jtoken.Error) {
	token := jtoken.StringToToken(str)
	return token.Valid([]byte(verifyAdminTokenKey))
}

func (my Admin) getToken() map[string]interface{} {
	return map[string]interface{}{
		"is_root": my.IsRoot,
	}
}

func tokenToAdmin(token jtoken.AuthToken) Admin {
	mapper := token.ETC().(map[string]interface{})
	admin := Admin{
		IsRoot: mapper["is_root"].(bool),
	}
	return admin
}

//MakeTokenForAdmin :
func MakeTokenForAdmin(admin *Admin, tokentype jtoken.Type) jtoken.Token {
	var exp time.Time
	if tokentype == jtoken.AccessToken {
		exp = mms.NowTime().Add(time.Minute * 30)
	} else {
		exp = mms.NowTime().Add(time.Minute * 90)
	}
	token, _ := jtoken.Make(
		Platform,
		[]byte(verifyAdminTokenKey),
		tokentype,
		admin.Name,
		exp,
		admin.getToken(),
	)
	return token
}

const cUSERTOKENADMIN = "jwt_token_admin_user"

//TokenRequestAdmin :
func TokenRequestAdmin(req *http.Request, token jtoken.AuthToken) *http.Request {
	cv := tokenToAdmin(token)
	ctx := context.WithValue(req.Context(), cUSERTOKENADMIN, cv)
	return req.WithContext(ctx)
}

//GetTokenAdmin :
func GetTokenAdmin(req *http.Request) Admin {
	return req.Context().Value(cUSERTOKENADMIN).(Admin)
}
