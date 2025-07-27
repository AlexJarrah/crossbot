package crossbot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/go-telegram/bot"
)

type (
	// Command metadata
	Command struct {
		// Commands done via text
		Text TextCommand

		// Platform-specific commands
		Telegram TelegramCommand
		Discord  DiscordCommand

		// Function to be ran once command is called.
		Handler func(fields map[string]string) *Message
	}

	// TextCommand is a command's text configuration
	TextCommand struct {
		// Command aliases allow multiple ways to activate the same function
		Aliases []string

		// Required fields that the command must have in the correct order. All
		// arguments specified by the user will have a key and value in the map
		// provided to the handler.
		Arguments []string

		// Optional fields that the command does not require. All options specified by
		// the user will have a key and value in the map provided to the handler.
		Options []string

		// String to search for in the user message to seperate fields on. By default,
		// this is set to " ", meaning any space in the user's message will mark a new
		// field. This can be disabled by setting it to "" to allow the handler to get
		// the entire message without parsing fields (i.e. 'echo hello world').
		SplitFieldsOn *string
	}

	// TelegramCommand is a command's Telegram configuration
	TelegramCommand struct {
		// Function that gets called on every chat message
		TextMiddleware func(next bot.HandlerFunc) bot.HandlerFunc
	}

	// DiscordCommand is a command's Discord configuration
	DiscordCommand struct {
		// Function that gets called on every chat message
		TextMiddleware func(s *discordgo.Session, m *discordgo.MessageCreate)

		// Discord slash command
		ApplicationCommand discordgo.ApplicationCommand
	}
)
