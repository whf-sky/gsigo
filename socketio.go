package gsigo

import (
	"sync"
	"time"
	"github.com/googollee/go-engine.io"
	gsocketio "github.com/googollee/go-socket.io"
)

//NewApp returns a new wsigo application.
func newSocketio() *socketio {
	s := &socketio{
		nsp: "/",
		users: map[string]map[string]int{},
		cids: map[string]string{},
	}
	s.newServer()
	return s
}

type socketio struct {
	nsp string //namespace
	//lock is users and cids lock
	lock sync.RWMutex

	//users is user and socketio connect id
	//map[uid]map[cid]conn num
	users map[string]map[string]int

	//cids is socketio connect id and user id
	//map[cid]uid
	cids map[string]string

	// socketio server
	Server *gsocketio.Server
}

//newServer new socketio sever
func (s *socketio) newServer() {
	var err error
	s.Server, err = gsocketio.NewServer(&engineio.Options{
		PingInterval:time.Duration(Config.Socket.PingInterval) * time.Second,
		PingTimeout:time.Duration(Config.Socket.PingTimeout) * time.Second,
	})
	if err != nil {
		Log.Error(err)
	}
	s.register()
}

//register onConnect, onEvent, onError, onDisconnect router
func (s *socketio) register(){
	for nsp, events:=range routerObj.socketio  {
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

//serve socketio server
func (s *socketio) serve()  {
	err := s.Server.Serve()
	if err != nil {
		Log.Error(err)
	}
}

//close socketio close
func (s *socketio) close()  {
	err := s.Server.Close()
	if err != nil {
		Log.Error(err)
	}
}

//groupHandle is onConnect, onEvent, onError, onDisconnect handle
func (s *socketio) funcHandle(eventType string, e EventInterface, conn gsocketio.Conn,  message string, err error) string {
	e.Init(eventType, conn, message, err)
	e.Prepare()
	e.Execute()
	e.Finish()
	return e.GetAckMsg()
}

//getNspHandler get namespace handle
func (s *socketio) getNspHandler(nsp string) EventInterface{
	if event, ok := routerObj.nsp[nsp];ok{
		return event
	}
	return nil
}

//groupHandle add group handle before execute onConnect, onEvent, onError, onDisconnect
func (s *socketio) groupHandle(eventType string, nsp string, conn gsocketio.Conn,  message string, err error) {
	if event := s.getNspHandler(nsp); event != nil{
		s.funcHandle(eventType, event, conn, message, err)
	}
}

//onConnect add connect event
func (s *socketio) onConnect(event EventInterface){
	s.Server.OnConnect(s.nsp, func(conn gsocketio.Conn)  error{
		s.groupHandle("connect", s.nsp, conn, "", nil)
		s.funcHandle("connect", event, conn, "", nil)
		return event.GetError()
	})
}

//onEvent add evnet event
func (s *socketio)onEvent(eventName string, event EventInterface){
	var f interface{}
	if event.IsAck() == false {
		f = func(conn gsocketio.Conn, message string) {
			s.groupHandle("connect", s.nsp, conn, message, nil)
			s.funcHandle("event", event, conn, message, nil)
		}
	} else {
		f = func(conn gsocketio.Conn, message string) string {
			s.groupHandle("connect", s.nsp, conn, message, nil)
			return s.funcHandle("event", event, conn, message, nil)
		}
	}
	s.Server.OnEvent(s.nsp, eventName, f)
}

//onError add error event
func (s *socketio) onError(event EventInterface){
	s.Server.OnError(s.nsp, func(err error) {
		s.groupHandle("error", s.nsp, nil, "", err)
		s.funcHandle("error", event, nil, "", err)
	})
}

//onDisconnect add disconnect event
func (s *socketio) onDisconnect(event EventInterface){
	s.Server.OnDisconnect(s.nsp, func(conn gsocketio.Conn, message string) {
		s.groupHandle("disconnect", s.nsp, conn, message, nil)
		s.funcHandle("disconnect", event, conn, message, nil)
	})
}