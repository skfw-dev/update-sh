package tests

import (
	"testing"
	"update-sh/nextgen/log/common"
	"update-sh/nextgen/log/logger"
)

func TestLogger(t *testing.T) {
	config := logger.NewConfig(common.InfoLevel, common.ModeConsole, "", 0)
	formatter := config.Formatter
	formatter.SetTextMode(logger.XMLLogFmt)
	formatter.SetFormatEncoder(logger.TextFormatEncoder)
	customLogger := logger.NewLogger(config)

	for i := 0; i < 10; i++ {
		customLogger.
			Field("Name", "John, Doe").
			Field("Age", 20).
			Field("Blob", []byte{255, 255, 255, 0}).
			Print("This is message.")
	}
}
