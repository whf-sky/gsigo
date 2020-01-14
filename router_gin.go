package gsigo

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
