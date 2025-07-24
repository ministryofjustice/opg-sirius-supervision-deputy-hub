package util

import (
	"strings"
)

func StringToArray(input string) []string {
	items := strings.Split(input, ",")
	var result []string

	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
