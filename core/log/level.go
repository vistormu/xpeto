package log

type Level uint8

const (
	Debug Level = iota
	Info
	Warning
	Error
	Fatal
	MaxLevel = Fatal
)

func (l Level) String() string {
	switch l {
	case Debug:
		return "debug"
	case Info:
		return "info"
	case Warning:
		return "warning"
	case Error:
		return "error"
	case Fatal:
		return "fatal"
	default:
		return "unknown"
	}
}
