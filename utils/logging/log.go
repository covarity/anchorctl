package logging

import (
	"github.com/sirupsen/logrus"
)

var Log = &Logger{}

func init() {
	Log.Log = logrus.New()
}

func (logger *Logger) SetVerbosity(verbosity int) {
	Log.Log.SetLevel(convertLevel(verbosity))
	Log.Verbosity = verbosity
}

func (logger *Logger) GetLogger() *logrus.Logger {
	return Log.Log
}

func convertLevel(verbosity int) logrus.Level {
	logLvl := map[int]logrus.Level{
		1: logrus.PanicLevel,
		2: logrus.FatalLevel,
		3: logrus.ErrorLevel,
		4: logrus.WarnLevel,
		5: logrus.InfoLevel,
		6: logrus.DebugLevel,
		7: logrus.TraceLevel,
	}

	if verbosity < 1 {
		verbosity = 1
	} else if verbosity > 7 {
		verbosity = 7
	}

	return logLvl[verbosity]
}
