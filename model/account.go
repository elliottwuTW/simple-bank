package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	ID        primitive.ObjectID `bson:"_id"        json:"_id"`
	Owner     string             `bson:"owner"     json:"owner"`
	Balance   int64              `bson:"balance"   json:"balance"`
	Currency  string             `bson:"currency"  json:"currency"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
