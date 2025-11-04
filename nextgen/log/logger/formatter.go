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
	Level   string    `json:"level"`
	Message string    `json:"message"`
	Data    any       `json:"data"`
}

type Formatter interface {
	SetTextMode(textMode TextMode)
	SetFormatEncoder(formatEncode FormatEncoderKind)
	SetEncodeTimeLayout(encodeTimeLayout string)
	Format(time time.Time, level common.Level, message string, fields common.Fields) string
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

func (s *CustomFormatter) formatString(t time.Time, level common.Level, message string, fields common.Fields) string {
	formattedTime := t.UTC().Format(s.encodeTimeLayout)
	formattedMessage := formatMessageByTextMode(s.textMode, message, fields)
	return fmt.Sprintf("[%s] [%s] %s", formattedTime, level.String(), formattedMessage)
}

func (s *CustomFormatter) formatJSON(t time.Time, level common.Level, message string, fields common.Fields) string {
	cleanedMessage := strings.Trim(message, "\r\n")

	data := &FormatEntry{
		Time:    t.UTC(),
		Level:   level.String(),
		Message: cleanedMessage,
		Data:    fields.ToMap(),
	}

	b, _ := json.Marshal(data)
	return string(b)
}

func (s *CustomFormatter) formatXML(t time.Time, level common.Level, message string, fields common.Fields) string {
	formattedTime := t.UTC().Format(s.encodeTimeLayout)
	cleanedMessage := strings.Trim(message, "\r\n")

	var b strings.Builder
	b.WriteString(fmt.Sprintf("<date>%s</date>", formattedTime))
	b.WriteString(fmt.Sprintf("<level>%s</level>", level.String()))
	b.WriteString(fmt.Sprintf("<message>%s</message>", cleanedMessage))
	b.WriteString(fmt.Sprintf("<data>%s</data>", toFieldsJSON(fields)))
	return b.String()
}

// Format function
func (s *CustomFormatter) Format(t time.Time, level common.Level, message string, fields common.Fields) string {
	switch s.formatEncode {
	case TextFormatEncoder:
		return s.formatString(t, level, message, fields)
	case JSONFormatEncoder:
		return s.formatJSON(t, level, message, fields)
	case XMLFormatEncoder:
		return s.formatXML(t, level, message, fields)
	default:
		return s.formatString(t, level, message, fields)
	}
}
