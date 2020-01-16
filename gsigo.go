package gsigo

import (
	"flag"
	ggin "github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/whf-sky/gsigo/log"
)

const (
	// VERSION represent gsigo gin socketio framework version.
	VERSION = "1.0.0"
	//mode
	ModeDefault = "default"
	ModeGin = "gin"
	ModeCmd = "cmd"
	ModeInit = "init"
)

var (
	// Wsi is an application instance
	ENV string
	//system yml config
	Config gsigoCnf
	//logrus Logger
	Log *logrus.Logger
	//new socketio
	Socketio *socketio
	//new gin
	Gin *gin
	//config path
	ConfigPath string
)

//run gsigo
func Run(file ...string) {
	loadConfig(file...)
	loadLog()
	switch Config.APP.Mode {
	case ModeDefault:
		defaultRun()
	case ModeGin:
		ginRun()
	case ModeCmd:
		cmdRun()
	case ModeInit:
		//todo 待实现
	default:
		defaultRun()
	}
}

//load log
func loadLog (){
	//new log
	Log = log.Newlog(
		Config.Log.Hook,
		Config.Log.Formatter,
		Config.Log.Params,
		Config.APP.Debug)
}

//defaultRun run socketio, gin server
func defaultRun(){
	Socketio = newSocketio()
	go Socketio.serve()
	defer Socketio.close()
	ginRun()
}

//ginRun run gin server
func ginRun ()  {
	Gin = newGin()
	if Config.APP.Mode == ModeDefault {
		Gin.Server.GET("/socket.io/*any", ggin.WrapH(Socketio.Server))
		Gin.Server.POST("/socket.io/*any", ggin.WrapH(Socketio.Server))
	}
	addr := Config.APP.Host + ":" + Config.APP.Port
	if addr == ":" {
		addr = "0.0.0.0:8080"
	}
	Gin.run(addr)
}

//cmdRun run cmd server
func cmdRun(){
	for _,cmd := range routerObj.cmd {
		go func() {
			cmd.Init()
			cmd.Execute()
		}()
	}
	select {}
}

func init() {
	//get environment variable
	flag.StringVar(&ENV, "env", "production","input environment variable")
	flag.Parse()
	//new gin, socketio, cmd router
	newRouter()
}