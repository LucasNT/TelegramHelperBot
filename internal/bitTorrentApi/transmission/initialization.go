package transmission

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	bittorrentapi "github.com/LucasNT/TorrentHelperBot/internal/bitTorrentApi"
	transmissionrpc "github.com/hekmon/transmissionrpc/v3"
)

type transmissionBitTorrent struct {
	transmission *transmissionrpc.Client
}

func New(ctx context.Context, transmissionUrl *url.URL) (bittorrentapi.BitTorrentClient, error) {
	ret := &transmissionBitTorrent{}
	tbt, err := transmissionrpc.New(transmissionUrl, nil)
	if err != nil {
		slog.Error("Failed to create transmission client", "error", err)
		return nil, err
	}
	ok, serverVersion, serverMinimumVersion, err := tbt.RPCVersion(ctx)
	if err != nil {
		slog.Error("Failed to create transmission client", "error", err)
		return nil, err
	}
	if !ok {
		err = fmt.Errorf("Incompatible version of transmission with this lib remote lib is v%d, library supports v%d and minimum version required is v%d", serverVersion, transmissionrpc.RPCVersion, serverMinimumVersion)
		slog.Default().Error("Failed to create transmission client", "error", err)
		return nil, err
	}
	ret.transmission = tbt
	return ret, nil
}
