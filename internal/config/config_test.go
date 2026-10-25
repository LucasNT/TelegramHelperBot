package config

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

// writeTokenFile cria um arquivo de token temporário e devolve seu caminho.
func writeTokenFile(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "token")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("escrevendo o arquivo de token: %v", err)
	}

	return path
}

func TestReadTokenFile(t *testing.T) {
	validPath := writeTokenFile(t, "123456:FAKE-TOKEN\n")
	emptyPath := writeTokenFile(t, "   \n\t ")
	missingPath := filepath.Join(t.TempDir(), "inexistente")

	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{name: "arquivo válido com newline final", path: validPath, want: "123456:FAKE-TOKEN"},
		{name: "caminho não informado", path: "", wantErr: true},
		{name: "arquivo inexistente", path: missingPath, wantErr: true},
		{name: "arquivo só com whitespace", path: emptyPath, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readTokenFile(tt.path)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("esperava erro, recebeu token %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if got != tt.want {
				t.Errorf("token = %q, esperado %q", got, tt.want)
			}
		})
	}
}

// O conteúdo do token não pode vazar em mensagem de erro.
func TestReadTokenFileErrorOmitsSecret(t *testing.T) {
	const secret = "123456:SUPER-SECRET"

	path := writeTokenFile(t, secret)
	if err := os.Chmod(path, 0o000); err != nil {
		t.Skipf("não foi possível remover a permissão de leitura: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(path, 0o600) })

	_, err := readTokenFile(path)
	if err == nil {
		t.Skip("arquivo ainda legível (provavelmente rodando como root)")
	}
	if strings.Contains(err.Error(), secret) {
		t.Errorf("mensagem de erro vazou o token: %v", err)
	}
}

func TestParseAllowedUsers(t *testing.T) {
	tests := []struct {
		name    string
		ids     []string
		want    []int64
		wantErr bool
	}{
		{name: "lista válida", ids: []string{"1", "2", "3"}, want: []int64{1, 2, 3}},
		{name: "ID negativo é válido", ids: []string{"-1001234"}, want: []int64{-1001234}},
		{name: "preserva a ordem de entrada", ids: []string{"3", "1", "2"}, want: []int64{3, 1, 2}},
		{name: "ID maior que int32", ids: []string{"7000000000"}, want: []int64{7000000000}},
		{name: "lista vazia", ids: nil, wantErr: true},
		{name: "ID não numérico", ids: []string{"1", "abc"}, wantErr: true},
		{name: "ID duplicado", ids: []string{"111", "222", "111"}, wantErr: true},
		{name: "vírgula não separada pelo caller", ids: []string{"1,2,3"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAllowedUsers(tt.ids)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("esperava erro, recebeu %v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("IDs = %v, esperado %v", got, tt.want)
			}
		})
	}
}

// O erro deve nomear qual ID é inválido.
func TestParseAllowedUsersErrorNamesID(t *testing.T) {
	_, err := parseAllowedUsers([]string{"1", "abc", "3"})
	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}
	if !strings.Contains(err.Error(), "abc") {
		t.Errorf("erro não cita o ID inválido: %v", err)
	}
}

// O erro deve nomear qual ID está repetido.
func TestParseAllowedUsersErrorNamesDuplicate(t *testing.T) {
	_, err := parseAllowedUsers([]string{"111", "222", "111"})
	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}
	if !strings.Contains(err.Error(), "111") {
		t.Errorf("erro não cita o ID duplicado: %v", err)
	}
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    slog.Level
		wantErr bool
	}{
		{name: "debug", input: "debug", want: slog.LevelDebug},
		{name: "maiúsculas", input: "INFO", want: slog.LevelInfo},
		{name: "warn", input: "warn", want: slog.LevelWarn},
		{name: "error", input: "error", want: slog.LevelError},
		{name: "vazio usa info", input: "", want: slog.LevelInfo},
		{name: "nível desconhecido", input: "verbose", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLogLevel(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("esperava erro, recebeu nível %v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if got != tt.want {
				t.Errorf("nível = %v, esperado %v", got, tt.want)
			}
		})
	}
}

func TestParseTransmissionRPCURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "URL completa", input: "http://localhost:9091/transmission/rpc", want: "http://localhost:9091/transmission/rpc"},
		{name: "https sem porta", input: "https://nas.lan/transmission/rpc", want: "https://nas.lan/transmission/rpc"},
		{name: "host IPv6", input: "http://[::1]:9091/transmission/rpc", want: "http://[::1]:9091/transmission/rpc"},
		{name: "caminho customizado", input: "https://nas.lan:8443/custom/path", want: "https://nas.lan:8443/custom/path"},
		{name: "raiz como caminho", input: "http://nas.lan/", want: "http://nas.lan/"},
		{name: "vazio usa o padrão", input: "", want: DefaultTransmissionRPCURL},
		{name: "whitespace ao redor", input: "  http://nas.lan/rpc\n", want: "http://nas.lan/rpc"},
		{name: "porta não numérica", input: "http://localhost:abc/transmission/rpc", wantErr: true},
		{name: "sem scheme", input: "//localhost:9091/transmission/rpc", wantErr: true},
		{name: "caminho relativo", input: "transmission/rpc", wantErr: true},
		{name: "scheme desconhecido", input: "ftp://host/transmission/rpc", wantErr: true},
		{name: "host:porta lido como scheme", input: "localhost:9091/transmission/rpc", wantErr: true},
		{name: "sem host", input: "http:///transmission/rpc", wantErr: true},
		{name: "scheme opaco", input: "http:transmission/rpc", wantErr: true},
		{name: "sem caminho", input: "http://localhost:9091", wantErr: true},
		{name: "com query string", input: "http://localhost:9091/transmission/rpc?x=1", wantErr: true},
		{name: "com fragment", input: "http://localhost:9091/transmission/rpc#frag", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTransmissionRPCURL(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("esperava erro, recebeu URL %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if got.String() != tt.want {
				t.Errorf("URL = %q, esperado %q", got, tt.want)
			}
		})
	}
}

