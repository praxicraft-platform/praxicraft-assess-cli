package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/config"
)

func TestResolvePrecedence(t *testing.T) {
	dir := t.TempDir()
	cfg := filepath.Join(dir, "config.toml")
	t.Setenv("PRAXICRAFT_CONFIG", cfg)
	t.Setenv("PRAXICRAFT_API_KEY", "")
	t.Setenv("PRAXICRAFT_API_BASE_URL", "")
	t.Setenv("PRAXICRAFT_PROFILE", "")

	f := &config.File{
		DefaultProfile: "default",
		Profiles: map[string]config.Profile{
			"default": {APIKey: "from_file", BaseURL: "https://file.example"},
		},
	}
	if err := config.Save(f); err != nil {
		t.Fatal(err)
	}

	res, err := config.Resolve("", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if res.APIKey != "from_file" {
		t.Fatalf("got %q", res.APIKey)
	}

	t.Setenv("PRAXICRAFT_API_KEY", "from_env")
	res, err = config.Resolve("", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if res.APIKey != "from_env" {
		t.Fatalf("got %q", res.APIKey)
	}

	res, err = config.Resolve("", "from_flag", "https://flag.example")
	if err != nil {
		t.Fatal(err)
	}
	if res.APIKey != "from_flag" || res.BaseURL != "https://flag.example" {
		t.Fatalf("%+v", res)
	}

	_ = os.Remove(cfg)
}
