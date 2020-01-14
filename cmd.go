package gsigo

// Controller defines some basic http request handler operations, such as
// http context, template and view, session and xsrf.
type Cmd struct {}

// ControllerInterface is an interface to uniform all controller handler.
type CmdInterface interface {
	Init()
	Execute()
}

// Init generates default values of controller operations.
func (c *Cmd) Init() {}
