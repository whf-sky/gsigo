package gsigo

import (
	"errors"
	gsocketio "github.com/googollee/go-socket.io"
)

type Event struct {
	// is ack
	Ack bool
	//socketio connect
	Conn gsocketio.Conn
	//socketio event type
	EventType string

	//socketio message
	message string
	//socketio namespace
	namespace string
	//socketio error
	error error
	//socketio ack message
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
}

//SetUser is binding user
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

//GetUser is get user by connect id
func (e *Event) GetUser() string {
	// get read lock
	Socketio.lock.RLock()
	defer Socketio.lock.RUnlock()

	if val,ok := Socketio.cids[e.Conn.ID()];ok {
		return val
	}
	return ""
}

//GetCidsByUser get connect ids by user
func (e *Event) GetCidsByUser(uid string) map[string]int {
	// get read lock
	Socketio.lock.RLock()
	defer Socketio.lock.RUnlock()
	if val,ok := Socketio.users[uid];ok {
		return val
	}
	return nil
}

//IsAck is ack
func (e *Event) IsAck() bool {
	return e.Ack
}

//GetMessage get event message
func (e *Event) GetMessage() string {
	return e.message
}

//SetAckMsg set socketio 'event' event ack msg
func (e *Event) SetAckMsg(msg string) {
	e.ackMsg = msg
}

//GetAckMsg get socketio 'event' event ack msg
func (e *Event) GetAckMsg() string{
	return e.ackMsg
}

//GetNamespace get socketio namespace
func (e *Event) GetNamespace() string{
	return e.namespace
}

//SetError set socketio error
func (e *Event) SetError(text string) {
	e.error = errors.New(text)
}

//GetError is get socketio event error
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