package model

import "time"

type Transfer struct {
	ID            int64 `bson:"id"        json:"id"`
	FromAccountID int64 `bson:"fromAccountID" json:"fromAccountID"`
	ToAccountID   int64 `bson:"toAccountID" json:"toAccountID"`
	// must be positive
	Amount    int64     `bson:"amount"    json:"amount"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}
