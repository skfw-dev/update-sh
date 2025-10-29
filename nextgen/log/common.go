package log

import (
	"context"
	"fmt"
	"strconv"
	"update-sh/nextgen/cores"
	"update-sh/nextgen/cores/caseconv"
)

// Mode defines the desired output.
type Mode int

const (
	ModeConsole  Mode = iota // Silent only
	ModeFileOnly             // File only
	ModeBoth                 // Console and File
)

type Level int8

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
)

func (s Level) String() string {
	switch s {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	case PanicLevel:
		return "PANIC"
	default:
		return "ERROR"
	}
}

type Field struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

func NewField(name string, value any) *Field {
	return &Field{
		Key:   name,
		Value: value,
	}
}

func (s *Field) GetKey() string {
	return caseconv.ToCamelCase(s.Key)
}

func (s *Field) GetValue() any {
	return s.Value
}

func (s *Field) GetEscapeValue() any {
	switch v := s.Value.(type) {
	case string:
		return strconv.Quote(v)
	case []byte:
		return cores.EscapeBytes(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (s *Field) ToMap() map[string]any {
	return map[string]any{
		s.Key: s.Value,
	}
}

type Fields []Field

func (s Fields) Count() int {
	return len(s)
}

func (s Fields) ToMap() map[string]any {
	result := make(map[string]any)
	for _, field := range s {
		result[(&field).GetKey()] = field.Value
	}
	return result
}

type Logger interface {
	WithContext(ctx context.Context) Logger
	Field(name string, value any) Logger
	Fields(fields ...Field) Logger
	Log(level Level, msg string, fields ...Field) Logger
	Logf(level Level, format string, args ...any) Logger
	Print(msg string, fields ...Field) Logger
	Printf(format string, args ...any) Logger
	Debug(msg string, fields ...Field) Logger
	Debugf(format string, args ...any) Logger
	Info(msg string, fields ...Field) Logger
	Infof(format string, args ...any) Logger
	Warn(msg string, fields ...Field) Logger
	Warnf(format string, args ...any) Logger
	Error(msg string, fields ...Field) Logger
	Errorf(format string, args ...any) Logger
	Fatal(msg string, fields ...Field) Logger
	Fatalf(format string, args ...any) Logger
	Panic(msg string, fields ...Field) Logger
	Panicf(format string, args ...any) Logger
}
