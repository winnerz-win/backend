package cnet

import (
	"fmt"

	"github.com/gorilla/websocket"
)

//Socket : Web-Socket
type Socket struct {
	//ID : Socket Session ID
	ID int32 //Socket Session ID
	*nSession
	// *websocket.Conn
	// writeMu sync.Mutex
}

func newServerSocket(c *websocket.Conn) *Socket {
	return &Socket{
		ID:       makeSocketID(),
		nSession: newSession(c),
	}
}

//ToString :
func (my *Socket) ToString() string {
	if my.Conn == nil {
		return "nil"
	}
	return fmt.Sprintf("[%v]%v", my.ID, my.RemoteAddr().String())
}
