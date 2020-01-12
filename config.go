package gsigo

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type systemCnf struct {
	Debug bool `yaml:"Debug"`
	Static staticCnf `yaml:"Static"`
	Gin GinCnf `yaml:"Gin"`
	SocketIo socketioCnf `yaml:"SocketIO"`
	Log logCnf `yaml:"Log"`
}

type GinCnf struct {
	Origin string `yaml:"Origin"`
}

type staticCnf struct {
	UrlPath string `yaml:"UrlPath"`
	FilePath string `yaml:"FilePath"`
}

type socketioCnf struct {
	PingTimeout int `yaml:"PingTimeout"`
	PingInterval int `yaml:"PingInterval"`
}

type logCnf struct {
	Formatter string `yaml:"Formatter"`
	Syslog syslogCnf `yaml:"Syslog"`
}

type syslogCnf struct {
	Network string `yaml:"Network"`
	Raddr string `yaml:"Raddr"`
	Tag string `yaml:"Tag"`
	Priority string `yaml:"Priority"`
}

func systemYmlParse() {
	Config = &systemCnf{}
	data, _ := ioutil.ReadFile("./config/"+ENV+"/system.yml")
	err := yaml.Unmarshal(data, Config)
	if err != nil {
		Log.Error(err)
	}
}