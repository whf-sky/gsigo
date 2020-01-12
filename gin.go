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

type gin struct {
	Server *ggin.Engine
}

func (g *gin) newServer() {
	if !Config.Debug {
		ggin.SetMode(ggin.ReleaseMode)
	}
	g.Server = ggin.New()
	g.register()
}

func (g *gin) run(addr ...string)  {
	address := g.resolveAddress(addr)
	g.debugPrint("Listening and serving HTTP on "+address)
	err := gracehttp.ListenAndServe(address, g.Server)
	if err != nil {
		Log.Error(err)
	}
	return
}

func (g *gin) register(){
	for relativePath, group:=range routerObj.gin  {
		if handle, ok := routerObj.group[relativePath]; ok {
			g.Group(relativePath, handle)
		}
		for method, request :=  range group  {
			switch method {
			case "Post":
				g.POST(request.relativePath, request.controller)
			case "Get":
				g.GET(request.relativePath, request.controller)
			case "Delete":
				g.DELETE(request.relativePath, request.controller)
			case "Put":
				g.PUT(request.relativePath, request.controller)
			case "Head":
				g.HEAD(request.relativePath, request.controller)
			case "Patch":
				g.PATCH(request.relativePath, request.controller)
			case "Options":
				g.OPTIONS(request.relativePath, request.controller)
			case "Any":
				g.Any(request.relativePath, request.controller)
			case "Use":
				g.Use(request.controller)
			default:
				g.GET(request.relativePath, request.controller)
			}

		}
	}
}

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

func (g *gin) debugPrint(format string, values ...interface{}) {
	if Config.Debug {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		_,_ = fmt.Fprintf(os.Stdout, "[GIN-debug] "+format, values...)
	}
}

func (g *gin) getController(c ControllerInterface) string {
	reflectVal := reflect.ValueOf(c)
	ct := reflect.Indirect(reflectVal).Type()
	controllerName := strings.TrimSuffix(ct.Name(), "Controller")
	return controllerName
}

func (g *gin) funcHandle(c ControllerInterface, ctx *ggin.Context, actionName string) {
	c.Init(ctx, g.getController(c), actionName)
	c.Prepare()
	g.excuteAction(c, actionName)
	c.Finish()
}

func (g *gin) excuteAction(c ControllerInterface, actionName string){
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

// Group creates a new router group. You should add all the routes that have common middlewares or the same path prefix.
// For example, all the routes that use a common middleware for authorization could be grouped.
func (g *gin) Group(relativePath string, c ControllerInterface) *gin {
	if c == nil {
		g.Server.Group(relativePath)
	} else {
		g.Server.Group(relativePath, func(context *ggin.Context) {
			g.funcHandle(c, context, "Group")
		})
	}
	return g
}

// Use adds middleware to the group.
func (g *gin) Use(middleware ControllerInterface) *gin {
	g.Server.Use(func(context *ggin.Context) {
		g.funcHandle(middleware, context, "Use")
	})
	return g
}

// POST is a shortcut for router.Handle("POST", path, handle).
func (g *gin) POST(relativePath string, c ControllerInterface) *gin {
	g.Server.POST(relativePath, func(context *ggin.Context) {
		g.funcHandle(c, context, "Post")
	})
	return g
}

// GET is a shortcut for router.Handle("GET", path, handle).
func (g *gin) GET(relativePath string, c ControllerInterface) *gin {
	g.Server.GET(relativePath, func(context *ggin.Context) {
		g.funcHandle(c, context, "Get")
	})
	return g
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle).
func (g *gin) DELETE(relativePath string, c ControllerInterface) *gin {
	g.Server.DELETE(relativePath, func(context *ggin.Context) {
		g.funcHandle(c, context, "Delete")
	})
	return g
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle).
func (g *gin) PATCH(relativePath string, c ControllerInterface) *gin {
	g.Server.PATCH(relativePath, func(context *ggin.Context) {
		g.funcHandle(c, context, "Patch")
	})
	return g
}

// PUT is a shortcut for router.Handle("PUT", path, handle).
func (g *gin) PUT(relativePath string, c ControllerInterface) *gin {
	g.Server.PUT(relativePath, func(context *ggin.Context) {
		g.funcHandle(c, context, "Put")
	})
	return g
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle).
func (g *gin) OPTIONS(relativePath string, c ControllerInterface) *gin {
	g.Server.OPTIONS(relativePath, func(context *ggin.Context) {
		g.funcHandle(c, context, "Options")
	})
	return g
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle).
func (g *gin) HEAD(relativePath string, c ControllerInterface) *gin {
	g.Server.HEAD(relativePath, func(context *ggin.Context) {
		g.funcHandle(c, context, "Head")
	})
	return g
}

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (g *gin) Any(relativePath string, c ControllerInterface) *gin {
	g.Server.Any(relativePath, func(context *ggin.Context) {
		g.funcHandle(c, context, "Any")
	})
	return g
}