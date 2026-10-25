package telegrambot

import (
	"context"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var magnetRE = regexp.MustCompile(`^(?i)magnet:\?(?:[^\s&]*&)*?xt=urn:btih:([a-f0-9]{40}|[a-z2-7]{32})(?:&[^\s]*)?$`)

func magnetLinkIdentifierHanler(ctx context.Context, tbt *bot.Bot, update *models.Update) {
	messageParams := &bot.SendMessageParams{}
	if update.Message != nil {
		messageParams.ChatID = update.Message.Chat.ID

		ok, err := btClient.AddTorrentByMagnet(ctx, update.Message.Text)

		if err != nil || !ok {
			slog.Default().Error("Failed to add torrent", "error", err)
			messageParams.Text = "Failed to add torrent"
		}

		if ok {
			messageParams.Text = "Torrent was added"
		}

		tbt.SendMessage(ctx, messageParams)
	}
}

func torrentFileSendHandler(ctx context.Context, tbt *bot.Bot, update *models.Update) {
	log := slog.With("command", "Torrent File Handler")
	if update.Message != nil && update.Message.Document != nil {
		if update.Message.Document.MimeType == "application/x-bittorrent" {
			doc := update.Message.Document
			params := &bot.GetFileParams{}
			params.FileID = doc.FileID
			modelFile, err := tbt.GetFile(ctx, params)
			if err != nil {
				log.Error("Failed to add torrent file", "error", err)
				return
			}

			downloadLink := tbt.FileDownloadLink(modelFile)
			resp, err := http.Get(downloadLink)
			if err != nil {
				log.Error("Faeild to add torrent file", "error", err)
				return
			}
			defer resp.Body.Close()

			ok, err := btClient.AddTorrentByTorrentFile(ctx, resp.Body)
			if err != nil {
				log.Error("Failed to add torrent file", "error", err)
				return
			}

			if !ok {
				log.Error("Failed to add torrent file")
				return
			}
			return
		}
	}
	log.Error("Failed to add torrent file, because the message doesn't have a torrent file")
}
