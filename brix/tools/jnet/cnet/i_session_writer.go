package cnet

//SocketWriter :
type SocketWriter interface {
	SendMessage(int, []byte) error
}
