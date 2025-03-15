package config

import (
	"fmt"
	"strings"
)

func flagNameFromConfigKey(key string) string {
	return strings.ReplaceAll(key, ".", "_")
}

func getFlagUsage(key string, usage string) string {
	return fmt.Sprintf("Env: %s\n\t\t%s", envNameFromConfigKey(key), usage)
}
func envNameFromConfigKey(key string) string {
	return envPrefix + "_" + strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
}
