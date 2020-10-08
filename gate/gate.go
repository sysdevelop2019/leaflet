package gate

import (
	"gitee.com/aarlin/leaflet/chanrpc"
	"gitee.com/aarlin/leaflet/log"
	"gitee.com/aarlin/leaflet/network"
	"net"
	"reflect"
)


type agent struct {
	conn     network.Conn
	closeAgentName string
	processor       network.Processor
	agentChanRPC    *chanrpc.Server
	userData interface{}
}

func (a *agent) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}
		//log.Debug("read message: %v", data)
		if a.processor != nil {
			msg, err := a.processor.Unmarshal(data)
			if err != nil {
				log.Debug("unmarshal message error: %v", err)
				break
			}
			err = a.processor.Route(msg, a)
			if err != nil {
				log.Debug("route message error: %v", err)
				break
			}
		}
	}
}

func (a *agent) OnClose() {
	if a.agentChanRPC != nil {
		err := a.agentChanRPC.Call0(a.closeAgentName, a)
		if err != nil {
			log.Error("chanrpc error: %v", err)
		}
	}
}

func (a *agent) WriteMsg(msg interface{}) {
	var data [][]byte
	var err error
	if a.processor != nil {
		data, err = a.processor.Marshal(msg)
		if err != nil {
			log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
			return
		}

	}else{
		if data2,ok := msg.([][]byte); !ok {
			log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
			return
		}else{
			data = data2
		}
	}
	err = a.conn.WriteMsg(data...)
	if err != nil {
		log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
	}
}

func (a *agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *agent) Close() {
	a.conn.Close()
}

func (a *agent) Destroy() {
	a.conn.Destroy()
}

func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}
