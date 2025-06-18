package database

import "go.mongodb.org/mongo-driver/mongo"

type dao struct {
	coll *mongo.Collection
}

func newDAO(db *mongo.Database, collName string) dao {
	return dao{
		coll: db.Collection(collName),
	}
}
