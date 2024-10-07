package jtoken

import (
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	//ErrorTokenExpired : error message
	/*
		{
			"Inner": {
				"Now": 1570528934,
				"ExpiredBy": 3600259492000,
				"Claims": {
				"exp": 1570525334,
				"foo": "bar"
				}
			},
			"Errors": 16
		}
	*/
	ErrorTokenExpired = "Token is expired"

	//ErrorTokenInvalid : error message
	/*
		{
			"Inner": {},
			"Errors": 4
		}
	*/
	ErrorTokenInvalid = "signature is invalid"

	//SubjectAccess : 엑세스 토큰
	SubjectAccess = "access"
	//SubjectRefresh :	갱신용 토큰
	SubjectRefresh = "refresh"
)

//CheckError :
func CheckError(errString string) string {
	if strings.HasPrefix(errString, "token is expired") {
		return ErrorTokenExpired
	}
	return ErrorTokenInvalid
}

//ExpiredError :
func ExpiredError(err Error) bool {
	if strings.HasPrefix(err.ToString(), "token is expired") {
		return true
	}
	return false
}

//Type :
type Type string

//ToString :
func (my Type) ToString() string {
	return string(my)
}

//AccessToken :
var AccessToken = Type(SubjectAccess)

//RefreshToken :
var RefreshToken = Type(SubjectRefresh)

//AuthToken : AuthClaims
type AuthToken interface {
	Claims() *jwt.StandardClaims
	ETC() interface{}
	Type() Type
	GetID() string
	IssuedTime() time.Time
	ExpiresTime() time.Time
}

//Token : cToken
type Token interface {
	ToString() string
	Valid(verifyKey []byte) (AuthToken, Error)
	ViewToken(isLog ...bool) ViewToken
	Type() Type
}

//Error : tokenError
type Error interface {
	ToString() string
	ID() string
	TokenType() Type
	ExpiresTime() time.Time
}

//ViewToken :
type ViewToken interface {
	Header() *viewHeader
	Body() AuthToken
}

//MakeToken :
func MakeToken(verifyKey []byte, subject Type, id string, exp time.Time, etc interface{}) (Token, error) {
	return Make("txscheduler/brix-cc", verifyKey, subject, id, exp, etc)
}

//Make :
func Make(issuer string, verifyKey []byte, subject Type, id string, exp time.Time, etc interface{}) (Token, error) {
	claims := AuthClaims{
		StandardClaims: &jwt.StandardClaims{
			Issuer:    issuer,
			IssuedAt:  time.Now().UTC().Unix(),
			ExpiresAt: exp.Unix(),
			Id:        id,
			Subject:   string(subject),
		},
		Etc: etc,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(verifyKey)
	if err != nil {
		return nil, err
	}
	return cToken(tokenString), nil
}

//StringToToken :
func StringToToken(str string) Token {
	str = strings.TrimSpace(str)
	return cToken(str)
}
