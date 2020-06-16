package gsigo

import (
	"github.com/googollee/go-engine.io"
	"github.com/googollee/go-socket.io"
	"sync"
	"time"
)

type evtFunc func() EventInterface

//实例化一个socketio服务
func newGsocketio() *gsocketio {
	Gsocketio = &gsocketio{
		nsp: "/",
		users: map[string]map[string]socketio.Conn{},
		cids: map[string]string{},
		conns: map[string]socketio.Conn{},
	}
	return Gsocketio.newServer()
}

type gsocketio struct {
	//命名空间
	nsp string
	//读锁
	lock sync.RWMutex
	//用户与连接绑定关系关系
	//map[用户编号]map[连接编号]Conn
	users map[string]map[string]socketio.Conn
	//连接编号与用户编号关系
	//map[连接编号]用户编号
	cids map[string]string
	//socket连接
	conns map[string]socketio.Conn
	//socketio服务
	Server *socketio.Server
}

//运行服务
func (s *gsocketio) run() *gsocketio {
	defer s.close()
	go s.serve()
	newGgin().socketioRouter(s.Server).run()
	return s
}

//newServer 实例一个 socketio 服务
func (s *gsocketio) newServer() *gsocketio {
	var err error
	s.Server, err = socketio.NewServer(&engineio.Options{
		PingInterval:time.Duration(Config.Socket.PingInterval) * time.Second,
		PingTimeout:time.Duration(Config.Socket.PingTimeout) * time.Second,
	})
	if err != nil {
		Log.Error(err)
	}
	s.registerRouter()
	return s
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
func (s *gsocketio) onConnect(event evtFunc){
	e := event()
	s.Server.OnConnect(s.nsp, func(conn socketio.Conn)  error {
		s.lock.Lock()//获取写锁
		//添加连接
		if _,ok := s.conns[conn.ID()];!ok{
			s.conns[conn.ID()] = conn
		}
		s.lock.Unlock()//写锁解锁
		s.groupHandle("connect", s.nsp, conn, "", nil)
		s.funcHandle("connect", e, conn, "", nil)
		return e.GetError()
	})
}

//onEvent add evnet event
//eventName socketio事件名称
//event 事件控制器
func (s *gsocketio) onEvent(eventName string, event evtFunc){
	e := event()
	var f interface{}
	if e.IsAck() == false {
		f = func(conn socketio.Conn, message string) {
			s.groupHandle("event", s.nsp, conn, message, nil)
			s.funcHandle("event", e, conn, message, nil)
		}
	} else {
		f = func(conn socketio.Conn, message string) string {
			s.groupHandle("event", s.nsp, conn, message, nil)
			return s.funcHandle("event", e, conn, message, nil)
		}
	}
	s.Server.OnEvent(s.nsp, eventName, f)
}

//onError 添加一个错误事件路由
//event 事件控制器
func (s *gsocketio) onError(event evtFunc){
	e := event()
	s.Server.OnError(s.nsp, func(err error) {
		s.groupHandle("error", s.nsp, nil, "", err)
		s.funcHandle("error", e, nil, "", err)
	})
}

//onDisconnect 添加一个关闭事件路由
//event 事件控制器
func (s *gsocketio) onDisconnect(event evtFunc){
	e := event()
	s.Server.OnDisconnect(s.nsp, func(conn socketio.Conn, message string) {
		s.lock.Lock()//获取写锁
		cid := conn.ID()
		//删除cid与uid映射关系
		uid, ok := s.cids[cid];
		if ok {
			delete(s.cids, cid)
		}
		//从socket连接集合里删除连接
		if _,ok := s.conns[cid]; ok {
			delete(s.conns, cid)
		}
		//删除用户绑定的连接
		if conns, ok := s.users[uid]; ok {
			delete(conns, cid)
			if len(conns) == 0 {
				delete(s.users, uid)
			}
		}
		s.lock.Unlock()//写锁解锁
		s.groupHandle("disconnect", s.nsp, conn, message, nil)
		s.funcHandle("disconnect", e, conn, message, nil)
	})
}