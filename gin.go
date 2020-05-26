package gsigo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"net/http"
	"os"
	"reflect"
	"strings"
)

// NewApp 实例化 ggin
func newGgin() *ggin {
	Ggin = &ggin{}
	return Ggin.server()
}

//gin
type ggin struct {
	Engine *gin.Engine
	RouterGroup *gin.RouterGroup
}

//newServer create gin
func (g *ggin) server() *ggin {
	//设置gin模式信息
	if !Config.APP.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	//获取gin的引擎
	g.Engine = gin.New()
	//注册路由信息
	g.registerRouter()
	return g
}

//socketio路由规则
func (g *ggin) socketioRouter(server *socketio.Server) *ggin {
	g.Engine.GET("/socket.io/*any", gin.WrapH(server))
	g.Engine.POST("/socket.io/*any", gin.WrapH(server))
	return g
}

//run 使用可平滑重启方式运行
func (g *ggin) run() *ggin {
	addr := []string{}
	if Config.APP.Host == "" {
		addr = []string{":" + Config.APP.Port}
	}
	if Config.APP.Host != "" && Config.APP.Port != ""{
		addr = []string{Config.APP.Host + ":" + Config.APP.Port}
	}
	address := g.resolveAddress(addr)
	g.debugPrint("Listening and serving HTTP on "+address)
	err := http.ListenAndServe(address, g.Engine)
	if err != nil {
		Log.Error(err)
	}
	return g
}

