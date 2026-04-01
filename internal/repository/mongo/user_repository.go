package mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"ppharma/backend/internal/domain/user"
	mongowrap "ppharma/backend/support-pkg/db/mongo"
)

const usersCollection = "users"

type UserRepository struct {
	collection mongowrap.Collection
}

var _ user.Repository = (*UserRepository)(nil)

func NewUserRepository(db mongowrap.Database) *UserRepository {
	return &UserRepository{collection: db.Collection(usersCollection)}
}

func (r *UserRepository) Create(u *user.User) error {
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": u.ID},
		bson.M{"$set": u},
		mongowrap.UpdateOptions{Upsert: true},
	)
	return err
}

func (r *UserRepository) GetByID(id string) (*user.User, error) {
	var u user.User
	err := r.collection.GetOne(context.Background(), bson.M{"_id": id}, &u)
	if errors.Is(err, mongowrap.ErrNotFound) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByEmail(email string) (*user.User, error) {
	var u user.User
	err := r.collection.GetOne(context.Background(), bson.M{"email": email}, &u)
	if errors.Is(err, mongowrap.ErrNotFound) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByMobile(mobile string) (*user.User, error) {
	var u user.User
	err := r.collection.GetOne(context.Background(), bson.M{"mobile": mobile}, &u)
	if errors.Is(err, mongowrap.ErrNotFound) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Update(u *user.User) error {
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": u.ID},
		bson.M{"$set": u},
	)
	return err
}
