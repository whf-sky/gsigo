package gsigo

import (
	"flag"
	"net/http"
	"github.com/sirupsen/logrus"
	"github.com/whf-sky/gsigo/log"
	ggin "github.com/gin-gonic/gin"
)

var (
	// Wsi is an application instance
	ENV string
	//system yml config
	Config *systemCnf
	//logrus Logger
	Log *logrus.Logger
	//new socketio
	Socketio *socketio
	//new gin
	Gin *gin
)


func init() {
	//get environment variable
	flag.StringVar(&ENV, "env", "production","input environment variable")
	flag.Parse()

	//parse system yml config
	systemYmlParse()
	//new log
	Log = log.Newlog(
		Config.Log.Formatter,
		Config.Log.Syslog.Network,
		Config.Log.Syslog.Raddr,
		Config.Log.Syslog.Priority,
		Config.Log.Syslog.Tag,
		Config.Debug)
	//new gin, socketio, cmd router
	newRouter()
}

//defaultRun run socketio, gin server
func defaultRun(addr ...string){
	Socketio = newSocketio()
	go Socketio.serve()
	defer Socketio.close()
	ginRun(addr...)
}

//ginRun run gin server
func ginRun (addr ...string)  {
	Gin = newGin()
	Gin.Server.Use(func(c *ggin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", Config.Gin.Origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")
		if gsigoMode == ModeDefault && c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Request.Header.Del("Origin")
		c.Next()
	})
	if gsigoMode == ModeDefault {
		Gin.Server.GET("/socket.io/*any", ggin.WrapH(Socketio.Server))
		Gin.Server.POST("/socket.io/*any", ggin.WrapH(Socketio.Server))
	}
	Gin.Server.StaticFS(Config.Static.UrlPath, http.Dir(Config.Static.FilePath))
	Gin.run(addr...)
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