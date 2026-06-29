package main_test

import (
	"testing"

	"github.com/omcrgnt/ecfg/config"
	"github.com/omcrgnt/ecfg/internal/testdata"
)

func TestExampleParse(t *testing.T) {
	t.Setenv("APP_SERVER_LABEL", "demo")
	cfg, err := config.Parse[testdata.AppConfig](config.WithPrefix("APP"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Label != "demo" {
		t.Fatalf("label: %q", cfg.Server.Label)
	}
}
