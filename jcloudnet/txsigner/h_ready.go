package txsigner

import (
	"jtools/dbg"
	"jtools/jnet/chttp"
	"net/http"
)

var (
	handle = chttp.HANDLERLIST{}
)

func _ready(classic *chttp.Classic) chttp.Starter {
	rtx := chttp.NewStarter("api")
	classic.SetHandlerFunc(handlerFunc(classic))
	classic.SetHandlers(handle)
	return rtx
}

func handlerFunc(classic *chttp.Classic) chttp.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		defer dbg.VRecover(func() {
			chttp.Fail(
				classic.ResponsWriter(w),
				ERROR_HandlerPanic,
				r.URL.Path,
			)
		})

		next(w, r)
	}

}
