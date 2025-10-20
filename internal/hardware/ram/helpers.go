package ram

import (
	"strconv"
	"strings"
)

func cutPrefixString(line, prefix string, field *string) {
	if after, ok := strings.CutPrefix(line, prefix); ok {
		*field = strings.TrimSpace(after)
	}
}

func parseIntField(line, prefix string, field *int, mult int, suffixes ...string) bool {
	if after, ok := strings.CutPrefix(line, prefix); ok {
		val := strings.TrimSpace(after)
		for _, s := range suffixes {
			val = strings.TrimSuffix(val, s)
		}
		val = strings.TrimSpace(val)

		if val == "" || val == "Unknown" || val == "0" || strings.Contains(val, "No Module") {
			return false
		}

		num, err := strconv.Atoi(val)
		if err == nil {
			*field = num * mult
			return true
		}
	}
	return false
}
