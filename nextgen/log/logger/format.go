package logger

import (
	"fmt"
	"strings"
	"time"
	"update-sh/nextgen/log"
)

type Formatter interface {
	Format(time time.Time, level log.Level, message string, fields log.Fields) string
}

type CustomFormatter struct {
	Formatter
}

func NewCustomFormatter() *CustomFormatter {
	return &CustomFormatter{}
}

// Format function
func (s *CustomFormatter) Format(time time.Time, level log.Level, message string, fields log.Fields) string {
	// Formatted time with time layout "<YYYY>/<MM>/<DD>T<HH>:<mm>:<ss>.<sss>Z"
	formattedTime := time.UTC().Format("2006/01/02T15:04:05.000Z")

	// Convert level into string
	levelString := level.String()

	if len(fields) > 0 {
		var b strings.Builder
		for _, field := range fields {
			b.WriteString(fmt.Sprintf(" %s=%v", field.Name, field.Value))
		}

		if len(message) > 0 {
			return fmt.Sprintf("[%s] [%s] %s%s", formattedTime, levelString, message, b.String())
		}

		return fmt.Sprintf("[%s] [%s]%s", formattedTime, levelString, b.String())
	}

	// Combined formatted time, level string, and message
	return fmt.Sprintf("[%s] [%s] %s", formattedTime, levelString, message)
}
