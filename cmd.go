package gsigo

import (
	"fmt"
)

// NewApp 实例化 ggin
func newCmd() *cmd {
	return &cmd{}
}

type cmd struct {}

func (c *cmd) run()  {
	if CmdRequestUri == "" {
		fmt.Println("Please input request_uri ")
		return
	}
	if cmd ,ok := routerObj.cmdRouters[CmdRequestUri]; ok {
		cmd.Init()
		cmd.Prepare()
		cmd.Execute()
		cmd.Finish()
		return
	}
	fmt.Println("'"+CmdRequestUri+"' request_uri not exist ")
}