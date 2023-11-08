package ebcm

import "strings"

func (my Sender) Coins(address string) string {
	wei := my.Balance(address)
	return WeiToETH(wei)
}

func GetHostURL(host string, key ...string) string {
	connURL := host
	if len(key) > 0 {
		if key[0] != "" {
			if strings.HasSuffix(connURL, "/v3") || strings.HasSuffix(connURL, "/v3/") {
				if !strings.HasSuffix(connURL, "/") {
					connURL += "/"
				}
				connURL += strings.TrimSpace(key[0])
			}
		}
	}
	return connURL
}
