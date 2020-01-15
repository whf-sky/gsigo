package gsigo

import "github.com/whf-sky/gsigo/log"

// VERSION represent gsigo gin socketio framework version.
const VERSION = "1.0.0"

const (
	ModeDefault = "default"
	ModeGin = "gin"
	ModeCmd = "cmd"
	ModeInit = "init"
)

//run gsigo
func Run(config ...string) {
	if len(config) > 0 {
		configfile = config[0]
	}
	//parse app yml config
	appYmlParse()

	//init mode
	if Config.Mode == "" {
		Config.Mode = ModeDefault
	}

	if configfile == "" {
		Config.Debug = true
	}

	//new log
	Log = log.Newlog(
		Config.Log.Hook,
		Config.Log.Formatter,
		Config.Log.Params,
		Config.Debug)

	//run
	switch Config.Mode {
	case ModeDefault:
		defaultRun()
	case ModeGin:
		ginRun()
	case ModeCmd:
		cmdRun()
	case ModeInit:
		createProject()
	default:
		defaultRun()
	}
}