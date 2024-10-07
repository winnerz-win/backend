package cnet

import (
	"jtools/jnet/chttp"
	"net/http"
)

type Header struct {
	header http.Header
}

func (my Header) String() string {
	return chttp.ToJsonString(my.header)
}
func (my Header) Get(key string) string {
	return my.header.Get(key)
}
func (my Header) ContentType() string {
	return my.Get("Content-Type")
}
func (my Header) ContentLength() string {
	return my.Get("Content-Length")
}
func (my Header) IpfsPath() string {
	return my.Get("X-Ipfs-Path")
}
func (my Header) OriginHeader() http.Header {
	return my.header
}
