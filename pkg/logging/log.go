package logging

import (
	"github.com/sirupsen/logrus"
)

// Log variable initialises Top level logging object
var Log = &Logger{}

// Initialises the logging tool
func init() {
	Log.Log = logrus.New()
}

// SetVerbosity sets the verbosity level of the logging tool
func (logger *Logger) SetVerbosity(verbosity int) {
	Log.Log.SetLevel(convertLevel(verbosity))
	Log.Verbosity = verbosity
}

// GetLogger gets the logger object
func (logger *Logger) GetLogger() *logrus.Logger {
	return Log.Log
}

// Info outputs info log with single field and value
func (logger *Logger) Info(field, value, info string) {
	Log.Log.WithFields(logrus.Fields{field: value}).Infoln(info)
}

// InfoWithFields outputs info log with multiple field and values
func (logger *Logger) InfoWithFields(fields map[string]interface{}, info string) {
	var field logrus.Fields = fields
	Log.Log.WithFields(field).Infoln(info)
}

// Warn outputs warning with single field and value
func (logger *Logger) Warn(field, value, warn string) {
	Log.Log.WithFields(logrus.Fields{field: value}).Warnln(warn)
}

// WarnWithFields outputs warning log with multiple field and values
func (logger *Logger) WarnWithFields(fields map[string]interface{}, warn string) {
	var field logrus.Fields = fields
	Log.Log.WithFields(field).Warnln(warn)
}

// Error outputs error with err object
func (logger *Logger) Error(err error, error string) {
	Log.Log.WithError(err).Errorln(error)
}

// Fatal outputs error and os.Exits(1)
func (logger *Logger) Fatal(err error, fatal string) {
	Log.Log.WithError(err).Fatalln(fatal)
}

// convertLevel takes in a verbosity level and set it on the logging tool
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
