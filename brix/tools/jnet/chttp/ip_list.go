package chttp

import (
	"net"
	"net/http"
	"strings"
)

//IPList :
func IPList() []string {
	iplist := []string{}

	ifaces, err := net.Interfaces()
	if err != nil {
		return iplist
	}
	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// process IP address
			iplist = append(iplist, ip.String())
		}
	}
	return iplist
}

//IPDO :
func IPDO(ip string) bool {
	iplist := IPList()
	for _, v := range iplist {
		if v == ip {
			return true
		}
	}
	return false
}

func RemoteIPPort(req *http.Request) []string {
	if strings.HasPrefix(req.RemoteAddr, "[::1]") {
		ss := []string{
			"localhost",
			strings.ReplaceAll(req.RemoteAddr, "[::1]:", ""),
		}
		return ss
	}
	ss := strings.Split(req.RemoteAddr, ":")
	return ss
}

func IsLocalHost(ip string) bool {
	return true
}
