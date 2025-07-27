package crossbot

import (
	"fmt"
	"strings"

	"github.com/go-telegram/bot/models"
)

// Run validates & runs the specified command
func (c *Config) Run(cmd *Command, user, msg, command string, platform Platform) (text string, markup models.ReplyMarkup) {
	fields := c.parseFields(msg, cmd.Text)
	fields["user"] = user
	fields["platform"] = fmt.Sprint(platform)

	for _, a := range cmd.Text.Arguments {
		if _, ok := fields[a]; !ok {
			msg := fmt.Sprintf("Invalid usage!\n\nex: /%s %s\n\noptionally: %s",
				command,
				strings.Join(cmd.Text.Arguments, " "),
				strings.Join(cmd.Text.Options, ", "),
			)
			return msg, nil
		}
	}

	return cmd.Handler(fields).Telegram()
}
