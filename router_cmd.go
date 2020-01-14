package gsigo
//CmdRouter add cmd handler
func CmdRouter(cmd CmdInterface) *router{
	routerObj.cmd = append(routerObj.cmd, cmd)
	return routerObj
}