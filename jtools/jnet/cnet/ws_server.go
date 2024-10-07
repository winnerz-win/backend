package cnet

import (
	"net/http"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

const (
	defaultSocketBufferSize = 1024
)

var (
	socketNetID int64 = 1000
)

func makeSocketID() int64 {
	return atomic.AddInt64(&socketNetID, 1)
}

//SocketOption :
type SocketOption struct {
	W           http.ResponseWriter
	Request     *http.Request
	BufferSize  int
	CheckOrigin func(req *http.Request) bool
}

//SocketUpgradeOpt :
func SocketUpgradeOpt(opt *SocketOption) (*Socket, error) {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  opt.BufferSize,
		WriteBufferSize: opt.BufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	if opt.CheckOrigin != nil {
		upgrader.CheckOrigin = opt.CheckOrigin
	}
	conn, err := upgrader.Upgrade(opt.W, opt.Request, nil)
	if err != nil {
		return nil, err
	}

	client := newServerSocket(conn)

	return client, nil
}

//SocketUpgrade :
func SocketUpgrade(w http.ResponseWriter, req *http.Request,
	cHandle ...func(req *http.Request) bool,
) (*Socket, error) {
	opt := &SocketOption{
		W:          w,
		Request:    req,
		BufferSize: defaultSocketBufferSize,
	}
	if len(cHandle) > 0 {
		opt.CheckOrigin = cHandle[0]
	}
	return SocketUpgradeOpt(opt)
}
