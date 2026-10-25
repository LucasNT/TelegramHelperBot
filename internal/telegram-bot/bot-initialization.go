package telegrambot

import (
	"log/slog"

	"github.com/LucasNT/TorrentHelperBot/internal/config"
	"github.com/go-telegram/bot"
)

func Initialization() (*bot.Bot, error) {
	opts := []bot.Option{
		bot.WithMiddlewares(allowUsers(config.DefaultConfig.AllowedUsersID...)),
	}
	b, err := bot.New(config.DefaultConfig.TelegramBotToken, opts...)
	if err != nil {
		slog.Default().Error("Failed to create Bot", "error", err)
		return nil, err
	}
	b.RegisterHandlerRegexp(bot.HandlerTypeMessageText, magnetRE, magnetLinkIdentifierHanler)
	return b, err
}
