package gsigo

var routerObj *router

//router gin、socketio 路由
type router struct {
	//gin 组的网站相对路径
	groupRelativePath string
	//socketio 组的网站相对路径
	nsp string
	//socketio 命名空间
	nspRouters map[string]EventInterface
	//gin 组路由规则
	groupRouters map[string]ControllerInterface
	//socketio 路由
	socketioRouters map[string]map[string]map[string]EventInterface
	// gin 路由
	ginRouters map[string][]gginRouter
	//cmd 路由
	cmdRouters map[string]CmdInterface
}

//添加cmd路由信息
//requestUri 请求URI
//gsigo.Cmd("test", &TestCmd{})
//go run main.go -request_uri='test'
//cmd 命令行控制器
func (r *router) Cmd(requestUri string, cmd CmdInterface) {
	//命名空间不存在时初始化信息
	r.cmdRouters[requestUri] = cmd
}

//event 添加事件路由
//事件：OnConnect, OnEvent, OnError, OnDisconnect
//name 事件
//event 事件控制器
//eventName 事件名称
func (r *router) event(name string, event EventInterface, eventName string) {
	//命名空间不存在时初始化信息
	_, ok := r.socketioRouters[routerObj.nsp]
	if !ok {
		r.socketioRouters[routerObj.nsp] = map[string]map[string]EventInterface{}
	}
	//事件不存在时初始化信息
	_, ok = r.socketioRouters[routerObj.nsp][name]
	if !ok{
		r.socketioRouters[routerObj.nsp][name] = map[string]EventInterface{}
	}
	//添加事件
	r.socketioRouters[routerObj.nsp][name][eventName] = event
}

//OnConnect 添加连接事件
//event 事件控制器
func (r *router) OnConnect(event EventInterface)  {
	r.event("onConnect", event, "_")
}

//OnEvent 添加socketio事件
//eventName socketio事件名称
//event 事件控制器
func (r *router) OnEvent(eventName string, event EventInterface){
	r.event("onEvent", event, eventName)
}

//OnError 添加错误事件
//event 事件控制器
func (r *router) OnError(event EventInterface){
	r.event("onError", event, "_")
}

//OnDisconnect 添加关闭连接事件
//event 事件控制器
func (r *router) OnDisconnect(event EventInterface){
	r.event("onDisconnect", event, "_")
}

//ggin 路由数据
type gginRouter struct {
	//方法
	method string
	//网站相对路径
	relativePath string
	//文件路径
	filePath string
	//控制器
	controller ControllerInterface
}

//请求用给定的路径和方法注册一个新的请求句柄和中间件。
//最后一个处理程序应该是真正的处理程序，其他的应该是可以而且应该在不同的路由之间共享的中间件。
// 事件： GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, Any requests the respective shortcut
func (r *router) Request(method string, relativePath string, controller ControllerInterface, filePaths ...string) {
	//初始化gin组路由
	_, ok := r.ginRouters[r.groupRelativePath]
	if !ok {
		r.ginRouters[r.groupRelativePath] = []gginRouter{}
	}
	//添加文件信息
	filePath := ""
	if len(filePaths) == 1{
		filePath = filePaths[0]
	}
	//添加gin的路由信息
	r.ginRouters[r.groupRelativePath] = append(r.ginRouters[r.groupRelativePath], gginRouter{
		method:method,
		relativePath:relativePath,
		controller: controller,
		filePath: filePath,
	})
}

// Use 在组中添加一个中间件
//middleware 中间件控制器
func (r *router) Use(middleware ControllerInterface) {
	r.Request("USE", "",  middleware)
}

// POST 添加 POST 路由信息
//relativePath 网站的相对路径
//c 中间件控制器
func (r *router) POST(relativePath string, controller ControllerInterface) {
	r.Request("POST", relativePath, controller)
}

// GET 添加 GET 路由信息
//relativePath 网站的相对路径
//c 中间件控制器
func (r *router) GET(relativePath string, controller ControllerInterface) {
	r.Request("GET", relativePath, controller)
}

// DELETE 添加 DELETE 路由信息
//relativePath 网站的相对路径
//c 中间件控制器
func (r *router) DELETE(relativePath string, controller ControllerInterface) {
	r.Request("DELETE", relativePath, controller)
}

// PATCH 添加 PATCH 路由信息
//relativePath 网站的相对路径
//c 中间件控制器
func (r *router) PATCH(relativePath string, controller ControllerInterface) {
	r.Request("PATCH", relativePath, controller)
}

// PUT 添加 PUT 路由信息
//relativePath 网站的相对路径
//c 中间件控制器
func (r *router) PUT(relativePath string, controller ControllerInterface) {
	r.Request("PUT", relativePath, controller)
}

// OPTIONS 添加 OPTIONS 路由信息
//relativePath 网站的相对路径
//c 中间件控制器
func (r *router) OPTIONS(relativePath string, controller ControllerInterface) {
	r.Request("OPTIONS", relativePath, controller)
}

// HEAD 添加 HEAD 路由信息
//relativePath 网站的相对路径
//c 中间件控制器
func (r *router) HEAD(relativePath string, controller ControllerInterface) {
	r.Request("HEAD", relativePath, controller)
}

// Any 注册一个匹配所有HTTP方法的路由。
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
//relativePath 网站的相对路径
//c 控制器方法
func (r *router) Any(relativePath string, controller ControllerInterface) {
	r.Request("ANY", relativePath, controller)
}

// Static 添加静态资源路由
//relativePath 网站的相对路径
//filePath 文件的路径
func (r *router) Static(relativePath string, filePath string) {
	r.Request("STATIC", relativePath,nil, filePath)
}

func init()  {
	//实例化路由
	routerObj = &router{
		groupRelativePath: "/",
		nsp: "/",
		nspRouters: nil,
		groupRouters: nil,
		socketioRouters: map[string]map[string]map[string]EventInterface{},
		ginRouters: map[string][]gginRouter{},
		cmdRouters: map[string]CmdInterface{},
	}
}
