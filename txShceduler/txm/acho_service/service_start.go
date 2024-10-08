package acho_service

import (
	"jtools/jmath"
	"net/http"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/txm/inf"
)

const (
	APIDeposit  = "/deposit"
	APIWithdraw = "/withdraw"
)

// StartService :
func StartService(c *inf.IConfig) {
	if c.Mainnet == true {
		return
	}
	go run(c)
}

func run(config *inf.IConfig) {
	dbg.PrintForce("acho_service_server ::::: run")

	port := 80
	portString := config.ClientHost[config.Mainnet][1]
	if portString != "" {
		port = jmath.Int(portString)
	}

	corsHeader := []string{
		"Access-Control-Allow-Origin",
	}
	classic := chttp.NewSimple(port, corsHeader, "Service")

	handle := &chttp.PContexts{}
	hDeposit(handle)
	hWithdraw(handle)
	classic.SetContextHandles(*handle)

	classic.StartSimple()
}

func hDeposit(handle *chttp.PContexts) {
	handle.Append(
		chttp.POST, APIDeposit,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := chttp.JsonType{}
			chttp.BindingJSON(req, &cdata)

			dbg.Green("< service.Deposit >")
			dbg.Green(cdata)
			chttp.ResultJSON(w, chttp.StatusOK, nil)
		},
	)
}

func hWithdraw(handle *chttp.PContexts) {
	handle.Append(
		chttp.POST, APIWithdraw,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := chttp.JsonType{}
			chttp.BindingJSON(req, &cdata)

			dbg.Green("< service.Withdraw >")
			dbg.Green(cdata)
			chttp.ResultJSON(w, chttp.StatusOK, nil)
		},
	)
}
