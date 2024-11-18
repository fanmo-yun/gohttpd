package logger

import (
	"fmt"
	"gohttpd/utils"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func ZapSync() {
	zap.L().Sync()
}

func NewLogger(lc utils.LoggerConfig) func() {
	switch lc.Out {
	case "console":
		NewStdOutLogger(lc.Level)
	case "file":
		NewFileLogger(lc.Level)
	default:
		zap.L().Fatal("gohttpd: Log Output Error type")
	}
	return ZapSync
}

func NewStdOutLogger(l string) {
	var logLevel zapcore.Level
	err := logLevel.UnmarshalText([]byte(l))
	if err != nil {
		zap.L().Fatal("gohttpd: Log Level Init Fatal", zap.String("unmarshaltext", err.Error()))
	}

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(logLevel),
		Development:      true,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := cfg.Build()
	if err != nil {
		zap.L().Fatal("gohttpd: Log Build Fatal", zap.String("build", err.Error()))
	}

	zap.ReplaceGlobals(logger)
}

func NewFileLogger(l string) {
	var logLevel zapcore.Level
	err := logLevel.UnmarshalText([]byte(l))
	if err != nil {
		zap.L().Fatal("gohttpd: Log Level Init Fatal", zap.String("unmarshaltext", err.Error()))
	}

	timestamp := time.Now().Format("2006-01-02T15-04-05")
	logFile := fmt.Sprintf("%s.log", timestamp)

	lumberjackLogger := &lumberjack.Logger{
		Filename:   filepath.Join("log", logFile),
		MaxSize:    128,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(lumberjackLogger),
		logLevel,
	)

	logger := zap.New(core)

	zap.ReplaceGlobals(logger)
}
