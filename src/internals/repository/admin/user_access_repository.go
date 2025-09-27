package admin_repo

import (
	"context"
	"errors"

	"github.com/Chethu16/Chethu/src/internals/domain/admin"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct{
		collection *mongo.Collection
}

func NewUser(db *mongo.Database)*UserRepository{
	return &UserRepository{
		collection: db.Collection("users"),
	}
}
func(repo *UserRepository)CreateUser(ctx context.Context, user admin.User)error{
	_, err := repo.collection.InsertOne(ctx,user)
	if err != nil{
		return errors.New("failed to cerate user")
	}
	return nil
}

