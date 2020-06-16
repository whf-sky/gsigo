package gsigo

//组创建一个新的路由器组。您应该添加所有具有公共中间件或相同路径前缀的路由。
//例如，所有使用公共中间件进行授权的路由都可以分组。
//relativePath 网站的相对路径
//c 组控制器
func Group(relativePath string, controller ...ctrFunc) *router{
	//初始化gin组路由信息
	routerObj.groupRelativePath = relativePath
	if routerObj.groupRouters == nil {
		routerObj.groupRouters = map[string]ctrFunc{}
	}
	//添加gin组路由信息
	routerObj.groupRouters[relativePath] = nil
	if len(controller) == 1 {
		routerObj.groupRouters[relativePath] = controller[0]
	}
	return routerObj
}

// Use 在组中添加一个中间件
//middleware 中间件控制器
func Use(middleware ctrFunc) {
	routerObj.Use(middleware)
}

// POST 添加 POST 路由信息
//relativePath 网站的相对路径
//c 中间件控制器
func POST(relativePath string, controller ctrFunc) {
	routerObj.POST(relativePath, controller)
}

// GET 添加 GET 路由信息
//relativePath 网站的相对路径
//c 中间件控制器
func GET(relativePath string, controller ctrFunc) {
	routerObj.GET(relativePath, controller)
}

// DELETE 添加 DELETE 路由信息
//relativePath 网站的相对路径
//c 中间件控制器
func DELETE(relativePath string, controller ctrFunc) {
	routerObj.DELETE(relativePath, controller)
}

// PATCH 添加 PATCH 路由信息
//relativePath 网站的相对路径
//c 中间件控制器
func PATCH(relativePath string, controller ctrFunc) {
	routerObj.PATCH(relativePath, controller)
}

// PUT 添加 PUT 路由信息
//relativePath 网站的相对路径
//c 中间件控制器
func PUT(relativePath string, controller ctrFunc) {
	routerObj.PUT(relativePath, controller)
}

// OPTIONS 添加 OPTIONS 路由信息
//relativePath 网站的相对路径
//c 中间件控制器
func OPTIONS(relativePath string, controller ctrFunc) {
	routerObj.OPTIONS(relativePath, controller)
}

// HEAD 添加 HEAD 路由信息
//relativePath 网站的相对路径
//c 中间件控制器
func HEAD(relativePath string, controller ctrFunc) {
	routerObj.HEAD(relativePath, controller)
}

// Any 注册一个匹配所有HTTP方法的路由。
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
//relativePath 网站的相对路径
//c 控制器方法
func Any(relativePath string, controller ctrFunc) {
	routerObj.Any(relativePath, controller)
}

// Static 添加静态资源路由
//relativePath 网站的相对路径
//filePath 文件的路径
func Static(relativePath string, filePath string) {
	routerObj.Static(relativePath, filePath)
}
