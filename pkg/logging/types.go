package logging

import "github.com/sirupsen/logrus"

type Logger struct {
	Log       *logrus.Logger
	Verbosity int
}
