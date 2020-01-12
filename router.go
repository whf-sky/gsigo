package gsigo

var routerObj *router

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

func OnConnect(event EventInterface){
	routerObj.OnConnect(event)
}

func OnEvent(eventName string, event EventInterface){
	routerObj.OnEvent(eventName, event)
}

func OnError(event EventInterface){
	routerObj.OnError(event)
}

func OnDisconnect(event EventInterface){
	routerObj.OnDisconnect(event)
}

// Group creates a new router group. You should add all the routes that have common middlewares or the same path prefix.
// For example, all the routes that use a common middleware for authorization could be grouped.
func Group(relativePath string, controller ...ControllerInterface) *router{
	routerObj.relativePath = relativePath
	if routerObj.group == nil{
		routerObj.group = map[string]ControllerInterface{}
	}
	routerObj.group[relativePath] = nil
	if len(controller) == 1 {
		routerObj.group[relativePath] = controller[0]
	}
	return routerObj
}

// Use adds middleware to the group, see example code in GitHub.
func Use(middleware ControllerInterface) {
	routerObj.Use(middleware)
}

// POST is a shortcut for router.Handle("POST", path, handle).
func POST(relativePath string, controller ControllerInterface) {
	routerObj.POST(relativePath, controller)
}

// GET is a shortcut for router.Handle("GET", path, handle).
func GET(relativePath string, controller ControllerInterface) {
	routerObj.GET(relativePath, controller)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle).
func DELETE(relativePath string, controller ControllerInterface) {
	routerObj.DELETE(relativePath, controller)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle).
func PATCH(relativePath string, controller ControllerInterface) {
	routerObj.PATCH(relativePath, controller)
}

// PUT is a shortcut for router.Handle("PUT", path, handle).
func PUT(relativePath string, controller ControllerInterface) {
	routerObj.PUT(relativePath, controller)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle).
func OPTIONS(relativePath string, controller ControllerInterface) {
	routerObj.OPTIONS(relativePath, controller)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle).
func HEAD(relativePath string, controller ControllerInterface) {
	routerObj.HEAD(relativePath, controller)
}

// Any registers a route that matches all the HTTP methodWsiApp.socketio.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func Any(relativePath string, controller ControllerInterface) {
	routerObj.Any(relativePath, controller)
}

func CmdRouter(cmd CmdInterface) *router{
	routerObj.cmd = append(routerObj.cmd, cmd)
	return routerObj
}

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

type router struct {
	relativePath string
	nsp map[string]EventInterface
	group map[string]ControllerInterface
	socketio map[string]map[string]map[string]EventInterface
	gin map[string]map[string]gRouter
	cmd []CmdInterface
}

func (r *router) event(name string, event EventInterface, eventName string){
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

func (r *router) OnConnect(event EventInterface)  {
	r.event("onConnect", event, "_")
}

func (r *router) OnEvent(eventName string, event EventInterface){
	r.event("onEvent", event, eventName)
}

func (r *router) OnError(event EventInterface){
	r.event("onError", event, "_")
}

func (r *router) OnDisconnect(event EventInterface){
	r.event("onDisconnect", event, "_")
}


type gRouter struct {
	relativePath string
	controller ControllerInterface
}

func (r *router) request(method string, relativePath string, controller ControllerInterface){
	_, ok := r.gin[r.relativePath]
	if !ok {
		r.gin[r.relativePath] = map[string]gRouter{}
	}
	r.gin[r.relativePath][method] = gRouter{
		relativePath:relativePath,
		controller: controller,
	}
}

// Use adds middleware to the group, see example code in GitHub.
func (r *router) Use(middleware ControllerInterface) {
	r.request("Use", "", middleware)
}

// POST is a shortcut for router.Handle("POST", path, handle).
func (r *router) POST(relativePath string, controller ControllerInterface) {
	r.request("POST", relativePath, controller)
}

// GET is a shortcut for router.Handle("GET", path, handle).
func (r *router) GET(relativePath string, controller ControllerInterface) {
	r.request("GET", relativePath, controller)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle).
func (r *router) DELETE(relativePath string, controller ControllerInterface) {
	r.request("DELETE", relativePath, controller)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle).
func (r *router) PATCH(relativePath string, controller ControllerInterface) {
	r.request("PATCH", relativePath, controller)
}

// PUT is a shortcut for router.Handle("PUT", path, handle).
func (r *router) PUT(relativePath string, controller ControllerInterface) {
	r.request("PUT", relativePath, controller)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle).
func (r *router) OPTIONS(relativePath string, controller ControllerInterface) {
	r.request("OPTIONS", relativePath, controller)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle).
func (r *router) HEAD(relativePath string, controller ControllerInterface) {
	r.request("HEAD", relativePath, controller)
}

// Any registers a route that matches all the HTTP methodWsiApp.socketio.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (r *router) Any(relativePath string, controller ControllerInterface) {
	r.request("Any", relativePath, controller)
}
