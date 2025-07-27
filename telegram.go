package crossbot

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (c *Config) Telegram(cmds *[]*Command) error {
	var middlewares []bot.Middleware
	for _, cmd := range *cmds {
		if cmd.Telegram.TextMiddleware != nil {
			middlewares = append(middlewares, cmd.Telegram.TextMiddleware)
		}
	}
	middlewares = append(middlewares, middlewareLogger)

	opts := []bot.Option{
		bot.WithDefaultHandler(func(ctx context.Context, bot *bot.Bot, update *models.Update) {}),
		bot.WithMiddlewares(middlewares...),
		bot.WithCallbackQueryDataHandler("", bot.MatchTypePrefix, func(ctx context.Context, b *bot.Bot, update *models.Update) {
			if update.Message != nil {
				f := update.Message.From
				log.Printf("%s %s callback: %s", f.FirstName, f.LastName, update.CallbackQuery.Data)
			}

			id := update.CallbackQuery.Data
			cb, ok := CallbackCache[id]
			if !ok {
				return
			}

			user := getUserFromUpdate(update)

			var text string
			var markup models.ReplyMarkup
			if cb.Function != nil {
				msg := cb.Run(user, PlatformTelegram)
				text, markup = msg.Telegram()
			}

			switch cb.Action {
			case CallbackActionEditMessage:
				b.EditMessageText(ctx, &bot.EditMessageTextParams{
					ChatID:          update.CallbackQuery.Message.Message.Chat.ID,
					MessageID:       update.CallbackQuery.Message.Message.ID,
					InlineMessageID: update.CallbackQuery.InlineMessageID,
					Text:            text,
					ReplyMarkup:     markup,
				})

			case CallbackActionCreateMessage:
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.CallbackQuery.Message.Message.Chat.ID,
					Text:   text,
					ReplyParameters: &models.ReplyParameters{
						MessageID:                update.Message.ID,
						ChatID:                   update.Message.Chat.ID,
						AllowSendingWithoutReply: true,
					},
					ReplyMarkup: markup,
				})

			case CallbackActionDeleteMessage:
				b.DeleteMessage(ctx, &bot.DeleteMessageParams{
					ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
					MessageID: update.CallbackQuery.Message.Message.ID,
				})

			case CallbackActionAlert:
				fields, err := cb.ParseFields(user, PlatformTelegram)
				if err != nil {
					return
				}

				b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
					CallbackQueryID: update.CallbackQuery.ID,
					Text:            cb.AlertMessage(fields),
					ShowAlert:       true,
				})

				return
			}

			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				ShowAlert:       false,
			})
		}),
	}

	b, err := bot.New(c.TelegramConfig.BotToken, opts...)
	if err != nil {
		return fmt.Errorf("failed to create new bot instance: %w", err)
	}

	for _, cmd := range *cmds {
		if cmd.Telegram.TextMiddleware != nil {
			continue
		}

		for _, name := range cmd.Text.Aliases {
			cmdCpy, nameCpy := cmd, name

			fn := func(ctx context.Context, b *bot.Bot, update *models.Update) {
				txt := update.Message.Text
				txt = strings.TrimLeft(txt, fmt.Sprintf("@%s", c.TelegramConfig.BotUsername))
				txt = strings.TrimSpace(txt)

				user := getUserFromUpdate(update)
				text, markup := c.Run(cmdCpy, user, txt, nameCpy, PlatformTelegram)
				_, err := b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   text,
					ReplyParameters: &models.ReplyParameters{
						MessageID:                update.Message.ID,
						ChatID:                   update.Message.Chat.ID,
						AllowSendingWithoutReply: true,
					},
					ReplyMarkup: markup,
				})
				if err != nil {
					log.Println("Failed to send Telegram message:", err)
				}
			}

			b.RegisterHandler(bot.HandlerTypeMessageText, "/"+nameCpy, bot.MatchTypeContains, fn)
		}
	}

	b.Start(context.Background())
	return nil
}

func middlewareLogger(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message != nil {
			f := update.Message.From
			log.Printf("%s %s (%s): %s", f.FirstName, f.LastName, f.Username, update.Message.Text)
		}
		next(ctx, b, update)
	}
}

func getUserFromUpdate(update *models.Update) string {
	var from *models.User
	if update.Message != nil {
		from = update.Message.From
	} else {
		from = &update.CallbackQuery.From
	}

	if from.Username == "" {
		return fmt.Sprintf("%s %s", from.FirstName, from.LastName)
	} else {
		return from.Username
	}
}
