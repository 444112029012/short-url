package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Config is process configuration (ENG-002).
// Rate-limit fields are loaded for ENG-011; this batch does not enforce them.
type Config struct {
	AppEnv                  string
	HTTPAddr                string
	DatabasePath            string
	BaseURL                 string
	RateLimitCreatePerMin   int
	RateLimitRedirectPerMin int
	TrustedProxyCIDRs       string
	LogLevel                string
	// HSTSEnabled opts in to Strict-Transport-Security. The HTTP layer still
	// sends the header only when the request is TLS (SEC-008). Default false
	// so local HTTP does not advertise HSTS.
	HSTSEnabled bool
}

// Load reads the process environment. Empty or unset values fall back to local defaults.
// The server does not open .env itself; see .env.example.
func Load() (Config, error) {
	cfg := Config{
		AppEnv:                  getenv("APP_ENV", "local"),
		HTTPAddr:                getenv("HTTP_ADDR", "127.0.0.1:8080"),
		DatabasePath:            getenv("DATABASE_PATH", "./data/shorturl.local.db"),
		BaseURL:                 strings.TrimRight(getenv("BASE_URL", "http://localhost:8080"), "/"),
		RateLimitCreatePerMin:   30,
		RateLimitRedirectPerMin: 120,
		TrustedProxyCIDRs:       strings.TrimSpace(os.Getenv("TRUSTED_PROXY_CIDRS")),
		LogLevel:                getenv("LOG_LEVEL", "info"),
	}
	if strings.TrimSpace(cfg.DatabasePath) == "" {
		return Config{}, fmt.Errorf("DATABASE_PATH is empty")
	}
	if strings.TrimSpace(cfg.HTTPAddr) == "" {
		return Config{}, fmt.Errorf("HTTP_ADDR is empty")
	}
	if err := validateBaseURL(cfg.BaseURL); err != nil {
		return Config{}, err
	}
	var err error
	cfg.RateLimitCreatePerMin, err = parsePositive("RATE_LIMIT_CREATE_PER_MIN", cfg.RateLimitCreatePerMin)
	if err != nil {
		return Config{}, err
	}
	cfg.RateLimitRedirectPerMin, err = parsePositive("RATE_LIMIT_REDIRECT_PER_MIN", cfg.RateLimitRedirectPerMin)
	if err != nil {
		return Config{}, err
	}
	if err := validateCIDRs(cfg.TrustedProxyCIDRs); err != nil {
		return Config{}, err
	}
	cfg.HSTSEnabled, err = parseBool("HSTS_ENABLED", false)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func getenv(key, fallback string) string {
	v, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(v) == "" {
		return fallback
	}
	return strings.TrimSpace(v)
}

func parseBool(key string, fallback bool) (bool, error) {
	raw, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(raw) == "" {
		return fallback, nil
	}
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "1", "true", "yes", "on":
		return true, nil
	case "0", "false", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("%s must be true or false", key)
	}
}

func parsePositive(key string, fallback int) (int, error) {
	raw, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(raw) == "" {
		return fallback, nil
	}
	n, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", key)
	}
	return n, nil
}

func validateBaseURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") || !u.IsAbs() {
		return fmt.Errorf("BASE_URL must be an absolute http or https URL")
	}
	if u.User != nil {
		return fmt.Errorf("BASE_URL must not include user info")
	}
	if u.Path != "" && u.Path != "/" {
		return fmt.Errorf("BASE_URL must not include a path prefix")
	}
	if u.RawQuery != "" || u.Fragment != "" {
		return fmt.Errorf("BASE_URL must not include a query or fragment")
	}
	return nil
}

func validateCIDRs(raw string) error {
	if raw == "" {
		return nil
	}
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if _, _, err := net.ParseCIDR(part); err != nil {
			return fmt.Errorf("TRUSTED_PROXY_CIDRS contains an invalid CIDR")
		}
	}
	return nil
}
