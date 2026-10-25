package telegrambot

import (
	"context"
	"log/slog"

	bittorrentapi "github.com/LucasNT/TorrentHelperBot/internal/bitTorrentApi"
	"github.com/LucasNT/TorrentHelperBot/internal/bitTorrentApi/transmission"
	"github.com/LucasNT/TorrentHelperBot/internal/config"
	"github.com/go-telegram/bot"
)

var btClient bittorrentapi.BitTorrentClient

func Initialization(ctx context.Context) (*bot.Bot, error) {
	opts := []bot.Option{
		bot.WithMiddlewares(allowUsers(config.DefaultConfig.AllowedUsersID...)),
	}
	b, err := bot.New(config.DefaultConfig.TelegramBotToken, opts...)
	if err != nil {
		slog.Default().Error("Failed to create Bot", "error", err)
		return nil, err
	}
	b.RegisterHandlerRegexp(bot.HandlerTypeMessageText, magnetRE, magnetLinkIdentifierHanler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypePrefix, torrentFileSendHandler)

	btClient, err = transmission.New(ctx, config.DefaultConfig.TransmissionRPCURL)
	if err != nil {
		slog.Error("Failed to create Bot", "error", err)
	}

	return b, err
}
