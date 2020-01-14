package gsigo

import (
	ggin "github.com/gin-gonic/gin"
	"net/http"
)

// Controller defines some basic http request handler operations, such as
// http context, template and view, session and xsrf.
type Controller struct {
	// context data
	Ctx  *ggin.Context

	// route controller info
	controllerName string
	actionName     string
	groupName	   string
}

// ControllerInterface is an interface to uniform all controller handler.
type ControllerInterface interface {
	Init(ct *ggin.Context, controllerName, actionName string)
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

// Init generates default values of controller operations.
func (c *Controller) Init(ctx *ggin.Context, controllerName, actionName string) {
	c.controllerName = controllerName
	c.actionName = actionName
	c.Ctx = ctx
}

// Prepare runs after Init before request function execution.
func (c *Controller) Prepare() {}

// Finish runs after request function execution.
func (c *Controller) Finish() {}

// Get adds a request function to handle GET request.
func (c *Controller) Get() {
	http.Error(c.Ctx.Writer, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// Post adds a request function to handle POST request.
func (c *Controller) Post() {
	http.Error(c.Ctx.Writer, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// Delete adds a request function to handle DELETE request.
func (c *Controller) Delete() {
	http.Error(c.Ctx.Writer, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// Put adds a request function to handle PUT request.
func (c *Controller) Put() {
	http.Error(c.Ctx.Writer, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// Head adds a request function to handle HEAD request.
func (c *Controller) Head() {
	http.Error(c.Ctx.Writer, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// Patch adds a request function to handle PATCH request.
func (c *Controller) Patch() {
	http.Error(c.Ctx.Writer, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// Options adds a request function to handle OPTIONS request.
func (c *Controller) Options() {
	http.Error(c.Ctx.Writer, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// Any adds a request function to handle all request.
func (c *Controller) Any() {}

// Group adds a request function to handle group request.
func (c *Controller) Group() {}

// Use adds middleware to the group.
func (c *Controller) Use() {}

//GetGroup gets the executing group name.
func (c *Controller) GetGroup() string {
	return c.groupName
}

//GetController gets the executing controller name.
func (c *Controller) GetController() string {
	return c.controllerName
}

//GetAction gets the executing action name.
func (c *Controller) GetAction() string {
	return c.actionName
}