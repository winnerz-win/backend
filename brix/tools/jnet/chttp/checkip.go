package chttp

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func GetIP(r *http.Request) (string, error) {
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

func ParseIP(ipstring string) (string, error) {
	//Get IP from the X-REAL-IP header
	netIP := net.ParseIP(ipstring)
	if netIP != nil {
		return ipstring, nil
	}

	//Get IP from X-FORWARDED-FOR header
	splitIps := strings.Split(ipstring, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(ipstring)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("[chttp.check_ip]No valid ip found")
}

func IsLocalhost(r *http.Request, ip string) bool {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return false
	}

	host := strings.Split(r.Host, ":")[0]
	if host == ip {
		return true
	}
	if ip == "::1" {
		return true
	}
	return false
}
