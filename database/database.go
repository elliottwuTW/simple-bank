package database

import (
	"context"

	"github.com/simple_bank/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Database struct {
	db         *mongo.Database
	accountDAO AccountDAO
}

func New(ctx context.Context, cfg config.DBConfig) (*Database, error) {
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
	return &Database{
		db:         db,
		accountDAO: newAccountDAO(db),
	}, nil
}

func (d *Database) execTx(
	ctx context.Context,
	fn func(
		ctx context.Context,
		d *Database,
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
