package database

import (
	"context"
	"time"

	"github.com/simple_bank/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AccountDAO struct {
	dao
}

func newAccountDAO(db *mongo.Database) AccountDAO {
	return AccountDAO{
		dao: newDAO(db, "accounts"),
	}
}

type CreateAccountParams struct {
	Owner    string `bson:"owner"     json:"owner"`
	Balance  int64  `bson:"balance"   json:"balance"`
	Currency string `bson:"currency"  json:"currency"`
}

func (dao *AccountDAO) CreateAccount(
	ctx context.Context, params CreateAccountParams,
) (model.Account, error) {
	account := model.Account{
		ID:        primitive.NewObjectID(),
		Owner:     params.Owner,
		Balance:   params.Balance,
		Currency:  params.Currency,
		CreatedAt: time.Now(),
	}
	_, err := dao.coll.InsertOne(ctx, account)
	if err != nil {
		return model.Account{}, err
	}
	return account, nil
}

func (dao *AccountDAO) GetAccount(
	ctx context.Context, id primitive.ObjectID,
) (model.Account, error) {
	filter := bson.M{"_id": id}

	account := model.Account{}
	err := dao.coll.FindOne(ctx, filter).Decode(&account)
	return account, err
}

func (dao *AccountDAO) ListAccounts(
	ctx context.Context, limit int64, offset int64,
) ([]model.Account, error) {
	opt := options.Find()
	opt.SetLimit(limit)
	opt.SetSkip(offset)

	cur, err := dao.coll.Find(ctx, bson.M{}, opt)
	if err != nil {
		return nil, err
	}

	accounts := []model.Account{}
	err = cur.All(ctx, &accounts)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

type UpdateAccountParams struct {
	ID      primitive.ObjectID `bson:"_id"     json:"_id"`
	Balance int64              `bson:"balance"   json:"balance"`
}

func (dao *AccountDAO) UpdateAccount(
	ctx context.Context, params UpdateAccountParams,
) (model.Account, error) {
	filter := bson.M{"_id": params.ID}
	update := bson.M{
		"$set": bson.M{
			"balance": params.Balance,
		},
	}

	account := model.Account{}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := dao.coll.FindOneAndUpdate(ctx, filter, update, opt).Decode(&account)
	return account, err
}

func (dao *AccountDAO) DeleteAccount(
	ctx context.Context, id primitive.ObjectID,
) error {
	filter := bson.M{"_id": id}

	return dao.coll.FindOneAndDelete(ctx, filter).Err()
}
