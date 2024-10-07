package api

import (
	"net/http"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/jnet/doc"
	"txscheduler/txm/inf"
)

var dc doc.Object

// Doc :
func Doc() doc.Object {
	if dc == nil {
		dc = doc.NewObjecter("TX_SCHEDULER", "API LIST", inf.CoreVersion)
	}
	return dc
}

// DocEnd :
func DocEnd(classic *chttp.Classic) {
	if dc == nil {
		return
	}
	// dc.Update()
	// dc = nil

	if !inf.Mainnet() {
		classic.SetHandler(
			chttp.GET, "/doc/api",
			func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
				w.Write(dc.Bytes())
			},
		)
	}

}