//resolveAddress 分析地址
func (g *ggin) resolveAddress(addr []string) string {
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

//debugPrint 打印调试信息
func (g *ggin) debugPrint(format string, values ...interface{}) {
	if Config.APP.Debug {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		_,_ = fmt.Fprintf(os.Stdout, "[GIN-debug] "+format, values...)
	}
}

//registerRouter 注册gin路由信息
// For GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, Any, Use requests the respective shortcut
func (g *ggin) registerRouter() {
	for groupRelativePath, group := range routerObj.ginRouters  {
		//注册路由组
		if handle, ok := routerObj.groupRouters[groupRelativePath]; ok {
			g.RouterGroup = g.group(groupRelativePath, handle)
		} else {
			g.RouterGroup = g.group("/", nil)
		}

		//注册路由信息
		for _, request :=  range group  {
			switch request.method {
			case "POST":
				g.post(request.relativePath, request.controller)
			case "GET":
				g.get(request.relativePath, request.controller)
			case "DELETE":
				g.delete(request.relativePath, request.controller)
			case "PUT":
				g.put(request.relativePath, request.controller)
			case "HEAD":
				g.head(request.relativePath, request.controller)
			case "PATCH":
				g.patch(request.relativePath, request.controller)
			case "OPTIONS":
				g.options(request.relativePath, request.controller)
			case "ANY":
				g.any(request.relativePath, request.controller)
			case "USE":
				g.use(request.controller)
			case "STATIC":
				g.Static(request.relativePath, request.filePath)
			default:
				g.get(request.relativePath, request.controller)
			}

		}
	}
}

//getController 获取控制器名称
func (g *ggin) getController(c ControllerInterface) string {
	reflectVal := reflect.ValueOf(c)
	ct := reflect.Indirect(reflectVal).Type()
	controllerName := strings.TrimSuffix(ct.Name(), "Controller")
	return controllerName
}

//handle is handle for action
func (g *ggin) handle(c ControllerInterface, ctx *gin.Context, actionName string) {
	c.Init(ctx, g.getController(c), actionName)
	c.Prepare()
	g.execute(c, actionName)
	c.Finish()
}

//execute 控制器的执行函数
func (g *ggin) execute(c ControllerInterface, actionName string){
	var actions = map[string]func() {
		"Post":c.Post,
		"Get":c.Get,
		"Delete":c.Delete,
		"Put":c.Put,
		"Head":c.Head,
		"Patch":c.Patch,
		"Options":c.Options,
		"Any":c.Any,
		"Group":c.Group,
		"Use":c.Use,
	}
	if action, ok := actions[actionName]; ok {
		action()
	} else {
		c.Get()
	}
}

//组创建一个新的路由器组。您应该添加所有具有公共中间件或相同路径前缀的路由。
//例如，所有使用公共中间件进行授权的路由都可以分组。
//relativePath 网站的相对路径
//c 组控制器
func (g *ggin) group(relativePath string, c ControllerInterface) *gin.RouterGroup {
	if c == nil {
		return g.Engine.Group(relativePath)
	}
	return g.Engine.Group(relativePath, func(context *gin.Context) {
		g.handle(c, context, "Group")
	})
}

// Use 在组中添加一个中间件
//middleware 中间件控制器
func (g *ggin) use(middleware ControllerInterface) *ggin {
	g.Engine.Use(func(context *gin.Context) {
		g.handle(middleware, context, "Use")
	})
	return g
}

// POST 添加 POST 路由信息
//relativePath 网站的相对路径
//c 中间件控制器
func (g *ggin) post(relativePath string, c ControllerInterface) *ggin {
	g.RouterGroup.POST(relativePath, func(context *gin.Context) {
		g.handle(c, context, "Post")
	})
	return g
}

// GET 添加 GET 路由信息
//relativePath 网站的相对路径
//c 中间件控制器
func (g *ggin) get(relativePath string, c ControllerInterface) *ggin {
	g.RouterGroup.GET(relativePath, func(context *gin.Context) {
		g.handle(c, context, "Get")
	})
	return g
}

// DELETE 添加 DELETE 路由信息
//relativePath 网站的相对路径
//c 中间件控制器
func (g *ggin) delete(relativePath string, c ControllerInterface) *ggin {
	g.RouterGroup.DELETE(relativePath, func(context *gin.Context) {
		g.handle(c, context, "Delete")
	})
	return g
}

// PATCH 添加 PATCH 路由信息
//relativePath 网站的相对路径
//c 中间件控制器
func (g *ggin) patch(relativePath string, c ControllerInterface) *ggin {
	g.RouterGroup.PATCH(relativePath, func(context *gin.Context) {
		g.handle(c, context, "Patch")
	})
	return g
}

// PUT 添加 PUT 路由信息
//relativePath 网站的相对路径
//c 中间件控制器
func (g *ggin) put(relativePath string, c ControllerInterface) *ggin {
	g.RouterGroup.PUT(relativePath, func(context *gin.Context) {
		g.handle(c, context, "Put")
	})
	return g
}

// OPTIONS 添加 OPTIONS 路由信息
//relativePath 网站的相对路径
//c 中间件控制器
func (g *ggin) options(relativePath string, c ControllerInterface) *ggin {
	g.RouterGroup.OPTIONS(relativePath, func(context *gin.Context) {
		g.handle(c, context, "Options")
	})
	return g
}

// HEAD 添加 HEAD 路由信息
//relativePath 网站的相对路径
//c 中间件控制器
func (g *ggin) head(relativePath string, c ControllerInterface) *ggin {
	g.RouterGroup.HEAD(relativePath, func(context *gin.Context) {
		g.handle(c, context, "Head")
	})
	return g
}

// Any 注册一个匹配所有HTTP方法的路由。
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
//relativePath 网站的相对路径
//c 控制器方法
func (g *ggin) any(relativePath string, c ControllerInterface) *ggin {
	g.RouterGroup.Any(relativePath, func(context *gin.Context) {
		g.handle(c, context, "Any")
	})
	return g
}

// Static 添加静态资源路由
//relativePath 网站的相对路径
//filePath 文件的路径
func (g *ggin) Static(relativePath string, filePath string) *ggin {
	g.RouterGroup.StaticFS(relativePath, http.Dir(filePath))
	return g
}