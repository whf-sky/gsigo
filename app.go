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
	Config *systemCnf
	Log *logrus.Logger
	Socketio *socketio
	Gin *gin
)

func init() {
	flag.StringVar(&ENV, "env", "production","input environment variable")
	flag.Parse()

	systemYmlParse()
	Log = log.Newlog(
		Config.Log.Formatter,
		Config.Log.Syslog.Network,
		Config.Log.Syslog.Raddr,
		Config.Log.Syslog.Priority,
		Config.Log.Syslog.Tag,
		Config.Debug)
	newRouter()
}

func defaultRun(addr ...string){
	Socketio = newSocketio()
	go Socketio.serve()
	defer Socketio.close()
	ginRun(addr...)
}

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

func cmdRun(){
	for _,cmd := range routerObj.cmd {
		//go func() {
		//	cmd.Init()
		//	cmd.Execute()
		//}()
		cmd.Init()
		cmd.Execute()
	}
	//select {}
}

//createProject is create a gsigo project
func createProject(){}