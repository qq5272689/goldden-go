package zap_logger

import (
	"github.com/qq5272689/goldden-go/pkg/utils/log_writer"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Closer func() error

func GetProdLogger(path, service, when string) (logger *zap.Logger, closer Closer, err error) {
	core, err := buildZapCore(path, service, when, false)
	if err != nil {
		return nil, nil, err
	}
	hook, errLogger, err := zapLevelFileHook(path, service, when, zap.ErrorLevel)
	//zap.AddCallerSkip(1),
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2), zap.Development(), zap.Hooks(hook))
	return logger, func() error {
		return multierr.Combine(errLogger.Sync(), logger.Sync())
	}, nil
}

func GetDevLogger(path, service, when string) (logger *zap.Logger, closer Closer, err error) {
	core, err := buildZapCore(path, service, when, true)
	if err != nil {
		return nil, nil, err
	}
	//zap.AddCallerSkip(1),
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2), zap.Development())
	return logger, func() error {
		return logger.Sync()
	}, nil
}

func buildZapCore(path, service, when string, isdev bool) (core zapcore.Core, err error) {
	lw, err := log_writer.NewLogWriter(service, path, when)
	if err != nil {
		return nil, err
	}
	var ws zapcore.WriteSyncer
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	if isdev {
		atomicLevel.SetLevel(zapcore.DebugLevel)
		ws = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), lw)
	} else {
		atomicLevel.SetLevel(zapcore.InfoLevel)
		ws = lw
	}

	//return zap.NewDevelopment()

	//公用编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		//EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeCaller: zapcore.ShortCallerEncoder, //
		EncodeName:   zapcore.FullNameEncoder,
	}

	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), // 编码器配置
		ws,
		atomicLevel, // 日志级别
	), nil
}

func zapLevelFileHook(path, service, when string, level zapcore.Level) (hook func(zapcore.Entry) error, logger *zap.Logger, err error) {
	core, err := buildZapCore(path, service+"-"+level.String(), when, false)
	if err != nil {
		return nil, nil, err
	}
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(3), zap.Development())
	return func(entry zapcore.Entry) error {
		if entry.Level < level {
			return nil
		}
		switch entry.Level {
		case zap.InfoLevel:
			logger.Info(entry.Message)
		case zap.WarnLevel:
			logger.Warn(entry.Message)
		case zap.ErrorLevel:
			logger.Error(entry.Message)
		}
		return nil
	}, logger, nil

}
