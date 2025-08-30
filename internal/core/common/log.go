package common

import (
	"github.com/sirupsen/logrus"
)

var globalLogger MyLogger

func Log() MyLogger {
	return globalLogger
}

func InitLog(log MyLogger) {
	globalLogger = log
}

type MyLogger interface {
	Info(args ...interface{})
	Infof(fmt string, args ...interface{})
	Debug(args ...interface{})
	Debugf(fmt string, args ...interface{})
	Warn(args ...interface{})
	Warnf(fmt string, args ...interface{})
	Error(args ...interface{})
	Errorf(fmt string, args ...interface{})
}

type DefaultLogger struct {
	MyLogger
}

func (l *DefaultLogger) Info(args ...interface{}) {
	logrus.Info(args)
}

func (l *DefaultLogger) Infof(fmt string, args ...interface{}) {
	logrus.Infof(fmt, args)
}

func (l *DefaultLogger) Debug(args ...interface{}) {
	logrus.Debug(args)
}

func (l *DefaultLogger) Debugf(fmt string, args ...interface{}) {
	logrus.Debugf(fmt, args)
}

func (l *DefaultLogger) Warn(args ...interface{}) {
	logrus.Warn(args)
}

func (l *DefaultLogger) Warnf(fmt string, args ...interface{}) {
	logrus.Warnf(fmt, args)
}

func (l *DefaultLogger) Error(args ...interface{}) {
	logrus.Error(args)
}

func (l *DefaultLogger) Errorf(fmt string, args ...interface{}) {
	logrus.Errorf(fmt, args)
}

func InitDefaultLogger() {
	globalLogger = &DefaultLogger{}
}
