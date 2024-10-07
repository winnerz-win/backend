package chttp

import (
	"context"
	"jtools/jnet/chttp/cors"
	"jtools/jnet/console"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

const (
	SystemWriteKey = "classic_system_writer"
)

type SystemWriterFunc func(v ...interface{})

func SystemWriteln(r *http.Request, v ...interface{}) {
	if len(v) == 0 {
		return
	}
	if val := r.Context().Value(SystemWriteKey); val != nil {
		if writeln, do := val.(SystemWriterFunc); do {
			writeln(v...)
		}
	}
}

type Classic struct {
	context *iContext

	n        *negroni.Negroni
	router   *httprouter.Router
	renderer *render.Render

	apiCount int
	items    HANDLERLIST

	callNav  map[string]int
	callCnt  []int32
	call_log TextAlignTapperKey

	isStart bool
	isEnd   bool
	isSSL   bool

	endCtx context.Context

	mu sync.RWMutex

	//logger     *log.Logger
	systemWriter SystemWriterFunc

	console   console.IConsole
	f_version func() string
}

func (my *Classic) String() string {
	defer my.mu.RUnlock()
	my.mu.RLock()

	v := my.context.getMap()
	v["apiCount"] = my.apiCount

	return ToJsonString(v)
}

func (my *Classic) Context() *iContext {
	return my.context
}
func (my *Classic) SystemWriter() SystemWriterFunc {
	if my.systemWriter == nil {
		return LogGray
	}
	return my.systemWriter
}

func (my *Classic) Console() console.IConsole {
	return my.console
}

func (my *Classic) SetVersion(f func() string) {
	my.f_version = f
}
func (my *Classic) GetVersion() string {
	ver_msg := ""
	if my.f_version != nil {
		ver_msg = my.f_version()
	}
	return ver_msg
}

func New(ctx *iContext, log_handler Logger) *Classic {
	cpuNum := runtime.NumCPU()
	goCnt := runtime.GOMAXPROCS(cpuNum)
	LogPurple("CPU[", cpuNum, "] :", goCnt)

	nRecovery := negroni.NewRecovery()

	var systemWriter func(v ...interface{})
	var logHandler negroni.Handler
	if log_handler == nil {
		defaultLogger := negroni.NewLogger()
		defaultLogger.ALogger = log.New(os.Stdout, Cat("[", ctx.Name(), "] "), 0)
		logHandler = defaultLogger
		systemWriter = func(v ...interface{}) {
			LogPurple(v...)
		}
	} else {
		logHandler = log_handler
		systemWriter = log_handler.Writeln
	}
	nStatic := negroni.NewStatic(http.Dir("public"))
	n := negroni.New(
		nRecovery,
		logHandler,
		nStatic,
	)

	router, renderer := newRouterRenderer(ctx.ServeFilePath())
	classic := &Classic{
		context:  ctx,
		n:        n,
		router:   router,
		renderer: renderer,
		callNav:  map[string]int{},
		call_log: NewTextAlignTapperKey(),

		systemWriter: systemWriter,

		console: console.New().
			SetLogFunc(LogWhite).
			SetErrFunc(LogError),
	}

	//console area -- start
	{
		con := classic.console
		con.AppendList(
			console.CMDList{
				{
					Name:  "version",
					HelpS: "classic version info",
					Action: func(ps []string) {
						runtime.GC()

						defer con.Atap()
						con.Atap()
						con.Log(" Brix classic web server infos.")
						con.Atap()
						con.Log(" os         :", runtime.GOOS)
						con.Log(" https      :", classic.isSSL)
						con.Log(" Port       :", ctx.Port())
						con.Log(" APICount   :", classic.apiCount)
						con.Log(" Goroutines :", runtime.NumGoroutine())
						if classic.f_version != nil {
							con.Atap()
							con.Log(classic.f_version())
						}
					},
				},
				{
					Name:  "classic",
					HelpS: "classic view sub infos..",
					CMDS: console.CMDList{
						{
							Name:  "apis",
							HelpS: "classic api count view",
							Action: func(ps []string) {
								con.Btap()
								classic.call_log.For("    ", func(i int, item, tap string) {
									con.Log(item, tap, classic.callCnt[i])
								})
								con.Atap()
								con.Log("total :", classic.call_log.Count())
								con.Btap()
							},
						},
					},
				},
			},
		)
	}
	//console area -- end

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing

	basic, exHeaders := cors.Robin(ctx.Headers())
	nhf := basic.NegroniHandlerFunc(exHeaders)
	classic.SetHandlerFunc(nhf)
	classic.SetHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if classic.isEnd {
			Fail(classic.ResponsWriter(w), serverUnderMaintenance)
			return
		}

		if index, do := classic.callNav[r.URL.Path]; do {
			atomic.AddInt32(&classic.callCnt[index], 1)
		}
		next(w, r)
	})

	return classic
}
func (my *Classic) setCall(pathName string) {
	my.callNav[pathName] = my.call_log.Count() - 1
	my.callCnt = append(my.callCnt, 0)

}

