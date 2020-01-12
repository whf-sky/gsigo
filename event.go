package gsigo

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/whf-sky/gsigo/log"
	gsocketio "github.com/googollee/go-socket.io"
)

type Event struct {
	Ack bool
	Conn gsocketio.Conn
	Log *logrus.Entry
	EventType string

	message string
	namespace string
	error error
	ackMsg string
}

type EventInterface interface {
	Init(eventType string, Conn gsocketio.Conn, Message string, Err error)
	Prepare()
	Execute()
	Finish()

	IsAck() bool
	GetError() error
	GetAckMsg() string
}

func (e *Event)Init(eventType string, conn gsocketio.Conn, message string, err error)  {
	e.Conn = conn
	e.EventType = eventType
	e.message = message
	e.error = err
	e.ackMsg = ""
	fields := logrus.Fields{
		"logid":log.GenerateLogid(),
	}
	if conn != nil {
		fields["ip"] = conn.RemoteAddr()
		fields["cid"] = conn.ID()
		e.namespace = conn.Namespace()
	}
	e.Log = Log.WithFields(fields)
}

func (e *Event) SetUser(uid string)  {
	cid := e.Conn.ID()
	// get write lock
	Socketio.lock.Lock()
	defer Socketio.lock.Unlock()
	if _,ok := Socketio.cids[cid];!ok{
		Socketio.cids[cid] = uid
	}
	if _,ok := Socketio.users[uid];!ok{
		Socketio.users[uid] = map[string]int{}
	}
	if _,ok := Socketio.users[uid][cid];!ok{
		Socketio.users[uid][cid] = len(Socketio.users[uid])+1
	}
}

func (e *Event) GetUser() string {
	// get read lock
	Socketio.lock.RLock()
	defer Socketio.lock.RUnlock()

	if val,ok := Socketio.cids[e.Conn.ID()];ok {
		return val
	}
	return ""
}

func (e *Event) GetCidsByUser(uid string) map[string]int {
	// get read lock
	Socketio.lock.RLock()
	defer Socketio.lock.RUnlock()
	if val,ok := Socketio.users[uid];ok {
		return val
	}
	return nil
}

func (e *Event) IsAck() bool {
	return e.Ack
}

func (e *Event) GetMessage() string {
	return e.message
}

func (e *Event) SetAckMsg(msg string) {
	e.ackMsg = msg
}

func (e *Event) GetAckMsg() string{
	return e.ackMsg
}

func (e *Event) GetNamespace() string{
	return e.namespace
}

func (e *Event) SetError(text string) {
	e.error = errors.New(text)
}

func (e *Event) GetError() error {
	return e.error
}

// Prepare runs after Init before request function execution.
func (e *Event) Prepare() {}


// Execute runs event function.
func (e *Event) Execute() {
	e.Conn.Emit("error", "Event Not Allowed")
}

// Finish runs after request function execution.
func (e *Event) Finish() {}