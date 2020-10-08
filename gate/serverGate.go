package gate

import (
	"gitee.com/aarlin/leaflet/chanrpc"
	"gitee.com/aarlin/leaflet/network"
	"time"
)

type ServerGate struct {
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	Processor       network.Processor
	AgentChanRPC    *chanrpc.Server

	NewAgentName	string
	CloseAgentName	string

	// websocket
	WSAddr      string
	HTTPTimeout time.Duration
	CertFile    string
	KeyFile     string

	// tcp
	TCPAddr      string
	TcpParser network.TcpParser
}

func (gate *ServerGate) Run(closeSig chan bool) {
	gate.init()

	var wsServer *network.WSServer
	if gate.WSAddr != "" {
		wsServer = new(network.WSServer)
		wsServer.Addr = gate.WSAddr
		wsServer.MaxConnNum = gate.MaxConnNum
		wsServer.PendingWriteNum = gate.PendingWriteNum
		wsServer.MaxMsgLen = gate.MaxMsgLen
		wsServer.HTTPTimeout = gate.HTTPTimeout
		wsServer.CertFile = gate.CertFile
		wsServer.KeyFile = gate.KeyFile
		wsServer.NewAgent = func(conn *network.WSConn) network.Agent {
			a := &agent{conn: conn, closeAgentName: gate.CloseAgentName,processor:gate.Processor,agentChanRPC:gate.AgentChanRPC}
			if gate.AgentChanRPC != nil {
				gate.AgentChanRPC.Go(gate.NewAgentName, a)
			}
			return a
		}
	}

	var tcpServer *network.TCPServer
	if gate.TCPAddr != "" {
		tcpServer = new(network.TCPServer)
		tcpServer.Addr = gate.TCPAddr
		tcpServer.MaxConnNum = gate.MaxConnNum
		tcpServer.PendingWriteNum = gate.PendingWriteNum
		tcpServer.TcpParser = gate.TcpParser
		//tcpServer.LenMsgLen = gate.LenMsgLen
		//tcpServer.MaxMsgLen = gate.MaxMsgLen
		//tcpServer.LittleEndian = gate.LittleEndian
		tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
			a := &agent{conn: conn, closeAgentName: gate.CloseAgentName,processor:gate.Processor,agentChanRPC:gate.AgentChanRPC}
			if gate.AgentChanRPC != nil {
				gate.AgentChanRPC.Go(gate.NewAgentName, a)
			}
			return a
		}
	}

	if wsServer != nil {
		wsServer.Start()
	}
	if tcpServer != nil {
		tcpServer.Start()
	}
	<-closeSig
	if wsServer != nil {
		wsServer.Close()
	}
	if tcpServer != nil {
		tcpServer.Close()
	}
}


func (gate *ServerGate) init()  {
	if  gate.CloseAgentName == ""{
		gate.CloseAgentName = "CloseAgent"
	}
	if  gate.NewAgentName == ""{
		gate.NewAgentName = "NewAgent"
	}
}

func (gate *ServerGate) OnDestroy() {}

