package main

import (
	"github.com/arendi-project/ba-version-2/config"
	"github.com/arendi-project/ba-version-2/internal/app"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	app.Run(cfg)
}
