package log

import "context"

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
	Name  string
	Value any
}

func NewField(name string, value any) *Field {
	return &Field{
		Name:  name,
		Value: value,
	}
}

type Fields []Field

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
