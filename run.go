package crossbot

import (
	"fmt"
	"strings"

	"github.com/go-telegram/bot/models"
)

// Run validates & runs the specified command
func (c *Config) Run(cmd *Command, user, msg, command string, platform Platform) (text string, markup models.ReplyMarkup) {
	msg = strings.TrimPrefix(msg, "/"+command)
	fields := c.parseFields(msg, cmd.Text)
	fields["user"] = user
	fields["platform"] = fmt.Sprint(platform)

	for _, a := range cmd.Text.Arguments {
		if _, ok := fields[a]; ok {
			continue
		}

		var parts []string

		if len(cmd.Text.Arguments) > 0 {
			required := fmt.Sprintf("Required: %s", strings.Join(cmd.Text.Arguments, " | "))
			parts = append(parts, required)
		}

		if len(cmd.Text.Options) > 0 {
			optional := fmt.Sprintf("Optional: %s", strings.Join(cmd.Text.Options, " | "))
			parts = append(parts, optional)
		}

		exampleParts := []string{fmt.Sprintf("/%s", cmd.Text.Aliases[0])}
		exampleParts = append(exampleParts, cmd.Text.Arguments...)
		for _, opt := range cmd.Text.Options {
			exampleParts = append(exampleParts, fmt.Sprintf("--%s=value", opt))
		}

		example := fmt.Sprintf("Example: `%s`", strings.Join(exampleParts, " "))
		parts = append(parts, example)

		msg := strings.Join(parts, "\n")
		return msg, nil
	}

	return cmd.Handler(fields).Telegram()
}
