package log

import (
	"github.com/sirupsen/logrus"
	"github.com/whf-sky/gsigo/log/hooks"
	"io/ioutil"
	"log/syslog"
	"math/rand"
	"strconv"
	"time"
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
	switch hook {
	case "syslog":
		l.syslogHook(params)
	default:
		l.defaultHook()
	}
}

//developHook develop hook
func (l *log) defaultHook()  {
	hook := hooks.NewDefault()
	l.logrus.Hooks.Add(hook)
}

//syslogHook is syslog hook
func (l *log) syslogHook(params map[string]string)  {
	l.logrus.Out = ioutil.Discard
	var priority syslog.Priority
	switch params["priority"] {
	case "LOG_LOCAL0":
		priority = syslog.LOG_LOCAL0
	case "LOG_LOCAL1":
		priority = syslog.LOG_LOCAL1
	case "LOG_LOCAL2":
		priority = syslog.LOG_LOCAL2
	case "LOG_LOCAL3":
		priority = syslog.LOG_LOCAL3
	case "LOG_LOCAL4":
		priority = syslog.LOG_LOCAL4
	case "LOG_LOCAL5":
		priority = syslog.LOG_LOCAL5
	case "LOG_LOCAL6":
		priority = syslog.LOG_LOCAL6
	case "LOG_LOCAL7":
		priority = syslog.LOG_LOCAL7
	default:
		priority = syslog.LOG_LOCAL0
	}
	hook, err := hooks.NewSyslogHook(params["network"], params["addr"], priority, params["tag"])
	if err == nil {
		l.logrus.Hooks.Add(hook)
	} else {
		l.logrus.Error(err)
	}
}

//GenerateLogid is generate log id
func GenerateLogid() string {
	unixNano := time.Now().UnixNano()
	rand.Seed(unixNano)
	randNum := rand.Int63n(999)
	return strconv.FormatInt(unixNano+randNum,16)
}