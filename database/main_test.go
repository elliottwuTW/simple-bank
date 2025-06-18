package database

import (
	"context"
	"log"
	"os"
	"testing"
)

var testDB *Database

func TestMain(m *testing.M) {
	db, err := New(context.Background())
	if err != nil {
		log.Fatal("database initialization fail", err)
	}
	testDB = db

	os.Exit(m.Run())
}
