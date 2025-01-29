package crossbot

import (
	"fmt"
	"strings"

	"github.com/go-telegram/bot/models"
)

// Validate & run specified command
func (c Command) Run(user, message, command string, platform Platform) (text string, markup models.ReplyMarkup) {
	fields := parseFields(message, c.Text.Arguments)
	fields["user"] = user
	fields["platform"] = fmt.Sprint(platform)

	for _, a := range c.Text.Arguments {
		if _, ok := fields[a]; !ok {
			return fmt.Sprintf("Invalid usage!\n\nex: /%s %s\n\noptionally: %s", command, strings.Join(c.Text.Arguments, " "), strings.Join(c.Text.Options, ", ")), nil
		}
	}

	return c.Handler(fields).Telegram()
}
