package gsigo

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type AppCnf struct
{
	//is open debug
	//value is true/false
	Debug bool `yaml:"debug"`
	//gsigo mode
	Mode string `yaml:"mode"`
	//server address
	//default port: 8080
	//example :
	//:8080, 0.0.0.0:8080
	Addr string `yaml:"addr"`
	//socketio ping time out
	//Unit second
	PingTimeout int `yaml:"pingTimeout"`
	//socketio ping interval
	//Unit second
	PingInterval int `yaml:"pingInterval"`
	//log config
	Log logCnf `yaml:"log"`
	//web static dir
	Static staticCnf `yaml:"static"`
}

//static config
type staticCnf struct {
	//web url path
	UrlPath string `yaml:"urlpath"`
	//file dir
	FilePath string `yaml:"filepath"`
}

//log config
type logCnf struct {
	//logrus log hook
	//hook: syslog,
	Hook  string `yaml:"hook"`
	//Formatter json/text
	Formatter string `yaml:"formatter"`
	//log params config
	Params map[string]string `yaml:"params"`
}

//appYmlParse parse application yml config
func appYmlParse() {
	if configfile == "" {
		return
	}
	var configs map[string]AppCnf
	data, _ := ioutil.ReadFile(configfile)
	err := yaml.Unmarshal(data, &configs)
	if err != nil {
		Log.Error(err)
	}

	//Common environment variable
	if cnf, ok := configs["gsigo"];ok {
		Config = cnf
	}

	//environment variable for ENV
	cnf, ok := configs[ENV];
	if !ok {
		return
	}

	if cnf.Debug == true {
		Config.Debug = true
	}

	if cnf.Addr != "" {
		Config.Addr = cnf.Addr
	}

	if cnf.Mode != "" {
		Config.Mode = cnf.Mode
	}

	if cnf.PingTimeout != 0 {
		Config.PingTimeout = cnf.PingTimeout
	}

	if cnf.PingInterval != 0 {
		Config.PingInterval = cnf.PingInterval
	}

	if cnf.Static.FilePath != "" &&  cnf.Static.UrlPath != ""{
		Config.Static.FilePath = cnf.Static.FilePath
	}

	if cnf.Log.Hook != ""{
		Config.Log.Hook = cnf.Log.Hook
	}

	if cnf.Log.Formatter != "" {
		Config.Static.FilePath = cnf.Static.FilePath
	}

	if cnf.Log.Params != nil {
		Config.Log.Params = cnf.Log.Params
	}
}