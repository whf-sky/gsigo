package orm

import (
	"github.com/jinzhu/gorm"
	"time"
)

//实力化组
func newGroup() *group {
	return &group{}
}

//数据库组
type group struct {
	Db 		*gorm.DB
	Master 	*gorm.DB
	Slave 	[]*gorm.DB
}

//启动一个组的数据库连接
func (g *group) openGroup(config dbGroupsCnf) *group {
	//不存在主从配置时直接打开一个数据库连接
	if config.Dsn != "" {
		g.Db = g.open(config.Driver,
			config.Dsn,
			config.MaxIdle,
			config.MaxOpen,
			config.MaxLifetime)
		return g
	}

	//打开主库的数据库连接
	g.Master = g.open(config.Driver,
		config.Master.Dsn,
		config.Master.MaxIdle,
		config.Master.MaxOpen,
		config.MaxLifetime)

	//打开从库的数据库连接
	for i, dsn := range config.Slave.Dsn  {
		g.Slave[i] = g.open(config.Driver,
			dsn,
			config.Slave.MaxIdle,
			config.Slave.MaxOpen,
			config.MaxLifetime)
	}
	return g
}

//Open 初始化一个数据库连接
func (g *group) open(driver string, dsn string, maxIdle int, maxOpen int, maxLifetime int) *gorm.DB {
	// 打开一个新的数据库连接
	db, err := Open(driver, dsn)
	if err != nil {
		panic(err)
	}

	// 设置空闲连接池中的最大连接数。
	db.DB().SetMaxIdleConns(maxIdle)
	// 设置到数据库的最大打开连接数。
	db.DB().SetMaxOpenConns(maxOpen)
	// 设置可重用连接的最大时间量。
	db.DB().SetConnMaxLifetime(time.Duration(maxLifetime) * time.Hour)

	return db
}

