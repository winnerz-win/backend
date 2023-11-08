package tools

import (
	"fmt"
	"log"
	"net/http"

	_ "net/http/pprof"

	"txscheduler/brix/tools/dbg"
)

var (
	defaultPort   = 60000
	isProfileMode = false
)

//profileLog :
func profileLog() {
	if isProfileMode {
		fmt.Println()
		dbg.RedItalicBG("############################################")
		dbg.RedItalicBG("  Brix.Go Profile Mode                      ")
		dbg.RedItalicBG("  url - http://localhost:" + fmt.Sprint(defaultPort) + "/debug/pprof")
		dbg.RedItalicBG("############################################")
	}
}

//StartProfile :  http://localhost:60000/debug/pprof/
func StartProfile(port ...int) {
	isProfileMode = true
	go func() {
		defaultPort = 60000
		if len(port) > 0 {
			defaultPort = port[0]
		}
		fmt.Println()
		dbg.RedItalicBG("############################################")
		dbg.RedItalicBG("           chttp.StartProfile               ")
		dbg.RedItalicBG("############################################")
		log.Println(http.ListenAndServe(fmt.Sprintf(":%v", defaultPort), nil))

	}()
}
