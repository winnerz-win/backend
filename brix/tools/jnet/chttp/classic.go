package chttp

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"sync"
	"syscall"
	"time"

	"txscheduler/brix/tools"
	"txscheduler/brix/tools/console"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp/cors"
	"txscheduler/brix/tools/jpath"
	"txscheduler/brix/tools/runtext"

	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

const (
	_RWTimeoutDuration = time.Duration(10) * time.Second
)

// DumbLogger : negroni log dumb
type DumbLogger struct{}

// Println :
func (my DumbLogger) Println(v ...interface{}) {
	//fmt.Println(v...)
}

// Printf :
func (my DumbLogger) Printf(format string, v ...interface{}) {
	//fmt.Println(format, v)
}

type pathCntData struct {
	count int64
	duSum time.Duration
}

const (
	SystemWriteKey = "classic_system_writer"
)

type SystemWriter func(v ...interface{})

func SystemWriteln(r *http.Request, v ...interface{}) {
	if len(v) == 0 {
		return
	}
	if val := r.Context().Value(SystemWriteKey); val != nil {
		if writeln, do := val.(SystemWriter); do {
			writeln(v...)
		}
	}
}

// Classic :
type Classic struct {
	n              *negroni.Negroni
	router         *httprouter.Router
	port           string
	apiCount       int
	apiList        []string
	isConsoleStart bool
	startC         chan struct{}
	isExit         bool

	isAPISet      bool
	isServerStart bool
	isSSL         bool
	mu            sync.RWMutex

	//------------------------------
	toinfo []func() interface{}

	testAPIS map[string]PContext

	pathCount map[string]pathCntData
	pathMu    sync.RWMutex

	isKeepAlive bool

	isLogViewSkip bool

	starterList runtext.StarterList

	systemWriter SystemWriter
}

func (my *Classic) IsLogView() bool { return !my.isLogViewSkip }

func (my *Classic) SystemWriter() SystemWriter {
	return my.systemWriter
}

// SetKeepAlive :
func (my *Classic) SetKeepAlive(keep bool) {
	my.isKeepAlive = keep
}

// TestAPIHandle :
func (my *Classic) TestAPIHandle(api string, w http.ResponseWriter, req *http.Request, ps Params) error {
	defer my.mu.RUnlock()
	my.mu.RLock()
	if _, do := my.testAPIS[api]; do {
		my.testAPIS[api].Handle(w, req, ps)
		return nil
	}
	return fmt.Errorf("not found handle [ %v ]", api)
}

// AddInfo :
func (my *Classic) AddInfo(f func() interface{}) {
	defer my.mu.Unlock()
	my.mu.Lock()
	my.toinfo = append(my.toinfo, f)
}

// ToString :
func (my *Classic) ToString() string {
	msg := ""
	for _, v := range my.toinfo {
		msg += fmt.Sprintf("%v\n", v())
	}
	return msg
}

// NewRouter :
func NewRouter(sfp ...string) *httprouter.Router {
	serveFilePath := ServeFilesPath
	if len(sfp) > 0 {
		serveFilePath = sfp[0]
	}

	renderPath := jpath.NowPath() + "\\" + RenderRootPath
	Renderer = render.New(render.Options{
		Directory:  renderPath,
		Extensions: []string{".tmpl", ".html"},
	})

	router := httprouter.New()
	//KKJJSS
	staticDir := jpath.NowPath() + "\\" + serveFilePath

	dbg.PurpleBold("router.filepath :", staticDir)
	dbg.PurpleBold("renderer.filepath :", renderPath)

	//router.ServeFiles(fmt.Sprintf("/%v/*filepath", serveFilePath), http.Dir(serveFilePath))
	router.ServeFiles(fmt.Sprintf("/%v/*filepath", serveFilePath), http.Dir(staticDir))
	return router
}

// ServeFiles :
func (my *Classic) ServeFiles(path string, staticPath string) {
	router := my.router
	staticDir := jpath.NowPath() + "\\" + staticPath

	router.ServeFiles(fmt.Sprintf("/%v/*filepath", path), http.Dir(staticDir))

}

// NewClassic :
func NewClassic(port int, corsHeader []string, tag ...string) *Classic {
	taglist := []string{}
	for _, v := range tag {
		taglist = append(taglist, v)
	}
	return newClassic(port, corsHeader, taglist, true, nil)
}
func NewClassicLogger(port int, corsHeader []string, tag string, logger Logger) *Classic {
	return newClassic(port, corsHeader, []string{tag}, true, logger)
}

