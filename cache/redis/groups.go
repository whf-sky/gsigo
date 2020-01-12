package redis

import (
	"github.com/gomodule/redigo/redis"
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
	Config 	GroupYml
	Pool 	*redis.Pool
	Master 	*redis.Pool
	Slave 	[]*redis.Pool
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
	if g.Config.Address != "" {
		g.Pool = g.dial(g.Config, g.Config.Address, g.Config.MaxIdle)
		return g
	}

	//get master gorm
	g.Master = g.dial(g.Config, g.Config.Master.Address, g.Config.Master.MaxIdle)

	//get slave gorm
	for i, addr := range g.Config.Slave.Address  {
		g.Slave[i] = g.dial(g.Config, addr, g.Config.Slave.MaxIdle)
	}
	return g
}

//Open initialize a new db connection
func (g *group) dial(config GroupYml, address string, maxIdle int) *redis.Pool {
	dial := func() (redis.Conn, error) {
		r, err  := Dial(address,
						redis.DialPassword(config.Password),
						redis.DialDatabase(config.Db),
						redis.DialKeepAlive(time.Second * time.Duration(config.KeepAlive)))
		return r, err
	}
	// initialize a new pool
	return &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: 180 * time.Second,
		Dial:        dial,
	}
}

func init()  {
	groups = map[string]*group{}
}

