package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Load snapshots the process environment once. A local .env supplies missing
// values only in development. It never mutates the process environment.
func Load() (*Config, error) {
	values := map[string]string{}
	for _, entry := range os.Environ() {
		key, value, ok := strings.Cut(entry, "=")
		if ok {
			values[key] = value
		}
	}
	return loadWithFile(values, ".env")
}
func loadWithFile(values map[string]string, path string) (*Config, error) {
	// Explicit production/test/unknown environments never read local dotenv.
	if env, present := values["APP_ENV"]; !present || env == string(Development) {
		f, err := os.Open(path)
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf(".env: cannot read configuration file")
		}
		if err == nil {
			defer f.Close()
			sc := bufio.NewScanner(f)
			for sc.Scan() {
				line := strings.TrimSpace(sc.Text())
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				key, value, ok := splitEnvLine(line)
				if !ok {
					return nil, fmt.Errorf(".env: invalid assignment")
				}
				if _, present := values[key]; !present {
					values[key] = value
				}
			}
			if sc.Err() != nil {
				return nil, fmt.Errorf(".env: cannot read configuration file")
			}
		}
	}
	return Parse(values)
}
func splitEnvLine(line string) (string, string, bool) {
	key, value, ok := strings.Cut(strings.TrimPrefix(line, "export "), "=")
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	if !ok || key == "" {
		return "", "", false
	}
	for _, c := range key {
		if !(c == '_' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9') {
			return "", "", false
		}
	}
	if len(value) > 0 && (value[0] == '\'' || value[0] == '"') {
		if len(value) < 2 || value[len(value)-1] != value[0] {
			return "", "", false
		}
		value = value[1 : len(value)-1]
	}
	return key, value, true
}
