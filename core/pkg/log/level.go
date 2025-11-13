package log

import "github.com/vistormu/go-dsa/ansi"

type Level uint8

const (
	Debug Level = iota
	Info
	Warning
	Error
	Fatal
)

var levelToString = map[Level]string{
	Debug:   "  [debug]",
	Info:    "   [info]",
	Warning: "[warning]",
	Error:   "  [error]",
	Fatal:   "  [fatal]",
}

var levelToColor = map[Level]string{
	Debug:   ansi.Blue,
	Info:    ansi.Green,
	Warning: ansi.Yellow,
	Error:   ansi.Red,
	Fatal:   ansi.BgRed,
}
