package config

import (
	"os"
	"testing"
)

func TestNewDefaults(t *testing.T) {
	cfg := New("", "", 0, "")
	if cfg.Profile != "default" {
		t.Errorf("expected default profile, got %s", cfg.Profile)
	}
	if cfg.Timeout != DefaultTimeout {
		t.Errorf("expected timeout %d, got %d", DefaultTimeout, cfg.Timeout)
	}
}

func TestNewFlagsOverride(t *testing.T) {
	cfg := New("https://gateway.example.com", "key-from-flag", 30, "prod")
	if cfg.BaseURL != "https://gateway.example.com" {
		t.Errorf("expected flag base URL, got %s", cfg.BaseURL)
	}
	if cfg.APIKey != "key-from-flag" {
		t.Errorf("expected flag API key, got %s", cfg.APIKey)
	}
	if cfg.Timeout != 30 {
		t.Errorf("expected timeout 30, got %d", cfg.Timeout)
	}
	if cfg.Profile != "prod" {
		t.Errorf("expected profile prod, got %s", cfg.Profile)
	}
}

func TestNewEnvOverride(t *testing.T) {
	os.Setenv("CARBONSTOP_API_KEY", "key-from-env")
	os.Setenv("CARBONSTOP_BASE_URL", "https://env.example.com")
	os.Setenv("CARBONSTOP_TIMEOUT", "45")
	defer func() {
		os.Unsetenv("CARBONSTOP_API_KEY")
		os.Unsetenv("CARBONSTOP_BASE_URL")
		os.Unsetenv("CARBONSTOP_TIMEOUT")
	}()

	cfg := New("", "", 0, "")
	if cfg.APIKey != "key-from-env" {
		t.Errorf("expected env API key, got %s", cfg.APIKey)
	}
	if cfg.BaseURL != "https://env.example.com" {
		t.Errorf("expected env base URL, got %s", cfg.BaseURL)
	}
	if cfg.Timeout != 45 {
		t.Errorf("expected timeout 45, got %d", cfg.Timeout)
	}
}

func TestNewFlagWinsOverEnv(t *testing.T) {
	os.Setenv("CARBONSTOP_API_KEY", "key-from-env")
	defer os.Unsetenv("CARBONSTOP_API_KEY")

	cfg := New("", "key-from-flag", 0, "")
	if cfg.APIKey != "key-from-flag" {
		t.Errorf("expected flag key to win, got %s", cfg.APIKey)
	}
}

func TestBaseURLTrimTrailingSlash(t *testing.T) {
	cfg := New("https://example.com/", "", 0, "")
	if cfg.BaseURL != "https://example.com" {
		t.Errorf("expected trailing slash trimmed, got %s", cfg.BaseURL)
	}
}

func TestFirstNonEmpty(t *testing.T) {
	if got := firstNonEmpty("", "b", "c"); got != "b" {
		t.Errorf("expected b, got %s", got)
	}
	if got := firstNonEmpty("a", "b"); got != "a" {
		t.Errorf("expected a, got %s", got)
	}
	if got := firstNonEmpty("", "", ""); got != "" {
		t.Errorf("expected empty, got %s", got)
	}
}
