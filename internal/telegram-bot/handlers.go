package telegrambot

import (
	"context"
	"regexp"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var magnetRE = regexp.MustCompile(`^(?i)magnet:\?(?:[^\s&]*&)*?xt=urn:btih:([a-f0-9]{40}|[a-z2-7]{32})(?:&[^\s]*)?$`)

func magnetLinkIdentifierHanler(ctx context.Context, tbt *bot.Bot, update *models.Update) {
	messageParams := &bot.SendMessageParams{}
	if update.Message != nil {
		messageParams.ChatID = update.Message.Chat.ID
		messageParams.Text = "It's a magnet Link"
		tbt.SendMessage(ctx, messageParams)
	}
}
