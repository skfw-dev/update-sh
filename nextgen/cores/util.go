package cores

import (
	"bytes"
	"fmt"
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

// EscapeBytes function escapes a byte slice into a string with hexadecimal escape sequences.
func EscapeBytes(data []byte) string {
	var buf bytes.Buffer
	buf.WriteString(`"`)
	for _, b := range data {
		// Write each byte as a two-digit hex escape
		digit := fmt.Sprintf(`\x%02x`, b)
		buf.WriteString(digit)
	}
	buf.WriteString(`"`)
	return buf.String()
}
