package cnet

import (
	"sync"

	"txscheduler/brix/tools/dbg"

	"github.com/gorilla/websocket"
)

type nSession struct {
	*websocket.Conn
	writeMu sync.Mutex
}

func newSession(c *websocket.Conn) *nSession {
	return &nSession{
		Conn: c,
	}
}

//SendMessage : lock
func (my *nSession) SendMessage(ctype int, buf []byte) error {
	defer my.writeMu.Unlock()
	my.writeMu.Lock()
	return my.WriteMessage(ctype, buf)
}

//dSession : Dialer.Session
type dSession struct {
	*nSession
	isConnect bool
	mu        sync.RWMutex
}

func newSessionD(c *websocket.Conn) *dSession {
	return &dSession{
		nSession:  newSession(c),
		isConnect: true,
	}
}

//Close : [brix.cc] dSession Disconnect
func (my *dSession) Close(logs ...bool) {
	defer my.mu.Unlock()
	my.mu.Lock()
	if my.isConnect {
		my.isConnect = false
		my.nSession.Close()
		if dbg.IsTrue2(logs...) {
			dbg.PurpleItalic("Disconnect.")
		}
	}
}
func (my *dSession) IsConnect() bool {
	defer my.mu.RUnlock()
	my.mu.RLock()
	return my.isConnect
}
