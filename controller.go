package gsigo

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Controller 定义了一些基本的http请求处理程序操作，例如
// http上下文，模板和视图，会话和xsrf。
type Controller struct {
	//gin Context
	Ctx  *gin.Context
	//控制器名称
	controllerName string
	//操作名称
	actionName     string
	//组名
	groupName	   string
}

// ControllerInterface是一个统一所有控制器处理程序的接口。
type ControllerInterface interface {
	Init(ct *gin.Context, controllerName, actionName string)
	Prepare()
	Get()
	Post()
	Delete()
	Put()
	Head()
	Patch()
	Options()
	Any()
	Finish()
	Group()
	Use()
}

//Init 初始化控制器操作的默认值。
func (c *Controller) Init(ctx *gin.Context, controllerName, actionName string) {
	c.controllerName = controllerName
	c.actionName = actionName
	c.Ctx = ctx
}

// Prepare 在请求函数执行之前，在Init之后运行。
func (c *Controller) Prepare() {}

// Finish 在请求函数执行后运行。
func (c *Controller) Finish() {}

// Get 添加一个请求函数来处理GET请求。
func (c *Controller) Get() {
	http.Error(c.Ctx.Writer, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// Post 添加一个请求函数来处理POST请求。
func (c *Controller) Post() {
	http.Error(c.Ctx.Writer, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// Delete 添加一个请求函数来处理DELETE请求
func (c *Controller) Delete() {
	http.Error(c.Ctx.Writer, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// Put  添加一个请求函数来处理PUT请求
func (c *Controller) Put() {
	http.Error(c.Ctx.Writer, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// Put  添加一个请求函数来处理 HEAD 请求
func (c *Controller) Head() {
	http.Error(c.Ctx.Writer, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// Patch  添加一个请求函数来处理 PATCH 请求
func (c *Controller) Patch() {
	http.Error(c.Ctx.Writer, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// Options  添加一个请求函数来处理 OPTIONS 请求
func (c *Controller) Options() {
	http.Error(c.Ctx.Writer, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// Any  添加一个请求函数来处理所有请求
func (c *Controller) Any() {}

// Group 添加一个请求函数来处理 Group 请求。
func (c *Controller) Group() {}

// Use 添加一个中间件。
func (c *Controller) Use() {}

//GetGroup 获取执行的组名称。
func (c *Controller) GetGroup() string {
	return c.groupName
}

//GetController 获取执行的控制器名称。
func (c *Controller) GetController() string {
	return c.controllerName
}

//GetAction 获取执行的操作名称。
func (c *Controller) GetAction() string {
	return c.actionName
}