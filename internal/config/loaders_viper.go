package config

import (
	"context"
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

// Load lê os valores do Viper, resolve o arquivo de token e valida tudo.
func Load(ctx context.Context, v *viper.Viper) (*Config, error) {
	if ctx.Err() != nil {
		slog.Default().Error("Context Finished", "error", ctx.Err())
		return nil, ctx.Err()
	}
	tokenPath := strings.TrimSpace(v.GetString(KeyTelegramToken))
	token, err := readTokenFile(tokenPath)
	if err != nil {
		return nil, err
	}

	// GetStringSlice é o único getter que atende flag e variável de ambiente:
	// o Viper não tem caso para int64Slice, e o cast só separa string em slice
	// para []string. A conversão para int64 fica em parseAllowedUsers.
	allowedUsers, err := parseAllowedUsers(v.GetStringSlice(KeyAllowedUsers))
	if err != nil {
		return nil, err
	}

	level, err := parseLogLevel(v.GetString(KeyLogLevel))
	if err != nil {
		return nil, err
	}

	rpcURL, err := parseTransmissionRPCURL(v.GetString(KeyTransmissionRPCURL))
	if err != nil {
		return nil, err
	}

	return &Config{
		TelegramBotToken:   token,
		TokenFilePath:      tokenPath,
		AllowedUsersID:     allowedUsers,
		LogLevel:           level,
		TransmissionRPCURL: rpcURL,
	}, nil
}
