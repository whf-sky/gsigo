package hooks

import (
	"github.com/sirupsen/logrus"
	"runtime"
)

// StdoutHook to send logs via syslog.
type Stdout struct {
}

// Creates a hook to be added to an instance of logger. This is called with
// `hook, err := NewSyslogHook("udp", "localhost:514", syslog.LOG_DEBUG, "")`
// `if err == nil { log.Hooks.Add(hook) }`
func DevelopNew() (*Stdout) {
	return &Stdout{}
}

func (hook *Stdout) Fire(entry *logrus.Entry) error {
	_, file, line, ok := runtime.Caller(7)
	if ok{
		entry.Data["file"] = file
		entry.Data["line"] = line
	}
	return nil
}

func (hook *Stdout) Levels() []logrus.Level {
	return logrus.AllLevels
}
