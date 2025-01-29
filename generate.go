package crossbot

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/go-telegram/bot/models"
	"github.com/itschip/guildedgo"
)

func (m *Message) Discord() *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Content: m.Content,
		Embeds: []*discordgo.MessageEmbed{{
			URL:         m.URL,
			Title:       m.Title,
			Description: m.Description,
			Color:       m.Color,
			Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: m.ThumbnailImageURL},
			Footer:      &discordgo.MessageEmbedFooter{Text: m.Footer.Text, IconURL: m.Footer.IconURL},
			Fields: func() (res []*discordgo.MessageEmbedField) {
				for _, f := range m.Fields {
					if len(f.Value) > 1024 {
						f.Value = f.Value[:1024]
					}

					res = append(res, &discordgo.MessageEmbedField{Name: f.Name, Value: f.Value, Inline: f.Inline})
				}

				return res
			}(),
		}},
		Components: func() (res []discordgo.MessageComponent) {
			for _, r := range m.Buttons {
				var row discordgo.ActionsRow
				for _, b := range r {
					id := b.Callback.Register()
					CallbackCache[id] = b.Callback

					button := discordgo.Button{
						Style:    discordgo.PrimaryButton,
						Label:    b.Label,
						CustomID: id,
						Emoji:    &discordgo.ComponentEmoji{Name: b.Emoji},
					}
					row.Components = append(row.Components, button)
				}
				res = append(res, row)
			}

			return res
		}(),
	}
}

func (m *Message) Guilded() *guildedgo.MessageObject {
	return &guildedgo.MessageObject{
		Content: m.Content,
		Embeds: []guildedgo.ChatEmbed{{
			Title:       m.Title,
			Description: m.Description,
			URL:         m.URL,
			Color:       m.Color,
			Footer: guildedgo.ChatEmbedFooter{
				Text:    m.Footer.Text,
				IconURL: m.Footer.IconURL,
			},
			Thumbnail: guildedgo.ChatEmbedThumbnail{
				URL: m.ThumbnailImageURL,
			},
			Fields: func() (res []guildedgo.ChatEmbedField) {
				for _, f := range m.Fields {
					res = append(res, guildedgo.ChatEmbedField{Name: f.Name, Value: f.Value, Inline: f.Inline})
				}

				return res
			}(),
		}},
	}
}

func (m *Message) Telegram() (text string, markup models.ReplyMarkup) {
	var sb strings.Builder

	// Add Content
	if m.Content != "" {
		sb.WriteString(m.Content + "\n\n")
	}

	// Add Title and URL
	if m.URL != "" {
		sb.WriteString(fmt.Sprintf("[%s](%s)\n", m.Title, m.URL))
	} else {
		sb.WriteString(m.Title + "\n")
	}

	// Add Description
	if m.Description != "" {
		sb.WriteString(m.Description + "\n\n")
	}

	// Add Fields
	if len(m.Fields) > 0 {
		for _, v := range m.Fields {
			sb.WriteString(fmt.Sprintf("**%s**\n%s\n\n", v.Name, v.Value))
		}
	}

	// Add Footer
	if m.Footer.Text != "" {
		sb.WriteString(m.Footer.Text)
	}

	// Remove any unnecessary whitespace
	text = strings.TrimSpace(sb.String())

	// Parse buttons as valid markup
	kb := models.InlineKeyboardMarkup{}
	for _, r := range m.Buttons {
		var row []models.InlineKeyboardButton
		for _, b := range r {
			id := b.Callback.Register()
			CallbackCache[id] = b.Callback

			var button models.InlineKeyboardButton
			if b.Callback.Action != CallbackActionPrompt {
				button = models.InlineKeyboardButton{
					Text:         fmt.Sprintf("%s %s", b.Emoji, b.Label),
					CallbackData: id,
				}
			} else {
				prompt := fmt.Sprintf("%s\n\n", b.Callback.Prompt.Prefix)
				for _, p := range b.Callback.Prompt.Fields {
					key := strings.ReplaceAll(strings.ToLower(p.Key), " ", "_")
					value := strings.ReplaceAll(strings.ToLower(p.Value), " ", "_")
					prompt += fmt.Sprintf("--%s=\"%s\"\n", key, value)
				}

				button = models.InlineKeyboardButton{
					Text:                         fmt.Sprintf("%s %s", b.Emoji, b.Label),
					SwitchInlineQueryCurrentChat: prompt,
				}
			}

			row = append(row, button)
		}
		kb.InlineKeyboard = append(kb.InlineKeyboard, row)
	}

	if len(kb.InlineKeyboard) == 0 {
		markup = nil
	} else {
		markup = models.ReplyMarkup(kb)
	}

	return text, markup
}

func (cb Callback) Run(user string, platform Platform) *Message {
	fields, err := cb.ParseFields(user, platform)
	if err != nil {
		return &Message{Title: err.Error()}
	}

	return cb.Function(fields)
}

func (cb Callback) ParseFields(user string, platform Platform) (map[string]string, error) {
	var fields map[string]string
	if err := json.Unmarshal([]byte(cb.Fields), &fields); err != nil {
		return nil, err
	}

	fields["user"] = user
	fields["platform"] = fmt.Sprint(platform)

	return fields, nil
}

func (cb Callback) Register() string {
	index := fmt.Sprint(len(CallbackCache))
	CallbackCache[index] = cb
	return index
}
