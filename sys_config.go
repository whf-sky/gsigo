package gsigo

import (
	"github.com/whf-sky/gsigo/config"
	"path/filepath"
)

//gsigo config
type gsigoCnf struct {
	//application config
	APP APPCnf `ini:"app"`
	//socket config
	Socket SocketCnf `ini:"socket"`
	//log config
	Log LogCnf `ini:"log"`
}

//application config for gsigo config
type APPCnf struct {
	//application name
	Name string `ini:"name"`
	//server address
	Host string `ini:"host"`
	//server port
	Port string `ini:"port"`
	//gsigo mode
	Mode string `ini:"mode"`
	//is open debug
	//value is true/false
	Debug bool `ini:"debug"`
	//关闭配置文件 默认关闭
	CloseConfigFile bool `ini:"close_config_file"`
}

//socket config for gsigo config
type SocketCnf struct {
	//socketio ping time out
	//Unit second
	PingTimeout int `ini:"ping_timeout"`
	//socketio ping interval
	//Unit second
	PingInterval int `ini:"ping_interval"`
}

//log config for gsigo config
type LogCnf struct {
	//logrus log hook
	//hook: syslog,
	Hook string `ini:"hook"`
	//Formatter json/text
	Formatter string `ini:"formatter"`
	//log params config
	Params map[string]string `ini:"params"`
}

//加载配置信息
func sysConfig() error {
	if sysConfigFile == "" {
		Config.APP.Debug = true
		return nil
	}
	ConfigPath = filepath.Dir(sysConfigFile)
	return config.NewIni().ReadMerge( &Config, sysConfigFile, ENV)
}