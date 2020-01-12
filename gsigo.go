package gsigo

// VERSION represent gsigo gin socketio framework version.
const VERSION = "1.0.0"

const (
	ModeDefault = iota//gin+socketio
	ModeGin
	ModeCmd
	ModeInit
)

var gsigoMode = ModeDefault

func SetMode(mode int)  {
	gsigoMode = mode
}

func Run(addr ...string) {
	switch gsigoMode {
	case ModeCmd:
		cmdRun()
	case ModeGin:
		ginRun(addr...)
	case ModeDefault:
		defaultRun(addr...)
	case ModeInit:
		createProject()
	default:
		defaultRun(addr...)
	}
}

