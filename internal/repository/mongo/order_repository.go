package mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"ppharma/backend/internal/domain/order"
	mongowrap "ppharma/backend/internal/infra/db/mongo"
)

const ordersCollection = "orders"

type OrderRepository struct {
	collection mongowrap.Collection
}

var _ order.Repository = (*OrderRepository)(nil)

func NewOrderRepository(db mongowrap.Database) *OrderRepository {
	return &OrderRepository{collection: db.Collection(ordersCollection)}
}

func (r *OrderRepository) GetByID(orderID string) (*order.Order, error) {
	var ord order.Order
	err := r.collection.GetOne(context.Background(), bson.M{"order_id": orderID}, &ord)
	if errors.Is(err, mongowrap.ErrNotFound) {
		return nil, order.ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	return &ord, nil
}

func (r *OrderRepository) GetByIDForCustomer(orderID, customerID string) (*order.Order, error) {
	var ord order.Order
	err := r.collection.GetOne(context.Background(), bson.M{"order_id": orderID, "customer_id": customerID}, &ord)
	if errors.Is(err, mongowrap.ErrNotFound) {
		return nil, order.ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	return &ord, nil
}

func (r *OrderRepository) Save(ord *order.Order) error {
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"order_id": ord.OrderID},
		bson.M{"$set": ord},
		mongowrap.UpdateOptions{Upsert: true},
	)
	return err
}
