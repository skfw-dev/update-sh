package tests

import (
	"testing"
	"update-sh/nextgen/log"
	"update-sh/nextgen/log/logger"
)

func TestLogger(t *testing.T) {
	config := logger.NewConfig(log.InfoLevel, log.ModeConsole, "", 0)
	formatter := config.Formatter
	formatter.SetLogFmt(false)
	formatter.SetFormatEncoder(logger.TextFormatEncoder)
	customLogger := logger.NewLogger(config)

	for i := 0; i < 10; i++ {
		customLogger.
			Field("Name", "John, Doe").
			Field("Age", 20).
			Print("This is message.")
	}
}
