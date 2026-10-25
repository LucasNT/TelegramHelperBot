// Package config resolve e valida a configuração da aplicação a partir de um
// *viper.Viper já populado com flags e variáveis de ambiente.
package config

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Prefixo das variáveis de ambiente e chaves de configuração. As chaves são
// idênticas aos nomes das flags; o replacer do Viper cuida da tradução para
// TORRENT_HELPER_BOT_TELEGRAM_BOT_TOKEN e afins.
const (
	EnvPrefix = "TORRENT_HELPER_BOT"

	KeyTelegramToken      = "telegram-bot-token"
	KeyAllowedUsers       = "allowed-users-id"
	KeyLogLevel           = "log-level"
	KeyTransmissionRPCURL = "transmission-rpc-url"

	// DefaultTransmissionRPCURL é o endpoint padrão de um Transmission local.
	DefaultTransmissionRPCURL = "http://localhost:9091/transmission/rpc"
)

// Config é a configuração já resolvida e validada.
type Config struct {
	TelegramBotToken   string  // conteúdo do arquivo, já trimado — nunca logar
	TokenFilePath      string  // caminho informado, seguro para log
	AllowedUsersID     []int64 // IDs numéricos dos usuários autorizados
	LogLevel           slog.Level
	TransmissionRPCURL *url.URL // endpoint RPC já validado
}

var (
	DefaultConfig *Config
)

// readTokenFile lê o token do arquivo apontado por path. O conteúdo nunca é
// incluído em mensagens de erro — apenas o caminho.
func readTokenFile(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("--%s is required: provide the path to the file holding the token", KeyTelegramToken)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading token file %q: %w", path, err)
	}

	token := strings.TrimSpace(string(content))
	if token == "" {
		return "", fmt.Errorf("token file %q is empty", path)
	}

	return token, nil
}

// parseAllowedUsers exige ao menos um ID, que todos sejam numéricos e que não
// haja repetição. IDs negativos são válidos: grupos e canais do Telegram usam
// IDs negativos.
func parseAllowedUsers(ids []string) ([]int64, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("--%s is required: provide at least one user ID", KeyAllowedUsers)
	}

	parsed := make([]int64, 0, len(ids))
	seen := make(map[int64]struct{}, len(ids))

	for _, id := range ids {
		userID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid user ID %q in --%s: must be numeric", id, KeyAllowedUsers)
		}

		if _, dup := seen[userID]; dup {
			return nil, fmt.Errorf("duplicate user ID %d in --%s", userID, KeyAllowedUsers)
		}
		seen[userID] = struct{}{}

		parsed = append(parsed, userID)
	}

	return parsed, nil
}

// parseLogLevel aceita debug, info, warn e error, sem diferenciar maiúsculas.
// String vazia resulta em slog.LevelInfo.
func parseLogLevel(s string) (slog.Level, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return slog.LevelInfo, nil
	}

	var level slog.Level
	if err := level.UnmarshalText([]byte(s)); err != nil {
		return 0, fmt.Errorf("invalid log level %q in --%s: use debug, info, warn or error", s, KeyLogLevel)
	}

	return level, nil
}

// parseTransmissionRPCURL exige uma URL http(s) absoluta, com host e caminho. A porta é
// opcional: sem ela vale a porta padrão do scheme. String vazia resulta em
// DefaultTransmissionRPCURL.
func parseTransmissionRPCURL(s string) (*url.URL, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		s = DefaultTransmissionRPCURL
	}

	parsed, err := url.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("invalid URL %q in --%s: %w", s, KeyTransmissionRPCURL, err)
	}

	// Sem scheme, url.Parse aceita caminho relativo; com um scheme desconhecido, aceita
	// qualquer coisa. "localhost:9091/rpc" cai aqui como scheme "localhost".
	switch parsed.Scheme {
	case "http", "https":
	case "":
		return nil, fmt.Errorf("invalid URL %q in --%s: missing scheme, use http:// or https://",
			s, KeyTransmissionRPCURL)
	default:
		return nil, fmt.Errorf("invalid scheme %q in --%s: use http or https",
			parsed.Scheme, KeyTransmissionRPCURL)
	}

	if parsed.Host == "" {
		return nil, fmt.Errorf("invalid URL %q in --%s: missing host", s, KeyTransmissionRPCURL)
	}

	// Um caminho só com "/" é aceito: um reverse proxy pode mapear a raiz para o endpoint.
	if parsed.Path == "" {
		return nil, fmt.Errorf("invalid URL %q in --%s: missing the RPC path (for example /transmission/rpc)",
			s, KeyTransmissionRPCURL)
	}

	if parsed.RawQuery != "" || parsed.Fragment != "" {
		return nil, fmt.Errorf("invalid URL %q in --%s: must not carry a query string or fragment",
			s, KeyTransmissionRPCURL)
	}

	return parsed, nil
}
