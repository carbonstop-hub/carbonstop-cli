// Package config manages CLI configuration from env vars, config files, and CLI flags.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/carbonstop/carbonstop-cli/internal/version"
)

const (
	DefaultTimeout = 60
	configDirName  = ".config/carbonstop-cli"
	configFileName = "config.json"
)

var (
	ConfigDir  = filepath.Join(os.Getenv("HOME"), configDirName)
	ConfigFile = filepath.Join(ConfigDir, configFileName)
)

// Config holds the resolved configuration.
type Config struct {
	BaseURL string
	APIKey  string
	Timeout int
	Profile string
}

// New creates a Config resolving from: CLI flags > env > config file (api key only) > built-in defaults.
// Base URL is never read from config file — only CLI flags, env var, or built-in.
func New(flagBaseURL, flagAPIKey string, flagTimeout int, flagProfile string) *Config {
	cfg := &Config{
		Profile: flagProfile,
	}
	if cfg.Profile == "" {
		cfg.Profile = os.Getenv("CARBONSTOP_PROFILE")
	}
	if cfg.Profile == "" {
		cfg.Profile = "default"
	}

	cfg.BaseURL = firstNonEmpty(
		flagBaseURL,
		os.Getenv("CARBONSTOP_BASE_URL"),
		version.GetBaseURL(),
	)
	cfg.BaseURL = strings.TrimRight(cfg.BaseURL, "/")

	cfg.APIKey = firstNonEmpty(
		flagAPIKey,
		os.Getenv("CARBONSTOP_API_KEY"),
		getProfileValue("api_key", cfg.Profile),
	)

	cfg.Timeout = DefaultTimeout
	if flagTimeout > 0 {
		cfg.Timeout = flagTimeout
	} else if t := os.Getenv("CARBONSTOP_TIMEOUT"); t != "" {
		fmt.Sscanf(t, "%d", &cfg.Timeout)
	} else if t := getProfileValue("timeout", cfg.Profile); t != "" {
		fmt.Sscanf(t, "%d", &cfg.Timeout)
	}

	return cfg
}

// SaveProfile saves a named profile to the config file.
// Only api_key is persisted; base_url is never written.
func (c *Config) SaveProfile(apiKey string) error {
	if err := os.MkdirAll(ConfigDir, 0700); err != nil {
		return err
	}

	data := loadConfigFile()
	if data == nil {
		data = make(map[string]interface{})
	}

	profiles, ok := data["profiles"].(map[string]interface{})
	if !ok {
		profiles = make(map[string]interface{})
		data["profiles"] = profiles
	}

	profile := make(map[string]string)
	if existing, ok := profiles[c.Profile].(map[string]interface{}); ok {
		for k, v := range existing {
			if s, ok := v.(string); ok {
				profile[k] = s
			}
		}
	}

	if apiKey != "" {
		profile["api_key"] = apiKey
	}

	profiles[c.Profile] = profile
	data["api_key"] = profile["api_key"]

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigFile, bytes, 0600)
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func getProfileValue(key, profile string) string {
	data := loadConfigFile()
	if data == nil {
		return ""
	}
	profiles, ok := data["profiles"].(map[string]interface{})
	if !ok {
		return ""
	}
	p, ok := profiles[profile].(map[string]interface{})
	if !ok {
		return ""
	}
	if v, ok := p[key].(string); ok {
		return v
	}
	return ""
}

func loadConfigFile() map[string]interface{} {
	bytes, err := os.ReadFile(ConfigFile)
	if err != nil {
		return nil
	}
	var data map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil
	}
	return data
}
