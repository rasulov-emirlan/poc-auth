package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/rasulov-emirlan/poc-auth/internal/domains/auth"
	"github.com/rasulov-emirlan/poc-auth/internal/entities"
)

type UsersRepository struct {
	conn *mongo.Client
}

func (r UsersRepository) Create(ctx context.Context, user entities.User) (entities.User, error) {
	res, err := r.conn.Database("poc-auth").Collection("users").InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return entities.User{}, auth.ErrEmailTaken
		}
		return entities.User{}, err
	}

	newID, ok := res.InsertedID.(string)
	if !ok {
		return entities.User{}, errors.New("failed to convert inserted id to string")
	}

	user.ID = newID

	return user, nil
}

func (r UsersRepository) GetByEmail(ctx context.Context, email string) (entities.User, error) {
	var user entities.User

	err := r.conn.Database("poc-auth").Collection("users").FindOne(ctx, entities.User{Email: email}).Decode(&user)
	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func (r UsersRepository) Update(ctx context.Context, user entities.User) (entities.User, error) {
	_, err := r.conn.Database("poc-auth").Collection("users").UpdateOne(ctx, entities.User{Email: user.Email}, user)
	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}
