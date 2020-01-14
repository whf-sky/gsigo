package gsigo

import (
	"fmt"
	ggin "github.com/gin-gonic/gin"
	"github.com/tabalt/gracehttp"
	"os"
	"reflect"
	"strings"
)

// NewApp returns a new wsigo application.
func newGin() *gin {
	s := &gin{}
	s.newServer()
	return s
}

//gin
//Server is *gin.Engine
//rgroup *gin.RouterGroup
type gin struct {
	Server *ggin.Engine
	rgroup *ggin.RouterGroup
}

//newServer create gin
func (g *gin) newServer() {
	if !Config.Debug {
		ggin.SetMode(ggin.ReleaseMode)
	}
	g.Server = ggin.New()
	g.register()
}

//run http.ListenAndServe
func (g *gin) run(addr ...string)  {
	address := g.resolveAddress(addr)
	g.debugPrint("Listening and serving HTTP on "+address)
	err := gracehttp.ListenAndServe(address, g.Server)
	if err != nil {
		Log.Error(err)
	}
	return
}

//resolveAddress resolve address
func (g *gin) resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			g.debugPrint("Environment variable PORT=\""+port+"\"")
			return ":" + port
		}
		g.debugPrint("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}

//debugPrint print debug info
func (g *gin) debugPrint(format string, values ...interface{}) {
	if Config.Debug {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		_,_ = fmt.Fprintf(os.Stdout, "[GIN-debug] "+format, values...)
	}
}

//register is register router
// For GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, Any, Use requests the respective shortcut
func (g *gin) register() {
	for relativePath, group:=range routerObj.gin  {
		if handle, ok := routerObj.group[relativePath]; ok {
			g.rgroup = g.group(relativePath, handle)
		} else {
			g.rgroup = g.group("/", nil)
		}
		for method, request :=  range group  {
			switch method {
			case "Post":
				g.post(request.relativePath, request.controller)
			case "Get":
				g.get(request.relativePath, request.controller)
			case "Delete":
				g.delete(request.relativePath, request.controller)
			case "Put":
				g.put(request.relativePath, request.controller)
			case "Head":
				g.head(request.relativePath, request.controller)
			case "Patch":
				g.patch(request.relativePath, request.controller)
			case "Options":
				g.options(request.relativePath, request.controller)
			case "Any":
				g.any(request.relativePath, request.controller)
			case "Use":
				g.use(request.controller)
			default:
				g.get(request.relativePath, request.controller)
			}

		}
	}
}

//getController is get controller name
func (g *gin) getController(c ControllerInterface) string {
	reflectVal := reflect.ValueOf(c)
	ct := reflect.Indirect(reflectVal).Type()
	controllerName := strings.TrimSuffix(ct.Name(), "Controller")
	return controllerName
}

//handle is handle for action
func (g *gin) handle(c ControllerInterface, ctx *ggin.Context, actionName string) {
	c.Init(ctx, g.getController(c), actionName)
	c.Prepare()
	g.execute(c, actionName)
	c.Finish()
}

//execute is execute action for controller
func (g *gin) execute(c ControllerInterface, actionName string){
	switch actionName {
	case "Post":
		c.Post()
	case "Get":
		c.Get()
	case "Delete":
		c.Delete()
	case "Put":
		c.Put()
	case "Head":
		c.Head()
	case "Patch":
		c.Patch()
	case "Options":
		c.Options()
	case "Any":
		c.Any()
	case "Group":
		c.Group()
	case "Use":
		c.Use()
	default:
		c.Get()
	}
}

// group creates a new router group. You should add all the routes that have common middlewares or the same path prefix.
// For example, all the routes that use a common middleware for authorization could be grouped.
func (g *gin) group(relativePath string, c ControllerInterface) *ggin.RouterGroup {
	if c == nil {
		return g.Server.Group(relativePath)
	}
	return g.Server.Group(relativePath, func(context *ggin.Context) {
		g.handle(c, context, "Group")
	})
}

// Use adds middleware to the group.
func (g *gin) use(middleware ControllerInterface) *gin {
	g.rgroup.Use(func(context *ggin.Context) {
		g.handle(middleware, context, "Use")
	})
	return g
}

// POST is a shortcut for router.Handle("POST", path, handle).
func (g *gin) post(relativePath string, c ControllerInterface) *gin {
	g.rgroup.POST(relativePath, func(context *ggin.Context) {
		g.handle(c, context, "Post")
	})
	return g
}

// GET is a shortcut for router.Handle("GET", path, handle).
func (g *gin) get(relativePath string, c ControllerInterface) *gin {
	g.rgroup.GET(relativePath, func(context *ggin.Context) {
		g.handle(c, context, "Get")
	})
	return g
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle).
func (g *gin) delete(relativePath string, c ControllerInterface) *gin {
	g.rgroup.DELETE(relativePath, func(context *ggin.Context) {
		g.handle(c, context, "Delete")
	})
	return g
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle).
func (g *gin) patch(relativePath string, c ControllerInterface) *gin {
	g.rgroup.PATCH(relativePath, func(context *ggin.Context) {
		g.handle(c, context, "Patch")
	})
	return g
}

// PUT is a shortcut for router.Handle("PUT", path, handle).
func (g *gin) put(relativePath string, c ControllerInterface) *gin {
	g.rgroup.PUT(relativePath, func(context *ggin.Context) {
		g.handle(c, context, "Put")
	})
	return g
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle).
func (g *gin) options(relativePath string, c ControllerInterface) *gin {
	g.rgroup.OPTIONS(relativePath, func(context *ggin.Context) {
		g.handle(c, context, "Options")
	})
	return g
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle).
func (g *gin) head(relativePath string, c ControllerInterface) *gin {
	g.rgroup.HEAD(relativePath, func(context *ggin.Context) {
		g.handle(c, context, "Head")
	})
	return g
}

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (g *gin) any(relativePath string, c ControllerInterface) *gin {
	g.rgroup.Any(relativePath, func(context *ggin.Context) {
		g.handle(c, context, "Any")
	})
	return g
}