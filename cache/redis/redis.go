package redis

import (
	"github.com/gomodule/redigo/redis"
	"strings"
)

const (
	ModeDefault = iota
	ModeMaster
	ModeSlave
)

var cmds = []string{"TYPE" ,"RANDOMKEY", "DUMP", "TTL", "PTTL", "EXISTS", "KEYS", "GETRANGE", "GET", "GETBIT",
	"STRLEN", "MGET", "LINDEX", "LRANGE", "LLEN", "HMGET", "HGETALL", "HGET","HEXISTS", "HLEN", "HVALS",
	"HKEYS", "SUNION", "SCARD", "SRANDMEMBER", "SMEMBERS", "SINTER", "SISMEMBER", "SDIFF", "SSCAN", "ZREVRANK",
	"ZLEXCOUNT", "ZCARD", "ZRANK", "ZRANGEBYSCORE", "ZRANGEBYLEX", "ZSCORE", "ZSCAN", "ZREVRANGEBYSCORE", "ZREVRANGE",
	"ZRANGE", "ZCOUNT", "PFCOUNT", "GEOHASH", "GEOPOS", "GEODIST", "GEORADIUS", "GEORADIUSBYMEMBER"}

var slaveCmd = map[string]int{}

func init()  {
	for _, cmd := range cmds  {
		slaveCmd[cmd] = ModeSlave
	}
}

// NewRedis create new redis cache with default collection name.
func NewRedis(gname ...string) *Redis {
	return (&Redis{}).Using(gname...)
}
//Redis redis struct
type Redis struct {
	gnames 		[]string
	mode 		int
	curSlave 	int
	pool 		*redis.Pool
	slave 		[]*redis.Pool
	groups 		map[string]*group // redis connection pool
}

//open database group
func (d *Redis) Using(gname ...string) *Redis {
	d.gnames = gname
	d.groups = using(gname...)
	return d
}

func (r *Redis) Master() *Redis{
	r.mode = ModeMaster
	return r
}

func (r *Redis) Slave() *Redis{
	r.mode = ModeSlave
	return r
}

func (r *Redis) Do(cmd string, args ...interface{}) (reply interface{}, err error){
	mode := ModeMaster
	cmd = strings.ToUpper(cmd)
	if _, ok := slaveCmd[cmd]; ok {
		mode = ModeSlave
	}
	c := r.Pool(mode).Get()
	defer c.Close()
	return c.Do(cmd, args...)
}

func (r *Redis) Send(commandName string, args ...interface{}) error{
	c := r.Pool(1).Get()
	defer c.Close()
	return c.Send(commandName, args...)
}

func (r *Redis) Flush() error{
	c := r.Pool(ModeMaster).Get()
	defer c.Close()
	return c.Flush()
}

func (r *Redis) Receive() (reply interface{}, err error) {
	c := r.Pool(ModeMaster).Get()
	defer c.Close()
	return c.Receive()
}

func (r *Redis) PubSub() redis.PubSubConn {
	c := r.Pool(ModeMaster).Get()
	return redis.PubSubConn{Conn: c}
}

func (r *Redis) Script(keyCount int, src string) *script {
	c := r.Pool(ModeMaster).Get()
	s := redis.NewScript(keyCount, src)
	return &script{conn: c, script: s}
}

func (r *Redis) Pool(mode int) *redis.Pool {
	defer func() {
		r.mode = ModeDefault
	}()

	//set mode
	if r.mode != ModeDefault {
		mode = r.mode
	}

	//get group
	group, ok := r.groups[r.gnames[0]]
	if !ok {
		panic("this dao without using group")
	}

	//get gorm Db without set read write database
	if group.Master == nil {
		r.pool = group.Pool
		return r.pool
	}
	//get master gorm Db
	if mode == ModeMaster {
		r.pool = group.Master
		return r.pool
	}

	//get slave  gorm Db
	sCnt := len(group.Slave)
	if r.curSlave > sCnt-1 {
		r.curSlave = 0
	}
	r.pool = group.Slave[r.curSlave]
	r.curSlave++
	return r.pool
}

// Dial connects to the Redis server at the given network and
// address using the specified options.
func Dial(address string, options ...redis.DialOption) (redis.Conn, error) {
	return redis.Dial("tcp", address, options...)
}

type script struct {
	script *redis.Script
	conn redis.Conn
}

func (s *script) Do(keysAndArgs ...interface{}) (interface{}, error) {
	defer s.conn.Close()
	return s.script.Do(s.conn, keysAndArgs...)
}

func (s *script) Hash() string {
	return s.Hash()
}
func (s *script) Load() error {
	defer s.conn.Close()
	return s.script.Load(s.conn)
}
func (s *script) Send(keysAndArgs ...interface{}) error {
	defer s.conn.Close()
	return s.script.Send(s.conn, keysAndArgs...)
}

func (s *script) SendHash(keysAndArgs ...interface{}) error {
	defer s.conn.Close()
	return s.script.SendHash(s.conn, keysAndArgs...)
}