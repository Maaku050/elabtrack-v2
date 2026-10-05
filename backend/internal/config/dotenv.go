package config

import (
	"bufio"
	"os"
	"strings"
)

// LoadDotEnv reads a .env file (if present) from the working directory and
// sets any variables that are not already present in the environment. This
// is a minimal, dependency-free loader modeled on godotenv semantics.
//
// Existing environment variables always win over .env values, so production
// deployments can override behaviour without editing files.
func LoadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		// Missing .env is fine; rely on real environment variables.
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := splitEnvLine(line)
		if !ok {
			continue
		}
		// Don't override existing env vars.
		if _, present := os.LookupEnv(key); present {
			continue
		}
		_ = os.Setenv(key, val)
	}
}

func splitEnvLine(line string) (string, string, bool) {
	// Handle optional `export ` prefix.
	line = strings.TrimPrefix(line, "export ")
	idx := strings.Index(line, "=")
	if idx <= 0 {
		return "", "", false
	}
	key := strings.TrimSpace(line[:idx])
	val := strings.TrimSpace(line[idx+1:])
	val = unquote(val)
	return key, val, true
}

func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
