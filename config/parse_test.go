package config_test

import (
	"errors"
	"testing"

	"github.com/omcrgnt/ecfg/config"
	"github.com/omcrgnt/ecfg/internal/testdata"
)

func TestParse_delegates(t *testing.T) {
	t.Setenv("APP_SERVER_LABEL", "ok")
	cfg, err := config.Parse[testdata.AppConfig](config.WithPrefix("APP"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Label != "ok" {
		t.Fatalf("got %q", cfg.Server.Label)
	}
}

func TestErrors_alias(t *testing.T) {
	t.Setenv("APP_SERVER_LABEL", "")
	_, err := config.Parse[testdata.AppConfig](config.WithPrefix("APP"))
	if !errors.Is(err, config.ErrEmptyEnv) {
		t.Fatalf("got %v", err)
	}
}
