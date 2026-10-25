package telegrambot

import (
	"context"
	"log/slog"
	"slices"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Sender struct {
	UserID   int64
	Username string
	ChatID   int64
	ChatType models.ChatType
}

func allowUsers(ids ...int64) bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			s, ok := senderOf(update)
			if !ok || !slices.Contains(ids, s.UserID) {
				slog.Warn("update rejected",
					"user", s.UserID,
					"username", s.Username,
					"chat", s.ChatID,
				)
				return
			}

			next(ctx, b, update)
		}
	}
}

func senderOf(update *models.Update) (Sender, bool) {
	for _, msg := range []*models.Message{
		update.Message,
		update.EditedMessage,
		update.ChannelPost,
		update.EditedChannelPost,
		update.BusinessMessage,
		update.EditedBusinessMessage,
	} {
		if msg != nil {
			return senderOfMessage(msg), true
		}
	}

	if cq := update.CallbackQuery; cq != nil {
		s := Sender{UserID: cq.From.ID, Username: cq.From.Username}
		// The originating message is absent for callbacks on inline messages.
		if m := cq.Message.Message; m != nil {
			s.ChatID, s.ChatType = m.Chat.ID, m.Chat.Type
		} else if m := cq.Message.InaccessibleMessage; m != nil {
			s.ChatID, s.ChatType = m.Chat.ID, m.Chat.Type
		}
		return s, true
	}

	if iq := update.InlineQuery; iq != nil && iq.From != nil {
		return Sender{UserID: iq.From.ID, Username: iq.From.Username}, true
	}

	return Sender{}, false
}

func senderOfMessage(msg *models.Message) Sender {
	s := Sender{ChatID: msg.Chat.ID, ChatType: msg.Chat.Type}

	switch {
	case msg.From != nil:
		s.UserID, s.Username = msg.From.ID, msg.From.Username
	case msg.SenderChat != nil:
		// Channel posts and anonymous group admins have no user, only a sender chat.
		s.Username = msg.SenderChat.Username
	}

	return s
}
