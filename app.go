package gsigo

import (
	"flag"
	ggin "github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	// Wsi is an application instance
	ENV string
	//system yml config
	Config AppCnf
	//logrus Logger
	Log *logrus.Logger
	//new socketio
	Socketio *socketio
	//new gin
	Gin *gin
	//config file path
	configfile string
)


func init() {
	//get environment variable
	flag.StringVar(&ENV, "env", "production","input environment variable")
	flag.Parse()

	//new gin, socketio, cmd router
	newRouter()
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
	if Config.Mode == ModeDefault {
		Gin.Server.GET("/socket.io/*any", ggin.WrapH(Socketio.Server))
		Gin.Server.POST("/socket.io/*any", ggin.WrapH(Socketio.Server))
	}
	if Config.Static.UrlPath != "" && Config.Static.FilePath != "" {
		Gin.Server.StaticFS(Config.Static.UrlPath, http.Dir(Config.Static.FilePath))
	}
	if Config.Addr == "" {
		Config.Addr = ":8080"
	}
	Gin.run(Config.Addr)
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

//createProject is create a gsigo project
func createProject(){}