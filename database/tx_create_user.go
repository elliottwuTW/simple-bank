package database

import (
	"context"

	"github.com/simple_bank/model"
)

// CreateUserTxParams contains the input parameters of the create user transaction
type CreateUserTxParams struct {
	CreateUserParams
	// 等於傳了 callback function 進來，讓 create user 完直接呼叫！！！！！
	AfterCreate func(user model.User) error
}

// CreateUserTxResult is the result of the create user transaction
type CreateUserTxResult struct {
	User model.User
}

// CreateUserTx 可以在 create user 完，又執行其它的操作，當任一操作有錯誤時，仍可以 rollback
// 而不會在 DB 裡面產生出 user 的資料!
func (d *MongoDB) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := d.execTx(ctx, func(ctx context.Context, d Database) error {
		var err error
		result.User, err = d.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)
	})

	return result, err
}
