package logging

import "github.com/sirupsen/logrus"

// Logger struct used to set log level across packages
type Logger struct {
	Log       *logrus.Logger
	Verbosity int
}
