package mongo

import (
	"context"
	"errors"
	"time"

	gomongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultOperationTimeout = 10 * time.Second

var ErrNotFound = errors.New("document not found")

type Collection interface {
	Create(ctx context.Context, document any) error
	Get(ctx context.Context, filter any, result any) error
	GetOne(ctx context.Context, filter any, result any) error
	GetAll(ctx context.Context, filter any, results any, opts ...FindOptions) error
	UpdateOne(ctx context.Context, filter any, update any, opts ...UpdateOptions) (int64, error)
	DeleteOne(ctx context.Context, filter any) (int64, error)
	CountDocuments(ctx context.Context, filter any) (int64, error)
}

type Database interface {
	Collection(name string) Collection
}

type FindOptions struct {
	Sort       any
	Projection any
	Limit      int64
	Skip       int64
}

type UpdateOptions struct {
	Upsert bool
}

type ClientOption func(*Client)

type Client struct {
	db               *gomongo.Database
	operationTimeout time.Duration
}

func NewClient(db *gomongo.Database, opts ...ClientOption) *Client {
	c := &Client{db: db, operationTimeout: defaultOperationTimeout}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithOperationTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		if timeout > 0 {
			c.operationTimeout = timeout
		}
	}
}

func (c *Client) Collection(name string) Collection {
	return &collection{raw: c.db.Collection(name), operationTimeout: c.operationTimeout}
}

type collection struct {
	raw              *gomongo.Collection
	operationTimeout time.Duration
}

func (c *collection) Create(ctx context.Context, document any) error {
	opCtx, cancel := c.withTimeout(ctx)
	defer cancel()

	_, err := c.raw.InsertOne(opCtx, document)
	return err
}

func (c *collection) Get(ctx context.Context, filter any, result any) error {
	return c.GetOne(ctx, filter, result)
}

func (c *collection) GetOne(ctx context.Context, filter any, result any) error {
	opCtx, cancel := c.withTimeout(ctx)
	defer cancel()

	err := c.raw.FindOne(opCtx, filter).Decode(result)
	if errors.Is(err, gomongo.ErrNoDocuments) {
		return ErrNotFound
	}
	return err
}

func (c *collection) GetAll(ctx context.Context, filter any, results any, opts ...FindOptions) error {
	opCtx, cancel := c.withTimeout(ctx)
	defer cancel()

	findOptions := options.Find()
	if len(opts) > 0 {
		applied := opts[0]
		if applied.Sort != nil {
			findOptions.SetSort(applied.Sort)
		}
		if applied.Projection != nil {
			findOptions.SetProjection(applied.Projection)
		}
		if applied.Limit > 0 {
			findOptions.SetLimit(applied.Limit)
		}
		if applied.Skip > 0 {
			findOptions.SetSkip(applied.Skip)
		}
	}

	cursor, err := c.raw.Find(opCtx, filter, findOptions)
	if err != nil {
		return err
	}
	defer cursor.Close(opCtx)

	return cursor.All(opCtx, results)
}

func (c *collection) UpdateOne(ctx context.Context, filter any, update any, opts ...UpdateOptions) (int64, error) {
	opCtx, cancel := c.withTimeout(ctx)
	defer cancel()

	updateOptions := options.Update()
	if len(opts) > 0 {
		updateOptions.SetUpsert(opts[0].Upsert)
	}

	result, err := c.raw.UpdateOne(opCtx, filter, update, updateOptions)
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}

func (c *collection) DeleteOne(ctx context.Context, filter any) (int64, error) {
	opCtx, cancel := c.withTimeout(ctx)
	defer cancel()

	result, err := c.raw.DeleteOne(opCtx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func (c *collection) CountDocuments(ctx context.Context, filter any) (int64, error) {
	opCtx, cancel := c.withTimeout(ctx)
	defer cancel()

	return c.raw.CountDocuments(opCtx, filter)
}

func (c *collection) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithTimeout(ctx, c.operationTimeout)
}
