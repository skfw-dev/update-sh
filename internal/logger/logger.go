package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var once sync.Once
var levelNames map[zerolog.Level]string
var consoleWriter zerolog.ConsoleWriter

// init initializes the logger by setting level names and configuring a custom console writer.
// It maps zerolog levels to human-readable strings and defines the format for log output,
// including timestamp, level, and message formatting.
func init() {
	// Initialize level names
	levelNames = map[zerolog.Level]string{
		zerolog.NoLevel:    "NONE",
		zerolog.Disabled:   "DISABLED",
		zerolog.TraceLevel: "TRACE",
		zerolog.DebugLevel: "DEBUG",
		zerolog.InfoLevel:  "INFO",
		zerolog.WarnLevel:  "WARN",
		zerolog.ErrorLevel: "ERROR",
		zerolog.FatalLevel: "FATAL",
		zerolog.PanicLevel: "PANIC",
	}

	// Custom console writer
	consoleWriter = zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006/01/02T15:04:05.000Z",
		FormatLevel: func(i any) string {
			return fmt.Sprintf("[%s]", i)
		},
		FormatMessage: func(i any) string {
			return fmt.Sprintf("%s", i)
		},
		FormatTimestamp: func(i any) string {
			return fmt.Sprintf("[%s]", i)
		},
	}
}

// Init initializes zerolog with custom formatting
func Init(verbose, quiet bool, logFilePath string) {
	once.Do(func() {
		// Custom level formatter
		zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
			if name, ok := levelNames[l]; ok {
				return name
			}
			return l.String()
		}

		// Set global level
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		if verbose {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else if quiet {
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		}

		// Log file setup
		logDir := filepath.Dir(logFilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.Warn().Err(err).Msgf("Failed to create log directory %s, logging only to console.", logDir)
		}

		fileWriter, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Warn().Err(err).Msgf("Failed to open log file %s, logging only to console.", logFilePath)
		}

		if fileWriter != nil {
			multi := io.MultiWriter(consoleWriter, fileWriter)
			log.Logger = zerolog.New(multi).With().Timestamp().Logger()
		} else {
			log.Logger = zerolog.New(consoleWriter).With().Timestamp().Logger()
		}

		if zerolog.GlobalLevel() <= zerolog.DebugLevel {
			log.Logger = log.With().Caller().Logger()
		}

		log.Info().Msgf("Zerolog initialized. Log level: %s", zerolog.GlobalLevel().String())
	})
}

func Log(msg string, args ...any)  { log.Info().Msgf(msg, args...) }
func Logf(msg string, args ...any) { log.Info().Msgf(msg, args...) }

func Info(msg string, args ...any)  { log.Info().Msgf(msg, args...) }
func Infof(msg string, args ...any) { log.Info().Msgf(msg, args...) }

func Debug(msg string, args ...any)  { log.Debug().Msgf(msg, args...) }
func Debugf(msg string, args ...any) { log.Debug().Msgf(msg, args...) }

func Warn(msg string, args ...any)  { log.Warn().Msgf(msg, args...) }
func Warnf(msg string, args ...any) { log.Warn().Msgf(msg, args...) }

func Error(msg string, args ...any)  { log.Error().Msgf(msg, args...) }
func Errorf(msg string, args ...any) { log.Error().Msgf(msg, args...) }

func Fatal(msg string, args ...any)  { log.Fatal().Msgf(msg, args...) }
func Fatalf(msg string, args ...any) { log.Fatal().Msgf(msg, args...) }

func Panic(msg string, args ...any)  { log.Panic().Msgf(msg, args...) }
func Panicf(msg string, args ...any) { log.Panic().Msgf(msg, args...) }
