package logger

import (
	"github.com/qq5272689/goutils/zap_logger"
	"go.uber.org/zap"
	"runtime/debug"
)

var logger *zap.Logger
var Closer zap_logger.Closer

func init() {
	logger, _ = zap.NewDevelopment()
	Closer = func() error {
		return logger.Sync()
	}
}

func LoggerInit(env, dir, service, when string) {
	logger.Sync()
	var err error
	if env == "dev" || env == "local" {
		logger, Closer, err = zap_logger.GetDevLogger(dir, service, when)
	} else {
		logger, Closer, err = zap_logger.GetProdLogger(dir, service, when)
	}
	if err != nil {
		l, _ := zap.NewDevelopment()
		l.Sugar().Fatal("logger init fail!!!", zap.Error(err))
	}
	logger.Debug("logger init ok", zap.String("dir", dir))
}

func GetLogger() *zap.Logger {
	return logger
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg+"\n---------------------------->stack:\n"+string(debug.Stack())+"\n<----------------------------stack", fields...)
}
