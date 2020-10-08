package gate

import (
	"gitee.com/aarlin/leaflet/chanrpc"
	"gitee.com/aarlin/leaflet/network"
	"time"
)

type ClientGate struct {
	wsClient *network.WSClient
	tcpClient *network.TCPClient

	ConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	Processor       network.Processor
	AgentChanRPC    *chanrpc.Server

	NewAgentName	string
	CloseAgentName	string

	// websocket
	WSAddr      string
	HTTPTimeout time.Duration
	ReadTimeOut time.Duration
	ConnectInterval  time.Duration
	AutoReconnect    bool

	// tcp
	TCPAddr      string
	TcpParser network.TcpParser
	//LenMsgLen    int
	//LittleEndian bool
}

func (c *ClientGate) Run() {
	c.init()

	var wsClient *network.WSClient
	if c.WSAddr != "" {
		wsClient = new(network.WSClient)
		wsClient.Addr = c.WSAddr
		wsClient.ConnNum = c.ConnNum
		wsClient.PendingWriteNum = c.PendingWriteNum
		wsClient.MaxMsgLen = c.MaxMsgLen
		wsClient.ConnectInterval = c.ConnectInterval
		wsClient.AutoReconnect = c.AutoReconnect
		wsClient.HandshakeTimeout = c.HTTPTimeout
		wsClient.ReadTimeOut = c.ReadTimeOut
		wsClient.NewAgent = func(conn *network.WSConn) network.Agent {
			a := &agent{conn: conn, closeAgentName: c.CloseAgentName,processor: c.Processor,agentChanRPC: c.AgentChanRPC}
			if c.AgentChanRPC != nil {
				c.AgentChanRPC.Go(c.NewAgentName, a)
			}
			return a
		}
	}

	var tcpClient *network.TCPClient
	if c.TCPAddr != "" {
		tcpClient = new(network.TCPClient)
		tcpClient.Addr = c.TCPAddr
		tcpClient.ConnNum = c.ConnNum
		tcpClient.PendingWriteNum = c.PendingWriteNum
		tcpClient.ConnectInterval = c.ConnectInterval
		tcpClient.AutoReconnect = c.AutoReconnect
		tcpClient.TcpParser = c.TcpParser
		tcpClient.NewAgent = func(conn *network.TCPConn) network.Agent {
			a := &agent{conn: conn, closeAgentName: c.CloseAgentName,processor: c.Processor,agentChanRPC: c.AgentChanRPC}
			if c.AgentChanRPC != nil {
				c.AgentChanRPC.Go(c.NewAgentName, a)
			}
			return a
		}
	}

	if wsClient != nil {
		wsClient.Start()
	}
	c.wsClient = wsClient

	if tcpClient != nil {
		tcpClient.Start()
	}
	c.tcpClient = tcpClient
}

func (c *ClientGate) init()  {
	if  c.CloseAgentName == ""{
		c.CloseAgentName = "CloseAgent"
	}
	if  c.NewAgentName == ""{
		c.NewAgentName = "NewAgent"
	}
}

func (c *ClientGate) Stop()  {
	if c.wsClient != nil {
		c.wsClient.Close()
	}
	c.wsClient = nil
}

func (c *ClientGate) OnDestroy() {}