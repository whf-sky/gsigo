package gsigo

import (
	"sync"
	"time"
	"github.com/googollee/go-engine.io"
	gsocketio "github.com/googollee/go-socket.io"
)

// NewApp returns a new wsigo application.
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
	nsp string
	lock sync.RWMutex
	users map[string]map[string]int //map[uid]map[cid]conn num
	cids map[string]string //map[cid]uid
	Server *gsocketio.Server
}

func (s *socketio) newServer() {
	var err error
	s.Server, err = gsocketio.NewServer(&engineio.Options{
		PingInterval:time.Duration(Config.SocketIo.PingInterval) * time.Second,
		PingTimeout:time.Duration(Config.SocketIo.PingTimeout) * time.Second,
	})

	if err != nil {
		Log.Error(err)
	}
	s.register()
}

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

func (s *socketio)serve()  {
	err := s.Server.Serve()
	if err != nil {
		Log.Error(err)
	}
}

func (s *socketio)close()  {
	err := s.Server.Close()
	if err != nil {
		Log.Error(err)
	}
}

func (s *socketio) onConnect(event EventInterface){
	s.Server.OnConnect(s.nsp, s.connectHandle(s.nsp, event))
}

func (s *socketio)onEvent(eventName string, event EventInterface){
	s.Server.OnEvent(s.nsp, eventName, s.eventHandle(s.nsp, event ))
}

func (s *socketio) onError(event EventInterface){
	s.Server.OnError(s.nsp, s.errorHandle(s.nsp, event))
}

func (s *socketio) onDisconnect(event EventInterface){
	s.Server.OnDisconnect(s.nsp, s.disconnectHandle(s.nsp, event))
}

func (s *socketio) getNspHandler(nsp string) EventInterface{
	if evnet, ok := routerObj.nsp[nsp];ok{
		return evnet
	}
	return nil
}

func (s *socketio) funcHandle(eventType string, e EventInterface, conn gsocketio.Conn,  message string, err error) string {
	e.Init(eventType, conn, message, err)
	e.Prepare()
	e.Execute()
	e.Finish()
	return e.GetAckMsg()
}

func (s *socketio) groupHandle(eventType string, nsp string, conn gsocketio.Conn,  message string, err error) {
	if event := s.getNspHandler(nsp); event != nil{
		s.funcHandle(eventType, event, conn, message, err)
	}
}

func (s *socketio) connectHandle(nsp string, event EventInterface) func(conn gsocketio.Conn)error {
	return func(conn gsocketio.Conn)  error{
		s.groupHandle("connect", nsp, conn, "", nil)
		s.funcHandle("connect", event, conn, "", nil)
		return event.GetError()
	}
}

func (s *socketio) eventHandle(nsp string, event EventInterface) interface{} {
	var f interface{}
	if event.IsAck() == false {
		f = func(conn gsocketio.Conn, message string) {
			s.groupHandle("connect", nsp, conn, message, nil)
			s.funcHandle("event", event, conn, message, nil)
		}
	} else {
		f = func(conn gsocketio.Conn, message string) string {
			s.groupHandle("connect", nsp, conn, message, nil)
			return s.funcHandle("event", event, conn, message, nil)
		}
	}
	return f
}

func (s *socketio) errorHandle(nsp string, event EventInterface) func(err error) {
	return func(err error) {
		s.groupHandle("error", nsp, nil, "", err)
		s.funcHandle("error", event, nil, "", err)
	}
}

func (s *socketio) disconnectHandle(nsp string, e EventInterface) func(conn gsocketio.Conn, message string){
	return func(conn gsocketio.Conn, message string) {
		s.groupHandle("disconnect", nsp, conn, message, nil)
		s.funcHandle("disconnect", e, conn, message, nil)
	}
}