package network

import (
	"net"
)

type Conn interface {
	ReadMsg() ([]byte, error)
	WriteMsg(args ...[]byte) error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()

	Read(p []byte) (n int, err error)
	Write(p []byte)
}
