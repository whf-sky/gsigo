package log

import (
	"github.com/sirupsen/logrus"
	"github.com/whf-sky/gsigo/log/hooks"
	"io/ioutil"
	"log/syslog"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

type log struct {
	logrus *logrus.Logger
}

func Newlog(formatter string, network string, raddr string, priority string, tag string, debug bool) *logrus.Logger {
	log := &log{}
	log.newLogrus(formatter, debug)
	log.hook(network, raddr, priority, tag)
	return log.logrus
}

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

func (l *log) hook(network string, raddr string, priority string, tag string)  {
	if runtime.GOOS == "linux" {
		l.syslogHook(network, raddr, priority, tag)
	} else {
		l.developHook()
	}
}

func (l *log) developHook()  {
	hook := hooks.DevelopNew()
	l.logrus.Hooks.Add(hook)
}

func (l *log) syslogHook(network string, raddr string, priorityStr string, tag string)  {
	l.logrus.Out = ioutil.Discard
	var priority syslog.Priority
	switch priorityStr {
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
	hook, err := hooks.NewSyslogHook(network, raddr, priority, tag)
	if err == nil {
		l.logrus.Hooks.Add(hook)
	} else {
		l.logrus.Error(err)
	}
}

func GenerateLogid() string {
	unixNano := time.Now().UnixNano()
	rand.Seed(unixNano)
	randNum := rand.Int63n(999)
	return strconv.FormatInt(unixNano+randNum,16)
}