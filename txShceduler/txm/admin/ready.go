package admin

import (
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/runtext"
)

var handle = chttp.PContexts{}

// Ready :
func Ready(classic *chttp.Classic) runtext.Starter {
	rtx := runtext.New("admin")

	classic.SetHandlerFunc(handlerFunc(classic))
	classic.SetContextHandles(handle)

	DocEnd(classic)
	return rtx
}
