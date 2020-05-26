package gsigo

import (
	"github.com/whf-sky/gsigo/cache/redis"
	"github.com/whf-sky/gsigo/config"
	"github.com/whf-sky/gsigo/log"
	"github.com/whf-sky/gsigo/orm"
	"github.com/whf-sky/gsigo/utils/fun"
	"os"
)

//钩子函数
type hooksFunc func() error

func newHook()  *hook{
	return &hook{hooks: []hooksFunc{}}
}

type hook struct {
	//钩子
	hooks []hooksFunc
}

//添加钩子
//hfuncs 钩子执行的函数
func(h *hook) add(hfuncs ...hooksFunc) *hook{
	h.hooks = append(h.hooks ,hfuncs...)
	return h
}

//运行某个组的钩子
func(h *hook) run() error {
	for _, hook := range h.hooks {
		err := hook()
		if err != nil {
			return err
		}
	}
	return nil
}

//获取环境变量
func registerGetEnvHook()  error {
	//环境变量中获取环境变量
	ENV = os.Getenv("GSIGO_ENV")
	if ENV == "" {
		ENV = "production"
	}
	return  nil
}

//注册log
func registerLogHook () error {
	Log = log.Newlog(
		Config.Log.Hook,
		Config.Log.Formatter,
		Config.Log.Params,
		Config.APP.Debug,
	)
	return nil
}

//注册 redis hook
func registerRedisHook () error {
	if Config.APP.CloseConfigFile {
		return nil
	}
	cfile := ConfigPath+"/"+ENV+"/"+"redis.ini"
	if !fun.DirExists(cfile) {
		return nil
	}
	groups, err := redis.DialGroup(func(out interface{}) error {
		return config.Read( out, cfile)
	})
	if err != nil {
		return err
	}
	if groups != nil {
		GRedis = redis.NewRedis().SetGroups(groups)
	}
	return nil
}

//注册orm hook
func registerOrmHook () error {
	if Config.APP.CloseConfigFile {
		return nil
	}
	cfile := ConfigPath+"/"+ENV+"/"+"database.ini"
	if !fun.DirExists(cfile) {
		return nil
	}
	groups, err := orm.OpenGroup(func(out interface{}) error {
		return config.Read( out, cfile)
	})
	if err != nil {
		return err
	}
	if groups != nil {
		GOrm = orm.NewDB().SetGroups(groups)
	}
	return nil
}
