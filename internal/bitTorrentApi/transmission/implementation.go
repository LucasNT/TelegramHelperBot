package transmission

import (
	"context"
	"encoding/base64"
	"io"
	"log/slog"
	"strings"

	"github.com/hekmon/transmissionrpc/v3"
)

func (self *transmissionBitTorrent) AddTorrentByMagnet(ctx context.Context, magnetLink string) (bool, error) {

	payload := transmissionrpc.TorrentAddPayload{}
	payload.Filename = &magnetLink

	torrent, err := self.transmission.TorrentAdd(ctx, payload)
	if err != nil {
		slog.Default().Error("Failed to add torrent in transmission", "error", err)
		return false, err
	}

	slog.Info("Torrent added to the transimssion", "name", *torrent.Name, "ID", *torrent.ID)

	return true, nil
}

func (self *transmissionBitTorrent) AddTorrentByTorrentFile(ctx context.Context, file io.Reader) (bool, error) {
	base64StringBuilder := strings.Builder{}
	encoder := base64.NewEncoder(base64.StdEncoding, &base64StringBuilder)
	_, err := io.Copy(encoder, file)
	if err != nil {
		slog.Error("Failed to add torrent file in transmission", "error", err)
	}
	encoder.Close()
	metaDataInfo := base64StringBuilder.String()

	payload := transmissionrpc.TorrentAddPayload{}
	payload.MetaInfo = &metaDataInfo

	torrent, err := self.transmission.TorrentAdd(ctx, payload)
	if err != nil {
		slog.Default().Error("Failed to add torrent file in transmission", "error", err)
		return false, err
	}

	slog.Info("Torrent added to the transimssion", "name", *torrent.Name, "ID", *torrent.ID)

	return true, nil
}

func (*transmissionBitTorrent) Close(ctx context.Context) error {
	return nil
}
