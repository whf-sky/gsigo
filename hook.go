package gsigo

import (
	"github.com/sirupsen/logrus"
	"github.com/whf-sky/gsigo/log"
)

//RegisterLogHook register a log hook
//example :
//package hooks
//
//import (
//"fmt"
//"github.com/sirupsen/logrus"
//"github.com/whf-sky/gsigo"
//)
//
//type Stdout struct {}
//
//func (hook *Stdout) Fire(entry *logrus.Entry) error {
//	fmt.Println("xxxxxxxxxxxx")
//	return nil
//}
//
//func (hook *Stdout) Levels() []logrus.Level {
//	return logrus.AllLevels
//}
//
//func init()  {
//	gsigo.RegisterLogHook("stdout", func(params map[string]string) logrus.Hook {
//		return &Stdout{}
//	})
//}
func RegisterLogHook(name string, hook func(params map[string]string) logrus.Hook)  {
	if log.Hooks == nil {
		log.Hooks = map[string]func(params map[string]string) logrus.Hook{}
	}
	log.Hooks[name] = hook
}
