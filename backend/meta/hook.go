package meta

import (
	"github.com/sirupsen/logrus"
	"os"
)

type LogHook struct {
	Path string
}

var Size5m int64 = 1024 * 1024 * 5
var step = 0

func (hook *LogHook) Fire(_ *logrus.Entry) error {
	if step < 1024 {
		step++
		return nil
	} else {
		step = 0
	}
	logCheck(hook.Path)
	return nil
}

func (hook *LogHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.TraceLevel,
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

func logCheck(path string) {
	fi, err := os.Stat(path)
	if err == nil {
		if fi.Size() > Size5m {
			_ = os.Truncate(path, 0)
		}
	}
}
