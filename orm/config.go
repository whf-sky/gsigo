package orm

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