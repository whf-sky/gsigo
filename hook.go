package gsigo

import (
	"github.com/sirupsen/logrus"
	"github.com/whf-sky/gsigo/log"
)

//RegisterLogHook register a log hook
func RegisterLogHook(name string, hook func(params map[string]string) logrus.Hook)  {
	if log.Hooks == nil {
		log.Hooks = map[string]func(params map[string]string) logrus.Hook{}
	}
	log.Hooks[name] = hook
}
