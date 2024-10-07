package cnet

import (
	"fmt"

	"github.com/gorilla/websocket"
)

//Socket : Web-Socket
type Socket struct {
	id int64
	*nSession
}

func (my Socket) ID() int64 { return my.id }

func newServerSocket(c *websocket.Conn) *Socket {
	return &Socket{
		id:       makeSocketID(),
		nSession: newSession(c),
	}
}

//ToString :
func (my *Socket) ToString() string {
	if my.Conn == nil {
		return "nil"
	}
	return fmt.Sprintf("[%v]%v", my.id, my.RemoteAddr().String())
}
