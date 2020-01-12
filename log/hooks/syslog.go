package hooks

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log/syslog"
	"os"
	"runtime"
)

// SyslogHook to send logs via syslog.
type SyslogHook struct {
	Writer        *syslog.Writer
	SyslogNetwork string
	SyslogRaddr   string
}

// Creates a hook to be added to an instance of logger. This is called with
// `hook, err := NewSyslogHook("udp", "localhost:514", syslog.LOG_DEBUG, "")`
// `if err == nil { log.Hooks.Add(hook) }`
func NewSyslogHook(network, raddr string, priority syslog.Priority, tag string) (*SyslogHook, error) {
	w, err := syslog.Dial(network, raddr, priority, tag)
	return &SyslogHook{w, network, raddr}, err
}

func (hook *SyslogHook) Fire(entry *logrus.Entry) error {
	_, file, line, ok := runtime.Caller(7)
	if ok{
		entry.Data["file"] = file
		entry.Data["line"] = line
	}
	text, err := entry.String()
	if text == "" {
		return nil
	}
	if err != nil {
		_,_ = fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	switch entry.Level {
	case logrus.PanicLevel:
		return hook.Writer.Crit(text)
	case logrus.FatalLevel:
		return hook.Writer.Crit(text)
	case logrus.ErrorLevel:
		return hook.Writer.Err(text)
	case logrus.WarnLevel:
		return hook.Writer.Warning(text)
	case logrus.InfoLevel:
		return hook.Writer.Info(text)
	case logrus.DebugLevel, logrus.TraceLevel:
		return hook.Writer.Debug(text)
	default:
		return nil
	}
}

func (hook *SyslogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
