package gsigo


//添加cmd路由信息
//name 控制器名字
//gsigo.Cmd("test", &TestCmd{})
//go run main.go -request_uri='test'
//cmd 命令行控制器
func Cmd(requestUri string, cmd cmdFunc) {
	routerObj.Cmd(requestUri, cmd)
}
