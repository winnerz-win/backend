package api

import (
	"net/http"
	"strings"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/txm/ack"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
)

func handlerFunc(classic *chttp.Classic) chttp.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		defer dbg.VRecover(func() {
			chttp.Fail(w, ack.HandlerPanic, r.URL.Path)
		})

		path := r.URL.Path
		isAccessCheck := strings.HasPrefix(path, model.V1)

		if isAccessCheck {
			config := inf.Config()
			if config.IPCheck {
				remoteaddress := chttp.RemoteIPPort(r)
				remoteIP := remoteaddress[0]
				remotePort := remoteaddress[1]

				if remoteIP != inf.ClientHostIP() {
					dbg.Red(remoteIP, ":", remotePort, " is blocked!")
					return
				}

				dbg.Purple("ip_allow :", remoteIP)
			}
		}

		next(w, r)

	}
}
