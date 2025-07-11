package database

import (
	"context"

	"github.com/simple_bank/config"
	"github.com/simple_bank/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	db             *mongo.Database
	accountDAO     AccountDAO
	userDAO        UserDAO
	verifyEmailDAO VerifyEmailDAO
}

func New(ctx context.Context, cfg config.DBConfig) (Database, error) {
	opt := options.Client().ApplyURI(cfg.URI)
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		client.Disconnect(ctx)
		return nil, err
	}

	db := client.Database(cfg.Name)
	return &MongoDB{
		db:             db,
		accountDAO:     newAccountDAO(db),
		userDAO:        newUserDAO(db),
		verifyEmailDAO: newVerifyEmailDAO(db),
	}, nil
}

func (d *MongoDB) execTx(
	ctx context.Context,
	fn func(
		ctx context.Context,
		d Database,
	) error,
) error {
	rc := readconcern.Majority()
	txnOpts := options.Transaction().SetReadConcern(rc)
	sessOpts := options.Session().SetDefaultReadConcern(rc)

	// https://www.mongodb.com/docs/drivers/go/current/fundamentals/transactions/
	sess, err := d.db.Client().StartSession(sessOpts)
	if err != nil {
		return err
	}

	defer sess.EndSession(ctx)
	_, err = sess.WithTransaction(
		ctx,
		func(sessCtx mongo.SessionContext) (interface{}, error) {
			return nil, fn(sessCtx, d)
		},
		txnOpts,
	)
	return err
}

func (d *MongoDB) CreateUser(
	ctx context.Context, params CreateUserParams,
) (model.User, error) {
	return d.userDAO.CreateUser(ctx, params)
}

func (d *MongoDB) CreateAccount(
	ctx context.Context, params CreateAccountParams,
) (model.Account, error) {
	return d.accountDAO.CreateAccount(ctx, params)
}

func (d *MongoDB) GetAccount(
	ctx context.Context, id primitive.ObjectID,
) (model.Account, error) {
	return d.accountDAO.GetAccount(ctx, id)
}

func (d *MongoDB) UpdateAccount(
	ctx context.Context, params UpdateAccountParams,
) (model.Account, error) {
	return d.accountDAO.UpdateAccount(ctx, params)
}

func (d *MongoDB) DeleteAccount(
	ctx context.Context, id primitive.ObjectID,
) error {
	return d.accountDAO.DeleteAccount(ctx, id)
}

func (d *MongoDB) CreateVerifyEmail(
	ctx context.Context, params CreateVerifyEmailParams,
) (model.VerifyEmail, error) {
	return d.verifyEmailDAO.CreateVerifyEmail(ctx, params)
}
