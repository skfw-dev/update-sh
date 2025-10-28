package logger

import (
	"strings"
	"update-sh/nextgen/cores"
	"update-sh/nextgen/log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level interface {
	log.Level | zapcore.Level | string
}

// toLevel converts a level to a log.Level.
func toLevel[T Level](level T) log.Level {
	switch mod := any(level); v := mod.(type) {
	case log.Level:
		return v
	case zapcore.Level:
		return zapToLevel(v)
	case string:
		return strToLevel(v)
	default:
		panic("invalid level type")
	}
}

// zapToLevel converts a zapcore.Level to a log.Level.
func zapToLevel(level zapcore.Level) log.Level {
	switch level {
	case zapcore.DebugLevel:
		return log.DebugLevel
	case zapcore.InfoLevel:
		return log.InfoLevel
	case zapcore.WarnLevel:
		return log.WarnLevel
	case zapcore.ErrorLevel:
		return log.ErrorLevel
	case zapcore.FatalLevel:
		return log.FatalLevel
	case zapcore.PanicLevel:
		return log.PanicLevel
	default:
		panic("invalid level type")
	}
}

// strToLevel converts a string to a log.Level.
func strToLevel(level string) log.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug", "debugging":
		return log.DebugLevel
	case "info", "information":
		return log.InfoLevel
	case "warn", "warning":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	case "panic":
		return log.PanicLevel
	default:
		panic("invalid level type")
	}
}

// toZapLevel converts the custom log.Level to a zapcore.Level.
func toZapLevel(level log.Level) zapcore.Level {
	switch level {
	case log.DebugLevel:
		return zapcore.DebugLevel
	case log.InfoLevel:
		return zapcore.InfoLevel
	case log.WarnLevel:
		return zapcore.WarnLevel
	case log.ErrorLevel:
		return zapcore.ErrorLevel
	case log.FatalLevel:
		return zapcore.FatalLevel
	case log.PanicLevel:
		return zapcore.PanicLevel
	default:
		panic("invalid level type")
	}
}

// toFields converts the zap.Field slice to the custom log.Field slice.
func toFields(fields []zapcore.Field) []log.Field {
	result := make([]log.Field, len(fields))
	for i, f := range fields {
		var value any
		switch {
		case f.Integer > 0:
			value = f.Integer
		case len(f.String) > 0:
			value = f.String
		default:
			value = f.Interface
		}
		if field := log.NewField(f.Key, value); field != nil {
			result[i] = *field
		}
	}

	return result
}

// toZapFields converts the custom log.Field slice to a zap.Field slice.
func toZapFields(fields []log.Field) []zap.Field {
	result := make([]zap.Field, len(fields))
	for i, f := range fields {
		result[i] = zap.Any(f.Name, f.Value)
	}

	return result
}

// createCustomEncoderConfig creates a custom encoder configuration based on the provided configuration.
func createCustomEncoderConfig(config *Config) zapcore.EncoderConfig {
	switch config.AppEnv {
	case cores.Production:
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		return encoderConfig

	case cores.Development:
		return zap.NewDevelopmentEncoderConfig()

	default:
		panic("invalid app env")
	}
}

// createCustomConsoleEncoder creates a custom console encoder based on the provided configuration.
func createCustomConsoleEncoder(config *Config) (zapcore.Encoder, zapcore.EncoderConfig) {
	encoderConfig := createCustomEncoderConfig(config)

	if config.Formatter != nil {
		// Use custom encoder which calls the provided Formatter method
		return NewCustomConsoleEncoder(config.Formatter, encoderConfig), encoderConfig
	}

	// Fallback to Zap's standard console format
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	return consoleEncoder, encoderConfig
}
