package chttp

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/urfave/negroni"
)

type Logger interface {
	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	Writeln(v ...interface{})
}

type dayLogger struct {
	folderPath string
	prefixName string
	tagName    string
	nDay       int
	fd         *os.File
	l          *log.Logger
	mu         sync.Mutex
}

func NewDayLogger(folder, prefix string, tagName string) Logger {
	my := &dayLogger{
		folderPath: folder,
		prefixName: prefix,
		tagName:    tagName,
		nDay:       time.Now().Day(),
	}
	fd, err := my.getFile()
	if err != nil {
		os.MkdirAll(my.folderPath, os.ModePerm)
		fd, _ = my.getFile()
		//panic(err)
	}
	my.fd = fd
	my.l = log.New(my.fd, "["+my.tagName+"]", 0)
	return my
}

func (my *dayLogger) getFile() (*os.File, error) {
	ymdstring := YMDString(NowTime())
	return os.OpenFile(
		my.folderPath+"/"+my.prefixName+"_"+ymdstring+".txt",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
}

func (my *dayLogger) makeNextDayFile() error {

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

func (my *dayLogger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := NowTime()

	next(rw, r)

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
	logMessage := Cat(
		start.Format("2006-01-02 15:04:05.999"),
		" | ", res.Status(),
		" | ", time.Since(start),
		" | ", client_ip,
		" | ", r.Host,
		" | ", r.Method,
		" | ", r.URL.Path,
		//" | ", r,
	)
	my.l.Println(logMessage)
	fmt.Println(logMessage)
}

func (my *dayLogger) Writeln(v ...interface{}) {
	start := NowTime()
	_ = start
	if my.nDay != start.Day() {
		my.mu.Lock()
		{
			my.makeNextDayFile()
			my.nDay = start.Day()
		}
		my.mu.Unlock()
	}
	my.l.Println(v...)
}
