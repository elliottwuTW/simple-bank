package database

import (
	"context"

	"github.com/simple_bank/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Database interface {
	CreateAccount(
		ctx context.Context, params CreateAccountParams,
	) (model.Account, error)
	GetAccount(
		ctx context.Context, id primitive.ObjectID,
	) (model.Account, error)
	UpdateAccount(
		ctx context.Context, params UpdateAccountParams,
	) (model.Account, error)
	DeleteAccount(
		ctx context.Context, id primitive.ObjectID,
	) error
	CreateUser(
		ctx context.Context, params CreateUserParams,
	) (model.User, error)
	CreateUserTx(
		ctx context.Context, arg CreateUserTxParams,
	) (CreateUserTxResult, error)
}
