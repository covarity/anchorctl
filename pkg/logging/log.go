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

func (logger *Logger) Info(field, value, info string) {
	Log.Log.WithFields(logrus.Fields{field: value}).Infoln(info)
}

func (logger *Logger) InfoWithFields(fields map[string]interface{}, info string) {
	var field logrus.Fields = fields
	Log.Log.WithFields(field).Infoln(info)
}

func (logger *Logger) Warn(field, value, warn string) {
	Log.Log.WithFields(logrus.Fields{field: value}).Warnln(warn)
}

func (logger *Logger) WarnWithFields(fields map[string]interface{}, warn string) {
	var field logrus.Fields = fields
	Log.Log.WithFields(field).Warnln(warn)
}

func (logger *Logger) Error(err error, error string) {
	Log.Log.WithError(err).Errorln(error)
}

func (logger *Logger) Fatal(err error, fatal string) {
	Log.Log.WithError(err).Fatalln(fatal)
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
