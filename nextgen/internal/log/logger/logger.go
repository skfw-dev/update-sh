package logger

import (
	"context"
	"os"
	"update-sh/nextgen/internal/log"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	log.Logger
	zapLogger *zap.Logger
	level     log.Level
}

// NewZapLogger creates a new production-ready Zap logger
func NewZapLogger(config *Config) *Logger {
	zapLevel := toZapLevel(config.Level)

	// Define Zap configuration for console and file outputs
	consoleEncoder, encoderConfig := createCustomConsoleEncoder(config)
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// Create a dynamic list of sinks (outputs)
	var cores []zapcore.Core

	// Console Output
	if config.Mode == log.ModeConsole || config.Mode == log.ModeBoth {
		consoleWriter := zapcore.Lock(os.Stdout)
		cores = append(cores, zapcore.NewCore(
			consoleEncoder,
			consoleWriter,
			zapLevel,
		))
	}

	// File Output (using an io.WriteSyncer for Zap)
	if config.Mode == log.ModeFileOnly || config.Mode == log.ModeBoth {
		fileWriter, _ := os.OpenFile(config.Filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		cores = append(cores, zapcore.NewCore(
			fileEncoder,
			zapcore.AddSync(fileWriter),
			zapLevel,
		))
	}

	// Combine cores into a single Zap logger
	core := zapcore.NewTee(cores...)

	// Skip 1 to point to the caller of the wrapper
	zapLogger := zap.New(core, zap.AddCallerSkip(config.CallerSkip+1))

	// Create a new Logger instance
	logger := &Logger{
		zapLogger: zapLogger,
		level:     config.Level,
	}

	return logger
}

// WithContext returns a new Logger with the given context
func (s *Logger) WithContext(ctx context.Context) log.Logger {
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
func (s *Logger) Field(name string, value any) log.Logger {
	zapLogger := s.zapLogger.With(zap.Any(name, value))
	logger := &Logger{
		zapLogger: zapLogger,
		level:     s.level,
	}

	return logger
}

// Fields returns a new logger with the added fields (contextual logging)
func (s *Logger) Fields(fields ...log.Field) log.Logger {
	zapLogger := s.zapLogger.With(toZapFields(fields)...)
	logger := &Logger{
		zapLogger: zapLogger,
		level:     s.level,
	}

	return logger
}

func (s *Logger) Log(level log.Level, msg string, fields ...log.Field) log.Logger {
	s.zapLogger.Log(toZapLevel(level), msg, toZapFields(fields)...)
	return s
}

func (s *Logger) Logf(level log.Level, format string, args ...any) log.Logger {
	s.zapLogger.Sugar().Logf(toZapLevel(level), format, args...)
	return s
}

// Print map to Info level
func (s *Logger) Print(msg string, fields ...log.Field) log.Logger {
	s.zapLogger.Info(msg, toZapFields(fields)...)
	return s
}

// Printf map to Info level
func (s *Logger) Printf(format string, args ...any) log.Logger {
	s.zapLogger.Sugar().Infof(format, args...)
	return s
}

// Debug logs a message at Debug level
func (s *Logger) Debug(msg string, fields ...log.Field) log.Logger {
	s.zapLogger.Debug(msg, toZapFields(fields)...)
	return s
}

// Debugf logs a message at Debug level
func (s *Logger) Debugf(format string, args ...any) log.Logger {
	s.zapLogger.Sugar().Debugf(format, args...)
	return s
}

// Info logs a message at Info level
func (s *Logger) Info(msg string, fields ...log.Field) log.Logger {
	s.zapLogger.Info(msg, toZapFields(fields)...)
	return s
}

// Infof logs a message at Info level
func (s *Logger) Infof(format string, args ...any) log.Logger {
	s.zapLogger.Sugar().Infof(format, args...)
	return s
}

// Warn logs a message at Warn level
func (s *Logger) Warn(msg string, fields ...log.Field) log.Logger {
	s.zapLogger.Warn(msg, toZapFields(fields)...)
	return s
}

// Warnf logs a message at Warn level
func (s *Logger) Warnf(format string, args ...any) log.Logger {
	s.zapLogger.Sugar().Warnf(format, args...)
	return s
}

// Error logs a message at Error level
func (s *Logger) Error(msg string, fields ...log.Field) log.Logger {
	s.zapLogger.Error(msg, toZapFields(fields)...)
	return s
}

// Errorf logs a message at Error level
func (s *Logger) Errorf(format string, args ...any) log.Logger {
	s.zapLogger.Sugar().Errorf(format, args...)
	return s
}

// Fatal logs a message at Fatal level
func (s *Logger) Fatal(msg string, fields ...log.Field) log.Logger {
	s.zapLogger.Fatal(msg, toZapFields(fields)...)
	return s
}

// Fatalf logs a message at Fatal level
func (s *Logger) Fatalf(format string, args ...any) log.Logger {
	s.zapLogger.Sugar().Fatalf(format, args...)
	return s
}

// Panic logs a message at Panic level and then panics
func (s *Logger) Panic(msg string, fields ...log.Field) log.Logger {
	s.zapLogger.Panic(msg, toZapFields(fields)...)
	return s
}

// Panicf logs a message at Panic level and then panics
func (s *Logger) Panicf(format string, args ...any) log.Logger {
	s.zapLogger.Sugar().Panicf(format, args...)
	return s
}
