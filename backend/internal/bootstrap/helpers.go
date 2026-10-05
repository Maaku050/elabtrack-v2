package bootstrap

import (
	"strconv"
	"strings"
)

// bodyLimitBytes converts a human-readable body limit (e.g. "1MB", "512KB")
// into bytes. Defaults to 1MB on parse failure.
func bodyLimitBytes(s string) int {
	t := strings.TrimSpace(strings.ToUpper(s))
	if t == "" {
		return 1 << 20
	}
	mult := 1
	num := t
	switch {
	case strings.HasSuffix(t, "KB"):
		mult = 1 << 10
		num = strings.TrimSuffix(t, "KB")
	case strings.HasSuffix(t, "MB"):
		mult = 1 << 20
		num = strings.TrimSuffix(t, "MB")
	case strings.HasSuffix(t, "GB"):
		mult = 1 << 30
		num = strings.TrimSuffix(t, "GB")
	}
	n, err := strconv.Atoi(strings.TrimSpace(num))
	if err != nil || n <= 0 {
		return 1 << 20
	}
	return n * mult
}
