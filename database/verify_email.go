package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/simple_bank/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type VerifyEmailDAO struct {
	dao
}

func newVerifyEmailDAO(db *mongo.Database) VerifyEmailDAO {
	return VerifyEmailDAO{
		dao: newDAO(db, "verify_emails"),
	}
}

type CreateVerifyEmailParams struct {
	Username string `bson:"Username"   json:"Username"`
	Secret   string `bson:"secret"     json:"secret"`
	Email    string `bson:"email"      json:"email"`
}

func (dao *VerifyEmailDAO) CreateVerifyEmail(
	ctx context.Context, params CreateVerifyEmailParams,
) (model.VerifyEmail, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return model.VerifyEmail{}, err
	}

	user := model.VerifyEmail{
		ID:        id.String(),
		Username:  params.Username,
		Secret:    params.Secret,
		Email:     params.Email,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(15 * time.Minute),
	}
	_, err = dao.coll.InsertOne(ctx, user)
	if err != nil {
		return model.VerifyEmail{}, err
	}
	return user, nil
}
