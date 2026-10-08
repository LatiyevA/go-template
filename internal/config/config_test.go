package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func configEnvironment(t *testing.T) {
	t.Helper()
	t.Chdir(t.TempDir())
	t.Setenv("SERVICE_NAME", "test")
	t.Setenv("PG_URL", "postgres://postgres:postgres@localhost/app_db")
	for _, key := range []string{
		"HTTP_PORT", "HTTP_SHUTDOWN_TIMEOUT", "PG_MAX_OPEN_CONNS", "PG_MAX_IDLE_CONNS",
		"PG_CONN_MAX_LIFETIME", "PG_CONN_MAX_IDLE_TIME", "PG_CONNECT_TIMEOUT", "MIGRATIONS_PATH",
	} {
		unset(t, key)
	}
}

func unset(t *testing.T, key string) {
	t.Helper()
	t.Setenv(key, "")
	if err := os.Unsetenv(key); err != nil {
		t.Fatal(err)
	}
}

func TestEnvironmentOnlyConfiguration(t *testing.T) {
	configEnvironment(t)
	cfg, err := Get()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.HTTPPort != 8080 || cfg.MaxOpenConns != 10 || cfg.ConnectTimeout != 5*time.Second {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
}

func TestEnvironmentOverridesDotenv(t *testing.T) {
	configEnvironment(t)
	t.Setenv("HTTP_PORT", "9090")
	if err := os.WriteFile(".env", []byte("SERVICE_NAME=dotenv\nHTTP_PORT=8081\n"), 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Get()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ServiceName != "test" || cfg.HTTPPort != 9090 {
		t.Fatalf("dotenv overwrote process environment: %+v", cfg)
	}
}

func TestDotenvSuppliesMissingValues(t *testing.T) {
	configEnvironment(t)
	unset(t, "SERVICE_NAME")
	if err := os.WriteFile(".env", []byte("SERVICE_NAME=local\nHTTP_PORT=8081\n"), 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Get()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ServiceName != "local" || cfg.HTTPPort != 8081 {
		t.Fatalf("dotenv values were not loaded: %+v", cfg)
	}
}

func TestMalformedDotenvIsReported(t *testing.T) {
	configEnvironment(t)
	if err := os.WriteFile(filepath.Join(".", ".env"), []byte("SERVICE_NAME=\"unterminated"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := Get(); err == nil {
		t.Fatal("expected dotenv parsing error")
	}
}

func TestConfigurationRejectsInvalidValues(t *testing.T) {
	for _, tc := range []struct{ key, value string }{
		{"SERVICE_NAME", ""}, {"PG_URL", ""}, {"HTTP_PORT", "0"},
		{"HTTP_PORT", "65536"}, {"HTTP_PORT", "invalid"},
		{"HTTP_SHUTDOWN_TIMEOUT", "0s"}, {"PG_MAX_OPEN_CONNS", "0"},
		{"PG_MAX_IDLE_CONNS", "-1"}, {"PG_MAX_IDLE_CONNS", "11"},
		{"PG_CONNECT_TIMEOUT", "-1s"}, {"PG_CONN_MAX_LIFETIME", "0s"},
		{"PG_CONN_MAX_IDLE_TIME", "0s"},
	} {
		t.Run(tc.key+"="+tc.value, func(t *testing.T) {
			configEnvironment(t)
			t.Setenv(tc.key, tc.value)
			if _, err := Get(); err == nil {
				t.Fatal("expected invalid configuration error")
			}
		})
	}
}

func TestConfigurationRequiresDatabaseURL(t *testing.T) {
	configEnvironment(t)
	unset(t, "PG_URL")
	if _, err := Get(); err == nil {
		t.Fatal("expected missing required variable error")
	}
}
