package orm

import (
	"github.com/whf-sky/gsigo"
	"github.com/whf-sky/gsigo/config"
)

var configs map[string]dbGroupsCnf

type dbGroupsCnf struct {
	Driver 		string  	`ini:"driver"`
	Dsn 		string  	`ini:"dsn"`
	MaxOpen 	int 		`ini:"max_open"`
	MaxIdle 	int 		`ini:"max_idle"`
	//unit hour
	MaxLifetime int 		`ini:"max_lifetime"`
	Master 		MasterCnf 	`ini:"master"`
	Slave  		SlaveCnf 	`ini:"slave"`
}

type MasterCnf struct {
	Dsn 	string  `ini:"dsn"`
	MaxOpen int 	`ini:"max_open"`
	MaxIdle int 	`ini:"max_idle"`
}

type SlaveCnf struct {
	Dsn 	[]string  	`ini:"dsn"`
	MaxOpen int 		`ini:"max_open"`
	MaxIdle int 		`ini:"max_idle"`
}

func init() {
	err := config.NewIni().Read( &configs, "./config/"+gsigo.ENV+"/database.ini")
	if err != nil {
		panic(err)
	}
}
