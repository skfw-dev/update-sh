package logger

import (
	"time"

	"go.skfw.net/shikalog/common"
)

type Formatter interface {
	Format(time time.Time, level common.Level, message string, fields common.Fields) string
}
