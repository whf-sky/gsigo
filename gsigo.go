package gsigo

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/whf-sky/gsigo/log"
	"os"
	"strings"
)

const (
	//框架版本信息
	VERSION = "1.0.0"

	//模式: default/gin
	//socketio+gin
	ModeDefault = "default"
	//gin
	ModeGin = "gin"
	//cmd
	ModeCmd = "cmd"
)

var (
	//环境变量
	ENV string
	//系统配置信息
	Config gsigoCnf
	//logrus 日志记录器
	Log *logrus.Logger
	//gsocketio 实例
	Gsocketio *gsocketio
	//gin 实例
	Ggin *ggin
	//配置文件路径
	ConfigPath string
	//命令行请求URI
	CmdRequestUri string
)

//Run 运行
func Run(file ...string) {
	loadConfig(file...)
	flagParse(Config.APP.Mode)
	registerLog()
	switch Config.APP.Mode {
	case ModeDefault:
		defaultRun()
	case ModeGin:
		ginRun()
	case ModeCmd:
		newCmd().run()
	default:
		defaultRun()
	}
}

//load 注册日志差价
func registerLog () {
	Log = log.Newlog(
		Config.Log.Hook,
		Config.Log.Formatter,
		Config.Log.Params,
		Config.APP.Debug)
}

//defaultRun 运行socketio和gin 服务
func defaultRun() {
	Gsocketio = newGsocketio()
	go Gsocketio.serve()
	defer Gsocketio.close()
	ginRun()
}

//ginRun 运行gin服务
func ginRun () {
	Ggin = newGgin()
	if Config.APP.Mode == ModeDefault {
		Ggin.Engine.GET("/socket.io/*any", gin.WrapH(Gsocketio.Server))
		Ggin.Engine.POST("/socket.io/*any", gin.WrapH(Gsocketio.Server))
	}
	addr := Config.APP.Host + ":" + Config.APP.Port
	if addr == ":" {
		addr = "0.0.0.0:8080"
	}
	Ggin.run(addr)
}

//Version 打印版本信息
func Version () {
	fmt.Println(VERSION)
}

//命令行参数中获取变量信息
//mode 模式 default/gin/cmd
func flagParse(mode string)  {
	parse := false
	if ENV == "" {
		parse = true
		flag.StringVar(&ENV, "env", "production","input environment variable")
	}
	if mode == ModeCmd {
		parse = true
		flag.StringVar(&CmdRequestUri, "request_uri", "","Please input request_uri")
	}
	if parse {
		flag.Parse()
		ENV = strings.TrimSpace(ENV)
		CmdRequestUri = strings.TrimSpace(CmdRequestUri)
	}
}

func init() {
	//环境变量中获取环境变量
	ENV = os.Getenv("GSIGO_ENV")
}