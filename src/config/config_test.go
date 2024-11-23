package config_test

import (
	"testing"

	"log/slog"

	"github.com/vdbulcke/confkv/src/config"
)

func TestConfig(t *testing.T) {

	filename := "../../example/config.toml"

	cfg := config.MustOpen(filename)

	slog.Info("parsed config", "config", cfg)

	if !cfg.Validate() {

		t.Fail()
	}

}
