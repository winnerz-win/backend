package txsigner

import "jtools/jnet/chttp"

var (
	ERROR_HandlerPanic  = chttp.Error(1, "Server handler panic error")
	ERROR_SignTx        = chttp.Error(400, "Fail SignTx")
	ERROR_MarshalBinary = chttp.Error(401, "Fail MarshalBinary")
)
