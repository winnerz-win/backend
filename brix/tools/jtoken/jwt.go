package jtoken

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"txscheduler/brix/tools/dbg"

	jwt "github.com/dgrijalva/jwt-go"
)

//AuthClaims :
type AuthClaims struct {
	*jwt.StandardClaims
	Etc interface{} `json:"etc,omitempty"`
}

//Claims :
func (my *AuthClaims) Claims() *jwt.StandardClaims {
	return my.StandardClaims
}

//ETC :
func (my *AuthClaims) ETC() interface{} {
	return my.Etc
}

//Type :
func (my *AuthClaims) Type() Type {
	return Type(my.Subject)
}

//GetID :
func (my AuthClaims) GetID() string {
	return my.Id
}

//IssuedTime :
func (my *AuthClaims) IssuedTime() time.Time {
	return time.Unix(my.IssuedAt, 0)
}

//ExpiresTime :
func (my *AuthClaims) ExpiresTime() time.Time {
	return time.Unix(my.ExpiresAt, 0)
}

/////////////////////////////////////////////////////////////////////////////////

//STRToken :
type cToken string

//ToString :
func (my cToken) ToString() string {
	return string(my)
}

type viewToken struct {
	PHeader *viewHeader `json:"header"`
	PBody   *AuthClaims `json:"body"`
}

//Header :
func (my viewToken) Header() *viewHeader {
	return my.PHeader
}

//Body :
func (my viewToken) Body() AuthToken {
	return my.PBody
}

type viewHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

func (my cToken) ViewToken(isLog ...bool) ViewToken {
	logView := false
	if len(isLog) > 0 && isLog[0] {
		logView = true
	}
	vt := &viewToken{}

	ss := strings.Split(my.ToString(), ".")
	if len(ss) != 3 {
		return nil
	}

	header := ss[0]
	body := ss[1]
	secret := ss[2]
	_ = secret

	hbuf, _ := base64.StdEncoding.DecodeString(header)
	bbuf, _ := base64.StdEncoding.DecodeString(body)

	vHeader := &viewHeader{}
	if err := json.Unmarshal(hbuf, vHeader); err != nil {
		return nil
	}

	if logView {
		dbg.ViewJSONTag(vHeader, "vHeader")
	}

	chkJSON := func(msg []byte) []byte {
		iCnt := 0
		oCnt := 0
		for i := 0; i < len(msg); i++ {
			if msg[i] == '{' {
				iCnt++
			} else if msg[i] == '}' {
				oCnt++
			}
		} //for
		if iCnt > oCnt {
			zCnt := iCnt - oCnt
			for ; zCnt > 0; zCnt-- {
				msg = append(msg, '}')
			}
		}
		return msg
	}
	bbuf = chkJSON(bbuf)
	vBody := &AuthClaims{
		StandardClaims: &jwt.StandardClaims{},
	}
	if err := json.Unmarshal(bbuf, vBody); err != nil {
		return nil
	}

	if logView {
		dbg.ViewJSONTag(vBody, "vBody")
	}

	vt.PHeader = vHeader
	vt.PBody = vBody

	return vt
}

func (my cToken) Type() Type {
	if my.ViewToken().Body().Claims().Subject == SubjectAccess {
		return AccessToken
	}
	return RefreshToken
}

//Valid :
func (my cToken) Valid(verifyKey []byte) (AuthToken, Error) {
	//tokenString := my.ToString()
	// if tokenString == "" {
	// 	return nil, tokenError{
	// 		ErrorMsg:     "string is Empty",
	// 		ErrorID:      "error",
	// 		ErrorSubject: Type("error"),
	// 		ErrorExpire:  common.NowTime(),
	// 		claims:       nil,
	// 	}
	// }
	token, err := jwtParseWithClaims(my.ToString(), verifyKey, &AuthClaims{})
	if err != nil {
		terror := tokenError{
			ErrorMsg:     err.Error(),
			ErrorID:      "error",
			ErrorSubject: Type("error"),
			ErrorExpire:  time.Now().UTC(),
			claims:       nil,
		}
		if err.Error() == ErrorTokenExpired || err.Error() == ErrorTokenInvalid {
			terror.claims = token.Claims.(*AuthClaims)
		}

		return nil, terror
	}
	ac := token.Claims.(*AuthClaims)
	return ac, nil
}

func jwtParseWithClaims(tokenString string, verifyKey []byte, claims jwt.Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(*jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
}

/////////////////////////////////////////////////////////////////////////////////

//TokenError :
type tokenError struct {
	ErrorMsg     string    `json:"error"`
	ErrorID      string    `json:"ID"`
	ErrorSubject Type      `json:"Subject"`
	ErrorExpire  time.Time `json:"ExpireAt"`
	claims       *AuthClaims
}

//ToString :
func (my tokenError) ToString() string {
	return my.ErrorMsg
}

//ID :
func (my tokenError) ID() string {
	return my.ErrorID
}

//TokenType :
func (my tokenError) TokenType() Type {
	return Type(my.ErrorSubject)
}

//ExpiresTime :
func (my tokenError) ExpiresTime() time.Time {
	return my.ErrorExpire
}
