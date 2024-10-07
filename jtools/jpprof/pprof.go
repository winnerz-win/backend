package jpprof

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
)

var (
	is_profile_mode = false
	run_port        = 60000
	_println        func(a ...interface{})
)

func Start(port ...int) {
	_start(
		func(a ...interface{}) {
			fmt.Println(a...)
		},
		port...,
	)
}

func Start2(println func(a ...interface{}), port ...int) {
	_start(println, port...)
}

func _start(println func(a ...interface{}), port ...int) {
	go func() {
		_println = println
		if len(port) > 0 && port[0] > 0 && port[0] < 65535 {
			run_port = port[0]
		}
		println()
		println("############################################")
		println("  jpprof.Start [ port:", run_port, "]")
		println("############################################")

		is_profile_mode = true
		err := http.ListenAndServe(
			fmt.Sprintf(":%v", run_port),
			nil,
		)
		if err != nil {
			is_profile_mode = false
			log.Println(err)
		}
	}()
}

func View() {
	if !is_profile_mode {
		return
	}
	_println()
	_println("############################################")
	_println(" jpprof.View                                ")
	_println(" http://localhost:" + fmt.Sprint(run_port) + "/debug/pprof")
	_println("############################################")
}
