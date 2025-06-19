package main

import (
	"context"
	"log"

	"github.com/simple_bank/api"
	"github.com/simple_bank/database"
)

const address = "0.0.0.0:8080"

func main() {
	db, err := database.New(context.Background())
	if err != nil {
		log.Fatal("cannot initialize database")
	}

	server := api.NewServer(db)
	err = server.Start(address)
	if err != nil {
		log.Fatal("cannot initialize server")
	}
}
