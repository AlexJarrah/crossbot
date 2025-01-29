package crossbot

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

type (
	Config struct {
		ID   string
		Name string

		TelegramConfig *TelegramConfig
		DiscordConfig  *DiscordConfig

		CacheDirectory string
	}

	TelegramConfig struct {
		BotToken    string
		BotUsername string
	}

	DiscordConfig struct {
		ApplicationID      string
		BotToken           string
		Intents            discordgo.Intent
		BotActivityType    discordgo.ActivityType
		BotActivityMessage string
	}
)

func (c *Config) Validate() error {
	switch {
	case c == nil:
		return errors.New("config undefined")
	case c.ID == "":
		return errors.New("id must be specified")
	case c.Name == "":
		return errors.New("name must be specified")
	case c.TelegramConfig == nil && c.DiscordConfig == nil:
		return errors.New("no platform configuration specified")
	default:
		return nil
	}
}
