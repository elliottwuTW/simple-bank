package main

import (
	"context"
	"log"
	"time"

	"github.com/simple_bank/api"
	"github.com/simple_bank/config"
	"github.com/simple_bank/database"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := database.New(ctx, cfg.DB)
	if err != nil {
		log.Fatal("cannot initialize database", err)
	}

	server := api.NewServer(db)
	err = server.Start(cfg.ServerAddress)
	if err != nil {
		log.Fatal("cannot initialize server")
	}
}
