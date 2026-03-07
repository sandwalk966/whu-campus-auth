package utils

import (
	"os"
	"strings"
	"whu-campus-auth/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func InitLogger(cfg *config.Config) error {
	level := parseLogLevel(cfg.Log.Level)

	encoder := getEncoder(cfg.Log.Format)

	var output zapcore.WriteSyncer
	if cfg.Log.Output == "stdout" {
		output = zapcore.AddSync(os.Stdout)
	} else {
		output = zapcore.AddSync(os.Stderr)
	}

	core := zapcore.NewCore(encoder, output, level)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	Logger = logger.Sugar()
	return nil
}

func parseLogLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func getEncoder(format string) zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if format == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func LogDebug(args ...interface{}) {
	Logger.Debug(args...)
}

func LogDebugf(template string, args ...interface{}) {
	Logger.Debugf(template, args...)
}

func LogInfo(args ...interface{}) {
	Logger.Info(args...)
}

func LogInfof(template string, args ...interface{}) {
	Logger.Infof(template, args...)
}

func LogWarn(args ...interface{}) {
	Logger.Warn(args...)
}

func LogWarnf(template string, args ...interface{}) {
	Logger.Warnf(template, args...)
}

func LogError(args ...interface{}) {
	Logger.Error(args...)
}

func LogErrorf(template string, args ...interface{}) {
	Logger.Errorf(template, args...)
}

func LogDPanic(args ...interface{}) {
	Logger.DPanic(args...)
}

func LogDPanicf(template string, args ...interface{}) {
	Logger.DPanicf(template, args...)
}

func LogPanic(args ...interface{}) {
	Logger.Panic(args...)
}

func LogPanicf(template string, args ...interface{}) {
	Logger.Panicf(template, args...)
}

func LogFatal(args ...interface{}) {
	Logger.Fatal(args...)
}

func LogFatalf(template string, args ...interface{}) {
	Logger.Fatalf(template, args...)
}

func Sync() error {
	return Logger.Sync()
}