type ISimple interface {
	StartSimple()
	StartSimpleTLS(certFile, keyFile string)
	SetContextHandles(contexts []PContext, countSkip ...bool)
	SetHandlerFunc(f HandlerFunc)
	WaitC()
	Exit()
	IsExit() bool
	IsServerRun() bool
	APICount() int
}

func NewSimple(port int, corsHeader []string, tag ...string) ISimple {
	taglist := []string{}
	for _, v := range tag {
		taglist = append(taglist, v)
	}
	return newClassic(port, corsHeader, taglist, false, nil)
}

func newClassic(
	port int,
	corsHeader []string,
	tag []string,
	isCmd bool,
	logger Logger,
) *Classic {
	dbg.PurpleBold("######### negroni.NewClassic #########")

	const acao = "Access-Control-Allow-Origin"
	if corsHeader == nil {
		corsHeader = []string{acao}
	} else {
		isDo := false
		for _, v := range corsHeader {
			if v == acao {
				isDo = true
				break
			}
		}
		if !isDo {
			corsHeader = append(corsHeader, acao)
		}
	}

	cpuNum := runtime.NumCPU()
	goCnt := runtime.GOMAXPROCS(cpuNum)
	dbg.PurpleBold("CPU :", cpuNum, ", ", goCnt)

	var systemWriter func(v ...interface{}) = nil
	var logHandler negroni.Handler
	var consoleLogger *negroni.Logger = nil
	cmdLogOn := console.Ccmd{}
	cmdLogOff := console.Ccmd{}
	cmdClassicLog := console.Ccmd{}

	isConsoleLog := logger == nil
	if !isConsoleLog {
		logHandler = logger
		systemWriter = logger.Writeln
	} else {
		nLogger := negroni.NewLogger()
		consoleLogger = nLogger
		logHandler = consoleLogger
	}

	nRecovery := negroni.NewRecovery()
	nStatic := negroni.NewStatic(http.Dir("public"))
	n := negroni.New(
		nRecovery,
		logHandler,
		nStatic,
	)
	classic := &Classic{
		n:              n,
		router:         NewRouter(),
		port:           fmt.Sprintf("%v", port),
		isConsoleStart: false,
		startC:         make(chan struct{}, 1),
		isAPISet:       false,
		isServerStart:  false,
		testAPIS:       map[string]PContext{},
		pathCount:      map[string]pathCntData{},
		systemWriter:   systemWriter,
	}

	if isConsoleLog {
		dumbLoger := DumbLogger{}

		logPrefix := "[robin] "
		if len(tag) > 0 {
			logPrefix = fmt.Sprintf("[%v] ", tag[0])
		}
		nLoger := log.New(os.Stdout, logPrefix, 0)
		consoleLogger.ALogger = nLoger

		cmdLogOn = console.Ccmd{
			Cmd:        "on",
			HeaderFunc: func() string { return "server console log on" },
			NoParams:   true,
			Help:       "ConsoleLog On.",
			Work: func(done chan<- bool, ps []string) {
				defer func() {
					done <- true
				}()
				classic.isLogViewSkip = false
				consoleLogger.ALogger = nLoger
				dbg.SetSkipLog(false)
				console.Println("Negroni Console Log On.")
			},
		}
		cmdLogOff = console.Ccmd{
			Cmd:        "off",
			HeaderFunc: func() string { return "server console log off" },
			NoParams:   true,
			Help:       "ConsoleLog Off.",
			Work: func(done chan<- bool, ps []string) {
				defer func() {
					done <- true
				}()
				classic.isLogViewSkip = true
				consoleLogger.ALogger = dumbLoger
				dbg.SetSkipLog(true)
				console.Println("Negroni Console Log Off.")
			},
		}

		classicLogToggle := true
		cmdClassicLog = console.Ccmd{
			Cmd:        "classiclog",
			HeaderFunc: func() string { return "classic console log on/off toggle" },
			NoParams:   true,
			Help:       "ConsoleLog Off.",
			Work: func(done chan<- bool, ps []string) {
				defer func() {
					done <- true
				}()

				classicLogToggle = !classicLogToggle
				classic.isLogViewSkip = classicLogToggle
				if classicLogToggle {
					consoleLogger.ALogger = dumbLoger
				} else {
					consoleLogger.ALogger = nLoger
				}
				console.Println("Negroni Classic loger :", !classicLogToggle)
			},
		}
	}

	if isCmd {
		cmds := console.Commands{
			console.ClearConsole(),
			cmdLogOn,
			cmdLogOff,
			cmdClassicLog,
			{
				Cmd:        "version",
				HeaderFunc: func() string { return "view server framework version" },
				NoParams:   true,
				Help:       "show server info",
				Work: func(done chan<- bool, ps []string) {
					defer console.DoneC(done)
					console.Log(classic.ToString())
				},
			},
			{
				Cmd:        "classic.apis",
				HeaderFunc: func() string { return "view http apis" },
				NoParams:   true,
				Help:       "view Server-API list.",
				Work: func(done chan<- bool, ps []string) {
					defer console.DoneC(done)

					console.Atap()
					apiList := classic.APIList()
					for _, v := range apiList {
						console.Log(v)
					} //for
					console.Atap()
					console.Log("Total :", len(apiList))
					console.Atap()

				},
			},
			{
				Cmd:        "classic.cnt",
				HeaderFunc: func() string { return "view http api-call count" },
				NoParams:   true,
				Help:       "view Server-API Call count list.",
				Work: func(done chan<- bool, ps []string) {
					defer console.DoneC(done)

					console.Atap()
					apiList := classic.PathCountView()
					for _, v := range apiList {
						console.Log(v)
					} //for
					console.Atap()
					console.Log("Total :", len(apiList))
					console.Atap()

				},
			},
		}
		console.SetCmd(cmds)
	}

	classic.AddInfo(func() interface{} {
		runtime.GC()
		msg := tools.Version()
		if classic.isSSL {
			msg += fmt.Sprintf("### HTTPS-SERVER ###\n")
		}
		msg += fmt.Sprintf("Port       : %v\n", classic.port)
		msg += fmt.Sprintf("APICount   : %v\n", classic.apiCount)
		msg += fmt.Sprintf("Goroutines : %v\n", runtime.NumGoroutine())
		msg += tools.BAR + tools.ENTER
		return msg
	})

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	basic, exHeaders := cors.Robin(corsHeader)
	nhf := basic.NegroniHandlerFunc(exHeaders)
	classic.SetHandlerFunc(nhf)

	classic.SetHandlerFunc(
		func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			if classic.IsExit() {
				if isResponseFormat == false {
					v := JsonType{"code": ErrorServerUnderMaintenance}
					Renderer.JSON(w, http.StatusBadRequest, v)
				} else {
					Fail(w, serverUnderMaintenance)
				}

				return
			}

			start := time.Now()

			next(w, r)

			du := time.Since(start)

			classic.incPathCount(r.URL.Path, du)
		},
	)

	return classic
}

