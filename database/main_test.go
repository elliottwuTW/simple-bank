package database

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/simple_bank/config"
)

var testDB *Database

func TestMain(m *testing.M) {
	cfg, err := config.LoadConfig("..")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	db, err := New(context.Background(), cfg.DB)
	if err != nil {
		log.Fatal("database initialization fail", err)
	}
	testDB = db

	os.Exit(m.Run())
}
