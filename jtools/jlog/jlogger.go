package jlog

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	OutLogConsole = OutLogType("console")
	OutLogFile    = OutLogType("file")
)

type OutLogType string

type LogEntry struct {
	logConfig ConfigLogYAML
	fields    logrus.Fields

	loggers []*logrus.Logger
}

func new(config ConfigLogYAML) *LogEntry {

	item := LogEntry{}
	item.logConfig = config
	item.loggers = []*logrus.Logger{}

	if config.File.IsOutLog {
		filePath := config.File.makeFilePath()
		filePtr, err := os.OpenFile(
			filePath,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			0666,
		)
		if err != nil {
			Panic(fmt.Errorf("NewLogEntry File Error : " + err.Error()))
		}

		logLevel, err := logrus.ParseLevel(config.File.LogLevel)
		if err != nil {
			Panic(fmt.Errorf("NewLogEntry File logLevel Error : " + err.Error()))
		}

		format := &Formatter{}
		format.TimestampFormat = "2006-01-02T15:04:05.999-UTC" //"2006-01-02T15:04:05.999+0900"
		format.NoColors = true
		format.NoFieldsColors = true

		logger := &logrus.Logger{
			Out:          filePtr,
			Formatter:    format,
			Level:        logLevel,
			ExitFunc:     os.Exit,
			ReportCaller: false,
		}
		item.loggers = append(item.loggers, logger)
	}

	if config.Console.IsOutLog {
		logLevel, err := logrus.ParseLevel(config.Console.LogLevel)
		if err != nil {
			Panic(fmt.Errorf("NewLogEntry Console logLevel Error : " + err.Error()))
		}

		format := &Formatter{
			IsConsole: true,
		}
		format.TimestampFormat = "2006-01-02T15:04:05.999-UTC" //"2006-01-02T15:04:05.999+0900"
		//format.NoFieldsColors = true

		logger := &logrus.Logger{
			Out:          os.Stdout,
			Formatter:    format,
			Level:        logLevel,
			ExitFunc:     os.Exit,
			ReportCaller: false,
		}
		item.loggers = append(item.loggers, logger)
	}

	item.fields = logrus.Fields{
		"Server":  config.FieldBase.ServerName,
		"Version": config.FieldBase.ServerVersion,
	}

	return &item
}

func logNowTime() time.Time {
	return time.Now().UTC()
}
func logFormat(logger *logrus.Logger, t time.Time, fields logrus.Fields) *logrus.Entry {
	return logger.WithTime(t).WithFields(fields)
}

func (my LogEntry) Panic(args ...interface{}) {
	t := logNowTime()
	for i, logger := range my.loggers {
		if i < len(my.loggers)-1 {
			func() {
				defer func() { recover() }()
				logFormat(logger, t, my.fields).Panic(args, string(debug.Stack()))
			}()
		} else {
			logFormat(logger, t, my.fields).Panic(args...)
		}
	}
}
func (my LogEntry) Error(args ...interface{}) {
	t := logNowTime()
	for _, logger := range my.loggers {
		logFormat(logger, t, my.fields).Error(args...)
	}
}
func (my LogEntry) Warn(args ...interface{}) {
	t := logNowTime()
	for _, logger := range my.loggers {
		logFormat(logger, t, my.fields).Warn(args...)
	}
}
func (my LogEntry) Info(args ...interface{}) {
	t := logNowTime()
	for _, logger := range my.loggers {
		logFormat(logger, t, my.fields).Info(args...)
	}
}
func (my LogEntry) Debug(args ...interface{}) {
	t := logNowTime()
	for _, logger := range my.loggers {
		logFormat(logger, t, my.fields).Debug(args...)
	}
}
func (my LogEntry) Trace(args ...interface{}) {
	t := logNowTime()
	for _, logger := range my.loggers {
		logFormat(logger, t, my.fields).Trace(args...)
	}
}
