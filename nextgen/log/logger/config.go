package logger

import (
	"update-sh/nextgen/cores"
	"update-sh/nextgen/log"
)

type Config struct {
	Level      log.Level `json:"level" yaml:"level"`
	Mode       log.Mode  `json:"mode" yaml:"mode"`
	Filename   string    `json:"filename" yaml:"filename"`
	CallerSkip int       `json:"callerSkip" yaml:"caller_skip"`
	Formatter  Formatter `json:"formatter" yaml:"formatter"`
	AppEnv     cores.Env `json:"appEnv" yaml:"app_env"`
}

func NewConfig(level log.Level, mode log.Mode, filename string, callerSkip int) *Config {
	return &Config{
		Level:      level,
		Mode:       mode,
		Filename:   filename,
		CallerSkip: callerSkip,
		Formatter:  NewCustomFormatter(),
		AppEnv:     cores.GetAppEnv(),
	}
}
