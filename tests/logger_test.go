package tests

import (
	"testing"
	"update-sh/nextgen/log"
	"update-sh/nextgen/log/logger"
)

func TestLogger(t *testing.T) {
	config := logger.NewConfig(log.InfoLevel, log.ModeConsole, "", 0)
	customLogger := logger.NewZapLogger(config)

	customLogger.
		Field("Name", "John, Doe").
		Field("Age", 20).
		Printf("Test: %s\n\n", "Hello World!")
}
