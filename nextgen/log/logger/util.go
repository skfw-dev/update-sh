package logger

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"update-sh/nextgen/cores"
	"update-sh/nextgen/log/common"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level interface {
	common.Level | zapcore.Level | string
}

// toLevel converts a level to a log.Level.
func toLevel[T Level](level T) common.Level {
	switch v := any(level).(type) {
	case common.Level:
		return v
	case zapcore.Level:
		return zapToLevel(v)
	case string:
		return strToLevel(v)
	default:
		return common.InfoLevel
	}
}

// zapToLevel converts a zapcore.Level to a log.Level.
func zapToLevel(level zapcore.Level) common.Level {
	switch level {
	case zapcore.DebugLevel:
		return common.DebugLevel
	case zapcore.InfoLevel:
		return common.InfoLevel
	case zapcore.WarnLevel:
		return common.WarnLevel
	case zapcore.ErrorLevel:
		return common.ErrorLevel
	case zapcore.FatalLevel:
		return common.FatalLevel
	case zapcore.PanicLevel:
		return common.PanicLevel
	default:
		return common.InfoLevel
	}
}

// strToLevel converts a string to a log.Level.
func strToLevel(level string) common.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug", "debugging":
		return common.DebugLevel
	case "info", "information":
		return common.InfoLevel
	case "warn", "warning":
		return common.WarnLevel
	case "error":
		return common.ErrorLevel
	case "fatal":
		return common.FatalLevel
	case "panic":
		return common.PanicLevel
	default:
		return common.InfoLevel
	}
}

// toZapLevel converts the custom log.Level to a zapcore.Level.
func toZapLevel(level common.Level) zapcore.Level {
	switch level {
	case common.DebugLevel:
		return zapcore.DebugLevel
	case common.InfoLevel:
		return zapcore.InfoLevel
	case common.WarnLevel:
		return zapcore.WarnLevel
	case common.ErrorLevel:
		return zapcore.ErrorLevel
	case common.FatalLevel:
		return zapcore.FatalLevel
	case common.PanicLevel:
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}

// toFields converts the zap.Field slice to the custom log.Field slice.
func toFields(fields []zapcore.Field) common.Fields {
	result := make(common.Fields, len(fields))
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
		if field := common.NewField(f.Key, value); field != nil {
			result[i] = *field
		}
	}

	return result
}

// toZapFields converts the custom log.Field slice to a zap.Field slice.
func toZapFields(fields common.Fields) []zap.Field {
	result := make([]zap.Field, len(fields))
	for i, f := range fields {
		result[i] = zap.Any(f.Key, f.Value)
	}

	return result
}

// formatMessageByTextMode formats a message with fields into a single string.
func formatMessageByTextMode(textMode TextMode, message string, fields common.Fields) string {
	message = strings.Trim(message, "\r\n")

	switch textMode {
	case TXTLogFmt:
		var b strings.Builder
		b.WriteString(message)
		if len(fields) > 0 {
			b.WriteString(" ")
			b.WriteString(toFieldsString(fields))
		}

		return b.String()

	case XMLLogFmt:
		var b strings.Builder
		b.WriteString(fmt.Sprintf("<message>%s</message>", message))
		if len(fields) > 0 {
			b.WriteString(" ")
			b.WriteString(fmt.Sprintf("<data>%s</data>", toFieldsJSON(fields)))
		}

		return b.String()

	default:
		return message
	}
}

// toFieldString converts a log.Field into a string representation.
func toFieldString(field *common.Field) string {
	return fmt.Sprintf("%s=%s", field.GetKey(), strconv.Quote(field.ToString()))
}

// toFieldsString converts a slice of log.Fields into a string representation.
func toFieldsString(fields common.Fields) string {
	result := make([]string, len(fields))
	for i, field := range fields {
		result[i] = toFieldString(&field)
	}
	return strings.Join(result, " ")
}

// toFieldsJSON converts a slice of log.Fields into a JSON-like string representation.
func toFieldsJSON(fields common.Fields) string {
	b, _ := json.Marshal(fields.ToMap())
	return string(b)
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
		return zap.NewDevelopmentEncoderConfig()
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
