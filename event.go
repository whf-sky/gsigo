package gsigo

import (
	"errors"
	"github.com/googollee/go-socket.io"
)

type Event struct {
	//是否应答
	Ack bool
	//socketio connect
	Conn socketio.Conn
	//socketio 事件类型
	EventType string
	//socketio 消息
	message string
	//socketio 命名空间
	namespace string
	//socketio 错误信息
	error error
	//socketio 应答消息
	ackMsg string
}

//定义了一些基本的socketio请求处理程序操作
type EventInterface interface {
	Init(eventType string, Conn socketio.Conn, Message string, Err error)
	Prepare()
	Execute()
	Finish()

	IsAck() bool
	GetError() error
	GetAckMsg() string
}

//初始化事件操作的默认值。
func (e *Event)Init(eventType string, conn socketio.Conn, message string, err error)  {
	e.Conn = conn
	e.EventType = eventType
	e.message = message
	e.error = err
	e.ackMsg = ""
}

//SetUser 绑定用户信息
//uid 用户标识
//绑的的是：用户与链接关系、链接与用户关系
func (e *Event) SetUser(uid string)  {
	cid := e.Conn.ID()
	//获取写锁
	Gsocketio.lock.Lock()
	//写锁解锁
	defer Gsocketio.lock.Unlock()
	//连接编号与用户编号关系绑定
	if _,ok := Gsocketio.cids[cid];!ok{
		Gsocketio.cids[cid] = uid
	}
	//用户编号与连接关系绑定
	if _,ok := Gsocketio.users[uid];!ok{
		Gsocketio.users[uid] = map[string]int{}
	}
	if _,ok := Gsocketio.users[uid][cid];!ok{
		Gsocketio.users[uid][cid] = len(Gsocketio.users[uid])+1
	}
}

//GetUser 根据链接编号获取用户
func (e *Event) GetUser() string {
	//获取读锁
	Gsocketio.lock.RLock()
	//读锁解锁
	defer Gsocketio.lock.RUnlock()
	//根据连接获取用户
	if val,ok := Gsocketio.cids[e.Conn.ID()];ok {
		return val
	}
	return ""
}

//GetCidsByUser 根据用户获取连接信息
func (e *Event) GetCidsByUser(uid string) map[string]int {
	//获取读锁
	Gsocketio.lock.RLock()
	//读锁解锁
	defer Gsocketio.lock.RUnlock()
	//根据用户获取多个连接信息
	if val,ok := Gsocketio.users[uid];ok {
		return val
	}
	return nil
}

//IsAck 获取是否应答信息
func (e *Event) IsAck() bool {
	return e.Ack
}

//GetMessage 获取事件信息
func (e *Event) GetMessage() string {
	return e.message
}

//SetAckMsg 设置事件的应答信息，前提是开启了应答机制
func (e *Event) SetAckMsg(msg string) {
	e.ackMsg = msg
}

//GetAckMsg 获取应答信息
func (e *Event) GetAckMsg() string{
	return e.ackMsg
}

//GetNamespace 获取命名空间
func (e *Event) GetNamespace() string{
	return e.namespace
}

//SetError 设置错误信息
func (e *Event) SetError(text string) {
	e.error = errors.New(text)
}

//GetError 获取错误信息
func (e *Event) GetError() error {
	return e.error
}

// Prepare 在请求函数执行之前，在Init之后运行。
func (e *Event) Prepare() {}


// Execute 件事的执行函数
func (e *Event) Execute() {
	e.Conn.Emit("error", "Event Not Allowed")
}

// Finish 在请求函数执行后运行。
func (e *Event) Finish() {}