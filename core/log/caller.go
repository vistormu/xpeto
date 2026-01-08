package log

import (
	"path/filepath"
	"runtime"
)

type Caller struct {
	File string
	Line int
	Func string
}

func caller(skip int) Caller {
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return Caller{}
	}

	fn := runtime.FuncForPC(pc)
	name := ""
	if fn != nil {
		name = fn.Name()
	}

	return Caller{File: filepath.Base(file), Line: line, Func: name}
}
