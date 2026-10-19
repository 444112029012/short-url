package config

import (
	"os"
	"strings"
	"testing"
)

func clearConfigEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"APP_ENV",
		"HTTP_ADDR",
		"DATABASE_PATH",
		"BASE_URL",
		"RATE_LIMIT_CREATE_PER_MIN",
		"RATE_LIMIT_REDIRECT_PER_MIN",
		"TRUSTED_PROXY_CIDRS",
		"LOG_LEVEL",
		"HSTS_ENABLED",
	} {
		t.Setenv(key, "")
	}
}

func TestLoadDefaults(t *testing.T) {
	clearConfigEnv(t)
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AppEnv != "local" || cfg.HTTPAddr != "127.0.0.1:8080" {
		t.Fatalf("addr/env: %+v", cfg)
	}
	if cfg.DatabasePath != "./data/shorturl.local.db" {
		t.Fatalf("db: %s", cfg.DatabasePath)
	}
	if cfg.BaseURL != "http://localhost:8080" {
		t.Fatalf("base: %s", cfg.BaseURL)
	}
	if cfg.RateLimitCreatePerMin != 30 || cfg.RateLimitRedirectPerMin != 120 {
		t.Fatalf("limits: %+v", cfg)
	}
	if cfg.HSTSEnabled {
		t.Fatal("HSTS must default off")
	}
}

func TestLoadOverrides(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("BASE_URL", "https://short.example/")
	t.Setenv("DATABASE_PATH", "/tmp/shorturl-test.db")
	t.Setenv("HTTP_ADDR", "127.0.0.1:9090")
	t.Setenv("RATE_LIMIT_CREATE_PER_MIN", "15")
	t.Setenv("RATE_LIMIT_REDIRECT_PER_MIN", "40")
	t.Setenv("HSTS_ENABLED", "true")
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.BaseURL != "https://short.example" || cfg.DatabasePath != "/tmp/shorturl-test.db" {
		t.Fatalf("%+v", cfg)
	}
	if cfg.HTTPAddr != "127.0.0.1:9090" || cfg.RateLimitCreatePerMin != 15 || cfg.RateLimitRedirectPerMin != 40 {
		t.Fatalf("%+v", cfg)
	}
	if !cfg.HSTSEnabled {
		t.Fatal("HSTS_ENABLED=true was not loaded")
	}
}

func TestLoadRejectsBadValues(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("BASE_URL", "ftp://example.com")
	if _, err := Load(); err == nil {
		t.Fatal("expected BASE_URL error")
	}
	t.Setenv("BASE_URL", "http://localhost:8080")
	t.Setenv("RATE_LIMIT_CREATE_PER_MIN", "0")
	if _, err := Load(); err == nil {
		t.Fatal("expected rate limit error")
	}
	t.Setenv("RATE_LIMIT_CREATE_PER_MIN", "30")
	t.Setenv("TRUSTED_PROXY_CIDRS", "not-a-cidr")
	if _, err := Load(); err == nil {
		t.Fatal("expected CIDR error")
	}
	t.Setenv("TRUSTED_PROXY_CIDRS", "")
	t.Setenv("HSTS_ENABLED", "maybe")
	if _, err := Load(); err == nil {
		t.Fatal("expected HSTS_ENABLED error")
	}
}

func TestEnvExampleDocumentsKeys(t *testing.T) {
	raw, err := os.ReadFile("../../.env.example")
	if err != nil {
		t.Fatal(err)
	}
	text := string(raw)
	for _, key := range []string{"BASE_URL", "DATABASE_PATH", "HTTP_ADDR", "RATE_LIMIT_CREATE_PER_MIN", "RATE_LIMIT_REDIRECT_PER_MIN", "HSTS_ENABLED"} {
		if !strings.Contains(text, key+"=") {
			t.Fatalf(".env.example missing %s", key)
		}
	}
}
