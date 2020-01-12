package gsigo

import (
	"github.com/sirupsen/logrus"
)

// Controller defines some basic http request handler operations, such as
// http context, template and view, session and xsrf.
type Cmd struct {
	//logrus log entry
	Log *logrus.Entry
}

// ControllerInterface is an interface to uniform all controller handler.
type CmdInterface interface {
	Init()
	Execute()
}

// Init generates default values of controller operations.
func (c *Cmd) Init() {
	c.Log = Log.WithFields(logrus.Fields{})
}
