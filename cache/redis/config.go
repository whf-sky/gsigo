package redis

import (
	"github.com/whf-sky/gsigo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Wsi is an application instance
var configs map[string]GroupYml

// redis config items
type GroupYml struct {
	Address string `yaml:"address"`
	Password string `yaml:"password"`
	ClientName string `yaml:"clientName"`
	Db int `yaml:"db"`
	KeepAlive int `yaml:"keepAlive"`
	MaxIdle int `yaml:"maxIdle"`
	Master MasterYml `yaml:"master"`
	Slave SlaveYml `yaml:"slave"`
}

type MasterYml struct {
	Address string `yaml:"address"`
	MaxIdle int `yaml:"maxIdle"`
}

type SlaveYml struct {
	Address []string `yaml:"address"`
	MaxIdle int `yaml:"maxIdle"`
}

//ymlParse parse redis config
func ymlParse(file string, YmlStruct map[string]GroupYml) map[string]GroupYml {
	data, _ := ioutil.ReadFile(file)
	err := yaml.Unmarshal(data, &YmlStruct)
	if err != nil {
		gsigo.Log.Error(err)
	}
	return YmlStruct
}

func init() {
	configs = ymlParse("./config/"+gsigo.ENV+"/redis.yml", map[string]GroupYml{})
}
