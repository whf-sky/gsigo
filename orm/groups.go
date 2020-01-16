package orm

import (
	"github.com/jinzhu/gorm"
	"time"
)

var groups map[string]*group

//get groups
func using(gname ...string) map[string]*group {
	// panic error not database group
	if len(gname) == 0 {
		panic("not database group")
	}

	tGroups := map[string]*group{}
	for _, name := range gname  {
		if _, ok := groups[name]; !ok {
			groups[name] = &group{}
			groups[name].get(name)
		}
		tGroups[name] = groups[name]
	}
	return tGroups
}

//database group
type group struct {
	Config 	dbGroupsCnf
	Db 		*gorm.DB
	Master 	*gorm.DB
	Slave 	[]*gorm.DB
}

//get gorm group
func (g *group) get(name string) *group {
	//get databases group
	var ok bool
	g.Config, ok = configs[name]
	if !ok {
		panic("database configs be short of '"+name+"'")
	}

	//Do not distinguish between master and slave
	if g.Config.Dsn != "" {
		g.Db = g.open(g.Config.Driver,
			g.Config.Dsn,
			g.Config.MaxIdle,
			g.Config.MaxOpen,
			g.Config.MaxLifetime)
		return g
	}

	//get master gorm
	g.Master = g.open(g.Config.Driver,
		g.Config.Master.Dsn,
		g.Config.Master.MaxIdle,
		g.Config.Master.MaxOpen,
		g.Config.MaxLifetime)

	//get slave gorm
	for i, dsn := range g.Config.Slave.Dsn  {
		g.Slave[i] = g.open(g.Config.Driver,
			dsn,
			g.Config.Slave.MaxIdle,
			g.Config.Slave.MaxOpen,
			g.Config.MaxLifetime)
	}
	return g
}

//Open initialize a new db connection
func (g *group) open(driver string, dsn string, maxIdle int, maxOpen int, maxLifetime int) *gorm.DB {
	// Open initialize a new db connection
	db, err := Open(driver, dsn)

	//panic error
	if err != nil {
		panic(err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	db.DB().SetMaxIdleConns(maxIdle)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	db.DB().SetMaxOpenConns(maxOpen)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	db.DB().SetConnMaxLifetime(time.Duration(maxLifetime) * time.Hour)

	return db
}


func init()  {
	groups = map[string]*group{}
}

