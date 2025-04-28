package logger

import "github.com/sirupsen/logrus"

type Logger interface {
	GetLogger() *logrus.Logger
}

func (l *logger) GetLogger() *logrus.Logger {
	return l.log
}

type logger struct {
	log *logrus.Logger
}

func NewLogger() Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetLevel(6)

	return &logger{log: l}
}
