package model

import "time"

type VerifyEmail struct {
	ID       string `bson:"id"       json:"id"`
	Username string `bson:"username" json:"username"`
	Email    string `bson:"email"    json:"email"`
	Secret   string `bson:"secret"   json:"secret"`
	// to tell the secret has been used or not for security reasons
	IsUsed    bool      `bson:"isUsed"   json:"isUsed"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	ExpiredAt time.Time `bson:"expiredAt" json:"expiredAt"`
}
