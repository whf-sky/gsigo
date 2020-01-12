package mongo

import (
	"github.com/whf-sky/gsigo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Wsi is an application instance
var configs map[string]GroupYml

// redis config items
type GroupYml struct {
	Uri string `yaml:"uri"`
	MaxIdle uint64 `yaml:"maxIdle"`
	Timeout int `yaml:"timeout"`
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
	configs = ymlParse("./config/"+gsigo.ENV+"/mongo.yml", map[string]GroupYml{})
}
