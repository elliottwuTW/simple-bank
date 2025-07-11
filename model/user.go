package model

type User struct {
	Username        string `bson:"username"        json:"username"`
	HashedPassword  string `bson:"hashedPassword"  json:"hashedPassword"`
	Email           string `bson:"email"           json:"email"`
	IsEmailVerified bool   `bson:"isEmailVerified" json:"isEmailVerified"`
}
