package logger

import (
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
)

type Logger struct {
	zap    *zap.Logger
	sugar  *zap.SugaredLogger
	prefix string
	level  zapcore.Level
}

type Config struct {
	Level      string `json:"level"`
	Encoding   string `json:"encoding"`
	OutputPath string `json:"output_path"`
	MaxSize    int    `json:"max_size"`
	MaxAge     int    `json:"max_age"`
	MaxBackups int    `json:"max_backups"`
	Compress   bool   `json:"compress"`
}

func DefaultConfig() Config {
	return Config{
		Level:      InfoLevel,
		Encoding:   "json",
		OutputPath: "logs/app.log",
		MaxSize:    100,
		MaxAge:     28,
		MaxBackups: 3,
		Compress:   true,
	}
}

func stringToZapLevel(level string) zapcore.Level {
	switch level {
	case DebugLevel:
		return zap.DebugLevel
	case WarnLevel:
		return zap.WarnLevel
	case ErrorLevel:
		return zap.ErrorLevel
	case FatalLevel:
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func NewLogger(config Config, prefix string) (*Logger, error) {
	logDir := filepath.Dir(config.OutputPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename:   config.OutputPath,
		MaxSize:    config.MaxSize,
		MaxAge:     config.MaxAge,
		MaxBackups: config.MaxBackups,
		Compress:   config.Compress,
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		FunctionKey:    zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	level := stringToZapLevel(config.Level)

	var core zapcore.Core
	if config.Encoding == "json" {
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberjackLogger)),
			level,
		)
	} else {
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberjackLogger)),
			level,
		)
	}

	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	defer zapLogger.Sync()

	if prefix != "" {
		zapLogger = zapLogger.Named(prefix)
	}

	sugarLogger := zapLogger.Sugar()

	return &Logger{
		zap:    zapLogger,
		sugar:  sugarLogger,
		prefix: prefix,
		level:  level,
	}, nil
}

func (l *Logger) WithPrefix(prefix string) *Logger {
	newLogger := l.zap.Named(prefix)
	return &Logger{
		zap:    newLogger,
		sugar:  newLogger.Sugar(),
		prefix: prefix,
		level:  l.level,
	}
}

func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.sugar.Debugw(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...interface{}) {
	l.sugar.Infow(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.sugar.Warnw(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...interface{}) {
	l.sugar.Errorw(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...interface{}) {
	l.sugar.Fatalw(msg, fields...)
}

func (l *Logger) Sync() error {
	return l.zap.Sync()
}

var defaultLogger *Logger

func InitDefaultLogger(config Config) error {
	var err error
	defaultLogger, err = NewLogger(config, "")
	return err
}

func Default() *Logger {
	if defaultLogger == nil {
		config := DefaultConfig()
		defaultLogger, _ = NewLogger(config, "")
	}
	return defaultLogger
}

func Debug(msg string, fields ...interface{}) {
	Default().Debug(msg, fields...)
}

func Info(msg string, fields ...interface{}) {
	Default().Info(msg, fields...)
}

func Warn(msg string, fields ...interface{}) {
	Default().Warn(msg, fields...)
}

func Error(msg string, fields ...interface{}) {
	Default().Error(msg, fields...)
}

func Fatal(msg string, fields ...interface{}) {
	Default().Fatal(msg, fields...)
}
