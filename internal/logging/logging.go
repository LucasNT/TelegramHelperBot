// Package logging constrói o logger estruturado da aplicação.
package logging

import (
	"log/slog"
	"os"
)

// New devolve um logger que escreve em stderr, deixando stdout livre para a
// saída de dados de subcomandos futuros.
func New(level slog.Level) *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
}
