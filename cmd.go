package gsigo

import (
	"flag"
	"fmt"
	"strings"
)

// NewApp 实例化 ggin
func newCmd() *cmd {
	Gcmd = &cmd{}
	return Gcmd
}

type cmdFunc func() CmdInterface

type cmd struct {}

func (c *cmd) run() *cmd {
	var request_uri string
	flag.StringVar(&request_uri, "request_uri", "","Please input request_uri")
	flag.Parse()
	request_uri = strings.TrimSpace(request_uri)
	if cFunc ,ok := routerObj.cmdRouters[request_uri]; ok {
		cmd := cFunc()
		cmd.Init()
		cmd.Prepare()
		cmd.Execute()
		cmd.Finish()
		return c
	}
	fmt.Println("'"+request_uri+"' request_uri not exist ")
	return c
}