// O erro deve nomear a flag, para o operador saber o que corrigir.
func TestParseTransmissionRPCURLErrorNamesFlag(t *testing.T) {
	_, err := parseTransmissionRPCURL("ftp://host/rpc")
	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}
	if !strings.Contains(err.Error(), KeyTransmissionRPCURL) {
		t.Errorf("erro não cita a flag: %v", err)
	}
}

func TestLoad(t *testing.T) {
	tokenPath := writeTokenFile(t, "123456:FAKE-TOKEN\n")

	v := viper.New()
	v.Set(KeyTelegramToken, tokenPath)
	v.Set(KeyAllowedUsers, []string{"111", "222"})
	v.Set(KeyLogLevel, "debug")
	v.Set(KeyTransmissionRPCURL, "http://nas.lan:8080/transmission/rpc")

	cfg, err := Load(context.Background(), v)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if cfg.TelegramBotToken != "123456:FAKE-TOKEN" {
		t.Errorf("token = %q", cfg.TelegramBotToken)
	}
	if cfg.TokenFilePath != tokenPath {
		t.Errorf("TokenFilePath = %q, esperado %q", cfg.TokenFilePath, tokenPath)
	}
	if !slices.Equal(cfg.AllowedUsersID, []int64{111, 222}) {
		t.Errorf("AllowedUsersID = %v", cfg.AllowedUsersID)
	}
	if cfg.LogLevel != slog.LevelDebug {
		t.Errorf("LogLevel = %v, esperado debug", cfg.LogLevel)
	}
	if got := cfg.TransmissionRPCURL.String(); got != "http://nas.lan:8080/transmission/rpc" {
		t.Errorf("TransmissionRPCURL = %q", got)
	}
}

// Sem --log-level, Load deve cair no padrão info.
func TestLoadDefaultsLogLevel(t *testing.T) {
	v := viper.New()
	v.Set(KeyTelegramToken, writeTokenFile(t, "token"))
	v.Set(KeyAllowedUsers, []string{"111"})

	cfg, err := Load(context.Background(), v)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if cfg.LogLevel != slog.LevelInfo {
		t.Errorf("LogLevel = %v, esperado info", cfg.LogLevel)
	}
}

// Sem --transmission-rpc-url, Load deve cair no endpoint local padrão.
func TestLoadDefaultsTransmissionRPCURL(t *testing.T) {
	v := viper.New()
	v.Set(KeyTelegramToken, writeTokenFile(t, "token"))
	v.Set(KeyAllowedUsers, []string{"111"})

	cfg, err := Load(context.Background(), v)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if got := cfg.TransmissionRPCURL.String(); got != DefaultTransmissionRPCURL {
		t.Errorf("TransmissionRPCURL = %q, esperado %q", got, DefaultTransmissionRPCURL)
	}
}

// A lista vinda de env var é separada por espaço pelo próprio Viper.
func TestLoadSplitsEnvListOnWhitespace(t *testing.T) {
	t.Setenv(EnvPrefix+"_ALLOWED_USERS_ID", "111 222 333")

	v := viper.New()
	v.SetEnvPrefix(EnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
	v.Set(KeyTelegramToken, writeTokenFile(t, "token"))

	cfg, err := Load(context.Background(), v)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if !slices.Equal(cfg.AllowedUsersID, []int64{111, 222, 333}) {
		t.Errorf("AllowedUsersID = %v", cfg.AllowedUsersID)
	}
}

// Um valor explícito (flag) tem precedência sobre a variável de ambiente.
func TestLoadFlagOverridesEnv(t *testing.T) {
	t.Setenv(EnvPrefix+"_ALLOWED_USERS_ID", "999")
	t.Setenv(EnvPrefix+"_LOG_LEVEL", "error")

	v := viper.New()
	v.SetEnvPrefix(EnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
	v.Set(KeyTelegramToken, writeTokenFile(t, "token"))
	v.Set(KeyAllowedUsers, []string{"111"})
	v.Set(KeyLogLevel, "debug")

	cfg, err := Load(context.Background(), v)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if !slices.Equal(cfg.AllowedUsersID, []int64{111}) {
		t.Errorf("AllowedUsersID = %v, esperado [111]", cfg.AllowedUsersID)
	}
	if cfg.LogLevel != slog.LevelDebug {
		t.Errorf("LogLevel = %v, esperado debug", cfg.LogLevel)
	}
}
