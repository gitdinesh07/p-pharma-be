package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := ensureCustomerIndexes(ctx, db); err != nil {
		return err
	}
	if err := ensureUserIndexes(ctx, db); err != nil {
		return err
	}
	if err := ensureOrderIndexes(ctx, db); err != nil {
		return err
	}
	if err := ensurePaymentIndexes(ctx, db); err != nil {
		return err
	}
	if err := ensureProductIndexes(ctx, db); err != nil {
		return err
	}
	if err := ensureSessionIndexes(ctx, db); err != nil {
		return err
	}
	if err := ensureConsultationIndexes(ctx, db); err != nil {
		return err
	}
	return nil
}

func ensureCustomerIndexes(ctx context.Context, db *mongo.Database) error {
	_, err := db.Collection("customers").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true).SetSparse(true)},
		{Keys: bson.D{{Key: "mobile", Value: 1}}, Options: options.Index().SetUnique(true).SetSparse(true)},
	})
	return err
}

func ensureUserIndexes(ctx context.Context, db *mongo.Database) error {
	_, err := db.Collection("users").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true).SetSparse(true)},
		{Keys: bson.D{{Key: "mobile", Value: 1}}, Options: options.Index().SetUnique(true).SetSparse(true)},
	})
	return err
}

func ensureOrderIndexes(ctx context.Context, db *mongo.Database) error {
	_, err := db.Collection("orders").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "customer_id", Value: 1}, {Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "items.item_id", Value: 1}}},
		{Keys: bson.D{{Key: "items.status", Value: 1}}},
		{Keys: bson.D{{Key: "derived_status", Value: 1}}},
	})
	return err
}

func ensurePaymentIndexes(ctx context.Context, db *mongo.Database) error {
	_, err := db.Collection("payments").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "order_id", Value: 1}, {Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
	})
	return err
}

func ensureProductIndexes(ctx context.Context, db *mongo.Database) error {
	_, err := db.Collection("products").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "sku", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "is_active", Value: 1}, {Key: "category", Value: 1}}},
	})
	return err
}

func ensureSessionIndexes(ctx context.Context, db *mongo.Database) error {
	expires := int32(0)
	_, err := db.Collection("sessions").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_type", Value: 1}, {Key: "principal_id", Value: 1}, {Key: "device_id", Value: 1}}},
		{Keys: bson.D{{Key: "refresh_token_hash", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "expires_at", Value: 1}}, Options: options.Index().SetExpireAfterSeconds(expires)},
	})
	return err
}

func ensureConsultationIndexes(ctx context.Context, db *mongo.Database) error {
	_, err := db.Collection("consultations").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "customer_id", Value: 1}, {Key: "scheduled_at", Value: 1}}},
		{Keys: bson.D{{Key: "meeting.provider_event_id", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
	})
	return err
}