func (my *Classic) incPathCount(path string, du time.Duration) {
	defer my.pathMu.Unlock()
	my.pathMu.Lock()

	if data, do := my.pathCount[path]; do {
		data.count = 1
		data.duSum += du
		my.pathCount[path] = data
	} else {

		my.pathCount[path] = pathCntData{
			count: 1,
			duSum: du,
		}
	}
}

type callPathCount struct {
	path string
	cnt  int64
	avg  time.Duration
}
type callPathCountList []callPathCount

func (my callPathCountList) toList() []string {
	list := []string{}
	sort.Slice(my, func(i, j int) bool { return my[i].cnt < my[j].cnt })

	tp := dbg.NewTextAlignTap()

	for _, c := range my {
		tp.AddText(c.path)
	} //for

	tp.For(func(i int, item, tap string) {
		data := my[i]
		list = append(
			list,
			dbg.Cat(item, tap, data.cnt, "  ", data.avg),
		)
	})

	return list
}

// PathCountView :
func (my *Classic) PathCountView() []string {
	list := callPathCountList{}
	my.pathMu.RLock()
	for path, data := range my.pathCount {
		item := callPathCount{
			path: path,
			cnt:  data.count,
			avg:  data.duSum / time.Duration(data.count),
		}
		list = append(list, item)
	} //for
	my.pathMu.RUnlock()
	return list.toList()
}

// SetServeFiles :
func (my *Classic) SetServeFiles(filePath string) {
	my.router.ServeFiles(fmt.Sprintf("/%v/*filepath", filePath), http.Dir(filePath))
}

// SetHandlerFunc :
func (my *Classic) SetHandlerFunc(f HandlerFunc) {
	defer my.mu.Unlock()
	my.mu.Lock()
	if my.isAPISet == true {
		panic("SetHandlerFunc:: after isSet")
	}
	my.n.Use(f.getHandlerFunc())
}

func (my *Classic) SetHandler(method, path string, handle RouterHandle) {
	defer my.mu.Unlock()
	my.mu.Lock()
	ctx := NewContext(method, path, handle)
	my.router.Handle(ctx.Method, ctx.Path, handle.getHandle())
	my.apiList = append(my.apiList, (*Context)(ctx).String())
	my.testAPIS[ctx.Path] = ctx
}

