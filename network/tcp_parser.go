package network

type TcpParser interface {
	SetMsgLen(lenMsgLen int, minMsgLen uint32, maxMsgLen uint32)
	SetByteOrder(littleEndian bool)
	Read(conn Conn) ([]byte, error)
	Write(conn Conn, args ...[]byte) error
}