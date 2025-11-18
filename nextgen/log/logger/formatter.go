package logger

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"update-sh/nextgen/log/common"
)

type TextMode string

const (
	TXTLogFmt TextMode = "txt"
	XMLLogFmt TextMode = "xml"
)

type FormatEncoderKind string

const (
	TextFormatEncoder FormatEncoderKind = "text"
	JSONFormatEncoder FormatEncoderKind = "json"
	XMLFormatEncoder  FormatEncoderKind = "xml"
)

type FormatEntry struct {
	Time    time.Time `json:"time"`
	Name    string    `json:"name"`
	Level   string    `json:"level"`
	Caller  string    `json:"caller"`
	Message string    `json:"message"`
	Stack   string    `json:"stack"`
	Data    any       `json:"data"`
}

type Formatter interface {
	SetTextMode(textMode TextMode)
	SetFormatEncoder(formatEncode FormatEncoderKind)
	SetEncodeTimeLayout(encodeTimeLayout string)
	Format(t time.Time, name string, level common.Level, caller string, message string, stack string, fields common.Fields) string
}

type CustomFormatter struct {
	Formatter
	textMode         TextMode
	formatEncode     FormatEncoderKind
	encodeTimeLayout string
}

func NewCustomFormatter() *CustomFormatter {
	return &CustomFormatter{
		textMode:         TXTLogFmt,
		formatEncode:     TextFormatEncoder,
		encodeTimeLayout: "2006-01-02T15:04:05.000Z07:00",
	}
}

func (s *CustomFormatter) SetTextMode(textMode TextMode) {
	s.textMode = textMode
}

func (s *CustomFormatter) SetFormatEncoder(formatEncode FormatEncoderKind) {
	s.formatEncode = formatEncode
}

func (s *CustomFormatter) SetEncodeTimeLayout(encodeTimeLayout string) {
	s.encodeTimeLayout = encodeTimeLayout
}

func (s *CustomFormatter) formatString(t time.Time, name string, level common.Level, caller string, message string, stack string, fields common.Fields) string {
	formattedTime := t.UTC().Format(s.encodeTimeLayout)
	formattedMessage := formatMessageByTextMode(s.textMode, message, fields)
	return fmt.Sprintf("[%s] [%s] [%s] %s %s %s", formattedTime, name, level.String(), caller, formattedMessage, stack)
}

func (s *CustomFormatter) formatJSON(t time.Time, name string, level common.Level, caller string, message string, stack string, fields common.Fields) string {
	cleanedMessage := strings.Trim(message, "\r\n")

	data := &FormatEntry{
		Time:    t.UTC(),
		Name:    name,
		Level:   level.String(),
		Caller:  caller,
		Message: cleanedMessage,
		Stack:   stack,
		Data:    fields.ToMap(),
	}

	b, _ := json.Marshal(data)
	return string(b)
}

func (s *CustomFormatter) formatXML(t time.Time, name string, level common.Level, caller string, message string, stack string, fields common.Fields) string {
	formattedTime := t.UTC().Format(s.encodeTimeLayout)
	cleanedMessage := strings.Trim(message, "\r\n")

	var b strings.Builder
	b.WriteString(fmt.Sprintf("<time>%s</time>", formattedTime))
	b.WriteString(fmt.Sprintf("<name>%s</name>", name))
	b.WriteString(fmt.Sprintf("<level>%s</level>", level.String()))
	b.WriteString(fmt.Sprintf("<caller>%s</caller>", caller))
	b.WriteString(fmt.Sprintf("<message>%s</message>", cleanedMessage))
	b.WriteString(fmt.Sprintf("<stack>%s</stack>", stack))
	b.WriteString(fmt.Sprintf("<data>%s</data>", toFieldsJSON(fields)))
	return b.String()
}

// Format function
func (s *CustomFormatter) Format(t time.Time, name string, level common.Level, caller string, message string, stack string, fields common.Fields) string {
	switch s.formatEncode {
	case TextFormatEncoder:
		return s.formatString(t, name, level, caller, message, stack, fields)
	case JSONFormatEncoder:
		return s.formatJSON(t, name, level, caller, message, stack, fields)
	case XMLFormatEncoder:
		return s.formatXML(t, name, level, caller, message, stack, fields)
	default:
		return s.formatString(t, name, level, caller, message, stack, fields)
	}
}
