package mongo

import (
	"context"
	"errors"

	"ppharma/backend/internal/domain/customer"
	mongowrap "ppharma/backend/support-pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

const customersCollection = "customers"

type CustomerRepository struct {
	collection mongowrap.Collection
}

var _ customer.Repository = (*CustomerRepository)(nil)

func NewCustomerRepository(db mongowrap.Database) *CustomerRepository {
	return &CustomerRepository{collection: db.Collection(customersCollection)}
}

func (r *CustomerRepository) Create(cust *customer.Customer) error {
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": cust.ID},
		bson.M{"$set": cust},
		mongowrap.UpdateOptions{Upsert: true},
	)
	return err
}

func (r *CustomerRepository) GetByID(id string) (*customer.Customer, error) {
	var cust customer.Customer
	err := r.collection.GetOne(context.Background(), bson.M{"_id": id}, &cust)
	if errors.Is(err, mongowrap.ErrNotFound) {
		return nil, customer.ErrCustomerNotFound
	}
	if err != nil {
		return nil, err
	}
	return &cust, nil
}

func (r *CustomerRepository) GetByEmail(email string) (*customer.Customer, error) {
	var cust customer.Customer
	err := r.collection.GetOne(context.Background(), bson.M{"email": email}, &cust)
	if errors.Is(err, mongowrap.ErrNotFound) {
		return nil, customer.ErrCustomerNotFound
	}
	if err != nil {
		return nil, err
	}
	return &cust, nil
}

func (r *CustomerRepository) GetByMobile(mobile string) (*customer.Customer, error) {
	var cust customer.Customer
	err := r.collection.GetOne(context.Background(), bson.M{"mobile": mobile}, &cust)
	if errors.Is(err, mongowrap.ErrNotFound) {
		return nil, customer.ErrCustomerNotFound
	}
	if err != nil {
		return nil, err
	}
	return &cust, nil
}

func (r *CustomerRepository) Update(cust *customer.Customer) error {
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": cust.ID},
		bson.M{"$set": cust},
		mongowrap.UpdateOptions{Upsert: false},
	)
	return err
}
