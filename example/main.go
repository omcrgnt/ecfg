package main

import (
	"fmt"
	"log"

	"github.com/omcrgnt/ecfg/config"
	"github.com/omcrgnt/ecfg/internal/testdata"
)

//go:generate go run github.com/omcrgnt/ecfg/cmd/ecfg-gen -type AppConfig -pkg github.com/omcrgnt/ecfg/internal/testdata -prefix APP -template env.template

func main() {
	cfg, err := config.Parse[testdata.AppConfig](config.WithPrefix("APP"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg.Server.Label)
}
