package chttp

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func _Cat(a ...interface{}) string {
	sl := []string{}
	for _, v := range a {
		sl = append(sl, fmt.Sprintf("%v", v))
	}
	return strings.Join(sl, "")
}

func IsEtagBuf(w ResponseWriter, req *http.Request, buf []byte, maxAge interface{}) bool {
	etag := fmt.Sprintf("%x", md5.Sum(buf))

	w.Header().Set("Etag", etag)
	w.Header().Set("Cache-Control", _Cat("private, max-age=", maxAge))
	w.Header().Set("Content-Length", _Cat(len(buf)))

	if match := req.Header.Get("If-None-Match"); match != "" {
		if match == etag {
			w.W().WriteHeader(http.StatusNotModified)
			return true
		}
	}

	return false
}

func IsEtagFile(w ResponseWriter, req *http.Request, file_name string, maxAge interface{}) bool {
	file, err := os.Stat(file_name)
	if err != nil {
		return false
	}
	modified_time := file.ModTime()

	etag := fmt.Sprintf("%x", md5.Sum([]byte(modified_time.String())))
	w.Header().Set("Etag", etag)
	w.Header().Set("Cache-Control", _Cat("private, max-age=", maxAge))
	w.Header().Set("Content-Length", _Cat(file.Size()))

	if match := req.Header.Get("If-None-Match"); match != "" {
		if match == etag {
			w.W().WriteHeader(http.StatusNotModified)
			return true
		}
	}

	return false
}
