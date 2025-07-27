package crossbot

import (
	"strings"
)

// parseFields parses a string message into a map based on positional arguments and flags
func (c *Config) parseFields(message string, cmd TextCommand) map[string]string {
	result := make(map[string]string)

	var delimiter string
	if cmd.SplitFieldsOn == nil {
		// Fallback to space
		delimiter = " "
	} else {
		delimiter = *cmd.SplitFieldsOn
	}

	// If delimiter is empty, return the entire message under an empty key
	if delimiter == "" {
		return map[string]string{"": message}
	}

	// Split message into parts using the delimiter
	parts := strings.Split(message, delimiter)

	// Track the current position in Arguments
	var argIndex int

	// Iterate through all parts of the split message
	for i := 0; i < len(parts); i++ {
		// Clean up whitespace and skip empty parts
		part := strings.TrimSpace(parts[i])
		if part == "" {
			continue
		}

		// Handle quoted strings that might span multiple parts
		if strings.HasPrefix(part, `"`) {
			var quotedValue strings.Builder
			quotedValue.WriteString(part[1:])

			// Continue reading parts until we find the closing quote
			for !strings.HasSuffix(part, `"`) && i+1 < len(parts) {
				i++
				part = parts[i]
				quotedValue.WriteString(delimiter + part)
			}

			// Extract the final value and remove trailing quote
			value := quotedValue.String()
			if strings.HasSuffix(value, `"`) {
				value = value[:len(value)-1]
			}

			// Store as positional argument if we haven't exhausted Arguments
			if argIndex < len(cmd.Arguments) {
				result[cmd.Arguments[argIndex]] = value
				argIndex++
			}
			continue
		}

		// Handle flags and key-value pairs (--key=value)
		if strings.HasPrefix(part, "--") || strings.HasPrefix(part, "-") || strings.HasPrefix(part, "—") {
			// Remove leading dashes
			flag := strings.TrimLeft(part, "-")
			flag = strings.TrimLeft(flag, "—")

			// Check if the flag (with or without value) matches an option
			for _, opt := range cmd.Options {
				// Check if flag starts with the option name
				if flag == opt || strings.HasPrefix(flag, opt+"=") {
					// Store the full flag (without dashes) as the key with empty value
					result[flag] = ""
					break
				}
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
