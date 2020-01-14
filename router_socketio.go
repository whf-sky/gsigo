package gsigo
//Nsp add socketio namespace
func Nsp(nsp string, event ...EventInterface) *router {
	routerObj.relativePath = nsp
	if routerObj.nsp == nil{
		routerObj.nsp = map[string]EventInterface{}
	}
	routerObj.nsp[nsp] = nil
	if len(event) == 1 {
		routerObj.nsp[nsp] = event[0]
	}
	return routerObj
}

//OnConnect add connect router
func OnConnect(event EventInterface){
	routerObj.OnConnect(event)
}

//OnEvent add event router
func OnEvent(eventName string, event EventInterface){
	routerObj.OnEvent(eventName, event)
}

//OnError add error router
func OnError(event EventInterface){
	routerObj.OnError(event)
}

//OnDisconnect add disconnect router
func OnDisconnect(event EventInterface){
	routerObj.OnDisconnect(event)
}

