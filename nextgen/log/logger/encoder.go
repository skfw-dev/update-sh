package logger

import (
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// CustomConsoleEncoder uses the provided Formatter interface to render logs.
type CustomConsoleEncoder struct {
	zapcore.Encoder
	formatter Formatter
}

// NewCustomConsoleEncoder creates a new CustomConsoleEncoder.
// Wrap a basic console encoder to handle Time formatting and basic fields
func NewCustomConsoleEncoder(formatter Formatter, encoderConfig zapcore.EncoderConfig) zapcore.Encoder {
	return &CustomConsoleEncoder{
		Encoder:   zapcore.NewConsoleEncoder(encoderConfig),
		formatter: formatter,
	}
}

func (s *CustomConsoleEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buff := buffer.NewPool().Get()

	// Call the user-provided Formatter method
	output := s.formatter.Format(entry.Time, toLevel(entry.Level), entry.Message, toFields(fields))

	// Append string
	buff.AppendString(output)
	buff.AppendByte('\n')
	return buff, nil
}

func (s *CustomConsoleEncoder) Clone() zapcore.Encoder {
	return &CustomConsoleEncoder{
		Encoder:   s.Encoder.Clone(),
		formatter: s.formatter,
	}
}
