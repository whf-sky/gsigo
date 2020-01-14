package gsigo

var routerObj *router

//newRouter function is new router
func newRouter() *router{
	routerObj = &router{
		relativePath: "/",
		nsp: nil,
		group: nil,
		socketio: map[string]map[string]map[string]EventInterface{},
		gin: map[string]map[string]gRouter{},
		cmd: []CmdInterface{},
	}
	return routerObj
}

//router For gin, socketio, cmd router
type router struct {
	relativePath string //gin relativePath
	nsp map[string]EventInterface //socketio nsp
	group map[string]ControllerInterface //gin group
	socketio map[string]map[string]map[string]EventInterface //socketio routers
	gin map[string]map[string]gRouter// gin routers
	cmd []CmdInterface //cmd routers
}

//event add router
//For OnConnect, OnEvent, OnError, OnDisconnect event
func (r *router) event(name string, event EventInterface, eventName string) {
	_, ok := r.socketio[routerObj.relativePath]
	if !ok {
		r.socketio[routerObj.relativePath] = map[string]map[string]EventInterface{}
	}

	_, ok = r.socketio[routerObj.relativePath][name]
	if !ok{
		r.socketio[routerObj.relativePath][name] = map[string]EventInterface{}
	}

	r.socketio[routerObj.relativePath][name][eventName] = event
}

//OnConnect add connect router
func (r *router) OnConnect(event EventInterface)  {
	r.event("onConnect", event, "_")
}

//OnEvent add event router
func (r *router) OnEvent(eventName string, event EventInterface){
	r.event("onEvent", event, eventName)
}

//OnError add error router
func (r *router) OnError(event EventInterface){
	r.event("onError", event, "_")
}

//OnDisconnect add disconnect router
func (r *router) OnDisconnect(event EventInterface){
	r.event("onDisconnect", event, "_")
}


type gRouter struct {
	relativePath string
	controller ControllerInterface
}

// Request registers a new request handle and middleware with the given path and method.
// The last handler should be the real handler, the other ones should be middleware that can and should be shared among different routes.
// For GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, Any requests the respective shortcut
func (r *router) Request(method string, relativePath string, controller ControllerInterface) {
	_, ok := r.gin[r.relativePath]
	if !ok {
		r.gin[r.relativePath] = map[string]gRouter{}
	}
	r.gin[r.relativePath][method] = gRouter{
		relativePath:relativePath,
		controller: controller,
	}
}

// Use adds middleware to the group.
func (r *router) Use(middleware ControllerInterface) {
	r.Request("Use", "",  middleware)
}

// POST is a shortcut for router.Request("POST", path, handle).
func (r *router) POST(relativePath string, controller ControllerInterface) {
	r.Request("POST", relativePath, controller)
}

// GET is a shortcut for router.Request("GET", path, handle).
func (r *router) GET(relativePath string, controller ControllerInterface) {
	r.Request("GET", relativePath, controller)
}

// DELETE is a shortcut for router.Request("DELETE", path, handle).
func (r *router) DELETE(relativePath string, controller ControllerInterface) {
	r.Request("DELETE", relativePath, controller)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle).
func (r *router) PATCH(relativePath string, controller ControllerInterface) {
	r.Request("PATCH", relativePath, controller)
}

// PUT is a shortcut for router.Request("PUT", path, handle).
func (r *router) PUT(relativePath string, controller ControllerInterface) {
	r.Request("PUT", relativePath, controller)
}

// OPTIONS is a shortcut for router.Request("OPTIONS", path, handle).
func (r *router) OPTIONS(relativePath string, controller ControllerInterface) {
	r.Request("OPTIONS", relativePath, controller)
}

// HEAD is a shortcut for router.Request("HEAD", path, handle).
func (r *router) HEAD(relativePath string, controller ControllerInterface) {
	r.Request("HEAD", relativePath, controller)
}

// Any registers a route that matches all the HTTP methodWsiApp.socketio.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (r *router) Any(relativePath string, controller ControllerInterface) {
	r.Request("Any", relativePath, controller)
}
