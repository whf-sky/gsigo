package redis

import (
	"github.com/whf-sky/gsigo"
	"github.com/whf-sky/gsigo/config"
)

// Wsi is an application instance
var configs map[string]GroupCnf

// redis config items
type GroupCnf struct {
	//host:port
	Address string 		`ini:"address"`
	//password
	Password string 	`ini:"password"`
	Select int 			`ini:"select"`
	//unit hour
	KeepAlive int 		`ini:"keep_alive"`
	MaxIdle int			`ini:"max_idle"`
	Master MasterCnf 	`ini:"master"`
	Slave SlaveCnf 		`ini:"slave"`
}

type MasterCnf struct {
	//host:port
	Address string 	`ini:"address"`
	MaxIdle int 	`ini:"max_idle"`
}

type SlaveCnf struct {
	//host:port
	Address []string 	`ini:"address"`
	MaxIdle int 		`ini:"max_idle"`
}

func init() {
	err := config.NewIni().Read( &configs, "./configs/"+gsigo.ENV+"/redis.ini")
	if err != nil {
		panic(err)
	}
}
