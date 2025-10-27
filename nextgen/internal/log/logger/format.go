package logger

import (
	"fmt"
	"time"
	"update-sh/nextgen/internal/log"
)

type Formatter interface {
	Format(time time.Time, level log.Level, message string, fields log.Fields) string
}

type CustomFormatter struct{}

func NewCustomFormatter() *CustomFormatter {
	return &CustomFormatter{}
}

// Format function
func (s *CustomFormatter) Format(time time.Time, level log.Level, message string, fields log.Fields) string {
	// Formatted time with time layout "YYYY/MM/DDTHH:mm:ss.sssZ"
	formattedTime := time.UTC().Format("2006/01/02T15:04:05.000Z")

	// Convert level into string
	levelString := level.String()

	// Combined formatted time, level string, and message
	return fmt.Sprintf("[%s] [%s] %s", formattedTime, levelString, message)
}
