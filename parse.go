package crossbot

import (
	"strings"
)

// parseFields parses a string message into a map based on positional arguments and flags
func (c *Config) parseFields(message string, cmd TextCommand) map[string]string {
	result := make(map[string]string)

	// Handle empty delimiter case
	if cmd.SplitFieldsOn == nil || *cmd.SplitFieldsOn == "" {
		return map[string]string{"": message}
	}

	// Track the current position in argsOrder
	var argIndex int

	// Split message into parts using the specified delimiter
	parts := strings.Split(message, *cmd.SplitFieldsOn)

	// Iterate through all parts of the split message
	for i := 0; i < len(parts); i++ {
		// Clean up whitespace and skip empty parts
		part := strings.TrimSpace(parts[i])
		if part == "" {
			continue
		}

		// Handle quoted strings that might span multiple parts
		if strings.HasPrefix(part, `"`) {
			// Initialize string builder for efficient concatenation
			var quotedValue strings.Builder
			quotedValue.WriteString(part[1:])

			// Continue reading parts until we find the closing quote
			for !strings.HasSuffix(part, `"`) && i+1 < len(parts) {
				i++
				part = parts[i]
				quotedValue.WriteString(*cmd.SplitFieldsOn + part)
			}

			// Extract the final value and remove trailing quote
			value := quotedValue.String()
			if strings.HasSuffix(value, `"`) {
				value = value[:len(value)-1]
			}

			// Store as positional argument if we haven't exhausted argsOrder
			if argIndex < len(cmd.Arguments) {
				result[cmd.Arguments[argIndex]] = value
				argIndex++
			}
			continue
		}

		// Handle flags (--flag) and key-value pairs (--key=value)
		if strings.HasPrefix(part, "--") || strings.HasPrefix(part, "-") {
			// Remove leading dashes
			part = strings.TrimLeft(part, "-")

			// Split on first equals sign if present
			if keyValue := strings.SplitN(part, "=", 2); len(keyValue) == 2 {
				value := keyValue[1]
				// Remove surrounding quotes if present
				if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
					value = value[1 : len(value)-1]
				}
				result[keyValue[0]] = value
			} else {
				// Handle flag without value
				result[part] = ""
			}
			continue
		}

		// Handle regular positional arguments
		if argIndex < len(cmd.Arguments) {
			result[cmd.Arguments[argIndex]] = part
			argIndex++
		}
	}

	return result
}
