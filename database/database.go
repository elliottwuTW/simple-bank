package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const uri = "mongodb://localhost:27017"
const db = "simple-bank"

type Database struct {
	db         *mongo.Database
	accountDAO AccountDAO
}

func New(ctx context.Context) (*Database, error) {
	opt := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		client.Disconnect(ctx)
		return nil, err
	}

	db := client.Database(db)
	return &Database{
		db:         db,
		accountDAO: newAccountDAO(db),
	}, nil
}
