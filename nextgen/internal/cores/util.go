package cores

import (
	"os"
	"strings"
)

const AppEnv = "APP_ENV"

func IsDevelopment() bool {
	if appEnv := os.Getenv(AppEnv); appEnv != "" {
		return strings.EqualFold(strings.TrimSpace(appEnv), string(Development))
	}

	return false
}

func IsProduction() bool {
	if appEnv := os.Getenv(AppEnv); appEnv != "" {
		return strings.EqualFold(strings.TrimSpace(appEnv), string(Production))
	}

	return false
}

func GetAppEnv() Env {
	if appEnv := os.Getenv(AppEnv); appEnv != "" {
		switch appEnv = strings.ToLower(strings.TrimSpace(appEnv)); Env(appEnv) {
		case Development:
			return Development

		case Production:
			return Production

		default:
			panic("invalid app env")
		}
	}

	return Development
}
