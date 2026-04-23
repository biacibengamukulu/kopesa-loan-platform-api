package main

import (
	"log"

	"github.com/biangacila/kopesa-loan-platform-api/internal/app"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/config"
)

func main() {
	cfg := config.Load()
	server, cleanup, err := app.Build(cfg)
	if err != nil {
		log.Fatalf("build app: %v", err)
	}
	defer cleanup()

	if err := server.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("listen: %v", err)
	}
}
