package logger

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"update-sh/nextgen/log"
)

type FormatEncoderKind string

const (
	JSONFormatEncoder FormatEncoderKind = "json"
	TextFormatEncoder FormatEncoderKind = "text"
)

type FormatEntry struct {
	Time    time.Time      `json:"time"`
	Level   string         `json:"level"`
	Message string         `json:"message"`
	Fields  map[string]any `json:"fields"`
}

type Formatter interface {
	SetLogFmt(logFmt bool)
	SetFormatEncoder(formatEncode FormatEncoderKind)
	SetEncodeTimeLayout(encodeTimeLayout string)
	Format(time time.Time, level log.Level, message string, fields log.Fields) string
}

type CustomFormatter struct {
	Formatter
	logFmt           bool
	formatEncode     FormatEncoderKind
	encodeTimeLayout string
}

func NewCustomFormatter() *CustomFormatter {
	return &CustomFormatter{
		logFmt:           true,
		formatEncode:     TextFormatEncoder,
		encodeTimeLayout: "2006-01-02T15:04:05.999Z07:00",
	}
}

func (s *CustomFormatter) SetLogFmt(logFmt bool) {
	s.logFmt = logFmt
}

func (s *CustomFormatter) SetFormatEncoder(formatEncode FormatEncoderKind) {
	s.formatEncode = formatEncode
}

func (s *CustomFormatter) SetEncodeTimeLayout(encodeTimeLayout string) {
	s.encodeTimeLayout = encodeTimeLayout
}

func (s *CustomFormatter) formatJSON(t time.Time, level log.Level, message string, fields log.Fields) string {
	levelString := level.String()
	cleanedMessage := strings.Trim(message, "\r\n")
	mapFields := fields.ToMap()

	data := &FormatEntry{
		Time:    t.UTC(),
		Level:   levelString,
		Message: cleanedMessage,
		Fields:  mapFields,
	}

	b, _ := json.Marshal(data)
	return string(b)
}

func (s *CustomFormatter) formatString(t time.Time, level log.Level, message string, fields log.Fields) string {
	formattedTime := t.UTC().Format(s.encodeTimeLayout)
	levelString := level.String()
	formattedMessage := formatLogMessageWithFields(s.logFmt, message, fields)
	return fmt.Sprintf("[%s] [%s] %s", formattedTime, levelString, formattedMessage)
}

// Format function
func (s *CustomFormatter) Format(t time.Time, level log.Level, message string, fields log.Fields) string {
	switch s.formatEncode {
	case JSONFormatEncoder:
		return s.formatJSON(t, level, message, fields)

	case TextFormatEncoder:
		return s.formatString(t, level, message, fields)

	default:
		return s.formatString(t, level, message, fields)
	}
}
