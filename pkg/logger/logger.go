package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init() error {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	log, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

func GetLogger() *zap.Logger {
	return log
}

// Fatal logs a fatal message and exits
func Fatal(msg string, err error) {
	if err != nil {
		log.Fatal(msg, zap.Error(err))
	} else {
		log.Fatal(msg)
	}
}

// Info logs an info message with fields
func Info(msg string, fields ...interface{}) {
	zapFields := make([]zap.Field, 0, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			zapFields = append(zapFields, zap.Any(fields[i].(string), fields[i+1]))
		}
	}
	log.Info(msg, zapFields...)
}

// Error logs an error message with fields
func Error(msg string, fields ...interface{}) {
	zapFields := make([]zap.Field, 0, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			zapFields = append(zapFields, zap.Any(fields[i].(string), fields[i+1]))
		}
	}
	log.Error(msg, zapFields...)
} 