package cnet

import (
	"jtools/jnet/chttp"
	"sync"

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

// SendMessage :
func (my *nSession) SendMessage(ctype int, buf []byte) error {
	defer my.writeMu.Unlock()
	my.writeMu.Lock()
	return my.WriteMessage(ctype, buf)
}

// dialerSession : Dialer.Session
type dialerSession struct {
	*nSession
	isConnect bool
	mu        sync.RWMutex
}

func newSessionD(c *websocket.Conn) *dialerSession {
	return &dialerSession{
		nSession:  newSession(c),
		isConnect: true,
	}
}

// Close :  dSession Disconnect
func (my *dialerSession) Close(logs ...bool) {
	defer my.mu.Unlock()
	my.mu.Lock()
	if my.isConnect {
		my.isConnect = false
		my.nSession.Close()
		if chttp.IsTrue(logs) {
			chttp.LogPurple("Disconnect.")
		}
	}
}
func (my *dialerSession) IsConnect() bool {
	defer my.mu.RUnlock()
	my.mu.RLock()
	return my.isConnect
}
