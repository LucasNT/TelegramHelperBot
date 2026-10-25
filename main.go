package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/LucasNT/TorrentHelperBot/cmd"
)

func main() {
	var rootContext context.Context
	var rootCancel context.CancelFunc

	h := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(h)

	rootContext, rootCancel = signal.NotifyContext(context.Background(), os.Interrupt)
	defer rootCancel()

	err := cmd.Execute(rootContext)
	if err != nil {
		slog.Default().Error("Programn finished with error", "error", err)
		os.Exit(1)
	}
}
