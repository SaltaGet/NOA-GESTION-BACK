package utils

import "strings"

func ParseUsername(raw string) (string, string) {
	parts := strings.SplitN(raw, "@", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return raw, ""
}