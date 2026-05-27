package util

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	loggerOnce sync.Once
	loggerMu sync.RWMutex
)

type LogConfig struct {
	Level      string `json:"level"`
	Filename   string `json:"filename"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	Compress   bool   `json:"compress"`
}

func InitLogger(config LogConfig) error {
	var err error
	loggerOnce.Do(func() {
		err = initLogger(config)
	})
	return err
}

func initLogger(config LogConfig) error {
	if config.Filename == "" {
		config.Filename = "./logs/app.log"
	}
	if config.MaxSize == 0 {
		config.MaxSize = 100
	}
	if config.MaxBackups == 0 {
		config.MaxBackups = 30
	}
	if config.MaxAge == 0 {
		config.MaxAge = 7
	}

	if err := os.MkdirAll(filepath.Dir(config.Filename), os.ModePerm); err != nil {
		logger = zap.NewNop()
		return err
	}

	level := zapcore.InfoLevel
	if config.Level == "debug" {
		level = zapcore.DebugLevel
	} else if config.Level == "warn" {
		level = zapcore.WarnLevel
	} else if config.Level == "error" {
		level = zapcore.ErrorLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	fileEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	file, err := os.OpenFile(config.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger = zap.NewNop()
		return err
	}

	fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(file), level),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	)

	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	return nil
}

func GetLogger() *zap.Logger {
	loggerMu.RLock()
	if logger != nil {
		defer loggerMu.RUnlock()
		return logger
	}
	loggerMu.RUnlock()

	loggerMu.Lock()
	defer loggerMu.Unlock()

	if logger != nil {
		return logger
	}

	err := InitLogger(LogConfig{
		Level:      "info",
		Filename:   "./logs/app.log",
		MaxSize:    100,
		MaxBackups: 30,
		MaxAge:     7,
	})
	if err != nil {
		logger = zap.NewNop()
	}

	return logger
}

func StringField(key string, value string) zap.Field {
	return zap.String(key, value)
}

func IntField(key string, value int) zap.Field {
	return zap.Int(key, value)
}

func Int64Field(key string, value int64) zap.Field {
	return zap.Int64(key, value)
}

func Uint64Field(key string, value uint64) zap.Field {
	return zap.Uint64(key, value)
}

func BoolField(key string, value bool) zap.Field {
	return zap.Bool(key, value)
}

func DurationField(key string, value time.Duration) zap.Field {
	return zap.Duration(key, value)
}

func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}

func LogError(err error, context string) {
	GetLogger().Error(context,
		zap.Error(err),
		zap.String("time", time.Now().Format(time.RFC3339)),
	)
}

func LogOrderOperation(orderNo string, action string, userID uint64, details map[string]interface{}) {
	GetLogger().Info("order_operation",
		zap.String("order_no", orderNo),
		zap.String("action", action),
		zap.Uint64("user_id", userID),
		zap.Any("details", details),
		zap.String("time", time.Now().Format(time.RFC3339)),
	)
}

func LogAuthOperation(action string, phone string, success bool, reason string) {
	GetLogger().Info("auth_operation",
		zap.String("action", action),
		zap.String("phone", maskPhone(phone)),
		zap.Bool("success", success),
		zap.String("reason", reason),
		zap.String("time", time.Now().Format(time.RFC3339)),
	)
}

func maskPhone(phone string) string {
	if len(phone) != 11 {
		return phone
	}
	return phone[:3] + "****" + phone[7:]
}

func Sync() {
	if logger != nil {
		logger.Sync()
	}
}
