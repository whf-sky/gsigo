package gsigo

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type systemCnf struct {
	//is open debug
	//value is true/false
	Debug bool `yaml:"Debug"`
	//web static dir
	Static staticCnf `yaml:"Static"`
	//gin config
	Gin GinCnf `yaml:"Gin"`
	//socketio config
	SocketIo socketioCnf `yaml:"SocketIO"`
	//log config
	Log logCnf `yaml:"Log"`
}

//gin config
type GinCnf struct {
	//Access-Control-Allow-Origin
	Origin string `yaml:"Origin"`
}

//static config
type staticCnf struct {
	//web url path
	UrlPath string `yaml:"UrlPath"`
	//file dir
	FilePath string `yaml:"FilePath"`
}

//socketio config
type socketioCnf struct {
	//ping time out
	//Unit second
	PingTimeout int `yaml:"PingTimeout"`
	//ping interval
	//Unit second
	PingInterval int `yaml:"PingInterval"`
}

//log config
type logCnf struct {
	//Formatter json/text
	Formatter string `yaml:"Formatter"`
	//syslog config
	Syslog syslogCnf `yaml:"Syslog"`
}

//syslogCnf syslog config
type syslogCnf struct {
	//syslog newwork
	Network string `yaml:"Network"`
	//syslog addr
	Raddr string `yaml:"Raddr"`
	//syslog tag
	Tag string `yaml:"Tag"`
	//syslog Priority
	Priority string `yaml:"Priority"`
}

//systemYmlParse parse system yml config
func systemYmlParse() {
	Config = &systemCnf{}
	data, _ := ioutil.ReadFile("./config/"+ENV+"/system.yml")
	err := yaml.Unmarshal(data, Config)
	if err != nil {
		Log.Error(err)
	}
}