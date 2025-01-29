package crossbot

import (
	"regexp"
	"strings"
)

func parseFields(message string, argsOrder []string) map[string]string {
	result := make(map[string]string)
	argCount := 0

	// Split message on each new word, quoted phrase, or key-value pair
	re := regexp.MustCompile(`--?[^=\s]+(?:=(?:"(?:\\.|[^"\\])*"|[^"\s]+))?|"(?:\\.|[^"\\])*"|\S+`)
	fields := re.FindAllString(message, -1)[1:]

	for _, field := range fields {
		if strings.HasPrefix(field, "--") || strings.HasPrefix(field, "-") {
			field = strings.TrimLeft(field, "-")
			if strings.Contains(field, "=") {
				parts := strings.SplitN(field, "=", 2)
				key := parts[0]
				value := parts[1]

				// Remove surrounding quotes if present
				if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
					value = value[1 : len(value)-1]
				}

				// Unescape any escaped quotes within the value
				value = strings.ReplaceAll(value, `\"`, `"`)

				result[key] = value
			} else {
				result[field] = ""
			}
		} else {
			// Remove surrounding quotes if present
			if strings.HasPrefix(field, `"`) && strings.HasSuffix(field, `"`) {
				field = field[1 : len(field)-1]
			}

			// Unescape any escaped quotes within the field
			field = strings.ReplaceAll(field, `\"`, `"`)

			if argCount < len(argsOrder) {
				key := argsOrder[argCount]
				result[key] = field
				argCount++
			}
		}
	}

	return result
}
