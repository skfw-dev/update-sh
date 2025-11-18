package logger

import (
	"context"
	"fmt"
	"os"
	"update-sh/nextgen/log/common"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	common.Logger
	zapLogger *zap.Logger
	level     common.Level
	fields    common.Fields
}

// NewLogger creates a new production-ready Zap logger
func NewLogger(config *Config) *Logger {
	zapLevel := toZapLevel(config.Level)

	// Define Zap configuration for console and file outputs
	consoleEncoder, encoderConfig := createCustomConsoleEncoder(config)
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// Create a dynamic list of sinks (outputs)
	var cores []zapcore.Core

	// Console Output
	if config.Mode == common.ModeConsole || config.Mode == common.ModeBoth {
		consoleWriter := zapcore.Lock(os.Stdout)
		cores = append(cores, zapcore.NewCore(
			consoleEncoder,
			consoleWriter,
			zapLevel,
		))
	}

	// File Output (using an io.WriteSyncer for Zap)
	if config.Mode == common.ModeFileOnly || config.Mode == common.ModeBoth {
		fileWriter, _ := os.OpenFile(config.Filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		cores = append(cores, zapcore.NewCore(
			fileEncoder,
			zapcore.AddSync(fileWriter),
			zapLevel,
		))
	}

	// Calculate the caller skip
	callerSkip := config.CallerSkip + 1

	// Combine cores into a single Zap logger
	core := zapcore.NewTee(cores...)
	opts := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(callerSkip),
		zap.AddStacktrace(zapcore.ErrorLevel),
	}

	// Skip 1 to point to the caller of the wrapper
	zapLogger := zap.New(core, opts...)

	// Create a new Logger instance
	logger := &Logger{
		zapLogger: zapLogger.Named("Logger"),
		level:     config.Level,
	}

	return logger
}

// WithContext returns a new Logger with the given context
func (s *Logger) WithContext(ctx context.Context) common.Logger {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		logger := &Logger{
			zapLogger: s.zapLogger.With(
				zap.String("trace_id", span.SpanContext().TraceID().String()),
				zap.String("span_id", span.SpanContext().SpanID().String()),
			),
		}

		return logger
	}

	return s
}

// Field returns a new logger with the added field (contextual logging)
func (s *Logger) Field(name string, value any) common.Logger {
	field := common.Field{Key: name, Value: value}
	fields := append(s.fields, field)
	logger := &Logger{
		zapLogger: s.zapLogger,
		level:     s.level,
		fields:    fields,
	}

	return logger
}

// Fields returns a new logger with the added fields (contextual logging)
func (s *Logger) Fields(fields ...common.Field) common.Logger {
	fields = append(s.fields, fields...)
	logger := &Logger{
		zapLogger: s.zapLogger,
		level:     s.level,
		fields:    fields,
	}

	return logger
}

func (s *Logger) Log(level common.Level, msg string, fields ...common.Field) common.Logger {
	fields = append(s.fields, fields...)
	zapLevel := toZapLevel(level)
	zapFields := toZapFields(fields)
	s.zapLogger.Log(zapLevel, msg, zapFields...)
	return s
}

func (s *Logger) Logf(level common.Level, format string, args ...any) common.Logger {
	zapLevel := toZapLevel(level)
	zapMsg := fmt.Sprintf(format, args...)
	zapFields := toZapFields(s.fields)
	s.zapLogger.Log(zapLevel, zapMsg, zapFields...)
	return s
}

// Print map to Info level
func (s *Logger) Print(msg string, fields ...common.Field) common.Logger {
	fields = append(s.fields, fields...)
	zapFields := toZapFields(fields)
	s.zapLogger.Info(msg, zapFields...)
	return s
}

// Printf map to Info level
func (s *Logger) Printf(format string, args ...any) common.Logger {
	zapMsg := fmt.Sprintf(format, args...)
	zapFields := toZapFields(s.fields)
	s.zapLogger.Info(zapMsg, zapFields...)
	return s
}

// Debug logs a message at Debug level
func (s *Logger) Debug(msg string, fields ...common.Field) common.Logger {
	fields = append(s.fields, fields...)
	zapFields := toZapFields(fields)
	s.zapLogger.Debug(msg, zapFields...)
	return s
}

// Debugf logs a message at Debug level
func (s *Logger) Debugf(format string, args ...any) common.Logger {
	zapMsg := fmt.Sprintf(format, args...)
	zapFields := toZapFields(s.fields)
	s.zapLogger.Debug(zapMsg, zapFields...)
	return s
}

// Info logs a message at Info level
func (s *Logger) Info(msg string, fields ...common.Field) common.Logger {
	fields = append(s.fields, fields...)
	zapFields := toZapFields(fields)
	s.zapLogger.Info(msg, zapFields...)
	return s
}

// Infof logs a message at Info level
func (s *Logger) Infof(format string, args ...any) common.Logger {
	zapMsg := fmt.Sprintf(format, args...)
	zapFields := toZapFields(s.fields)
	s.zapLogger.Info(zapMsg, zapFields...)
	return s
}

// Warn logs a message at Warn level
func (s *Logger) Warn(msg string, fields ...common.Field) common.Logger {
	fields = append(s.fields, fields...)
	zapFields := toZapFields(fields)
	s.zapLogger.Warn(msg, zapFields...)
	return s
}

// Warnf logs a message at Warn level
func (s *Logger) Warnf(format string, args ...any) common.Logger {
	zapMsg := fmt.Sprintf(format, args...)
	zapFields := toZapFields(s.fields)
	s.zapLogger.Warn(zapMsg, zapFields...)
	return s
}

// Error logs a message at Error level
func (s *Logger) Error(msg string, fields ...common.Field) common.Logger {
	fields = append(s.fields, fields...)
	zapFields := toZapFields(fields)
	s.zapLogger.Error(msg, zapFields...)
	return s
}

// Errorf logs a message at Error level
func (s *Logger) Errorf(format string, args ...any) common.Logger {
	zapMsg := fmt.Sprintf(format, args...)
	zapFields := toZapFields(s.fields)
	s.zapLogger.Error(zapMsg, zapFields...)
	return s
}

// Fatal logs a message at Fatal level
func (s *Logger) Fatal(msg string, fields ...common.Field) common.Logger {
	fields = append(s.fields, fields...)
	zapFields := toZapFields(fields)
	s.zapLogger.Fatal(msg, zapFields...)
	return s
}

// Fatalf logs a message at Fatal level
func (s *Logger) Fatalf(format string, args ...any) common.Logger {
	zapMsg := fmt.Sprintf(format, args...)
	zapFields := toZapFields(s.fields)
	s.zapLogger.Fatal(zapMsg, zapFields...)
	return s
}

// Panic logs a message at Panic level and then panics
func (s *Logger) Panic(msg string, fields ...common.Field) common.Logger {
	fields = append(s.fields, fields...)
	zapFields := toZapFields(fields)
	s.zapLogger.Panic(msg, zapFields...)
	return s
}

// Panicf logs a message at Panic level and then panics
func (s *Logger) Panicf(format string, args ...any) common.Logger {
	zapMsg := fmt.Sprintf(format, args...)
	zapFields := toZapFields(s.fields)
	s.zapLogger.Panic(zapMsg, zapFields...)
	return s
}
