package log

import (
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// ContextHook ..
type ContextHook struct {
}

// Levels ...
func (hook ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire ...
func (hook ContextHook) Fire(entry *logrus.Entry) error {
	var calldepth int
	switch entry.Level {
	case logrus.ErrorLevel, logrus.WarnLevel:
		calldepth = 6
	default:
		calldepth = 6
	}

	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}

	entry.Data["file"] = file
	entry.Data["line"] = line
	return nil
}
