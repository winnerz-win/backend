package jnet

import (
	"net"
	"net/url"
	"time"
)

//Ping : 3 sec
func Ping(address string, timeoutsec ...int) bool {
	path, err := url.Parse(address)
	if err != nil {
		//dbg.Red("jnet.Ping : ", err)
		return false
	}
	timeout := 3
	if len(timeoutsec) > 0 && timeoutsec[0] > 0 {
		timeout = timeoutsec[0]
	}
	con, err := net.DialTimeout("tcp", path.Host, time.Second*time.Duration(timeout))
	if err != nil {
		//dbg.Red(err)
		return false
	}
	_ = con.Close()
	return true
}
