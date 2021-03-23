package zap_logger

import (
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func GetProdJsonLogger(path, service, when string) (logger *zap.Logger, closer Closer, err error) {
	core, err := buildZapJsonCore(false)
	if err != nil {
		return nil, nil, err
	}
	//zap.AddCallerSkip(1),
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2), zap.Development())
	return logger, func() error {
		return multierr.Combine(logger.Sync())
	}, nil
}

func GetDevJsonLogger(path, service, when string) (logger *zap.Logger, closer Closer, err error) {
	core, err := buildZapJsonCore(true)
	if err != nil {
		return nil, nil, err
	}
	//zap.AddCallerSkip(1),
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2), zap.Development())
	return logger, func() error {
		return logger.Sync()
	}, nil
}

func buildZapJsonCore(isdev bool) (core zapcore.Core, err error) {
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	if isdev {
		atomicLevel.SetLevel(zapcore.DebugLevel)
	} else {
		atomicLevel.SetLevel(zapcore.InfoLevel)
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
		zapcore.NewJSONEncoder(encoderConfig), // 编码器配置
		zapcore.AddSync(os.Stdout),
		atomicLevel, // 日志级别
	), nil
}
