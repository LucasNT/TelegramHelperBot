package bittorrentapi

import "context"

// import "io"

type BitTorrentClient interface {
	AddTorrentByMagnet(ctx context.Context, magnetLink string) (bool, error)
	// AddTorrentByTorrentFile(file io.Reader) (bool, error)
	Close(ctx context.Context) error
}
