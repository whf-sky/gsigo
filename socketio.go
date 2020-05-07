package gsigo

import (
	"github.com/googollee/go-engine.io"
	"github.com/googollee/go-socket.io"
	"sync"
	"time"
)

//NewApp returns a new wsigo application.
func newGsocketio() *gsocketio {
	s := &gsocketio{
		nsp: "/",
		users: map[string]map[string]int{},
		cids: map[string]string{},
	}
	s.newServer()
	return s
}

type gsocketio struct {
	//命名空间
	nsp string
	//读锁
	lock sync.RWMutex
	//用户与连接绑定关系关系
	//map[用户编号]map[连接编号]conn num
	users map[string]map[string]int

	//连接编号与用户编号关系
	//map[连接编号]用户编号
	cids map[string]string

	//socketio服务
	Server *socketio.Server
}

//newServer 实例一个 socketio 服务
func (s *gsocketio) newServer() {
	var err error
	s.Server, err = socketio.NewServer(&engineio.Options{
		PingInterval:time.Duration(Config.Socket.PingInterval) * time.Second,
		PingTimeout:time.Duration(Config.Socket.PingTimeout) * time.Second,
	})
	if err != nil {
		Log.Error(err)
	}
	s.registerRouter()
}

//registerRouter 注册路由 onConnect, onEvent, onError, onDisconnect router
func (s *gsocketio) registerRouter(){
	for nsp, events:=range routerObj.socketioRouters  {
		s.nsp = nsp
		for name, eEvents :=  range events  {
			for eEventName, event :=  range eEvents {
				switch name {
				case "onConnect":
					s.onConnect(event)
				case "onEvent":
					s.onEvent(eEventName, event)
				case "onError":
					s.onError(event)
				case "onDisconnect":
					s.onDisconnect(event)
				}
			}

		}
	}
}

//serve socketio 服务
func (s *gsocketio) serve()  {
	err := s.Server.Serve()
	if err != nil {
		Log.Error(err)
	}
}

//close 关闭socketio连接
func (s *gsocketio) close()  {
	err := s.Server.Close()
	if err != nil {
		Log.Error(err)
	}
}

//groupHandle onConnect, onEvent, onError, onDisconnect句柄
func (s *gsocketio) funcHandle(eventType string, e EventInterface, conn socketio.Conn,  message string, err error) string {
	e.Init(eventType, conn, message, err)
	e.Prepare()
	e.Execute()
	e.Finish()
	return e.GetAckMsg()
}

//getNspHandler 获取命名空间句柄
func (s *gsocketio) getNspHandler(nsp string) EventInterface{
	if event, ok := routerObj.nspRouters[nsp];ok{
		return event
	}
	return nil
}

//groupHandle 添加组句柄 onConnect, onEvent, onError, onDisconnect
func (s *gsocketio) groupHandle(eventType string, nsp string, conn socketio.Conn,  message string, err error) {
	if event := s.getNspHandler(nsp); event != nil{
		s.funcHandle(eventType, event, conn, message, err)
	}
}

//onConnect add connect event
//event 事件控制器
func (s *gsocketio) onConnect(event EventInterface){
	s.Server.OnConnect(s.nsp, func(conn socketio.Conn)  error{
		s.groupHandle("connect", s.nsp, conn, "", nil)
		s.funcHandle("connect", event, conn, "", nil)
		return event.GetError()
	})
}

//onEvent add evnet event
//eventName socketio事件名称
//event 事件控制器
func (s *gsocketio)onEvent(eventName string, event EventInterface){
	var f interface{}
	if event.IsAck() == false {
		f = func(conn socketio.Conn, message string) {
			s.groupHandle("event", s.nsp, conn, message, nil)
			s.funcHandle("event", event, conn, message, nil)
		}
	} else {
		f = func(conn socketio.Conn, message string) string {
			s.groupHandle("event", s.nsp, conn, message, nil)
			return s.funcHandle("event", event, conn, message, nil)
		}
	}
	s.Server.OnEvent(s.nsp, eventName, f)
}

//onError 添加一个错误事件路由
//event 事件控制器
func (s *gsocketio) onError(event EventInterface){
	s.Server.OnError(s.nsp, func(err error) {
		s.groupHandle("error", s.nsp, nil, "", err)
		s.funcHandle("error", event, nil, "", err)
	})
}

//onDisconnect 添加一个关闭事件路由
//event 事件控制器
func (s *gsocketio) onDisconnect(event EventInterface){
	s.Server.OnDisconnect(s.nsp, func(conn socketio.Conn, message string) {
		s.groupHandle("disconnect", s.nsp, conn, message, nil)
		s.funcHandle("disconnect", event, conn, message, nil)
	})
}