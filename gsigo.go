package gsigo

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/whf-sky/gsigo/cache/redis"
	"github.com/whf-sky/gsigo/orm"
	"os"
)

const (
	//框架版本信息
	VERSION = "1.0.0"
	//模式: socketio/gin/cmd
	ModeGin = "gin"
	//socketio+gin
	ModeSocketio = "socketio"
	//cmd
	ModeCmd = "cmd"
)

var (
	//环境变量
	ENV string
	//配置文件路径
	ConfigPath string
	//系统配置文件路径
	sysConfigFile string
	//系统配置信息
	Config gsigoCnf
	//logrus 日志记录器
	Log *logrus.Logger
	//gsocketio 实例
	Gsocketio *gsocketio
	//gin 实例
	Ggin *ggin
	//cmd 实例
	Gcmd *cmd
	//Gsigo Orm
	GOrm *orm.DB
	//gsigo redis
	GRedis *redis.Redis
)

//Version 打印版本信息
func Version () {
	fmt.Println(VERSION)
}

//运行前初始化钩子
func initBeforeRun() {
	newHook().add(
		sysConfig,
		registerLogHook,
		registerRedisHook,
		registerOrmHook,
	).run()
}

//Run 运行
func Run(file ...string) {
	if len(file) > 0 {
		sysConfigFile = file[0]
	}
	//初始化加载项
	initBeforeRun();
	//模式
	switch Config.APP.Mode {
	case ModeSocketio:
		newGsocketio().run()
	case ModeCmd:
		newCmd().run()
	case ModeGin:
		newGgin().run()
	default:
		newGgin().run()
	}
}

func init()  {
	//环境变量中获取环境变量
	ENV = os.Getenv("GSIGO_ENV")
	if ENV == "" {
		ENV = "production"
	}
}