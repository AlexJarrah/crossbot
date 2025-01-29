package crossbot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/go-telegram/bot"
)

// Command metadata
type Command struct {
	// Commands done via text
	Text TextCommand

	// Platform-specific commands
	Telegram TelegramCommand
	Discord  DiscordCommand

	// Function to be ran once command is called.
	Handler func(fields map[string]string) *Message
}

type TelegramCommand struct {
	TextMiddleware func(next bot.HandlerFunc) bot.HandlerFunc
}

type DiscordCommand struct {
	TextMiddleware     func(s *discordgo.Session, m *discordgo.MessageCreate)
	ApplicationCommand discordgo.ApplicationCommand
}

type TextCommand struct {
	Names     []string
	Arguments []string
	Options   []string
}
