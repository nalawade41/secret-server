package logger

import (
	"github.com/sirupsen/logrus"
)

func Debug(msg ...interface{}) {
	// We Can add sentry here and send the error to sentry with error level
	logrus.Debug(msg...)
}

func Debugf(format string, args ...interface{}) {
	// We Can add sentry here and send the error to sentry with error level
	logrus.Debugf(format, args...)
}

func Info(msg ...interface{}) {
	// We Can add sentry here and send the error to sentry with error level
	logrus.Info(msg...)
}

func Infof(format string, args ...interface{}) {
	// We Can add sentry here and send the error to sentry with error level
	logrus.Infof(format, args...)
}

func Warn(msg ...interface{}) {
	// We Can add sentry here and send the error to sentry with error level
	logrus.Warn(msg...)
}

func Warnf(format string, args ...interface{}) {
	// We Can add sentry here and send the error to sentry with error level
	logrus.Warnf(format, args...)
}

func Error(msg ...interface{}) {
	// We Can add sentry here and send the error to sentry with error level
	logrus.Error(msg...)
}

func Errorf(format string, args ...interface{}) {
	// We Can add sentry here and send the error to sentry with error level
	logrus.Errorf(format, args...)
}
