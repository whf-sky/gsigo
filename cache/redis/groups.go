package redis

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

//实力化组
func newGroup() *group {
	return &group{}
}

//redis数据组
type group struct {
	//Config 	GroupCnf
	Pool 	*redis.Pool
	Master 	*redis.Pool
	Slave 	[]*redis.Pool
}

//获取redis数据组
func (g *group) dialGroup(config GroupCnf) *group {
	//主从不存在时获取单台redis连接
	if config.Address != "" {
		g.Pool = g.dial(config, config.Address, config.MaxIdle)
		return g
	}

	//打开写的连接
	g.Master = g.dial(config, config.Master.Address, config.Master.MaxIdle)

	//打开读的连接
	for i, addr := range config.Slave.Address  {
		g.Slave[i] = g.dial(config, addr, config.Slave.MaxIdle)
	}
	return g
}

//打开一个新的redis连接
//config 配置信息
//address 连接地址
//maxIdle 连接池的最大链接数
func (g *group) dial(config GroupCnf, address string, maxIdle int) *redis.Pool {
	dial := func() (redis.Conn, error) {
		r, err  := Dial(address,
			redis.DialPassword(config.Password),
			redis.DialDatabase(config.Select),
			redis.DialKeepAlive(time.Second * time.Duration(config.KeepAlive)))
		return r, err
	}
	// 初始化一个新的连接池
	return &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: 180 * time.Second,
		Dial:        dial,
	}
}