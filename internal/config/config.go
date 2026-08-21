package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

const (
	EnvAPIKey  = "PRAXICRAFT_API_KEY"
	EnvBaseURL = "PRAXICRAFT_API_BASE_URL"
	EnvProfile = "PRAXICRAFT_PROFILE"
)

// File is the on-disk config (~/.config/praxicraft/config.toml).
type File struct {
	DefaultProfile string             `toml:"default_profile"`
	Profiles       map[string]Profile `toml:"profiles"`
}

// Profile holds credentials for one named profile.
type Profile struct {
	APIKey  string `toml:"api_key"`
	BaseURL string `toml:"base_url"`
}

// Resolved is the effective runtime config after flag/env/profile merge.
type Resolved struct {
	Profile string
	APIKey  string
	BaseURL string
}

// ConfigPath returns the default config file path.
func ConfigPath() (string, error) {
	if x := os.Getenv("PRAXICRAFT_CONFIG"); x != "" {
		return x, nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "praxicraft", "config.toml"), nil
}

// Load reads the config file; missing file returns empty config.
func Load() (*File, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}
	f := &File{Profiles: map[string]Profile{}}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return f, nil
		}
		return nil, err
	}
	if _, err := toml.Decode(string(data), f); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if f.Profiles == nil {
		f.Profiles = map[string]Profile{}
	}
	return f, nil
}

// Save writes the config file, creating directories as needed.
func Save(f *File) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	out, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer out.Close()
	enc := toml.NewEncoder(out)
	return enc.Encode(f)
}

// Resolve merges flags > env > profile > defaults.
func Resolve(flagProfile, flagAPIKey, flagBaseURL string) (*Resolved, error) {
	f, err := Load()
	if err != nil {
		return nil, err
	}
	profile := firstNonEmpty(flagProfile, os.Getenv(EnvProfile), f.DefaultProfile, "default")
	p := f.Profiles[profile]

	apiKey := firstNonEmpty(flagAPIKey, os.Getenv(EnvAPIKey), p.APIKey)
	baseURL := firstNonEmpty(flagBaseURL, os.Getenv(EnvBaseURL), p.BaseURL, "https://assess.praxicraft.com")

	return &Resolved{Profile: profile, APIKey: strings.TrimSpace(apiKey), BaseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/")}, nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
