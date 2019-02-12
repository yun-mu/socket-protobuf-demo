/*
Package log get the common logger
*/
package log

import (
	"config"
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

func GetLogger() *logrus.Logger {
	return getLogger(config.Conf.AppInfo.Env)
}

func getLogger(env string) *logrus.Logger {
	logger := &logrus.Logger{
		Out:       os.Stdout,
		Hooks:     make(logrus.LevelHooks),
		Formatter: new(logrus.TextFormatter),
		Level:     logrus.WarnLevel,
	}
	if env == "prod" {
		logger.Formatter = new(logrus.JSONFormatter)
	}
	logConfig := config.Conf.Log

	baseLogPath := path.Join(logConfig.LogBasePath, logConfig.LogFileName)
	writer, err := rotatelogs.New(
		baseLogPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPath),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour), // Log cutting interval
	)
	if err != nil {
		logrus.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.JSONFormatter{})

	logger.AddHook(lfHook)
	return logger
}
