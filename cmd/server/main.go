package main

import (
	"log"

	"github.com/tranquangkha19/auth-service/internal/config"
	"github.com/tranquangkha19/auth-service/internal/server"
)

func main() {
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	s := server.NewServer(cfg)
	s.Run()
}
