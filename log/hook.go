package log

import (
	"github.com/sirupsen/logrus"
)

//RegisterLogHook 注册一个日志钩子
func RegisterLogHook(name string, hook func(params map[string]string) logrus.Hook)  {
	if Hooks == nil {
		Hooks = map[string]func(params map[string]string) logrus.Hook{}
	}
	Hooks[name] = hook
}
