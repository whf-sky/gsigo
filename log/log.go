package log

import (
	"github.com/sirupsen/logrus"
	"github.com/whf-sky/gsigo/log/hooks"
)

//log hooks
var Hooks map[string]func(params map[string]string)logrus.Hook

type log struct {
	logrus *logrus.Logger
}

//Newlog new log
func Newlog(hook string,formatter string, params map[string]string, debug bool) *logrus.Logger {
	log := &log{}
	log.newLogrus(formatter, debug)
	log.hook(hook, params)
	return log.logrus
}

//newLogrus new Logrus
func (l *log) newLogrus(Formatter string, Debug bool)  {
	l.logrus = logrus.New()
	switch Formatter {
	case "json":
		l.logrus.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		l.logrus.SetFormatter(&logrus.TextFormatter{})
	}
	if Debug {
		l.logrus.SetLevel(logrus.TraceLevel)
	}
}

//hook
func (l *log) hook(hook string, params map[string]string)  {
	if fun, ok := Hooks[hook]; ok {
		l.logrus.Hooks.Add(fun(params))
		return
	}
	l.logrus.Hooks.Add(hooks.NewDefault())
}