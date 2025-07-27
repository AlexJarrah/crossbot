package crossbot

import (
	"errors"
	"fmt"

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

// Validate ensures that the config is not missing any necessary fields.
// It also populates unspecified fields with default values.
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
	}

	if c.CacheDirectory == "" {
		dir, err := c.DefaultCacheDirectory()
		if err != nil {
			return fmt.Errorf("failed to populate missing field 'CacheDirectory': %w", err)
		}
		c.CacheDirectory = dir
	}

	return nil
}
