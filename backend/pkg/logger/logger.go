package logger

import (
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	global *zap.Logger
	once   sync.Once
)

// Config controls logger initialization.
type Config struct {
	Environment string
	AppName     string
}

// Init creates the global structured logger.
func Init(cfg Config) (*zap.Logger, error) {
	var err error
	once.Do(func() {
		global, err = buildLogger(cfg)
	})
	if err != nil {
		return nil, err
	}
	return global, nil
}

// L returns the global logger. Panics if Init was not called.
func L() *zap.Logger {
	if global == nil {
		panic("logger: Init must be called before L()")
	}
	return global
}

// Sync flushes buffered log entries.
func Sync() {
	if global != nil {
		_ = global.Sync()
	}
}

func buildLogger(cfg Config) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	if strings.EqualFold(cfg.Environment, "development") {
		level = zapcore.DebugLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.MessageKey = "message"
	encoderConfig.LevelKey = "level"

	var zapCfg zap.Config
	if strings.EqualFold(cfg.Environment, "development") {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapCfg = zap.Config{
			Level:            zap.NewAtomicLevelAt(level),
			Development:      true,
			Encoding:         "console",
			EncoderConfig:    encoderConfig,
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
	} else {
		zapCfg = zap.Config{
			Level:            zap.NewAtomicLevelAt(level),
			Development:      false,
			Encoding:         "json",
			EncoderConfig:    encoderConfig,
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
	}

	log, err := zapCfg.Build(zap.Fields(
		zap.String("app", cfg.AppName),
		zap.String("environment", cfg.Environment),
		zap.Int("pid", os.Getpid()),
	))
	if err != nil {
		return nil, err
	}
	return log, nil
}

// Info logs an info message.
func Info(msg string, fields ...zap.Field) {
	L().Info(msg, fields...)
}

// Warn logs a warning message.
func Warn(msg string, fields ...zap.Field) {
	L().Warn(msg, fields...)
}

// Error logs an error message.
func Error(msg string, fields ...zap.Field) {
	L().Error(msg, fields...)
}

// Fatal logs a fatal message and exits.
func Fatal(msg string, fields ...zap.Field) {
	L().Fatal(msg, fields...)
}

// Debug logs a debug message.
func Debug(msg string, fields ...zap.Field) {
	L().Debug(msg, fields...)
}

// Startup logs application startup information.
func Startup(component string, fields ...zap.Field) {
	all := append([]zap.Field{zap.String("component", component)}, fields...)
	L().Info("startup", all...)
}

// Request logs an HTTP request.
func Request(method, path string, status int, latencyMs float64, fields ...zap.Field) {
	all := append([]zap.Field{
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status", status),
		zap.Float64("latency_ms", latencyMs),
	}, fields...)
	L().Info("request", all...)
}

// With creates a child logger with additional fields.
func With(fields ...zap.Field) *zap.Logger {
	return L().With(fields...)
}
