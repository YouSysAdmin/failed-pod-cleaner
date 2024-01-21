package osenv

import (
	"os"
)

// GetEnv get os environment value
func GetEnv(name string, defaults string) string {
	env := os.Getenv(name)
	if env != "" {
		return env
	} else {
		return defaults
	}
}
