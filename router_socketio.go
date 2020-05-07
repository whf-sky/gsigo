package gsigo
//Nsp 添加一个命名空间
//event 事件控制器
func Nsp(nsp string, event ...EventInterface) *router {
	//设置命名空间
	routerObj.nsp = nsp
	//命名空间路由不存在时初始化信息
	if routerObj.nspRouters == nil{
		routerObj.nspRouters = map[string]EventInterface{}
	}
	//设置命名空间路由信息
	routerObj.nspRouters[nsp] = nil
	if len(event) == 1 {
		routerObj.nspRouters[nsp] = event[0]
	}
	return routerObj
}

//OnConnect 添加连接事件路由
//event 事件控制器
func OnConnect(event EventInterface){
	routerObj.OnConnect(event)
}

//OnEvent 添加事件路由
//eventName socketio事件名称
//event 事件控制器
func OnEvent(eventName string, event EventInterface){
	routerObj.OnEvent(eventName, event)
}

//OnError 添加错误事件路由
//event 事件控制器
func OnError(event EventInterface){
	routerObj.OnError(event)
}

//OnDisconnect 添加关闭事件路由
//event 事件控制器
func OnDisconnect(event EventInterface){
	routerObj.OnDisconnect(event)
}

