package jlog

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/urfave/negroni"
)

func (my *LogEntry) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := logNowTime()

	next(rw, r)

	res := rw.(negroni.ResponseWriter)

	client_ip, _ := getIP(r)
	logMessage := cat(
		res.Status(),
		" | ", time.Since(start),
		" | ", client_ip,
		" | ", r.Host,
		" | ", r.Method,
		" | ", r.URL.Path,
		//" | ", r,
	)
	Debug(logMessage)
}

func (my *LogEntry) Writeln(v ...interface{}) {
	Info(v...)
}

func getIP(r *http.Request) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("[chttp.check_ip]No valid ip found")
}

func cat(a ...interface{}) string {
	sl := []string{}
	for _, v := range a {
		sl = append(sl, fmt.Sprint(v))
	}
	return strings.Join(sl, "")
}
