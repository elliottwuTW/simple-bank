package database

import (
	"context"

	"github.com/simple_bank/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDAO struct {
	dao
}

func newUserDAO(db *mongo.Database) UserDAO {
	return UserDAO{
		dao: newDAO(db, "users"),
	}
}

type CreateUserParams struct {
	Username       string `bson:"Username"       json:"Username"`
	HashedPassword string `bson:"hashedPassword" json:"hashedPassword"`
	Email          string `bson:"email"          json:"email"`
}

func (dao *UserDAO) CreateUser(
	ctx context.Context, params CreateUserParams,
) (model.User, error) {
	user := model.User{
		Username:       params.Username,
		HashedPassword: params.HashedPassword,
		Email:          params.Email,
	}
	_, err := dao.coll.InsertOne(ctx, user)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
