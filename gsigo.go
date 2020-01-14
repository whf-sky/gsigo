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

//SetMode set mode
//For ModeDefault, ModeGin, ModeCmd, ModeInit
func SetMode(mode int)  {
	gsigoMode = mode
}

//run gsigo
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