// SetContextHandles :
func (my *Classic) SetContextHandles(contexts []PContext, countSkip ...bool) {
	defer my.mu.Unlock()
	my.mu.Lock()
	if my.isAPISet == true {
		panic("SetContextHandles:: after isSet")
	}

	if len(countSkip) == 0 {
		my.apiCount += len(contexts)
	}
	for _, ctx := range contexts {
		my.router.Handle(ctx.Method, ctx.Path, ctx.Handle.getHandle())
		my.apiList = append(my.apiList, (*Context)(ctx).String())

		my.testAPIS[ctx.Path] = ctx
	} //for
}

// APIList :
func (my *Classic) APIList() []string {
	return my.apiList
}

// SetRouting : 핸들러를 모두 등록 한 뒤 마지막에 라우팅을 연결한다.
func (my *Classic) SetRouting() {
	defer my.mu.Unlock()
	my.mu.Lock()

	my.n.UseHandler(my.router)

	my.isAPISet = true
}

// StartSimple :
func (my *Classic) StartSimple() {
	dbg.Green("[ StartSimple ]", my.ToString())
	portString := fmt.Sprintf(":%v", my.port)

	my.SetRouting()

	//my.n.Run(portString)
	startListening(portString, my)
}

// ConsoleStart : lock-func
func (my *Classic) ConsoleStart(isConsoleReadMode ...bool) {
	defer my.mu.Unlock()
	my.mu.Lock()
	if my.isConsoleStart == false {
		if len(isConsoleReadMode) > 0 && isConsoleReadMode[0] == true {
			console.ReadMode()
		} else {
			console.Start()
		}
	}
	my.isConsoleStart = true

	fmt.Println(my.ToString())
}

// ServerStart : isConsoleReadMode - true 일경우 콜솔에서 서버 제어
func (my *Classic) ServerStart(isConsoleReadMode ...bool) {

	my.ConsoleStart(isConsoleReadMode...)

	my.mu.Lock()
	{
		console.ServerStart()
		if my.isServerStart == false {
			close(my.startC)
		}
		my.isServerStart = true
	}
	my.mu.Unlock()

	portString := fmt.Sprintf(":%v", my.port)
	//my.n.Run(portString)
	startListening(portString, my)
}

func startListening(portString string, classic *Classic) {
	l := log.New(os.Stdout, "[brix.cc] ", 0)
	dbg.YellowItalicBG("############################################")
	dbg.YellowItalicBG("txscheduler/brix.CC Server Listening on", portString)
	dbg.YellowItalicBG("############################################")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			msg := <-sig
			if !classic.IsExit() {
				classic.Exit()
				for i, starter := range classic.starterList {
					starter.Close()
					dbg.Purple("i", i, "]process end")
				}
			}
			dbg.RedItalic("sig :", msg)
			os.Exit(0)
		}
	}()

	server := newHttpServer(
		portString,
		classic.n,
	)
	server.SetKeepAlivesEnabled(classic.isKeepAlive)
	l.Fatal(server.ListenAndServe())
	//l.Fatal(http.ListenAndServe(portString, classic.n))
}

func (my *Classic) prepareStart(starterList runtext.StarterList, startRun func()) {

	my.starterList = starterList
	console.SetCmd(console.Commands{
		{
			Cmd:        "exit",
			HeaderFunc: func() string { return "server process exit" },
			NoParams:   true,
			Work: func(done chan<- bool, ps []string) {
				defer console.DoneC(done)

				if !my.IsExit() {
					my.Exit()
					for _, starter := range starterList {
						starter.Close()
					} //for
					console.Atap()
					console.Log("server process exit")
					console.Atap()
				}

			},
		},
		{
			Cmd:        "start",
			HeaderFunc: func() string { return "server start" },
			NoParams:   true,
			Work: func(done chan<- bool, ps []string) {
				defer console.DoneC(done)
				if my.IsServerRun() {
					return
				}

				go func() {
					my.WaitC()
					for _, starter := range starterList {
						starter.Start()
					} //for
				}()
				my.SetRouting()

				if my.isServerStart == false {
					close(my.startC)
				}

				my.isServerStart = true
				go startRun()

			},
		},
	})
}

