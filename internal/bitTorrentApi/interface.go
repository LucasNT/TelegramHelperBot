package bittorrentapi

import (
	"context"
	"io"
)

// import "io"

type BitTorrentClient interface {
	AddTorrentByMagnet(ctx context.Context, magnetLink string) (bool, error)
	AddTorrentByTorrentFile(ctx context.Context, file io.Reader) (bool, error)
	Close(ctx context.Context) error
}
