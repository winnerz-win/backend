package admin

import (
	"net/http"
	"strings"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/jtoken"
	"txscheduler/txm/ack"
	"txscheduler/txm/model"
)

func handlerFunc(classic *chttp.Classic) chttp.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		defer dbg.VRecover(func() {
			chttp.Fail(w, ack.HandlerPanic, r.URL.Path)
		})

		path := r.URL.Path
		//dbg.Purple(path)
		isAccessCheck := strings.HasPrefix(path, model.V2)

		if isAccessCheck {
			tokenString := r.Header.Get(model.HeaderAdminToken)
			authtoken, err := model.ValidTokenAdmin(tokenString)
			if err != nil {
				if jtoken.ExpiredError(err) {
					chttp.Fail(w, ack.TokenExpired)
				} else { //jtoken.ErrorTokenInvalid
					chttp.Fail(w, ack.InvalidToeken)
				}
				return
			}
			r = model.TokenRequestAdmin(r, authtoken)
		}
		next(w, r)
	}
}
