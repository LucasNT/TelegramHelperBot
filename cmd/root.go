// Package cmd define a interface de linha de comando do TorrentHelperBot.
package cmd

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/LucasNT/TorrentHelperBot/internal/config"
	"github.com/LucasNT/TorrentHelperBot/internal/logging"
	telegrambot "github.com/LucasNT/TorrentHelperBot/internal/telegram-bot"
)

// newRootCmd monta o comando raiz com sua própria instância do Viper, evitando
// o singleton global.
func newRootCmd() *cobra.Command {
	v := viper.New()

	cmd := &cobra.Command{
		Use:   "torrenthelperbot",
		Short: "Telegram bot for managing torrents",
		Long: "TorrentHelperBot exposes torrent management through a Telegram bot,\n" +
			"restricted to a list of authorized users.\n\n" +
			"Every flag can also be set through an environment variable with the\n" +
			"prefix " + config.EnvPrefix + "_ (for example, " + config.EnvPrefix + "_LOG_LEVEL).",
		// Um erro de execução não deve despejar o help inteiro; o Cobra ainda
		// imprime a mensagem de erro.
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			v.SetEnvPrefix(config.EnvPrefix)
			// Sem o replacer, telegram-bot-token não mapeia para
			// TORRENT_HELPER_BOT_TELEGRAM_BOT_TOKEN.
			v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			v.AutomaticEnv()
			// Depois do parse das flags, para que flag tenha precedência sobre env.
			return v.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			cfg, err := config.Load(ctx, v)
			if err != nil {
				return err
			}
			config.DefaultConfig = cfg

			logger := logging.New(cfg.LogLevel)
			slog.SetDefault(logger)

			slog.DebugContext(ctx, "Creating the bot")
			tbt, err := telegrambot.Initialization(ctx)
			if err != nil {
				slog.Default().Error("Critical Error, failed to create bot", "error", err)
				os.Exit(1)
			}

			slog.Default().Info("Starting bot")
			tbt.Start(ctx)

			return nil
		},
	}

	flags := cmd.Flags()
	flags.String(config.KeyTelegramToken, "", "path to the file holding the Telegram bot token")
	flags.StringSlice(config.KeyAllowedUsers, nil,
		"numeric IDs of the authorized users (comma-separated or by repeating the flag; "+
			"space-separated when set through the environment variable)")
	flags.String(config.KeyLogLevel, "info", "log level: debug, info, warn or error")
	flags.String(config.KeyTransmissionRPCURL, config.DefaultTransmissionRPCURL,
		"full Transmission RPC URL: scheme, host, port and path "+
			"(for example http://localhost:9091/transmission/rpc)")

	return cmd
}

// Execute roda o comando raiz e encerra o processo em caso de erro.
func Execute(ctx context.Context) error {
	return newRootCmd().ExecuteContext(ctx)
}
