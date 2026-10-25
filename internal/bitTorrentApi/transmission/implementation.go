package transmission

import (
	"context"
	"log/slog"

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

func (*transmissionBitTorrent) Close(ctx context.Context) error {
	return nil
}
