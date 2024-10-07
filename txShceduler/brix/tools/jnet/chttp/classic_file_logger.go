package chttp

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/mms"

	"github.com/urfave/negroni"
)

type Logger interface {
	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	Writeln(v ...interface{})
}

type dayLoger struct {
	folderPath string
	prefixName string
	tagName    string
	nDay       int
	fd         *os.File
	l          *log.Logger
	mu         sync.Mutex
	isSkipLog  bool
}

func NewDayLogger(folder, prefix string, tagName string, isSkipLog ...bool) Logger {
	my := &dayLoger{
		folderPath: folder,
		prefixName: prefix,
		tagName:    tagName,
		nDay:       time.Now().Day(),
	}
	if len(isSkipLog) > 0 && isSkipLog[0] {
		my.isSkipLog = true
	}

	fd, err := my.getFile()
	if err != nil {
		os.MkdirAll(my.folderPath, os.ModePerm)
		//panic(err)
	}
	my.fd = fd
	my.l = log.New(my.fd, "["+my.tagName+"]", 0)
	return my
}
func (my *dayLoger) SetSkipLog(f bool) {
	my.isSkipLog = f
}

func (my *dayLoger) getFile() (*os.File, error) {
	ymdstring := mms.Now().YMDString()
	return os.OpenFile(
		my.folderPath+"/"+my.prefixName+"_"+ymdstring+".txt",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
}

func (my *dayLoger) makeNextDayFile() error {

	fd, err := my.getFile()
	if err != nil {
		return err
	}
	if my.fd != nil {
		my.fd.Close()
	}

	my.l.SetOutput(fd)
	my.fd = fd

	return nil
}

func (my *dayLoger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	next(rw, r)

	if !my.isSkipLog {

		if my.nDay != start.Day() {
			my.mu.Lock()
			{
				my.makeNextDayFile()
				my.nDay = start.Day()
			}
			my.mu.Unlock()
		}

		res := rw.(negroni.ResponseWriter)

		client_ip, _ := GetIP(r)
		logMessage := dbg.Cat(
			mms.FromTime(start).String2(),
			" [ ", client_ip, " / ", r.Host, " ] ",
			r.Method, " [ ", res.Status(), " ] ",
			"( ", time.Since(start), " ) ",
			r.URL.Path,
			//" | ", r,
		)
		my.l.Println(logMessage)
		fmt.Println(logMessage)

	}

}

func (my *dayLoger) Writeln(v ...interface{}) {
	start := time.Now()
	_ = start
	if my.nDay != start.Day() {
		my.mu.Lock()
		{
			my.makeNextDayFile()
			my.nDay = start.Day()
		}
		my.mu.Unlock()
	}

	sl := []interface{}{}
	sl = append(sl, mms.FromTime(start).String2(), " ")
	sl = append(sl, v...)

	my.l.Println(sl...)
}
