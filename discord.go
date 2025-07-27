package crossbot

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func (c *Config) Discord(cmds *[]*Command) (cancelFunc func() error, err error) {
	dg, err := discordgo.New("Bot " + c.DiscordConfig.BotToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create new Discord session: %w", err)
	}

	dg.Identify.Intents = c.DiscordConfig.Intents

	// Register all commands to Discord
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.Ready) {
		if err = c.RegisterDiscord(s, cmds); err != nil {
			log.Println("Failed to register Discord commands:", err)
		}
	})

	// Set bot status
	dg.AddHandler(func(s *discordgo.Session, event *discordgo.Ready) {
		s.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{{
				Name: c.DiscordConfig.BotActivityMessage,
				Type: c.DiscordConfig.BotActivityType,
			}},
		})
	})

	// Logger
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		log.Printf("%s (%s): %s", m.Author.Username, m.Author.ID, m.Content)
	})

	for _, cmd := range *cmds {
		if cmd.Discord.TextMiddleware != nil {
			dg.AddHandler(cmd)
		}
	}

	if err = dg.Open(); err != nil {
		return nil, fmt.Errorf("failed to establish Discord connection: %w", err)
	}

	return dg.Close, nil
}

// RegisterDiscord registers all provided commands with Discord
func (c *Config) RegisterDiscord(s *discordgo.Session, cmds *[]*Command) error {
	commandHandlers := make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))

	var dcmds []*discordgo.ApplicationCommand

	// Prepare commands for bulk overwrite
	for _, cmd := range *cmds {
		if cmd.Discord.TextMiddleware != nil {
			continue
		}

		cmdCpy := cmd
		dcmd := cmdCpy.Discord.ApplicationCommand
		dcmds = append(dcmds, &dcmd)
	}

	_, err := s.ApplicationCommandBulkOverwrite(c.DiscordConfig.ApplicationID, "", dcmds)
	if err != nil {

		rcmds, err := s.ApplicationCommands(c.DiscordConfig.ApplicationID, "")
		if err != nil {
			return fmt.Errorf("failed to overwrite commands: %w", err)
		}

		rcmdMap := make(map[*discordgo.ApplicationCommand]bool)
		for _, rcmd := range rcmds {
			rcmdMap[rcmd] = true
		}

		// Find elements in cmds that are not in rcmds
		var newCmds []*discordgo.ApplicationCommand
		for _, cmd := range *cmds {
			if !rcmdMap[&cmd.Discord.ApplicationCommand] {
				newCmds = append(newCmds, &cmd.Discord.ApplicationCommand)
			}
		}

		for _, ncmd := range newCmds {
			if _, err = s.ApplicationCommandCreate(c.DiscordConfig.ApplicationID, "", ncmd); err != nil {
				return fmt.Errorf("failed to create command: %w", err)
			}
		}
	}

	for _, cmd := range *cmds {
		if cmd.Discord.TextMiddleware != nil {
			continue
		}

		cmdCpy := cmd
		dcmd := cmdCpy.Discord.ApplicationCommand

		commandHandlers[dcmd.Name] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			fields := make(map[string]string)
			fields["user"] = i.User.Username
			fields["platform"] = fmt.Sprint(PlatformDiscord)

			for _, opt := range i.ApplicationCommandData().Options {
				fields[opt.Name] = opt.StringValue()
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: cmdCpy.Handler(fields).Discord(),
			})
			if err != nil {
				log.Println("Failed to respond to Discord interaction:", err)
			}
		}
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:
			id := i.Interaction.MessageComponentData().CustomID
			cb, ok := CallbackCache[id]
			if !ok {
				return
			}

			user := i.User.Username
			msg := cb.Run(user, PlatformDiscord)
			resp := msg.Discord()

			switch cb.Action {
			case CallbackActionEditMessage:
				s.ChannelMessageEditComplex(&discordgo.MessageEdit{
					ID:         i.Message.ID,
					Channel:    i.ChannelID,
					Embeds:     &resp.Embeds,
					Components: &resp.Components,
				})

			case CallbackActionCreateMessage:
				s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
					Embeds:     resp.Embeds,
					Components: resp.Components,
				})

			case CallbackActionDeleteMessage:
				s.ChannelMessageDelete(i.ChannelID, i.Message.ID)
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredMessageUpdate,
			})

		case discordgo.InteractionModalSubmit:
		}
	})

	return nil
}
