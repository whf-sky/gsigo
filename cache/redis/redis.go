package redis

import (
	"github.com/gomodule/redigo/redis"
)

//启动所有的组的redis连接
func DialGroup(configRead func(out interface{}) error) ( map[string]*group, error) {
	var configs map[string]GroupCnf
	err := configRead(&configs)
	if err != nil {
		return nil, err
	}
	if len(configs) == 0 {
		return nil, nil
	}
	groups  := map[string]*group{}
	//启动组的数据连接
	for gname, config := range configs {
		groups[gname] = newGroup().dialGroup(config)
	}
	return groups, nil
}

// Dial connects to the Redis server at the given network and
// address using the specified options.
func Dial(address string, options ...redis.DialOption) (redis.Conn, error) {
	return redis.Dial("tcp", address, options...)
}