// func (my *Classic) endCall() {
// 	my.call_log.SetAlign("  -->")
// }

func (my *Classic) SetHandlers(routerlist HANDLERLIST, skipCount ...bool) {
	defer my.mu.Unlock()
	my.mu.Lock()
	if my.isStart {
		panic("classic.SetRouter is Set//")
	}

	if !IsTrue(skipCount) {
		my.apiCount += len(routerlist)
	}
	for _, r := range routerlist {
		my.router.Handle(
			r.Method(),
			r.Path(),
			r.Handle().getHandle(my.renderer),
		)
		my.call_log.SetKey(r.Path())
		my.setCall(r.Path())
	} //for
	my.items = append(my.items, routerlist...)
}

func (my *Classic) SetHandlerList(list ...HANDLER) {
	defer my.mu.Unlock()
	my.mu.Lock()
	for _, r := range list {
		my.router.Handle(
			r.Method(),
			r.Path(),
			r.Handle().getHandle(my.renderer),
		)

		my.call_log.SetKey(r.Path())
		my.setCall(r.Path())
		my.items = append(my.items, r)
	}
}

func (my *Classic) SetHandler(r HANDLER, skipCount ...bool) {
	defer my.mu.Unlock()
	my.mu.Lock()
	if my.isStart {
		panic("classic.SetHandler is Set//")
	}

	if !IsTrue(skipCount) {
		my.apiCount++
	}
	my.router.Handle(
		r.Method(),
		r.Path(),
		r.Handle().getHandle(my.renderer),
	)

	my.call_log.SetKey(r.Path())
	my.setCall(r.Path())
	my.items = append(my.items, r)
}

func (my *Classic) SetHandlerFunc(f HandlerFunc) {
	defer my.mu.Unlock()
	my.mu.Lock()
	if my.isStart {
		panic("classic.SetHanderFunc is Set//")
	}
	my.n.Use(f.getHandlerFunc())
}

func (my *Classic) startHandler(f func()) {
	defer my.mu.Unlock()
	my.mu.Lock()
	if !my.isStart {
		f()
		my.n.UseHandler(my.router)
	}
	//my.endCall()

	my.isStart = true
}
func (my *Classic) IsStart() bool {
	defer my.mu.RUnlock()
	my.mu.RLock()
	return my.isStart
}
func (my *Classic) WaitStart() {
	for {
		if !my.IsStart() {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		break
	}
}

func (my *Classic) StartDirect() {

	con := my.console
	con.AppendList(
		console.CMDList{
			console.CMDClearWnd(),
			//console.CMDMMS(),
			{
				Name:  "exit",
				HelpS: "classic server is close try",
				Action: func(ps []string) {
					if !my.isEnd {
						list := my.context.StarterList()
						for _, starter := range list {
							starter.Close()
						}
						con.Atap()
						con.Error("server process end")
						con.SetTitle("")
						con.Atap()
					}
					my.isEnd = true
				},
			},
		},
	)

	con.SetTitle(my.context.Name())

	is_console_skip := my.context.ConsoleSkip()
	if !is_console_skip {
		LogPurple("console cli mode on")
		go con.Start()
	} else {
		LogPurple("console cli mode off")
	}

	my.startHandler(func() {
		list := my.context.StarterList()
		for _, starter := range list {
			starter.Start()
		}

	})
	startListening(my)
}

func startListening(classic *Classic) {
	l := log.New(os.Stdout, "[classic] ", 0)
	LogYellow(classic)
	LogYellowBG("==============================================")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			msg := <-sig
			if !classic.isEnd {
				for i, starter := range classic.context.StarterList() {
					starter.Close()
					LogPurple("i", i, "]process end")
				}
			}
			LogError("sig :", msg)
			os.Exit(0)
		}
	}()

	server := &http.Server{
		Addr:    Cat(":", classic.context.Port()),
		Handler: classic.n,
	}
	server.SetKeepAlivesEnabled(classic.context.KeepAlive())
	if ssl_auth := classic.context.SSL(); len(ssl_auth) == 0 {
		LogYellowBG(" : brixcore.Classic Server                        ")
		LogYellowBG("----------------------------------------------")
		LogYellowBG("listening on :", classic.context.Port())
		LogYellowBG("==============================================")
		if classic.f_version != nil {
			LogYellow(classic.f_version())
		}
		LogYellow("server_start_time :", time.Now())

		l.Fatal(server.ListenAndServe())
	} else {
		classic.isSSL = true
		certFile := ssl_auth[0]
		keyFile := ssl_auth[1]
		LogYellowBG(" : brixcore.Classic SSL Server                    ")
		LogYellowBG("----------------------------------------------")
		LogYellowBG(" cert :", certFile)
		LogYellowBG(" key  :", keyFile)
		LogYellowBG("listening on :", classic.context.Port())
		LogYellowBG("==============================================")
		if classic.f_version != nil {
			LogYellow(classic.f_version())
		}
		LogYellow("server_start_time :", time.Now())

		l.Fatal(server.ListenAndServeTLS(certFile, keyFile))
	}

}
