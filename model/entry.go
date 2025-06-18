package model

import "time"

type Entry struct {
	ID        int64 `bson:"id"        json:"id"`
	AccountID int64 `bson:"accountID" json:"accountID"`
	// can be positive or negative
	Amount    int64     `bson:"amount"    json:"amount"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}
