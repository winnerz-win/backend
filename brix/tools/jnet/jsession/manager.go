package jsession

import (
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
)

//var cookieHandler *securecookie.SecureCookie = nil

// CookieHandler :
type CookieHandler *securecookie.SecureCookie

// NewSession :
func NewSession() CookieHandler {
	cookieHandler := securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
	return cookieHandler
}

// GetSession :
func GetSession(handler CookieHandler, key string, value interface{}, req *http.Request) error {
	cHandler := (*securecookie.SecureCookie)(handler)
	if cookie, err := req.Cookie(key); err == nil {
		return cHandler.Decode(key, cookie.Value, value)
	} else {
		return err
	}
}

// SetSession :
func SetSession(handler CookieHandler, key string, value interface{}, res http.ResponseWriter) {
	cHandler := (*securecookie.SecureCookie)(handler)
	if encoded, err := cHandler.Encode(key, value); err == nil {
		cookie := &http.Cookie{
			Name:     key,
			Value:    encoded,
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
		}
		http.SetCookie(res, cookie)
	}
}

//SetSessionExpired :
func SetSessionExpired(handler CookieHandler, key string, value interface{}, expired time.Duration, res http.ResponseWriter) {
	cHandler := (*securecookie.SecureCookie)(handler)
	if encoded, err := cHandler.Encode(key, value); err == nil {
		cookie := &http.Cookie{
			Name:     key,
			Value:    encoded,
			Path:     "/",
			Expires:  time.Now().UTC().Add(expired),
			SameSite: http.SameSiteNoneMode,
		}
		http.SetCookie(res, cookie)
	}
}

// ClearSession :
func ClearSession(key string, res http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(res, cookie)
}
