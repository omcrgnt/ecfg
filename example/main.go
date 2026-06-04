package main

import (
	"fmt"
	"log"

	"github.com/omcrgnt/ecfg"
	"github.com/omcrgnt/ecfg/internal/testdata"
)

//go:generate go run github.com/omcrgnt/ecfg/cmd/ecfg-gen -type AppConfig -pkg github.com/omcrgnt/ecfg/internal/testdata -prefix APP -o env.template

func main() {
	cfg, err := ecfg.Parse[testdata.AppConfig](ecfg.WithPrefix("APP"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg.Server.Label)
}
