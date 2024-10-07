package txsigner

import (
	"jtools/cc"
	"jtools/dbg"
	"jtools/jargs"
	"jtools/jlog"
	"jtools/jmath"
	"jtools/jnet/chttp"
)

var (
	opt    Option
	config *Config
)

func StartServer(_option Option) {
	opt = _option
	cc.YellowItalic("txsigner.StartServer")
	cc.Yellow(_option)

	arg := jargs.New()

	port := 80
	if !arg.Next("--port", func(val string) {
		port = jmath.Int(val)
	}) {
		_Exit("port is emtpy")
	}
	config_file_name := ""
	if !arg.Next("--config", func(val string) {
		config_file_name = val
	}) {
		_Exit("config option is empty")
	}

	config = LoadConfig(config_file_name)

	ctx := chttp.Context()
	ctx.SetPort(port)
	ctx.SetName(config.TITLE)

	ctx.SetHeaders(
		"Orign",
	)

	chttp.SetErrorView()
	classic := chttp.New(
		ctx,
		jlog.GetEntry(),
	)
	ctx.SetStarter(
		_ready(classic),
	)

	classic.SetVersion(func() string {
		message := config.TITLE + dbg.ENTER
		message += `
		────────────────────────────────────────────────────
		 Option :
		 ` + _option.String() + `
		 Port : ` + dbg.Cat(port) + `
		`
		return message
	})

	jlog.Info("[", _option.InfraTag, "] SIGN_SERVER_START :::::::::::::::::::::::::::::::::::::")
	classic.StartDirect()
}