// CmdStart :
func (my *Classic) CmdStart(starterList runtext.StarterList, isConsoleReadMode ...bool) {
	my.prepareStart(
		starterList,
		func() {
			console.ServerStart()
			portString := fmt.Sprintf(":%v", my.port)
			//my.n.Run(portString)
			startListening(portString, my)
		},
	)

	if my.isConsoleStart == false {
		if len(isConsoleReadMode) > 0 && isConsoleReadMode[0] == true {
			console.ReadMode()
		} else {
			console.Start()
		}
	}
	my.isConsoleStart = true
	fmt.Println(my.ToString())

	waitC := make(chan struct{}, 1)
	<-waitC
}

/*
CmdStartTLS :
netsh interface portproxy add v4tov4 listenport=443 listenaddress=172.31.92.82 connectport=80 connectaddress=172.31.92.82 protocol=tcp
netsh interface portproxy show v4tov4
netsh interface portproxy delete v4tov4 listenport=443 listenaddress=172.31.92.82
*/
func (my *Classic) CmdStartTLS(certFile, keyFile string, starterList runtext.StarterList, isConsoleReadMode ...bool) {
	my.isSSL = true
	my.prepareStart(
		starterList,
		func() {
			console.ServerStart()
			portString := fmt.Sprintf(":%v", my.port)
			runTLS(my, certFile, keyFile, portString)
		},
	)

	if my.isConsoleStart == false {
		if len(isConsoleReadMode) > 0 && isConsoleReadMode[0] == true {
			console.ReadMode()
		} else {
			console.Start()
		}
	}
	my.isConsoleStart = true
	fmt.Println(my.ToString())

	waitC := make(chan struct{}, 1)
	<-waitC
}

// ServerStartTLS :
func (my *Classic) ServerStartTLS(certFile, keyFile string, isConsoleReadMode ...bool) {
	my.isSSL = true
	my.ConsoleStart(isConsoleReadMode...)
	my.mu.Lock()
	{
		console.ServerStart()
		if my.isServerStart == false {
			close(my.startC)
		}
		my.isServerStart = true
	}
	my.mu.Unlock()

	portString := fmt.Sprintf(":%v", my.port)
	runTLS(my, certFile, keyFile, portString)
}

// StartSimpleTLS :
func (my *Classic) StartSimpleTLS(certFile, keyFile string) {
	dbg.Green("[ StartSimpleTLS ]", my.ToString())
	portString := fmt.Sprintf(":%v", my.port)

	my.SetRouting()

	runTLS(my, certFile, keyFile, portString)
}

func runTLS(classic *Classic, certFile, keyFile string, addr ...string) {
	classic.isSSL = true
	detectAddress := func(addr ...string) string {
		if len(addr) > 0 {
			return addr[0]
		}
		if port := os.Getenv("PORT"); port != "" {
			return ":" + port
		}
		return negroni.DefaultAddress
	}
	_ = detectAddress

	l := log.New(os.Stdout, "[brix.cc-TLS] ", 0)
	dbg.YellowItalicBG("======================================================")
	dbg.YellowItalicBG(" : ON HTTPS SERVER                                    ")
	dbg.YellowItalicBG("======================================================")
	finalAddr := detectAddress(addr...)
	dbg.YellowItalicBG("listening on :", finalAddr)

	server := newHttpServer(
		finalAddr,
		classic.n,
	)

	server.SetKeepAlivesEnabled(classic.isKeepAlive)
	err := server.ListenAndServeTLS(certFile, keyFile)
	//err := http.ListenAndServeTLS(finalAddr, certFile, keyFile, classic.n)
	l.Fatal(err)
}

func newHttpServer(finalAddr string, handler http.Handler) *http.Server {
	/*
		https://syntaxsugar.tistory.com/entry/GoGolang-HTTP-%EC%84%B1%EB%8A%A5-%ED%8A%9C%EB%8B%9D
	*/
	server := &http.Server{
		Addr:    finalAddr,
		Handler: handler,
		// ReadTimeout:  RWTimeoutDuration,
		// WriteTimeout: RWTimeoutDuration,
	}
	return server
}

// WaitC :
func (my *Classic) WaitC() {
	<-my.startC
}

// Exit :
func (my *Classic) Exit() {
	defer my.mu.Unlock()
	my.mu.Lock()
	my.isExit = true
}

// IsExit :
func (my *Classic) IsExit() bool {
	defer my.mu.RUnlock()
	my.mu.RLock()
	return my.isExit
}

// IsServerRun :
func (my *Classic) IsServerRun() bool {
	defer my.mu.RUnlock()
	my.mu.RLock()
	return my.isServerStart
}

// APICount :
func (my *Classic) APICount() int {
	defer my.mu.Unlock()
	my.mu.Lock()
	return my.apiCount
}
