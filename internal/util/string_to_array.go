package util

import (
	"strings"
)

func StringToArray(newValue string) []string {
	items := strings.Split(newValue, ",")
	var result []string

	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
