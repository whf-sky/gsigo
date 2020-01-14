package orm

import (
	"github.com/whf-sky/gsigo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var configs map[string]dbGroupsYml

type dbGroupsYml struct {
	Driver 		string  	`yaml:"driver"`
	Dsn 		string  	`yaml:"dsn"`
	MaxOpen 	int 		`yaml:"maxOpen"`
	MaxIdle 	int 		`yaml:"maxIdle"`
	//unit hour
	MaxLifetime int 		`yaml:"maxLifetime"`
	Master 		MasterYml 	`yaml:"master"`
	Slave  		SlaveYml 	`yaml:"slave"`
}

type MasterYml struct {
	Dsn 	string  `yaml:"dsn"`
	MaxOpen int 	`yaml:"maxOpen"`
	MaxIdle int 	`yaml:"maxIdle"`
}

type SlaveYml struct {
	Dsn 	[]string  	`yaml:"dsn"`
	MaxOpen int 		`yaml:"maxOpen"`
	MaxIdle int 		`yaml:"maxIdle"`
}

func ymlParse(file string, YmlStruct map[string]dbGroupsYml) map[string]dbGroupsYml {
	data, _ := ioutil.ReadFile(file)
	err := yaml.Unmarshal(data, &YmlStruct)
	if err != nil {
		gsigo.Log.Error(err)
	}
	return YmlStruct
}

func init() {
	configs = ymlParse("./config/"+gsigo.ENV+"/database.yml", map[string]dbGroupsYml{})
}
