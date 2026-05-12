package workflows

import (
	"net/http"
	"strconv"
	"strings"
)

func parsePathUint64(r *http.Request, key string) (uint64, error) {
	raw := strings.TrimSpace(r.PathValue(key))
	if raw == "" {
		raw = fallbackPathValue(r.URL.Path, key)
	}
	return strconv.ParseUint(raw, 10, 64)
}

func parsePathInt64(r *http.Request, key string) (int64, error) {
	raw := strings.TrimSpace(r.PathValue(key))
	if raw == "" {
		raw = fallbackPathValue(r.URL.Path, key)
	}
	return strconv.ParseInt(raw, 10, 64)
}

func fallbackPathValue(path, key string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 {
		return ""
	}

	switch key {
	case "id":
		for i := 0; i < len(parts)-1; i++ {
			if parts[i] == "workflows" {
				return parts[i+1]
			}
		}
	case "version":
		for i := 0; i < len(parts)-1; i++ {
			if parts[i] == "rollback" {
				return parts[i+1]
			}
		}
	case "executionId":
		for i := 0; i < len(parts)-1; i++ {
			if parts[i] == "executions" {
				return parts[i+1]
			}
		}
	}

	return parts[len(parts)-1]
}
